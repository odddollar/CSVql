package main

import (
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
	printData := parser.Flag("", "print", &argparse.Options{Default: false, Help: "Display the loaded CSV file as a table"})
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
	if len(os.Args) == 1 {
		fmt.Print(parser.Usage(""))
	} else {
		err := parser.Parse(os.Args)
		if err != nil {
			fmt.Print(parser.Usage(err))
			return
		}
	}

	// check if creating file
	if *create != "" {
		_, _ = os.Create(*create)
	} else {
		// load data
		data := loadDataToFrame(*file)

		// if not creating file, check which arguments are given
		if *setCell != "" {
			x, _ := strconv.Atoi(strings.Split(*setCell, ";")[0])
			y, _ := strconv.Atoi(strings.Split(*setCell, ";")[1])
			value := strings.Split(*setCell, ";")[2]

			// remove old file and write new data
			setCellValue([2]int{x, y}, data, value)
			_ = os.Remove(*file)
			f, _ := os.Create(*file)
			writer := csv.NewWriter(f)
			_ = writer.WriteAll(*data)
		}
		if *appendColumn != "" {
			// append new column
			appendNewColumn(*appendColumn, data)

			// remove old file and write new data
			_ = os.Remove(*file)
			f, _ := os.Create(*file)
			writer := csv.NewWriter(f)
			_ = writer.WriteAll(*data)
		}
		if *appendRow {
			var commas []string
			for i := 0; i < len((*data)[0]); i++ {
				commas = append(commas, "")
			}
			*data = append(*data, commas)

			// remove old file and write new data
			_ = os.Remove(*file)
			f, _ := os.Create(*file)
			writer := csv.NewWriter(f)
			_ = writer.WriteAll(*data)
		}
		if *getColumn != "" {
			getColumnValues(*getColumn, data)
		}
		if *getRow != 0 {
			getRowValues(*getRow, data)
		}
		if *getCell != "" {
			x, _ := strconv.Atoi(strings.Split(*getCell, ";")[0])
			y, _ := strconv.Atoi(strings.Split(*getCell, ";")[1])
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

func loadDataToFrame(filePath string) *[][]string {
	// open file and create csv reader
	file, _ := os.Open(filePath)
	defer func() { _ = file.Close() }()
	r := csv.NewReader(file)

	// read through csv lines
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	if len(records) == 0 {
		records = append(records, []string{})
	}

	return &records
}

func getNumColumns(t *[][]string) {
	// print number of columns
	fmt.Println(len((*t)[0]))
}

func getNumRows(t *[][]string) {
	// print number of rows
	fmt.Println(len(*t))
}

func printTable(t *[][]string) {
	// get table headers
	headers := (*t)[0]

	// create table writer and set headers
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	// add data to table
	for i := 1; i < len(*t); i++ {
		table.Append((*t)[i])
	}

	// render table
	table.SetBorder(false)
	table.Render()
}

func getCellValue(position [2]int, t *[][]string) {
	// print value of cell at position
	fmt.Println((*t)[position[1]-1][position[0]-1])
}

func setCellValue(position [2]int, t *[][]string, value string) {
	// set the value of the cell
	(*t)[position[1]-1][position[0]-1] = value
}

func appendNewColumn(title string, t *[][]string) {
	// add new column to title row
	(*t)[0] = append((*t)[0], title)

	// iterate through empty last column and add empty data
	for i := 1; i < len(*t); i++ {
		(*t)[i] = append((*t)[i], "")
	}
}

func getRowValues(row int, t *[][]string) {
	// check if there are any rows
	if len(*t) > 0 {
		// print each row entry, separated with comma
		for i := 0; i < len((*t)[row-1]); i++ {
			fmt.Print((*t)[row-1][i], ", ")
		}
	}
}

func getColumnValues(col string, t *[][]string) {
	// print column title
	fmt.Println("1", col)

	// get column name position
	columnNumber := 0
	for i := 0; i < len((*t)[0]); i++ {
		if (*t)[0][i] == col {
			columnNumber = i
			break
		}
	}

	// print column values
	for j := 1; j < len(*t); j++ {
		fmt.Println(j+1, (*t)[j][columnNumber])
	}
}
