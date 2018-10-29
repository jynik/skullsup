// SPDX License Identifier: MIT
package psalm

import (
	"errors"
	"strings"

	"github.com/jynik/skullsup/go/src/frame"
)

type Psalm struct {
	Name     string
	Args     Range
	ArgNames []string
	Period   Range
	Luma     []Range // Luma range per argument
	impl     func([]string, uint) ([]frame.Frame, uint16, error)
}

var List []Psalm = []Psalm{
	{
		Name:     "hellivator",
		Args:     Range{0, 2},
		ArgNames: []string{"foreground", "background"},
		Period:   Range{65, 150},
		Luma:     []Range{{128, 255}, {0, 32}},
		impl:     hellivator,
	},

	{
		Name:     "pulse",
		Args:     Range{0, 1},
		ArgNames: []string{"color"},
		Period:   Range{50, 85},
		Luma:     []Range{{32, 255}},
		impl:     pulse,
	},

	{
		Name:     "vortex",
		Args:     Range{0, 2},
		ArgNames: []string{"foreground", "background"},
		Period:   Range{65, 125},
		Luma:     []Range{{64, 255}, {0, 64}},
		impl:     vortex,
	},
}

func Lookup(name string, args []string, ledCount uint) ([]frame.Frame, uint16, error) {
	if ledCount < 1 || ledCount&0x1 != 0 {
		return []frame.Frame{}, 0, errors.New("Invalid LED count")
	}

	nameLower := strings.ToLower(name)
	for _, psalm := range List {
		if nameLower == psalm.Name {
			return psalm.impl(args, ledCount)
		}
	}

	return []frame.Frame{}, 0, errors.New("No such psalm: " + name)
}
