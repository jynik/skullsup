![skull](https://user-images.githubusercontent.com/3046210/48930331-640f0f80-eeb5-11e8-826e-728b881eda93.png)

# Overview #

This project provides device firmware and a handful of programs that
can be used to animate and illuminate a glass :skull:.

Toss a `skullsup` invocation into your script of choice to:

  * Signal a broken build that results in a `git blame` driven witch hunt 
  * Alert you to a new email (that you probably didn't want to read anyway)
  * Just let you know when the form of the [Destructor] has been chosen

The hardware consists of two [Adafruit NeoPixel Sticks], an [ATtiny85-based
Digispark board], and a [Crystal Head Vodka] glass.

A little demo can be [found here]. 

[Adafruit NeoPixel Sticks]: https://www.adafruit.com/product/1426
[ATtiny85-based Digispark board]: http://digistump.com/products/1
[Crystal Head Vodka]: https://www.crystalheadvodka.com
[Destructor]: https://www.imdb.com/title/tt0087332/quotes/qt1767149
[found here]: https://www.youtube.com/watch?v=VBeKjXn2xuM

## Rot-Alone Mode ##

The setup can be operated in a rot-alone mode, with the :skull: connected
directly to a host via a [USB to UART](https://www.sparkfun.com/products/12731) adapter. 

To build the firmware, you'll need to follow the Digispark installation guide.
From there, you can just run a `make && make install` in the [firmware] directory.

[Digispark installation guide]: https://digistump.com/wiki/digispark/tutorials/connecting
[firmware]: firmware


To fetch and build the code host code, run the following:

~~~
go get github.com/jynik/skullsup/go/src/cmd/skullsup
go build github.com/jynik/skullsup/go/src/cmd/skullsup
~~~

Run `./skullsup --help` for a usage information.

To get started, try running the *pulse* animation with a diabolical purple
color. Note that the color is specified as a [hexadecimal string].

~~~
./skullsup --device /dev/ttyUSB0 incant pulse ff00ff
~~~

Or perhaps *vortex* with a irreverent red and ghastly green are more of your thing?

~~~
./skullsup --device /dev/ttyUSB0 incant vortex ff0000 000500
~~~

Note that the `--period` option configures the period of time between "frames"
in the animation and is specified in units of milliseconds. You can use this to
speed up or slow down animations.

If you just want to stop the animation and set the device to a melancholy blue:

~~~
./skullsup --device /dev/ttyUSB0 color 05052a
~~~

Note that you can use the `reanimate` command to experiment with your own
custom animations before committing them to code. 

The arguments to this command describe each frame of the desired animation.

* First, a hexidecimal color value is required.

* Optionally, the ID of the LED you want to change can be provided. If not
  specified, ```<color>:all``` is implicit. LEDs are addressed starting at 0
  and increment their way along the NeoPixel sticks.

* Finally, a "No-Delay" flag ("*N*") can be specified. This skips the delay
  after a frame (controlled by `--period`), allowing you to change
  specific LEDs between the animation frame period.

~~~
<color>[:LED_ID[:N]]
~~~

A simpler example of the *pulse* animation (in green) could be run as follows:

~~~
./skullsup --period 50 --device /dev/ttyUSB0 reanimate 000000 000500 000a00 001000 002000 003000 002000 001000 000a00 000500
~~~

Next, we could add a couple frames that add a tinge of orange at the brightest
point in the animation by specifically adjusting LEDs 6 and 9. Note that 
the "No-Delay" flag is set when we set LED 6 so that both LEDs appear to 
change at the same time. The inter-frame delay occurs after LED 9 is updated.

~~~
./skullsup --period 50 --device /dev/ttyUSB0 reanimate 000000 000500 000a00 001000 002000 003000 804000:6:N 804000:9 003000 002000 001000 000a00 000500
~~~

[hexadecimal string]: https://www.w3schools.com/colors/colors_picker.asp

## :fire: Internet of Terror :fire: ##

What good is it to awaken a sleeping demon in a colorful display of Hellish glory
if you have no one to share it with!? :fire:

In addition to the Rot-Alone mode, you can use the grossly overengineered IoT
setup in which the Skull is conntected to Raspberry Pi Zero W
running a custom Linux build. One or more skulls retrieve animations and
settings from a server, and the programs descriped below may be used to
submit animiations the the server.

Note that you'll need a [logic level converter] to connect the DigiSpark's 5V
UART to the RPi0's 3.3V UART.

The [yocto/meta-skullsup] directory contains a [Yocto] layer that can be used to
build a ready-to-go Linux image. See the README in that directory for how to 
configure the image prior to building it. 

If you're new to Yocto, you might want to start with [this tutorial] first.

[logic level converter]: https://www.sparkfun.com/products/12009
[yocto/meta-skullsup]: yocto/meta-skullsup 
[Yocto]: https://yoctoproject.org
[this tutorial]: https://github.com/jynik/ready-set-yocto


### skullsup-queue-server ###

`skullsup-queue-server` is an HTTPS server that presents a simple API
to allow users to queue up commands (i.e. `color`, `incant`, `reanimate`)
to one or more remote Skulls.

This supports multiple queues, which means you can deploy an armada of demonic
skulls or have multiple users queueing up commands each others' skulls.

Queues can be assigned on a per-user basis, or shared amongst multiple users.
In most cases, you'll assign each writer (a user submitting commands) their
own queue, and then grant the reader attached to the :skull: read access to the
corresponding queue.

Some example configuration files can be found in the [go/test/configs]
directory.

Authenication is performed via mutual TLS; each client must 
provide their own TLS certificate, signed by a CA that is trusted 
by `skullsup-queue-server`.

Check out the [gpgsm-as-ca] project for information on how to put together
your own little smart-card based CA for creating and/or signing client
certificates.

Why all the hassle? Well, what's the point of making an absurd IoT Skull if 
you don't have any excuse to put your [Nitrokey] to use?

[go/test/configs]: go/test/configs
[Nitrokey]: https://www.nitrokey.com
[gpgsm-as-ca]: https://github.com/jymigeon/gpgsm-as-ca

### skullsup-queue-reader ###

The `skullsup-queue-reader` program runs on the host (RPi0W) connected to the :skull:.

This application dequeuesa color and animation commands from one or more queues managed by a
`skullsup-queue-server` and display them on a device. 

### skullsup-queue-writer ###

`skullsup-queue-writer` is client that submits commands to a queue.
The usage of this is largely the same as `skullsup`; you simply just point
it to the remote `skullsup-queue-server` use the same commands you've already
been using.

### skullsup-queue-incantor ###

`skullsup-queue-incantor` is a variant of `skullsup-queue-writer` that
submits random incantations to a queue. 

This is intended to be used with a `skullsup-client.conf` is located in the
default location, allowing this to be used when arguments cannot be passed to
the program. (This is the case for some email clients that allow an application
to be run as a notification or inbox rule action.)

### skullsup-queue-randcolor ###

`skullsup-queue-randcolor` is similar to `skullsup-queue-incantor`, except
that it send submits random colors to the queue.
