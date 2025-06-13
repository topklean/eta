package main

import (
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

type columnLayoutInfo struct {
	validLen  bool
	lineLen   int
	colsArray []int
}

var columnsLayout []columnLayoutInfo
var columnInfoAllocated int = 0

// var ttyMaxCols int
var padding = 0

const MINCOLUMNWIDTH = 3

func getTTYMaxColumns() int {

	ttyLineLen := conf.ttySizeCol
	//     log.Printf("tty size col: %d\n", ttyLineLen)
	maxCols := ttyLineLen / MINCOLUMNWIDTH
	//     log.Printf("maxCols: %d\n", maxCols)

	return maxCols
}

func initInfoColumns(maxCols int) {
	//     log.Printf("maxCols: %d\n", maxCols)
	//     fmt.Scanf("%s")

	for i := range maxCols {
		//         log.Printf("init in col, i range: %d\n", i)

		var columnInfo columnLayoutInfo
		columnInfo.validLen = true
		columnInfo.lineLen = (i + 1) * MINCOLUMNWIDTH
		columnInfo.colsArray = make([]int, i+1)

		for j := range i + 1 {
			columnInfo.colsArray[j] = MINCOLUMNWIDTH
		}

		columnsLayout = append(columnsLayout, columnInfo)
	}

}

func (direntries *sliceDirEntries) getColumnsLayout(byColumn bool) int {

	ttyMaxCols := getTTYMaxColumns()
	if ttyMaxCols < 0 {
		return 1
	}
	//     log.Printf("ttyMAxCol: %d", ttyMaxCols)

	filesCount := len(*direntries)

	maxCols := 1
	//     maxCols := filesCount
	//     var maxCols int
	if 0 < ttyMaxCols && ttyMaxCols < filesCount {
		maxCols = ttyMaxCols
		//     } else {
		//         maxCols = filesCount
	}
	//     log.Printf(": %d", ttyMaxCols)

	//     log.Printf("before init col\n")
	initInfoColumns(maxCols)
	//     log.Printf("after init Info Colum, maxCols: %d\n", maxCols)

	//     log.Printf("before loop dir")
	for fileIndex, file := range *direntries {
		//         log.Printf("fileIndex: %d\n", fileIndex)
		nameLen := file.displayNameLen + padding

		for colIndex := range maxCols {
			var i int
			if columnsLayout[colIndex].validLen {
				if byColumn {
					i = fileIndex / ((filesCount + colIndex) / (colIndex + 1))
				} else {
					i = fileIndex % (colIndex + 1)
				}

				realLen := nameLen
				if i != colIndex {
					realLen = nameLen + 2
				}
				if columnsLayout[colIndex].colsArray[i] < realLen {
					columnsLayout[colIndex].lineLen += (realLen - columnsLayout[colIndex].colsArray[i])
					columnsLayout[colIndex].colsArray[i] = realLen
					columnsLayout[colIndex].validLen = (columnsLayout[colIndex].lineLen < conf.ttySizeCol)
				}
			}

		}

	}

	var cols int
	for cols = maxCols; 1 < cols; cols -= 1 {
		if columnsLayout[cols-1].validLen {
			break
		}
	}

	return cols
}

func (dirEntries *sliceDirEntries) printByColumn() {

	//     log.Printf("Before row loop, stdout type: %s\n", conf.tty.stdoutType)
	cols := dirEntries.getColumnsLayout(true)
	// TBD: adding error management for cols < 0
	if cols <= 0 {
		cols = 1
	}

	colLayout := columnsLayout[cols-1].colsArray

	// len list to display
	listCount := len(*dirEntries)
	// rows in each columns
	rows := listCount / cols
	// short column in the right
	if (listCount % cols) != 0 {
		rows += 1
	}

	var screenBuff strings.Builder

	for row := range rows {
		//         var str string
		col := 0
		filesno := row
		pos := 0
		//         var line string
		//         log.Printf("row: %d\n", row)
		for true {
			maxNameLen := colLayout[col]
			col++
			//             log.Printf("stdout: %s, Type: %s\n", conf.tty.stdoutType, TypeCharDevice)
			//             log.Println("(before tty check) in for true")
			//             runtime.Breakpoint()
			//             if conf.tty.stdoutType == TypeCharDevice {
			//                 log.Println("(screenBuf) in for true")
			screenBuff.WriteString((*dirEntries)[filesno].displayName)
			//             } else if conf.tty.stdoutType == TypeFile || conf.tty.stdoutType == TypePipe {
			//                 log.Println("(Sprintf) in for true")
			//                 line += fmt.Sprintf("%s", (*dirEntries)[filesno].displayName)
			//             }
			//             if err != nil {
			//                 log.Fatal("oups: " + err.Error())
			//             }
			p := maxNameLen - (*dirEntries)[filesno].displayNameLen
			//             log.Println("(before padding space) in for true")
			//             line += fmt.Sprintf("%s", (*dirEntries)[filesno].displayName)

			//             if conf.tty.stdoutType == TypeCharDevice {
			for i := 0; i < p; i++ {
				screenBuff.WriteRune(' ')
			}
			//             } else if conf.tty.stdoutType == TypeFile || conf.tty.stdoutType == TypePipe {
			//                 for i := 0; i < p; i++ {
			//                     line += " "
			//                 }
			if listCount-rows <= filesno {
				break
			}
			filesno += rows
			pos += maxNameLen
		}
		screenBuff.WriteString("\n")
		//             log.Println("(after padding space) in for true")
	}
	//         log.Println("after for true")
	//         if conf.tty.stdoutType == TypeCharDevice {
	//         } else if conf.tty.stdoutType == TypeFile || conf.tty.stdoutType == TypePipe {
	//             line += "\n"
	//         }

	//         if conf.tty.stdoutType == TypeFile || conf.tty.stdoutType == TypePipe {
	//             fmt.Println("coucou")
	//             log.Fatal("not in tty")
	//             spew.Dump(screenBuff)
	//             fmt.Printf("%s", line)
	//             fmt.Fprint(os.Stdout, screenBuff.String())
	//             fmt.Println("coucou")
	//             fmt.Println(screenBuff.String())
	//             screenBuff.Reset()
	//         }
	//     }
	// io.Copy(os.Stdout, byte(screenBuff)[])
	//             screenBuff.Reset()
	//     if conf.tty.stdoutType == TypeCharDevice {
	fmt.Println(screenBuff.String())
	//	    }
	//		screenBuff
	//
	// screenBuff.Reset()
	// fmt.Println(screenBuff.String())
	// fmt.Println(screenBuff.String())
}

func (dirEntries *sliceDirEntries) printLong() {
	for _, dirEntry := range *dirEntries {
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

		fmt.Printf(
			"%c%s%3s%s%3s%s%3s  %2d  %10s %10s %9d\t%v\t%s\n",
			//                                     "%c|%3s|%3s|%3s  %2d  %10s  %-10s %9d\t%v\t%s\n",
			dirEntry.typeEntry,
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
	}
}

func printListFiles() {
	// algo : don't ask me - I get it from coreutils ls.c

	sort.Slice(dirEntries, func(i, j int) bool {
		// Premier critère: comparer les types
		i_type := dirEntries[i].typeEntry
		j_type := dirEntries[j].typeEntry

		//         }

		// si l'un des deux est un rep
		if i_type == 'd' && j_type != 'd' {
			return true
		}
		if i_type != 'd' && j_type == 'd' {
			return false
		}

		// si les deux sont un reps ou si les deux ne le sont pas
		a := strings.ToLower(dirEntries[i].name)
		b := strings.ToLower(dirEntries[j].name)
		if a[0] == '.' {
			a = dirEntries[i].name[1:]
		}
		if b[0] == '.' {
			b = dirEntries[j].name[1:]
		}

		return a < b
	})

	//     log.Println("after sort: begin print by column")
	switch conf.format {
	case "long":
		dirEntries.printLong()

	default:
		dirEntries.printByColumn()
	}
	//         dirEntries.printByColumn()

}
