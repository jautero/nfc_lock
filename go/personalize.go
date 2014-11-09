package main

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "fmt"
    "encoding/hex"
    "encoding/binary"
    "github.com/fuzxxl/nfc/2.0/nfc"    
    "github.com/fuzxxl/freefare/0.3/freefare"
    "./keydiversification"
)

// TODO: move to a separate helper module
func string_to_aeskey(keydata_str string) (*freefare.DESFireKey, error) {
    keydata := new([16]byte)
    to_keydata, err := hex.DecodeString(keydata_str)
    if err != nil {
        key := freefare.NewDESFireAESKey(*keydata, 0)
        return key, err
    }
    copy(keydata[0:], to_keydata)
    key := freefare.NewDESFireAESKey(*keydata, 0)
    return key,nil
}

func bytes_to_aeskey(source []byte) (*freefare.DESFireKey) {
    keydata := new([16]byte)
    copy(keydata[0:], source)
    key := freefare.NewDESFireAESKey(*keydata, 0)
    return key
}

func string_to_byte(source string) (byte, error) {
    bytearray, err := hex.DecodeString(source)
    if err != nil {
        return 0x0, err
    }
    return bytearray[0], nil
}

func main() {
    keys_data, err := ioutil.ReadFile("keys.yaml")
    if err != nil {
        panic(err)
    }

    keymap := make(map[interface{}]interface{});
    err = yaml.Unmarshal([]byte(keys_data), &keymap);
    if err != nil {
        panic(err)
    }

    apps_data, err := ioutil.ReadFile("apps.yaml")
    if err != nil {
        panic(err)
    }

    appmap := make(map[interface{}]interface{});
    err = yaml.Unmarshal([]byte(apps_data), &appmap);
    if err != nil {
        panic(err)
    }

    // Application-id from config
    aidbytes, err := hex.DecodeString(appmap["hacklab_acl"].(map[interface{}]interface{})["aid"].(string))
    if err != nil {
        panic(err)
    }
    aidint, n := binary.Uvarint(aidbytes)
    if n <= 0 {
        panic(fmt.Sprintf("binary.Uvarint returned %d", n))
    }
    aid := freefare.NewDESFireAid(uint32(aidint))
    //fmt.Println(aid)
    // Needed for diversification
    sysid, err := hex.DecodeString(appmap["hacklab_acl"].(map[interface{}]interface{})["sysid"].(string))
    if err != nil {
        panic(err)
    }

    // Key id numbers from config
    uid_read_key_id, err := string_to_byte(appmap["hacklab_acl"].(map[interface{}]interface{})["uid_read_key_id"].(string))
    if err != nil {
        panic(err)
    }
    acl_read_key_id, err := string_to_byte(appmap["hacklab_acl"].(map[interface{}]interface{})["acl_read_key_id"].(string))
    if err != nil {
        panic(err)
    }
    acl_write_key_id, err := string_to_byte(appmap["hacklab_acl"].(map[interface{}]interface{})["acl_write_key_id"].(string))
    if err != nil {
        panic(err)
    }
    prov_key_id, err := string_to_byte(appmap["hacklab_acl"].(map[interface{}]interface{})["provisioning_key_id"].(string))
    if err != nil {
        panic(err)
    }

    // Defaul (null) key
    nullkeydata := new([8]byte)
    defaultkey := freefare.NewDESFireDESKey(*nullkeydata)

    // New card master key
    new_master_key, err := string_to_aeskey(keymap["card_master"].(string))
    if err != nil {
        panic(err)
    }
    //fmt.Println(new_master_key)

    // The static app key to read UID
    uid_read_key, err := string_to_aeskey(keymap["uid_read_key"].(string))
    if err != nil {
        panic(err)
    }
    //fmt.Println(uid_read_key)

    // Bases for the diversified keys    
    prov_key_base, err := hex.DecodeString(keymap["prov_master"].(string))
    if err != nil {
        panic(err)
    }
    acl_read_base, err := hex.DecodeString(keymap["acl_read_key"].(string))
    if err != nil {
        panic(err)
    }
    acl_write_base, err := hex.DecodeString(keymap["acl_write_key"].(string))
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

        fmt.Println("Authenticating");
        error = desfiretag.Authenticate(0,*defaultkey)
        if error != nil {
            fmt.Println("Failed, trying agin with new key")
            error = desfiretag.Authenticate(0,*new_master_key)
            if error != nil {
                panic(error)
            }
            fmt.Println("Changing key back to default")
            error = desfiretag.ChangeKey(0,  *defaultkey, *new_master_key);
            if error != nil {
                panic(error)
            }
            fmt.Println("Re-auth with default key")
            error = desfiretag.Authenticate(0,*defaultkey)
            if error != nil {
                panic(error)
            }
            fmt.Println("Formatting (to get a clean state)")
            error = desfiretag.FormatPICC()
            if error != nil {
                panic(error)
            }
            return
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

        // Calculate the diversified keys
        prov_key_bytes, err := keydiversification.AES128(prov_key_base, aidbytes, realuid, sysid)
        if err != nil {
            panic(err)
        }
        prov_key := bytes_to_aeskey(prov_key_bytes)
        acl_read_bytes, err := keydiversification.AES128(acl_read_base, aidbytes, realuid, sysid)
        if err != nil {
            panic(err)
        }
        acl_read_key := bytes_to_aeskey(acl_read_bytes)
        acl_write_bytes, err := keydiversification.AES128(acl_write_base, aidbytes, realuid, sysid)
        if err != nil {
            panic(err)
        }
        acl_write_key := bytes_to_aeskey(acl_write_bytes)


        fmt.Println("Changing default master key");
        error = desfiretag.ChangeKey(0, *new_master_key, *defaultkey);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        /**
         * This is not needed for creating the application and does not help when changing application keys
        fmt.Println("Re-auth with new key")
        error = desfiretag.Authenticate(0,*new_master_key)
        if error != nil {
            panic(error)
        }
         */

        fmt.Println("Creating application");
        // TODO:Figure out what the settings byte (now hardcoded to 0xFF as it was in libfreefare example code) actually does
        error = desfiretag.CreateApplication(aid, 0xFF, 6 | freefare.CryptoAES);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");


        fmt.Println("Selecting application");
        // TODO:Figure out what the settings byte (now hardcoded to 0xFF as it was in libfreefare exampkle code) actually does
        error = desfiretag.SelectApplication(aid);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        /**
         * Does not work
        fmt.Println("Re-auth with new master key")
        error = desfiretag.Authenticate(0,*new_master_key)
        if error != nil {
            panic(error)
        }
         */

        /**
         * Also does not work
        fmt.Println("Re-auth with default key")
        error = desfiretag.Authenticate(0,*defaultkey)
        if error != nil {
            panic(error)
        }
         */

        fmt.Println("Changing provisioning key");
        error = desfiretag.ChangeKey(prov_key_id, *prov_key, *defaultkey);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        fmt.Println("Re-auth with new provisioning key")
        error = desfiretag.Authenticate(prov_key_id,*prov_key)
        if error != nil {
            panic(error)
        }


        fmt.Println("Changing static UID reading key");
        error = desfiretag.ChangeKey(uid_read_key_id, *uid_read_key, *defaultkey);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        fmt.Println("Changing ACL reading key");
        error = desfiretag.ChangeKey(acl_read_key_id, *acl_read_key, *defaultkey);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        fmt.Println("Changing ACL writing key");
        error = desfiretag.ChangeKey(acl_write_key_id, *acl_write_key, *defaultkey);
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");


        fmt.Println("Creating ACL data file");
        error = desfiretag.CreateDataFile(0, freefare.Enciphered, freefare.MakeDESFireAccessRights(acl_read_key_id, acl_write_key_id, prov_key_id, prov_key_id), 8, false)
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");

        // Not sure if this is actually needed
        fmt.Println("Committing");
        error = desfiretag.CommitTransaction()
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");


        fmt.Println("Disconnecting");
        error = desfiretag.Disconnect()
        if error != nil {
            panic(error)
        }
        fmt.Println("Done");
    }

}