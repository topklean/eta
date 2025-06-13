package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// Escape Code
// ESC[FF;BBm
// FF = forground
// BB = backgroung
// ESC => \e \033
// Reset : ESC[39;49m
// Reset all attributs: ESC[0m

// TBD : use package color https://pkg.go.dev/github.com/fatih/color@v1.18.0

const (
	CC = "\033["
	CE = "m"

	BOLD      = "1"
	DIMFAINT  = "2"
	ITALIC    = "3"
	UNDERLINE = "4"
	BLINK     = "5"
	INVERSE   = "7"
	HIDDEN    = "8"
	STRIKE    = "9"

	// RESET ALL MODE
	RESET = CC + "0" + CE

	RESET_BOLD      = "22"
	RESET_DIMFAINT  = "22"
	RESET_ITALIC    = "23"
	RESET_UNDERLINE = "24"
	RESET_BLINK     = "25"
	RESET_INVERSE   = "27"
	RESET_HIDDEN    = "28"
	RESET_STRIKE    = "29"

	FG_Black          = "30"
	FG_Red            = "31"
	FG_Green          = "32"
	FG_Yellow         = "33"
	FG_Blue           = "34"
	FG_Magenta        = "35"
	FG_Cyan           = "36"
	FG_White          = "37"
	FG_Bright_Black   = "90"
	FG_Bright_Red     = "91"
	FG_Bright_Green   = "92"
	FG_Bright_Yellow  = "93"
	FG_Bright_Blue    = "94"
	FG_Bright_Magenta = "95"
	FG_Bright_Cyan    = "96"
	FG_Bright_White   = "97"
	BG_Black          = "40"
	BG_Red            = "41"
	BG_Green          = "42"
	BG_Yellow         = "43"
	BG_Blue           = "44"
	BG_Magenta        = "45"
	BG_Cyan           = "46"
	BG_White          = "47"
	BG_Bright_Black   = "100"
	BG_Bright_Red     = "101"
	BG_Bright_Green   = "102"
	BG_Bright_Yellow  = "103"
	BG_Bright_Blue    = "104"
	BG_Bright_Magenta = "105"
	BG_Bright_Cyan    = "106"
	BG_Bright_White   = "107"
)

// type color struct {
//     fg     string
//     bg     string
//     effect string
// }

// map pour les diff objets d'un rep
// "dir" "link" "char" "block" ...
// mode: read write exe suid guid etc..
//			pour les noms et pour les mod:

// TBD: do it clean please, with a struct...
var colorsMap map[string]string

func printColors() {
	//     fmt.Printf(CC + FG_Blue + CE + "█" + RESET);
	for i := range 8 {
		fg := strconv.Itoa(i + 30)
		for j := range 8 {
			bg := strconv.Itoa(j + 40)
			code := CC + fg + ";" + bg + CE
			fmt.Printf(code + "0" + RESET)
			//             fmt.Printf(code + "█" + RESET);
		}
		fmt.Println()
	}
	// BG Bright
	for i := range 8 {
		fg := strconv.Itoa(i + 30)
		for j := range 8 {
			bg := strconv.Itoa(j + 100)
			code := CC + fg + ";" + bg + CE
			fmt.Printf(code + "0" + RESET)
			//             fmt.Printf(code + "█" + RESET);
		}
		fmt.Println()
	}
	//FG Bright
	for i := range 8 {
		fg := strconv.Itoa(i + 90)
		for j := range 8 {
			bg := strconv.Itoa(j + 40)
			code := CC + fg + ";" + bg + CE
			fmt.Printf(code + "0" + RESET)
			//             fmt.Printf(code + "█" + RESET);
		}
		fmt.Println()
	}
	// bright
	for i := range 8 {
		fg := strconv.Itoa(i + 90)
		for j := range 8 {
			bg := strconv.Itoa(j + 100)
			code := CC + fg + ";" + bg + CE
			fmt.Printf(code + "0" + RESET)
			//             fmt.Printf(code + "█" + RESET);
		}
		fmt.Println()
	}
}

