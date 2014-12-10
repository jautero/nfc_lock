package main

import (
    "fmt"
    "os"
    "strconv"
    "encoding/hex"
    "encoding/binary"
    "github.com/fuzxxl/nfc/2.0/nfc"    
    "github.com/fuzxxl/freefare/0.3/freefare"
    "./keydiversification"
    "./helpers"
)



func main() {
    if (len(os.Args) < 3) {
        fmt.Println("Usage:")
        fmt.Println(fmt.Sprintf("  %s member-id-as-decimal acl-value-as-hex", os.Args[0]))
        os.Exit(1)
    }

    // This bugs out at least for 1001-1004 (results in 2028-2025)
    newmid64, error := strconv.ParseUint(os.Args[1], 10, 16)
    if (error != nil) {
         panic(error)
    }
    newmid := uint16(newmid64)

    newacl64, error := strconv.ParseUint(os.Args[2], 16, 64)
    if (error != nil) {
         panic(error)
    }
    newacl := uint32(newacl64)

    newaclbytes := make([]byte, 4)
    n := binary.PutUvarint(newaclbytes, uint64(newacl))
    if (n < 0) {
        panic(fmt.Sprintf("binary.PutUvarint returned %d", n))
    }
    newmidbytes := make([]byte, 2)
    n = binary.PutUvarint(newmidbytes, uint64(newmid))
    if (n < 0) {
        panic(fmt.Sprintf("binary.PutUvarint returned %d", n))
    }


    keymap, err := helpers.LoadYAMLFile("keys.yaml")
    if err != nil {
        panic(err)
    }

    appmap, err := helpers.LoadYAMLFile("apps.yaml")
    if err != nil {
        panic(err)
    }

    // Application-id
    aid, err := helpers.String2aid(appmap["hacklab_acl"].(map[interface{}]interface{})["aid"].(string))
    if err != nil {
        panic(err)
    }

    // Needed for diversification
    aidbytes := helpers.Aid2bytes(aid)
    sysid, err := hex.DecodeString(appmap["hacklab_acl"].(map[interface{}]interface{})["sysid"].(string))
    if err != nil {
        panic(err)
    }

    // Key id numbers from config
    uid_read_key_id, err := helpers.String2byte(appmap["hacklab_acl"].(map[interface{}]interface{})["uid_read_key_id"].(string))
    if err != nil {
        panic(err)
    }
    acl_read_key_id, err := helpers.String2byte(appmap["hacklab_acl"].(map[interface{}]interface{})["acl_read_key_id"].(string))
    if err != nil {
        panic(err)
    }
    prov_key_id, err := helpers.String2byte(appmap["hacklab_acl"].(map[interface{}]interface{})["provisioning_key_id"].(string))
    if err != nil {
        panic(err)
    }
    acl_file_id, err := helpers.String2byte(appmap["hacklab_acl"].(map[interface{}]interface{})["acl_file_id"].(string))
    if err != nil {
        panic(err)
    }
    mid_file_id, err := helpers.String2byte(appmap["hacklab_acl"].(map[interface{}]interface{})["mid_file_id"].(string))
    if err != nil {
        panic(err)
    }

    // The static app key to read UID
    uid_read_key, err := helpers.String2aeskey(keymap["uid_read_key"].(string))
    if err != nil {
        panic(err)
    }
    //fmt.Println(uid_read_key)

    // Bases for the diversified keys    
    acl_read_base, err := hex.DecodeString(keymap["acl_read_key"].(string))
    if err != nil {
        panic(err)
    }
    prov_key_base, err := hex.DecodeString(keymap["prov_master"].(string))
    if err != nil {
        panic(err)
    }

    // Open device and get tags list
    d, err := nfc.Open("");
    if err != nil {
        panic(err)
    }

    tags, err := freefare.GetTags(d);
    if err != nil {
        panic(err)
    }

    // Initialize each tag with our app
    for i := 0; i < len(tags); i++ {
        tag := tags[i]
        fmt.Println(tag.String(), tag.UID())

        // Skip non desfire tags
        if (tag.Type() != freefare.DESFire) {
            fmt.Println("Skipped");
            continue
        }
        
        desfiretag := tag.(freefare.DESFireTag)

        // Connect to this tag
        fmt.Println("Connecting");
        error := desfiretag.Connect()
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        fmt.Println("Selecting application");
        error = desfiretag.SelectApplication(aid);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        fmt.Println("Authenticating");
        error = desfiretag.Authenticate(uid_read_key_id,*uid_read_key)
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        // Get card real UID        
        realuid_str, error := desfiretag.CardUID()
        if error != nil {
            panic(error)
        }
        realuid, error := hex.DecodeString(realuid_str);
        if error != nil {
            panic(error)
        }
        fmt.Println("realuid:", hex.EncodeToString(realuid));

        // Calculate the diversified keys
        acl_read_bytes, err := keydiversification.AES128(acl_read_base, aidbytes, realuid, sysid)
        if err != nil {
            panic(err)
        }
        acl_read_key := helpers.Bytes2aeskey(acl_read_bytes)

        prov_key_bytes, err := keydiversification.AES128(prov_key_base, aidbytes, realuid, sysid)
        if err != nil {
            panic(err)
        }
        prov_key := helpers.Bytes2aeskey(prov_key_bytes)

        fmt.Println("Re-auth with ACL read key")
        error = desfiretag.Authenticate(acl_read_key_id,*acl_read_key)
        if error != nil {
            panic(error)
        }

        aclbytes := make([]byte, 4)
        fmt.Println("Reading ACL data file")
        bytesread, err := desfiretag.ReadData(acl_file_id, 0, aclbytes)
        if error != nil {
            panic(error)
        }
        if (bytesread <= 0) {
            fmt.Println(fmt.Sprintf("ReadData read %d bytes", bytesread))
            continue
        }
        if (bytesread < 4) {
            //panic(fmt.Sprintf("ReadData read %d bytes, 4 expected", bytesread))
            fmt.Println(fmt.Sprintf("ReadData read %d bytes, 4 expected", bytesread))
        }
        acl64, n := binary.Uvarint(aclbytes)
        if n <= 0 {
            panic(fmt.Sprintf("binary.Uvarint returned %d", n))
        }
        acl := uint32(acl64)
        fmt.Println("acl:", acl)

        midbytes := make([]byte, 2)
        fmt.Println("Reading member-id data file")
        bytesread, err = desfiretag.ReadData(mid_file_id, 0, midbytes)
        if error != nil {
            panic(error)
        }
        // TODO: make a helper that repeatedly reads until we're past the buffer or readbytes == 0
        if (bytesread < 2) {
            //panic(fmt.Sprintf("ReadData read %d bytes, 2 expected", bytesread))
            fmt.Println(fmt.Sprintf("ReadData read %d bytes, 2 expected", bytesread))
        }
        mid64, n := binary.Uvarint(midbytes)
        if n <= 0 {
            panic(fmt.Sprintf("binary.Uvarint returned %d", n))
        }
        mid := uint16(mid64)
        fmt.Println("mid:", mid)

        if (mid == newmid && acl == newacl) {
            // Do nothing
            fmt.Println("No need to update ACL and member-id")
        } else {

            fmt.Println("Re-auth with provisioning key")
            error = desfiretag.Authenticate(prov_key_id,*prov_key)
            if error != nil {
                panic(error)
            }
            fmt.Println("Done");
            
            fmt.Println("Writing ACL file")
            bytewritten, error := desfiretag.WriteData(acl_file_id, 0, newaclbytes)
            if error != nil {
                panic(error)
            }
            if (bytewritten < 4) {
                //panic(fmt.Sprintf("WriteData wrote %d bytes, 4 expected", bytewritten))
                fmt.Println(fmt.Sprintf("WriteData wrote %d bytes, 4 expected", bytewritten))
            }
            fmt.Println("Done");

            fmt.Println("Writing member-id file")
            bytewritten, error = desfiretag.WriteData(mid_file_id, 0, newmidbytes)
            if error != nil {
                panic(error)
            }
            if (bytewritten < 2) {
                //panic(fmt.Sprintf("WriteData wrote %d bytes, 2 expected", bytewritten))
                fmt.Println(fmt.Sprintf("WriteData wrote %d bytes, 2 expected", bytewritten))
            }
            fmt.Println("Done");

            /**
             * For some reason this gives the dreaded "unknown error"
            fmt.Println("Committing");
            error = desfiretag.CommitTransaction()
            if error != nil {
                panic(error)
            }
            fmt.Println("Done");
             */
        }



        fmt.Println("Disconnecting");
        error = desfiretag.Disconnect()
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        fmt.Println("*** BEGIN: SAVE THIS INFO ***");
        fmt.Println(fmt.Sprintf("  Member id: %d", newmid))
        fmt.Println(fmt.Sprintf("  Card UID: %s", realuid_str))
        fmt.Println("*** END: SAVE THIS INFO ***");

    }

}