package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/olekukonko/tablewriter"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	// create argparse
	parser := argparse.NewParser("CSV query language", "A command-line-interface for interacting with CSV files, such as getting and setting values and appending rows and columns. All get commands use base 1 indexing")

	// create argparse arguments
	file := parser.String("", "file", &argparse.Options{Help: "File to open"})
	printData := parser.Flag("", "print", &argparse.Options{Default:  false, Help: "Display the loaded CSV file as a table"})
	create := parser.String("", "create", &argparse.Options{Help: "Create a CSV file at the specified location"})
	getCell := parser.String("", "getcell", &argparse.Options{Help: "Get the value from a cell's x and y position, in the format of x|y"})
	getRow := parser.Int("", "getrow", &argparse.Options{Help: "Get all values in a row, given by it's x position"})
	getColumn := parser.String("", "getcolumn", &argparse.Options{Help: "Get all values in a column, given by it's name"})
	setCell := parser.String("", "setcell", &argparse.Options{Help: "Set the value of a cell, in the format of x|y|value. If value is more than one word then use \"\" around value"})
	appendColumn := parser.String("", "appendcolumn", &argparse.Options{Help: "Append a new column to the file. If more then one word then use \"\" around title"})
	appendRow := parser.Flag("", "appendrow", &argparse.Options{Default: false, Help: "Create an empty row that can be filled out with --setcell"})
	getNumC := parser.Flag("", "getnumcolumns", &argparse.Options{Help: "Get the number of columns"})
	getNumR := parser.Flag("", "getnumrows", &argparse.Options{Help: "Get the number of rows"})

	// run argparse
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	// check if creating file
	if *create != "" {
		_, _ = os.Create(*create)
	} else {
		// load data
		data := loadDataToFrame(*file)

		// if not creating file, check which arguments are given
		if *setCell != "" {
			x, _ := strconv.Atoi(strings.Split(*setCell, "|")[0])
			y, _ := strconv.Atoi(strings.Split(*setCell, "|")[1])
			value := strings.Split(*setCell, "|")[2]

			// remove old file and write new data
			newData := setCellValue([2]int{x, y}, data, value)
			_ = os.Remove(*file)
			f, _ := os.Create(*file)
			writer := csv.NewWriter(f)
			writer.WriteAll(newData)
		}
		if *appendColumn != "" {
			// remove old file and write new data
			_ = os.Remove(*file)
			f, _ := os.Create(*file)
			writer := csv.NewWriter(f)
			writer.WriteAll(appendNewColumn(*appendColumn, data))
		}
		if *appendRow {
			commas := []string{}
			for i := 0; i < len(data[0]); i++ {
				commas = append(commas, "")
			}
			data = append(data, commas)

			// remove old file and write new data
			_ = os.Remove(*file)
			f, _ := os.Create(*file)
			writer := csv.NewWriter(f)
			writer.WriteAll(data)
		}
		if *getColumn != "" {
			getColumnValues(*getColumn, data)
		}
		if *getRow != 0 {
			getRowValues(*getRow, *file)
		}
		if *getCell != "" {
			x, _ := strconv.Atoi(strings.Split(*getCell, "|")[0])
			y, _ := strconv.Atoi(strings.Split(*getCell, "|")[1])
			getCellValue([2]int{x, y}, data)
		}
		if *printData == true {
			printTable(data)
		}
		if *getNumC == true {
			getNumColumns(data)
		}
		if *getNumR == true {
			getNumRows(data)
		}
	}
}

func loadDataToFrame(filePath string) [][]string {
	// open file and create csv reader
	file, _ := os.Open(filePath)
	defer file.Close()
	r := csv.NewReader(file)

	// read through csv lines
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(records) == 0 {
		records = append(records, []string{})
	}

	return records
}

func getNumColumns(t [][]string) {
	fmt.Println(len(t[0]))
}

func getNumRows(t [][]string) {
	fmt.Println(len(t))
}

func printTable(t [][]string) {
	// get table headers
	headers := t[0]

	// create table writer and set headers
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for i := 1; i < len(t); i++ {
		table.Append(t[i])
	}

	// render table
	table.SetBorder(false)
	table.Render()
}

func getCellValue(position [2]int, t [][]string) {
	fmt.Printf("Value of cell at (%v,%v) is %v", position[0], position[1], t[position[1]-1][position[0]-1])
}

func setCellValue(position [2]int, t [][]string, value string) [][]string {
	t[position[1]-1][position[0]-1] = value

	return t
}

func appendNewColumn(title string, t [][]string) [][]string {
	t[0] = append(t[0], title)

	for i := 1; i < len(t); i++ {
		t[i] = append(t[i], "")
	}

	return t
}

func getRowValues(row int, file string) {
	// open file and create scanner
	f, _ := os.Open(file)
	scanner := bufio.NewScanner(f)

	// create file lines array
	lines := []string{}

	// iterate through file lines, appending to array
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println(lines[row-1])
}

func getColumnValues(col string, t [][]string) {
	// print column title
	fmt.Println("1", col)

	// get column name position
	columnNumber := 0
	for i := 0; i < len(t[0]); i++ {
		if t[0][i] == col {
			columnNumber = i
			break
		}
	}

	// print column values
	for j := 1; j < len(t); j++ {
		fmt.Println(j+1, t[j][columnNumber])
	}
}