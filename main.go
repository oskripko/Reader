package main

import (
	"fmt"
	"io/ioutil"
	"reader/counter"
)


func main() {
	dirName := "counter"
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		fmt.Print(err)
	}
	stats := counter.GetNumberOfAscii(files, dirName)
	stats.Range(func(k, v interface{}) bool {
		fmt.Printf("Count for %v: %v \n", k, v)
		return true
	})
}
