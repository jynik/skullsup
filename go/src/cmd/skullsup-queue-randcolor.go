// SPDX License Identifier: MIT
package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"../skullsup"
	"./client"
	c "./common"
)

func main() {
	httpClient, err := client.New(false)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(os.Args) != 1 {
		_, basename := filepath.Split(os.Args[0])
		fmt.Fprintf(os.Stderr, "usage: %s\n", basename)
		fmt.Fprintf(os.Stderr, "Writes a random color to the queue configured in:\n")
		fmt.Fprintf(os.Stderr, "  %s\n", httpClient.ConfigPath)
		os.Exit(1)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// Use a 1/3 scale maximum brightness as not to be obnoxious
	color := c.RandomColors([]skullsup.Range{{0, 85}})[0]
	fmt.Printf("Submitting color: %s\n", color)

	msg := c.Message{Command: "color", Args: []string{color}, Period: "0"}
	_, err = httpClient.WriteMessage(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
