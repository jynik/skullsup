// SPDX License Identifier: GPL-3.0
#include <Arduino.h>
#include <SoftSerial.h>
#include "Adafruit_NeoPixel.h"
#include "hw_cfg.h"
#include "version.h"

#define CMD_SUMMON      0xff    // Summon device for further commands.
                                // This used to transition from:
                                // (STATE_SLEEP|STATE_REANIMATED) -> STATE_IDLE
//                      0xfe       Reserved for future command
#define CMD_REANIMATE   0xfd    // Begin displaying frames
#define CMD_SET_COLOR   0xfc    // Reset and display a fixed color
#define CMD_FW_VERSION  0xfb    // Retrieve firmware version
#define CMD_NUM_STRIPS  0xfa    // Retrieve # of LED strips
#define CMD_STRIP_LEN   0xf9    // Retrieve # of LEDs per strip
#define CMD_LAYOUT      0xf8    // Retrieve physical LED layout
#define CMD_MAX_FRAMES  0xf7    // Retrieve MAX_FRAMES value

//  0xf6 - 0x80 are reserved for future commands
#define CMD_RESV_END    0xf6
#define CMD_RESV_START  0x80

// Do not include a frame delay, just update LEDs. OR this with LED "ID"
#define NO_FRAME_DELAY  0x40

// Addresses all LEDs when loading a frame. 0x3e - 0x00 address single LEDs
#define ALL_LEDS        0x3f

static enum {
    STATE_SLEEP = 0,
    STATE_IDLE,
    STATE_REANIMATED,
} state;

static struct frame {
  uint8_t   led_id;
  uint8_t   r;
  uint8_t   g;
  uint8_t   b;
} frames[MAX_FRAMES];               // Frame buffer


#define DEFAULT_FRAME_DUR_MS 100

static uint8_t frame_idx;           // Current animation frame
static uint8_t frame_count;         // Total # of animation frames
static uint16_t frame_dur_ms;       // frame duration in ms

#define CMD_BUF_LEN 4
static uint8_t cmd_idx = 0;         // Current index into command buffer
static uint8_t cmd_buf[CMD_BUF_LEN];   // Command buffer

#define RESP_BUF_LEN 4
static uint8_t resp_buf[RESP_BUF_LEN];

// Handle to software UART
static SoftSerial uart(PIN_RX, PIN_TX);

// Handle to LED strip
static Adafruit_NeoPixel leds(LED_COUNT, PIN_LED, NEO_GRB | NEO_KHZ800);

// Magic summon command sequance
static const uint8_t summon_cmd[] = {
    CMD_SUMMON, '1', '3', '8',
};

// Track where we are in the wake_cmd
static uint8_t summon_idx = 0;

#define SUMMON_CMD_ACK 0x9b

void neopixel_set_all(uint8_t r, uint8_t g, uint8_t b, bool show)
{
    for (unsigned int i = 0; i < LED_COUNT; i++) {
        leds.setPixelColor(i, r, g, b);
    }

    if (show) {
        leds.show();
    }
}

// Check the UART for our magic sequence. Return's true when we've gotten it.
bool is_summoned() {
    uint8_t c;
    bool ret;

    if (!uart.available()) return false;

    c = uart.read();
    if (c == summon_cmd[summon_idx]) {
        summon_idx++;

        if (summon_idx == sizeof(summon_cmd)) {
            summon_idx = 0;
            return true;
        }
    } else {
		summon_idx = 0;
	}

    return false;
}

inline void clear_frames()
{
    frame_idx = 0;
    frame_count = 0;
    frame_dur_ms = DEFAULT_FRAME_DUR_MS;
}

void setup()
{
    leds.begin();
    clear_frames();

    neopixel_set_all(24, 24, 24, true);

    uart.begin(UART_BAUDRATE);
    state = STATE_SLEEP;
}

inline void show_frame(const struct frame *f)
{
    if ((f->led_id & ALL_LEDS) == ALL_LEDS) {
      neopixel_set_all(f->r, f->g, f->b, false);
    } else {
      // This function will drop invalid LED ID value
      leds.setPixelColor(f->led_id & ALL_LEDS, f->r, f->g, f->b);
    }

    if (!(f->led_id & NO_FRAME_DELAY)) {
        leds.show();
        delay(frame_dur_ms);
    }
}

inline void enter_idle_state()
{
    clear_frames();
    state = STATE_IDLE;
}

inline uint8_t process_cmd() {
    uint8_t i, ack = 0;
    uint8_t resp_len = 0;

    for (i = 0; i < CMD_BUF_LEN; i++) {
        ack += cmd_buf[i];
    }

    switch (cmd_buf[0]) {
        case CMD_SUMMON:
            // Nothing to do other than ACK.
            break;

        case CMD_SET_COLOR:
            neopixel_set_all(cmd_buf[1], cmd_buf[2], cmd_buf[3], true);
            break;

        case CMD_REANIMATE:
            frame_idx = 0;
            frame_dur_ms = ((uint16_t) cmd_buf[1] << 8) | cmd_buf[2];
            state = STATE_REANIMATED;
            break;

            // Send FW Version in little-endian byte order
        case CMD_FW_VERSION:
            resp_buf[0] = FW_VERSION & 0xff;
            resp_buf[1] = FW_VERSION >> 8;
            resp_len = 2;
            break;

        case CMD_NUM_STRIPS:
            resp_buf[0] = NUM_STRIPS;
            resp_len = 1;
            break;

        case CMD_STRIP_LEN:
            resp_buf[0] = LEDS_PER_STRIP;
            resp_len = 1;
            break;

        case CMD_LAYOUT:
            resp_buf[0] = LED_LAYOUT;
            resp_len = 1;
            break;

        case CMD_MAX_FRAMES:
            resp_buf[0] = MAX_FRAMES;
            resp_len = 1;
            break;

        default:
            // Load Frame. The cmd nibble specifies the LED(s) to target
            if (cmd_buf[0] < CMD_RESV_START) {
                if (frame_count < MAX_FRAMES) {
                    frames[frame_count].led_id  = cmd_buf[0];
                    frames[frame_count].r       = cmd_buf[1];
                    frames[frame_count].g       = cmd_buf[2];
                    frames[frame_count].b       = cmd_buf[3];
                    frame_count++;
                }
            }
    }

    cmd_idx = 0;
    uart.write(ack);

    if (resp_len != 0) {
        uart.write(resp_buf, resp_len);
        resp_len = 0;
    }
}

void loop()
{
    switch (state) {
        case STATE_SLEEP:
            // Wait until we're summoned. This is intended to provide a bit
            // of resilience to spurious data if the host is still booting
            // and configuring pins.
            if (is_summoned()) {
                enter_idle_state();
                uart.write(SUMMON_CMD_ACK);
            }
            break;

        case STATE_IDLE:
            // Look for command and process it
            while (uart.available()) {
                char c = uart.read();
                cmd_buf[cmd_idx++] = c;
                if (cmd_idx >= CMD_BUF_LEN) {
                    process_cmd();
                }
            }
            break;

        case STATE_REANIMATED:
            if (is_summoned()) {
                // We need to be in the idle state to accept commands.
                // Otherwise, we may miss an byte while interrupts are
                // disabled in the time-critical NeoPixel update routine.
                enter_idle_state();
                uart.write(SUMMON_CMD_ACK);
            } else {
                // Display the current frame
                show_frame(&frames[frame_idx]);
                if (++frame_idx >= frame_count) {
                    frame_idx = 0;
                }
            }
			break;
    }
}
