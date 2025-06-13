package main

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"runtime"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/term"
)

// Parameters hof command line
type parameter struct {
	Name         string
	Opt          string // label
	OptLong      string // long version
	Help         string
	Value        any // interface{}
	DefaultValue any
}

type parameters map[string]parameter
type arguments []string

// Internal Conf
const (
	colorsAuto     string = "auto"
	colorsAlways   string = "always"
	colorsNever    string = "never"
	colorsNone     string = "none"
	TypeFile       string = "file"
	TypePipe       string = "pipe"
	TypeFifo       string = "fifo"
	TypeCharDevice string = "charDevice"
	indicatorExe   string = "*"
	indicatorDir   string = "/"
	indicatorLink  string = "@"
	indicatorPipe  string = "|"
	indicatorSock  string = "="
)

type structConfStty struct {
	stdoutType string // default, file, pipe, fifo
	stdinType  string // default, file, pipe, fifo
	stderrType string // default, file, pipe, fifo
	colors     bool   // is tty color capable
}

type configuration struct {
	progName     string // if called with an link, we can change behaviors
	progVersion  string
	os           string
	tty          structConfStty
	ttySizeCol   int
	dotFile      bool   // display or not hidden file
	dotDir       bool   // display or not the . .. file
	colorsWhen   string // auto (default), always, none
	colorsEnable bool
	inode        bool
	format       string
	sortReverse  bool
	oneperline   bool
	dirOnly      bool // just display the directories
	dirFirst     bool // display dir first
	indicator    bool
	sortKey      string

	cwd   string // current working directory
	debug bool
}

// func Conf
func (conf *configuration) configurationInit() {

	conf.progName = os.Args[0]
	conf.os = runtime.GOOS
	conf.progVersion = "0.0.1"
	conf.dotFile = false
	conf.dotDir = false
	// conf.colorsWhen = colorsAuto
	conf.colorsEnable = false
	conf.sortReverse = false
	conf.inode = false
	conf.oneperline = false
	conf.dirOnly = false
	conf.dirFirst = true
	conf.sortKey = "time"
	conf.format = "short"
	conf.indicator = false
	conf.ttySizeCol = 0

	// working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		conf.cwd = "Cannot get working directory."
	} else {
		conf.cwd = cwd
	}

	// stdout
	o, _ := os.Stdout.Stat()
	//     rl, _ := os.Readlink(os.Stdout.Name())
	//     log.Printf("(global) stdout name: %s, fd: %d, link: %s\n", os.Stdout.Name(), os.Stdout.Fd(), rl)
	//     log.Printf("mode       : %032b\n", o.Mode())

	//     log.Printf("mode dir   : %032b\n", os.ModeDir)
	//     log.Printf("mode append: %032b\n", os.ModeAppend)
	//     log.Printf("mode exclus: %032b\n", os.ModeExclusive)
	//     log.Printf("mode Temp  : %032b\n", os.ModeTemporary)
	//     log.Printf("mode Symlin: %032b\n", os.ModeSymlink)
	//     log.Printf("mode dev   : %032b\n", os.ModeDevice)
	//     log.Printf("mode pipe  : %032b\n", os.ModeNamedPipe)
	//     log.Printf("mode socket: %032b\n", os.ModeSocket)
	//     log.Printf("mode suid  : %032b\n", os.ModeSetuid)
	//     log.Printf("mode guid  : %032b\n", os.ModeSetgid)
	//     log.Printf("mode char d: %032b\n", os.ModeCharDevice)
	//     log.Printf("mode sitcky: %032b\n", os.ModeSticky)
	//     log.Printf("mode irregu: %032b\n", os.ModeIrregular)
	//     log.Printf("mode type  : %032b\n", os.ModeType)

	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		//Terminal
		//         conf.tty.stdoutType = TypeCharDevice
		//         rl, _ := os.Readlink(os.Stdout.Name())
		//         log.Printf("(term) stdout name: %s, fd: %d, link: %s\n", os.Stdout.Name(), os.Stdout.Fd(), rl)
		//         log.Printf("stdout: %s", o.Name())
		conf.tty.stdoutType = TypeCharDevice
		conf.ttySizeCol, _, err = term.GetSize(int(os.Stdout.Fd()))
		//         if err != nil {
		//             panic(err)
		//         }
		//     }else {
	}
	if (o.Mode() & os.ModeNamedPipe) == os.ModeNamedPipe {
		// pipe
		// ls algo when en grid mod, just one columns
		//         rl, _ := os.Readlink(os.Stdout.Name())
		//         log.Printf("(pipe)stdout name: %s, fd: %d, link: %s\n", os.Stdout.Name(), os.Stdout.Fd(), rl)
		//         log.Printf("stdout(pipe): %s", o.Name())
		//         conf.ttySizeCol = 80
		conf.tty.stdoutType = TypePipe
	}
	if (o.Mode() & os.ModeType) == 0 {
		// file redirection
		// ls algo when en grid mod, just one columns
		//         rl, _ := os.Readlink(os.Stdout.Name())
		//         log.Printf("(pipe)stdout name: %s, fd: %d, link: %s\n", os.Stdout.Name(), os.Stdout.Fd(), rl)
		//         log.Printf("stdout(pipe): %s", o.Name())
		//         conf.ttySizeCol = 80
		//         conf.ttySizeCol, _, err = term.GetSize(int(os.Stdout.Fd()))
		conf.tty.stdoutType = TypeFile
	}

	// stdin
	o, _ = os.Stdin.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
		// Terminal or file redirection
		conf.tty.stdinType = TypeCharDevice
	} else { //It is not the terminal
		// Display info to a pipe
		conf.tty.stdinType = TypePipe
	}

	// stderr
	o, _ = os.Stderr.Stat()
	if (o.Mode() & os.ModeCharDevice) == os.ModeCharDevice { //Terminal
		//Display info to the terminal
		conf.tty.stderrType = TypeCharDevice
	} else { //It is not the terminal
		// Display info to a pipe
		conf.tty.stderrType = TypePipe
	}

	// debug
	if d, _ := os.LookupEnv("DEBUG"); d == "true" {
		conf.debug = true
	}
}

