#ifndef PN512INTERFACE_H
#define PN512INTERFACE_H

#ifndef RUN_TEST
#include "mbed.h"


class PN512Interface
{
public:
    PN512Interface (SPI &spi, PinName SS);
    virtual ~PN512Interface();
    int read(int reg);
    void write(int reg,int value);

private:
    SPI *_spi;
    DigitalOut _ss;
    
    int regToSPIAddress(int reg, bool read);
    
    friend class TestPN512Interface;
};
#endif /* RUN_TEST */

#endif /* end of include guard: PN512INTERFACE_H */
