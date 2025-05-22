package main

// on fait quoi aujourd'hui ???
// ajout de logs (slog)
// ajout prise en charge des args ligne de commandes
// on ajoute les inodes
// formatage => couleurs
// organiser les sources (ajout de fonctions...)

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// struct fileInfo

//     type sFileInfo struct {
//         name string // pour lien, concatenation avec cyble
//                 attributs [10]rune
//         userName            string
//         groupname           string
//         userID int
//         groupID int
//         mode                 [10]rune
//         dateLastModification time.Time
//         size                 int

//                 link      string //
//                 inode     int16
//     }

// current dir

// init prog
// handle les args
// i
// func (receiver *type) name(para) err{}
// var = func(var)
type dirEntryInfo struct {
	name       string
	absName    string
	dirName    string
	mode       [10]rune
	userid     int
	groupid    int
	userName   string
	groupName  string
	btime      time.Time // birth (creation) time
	atime      time.Time // last acess time
	mtime      time.Time // last modification time
	ctime      time.Time // last status/metadata change
	size       int
	typeEntr   string
	inode      string
	hardLink   int
	softLink   int
	osIsDir    bool
	osFileMod  fs.FileMode
	osFileInfo fs.FileInfo
	//	dirEntry  os.DirEntry
	//
	// mountPoint string  // maybe
}

// contrainte du go: obligé de passer par def type pour pointer sur slice pour receiver de methode
type sliceDirEntries []dirEntryInfo

var dirEntries sliceDirEntries

// var dirEntries mde

// func (dirEntries *sliceDirEntries) add(dirEntry dirEntryInfo) {
//     *dirEntries = append(*dirEntries, dirEntry)
// }

func init() {
	//	dirEntries = append(dirEntries, dirEntryInfo{name: "toto"})
	//
	// my files (including dir)
}

func listDir() {

	// if no args, list current dir
	if len(args) == 0 {
		args = append(args, conf.cwd)
	}

	if conf.debug {
		fmt.Println("listDir: ")
		spew.Dump(args)
	}

	// get all entries
	//     var elms []os.DirEntry

	for _, arg := range args {
		// the names returned do not include the full path,
		// / /home /opt
		tmp, err := os.ReadDir(arg)
		if err != nil {
			// fichier
			// path
			// repertoire courant
			// chemin relatif ou abs
			//             traiter les fichiers
			fmt.Println(err)
			continue
		} else {
			//répertoires
		}

		// rebuild the full path
		for _, n := range tmp {
			fileInfo, err := n.Info()
			if err != nil {
				fmt.Println(err)
				continue
			}
			//             }

			dirEntries = append(dirEntries, dirEntryInfo{
				name:       n.Name(),
				dirName:    arg,
				absName:    arg + "/" + n.Name(),
				osFileMod:  n.Type(),
				osFileInfo: fileInfo,
			})
		}
	}

	if conf.debug {
		fmt.Println("==========")
		spew.Dump(dirEntries)
		fmt.Println("==========")
	}

	// dot dotdot Dir
	//     if conf.dotFile && !conf.dotDir {
	//         var dot, dotdot fs.DirEntry

	//         info, err := os.Lstat(".")
	//         if err != nil {
	//             fmt.Println("Cannot stat .")
	//         } else {
	//             dot = fs.FileInfoToDirEntry(info)
	//         }
	//         info, err = os.Stat("..")
	//         if err != nil {
	//             fmt.Println("Cannot stat .")
	//         } else {
	//             dot = fs.FileInfoToDirEntry(info)
	//         }
	//         elms = append([]fs.DirEntry{dot, dotdot}, elms...)
	//     }

	fmt.Printf("Mod\t\tSize\tTime\t\tName\n")
	for _, elm := range dirEntries {
		//         elmType := "None"
		if err := os.Chdir(elm.dirName); err != nil {
			fmt.Println(err)
			continue
		}

		if conf.debug == true {
			path, _ := os.Getwd()
			fmt.Println(path)
		}

		fileInfo, err := os.Lstat(elm.name)
		//         fileInfo, err := os.Lstat(elm.Name())
		//         if conf.debug {
		//             fmt.Println("file info: " + elm.Name())
		//             spew.Dump(fileInfo)
		//         }
		if err != nil {
			fmt.Println("Cannot stat: " + elm.name)
			continue
		}
		//         if fileInfo.Mode & os.ModePerm
		//         fileInfo.

		// setuid ?
		// - --- --- ---
		// - type (link, dir, file, ...)
		// --- user  mod
		// --- group mod
		// --- other mod

		// type ur uw ux gr gw gx or ow ox
		// 0	1	2 3  4  5  6  7  8  9
		//         var mode [10]rune

		mode := make([]rune, 10)
		for i := range mode {
			mode[i] = '-'
		}

		// iota for incremental int. first = 0
		const (
			element_type = iota
			user_read
			user_write
			user_execute
			group_read
			group_write
			group_execute
			other_read
			other_write
			other_execute
		)
		//         spew.Dump(
		for i, v := range fmt.Sprintf("%s", fileInfo.Mode()&fs.ModePerm) {
			if i == 0 {
				continue
			}
			mode[i] = v
		}

		fileName := fileInfo.Name()

		switch {

		// link
		case elm.osFileMod&fs.ModeSymlink != 0:
			mode[element_type] = 'l'
			// add the target
			v, _ := os.Readlink(fileInfo.Name())
			fileName += " -> " + v

			// dir ?
		case fileInfo.Mode()&fs.ModeDir != 0:
			mode[element_type] = 'd'

			// Charactere Device
		case fileInfo.Mode()&fs.ModeCharDevice != 0:
			mode[element_type] = 'c'

			// Device
		case fileInfo.Mode()&fs.ModeDevice != 0:
			mode[element_type] = 'b'

			// Pipe
		case fileInfo.Mode()&fs.ModeNamedPipe != 0:
			mode[element_type] = 'p'

			// Socket
		case fileInfo.Mode()&fs.ModeSocket != 0:
			mode[element_type] = 's'

		}

		// Append ?
		if fileInfo.Mode()&fs.ModeAppend != 0 {
			mode[element_type] = 'a'
		}

		if fileInfo.Mode()&fs.ModeSetuid != 0 {
			mode[user_execute] = 's'
		}

		if fileInfo.Mode()&fs.ModeSetgid != 0 {
			mode[group_execute] = 's'
		}

		if fileInfo.Mode()&fs.ModeSticky != 0 {
			mode[other_execute] = 'T'
		}

		//         fmt.Printf("%s\n", string(mode[:]))

		fmt.Printf(
			"%1s|%3s|%3s|%3s %9d\t%v\t%s\n",
			string(mode[:1]),
			string(mode[1:4]),
			string(mode[4:7]),
			string(mode[7:]),
			//             fileInfo.IsDir(),
			//             elm.Type(),
			fileInfo.Size(),
			fileInfo.ModTime().Format("06/01/02 15:04"),
			fileName,
		)
	}
}

func printListFiles() {

	return
}
