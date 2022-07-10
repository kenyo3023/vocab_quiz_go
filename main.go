package main

import (
	"fmt"
	"time"
	"vocab/quiz/data/process"
)

func main() {
	process.LoadVocabs("vocabs_20220612.csv")

	d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
	year, _, _ := d.Date()
	fmt.Printf("%T", year)
}
