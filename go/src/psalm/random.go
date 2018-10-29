// SPDX License Identifier: MIT
package psalm

import (
	"math/rand"
	"strings"

	"github.com/jynik/skullsup/go/src/color"
)

func randomArgs(p *Psalm) []string {
	// 10% of the time us the default args
	if rand.Intn(10) == 0 {
		return []string{p.Name}
	}

	// Otherwise get random colors for each arg, following
	// the recommended luma ranges
	args := make([]string, p.Args.Max+1)
	args[0] = p.Name

	for i := 0; i < p.Args.Max; i++ {
		lumaMin := p.Luma[i].Min
		lumaMax := p.Luma[i].Max
		args[i+1] = color.Random(lumaMin, lumaMax).String()
	}

	return args
}

// If a psalm name is provided, return that psalm name with random arguments.
// If the name is empty or invalid, choose a random psalm with random args.
func Random(name string) []string {
	name = strings.ToLower(name)
	for _, psalm := range List {
		if name == psalm.Name {
			return randomArgs(&psalm)
		}
	}

	numPsalms := len(List)
	idx := rand.Intn(numPsalms)
	psalm := &List[idx]
	return randomArgs(psalm)
}
