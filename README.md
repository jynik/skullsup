<pre>
                 ██████  ██ ▄█▀ █    ██  ██▓     ██▓      ██████  █    ██  ██▓███   ▐██▌
               ▒██    ▒  ██▄█▒  ██  ▓██▒▓██▒    ▓██▒    ▒██    ▒  ██  ▓██▒▓██░  ██▒ ▐██▌
               ░ ▓██▄   ▓███▄░ ▓██  ▒██░▒██░    ▒██░    ░ ▓██▄   ▓██  ▒██░▓██░ ██▓▒ ▐██▌
                 ▒   ██▒▓██ █▄ ▓▓█  ░██░▒██░    ▒██░      ▒   ██▒▓▓█  ░██░▒██▄█▓▒ ▒ ▓██▒
               ▒██████▒▒▒██▒ █▄▒▒█████▓ ░██████▒░██████▒▒██████▒▒▒▒█████▓ ▒██▒ ░  ░ ▒▄▄
               ▒ ▒▓▒ ▒ ░▒ ▒▒ ▓▒░▒▓▒ ▒ ▒ ░ ▒░▓  ░░ ▒░▓  ░▒ ▒▓▒ ▒ ░░▒▓▒ ▒ ▒ ▒▓▒░ ░  ░ ░▀▀▒
               ░ ░▒  ░ ░░ ░▒ ▒░░░▒░ ░ ░ ░ ░ ▒  ░░ ░ ▒  ░░ ░▒  ░ ░░░▒░ ░ ░ ░▒ ░      ░  ░
               ░  ░  ░  ░ ░░ ░  ░░░ ░ ░   ░ ░     ░ ░   ░  ░  ░   ░░░ ░ ░ ░░           ░
                     ░  ░  ░      ░         ░  ░    ░  ░      ░     ░               ░
</pre>

[![The Skull](https://www.dropbox.com/s/s33y501nyxj91fd/skull-link.jpg?raw=1)](https://youtu.be/DmYyMnP-sAg)

# Overview #

This project provides device firmware and a couple command-line programs that
can be used to animate and illuminate a glass :skull:.

Toss a `skullsup` invocation into your script of choice to signal a broken
build and kick off a `git blame` driven witch hunt (*you know it'll just be your
fault*), alert you to a new email that you probably didn't want to read anyway,
or just let you know when the form of the [Destructor] has been chosen.

[Destructor]: http://www.imdb.com/title/tt0087332/quotes

https://youtu.be/DmYyMnP-sAg## Stand Alone Mode ##

The setup can be operated in a stand-alone mode, with the :skull: connected
directly to a host via USB. Run `skullsup --help` for a list of commands, and
`skullsup help <command>` for more information about a specific command.

## Internet of Terror ##

Alternatively, you can use the grossly overengineered IoT setup
consisting of...

### skullsup-queue-server ###

`skullsup-queue-server` is an HTTPS server and presents a JSON API for queueing
up commands to a device. This supports multiple queues, which means you can
deploy an armada of demonic skulls or have multiple users queueing up commands
to a single skull.

Queues can be assigned on a per-user basis, or shared amongst multiple users.
In most cases, you'll assign each writer (a user submitting commands) a
fixed-length queue and then grant the reader attached to the :skull: read
access to each queue.

Accesses to queues are currently authenticated via Basic Authentication (over
TLS `¯\_(ツ)_/¯`). Credentials are stored in server-side configuration files
as bcrypt hashes (cost=11). Eventually this will (might?) all be replaced with
TLS client certificate-based authentication.

### skullsup-queue-writer ###

`skullsup-queue-writer` is HTTPS client that submits commands to a queue.
The usage of this is largely the same as `skullsup`; you simply just point
it to the remote `skullsup-queue-server` use the same commands you've already
been using.

### skullsup-queue-reader ###

`skullsup-queue-reader` should be run on a host or platform connected a :skull:
over USB. This we dequeue items from one or more queues managed by a
`skullsup-queue-server` and display them on a device. This currently polls the
server at a specified rate. (A WebSockets PR would be quite welcome.)

## Build ##

For build instructions and dependencies, see the `README.md` in each of the
following directories:

* [firmware](./firmware) - DigiSpark firmware
* [go](./go) - Application and library code written in Go
* [yocto/meta-skullsup](./yocto/meta-skullsup) - A Yocto layer that can be used
to build barebones Linux images for running a `skullsup-queue-reader` daemon on
your favorite embedded platform.

# Hardware #

The Skull consists of a [Crystal Head Vodka] glass containing 10 [NeoPixels]
and clear plastic [vase filler]. The [firmware](./firmware) running on the
[Digispark] presents USB Serial interface (via the DigiCDC library).

![Entrails](https://www.dropbox.com/s/1fr6voigxz2nb77/skullsup-hw.jpg?raw=1)

[Crystal Head Vodka]: https://www.crystalheadvodka.com
[NeoPixels]: https://www.adafruit.com/product/1655
[Digispark]: http://digistump.com/products/1
[vase filler]: http://www.michaels.com/ashland-clear-mini-discs/10998221.html
