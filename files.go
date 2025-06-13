/*
file.go

	relative do files
	TBD: utf8.RuneCountInString(string)
*/
package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
)

type dirEntryInfo struct {
	name string
	//     absName        string
	//     dirName        string
	nameLen        int
	displayName    string
	displayNameLen int
	mode           string
	userName       string
	groupName      string
	typeEntry      rune
	inode          string
	hardLink       uint64
	softLink       int
	osIsDir        bool
	osFileInfo     fs.FileInfo // have a mode to. but include rwx
}

// contrainte du go: obligé de passer par def type pour pointer sur slice pour receiver de methode
type sliceDirEntries []dirEntryInfo

var dirEntries sliceDirEntries

func (dirEntry *dirEntryInfo) addType(fileInfo *fs.FileInfo, arg string) {

	mode := (*fileInfo).Mode()
	var typeEntry rune

	switch {

	// link
	case mode&fs.ModeSymlink != 0:
		typeEntry = 'l'
		//         only do get link target in long format
		//         name -> target
		//         TBD: dereference link
		if conf.format == "long" {
			// we need to rebuild the filePath for sym link
			filePath, _ := filepath.Abs(arg + "/" + (*dirEntry).name)
			target, err := os.Readlink(filePath)
			if err != nil {
				(*dirEntry).name += " -> " + fmt.Sprintf("%s", err)
			} else {
				(*dirEntry).name += " -> " + target
			}
		}

	// dir ?
	case mode&fs.ModeDir != 0:
		typeEntry = 'd'

	// Charactere Device
	case mode&fs.ModeCharDevice != 0:
		typeEntry = 'c'
		// TBD: get major/minor number

	// Device
	case mode&fs.ModeDevice != 0:
		typeEntry = 'b'
		// TBD: get major/minor number

	// Pipe
	case mode&fs.ModeNamedPipe != 0:
		typeEntry = 'p'

	// Socket
	case mode&fs.ModeSocket != 0:
		typeEntry = 's'

	default:
		typeEntry = '-'
	}

	// Append ?
	if mode&fs.ModeAppend != 0 {
		typeEntry = 'a'
	}

	(*dirEntry).typeEntry = typeEntry
	// return typeEntry
}

func (dirEntry *dirEntryInfo) addUserGroupName(uid, gid uint32) {

	uidAsString := fmt.Sprintf("%d", uid)
	user, err := user.LookupId(uidAsString)
	if err != nil {
		(*dirEntry).userName = uidAsString
	} else {
		(*dirEntry).userName = user.Username
	}

	gidAsInt := int(gid)
	group_name, err := LookupId(gidAsInt) // opti with internal cache
	if err != nil {
		(*dirEntry).groupName = fmt.Sprintf("%d", gidAsInt)
	} else {
		(*dirEntry).groupName = group_name.Name
	}

}

func (dirEntries *sliceDirEntries) add(fileInfo *fs.FileInfo, arg string) {

	// stat the file for uid/gid/hard link/size...
	sys, ok := (*fileInfo).Sys().(*syscall.Stat_t)

	if !ok {
		fmt.Printf("Cannot stat file %s ...", arg)
		return
	}

	// element of dir (an entry)
	// just
	dirEntry := dirEntryInfo{
		name:       (*fileInfo).Name(),
		hardLink:   sys.Nlink,
		osFileInfo: (*fileInfo),
	}

	// len of file name
	dirEntry.nameLen = len(dirEntry.name)

	// add type link, dir, file, block, ...
	// if link and format long, add target in name file
	dirEntry.addType(fileInfo, arg)

	// user and group Name
	// only in long format
	//  TBD: add in other format (config flag)
	if conf.format == "long" {
		dirEntry.addUserGroupName(sys.Uid, sys.Gid)
	}

	dirEntry.setNameToDisplay()
	//     dirEntry.displayNameLen = len(dirEntry.displayName)

	// finaly, add entry to array
	*dirEntries = AppendDirEntry(*dirEntries, dirEntry)
}

//	var special []rune = []rune{
//	    ' ', '(', ')', '!', '"', '#', '$', '&', '\'', '*', '+', ',', '/', ':', ';', '<', '>', '?', '\\', '^', '`', '|', '~',
//	}
var special []rune = []rune{' '}

func (dirEntry *dirEntryInfo) quoteName() bool {
	// TBD
	// if name file contain:
	// espace ! " # $ % & ' () * + , - . : ; < = > ? / } @ [ \ ] ^ _ ` { | } ~ DEL
	// the firts 31 char of ascii table, have to be escaped or displayed in plain text
	//     for i, c := range dirEntry.name {
	//         if
	//         fmt.Printf("i: %d, c:%c", i, c)
	//     }
	//     _ = special
	// TBD: remove the string cast... because called in loop for every files... no good
	if strings.ContainsAny(dirEntry.name, string(special)) {
		//         fmt.Printf("to quote: %s\n", (*dirEntry).name)
		return true
	}

	return false
}

func (dirEntry *dirEntryInfo) getMode() []rune {
	mode := []rune(fmt.Sprintf("%s", (*dirEntry).osFileInfo.Mode()))
	mode_tmp := mode[:len(mode)-9]
	mode = mode[len(mode_tmp):]
	return mode
}

