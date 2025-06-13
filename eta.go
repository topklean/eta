/*
main.main
for git
*/
package main

import (
	//     "net/http/pprof"

	"github.com/davecgh/go-spew/spew"
)

// const Debug = true

var conf configuration
var args arguments

// var format formatStruct

func init() {
	var params parameters

	// conf
	conf.configurationInit()

	// my personal usage
	// params & args
	params = make(parameters)
	paramNums := params.paramsInit()
	args.argsInit()

	if paramNums > 0 {
		conf = paramsSetConf(conf, params)
	}

	// init colors
	getColors()

	// debug
	if conf.debug {
		spew.Dump(colorsMap)
		//         conf.configurationDump()
		//         params.paramsDump(os.Args)
		//         args.argsDump()
	}
}

func main() {
	// TBD : hyperlink on the terminal like eza

	//     f, err := os.Create("cpu.prof")
	//     if err != nil {
	//         log.Fatal("Could not create CPU profile: ", err)
	//     }
	//     defer f.Close()

	//     if err := pprof.StartCPUProfile(f); err != nil {
	//         log.Fatal("could not start CPU profile: ", err)
	//     }
	//     defer pprof.StopCPUProfile()

	//     listbegin := time.Now().UnixMilli()

	listDir()
	//     listend := time.Now().UnixMilli()

	//     begin := time.Now().UnixMilli()
	//     printColors()
	//     log.Println("begin print by column")
	if len(dirEntries) > 0 {
		printListFiles()
	}
	//     end := time.Now().UnixMilli()

	//     fmt.Printf("list global: %d\n", listend-listbegin)
	//     fmt.Printf("print: %d\n", end-begin)
}
