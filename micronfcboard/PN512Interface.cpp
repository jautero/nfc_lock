#include "PN512Interface.h"

PN512Interface::PN512Interface(SPI &spi, PinName ss):_ss(ss) 
{
    _spi = &spi;
    _ss = 1;
    _spi->format(8, 0);
    _spi->frequency(1000000);
}

PN512Interface::~PN512Interface(){
    _spi = NULL;
}

int PN512Interface::read(int reg) {
    int ret;
    _ss=0;
    _spi->write(regToSPIAddress(reg,true));
    ret=_spi->write(0);
    _ss=1;
    return ret;
}

void PN512Interface::write(int reg,int value) {
    _ss=0;
    _spi->write(regToSPIAddress(reg,false));
    _spi->write(value);
    _ss=1;
}

int PN512Interface::regToSPIAddress(int reg, bool read) {
    return ((reg & 0x3f) << 1) | ((read ? 1 : 0) << 7);
}