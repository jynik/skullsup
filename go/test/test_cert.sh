#!/bin/sh
yes "" | openssl req -nodes -new -x509 -keyout skullsup-queue-server.key -out skullsup-queue-server.pem

