#ifndef SKULL_HW_H__
#define SKULL_HW_H__

// Do LED addresses increment down a strip, or alternate front-back while
// increasing down the strip (i.e., one side is evens, the other is odds).
#define LAYOUT_INCREMENTING (0 << 0)
#define LAYOUT_ALTERNATING  (1 << 0)

// When one strip ends does addressing continue in the same direction, or
// invert and progress in the other direction?
#define LAYOUT_WRAP_NORMAL  (1 << 0)
#define LAYOUT_WRAP_INVERT  (1 << 1)

#if defined(PLATFORM_skull)
#   define LEDS_PER_STRIP   10
#   define NUM_STRIPS       1
#   define LED_LAYOUT       (LAYOUT_ALTERNATING | LAYOUT_WRAP_NORMAL)
#   define LED_PIN          0
#   define MAX_FRAMES       44

#elif defined(PLATFORM_bulb)
#   define LEDS_PER_STRIP   8
#   define NUM_STRIPS       2
#   define LED_PIN          0
#   define LED_LAYOUT       (LAYOUT_INCREMENTING | LAYOUT_WRAP_INVERT)
#   define MAX_FRAMES       38

#else
#   error "No target platform was specified."
#endif

#define LED_COUNT (LEDS_PER_STRIP * NUM_STRIPS)

// Used by NeoPixel modifications
#define FIXED_PIXEL_COUNT LED_COUNT

#endif