func getColors() {
	// init colors map
	colorsMap = make(map[string]string)
	colorsMap["reset"] = RESET
	// regex for validating the code format
	// //x;x;fg;bg;o

	reColorCode := regexp.MustCompile(`^\d?\d(;\d?\d)*`)

	const (
		key = 0
		val = 1
	)

	ENV_LS_COLORS, ok := os.LookupEnv("LS_COLORS")

	if ok {
		for lscolor := range strings.SplitSeq(ENV_LS_COLORS, ":") {
			// if last char is the separator
			if len(lscolor) == 0 {
				continue
				//                 break
			}
			// must have = as separator
			if strings.ContainsRune(lscolor, '=') {

				color := strings.Split(lscolor, "=")

				if !reColorCode.MatchString(color[1]) {
					//                     fmt.Printf("err: key=%s, color=%s\n", color[key], color[val])
					// just ignore the value and step to the next one
					continue
				}
				// get each one
				switch color[key] {
				//                 case "lc": // Left·of·color·sequence
				//                     fmt.Printf("%s (left code)   : %s\n", color[key], color[val])
				//                 case "rc": // Right·of·color·sequence
				//                     fmt.Printf("%s (right code)  : %s\n", color[key], color[val])
				//                 case "ec": // end color (replaces·lc+rs+rccolor[key], )
				//                     fmt.Printf("%s (end color)   : %s\n", color[key], color[val])
				//                 case "rs": // Reset to ordinary colors
				//                     fmt.Printf("%s (reset)       : %s\n", color[key], color[val])
				case "no": // normal
					colorsMap["normal"] = CC + color[val] + CE
				case "fi": // file default
					colorsMap["default"] = CC + color[val] + CE
				case "di": // directory
					colorsMap["dir"] = CC + color[val] + CE
				case "ln": // link
					colorsMap["link"] = CC + color[val] + CE
				case "pi": // pipe
					colorsMap["pipe"] = CC + color[val] + CE
				case "so": // socket
					colorsMap["socket"] = CC + color[val] + CE
				case "bd": // black device
					colorsMap["blak"] = CC + color[val] + CE
				case "cd": // char device
					colorsMap["charDevice"] = CC + color[val] + CE
					//                 case "mi": // Missing·file:·undefined
					//                     fmt.Printf("%s (missing file): %s\n", color[key], color[val])
					//                 case "or": // Or%saned·symlink:·undefined
					//                     fmt.Printf("%s (orphan ln)   : %s\n", color[key], color[val])
					//                 case "ex": // Executable
					//                     fmt.Printf("%s (exe)         : %s\n", color[key], color[val])
					//                 case "do": // Dooor
					//                     fmt.Printf("%s (door)        : %s\n", color[key], color[val])
					//                 case "su": // suid
					//                     fmt.Printf("%s (suid)        : %s\n", color[key], color[val])
					//                 case "sg": // guid
					//                     fmt.Printf("%s (guid)        : %s\n", color[key], color[val])
					//                 case "st": // sticky
					//                     fmt.Printf("%s (sitcky)      : %s\n", color[key], color[val])
					//                 case "ow": // other writable
					//                     fmt.Printf("%s (other writable): %s\n", color[key], color[val])
					//                 case "tw": // ow w/·sticky
					//                     fmt.Printf("%s (ow w/ sitcky): %s\n", color[key], color[val])
					//                 case "ca": // disabled·by·default
					//                     fmt.Printf("%s (disable??)   : %s\n", color[key], color[val])
					//                 case "mh": // disabled·by·default
					//                     fmt.Printf("%s (disable??)    : %s\n", color[key], color[val])
					//                 case "cl": // clear·to·end·of·line
					//                     fmt.Printf("%s (clear EOL)   : %s\n", color[key], color[val])

				default:
					continue
					//                                         fmt.Printf("Not handled: %s => %s\n", color[key], color[val])
				}

				//                 }
				//                 fmt.Printf("len: %d, key: %s, color: %s\n", len(key_color), key_color[0], key_color[1])
			}

		}

	} else {
		fmt.Printf("Do our own...")
	}
}
