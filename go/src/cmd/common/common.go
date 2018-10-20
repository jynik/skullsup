// SPDX License Identifier: MIT
package common

import (
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"../../skullsup"
)

// Command-line program flags and arguments
const (
	FLAG_CLIENT_CONFIG       = "client-config"
	FLAG_CLIENT_CONFIG_DESC  = "Path to client configuration file"
	FLAG_CLIENT_CONFIG_SHORT = 'c'

	FLAG_DEVICE       = "skull"
	FLAG_DEVICE_DESC  = "Path to our Skull Lord and Master"
	FLAG_DEVICE_SHORT = 's'

	FLAG_HOST       = "remote"
	FLAG_HOST_DESC  = "Remote server running skullsup-queue-server"
	FLAG_HOST_SHORT = 'r'

	FLAG_PERIOD             = "period"
	FLAG_PERIOD_DEV_DESC    = "Frame update period in ms"
	FLAG_PERIOD_CLIENT_DESC = "Poll server for new data every <period> seconds"
	FLAG_PERIOD_SHORT       = 'P'

	FLAG_PORT       = "port"
	FLAG_PORT_DESC  = "Queue server port number"
	FLAG_PORT_SHORT = 'p'

	FLAG_TLS_CERT       = "tls-cert"
	FLAG_TLS_CERT_DESC  = "Specifies path to TLS certificate"
	FLAG_TLS_CERT_SHORT = 't'

	FLAG_PRIVATE_KEY       = "private-key"
	FLAG_PRIVATE_KEY_DESC  = "Path to TLS certificate private key"
	FLAG_PRIVATE_KEY_SHORT = 'k'

	FLAG_LADDRESS       = "address"
	FLAG_LADDRESS_DESC  = "Local address to listen on. An empty string implies \"all interfaces\""
	FLAG_LADDRESS_SHORT = 'A'

	FLAG_QUIET       = "quiet"
	FLAG_QUIET_DESC  = "Disable logging. Overrides --log and --verbose"
	FLAG_QUIET_SHORT = 'q'

	FLAG_LOGFILE       = "log"
	FLAG_LOGFILE_DESC  = "Specifies log file to write to. May also be set to \"stderr\" or \"stdout\""
	FLAG_LOGFILE_SHORT = 'L'

	FLAG_VERBOSE       = "verbose"
	FLAG_VERBOSE_DESC  = "Increase log vebosity"
	FLAG_VERBOSE_SHORT = 'v'

	FLAG_VERSION      = "version"
	FLAG_VERSION_DESC = "Print SkullUp! version and exit"

	FLAG_TLS_FOOTGUN      = "insecure"
	FLAG_TLS_FOOTGUN_DESC = "Disable TLS certificate verification for testing or because you're an idiot"

	CMD_COLOR       = "color"
	CMD_COLOR_DESC  = "Cast light upon the Wicked and Almighty Skull"
	CMD_COLOR_ALIAS = "c"

	ARG_COLORVAL      = "value"
	ARG_COLORVAL_DESC = "Hexadecimal RGB color value. Example: ff8330"

	CMD_INCANT       = "incant"
	CMD_INCANT_DESC  = "Incant a prepared Unholy Psalm, with optional incantation-specific arguments"
	CMD_INCANT_ALIAS = "i"

	ARG_PSALM      = "psalm"
	ARG_PSALM_DESC = "The desired Psalm. Run list-psalms to view your options.."

	ARG_PSALM_ARGS      = "args"
	ARG_PSALM_ARGS_DESC = "Psalm-specific arguments"

	CMD_LIST       = "list"
	CMD_LIST_DESC  = "List the Unholy Psalms prepared for unworthy mortals"
	CMD_LIST_ALIAS = "l"

	CMD_REANIM       = "reanimate"
	CMD_REANIM_DESC  = "Reanimate His Dark Unholiness in a form of your choosing"
	CMD_REANIM_ALIAS = "r"

	ARG_FRAMESTR      = "frame"
	ARG_FRAMESTR_DESC = "Specify frames in the form: <RBG color>[:<LED #>[:n]. " +
		"Unless an LED # is specified, all LEDs will be set to " +
		"the provided color. The optional \"n\" flag suppresses" +
		" the intra-frame delay to allow a subset of LEDs to be" +
		" updated simultaneously."

	CMD_VERSION      = FLAG_VERSION
	CMD_VERSION_DESC = FLAG_VERSION_DESC
)

