#include "mbed.h"
#include "USBSerial.h"

DigitalOut myled(LED1);
USBSerial pc;

int main() {
    while(1) {
        myled = 1;
        wait(0.2);
        myled = 0;
        wait(0.2);
        pc.printf("HelloWorld!");
    }
}
