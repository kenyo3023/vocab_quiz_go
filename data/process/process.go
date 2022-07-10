package process

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var Pwd string
var Filepath string = "data/vocabs"

func init() {
	Pwd, _ = os.Getwd()
}

// Check if the file extension of the input document
// meets the requirements (currently only accepts ".csv").
func CheckFileExt(file string) {
	ext := filepath.Ext(file)
	switch ext {
	case ".csv":
		log.Printf("[INFO] The file is valid: %s", file)
	default:
		log.Fatalf("[ERROR] The file is invalid (must be .csv): %s", file)
	}
}

// Check if the file extension of the input document
// actually exists in the path.
func CheckFileExist(file string) {
	if _, err := os.Stat(file); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("[ERROR] The file does not exist: %s", file)
	}
}

func ProcessString2Date(s string) time.Time {
	const dateTemplate = "2/1/2006"
	date, error := time.Parse(dateTemplate, s)

	if error != nil {
		log.Fatalf("[ERROR] Failed to Parse the date string to time: %s", date)
	}
	return date
}

func JoinYearMonthDay2String(year string, month string, day string) string {
	return strings.Join([]string{year, month, day}, "/")
}

func ProcessMeaning2Dict(data []string, meanColIndex []int) map[string]string {
	dict := make(map[string]string)
	for i := 0; i < len(meanColIndex); i++ {
		speechIndex := meanColIndex[i]
		meanIndex := meanColIndex[i] + 1
		dict[data[speechIndex]] = data[meanIndex]
	}
	return dict
}

// func ProcessType2Dict(x int) {
// 	dict := map[interface{}]interface{}{
// 		1:     "hello",
// 		"hey": 2,
// 	}
// 	return dict
// }

// Read a row of data in the given csv.Reader.
func ReadARow(r csv.Reader, c2i *Col2Index) (string, Meta, bool) {
	data, err := r.Read()
	var empty bool = true

	// Determine if the current row is empty
	if c2i == nil {
		log.Fatalf("[ERROR] The index of the columns has not been configured.")
	}
	if len(data[c2i.vocab]) <= 0 {
		empty = true
	} else if err == io.EOF {
		empty = true
	} else {
		empty = false
	}

	// Storage the row data.
	var meta Meta
	meta.proficiency = 1

	date := ProcessString2Date(data[c2i.date])
	year := strconv.Itoa(date.Year())
	month := strconv.Itoa(int(date.Month()))
	day := strconv.Itoa(date.Day())
	meta.date = JoinYearMonthDay2String(year, month, day)
	meta.mean = ProcessMeaning2Dict(data, c2i.mean)

	return data[c2i.vocab], meta, empty
}

// Load the given vocabulary file.
func LoadVocabs(file string) {
	CheckFileExt(file)
	file = filepath.Join(Pwd, Filepath, file)

	CheckFileExist(file)

	// Load and read the given vocabulary file
	f, err := os.OpenFile(file, os.O_RDONLY, 0777)
	if err != nil {
		log.Fatalf("[ERROR] The file load failed: %s", file)
	}

	// Initialize the csv reader.
	r := csv.NewReader(f)
	r.Comma = ','

	// Get the column list of the given vocabulary file.
	var col2index Col2Index
	var prev_col string

	columns, _ := r.Read()
	for i, col := range columns {
		switch col {
		case "單字":
			col2index.vocab = i
		case "日期":
			col2index.date = i
		case "詞性":
			col2index.mean = append(col2index.mean, i)
		case "註解":
			if i-1 >= 0 {
				if prev_col != "詞性" {
					log.Fatalf("[ERROR] The current col is %s, and previous col must be %s", col, prev_col)
				}
			}
		default:
			log.Printf("[WARNING] The vocabulary file has unrecognized column. %s", col)
		}

		if i >= 1 {
			prev_col = col
		}
	}
	fmt.Println(columns)

	// Iterate all the row of given vocabulary file until the end.
	for {
		vocabDict := make(map[string]Meta)

		vocab, meta, isEmpty := ReadARow(*r, &col2index)
		if isEmpty {
			break
		}
		vocabDict[vocab] = meta
		fmt.Println(vocabDict)
	}
}
