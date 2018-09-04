package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tealeg/xlsx"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: xlsx pathname sheetname")
		os.Exit(1)
	}

	xlsxFile, err := xlsx.OpenFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	sheet := xlsxFile.Sheet[os.Args[2]]
	if sheet == nil {
		fmt.Println("表单名不存在")
		os.Exit(1)
	}

	for {
		title := getStdinInput("请输入列名：")
		if title == "" {
			fmt.Println(title)
			continue
		}

		titleColIndex := findColByTitle(sheet, title)
		if titleColIndex == -1 {
			fmt.Println("列名不存在")
			continue
		}

		rowLen := len(sheet.Rows)
		result := []string{}
		for rowIndex := 1; rowIndex < rowLen; rowIndex++ {
			content := sheet.Cell(rowIndex, titleColIndex).String()
			result = append(result, content)
		}

		fmt.Println(result)
	}
}

func getStdinInput(hint string) string {
	fmt.Print(hint)
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func findColByTitle(sheet *xlsx.Sheet, title string) int {
	titleRow := sheet.Rows[0]
	for titleIndex, col := range titleRow.Cells {
		if col.String() == title {
			return titleIndex
		}
	}

	return -1
}
