#include "PN512Interface.h"

PN512Interface::PN512Interface(SPI &spi, PinName ss):_ss(ss) 
{
    _spi = &spi;
}

PN512Interface::~PN512Interface(){
    
}