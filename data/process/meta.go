package process

// import (
// 	"time"
// )

// func ConvertData2Dict(x []string, col []string) map[string]string {
// 	var data map[string]string
// 	data[col[1]] = x[1]
// 	return data
// }

type Meta struct {
	mean        map[string]string
	date        string
	proficiency int
}

type Col2Index struct {
	vocab int
	date  int
	mean  []int
}
