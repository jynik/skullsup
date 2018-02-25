#ifndef SKULL_HW_H__
#define SKULL_HW_H__

#if defined(PLATFORM_skull)
#   define LED_COUNT 10
#   define LED_PIN 0
#   define MAX_FRAMES 44
#elif defined(PLATFORM_bulb)
#   define LED_COUNT 16
#   define LED_PIN 0
#   define MAX_FRAMES 38
#else
#   error "No target platform was specified."
#endif

#define FIXED_PIXEL_COUNT LED_COUNT

#endif
