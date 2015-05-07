#include "mbed.h"
#include "USBSerial.h"

static USBSerial test_output;

extern "C" void test_output_start() {
    test_output.getc();
}

extern "C" int test_output_char(int ch) {
    return test_output.putc(ch);
}