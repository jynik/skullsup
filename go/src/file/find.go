// SPDX License Identifier: MIT
package file

import (
	"fmt"
	"io/ioutil"
	"os"
)

// Search for file and read its contents
func FindAndRead(filename string) ([]byte, error) {

	// Try assuming CWD or full path first
	if data, err := ioutil.ReadFile(filename); err == nil {
		return data, nil
	}

	// Search OS-specific list of standard locations
	for _, pfx := range searchPrefixes {
		path := os.ExpandEnv(pfx + filename)
		if data, err := ioutil.ReadFile(path); err == nil {
			return data, nil
		}
	}

	return []byte{}, fmt.Errorf("Unable to find file: %s\n", filename)
}