// func (dirEntry *dirEntryInfo) setNameColor() string {
//     name := (*dirEntry).name
// }

func (dirEntry *dirEntryInfo) setNameToDisplay() {

	// TBD: quoting must be done first
	//      Indicator must be done last
	//      Colors must be done on the name without indicator
	//		the len of name do not take the colors codes!!!

	name := (*dirEntry).name

	if (*dirEntry).quoteName() {
		name = "'" + name + "'"
	}
	// TBD : colors must be donne afte column calculation !!!
	//     var colorsLenChar int = 0

	colorsLenChar := 0
	//     fmt.Printf("Colors: %v", conf.colorsEnable)
	if conf.colorsEnable {
		switch {

		case (*dirEntry).typeEntry == 'd':
			name = colorsMap["dir"] + name + colorsMap["reset"]
			colorsLenChar = len(colorsMap["dir"])

		case (*dirEntry).typeEntry == 'l':
			name = colorsMap["link"] + name + colorsMap["reset"]
			colorsLenChar = len(colorsMap["link"])

		case (*dirEntry).typeEntry == 'p':
			name = colorsMap["pipe"] + name + colorsMap["reset"]
			colorsLenChar = len(colorsMap["pipe"])

		case (*dirEntry).typeEntry == 's':
			name = colorsMap["socket"] + name + colorsMap["reset"]
			colorsLenChar = len(colorsMap["socket"])

		case (*dirEntry).typeEntry == 'c':
			name = colorsMap["charDevice"] + name + colorsMap["reset"]
			colorsLenChar = len(colorsMap["charDevice"])

			//         case (*dirEntry).typeEntry == 's':
			//             name = colorsMap["socket"] + name + colorsMap["reset"]
			//             colorsLenChar = len(colorsMap["socket"])
			//         case (*dirEntry).typeEntry == 's':
			//             name = colorsMap["socket"] + name + colorsMap["reset"]
			//             colorsLenChar = len(colorsMap["socket"])
		default:
			//             name = colorsMap["default"] + name + colorsMap["reset"]
			//             name = colorsMap["default"] + name + colorsMap["reset"]
			colorsLenChar = 0
		}
		colorsLenChar += len(colorsMap["reset"])
	}

	if conf.indicator {
		switch {
		case (*dirEntry).typeEntry == 'd':
			//             if conf.colorsEnable {
			name += indicatorDir
			//             }
		case (*dirEntry).typeEntry == 'l':
			name += indicatorLink
		case (*dirEntry).typeEntry == 'p':
			name += indicatorPipe
		case (*dirEntry).typeEntry == 's':
			name += indicatorSock
		default:
			// exec
			// rwx rwx rwx
			//   2   5   8
			//             fmt.Printf("%c\n", mode[3])
			mode := (*dirEntry).getMode()
			if mode[2] == 'x' || mode[5] == 'x' || mode[8] == 'x' {
				name += indicatorExe
			}
		}
	}
	(*dirEntry).displayName = name
	(*dirEntry).displayNameLen = len(name) - colorsLenChar
}

func listDir() {

	// if no args, work on current directory
	if len(args) == 0 {
		args = AppendString(args, conf.cwd)
	}

	// debug: dump args with spew
	if conf.debug {
		fmt.Println("listDir: ")
		spew.Dump(args)
	}

	// get all entries
	for _, arg := range args {

		// the names returned do not include the full path,
		// arg must be a directory

		// args : fichiers
		//			rep
		//			fichiers / args

		//         listbegin := time.Now().UnixMilli()
		direntries, err := os.ReadDir(arg)
		//         listend := time.Now().UnixMilli()
		//         fmt.Printf("args: %s, ReadDir: %d\n", arg, listend-listbegin)

		if err != nil {
			// arg is a file not dir
			// stat and add it if no error
			// we do not panic, continue en next arg if we got an error
			direntry, err := os.Lstat(arg)
			if err != nil {
				fmt.Println(err)
				// go to next arg
				continue
			} else {
				dirEntries.add(&direntry, arg)
			}
			continue
		}

		// dot(.) and otdot(..) Dir // -a or -A (not both)
		// TBD only add if args is dir
		// and ad it by dir
		if conf.dotFile && !conf.dotDir {

			for _, d := range []string{".", ".."} {

				direntry, err := os.Lstat(d)
				if err != nil {
					fmt.Println(err)
					continue
				} else {
					dirEntries.add(&direntry, arg)
				}

			}
		}

		// elements from readDir
		//         b := time.Now().UnixMilli()
		//         e := time.Now().UnixMilli()
		for _, direntry := range direntries {

			//             listbegin = time.Now().UnixMilli()
			fileInfo, err := direntry.Info()
			//             _ = fileInfo
			if err != nil {
				fmt.Println(err)
				continue
			}
			dirEntries.add(&fileInfo, arg)
			//             listend = time.Now().UnixMilli()

			//             fmt.Printf("add: %s, ReadDir: %d\n", fileInfo.Name(), listend-listbegin)
		}
		//         e = time.Now().UnixMilli()
		//         fmt.Printf("all: %d\n", e-b)

	}

	if conf.debug {
		fmt.Println("==========")
		spew.Dump(dirEntries)
		fmt.Println("==========")
	}
}
