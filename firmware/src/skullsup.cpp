// SPDX License Identifier: GPL-3.0
#include <Arduino.h>
#include <DigiCDC.h>
#include "Adafruit_NeoPixel.h"
#include "hw_cfg.h"

#define CMD_WAKE        0xff    // Bring device out of initial state
#define CMD_RESET       0xfe    // Clear frame buffer, and set fixed color
#define CMD_REANIMATE   0xfd	// Begin displaying frames
#define CMD_SET_COLOR   0xfc	// Reset and display a fixed color
								// fb - 0x81 are reserved for future commands

// Do not include a frame delay, just update LEDs. OR this with LED "ID"
#define NO_FRAME_DELAY	0x40

// Addresses all LEDs when loading a frame. 0x3e - 0x00 address single LEDs
#define ALL_LEDS        0x3f

static enum {
    STATE_SLEEP = 0,
    STATE_IDLE,
    STATE_REANIMATED,
} state;

#define MAX_FRAMES 44
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

// Handle to LED strip
static Adafruit_NeoPixel leds =
    Adafruit_NeoPixel(LED_COUNT, LED_PIN, NEO_GRB | NEO_KHZ800);

// Wake command
static const uint8_t wake_cmd[CMD_BUF_LEN] = {
    CMD_WAKE, '1', '3', '8',
};

static void neopixel_set_all(uint8_t r, uint8_t g, uint8_t b, bool show)
{
    for (unsigned int i = 0; i < 10; i++) {
        leds.setPixelColor(i, r, g, b);
    }

    if (show) {
        SerialUSB.refresh();
        leds.show();
        SerialUSB.refresh();
    }
}

static inline void ack_byte(uint8_t c) {
    SerialUSB.write(~c);
}

/* It seems our RX buffer gets filled with some data as the USB ACM driver
 * attaches and (presumably?) configures the device. Wait for a magic wake_cmd
 * string before handling data. */
void sleep_until_summoned()
{
    uint8_t c;
    bool summoned = false;
    uint8_t match = 0;

    while (!summoned) {
        if (!SerialUSB.available()) continue;

        c = SerialUSB.read();
        if (c == wake_cmd[match]) {
            match++;
        } else {
            match = 0;
        }

        summoned = (match >= sizeof(wake_cmd));
        ack_byte(c);
    }
}

static inline void reset_frames()
{
    frame_idx = 0;
    frame_count = 0;
    frame_dur_ms = DEFAULT_FRAME_DUR_MS;
}

void setup()
{
    leds.begin();
    reset_frames();

    neopixel_set_all(16, 24, 8, true);

    SerialUSB.begin();

    state = STATE_SLEEP;
    sleep_until_summoned();
    state = STATE_IDLE;

    neopixel_set_all(0, 0,0,true);
}

static inline void show_frame(const struct frame *f)
{
    if ((f->led_id & ALL_LEDS) == ALL_LEDS) {
      neopixel_set_all(f->r, f->g, f->b, false);
    } else {
      // This function will drop invalid LED ID value
      leds.setPixelColor(f->led_id & ALL_LEDS, f->r, f->g, f->b);
    }

    SerialUSB.refresh();

    if (!(f->led_id & NO_FRAME_DELAY)) {
        leds.show();
        SerialUSB.delay(frame_dur_ms);
    }
}


void loop()
{
    char c;
    while (SerialUSB.available()) {
        c = SerialUSB.read();
        cmd_buf[cmd_idx++] = c;
        if (cmd_idx >= CMD_BUF_LEN) {
            switch (cmd_buf[0]) {
                case CMD_WAKE:
                    break;

                case CMD_SET_COLOR:
                    neopixel_set_all(cmd_buf[1], cmd_buf[2], cmd_buf[3], true);
                    // Fall-through

                case CMD_RESET:
                    state = STATE_IDLE;
                    reset_frames();
                    break;

                case CMD_REANIMATE:
                    state = STATE_REANIMATED;
                    frame_idx = 0;
                    frame_dur_ms = ((uint16_t) cmd_buf[1] << 8) | cmd_buf[2];
                    break;

                // Load Frame. The cmd nibble specifies the LED(s) to target
                default:
                    state = STATE_IDLE;
                    if (frame_count < MAX_FRAMES) {
                        frames[frame_count].led_id  = cmd_buf[0];
                        frames[frame_count].r       = cmd_buf[1];
                        frames[frame_count].g       = cmd_buf[2];
                        frames[frame_count].b       = cmd_buf[3];
                        frame_count++;
                    }
            }

            cmd_idx = 0;
        }

        ack_byte(c);
    }

    if (state == STATE_REANIMATED) {
        show_frame(&frames[frame_idx]);
        if (++frame_idx >= frame_count) {
            frame_idx = 0;
        }
    }
}
