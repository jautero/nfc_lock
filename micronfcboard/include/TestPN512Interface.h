#ifndef TESTPN512INTERFACE_H
#define TESTPN512INTERFACE_H

#include "PN512Interface.h"

#ifndef RUN_TEST
class TestPN512Interface
{
public:
    TestPN512Interface(SPI &spi, PinName ss);
    virtual ~TestPN512Interface ();
    void testRegToSPIAddress();
    void testVersionReg(int reg, int version);
    void testFIFO();
    void testReset();

private:
    PN512Interface interfaceUnderTest;
};
#endif
#endif /* end of include guard: TESTPN512INTERFACE_H */
