# Overview #

This directory contains firmware for a Digistump Digispark device.

# Dependencies #

* avr-gcc >= 4.9.2
* Arduino libraries (shipped with the IDE) >= 1.8.3+
* [DigiCDC Library](https://digistump.com/wiki/digispark/tutorials/digicdc)

# Adafruit NeoPixel Modifications #

This firmware contains a modified version of the [Adafruit NeoPixel Library](https://github.com/adafruit/Adafruit_NeoPixel).
The changes to the library, discussed in this [issue tracker item], were needed to
reduce the size of the library in order to recover sufficient codespace.

[issue tracker item]: https://github.com/adafruit/Adafruit_NeoPixel/issues/142 

# Build #

~~~
make && make install
~~~

