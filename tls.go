package main

import (
	"os"
)

// const Debug = true

var conf configuration
var args arguments
var format formatStruct

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

	// format init
	//     format.formatInit()

	// debug
	if conf.debug {
		conf.configurationDump()
		params.paramsDump(os.Args)
		args.argsDump()
		//         spew.Dump(format)
	}

}

func main() {

	listDir()

	printListFiles()

}
