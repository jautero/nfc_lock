# Makefile for micronfcboard
GCC_BIN = 
RUBY=ruby
UNITYDIR=$(VIRTUAL_ENV)/Unity
TARGET = TARGET_APPNEARME_MICRONFCBOARD
TOOLCHAIN = TOOLCHAIN_GCC_ARM
PROJECT = MICRONFCBOARD
OBJECTS = main.o
TDD_OBJECTS = $(UNITYDIR)/src/unity.o Tests.o Tests_runner.o
SYS_OBJECTS = mbed/$(TARGET)/$(TOOLCHAIN)/board.o mbed/$(TARGET)/$(TOOLCHAIN)/cmsis_nvic.o \
	mbed/$(TARGET)/$(TOOLCHAIN)/retarget.o mbed/$(TARGET)/$(TOOLCHAIN)/startup_LPC11xx.o \
	mbed/$(TARGET)/$(TOOLCHAIN)/system_LPC11Uxx.o
INCLUDE_PATHS = -Imbed -Imbed/$(TARGET)/ -Imbed/$(TARGET)/TARGET_NXP/ -Imbed/$(TARGET)/TARGET_NXP/TARGET_LPC11UXX/ \
	-Imbed/$(TARGET)/TARGET_NXP/TARGET_LPC11UXX/$(TARGET)/ -Iusb/USBDevice/ -Iusb/USBSerial/ -Iusb/USBHID/ \
	-I$(UNITYDIR)/src
	
LIBRARY_PATHS = -Lmbed/$(TARGET)/$(TOOLCHAIN) -Lusb/$(TARGET)/$(TOOLCHAIN)
LIBRARIES = -lmbed -lUSBDevice
LINKER_SCRIPT = mbed/$(TARGET)/$(TOOLCHAIN)/LPC11U34.ld

############################################################################### 
AS      = $(GCC_BIN)arm-none-eabi-as
CC      = $(GCC_BIN)arm-none-eabi-gcc
CPP     = $(GCC_BIN)arm-none-eabi-g++
LD      = $(GCC_BIN)arm-none-eabi-gcc
OBJCOPY = $(GCC_BIN)arm-none-eabi-objcopy
OBJDUMP = $(GCC_BIN)arm-none-eabi-objdump
SIZE 	= $(GCC_BIN)arm-none-eabi-size

CPU = -mcpu=cortex-m0 -mthumb
CC_FLAGS = $(CPU) -c -g -fno-common -fmessage-length=0 -Wall -fno-exceptions -ffunction-sections -fdata-sections -fomit-frame-pointer
CC_FLAGS += -MMD -MP
CC_SYMBOLS = -DTARGET_APPNEARME_MICRONFCBOARD -DTARGET_M0 -DTARGET_CORTEX_M -DTARGET_NXP -DTARGET_LPC11UXX -D$(TOOLCHAIN) -DTOOLCHAIN_GCC -D__CORTEX_M0 -DARM_MATH_CM0 -DMBED_BUILD_TIMESTAMP=1427405354.12 -D__MBED__=1 

LD_FLAGS = $(CPU) -Wl,--gc-sections --specs=nano.specs -Wl,--wrap,main
#LD_FLAGS += -u _printf_float -u _scanf_float
LD_FLAGS += -Wl,-Map=$(PROJECT).map,--cref
LD_SYS_LIBS = -lstdc++ -lsupc++ -lm -lc -lgcc -lnosys

ifeq ($(DEBUG), 1)
  CC_FLAGS += -DDEBUG -O0
else
  CC_FLAGS += -DNDEBUG -Os
endif

all: $(PROJECT).bin $(PROJECT).hex 

clean:
	rm -f $(PROJECT).bin $(PROJECT).elf $(PROJECT).hex $(PROJECT).map $(PROJECT).lst $(OBJECTS) $(DEPS)

.s.o:
	$(AS) $(CPU) -o $@ $<

.c.o:
	$(CC)  $(CC_FLAGS) $(CC_SYMBOLS) -std=gnu99   $(INCLUDE_PATHS) -o $@ $<

.cpp.o:
	$(CPP) $(CC_FLAGS) $(CC_SYMBOLS) -std=gnu++98 -fno-rtti $(INCLUDE_PATHS) -o $@ $<

%_runner.c: %.cpp
	$(RUBY) $(UNITYDIR)/auto/generate_test_runner.rb $< $@

$(PROJECT).elf: $(OBJECTS) $(SYS_OBJECTS)
	$(LD) $(LD_FLAGS) -T$(LINKER_SCRIPT) $(LIBRARY_PATHS) -o $@ $^ $(LIBRARIES) $(LD_SYS_LIBS) $(LIBRARIES) $(LD_SYS_LIBS)
	$(SIZE) $@

TDD.elf: $(TDD_OBJECTS) $(SYS_OBJECTS)
	$(LD) $(LD_FLAGS) -T$(LINKER_SCRIPT) $(LIBRARY_PATHS) -o $@ $^ $(LIBRARIES) $(LD_SYS_LIBS) $(LIBRARIES) $(LD_SYS_LIBS)
	$(SIZE) $@
	
%.bin: %.elf
	@$(OBJCOPY) -O binary $< $@
	@python ./patchtool.py $@

%.hex: %.elf
	@$(OBJCOPY) -O ihex $< $@

%.lst: %.elf
	@$(OBJDUMP) -Sdh $< > $@

lst: $(PROJECT).lst

size:
	$(SIZE) $(PROJECT).elf

DEPS = $(OBJECTS:.o=.d) $(SYS_OBJECTS:.o=.d)
-include $(DEPS)
