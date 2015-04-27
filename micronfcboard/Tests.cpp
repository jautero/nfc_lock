#include "Tests.h"
extern "C" {
    #include "unity.h"
}
#include "PN512Interface.h"

void testLED() {
    DigitalOut led1(LED1);
    DigitalOut led2(LED2);
    led1=1;
    TEST_ASSERT_EQUAL_INT(1,led1);
    wait(0.5);
    led1=0;
    TEST_ASSERT_EQUAL_INT(0,led1);
    led2=1;
    TEST_ASSERT_EQUAL_INT(1,led2);
    wait(0.5);
    led2=0;
    TEST_ASSERT_EQUAL_INT(0,led2);
}

const int PN512VersionReg=0x37;
const int PN512Version=0x82; 
    
void testSPIVersion() 
{
    DigitalOut cs(P1_19);
    SPI spi(P0_21,P0_22,P1_15);
    int ret;
    cs=1;
    spi.format(8, 0);
    spi.frequency(1000000);
    cs=0;
    spi.write(1<<7|PN512VersionReg<<1);
    ret=spi.write(0);
    cs=1;
    TEST_ASSERT_EQUAL_INT(PN512Version,ret);
}

void testPN512Interface()
{
    PN512Interface iface;
}