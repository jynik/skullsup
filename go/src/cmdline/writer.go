// SPDX License Identifier: MIT
package cmdline

import "flag"

type WriterFlags struct {
	Cfg string

	Queue  string
	Period int
}

const cfgHelp = "Load the specified configuration file."
const queueHelp = "Write the the specified queue."
const periodHelp = "Intra-frame period, in ms."

func (f *WriterFlags) Init() {
	flag.StringVar(&f.Cfg, "cfg", "", cfgHelp)

	flag.StringVar(&f.Queue, "queue", "", queueHelp)

	if f.Period >= 0 {
		flag.IntVar(&f.Period, "period", -1, periodHelp)
	}
}

func (f *WriterFlags) LocalInit() {
	flag.StringVar(&f.Cfg, "cfg", "", cfgHelp)
	flag.IntVar(&f.Period, "period", -1, periodHelp)
}
