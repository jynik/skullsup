// SPDX License Identifier: MIT
package network

import (
	"strconv"
	"strings"
)

const QueueEndpoint = "hell"

const (
	ErrorQueueFull  = "There's no room left for the Damned in Hell."
	ErrorQueueEmpty = "We're fresh out of souls. Reap again later."
)

func QueueUrl(host string, port uint16, queue string) string {
	return "https://" + host + ":" + strconv.Itoa(int(port)) + "/" + QueueEndpoint + "/" + queue
}

func QueueFromURL(url string) string {
	pfx := "/" + QueueEndpoint + "/"
	if !strings.HasPrefix(url, pfx) {
		return ""
	}

	return url[len(pfx):]
}
