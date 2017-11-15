// SPDX License Identifier: MIT
package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"../skullsup"
	"./client"
	c "./common"
)

func randomInt(min, max int) int {
	return rand.Intn(max+1-min) + min
}

/*
 * Generate random colors.

 * Uses a rough YUV -> RGB conversion so that we can try to keep the
 * brightness associated with each argument within a reasonable range
 * (and in a manner that's sensitive to our color perception).
 */

func randomPsalmArgs(p *skullsup.Psalm) []string {
	var ret []string
	var n int

	// Go with defaults 10% of the time, and rand args 80%
	// It is undesirable to specify only a portion of args randomly,
	// as the luma for each arg won't be balanced nicely.
	if rand.Intn(10) == 0 {
		n = 0
	} else {
		n = p.Args.Max
	}

	for i := 0; i < n; i++ {

		y := float32(randomInt(p.Luma[i].Min, p.Luma[i].Max))
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

		ret = append(ret, fmt.Sprintf("%02x%02x%02x",
			uint(r)&0xff,
			uint(g)&0xff,
			uint(b)&0xff))
	}
	return ret
}

func main() {
	httpClient, err := client.New(false)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(os.Args) != 1 {
		_, basename := filepath.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "usage: %s\n", basename)
		fmt.Fprintf(os.Stderr, "Writes a random incantation to the queue configured in:\n")
		fmt.Fprintf(os.Stderr, "  %s\n", httpClient.ConfigPath)
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	psalms := skullsup.Psalms()
	psalm := psalms[rand.Intn(len(psalms))]
	period := strconv.Itoa(randomInt(psalm.Period.Min, psalm.Period.Max))
	args := []string{psalm.Name}
	args = append(args, randomPsalmArgs(&psalm)...)

	fmt.Printf("Submitting incantation { args=%s, period=%s ms }\n", args, period)

	msg := c.Message{Command: "incant", Args: args, Period: period}
	_, err = httpClient.WriteMessage(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
