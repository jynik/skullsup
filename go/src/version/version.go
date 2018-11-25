// SPDX License Identifier: MIT
package version

import "fmt"

// SkullsUp! API Major version
const Major = 1

// SkullsUp! API Minor version
const Minor = 0

// SkullsUp! API Patch version
const Patch = 1

// SkullsIp! API version description suffix
const Suffix = ""

var String string = fmt.Sprintf("%d.%d.%d%s", Major, Minor, Patch, Suffix)
