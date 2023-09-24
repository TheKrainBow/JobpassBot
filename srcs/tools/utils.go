package tools

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func ToPtr[K any](input K) *K {
	return &input
}

func LogDeleteLastNLines(n int) {
	for i := 0; i < n; i++ {
		fmt.Printf("\033[A\r\033[K")
	}
}

func CreateTextBox(str string, boxName string) string {
	strs := strings.Split(str, "\n")

	maxStrLen := 0
	for i, str := range strs {
		if len(str)+2 > maxStrLen {
			maxStrLen = len(str) + 2
		}
		strs[i] = str
	}
	if len(boxName)+2 > maxStrLen {
		maxStrLen = len(boxName) + 4
	}
	stringToPrint := "``` ┏━"
	i := 0
	if boxName != "" {
		stringToPrint += " " + boxName + " "
		i += len(boxName) + 2
	}
	for ; i < maxStrLen-1; i++ {
		stringToPrint += "━"
	}
	stringToPrint += "┓\n"
	for _, str := range strs {
		if str == "" {
			continue
		}
		stringToPrint += fmt.Sprintf(" ┃ %s%*s┃\n", str, maxStrLen-len(str)-1, " ")
	}
	color.Set(color.Reset)
	stringToPrint += " ┗"
	for i = 0; i < maxStrLen; i++ {
		stringToPrint += "━"
	}
	stringToPrint += "┛```"
	return stringToPrint
}
