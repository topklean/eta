package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
)

type dirEntryInfo struct {
	name       string
	absName    string
	dirName    string
	mode       string
	typeEntr   rune
	inode      string
	hardLink   int
	softLink   int
	osIsDir    bool
	osFileType fs.FileMode // onlyne type of file / ModeType =
	//                            ModeDir |
	//                            ModeSymlink |
	//                            ModeNamedPipe |
	//                            ModeSocket |
	//                            ModeDevice |
	//                            ModeCharDevice |
	//                            ModeIrregular
	osFileInfo fs.FileInfo // have a mode to. but include rwx
	// mountPoint string  // maybe in the futur
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
			path, _ := filepath.Abs(arg)
			direntry, err := os.Lstat(arg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			dirEntries = AppendDirEntry(dirEntries, dirEntryInfo{
				name:       direntry.Name(),
				dirName:    path,
				absName:    path + "/" + direntry.Name(),
				typeEntr:   '-',
				osFileInfo: direntry,
			})
			continue
		}
		// rebuild the full path
		// dot dotdot Dir
		if conf.dotFile && !conf.dotDir {
			// arg is a file not dir
			path, _ := filepath.Abs(".")
			direntry, err := os.Lstat(".")
			if err != nil {
				fmt.Println(err)
				continue
			}
			dirEntries = AppendDirEntry(dirEntries, dirEntryInfo{
				name:       direntry.Name(),
				dirName:    path,
				absName:    path + "/" + direntry.Name(),
				typeEntr:   '-',
				osFileInfo: direntry,
			})
			path, _ = filepath.Abs("..")
			direntry, err = os.Lstat("..")
			if err != nil {
				fmt.Println(err)
				continue
			}
			dirEntries = AppendDirEntry(dirEntries, dirEntryInfo{
				name:       direntry.Name(),
				dirName:    path,
				absName:    path + "/" + direntry.Name(),
				typeEntr:   '-',
				osFileInfo: direntry,
			})
			//             elms = append(, elms...)
		}

		for _, direntry := range direntries {
			fileInfo, err := direntry.Info()
			if err != nil {
				fmt.Println(err)
				continue
			}
			path, _ := filepath.Abs(arg)
			dirEntries = AppendDirEntry(dirEntries, dirEntryInfo{
				name:       direntry.Name(),
				dirName:    path,
				absName:    path + "/" + direntry.Name(),
				typeEntr:   '-',
				osFileType: direntry.Type(),
				osFileInfo: fileInfo,
			})
		}
	}

	if conf.debug {
		fmt.Println("==========")
		spew.Dump(dirEntries)
		fmt.Println("==========")
	}

	for i := range dirEntries {

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

		switch {

		// link
		case dirEntries[i].osFileType&fs.ModeSymlink != 0:
			dirEntries[i].typeEntr = 'l'
			// add the target
			target, _ := os.Readlink(dirEntries[i].osFileInfo.Name())
			dirEntries[i].name += " -> " + target

		// dir ?
		case dirEntries[i].osFileInfo.Mode()&fs.ModeDir != 0:
			dirEntries[i].typeEntr = 'd'

		// Charactere Device
		case dirEntries[i].osFileInfo.Mode()&fs.ModeCharDevice != 0:
			dirEntries[i].typeEntr = 'c'

		// Device
		case dirEntries[i].osFileInfo.Mode()&fs.ModeDevice != 0:
			dirEntries[i].typeEntr = 'b'

		// Pipe
		case dirEntries[i].osFileInfo.Mode()&fs.ModeNamedPipe != 0:
			dirEntries[i].typeEntr = 'p'

		// Socket
		case dirEntries[i].osFileInfo.Mode()&fs.ModeSocket != 0:
			dirEntries[i].typeEntr = 's'
		}

		// Append ?
		if dirEntries[i].osFileInfo.Mode()&fs.ModeAppend != 0 {
			dirEntries[i].typeEntr = 'a'
		}
	}
}

func printListFiles() {

	fmt.Printf("Mod\t\tSize\tTime\t\tName\n")
	for i := range dirEntries {
		switch {
		case conf.dotFile:
		case !conf.dotDir && dirEntries[i].name[0] == '.':
			continue
		}

		// get mod from osFile
		// => [ T u g t rwx rwx rwx ]
		mode := []rune(fmt.Sprintf("%s", dirEntries[i].osFileInfo.Mode()))
		// get only the last 9 char rwxwxrwx
		mode_tmp := mode[:len(mode)-9]
		mode = mode[len(mode_tmp):]
		//         fmt.Printf("%s\n", string(mode))

		if dirEntries[i].osFileInfo.Mode()&fs.ModeSetuid != 0 {
			if mode[2] == 'x' {
				mode[2] = 's'
			} else {
				mode[2] = 'S'
			}
		}

		if dirEntries[i].osFileInfo.Mode()&fs.ModeSetgid != 0 {
			if mode[5] == 'x' {
				mode[5] = 's'
			} else {
				mode[5] = 'S'
			}
		}
		//         spew.Dump(mode)

		if dirEntries[i].osFileInfo.Mode()&fs.ModeSticky != 0 {
			if mode[8] == 'x' {
				mode[8] = 't'
			} else {
				mode[8] = 'T'
			}
		}
		//         spew.Dump(mode)

		fmt.Printf(
			"%c|%3s|%3s|%3s %9d\t%v\t%s\n",
			dirEntries[i].typeEntr,
			string(mode[0:3]),
			string(mode[3:6]),
			string(mode[6:]),
			dirEntries[i].osFileInfo.Size(),
			dirEntries[i].osFileInfo.ModTime().Format("Jan 02 15:04"),
			dirEntries[i].name,
		)
	}
	return
}
