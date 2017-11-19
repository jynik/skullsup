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

/*
 * Go with incantation defaults 10% of the time.
 * It is undesirable to specify only a portion of args randomly,
 * as the luma for each arg won't be balanced nicely.
 */

func randomPsalmArgs(p *skullsup.Psalm) []string {
	if rand.Intn(10) == 0 {
		c.RandomColors([]skullsup.Range{})
	}

	return c.RandomColors(p.Luma)
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
	period := strconv.Itoa(c.RandomInt(psalm.Period.Min, psalm.Period.Max))
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
