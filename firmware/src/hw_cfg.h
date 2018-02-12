#ifndef SKULL_HW_H__
#define SKULL_HW_H__

#if defined(PLATFORM_skull)
#   define LED_COUNT 10
#   define LED_PIN 0
#else
#   error "No target platform was specified."
#endif

#define FIXED_PIXEL_COUNT LED_COUNT

#endif
