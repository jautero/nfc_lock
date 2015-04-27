#ifndef PN512INTERFACE_H
#define PN512INTERFACE_H

#include "mbed.h"

class PN512Interface
{
public:
    PN512Interface (SPI &spi, PinName SS);
    virtual ~PN512Interface ();

private:
    SPI *_spi;
    DigitalOut _ss;
};

#endif /* end of include guard: PN512INTERFACE_H */
