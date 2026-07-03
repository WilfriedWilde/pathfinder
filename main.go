package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// check number command line arguments
	args := os.Args
	if len(args) < 5 {
		fmt.Println("Error")
		return
	}

	// extract args
	mapPath, startStation, endStation := args[1], args[2], args[3]

	// read network map
	file, err := os.Open(mapPath)
	if err != nil {
		fmt.Println("Error opening file:",  err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan(){
		line := strings.TrimSpace(scanner.Text())
		hashtagIndex := strings.Index(line, "#")
		fmt.Println(hashtagIndex)
		if hashtagIndex > 0 {
			line = line[:hashtagIndex]
			continue
		}

		fmt.Println(line)
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Println("Scanner error:", err)
		return
	}

	// check if start == end
	if startStation == endStation {
		fmt.Println("Error")
		return
	}		
}