func (params *parameters) paramsInit() int {
	*params = map[string]parameter{
		"all": {
			Name:    "all",
			Opt:     "a",
			OptLong: "all",
			Help:    "List all files (include hidden)",
			Value:   new(bool),
			//             DefaultValue: false,
			DefaultValue: conf.dotFile,
		},
		"almost-all": {
			Name:    "almost-all",
			Opt:     "A",
			OptLong: "ALL",
			Help:    "Do not list . ..",
			Value:   new(bool),
			//             DefaultValue: false,
			DefaultValue: conf.dotDir,
		},
		// TBD : narrow the possible value.
		"color": {
			Name:    "color",
			Opt:     "c",
			OptLong: "color",
			Help:    "[auto|never|always] when to enable colors (default: auto)",
			Value:   new(string),
			//             DefaultValue: "never",
			DefaultValue: conf.colorsWhen,
		},
		"inode": {
			Name:    "inode",
			Opt:     "i",
			OptLong: "inode",
			Help:    "Print inode",
			Value:   new(bool),
			//             DefaultValue: false,
			DefaultValue: conf.inode,
		},
		"long": {
			Name:    "long",
			Opt:     "l",
			OptLong: "long",
			Help:    "long format: mod|user|group|size|date last modifications|name",
			Value:   new(bool),
			//             DefaultValue: false,
			DefaultValue: false,
		},
		"one": {
			Name:    "one",
			Opt:     "1",
			OptLong: "one",
			Help:    "list onen entry per line",
			Value:   new(bool),
			//             DefaultValue: false,
			DefaultValue: conf.oneperline,
		},
		"dirOnly": {
			Name:    "dirOnly",
			Opt:     "d",
			OptLong: "dir",
			Help:    "list directories not their contents",
			Value:   new(bool),
			//             DefaultValue: false,
			DefaultValue: conf.dirOnly,
		},
		"dirFirst": {
			Name:    "dirFirst",
			Opt:     "g",
			OptLong: "group-directories-first",
			Help:    "group directories first",
			Value:   new(bool),
			//             DefaultValue: true,
			DefaultValue: conf.dirFirst,
		},
		"indicator": {
			Name: "indicator",
			Opt:  "F",
			//             OptLong:      "",
			Help:         "append indicator (one of */=>@|) to entries",
			Value:        new(bool),
			DefaultValue: conf.indicator,
		},
		"sortKey": {
			Name:    "sortKey",
			Opt:     "k",
			OptLong: "sort-key",
			Help:    "key field for sorting",
			Value:   new(string),
			//             DefaultValue: "time",
			DefaultValue: conf.sortKey,
		},
		"reverse": {
			Name:    "revert sort",
			Opt:     "r",
			OptLong: "reverse",
			Help:    "Reverse sort",
			Value:   new(bool),
			//             DefaultValue: false,
			DefaultValue: conf.sortReverse,
		},
		"help": {
			Name:         "help",
			Opt:          "h",
			OptLong:      "help",
			Help:         "Print this help",
			Value:        new(bool),
			DefaultValue: false,
		},
	}

	// set Flags
	for k := range *params {
		if opt, ok := (*params)[k]; ok {

			switch reflect.TypeOf(opt.DefaultValue).String() {

			case "string":
				flag.StringVar(opt.Value.(*string), opt.Opt, opt.DefaultValue.(string), opt.Help)
				flag.StringVar(opt.Value.(*string), opt.OptLong, opt.DefaultValue.(string), opt.Help)
				(*params)[k] = opt

			case "bool":
				flag.BoolVar(opt.Value.(*bool), opt.Opt, opt.DefaultValue.(bool), opt.Help)
				flag.BoolVar(opt.Value.(*bool), opt.OptLong, opt.DefaultValue.(bool), opt.Help)
				(*params)[k] = opt

			case "int":

				flag.IntVar(opt.Value.(*int), opt.Opt, opt.DefaultValue.(int), opt.Help)
				flag.IntVar(opt.Value.(*int), opt.OptLong, opt.DefaultValue.(int), opt.Help)
				(*params)[k] = opt

			default:
				panic("type [ " + reflect.TypeOf(opt.DefaultValue).String() + " ] not implemented !!!")
			}
		}
	}

	// Usage
	flag.Usage = func() {
		flagSet := flag.CommandLine
		fmt.Printf("\n%s (%s) help :\n\n", conf.progName, conf.progVersion)
		order := []string{
			"a", "all",
			"A", "ALL",
			"c", "color",
			"i", "inode",
			"1", "one",
			"d", "dir",
			"F",
			"g", "group-directories-first",
			"r", "reverse",
			"k", "sort-key",
			"h", "help",
		}
		for _, name := range order {
			flag := flagSet.Lookup(name)
			if len(name) == 1 {
				fmt.Printf("  -%-10s%s\n", flag.Name, flag.Usage)
			} else {
				fmt.Printf("  --%-10s\n\n", flag.Name)
			}
		}
	}

	flag.Parse()

	return flag.NFlag()

}

