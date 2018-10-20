# Artisanal Skulls Up! build script

ifeq ($(ARDUINO_DIR),)
	ARDUINO_DIR := $(HOME)/.arduino15
endif

AVRGCC_VERSION := 4.8.1-arduino5
AVRGCC_DIR := $(ARDUINO_DIR)/packages/arduino/tools/avr-gcc/$(AVRGCC_VERSION)

CROSS_COMPILE := $(AVRGCC_DIR)/bin/avr-

ifeq ($(PLATFORM),)
	PLATFORM := default
endif

# Commands
AS			:= $(CROSS_COMPILE)gcc -x assembler-with-cpp
CC			:= $(CROSS_COMPILE)gcc
CXX			:= $(CROSS_COMPILE)g++
LD			:= $(CROSS_COMPILE)gcc
OBJCOPY		:= $(CROSS_COMPILE)objcopy
SIZE		:= $(CROSS_COMPILE)size
MKDIR		:= mkdir -p

# Compiler flags
ARCHDEFS	:= -mmcu=attiny85
CPUDEFS		:= $(ARCHDEFS) -DF_CPU=16500000L -DARDUINO=10803 -DARDUINO_AVR_DIGISPARK -DARDUINO_ARCH_AVR

ASFLAGS		:= -g $(CPUDEFS)
CFLAGS		:= -g -Os -w -fno-exceptions -ffunction-sections -fdata-sections -MMD $(CPUDEFS) -Wall -Wextra -DPLATFORM_$(PLATFORM)
CXXFLAGS	:= $(CFLAGS)
LDFLAGS		:= $(ARCHDEFS) -Os -Wl,-gc-sections
HEXFLAGS	:= -O ihex -R .eeprom
SIZEFLAGS	:= --format=avr --mcu=attiny85

# Build artifacts
TARGET 		:= skullsup
TARGET_ELF	:= $(TARGET)-$(PLATFORM).elf
TARGET_HEX  := $(TARGET)-$(PLATFORM).hex
TARGETS		:= $(TARGET_ELF) $(TARGET_HEX)

# Directories
ifeq ($(DIGISTUMP_VER),)
	DIGISTUMP_VER := 1.6.7
endif

OBJDIR := obj

DIGISTUMP_DIR := $(ARDUINO_DIR)/packages/digistump/hardware/avr/$(DIGISTUMP_VER)
LAUNCHER_DIR := $(ARDUINO_DIR)/packages/digistump/tools/micronucleus/2.0a4

SKULLSUP_OBJDIR	:= $(OBJDIR)/skullsup

DIGISPARK_LIB_OBJDIR := $(OBJDIR)/digispark

CORE_OBJDIR := $(OBJDIR)/core
CORE_SRCDIR := $(DIGISTUMP_DIR)/cores/tiny

OBJDIR_LIST := $(SKULLSUP_OBJDIR) $(DIGISPARK_LIB_OBJDIR) $(CORE_OBJDIR)

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
	$(addprefix $(CORE_OBJDIR)/, 					\
		$(notdir $(CORE_C_SRC:.c=.o) 				\
		$(notdir $(CORE_CPP_SRC:.cpp=.o))))

CORE_INC := -I$(DIGISTUMP_DIR)/cores/tiny 			\
			-I$(DIGISTUMP_DIR)/variants/digispark 	\

# Digispark Libraries
SOFTSERIAL_SRCDIR := \
	$(DIGISTUMP_DIR)/libraries/DigisparkSoftSerial

TINYPINCHANGE_SRCDIR := \
	$(DIGISTUMP_DIR)/libraries/DigisparkTinyPinChange

DIGISPARK_LIB_SRC := \
	$(SOFTSERIAL_SRCDIR)/SoftSerial.cpp \
	$(TINYPINCHANGE_SRCDIR)/TinyPinChange.cpp

DIGISPARK_LIB_OBJ := \
	$(addprefix $(DIGISPARK_LIB_OBJDIR)/,			\
		$(notdir $(DIGISPARK_LIB_SRC:.cpp=.o)))

DIGISPARK_INC := \
	-I$(DIGISTUMP_DIR)/libraries/DigisparkSoftSerial \
	-I$(DIGISTUMP_DIR)/libraries/DigisparkTinyPinChange

# Skulls Up! Application
SRC := $(wildcard src/*.cpp)
OBJ := $(addprefix $(SKULLSUP_OBJDIR)/,$(notdir $(SRC:.cpp=.o)))
INC := $(CORE_INC) $(DIGISPARK_INC)

all: $(TARGET_HEX) size

help:
	@echo "Usage: make PLATFORM=<platform> [clean]"
	@echo ""
	@echo "Available platforms: "
	@echo "  default"
	@echo ""

size: $(TARGET_ELF)
	@echo -e '\n'
	$(SIZE) $(SIZEFLAGS) $<

install: $(TARGET_HEX)
	$(LAUNCHER_DIR)/launcher -cdigispark --timeout 30 -Uflash:w:$<:i

$(TARGET_SYM): $(TARGET_HEX)
	@ln -s $(TARGET_HEX) $@


build-msg:
	@echo ---------------------------------------------------------------------
	@echo   Building for PLATFORM=$(PLATFORM)
	@echo ---------------------------------------------------------------------

$(TARGET_HEX): build-msg $(TARGET_ELF)
	$(OBJCOPY) $(HEXFLAGS) $(TARGET_ELF) $@

$(TARGET_ELF): $(OBJ) $(DIGISPARK_LIB_OBJ) $(CORE_OBJ)
	$(LD) $(LDFLAGS) $^ -o $@

$(OBJDIR_LIST):
	$(MKDIR) $@

# Application code
$(SKULLSUP_OBJDIR)/%.o: src/%.cpp $(SKULLSUP_OBJDIR)
	$(CXX) $(CXXFLAGS) $(INC) -c $< -o $@

# Digispark Libraries
$(DIGISPARK_LIB_OBJDIR)/%.o: $(SOFTSERIAL_SRCDIR)/%.cpp $(DIGISPARK_LIB_OBJDIR)
	$(CXX) $(CXXFLAGS) $(CORE_INC) $(DIGISPARK_INC) -c $< -o $@

$(DIGISPARK_LIB_OBJDIR)/%.o: $(TINYPINCHANGE_SRCDIR)/%.cpp $(DIGISPARK_LIB_OBJDIR)
	$(CXX) $(CXXFLAGS) $(CORE_INC) $(DIGISPARK_INC) -c $< -o $@

# Core library
$(CORE_OBJDIR)/%.o: $(CORE_SRCDIR)/%.c $(CORE_OBJDIR)
	$(CC) $(CFLAGS) $(CORE_INC) -c $< -o $@

$(CORE_OBJDIR)/%.o: $(CORE_SRCDIR)/%.cpp $(CORE_OBJDIR)
	$(CXX) $(CXXFLAGS) $(CORE_INC) -c $< -o $@

clean:
	rm -f *.elf *.hex
	rm -rf $(OBJDIR)

.PHONY: build-msg size help clean
