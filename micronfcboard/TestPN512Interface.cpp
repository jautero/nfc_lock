#include "TestPN512Interface.h"

extern "C" {
#include "unity.h"
}

TestPN512Interface::TestPN512Interface(SPI &spi, PinName ss):
    interfaceUnderTest(spi,ss) 
{
    
}
TestPN512Interface::~TestPN512Interface ()
{
    
}

void TestPN512Interface::testRegToSPIAddress()
{
    int reg=0;
    int addr=0;
    while (reg<0x3D) {
        TEST_ASSERT_EQUAL_INT(addr,interfaceUnderTest.regToSPIAddress(reg,false));
        TEST_ASSERT_EQUAL_INT(addr+0x80,interfaceUnderTest.regToSPIAddress(reg,true));
        reg++;
        addr++; addr++;
    };
}

void TestPN512Interface::testVersionReg(int reg, int version)
{
    TEST_ASSERT_EQUAL_INT(version,interfaceUnderTest.read(reg));
}

const int FIFODataReg=0x09;
const int FIFOLevelReg=0x10;
    
void TestPN512Interface::testFIFO()
{
    int level=interfaceUnderTest.read(FIFOLevelReg);
    interfaceUnderTest.write(FIFODataReg,0xbb);
    TEST_ASSERT_EQUAL_INT(level+1,interfaceUnderTest.read(FIFOLevelReg));
    TEST_ASSERT_EQUAL_INT(0xbb,interfaceUnderTest.read(FIFODataReg));
    TEST_ASSERT_EQUAL_INT(level,interfaceUnderTest.read(FIFOLevelReg));
}

void TestPN512Interface::testReset()
{
    interfaceUnderTest.write(0x01,0x0f);
    while (interfaceUnderTest.read(0x01)!=0x20);
    TEST_ASSERT_EQUAL_INT(0x3B,interfaceUnderTest.read(0x11));
}