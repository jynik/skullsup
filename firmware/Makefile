# Artisanal Skulls Up! build script

ifeq ($(ARDUINO_DIR),)
	ARDUINO_DIR := $(HOME)/.arduino15
endif

# Commands
AS			:= avr-gcc -x assembler-with-cpp
CC			:= avr-gcc
CXX			:= avr-g++
LD			:= avr-gcc
OBJCOPY		:= avr-objcopy
SIZE		:= avr-size
MKDIR		:= mkdir -p

# Compiler flags
ARCHDEFS	:= -mmcu=attiny85
CPUDEFS		:= $(ARCHDEFS) -DF_CPU=16500000L -DARDUINO=10803 -DARDUINO_AVR_DIGISPARK -DARDUINO_ARCH_AVR

ASFLAGS		:= -g $(CPUDEFS)
CFLAGS		:= -g -Os -w -fno-exceptions -ffunction-sections -fdata-sections -MMD $(CPUDEFS) -Wall -Wextra
CXXFLAGS	:= $(CFLAGS)
LDFLAGS		:= $(ARCHDEFS) -Os -Wl,-gc-sections
HEXFLAGS	:= -O ihex -R .eeprom
SIZEFLAGS	:= --format=avr --mcu=attiny85

# Build artifacts
TARGET 		:= skullsup
TARGET_ELF	:= $(TARGET).elf
TARGET_HEX  := $(TARGET).hex
TARGETS		:= $(TARGET_ELF) $(TARGET_HEX)

# Directories
ifeq ($(DIGISTUMP_VER),)
	DIGISTUMP_VER := 1.6.7
endif

OBJDIR := obj

DIGISTUMP_DIR := $(ARDUINO_DIR)/packages/digistump/hardware/avr/$(DIGISTUMP_VER)
LAUNCHER_DIR := $(ARDUINO_DIR)/packages/digistump/tools/micronucleus/2.0a4

SKULLSUP_OBJDIR	:= $(OBJDIR)/skullsup

DIGICDC_OBJDIR := $(OBJDIR)/DigiCDC
DIGICDC_SRCDIR := $(DIGISTUMP_DIR)/libraries/DigisparkCDC

CORE_OBJDIR := $(OBJDIR)/core
CORE_SRCDIR := $(DIGISTUMP_DIR)/cores/tiny

OBJDIR_LIST := $(SKULLSUP_OBJDIR) $(DIGICDC_OBJDIR) $(CORE_OBJDIR)

# Core library

# Skulls Up! Application
SRC := $(wildcard src/*.cpp)
OBJ := $(addprefix $(SKULLSUP_OBJDIR)/,$(notdir $(SRC:.cpp=.o)))
INC :=  -I$(DIGISTUMP_DIR)/cores/tiny 			\
		-I$(DIGISTUMP_DIR)/variants/digispark 	\
		-I$(DIGICDC_SRCDIR)

# Digispark CDC library
DIGICDC_ASM_SRC := $(DIGICDC_SRCDIR)/usbdrvasm.S

DIGICDC_C_SRC := \
	$(DIGICDC_SRCDIR)/oddebug.c \
	$(DIGICDC_SRCDIR)/osccal.c \
	$(DIGICDC_SRCDIR)/usbdrv.c

DIGICDC_CPP_SRC := $(DIGICDC_SRCDIR)/DigiCDC.cpp

DIGICDC_OBJ := \
	$(addprefix $(DIGICDC_OBJDIR)/,				\
		$(notdir $(DIGICDC_ASM_SRC:.S=.o)) 		\
		$(notdir $(DIGICDC_C_SRC:.c=.o)) 		\
		$(notdir $(DIGICDC_CPP_SRC:.cpp=.o)))

DIGICDC_INC :=  -I$(DIGISTUMP_DIR)/cores/tiny 			\
				-I$(DIGISTUMP_DIR)/variants/digispark 	\
				-I$(DIGICDC_SRCDIR)

# Arduino core library
CORE_C_SRC := \
	$(CORE_SRCDIR)/pins_arduino.c \
	$(CORE_SRCDIR)/wiring.c \
	$(CORE_SRCDIR)/wiring_digital.c \

CORE_CPP_SRC := \
	$(CORE_SRCDIR)/Print.cpp \
	$(CORE_SRCDIR)/main.cpp \
	$(CORE_SRCDIR)/new.cpp \

CORE_OBJ := \
	$(addprefix $(CORE_OBJDIR)/,	\
		$(notdir $(CORE_C_SRC:.c=.o) \
		$(notdir $(CORE_CPP_SRC:.cpp=.o))))


CORE_INC := -I$(DIGISTUMP_DIR)/cores/tiny 			\
			-I$(DIGISTUMP_DIR)/variants/digispark 	\

all: $(TARGET_HEX) size

size: $(TARGET_ELF)
	@echo -e '\n'
	$(SIZE) $(SIZEFLAGS) $<

install: $(TARGET_HEX)
	$(LAUNCHER_DIR)/launcher -cdigispark --timeout 30 -Uflash:w:$<:i


$(TARGET_HEX): $(TARGET_ELF)
	$(OBJCOPY) $(HEXFLAGS) $< $@

$(TARGET_ELF): $(OBJ) $(DIGICDC_OBJ) $(CORE_OBJ)
	$(LD) $(LDFLAGS) $^ -o $@

$(OBJDIR_LIST):
	$(MKDIR) $@

# Application code
$(SKULLSUP_OBJDIR)/%.o: src/%.cpp $(SKULLSUP_OBJDIR)
	$(CXX) $(CXXFLAGS) -DFIXED_PIXEL_COUNT=10 $(INC) -c $< -o $@

# Digspark CDC library
$(DIGICDC_OBJDIR)/%.o: $(DIGICDC_SRCDIR)/%.S $(DIGICDC_OBJDIR)
	$(AS) $(ASFLAGS) $(DIGICDC_INC) -c $< -o $@

$(DIGICDC_OBJDIR)/%.o: $(DIGICDC_SRCDIR)/%.c $(DIGICDC_OBJDIR)
	$(CC) $(CFLAGS) $(DIGICDC_INC) -c $< -o $@

$(DIGICDC_OBJDIR)/%.o: $(DIGICDC_SRCDIR)/%.cpp $(DIGICDC_OBJDIR)
	$(CXX) $(CXXFLAGS) $(DIGICDC_INC) -c $< -o $@

# Core library
$(CORE_OBJDIR)/%.o: $(CORE_SRCDIR)/%.c $(CORE_OBJDIR)
	$(CC) $(CFLAGS) $(CORE_INC) -c $< -o $@

$(CORE_OBJDIR)/%.o: $(CORE_SRCDIR)/%.cpp $(CORE_OBJDIR)
	$(CXX) $(CXXFLAGS) $(CORE_INC) -c $< -o $@

clean:
	rm -f $(TARGET_ELF) $(TARGET_HEX)
	rm -rf $(OBJDIR)

.PHONY: clean
