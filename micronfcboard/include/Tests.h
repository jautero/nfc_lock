#ifndef TESTS_H
#define TESTS_H 

#ifndef RUN_TEST
#include "mbed.h"
#endif

#  ifdef __cplusplus
extern  "C" {
#  endif                        // __cplusplus
    void testLED();
    void testSPIVersion();
    void testPN512Interface();

# ifdef __cplusplus
}
# endif


#endif