// Error strings
const (
	ERRTAG = "[Error]"

	ERR_INVALID_CMD = "Do not disturb The Skull with this nonsense."
	ERR_INVAL       = "You've been cast to the land of the Damned."
	ERR_TIMEOUT     = "\nThe Skull is not listening.\n  He hath foresaken us.\n    Abandoning all hope.\n"
	ERR_METHOD      = "The Skull does not understand your intentions."
	ERR_FAILURE     = "Chaos ensues; the world is ablaze!"
	ERR_FULL        = "Queue is full. Try again later."
	ERR_EMPTY       = "Queue is empty. Try again later."
	ERR_AUTH        = "The Skull forbids this. You are not worthy!"
	ERR_BUG         = "A plague is upon us!"
)

// HTTP server and client-related strings
const (
	ENDPOINT         = "/hell"
	HEADER_QUEUE     = "x-skull-queue"
	HEADER_QUEUE_SEP = ":"
)

type Message struct {
	Command string   `json:"cmd"`
	Args    []string `json:"args"`
	Period  string   `json:"period"`
}

func (m *Message) String() string {
	return fmt.Sprintf("{ %s %s (%s ms) }", m.Command, m.Args, m.Period)
}

var uuidRegexp *regexp.Regexp = regexp.MustCompile("[[:xdigit:]]{8}-[[:xdigit:]]{4}-4[[:xdigit:]]{3}-[89abAB][[:xdigit:]]{3}-[[:xdigit:]]{12}")

func IsValidQueueName(s string) bool {
	return uuidRegexp.MatchString(s)
}

func OpenDevice(device string) (*skullsup.Skull, error) {
	var skull *skullsup.Skull
	var err error

	retry := true
	max := 3
	for i := 0; i < max && retry; i++ {
		skull, err = skullsup.New(device)
		if err != nil {
			if err.Error() == skullsup.ErrorNotReady {
				if i < (max - 1) {
					fmt.Printf("%s. Retrying %d more time(s)...\n", err, max-i-1)
					time.Sleep(1 * time.Second)
				}
			} else {
				return nil, err
			}
		} else {
			return skull, err
		}
	}

	return nil, errors.New(ERR_TIMEOUT)
}

func PrintPsalms() {
	psalms := skullsup.Psalms()
	fmt.Println("Mortals may beckon The Skull using these Unholy Psalms:")
	for _, psalm := range psalms {
		fmt.Println("    " + psalm.Name)
	}
}

func EndpointURL(host string, port uint16) string {
	return "https://" + host + ":" + strconv.Itoa(int(port)) + ENDPOINT
}

func RandomInt(min, max int) int {
	return rand.Intn(max+1-min) + min
}

/* Generate random colors within the provided luma ranges.
 *
 * This uses a rough YUV -> RGB conversion so that we can try to keep the
 * brightness associated with each argument within a desired range
 * (and in a manner that's sensitive to our color perception).
 */
func RandomColors(lumas []skullsup.Range) []string {
	colors := []string{}
	for _, l := range lumas {
		y := float32(RandomInt(l.Min, l.Max))
		u := float32(rand.Intn(256))
		v := float32(rand.Intn(256))

		r := 1.164*(y-16) + 1.596*(v-128)
		if r > 255 {
			r = 255
		} else if r < 0 {
			r = 0
		}

		g := 1.164*(y-16) - 0.813*(v-128) - 0.391*(u-128)
		if g > 255 {
			g = 255
		} else if g < 0 {
			g = 0
		}

		b := 1.164*(y-16) + 2.018*(u-128)
		if b > 255 {
			b = 255
		} else if b < 0 {
			b = 0
		}

		colors = append(colors, fmt.Sprintf("%02x%02x%02x",
			uint(r)&0xff, uint(g)&0xff, uint(b)&0xff))
	}

	return colors
}
