# Makefile for micronfcboard
GCC_BIN = 
PROJECT = MICRONFCBOARD
OBJECTS = 
SYS_OBJECTS = 
INCLUDE_PATHS = 
LIBRARY_PATHS = 
LIBRARIES = 
LINKER_SCRIPT = 

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
CC_SYMBOLS = -DTARGET_LPC11U35_401 -DTARGET_M0 -DTARGET_CORTEX_M -DTARGET_NXP -DTARGET_LPC11UXX -DTOOLCHAIN_GCC_ARM -DTOOLCHAIN_GCC -D__CORTEX_M0 -DARM_MATH_CM0 -DMBED_BUILD_TIMESTAMP=1427405354.12 -D__MBED__=1 

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


$(PROJECT).elf: $(OBJECTS) $(SYS_OBJECTS)
	$(LD) $(LD_FLAGS) -T$(LINKER_SCRIPT) $(LIBRARY_PATHS) -o $@ $^ $(LIBRARIES) $(LD_SYS_LIBS) $(LIBRARIES) $(LD_SYS_LIBS)
	@echo ""
	@echo "*****"
	@echo "***** You must modify vector checksum value in *.bin and *.hex files."
	@echo "*****"
	@echo ""
	$(SIZE) $@

$(PROJECT).bin: $(PROJECT).elf
	@$(OBJCOPY) -O binary $< $@

$(PROJECT).hex: $(PROJECT).elf
	@$(OBJCOPY) -O ihex $< $@

$(PROJECT).lst: $(PROJECT).elf
	@$(OBJDUMP) -Sdh $< > $@

lst: $(PROJECT).lst

size:
	$(SIZE) $(PROJECT).elf

DEPS = $(OBJECTS:.o=.d) $(SYS_OBJECTS:.o=.d)
-include $(DEPS)
