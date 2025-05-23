package main

// on fait quoi aujourd'hui ???
// ajout de logs (slog)
// ajout prise en charge des args ligne de commandes
// on ajoute les inodes
// formatage => couleurs
// organiser les sources (ajout de fonctions...)
import (
	"os"
)

// const Debug = true

var conf configuration
var args arguments

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

	listDir()
	printListFiles()

}