// func args (remaning args)
func (args *arguments) argsInit() {
	if flag.Parsed() {
		*args = flag.Args()
	}
}

func paramsSetConf(confProvided configuration, params parameters) configuration {

	confProvided.dotFile = *params["all"].Value.(*bool)
	confProvided.dotDir = *params["almost-all"].Value.(*bool)

	colors := *params["color"].Value.(*string)
	if (colors == "auto" && conf.tty.stdoutType == TypeCharDevice) ||
		colors == "always" {
		confProvided.colorsEnable = true
	} else {
		confProvided.colorsEnable = false
	}

	confProvided.sortReverse = *params["reverse"].Value.(*bool)
	confProvided.inode = *params["inode"].Value.(*bool)
	confProvided.oneperline = *params["one"].Value.(*bool)
	confProvided.dirOnly = *params["dirOnly"].Value.(*bool)
	confProvided.dirFirst = *params["dirFirst"].Value.(*bool)
	confProvided.sortKey = *params["sortKey"].Value.(*string)
	confProvided.indicator = *params["indicator"].Value.(*bool)

	if *params["long"].Value.(*bool) {
		confProvided.format = "long"
	}
	if *params["one"].Value.(*bool) {
		confProvided.format = "one"
	}

	return confProvided
}

// debug
func (conf configuration) configurationDump() {
	fmt.Println("===================")
	fmt.Println("Configuration")
	fmt.Println("-------------")
	spew.Dump(conf)
	fmt.Println("===================")
}

func (params parameters) paramsDump(argsOS []string) {
	fmt.Println("===================")
	fmt.Println("parameters")
	fmt.Println("----------")
	spew.Dump(argsOS)
	spew.Dump(params)
	fmt.Println("===================")
}

func (args arguments) argsDump() {
	if len(args) > 0 {
		fmt.Println("===================")
		fmt.Println("Arguments")
		fmt.Println("---------")
	}
	for i, arg := range args {
		fmt.Printf("%d: %v\n", i, arg)
	}
	fmt.Println("==================")
}
