#ifndef NFCLOCK_MASTERKEY_CONFIG_H
#define NFCLOCK_MASTERKEY_CONFIG_H
#include <stdint.h>
#include "uid_config.h"
#include "acl_config.h"

// Application Master Key (ID 0x0 after application has been selected)
extern const uint8_t nfclock_amk[16];
// Card Master Key (ID 0x0 before application has been selected [or after selecting application 0x0])
extern const uint8_t nfclock_cmk[16];

#endif