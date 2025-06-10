/*
file.go

	relative do files
*/
package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"

	"github.com/davecgh/go-spew/spew"
)

type dirEntryInfo struct {
	name           string
	absName        string
	dirName        string
	nameLen        int
	displayName    string
	displayNameLen int
	mode           string
	userName       string
	groupName      string
	typeEntr       rune
	inode          string
	hardLink       uint64
	softLink       int
	osIsDir        bool
	osFileInfo     fs.FileInfo // have a mode to. but include rwx
}

// contrainte du go: obligé de passer par def type pour pointer sur slice pour receiver de methode
type sliceDirEntries []dirEntryInfo

var dirEntries sliceDirEntries

func (dirEntries *sliceDirEntries) add(fileInfo fs.FileInfo, arg string) {

	// we need to rebuild the path for sym link
	path, _ := filepath.Abs(arg)

	// stat the file for uid/gid/hard link/size...
	sys, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		fmt.Printf("Cannot stat file %s ...", arg)
		return
	}

	//     userUid := fmt.Sprintf("%d", sys.Uid)
	//     var userName string
	//     user_name, err := user.LookupId(userUid)
	//     if err != nil {
	//         userName = userUid
	//     } else {
	//         userName = user_name.Username
	//     }

	//     groupGid := int(sys.Gid)
	//     var groupName string
	//     group_name, err := LookupId(groupGid) // opti withh internal cache
	//     // so slow :(
	//     //             group_name, err := user.LookupGroupId(groupGid)
	//     if err != nil {
	//         groupName = fmt.Sprintf("%d", groupGid)
	//     } else {
	//         groupName = group_name.Name
	//     }

	hard_link := sys.Nlink

	// type of entry
	typeEntrie := '-'
	// target of link
	var target string
	// get mod of entry
	mode := fileInfo.Mode()
	//
	switch {

	// link
	case mode&fs.ModeSymlink != 0:
		typeEntrie = 'l'
		// add the target
		target, err := os.Readlink(path)
		if err != nil {
			target = "-> target not found..."
		} else {
			target = "->" + target
		}

		// dir ?
	case mode&fs.ModeDir != 0:
		typeEntrie = 'd'

	// Charactere Device
	case mode&fs.ModeCharDevice != 0:
		typeEntrie = 'c'
		// get major number

	// Device
	case mode&fs.ModeDevice != 0:
		typeEntrie = 'b'

	// Pipe
	case mode&fs.ModeNamedPipe != 0:
		typeEntrie = 'p'

	// Socket
	case mode&fs.ModeSocket != 0:
		typeEntrie = 's'
	}

	// Append ?
	if mode&fs.ModeAppend != 0 {
		typeEntrie = 'a'
	}

	// keep filename max len (for futur formating printing)
	fileName := fileInfo.Name()

	namelen := len(fileName)

	dirEntry := dirEntryInfo{
		name:     fileName + target,
		dirName:  path,
		absName:  path + "/" + fileName,
		nameLen:  namelen,
		typeEntr: typeEntrie,
		//         userName:   userName,
		//         groupName:  groupName,
		hardLink:   hard_link,
		osFileInfo: fileInfo,
	}
	dirEntry.displayName = dirEntry.getNameAndFrills()
	dirEntry.displayNameLen = len(dirEntry.displayName)

	*dirEntries = AppendDirEntry(*dirEntries, dirEntry)
}

// func (dirEntry *dirEntryInfo) quoteName() string {
//     // TBD
//     // if name file contain:
//     // espace ! " # $ % & ' () * + , - . : ; < = > ? / } @ [ \ ] ^ _ ` { | } ~ DEL
//     // the firts 31 char of ascii table, have to be escaped or displayed in plain text

// }

func (dirEntry *dirEntryInfo) getMode() []rune {
	mode := []rune(fmt.Sprintf("%s", (*dirEntry).osFileInfo.Mode()))
	mode_tmp := mode[:len(mode)-9]
	mode = mode[len(mode_tmp):]
	return mode
}

func (dirEntry *dirEntryInfo) getNameAndFrills() string {
	// file indictor
	//     var Indicator = make(map[string]string)
	name := (*dirEntry).name

	// do we have to quote the file ???
	typeEntry := (*dirEntry).typeEntr
	if conf.indicator {
		switch {
		case typeEntry == 'd':
			name += indicatorDir
		case typeEntry == 'l':
			name += indicatorLink
		case typeEntry == 'p':
			name += indicatorPipe
		case typeEntry == 's':
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
	return name
}

// func (dirEntry dirEntryInfo) String() string {
//     switch conf.format {

//     case "long":
//         // get mod from osFile
//         // => [ T u g t rwx rwx rwx ]
//         mode := []rune(fmt.Sprintf("%s", dirEntry.osFileInfo.Mode()))
//         // get only the last 9 char rwxwxrwx
//         mode_tmp := mode[:len(mode)-9]
//         mode = mode[len(mode_tmp):]

//         if dirEntry.osFileInfo.Mode()&fs.ModeSetuid != 0 {
//             if mode[2] == 'x' {
//                 mode[2] = 's'
//             } else {
//                 mode[2] = 'S'
//             }
//         }

//         if dirEntry.osFileInfo.Mode()&fs.ModeSetgid != 0 {
//             if mode[5] == 'x' {
//                 mode[5] = 's'
//             } else {
//                 mode[5] = 'S'
//             }
//         }
//         //         spew.Dump(mode)

//         if dirEntry.osFileInfo.Mode()&fs.ModeSticky != 0 {
//             if mode[8] == 'x' {
//                 mode[8] = 't'
//             } else {
//                 mode[8] = 'T'
//             }
//         }
//         //         spew.Dump(mode)
//         modsep := ""

//         return fmt.Sprintf(
//             "%c%s%3s%s%3s%s%3s  %2d  %10s %10s %9d\t%v\t%s\n",
//             //                                     "%c|%3s|%3s|%3s  %2d  %10s  %-10s %9d\t%v\t%s\n",
//             dirEntry.typeEntr,
//             modsep,
//             string(mode[0:3]),
//             modsep,
//             string(mode[3:6]),
//             modsep,
//             string(mode[6:]),
//             dirEntry.hardLink,
//             dirEntry.userName,
//             dirEntry.groupName,
//             dirEntry.osFileInfo.Size(),
//             dirEntry.osFileInfo.ModTime().Format("Jan 02 15:04"),
//             dirEntry.name,
//         )

//     case "one":
//         return fmt.Sprintf("%s\n", dirEntry.name)

//     //     default:
//     //         fmt.Println(conf.formatColRemaining)
//     //         if format.ColRemaining > 1 {
//     //             format.ColRemaining = format.ColRemaining - 1
//     //             return fmt.Sprintf(format.One+"│", dirEntry.name)
//     //         } else {
//     //             format.ColRemaining = format.MaxCol

//     //             return fmt.Sprintf(format.One+"│\n", dirEntry.name)
//     //         }
//     default:
//         return "toto"

//     }

//     // return
// }

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

		direntries, err := os.ReadDir(arg)

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
				dirEntries.add(direntry, arg)
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
					dirEntries.add(direntry, arg)
				}

			}
		}

		// elements from readDir
		for _, direntry := range direntries {

			fileInfo, err := direntry.Info()
			if err != nil {
				fmt.Println(err)
				continue
			}
			dirEntries.add(fileInfo, arg)
		}

	}

	if conf.debug {
		fmt.Println("==========")
		spew.Dump(dirEntries)
		fmt.Println("==========")
	}
}
