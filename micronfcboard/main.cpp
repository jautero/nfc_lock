#include "mbed.h"
#include "USBKeyboard.h"

DigitalOut myled(LED1);
USBKeyboard pc;

int main() {
    while(1) {
        pc.printf("1337 H4X0R!\n");
    }
}
