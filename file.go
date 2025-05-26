package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"syscall"

	"github.com/davecgh/go-spew/spew"
)

type dirEntryInfo struct {
	name       string
	absName    string
	dirName    string
	mode       string
	userName   string
	groupName  string
	typeEntr   rune
	inode      string
	hardLink   uint64
	softLink   int
	osIsDir    bool
	osFileInfo fs.FileInfo // have a mode to. but include rwx
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

	userUid := fmt.Sprintf("%d", sys.Uid)
	var userName string
	user_name, err := user.LookupId(userUid)
	if err != nil {
		userName = userUid
	} else {
		userName = user_name.Username
	}

	groupGid := int(sys.Gid)
	var groupName string
	group_name, err := LookupId(groupGid) // opti withh internal cache
	// so slow :(
	//             group_name, err := user.LookupGroupId(groupGid)
	if err != nil {
		groupName = fmt.Sprintf("%d", groupGid)
	} else {
		groupName = group_name.Name
	}

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

	//     var := len
	*dirEntries = AppendDirEntry(*dirEntries, dirEntryInfo{
		name:       fileInfo.Name() + target,
		dirName:    path,
		absName:    path + "/" + fileInfo.Name(),
		typeEntr:   typeEntrie,
		userName:   userName,
		groupName:  groupName,
		hardLink:   hard_link,
		osFileInfo: fileInfo,
	})
}

func (dirEntry dirEntryInfo) String() string {
	switch conf.format {

	case "long":
		// get mod from osFile
		// => [ T u g t rwx rwx rwx ]
		mode := []rune(fmt.Sprintf("%s", dirEntry.osFileInfo.Mode()))
		// get only the last 9 char rwxwxrwx
		mode_tmp := mode[:len(mode)-9]
		mode = mode[len(mode_tmp):]

		if dirEntry.osFileInfo.Mode()&fs.ModeSetuid != 0 {
			if mode[2] == 'x' {
				mode[2] = 's'
			} else {
				mode[2] = 'S'
			}
		}

		if dirEntry.osFileInfo.Mode()&fs.ModeSetgid != 0 {
			if mode[5] == 'x' {
				mode[5] = 's'
			} else {
				mode[5] = 'S'
			}
		}
		//         spew.Dump(mode)

		if dirEntry.osFileInfo.Mode()&fs.ModeSticky != 0 {
			if mode[8] == 'x' {
				mode[8] = 't'
			} else {
				mode[8] = 'T'
			}
		}
		//         spew.Dump(mode)
		modsep := ""
		return fmt.Sprintf(
			"%c%s%3s%s%3s%s%3s  %2d  %10s %10s %9d\t%v\t%s\n",
			//                 "%c|%3s|%3s|%3s  %2d  %10s  %-10s %9d\t%v\t%s\n",
			dirEntry.typeEntr,
			modsep,
			string(mode[0:3]),
			modsep,
			string(mode[3:6]),
			modsep,
			string(mode[6:]),
			dirEntry.hardLink,
			dirEntry.userName,
			dirEntry.groupName,
			dirEntry.osFileInfo.Size(),
			dirEntry.osFileInfo.ModTime().Format("Jan 02 15:04"),
			dirEntry.name,
		)

	case "one":
		return fmt.Sprintf("%s\n", dirEntry.name)
	default:
		return fmt.Sprintf("%s  ", dirEntry.name)
	}

	// return
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

		//         spew.Dump(args)
		direntries, err := os.ReadDir(arg)
		//         fmt.Println("ReadDir done")
		if err != nil {
			// arg is a file not dir
			// stat and add it if no error
			// we do not panic, continue en next arg if we got an error
			direntry, err := os.Lstat(arg)
			if err != nil {
				fmt.Println(err)
			} else {
				dirEntries.add(direntry, arg)
			}
			continue
		}

		// dot(.) and otdot(..) Dir // -a or -A (not both)
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

	// Check type for dir entrie
	// only needed in format :
	// + color
	// + -l --long

}

func printListFiles() {

	for i := range dirEntries {

		// -a -A (hidden file and dot dotdot dir entries)
		//         switch {
		//         case conf.dotFile:
		//         case !conf.dotDir && dirEntries[i].name[0] == '.':
		//                 case !conf.dotFile && (!conf.dotDir && dirEntries[i].name[0] == '.'):
		//             continue
		//         }
		if !conf.dotFile && !conf.dotDir && dirEntries[i].name[0] == '.' {
			continue
		}

		/*
			nombre de fichiers
			nombre de colonnes
			ex : 1000
			for I := range dirEntries { if len(dirEntry[i]) > maxDirEntryLen) { maxDirEntryLen = len(dirEntry[i]) } }
			long collon / plus long => nobr colonne

		*/

		//     switch conf.format {
		//         case ""
		//     }
		pf := "%s"
		fmt.Printf(pf, dirEntries[i])
	}
	return
}
