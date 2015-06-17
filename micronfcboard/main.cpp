#include "mbed.h"

#include "USBSerial.h"

DigitalOut cs(P1_19);
SPI spi(P0_21,P0_22,P1_15);
USBSerial pc;

int main() {
    int ret,out;
    cs=1;
    spi.format(8, 0);
    spi.frequency(1000000);

    while(1) {
        out=pc.getc();
        cs=0;
        ret=spi.write(out);
        cs=1;
        pc.putc(ret);
    }
}
