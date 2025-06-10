package main

import (
	"fmt"
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
	maxCols := ttyLineLen / MINCOLUMNWIDTH

	return maxCols
}

func initInfoColumns(maxCols int) {

	for i := range maxCols {

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

	filesCount := len(*direntries)

	maxCols := filesCount
	//     var maxCols int
	if 0 < ttyMaxCols && ttyMaxCols < filesCount {
		maxCols = ttyMaxCols
		//     } else {
		//         maxCols = filesCount
	}

	initInfoColumns(maxCols)

	for fileIndex, file := range *direntries {
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

// method on direntries
func (dirEntries *sliceDirEntries) printByColumn() {

	cols := dirEntries.getColumnsLayout(true)
	// TBD: adding error management for cols < 0
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

		for true {
			maxNameLen := colLayout[col]
			col++
			screenBuff.WriteString((*dirEntries)[filesno].displayName)
			p := maxNameLen - (*dirEntries)[filesno].displayNameLen
			for i := 0; i < p; i++ {
				screenBuff.WriteRune(' ')
			}
			if listCount-rows <= filesno {
				break
			}
			filesno += rows
			pos += maxNameLen
		}
		screenBuff.WriteString("\n")
	}
	// io.Copy(os.Stdout, byte(screenBuff)[])
	fmt.Println(screenBuff.String())
}

func printListFiles() {
	// algo : don't ask me - I get it from coreutils ls.c

	sort.Slice(dirEntries, func(i, j int) bool {
		// Premier critère: comparer les types
		i_type := dirEntries[i].typeEntr
		j_type := dirEntries[j].typeEntr

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
	dirEntries.printByColumn()
}
