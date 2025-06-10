/*
main.main
*/
package main

import (
	"fmt"
	"os"
	"time"
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

	// debug
	if conf.debug {
		conf.configurationDump()
		params.paramsDump(os.Args)
		args.argsDump()
	}
}

func main() {

	listbegin := time.Now().UnixMilli()
	listDir()
	listend := time.Now().UnixMilli()

	begin := time.Now().UnixMilli()
	if len(dirEntries) > 0 {
		printListFiles()
	}
	end := time.Now().UnixMilli()

	fmt.Printf("list: %d\n", listend-listbegin)
	fmt.Printf("print: %d\n", end-begin)

}
