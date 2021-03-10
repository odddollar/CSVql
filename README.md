# CSVql
CSV query language is a simple command-line program for taking some of the annoyance out of working with CSV files.

### Download and installation
The program source can be downloaded from GitHub and built into a binary using ```go build CSVql```. Prebuilt binary files are also available on GitHub.

### Features
The program supports
- Creating CSV files
- Pretty printing tables
- Getting the value of a specific cell
- Getting all values in a row
- Getting all values in a column
- Setting the value of a specific cell
- Appending columns
- Appending empty rows
- Getting number of rows
- Getting number of columns

**NOTE: All functions use base 1 indexing for x and y positions of cells, rows and columns. 1,1 is top-left**

### Usage

To create a CSV file

```CSVql --create [FILE NAME .CSV]```

To pretty print a CSV file

```CSVql --file [PATH TO CSV FILE] --print```

To get the value from a cell

```CSVql --file [PATH TO CSV FILE] --getcell x|y```

To get the values from a row

```CSVql --file [PATH TO CSV FILE] --getrow x```

To get the values from a column

```CSVql --file [PATH TO CSV FILE] --getcolumn [COLUMN NAME]```

To set the value of a specific cell

```CSVql --file [PATH TO CSV FILE] --setcell x|y|value```

To append a new column

```CSVql --file [PATH TO CSV FILE] --appendcolumn [NEW COLUMN NAME]```

To append an empty row

```CSVql --file [PATH TO CSV FILE] --appendrow```

To get number of rows

```CSVql --file [PATH TO CSV FILE] --getnumrows```

To get number of columns

```CSVql --file [PATH TO CSV FILE] --getnumcolumns```

### Example
```
CSVql --create people.csv
CSVql --file people.csv --appendcolumn name
CSVql --file people.csv --appendcolumn age
CSVql --file people.csv --appendcolumn gender
CSVql --file people.csv --appendrow
CSVql --file people.csv --setcell 1|2|Jeff
CSVql --file people.csv --setcell 2|2|34
CSVql --file people.csv --setcell 3|2|Male
CSVql --file people.csv --appendrow
CSVql --file people.csv --setcell 1|3|Andrea
CSVql --file people.csv --setcell 2|3|52
CSVql --file people.csv --setcell 3|3|Female
CSVql --file people.csv --print
```

Creates a file with contents
```
name,age,gender
Jeff,34,Male
Andrea,52,Female
```
And prints below to the console
```
   NAME  | AGE | GENDER
---------+-----+---------
  Jeff   |  34 | Male
  Andrea |  52 | Female
```
