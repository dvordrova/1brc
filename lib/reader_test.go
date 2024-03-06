package lib

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"
)

type TestArgs struct {
	testName   string
	fileName   string
	outputChan []chan []string
	numReaders int
}

func TestParser(t *testing.T) {
	tests := make([]TestArgs, 0)
	// find files from ./test_data folder with prefix "reader_"
	dir, err := os.Open("./test_data")
	if err != nil {
		panic(err)
	}
	dir_files, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}
	for _, file := range dir_files {
		name := file.Name()
		spliited_name := strings.Split(name, "_")
		if spliited_name[2] != "reader" {
			continue
		}
		readersNum, err := strconv.Atoi(spliited_name[0])
		if err != nil {
			panic(err)
		}
		outNum, err := strconv.Atoi(spliited_name[1])

		tests = append(tests, TestArgs{
			testName:   file.Name(),
			fileName:   "./test_data/" + file.Name(),
			outputChan: make([]chan []string, 0, outNum),
			numReaders: readersNum,
		})
		for i := 0; i < outNum; i++ {
			tests[len(tests)-1].outputChan = append(tests[len(tests)-1].outputChan, make(chan []string, 1000))
		}
	}
	fmt.Println("Prepared tests")

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			file, err := os.Open(tt.fileName)
			if err != nil {
				panic(err)
			}
			defer file.Close()
			stat, err := file.Stat()

			wg := sync.WaitGroup{}
			wg.Add(tt.numReaders)
			for i := 0; i < tt.numReaders; i++ {
				go func() {
					defer wg.Done()
					Reader(file, int(stat.Size()), i, tt.numReaders, tt.outputChan, 1000)
				}()
			}
			fmt.Println("Waiting for all readers to finish")
			wg.Wait()
			fmt.Println("Waiting for all readers to finish [done]")
			readedStrs := make(map[string]int, 0)
			for i := 0; i < len(tt.outputChan); i++ {
				close(tt.outputChan[i])
				for lineChunk := range tt.outputChan[i] {
					for _, line := range lineChunk {
						readedStrs[line]++
					}
				}
			}

			file2, err := os.Open(tt.fileName)
			if err != nil {
				panic(err)
			}
			defer file2.Close()
			expectedStrs := make(map[string]int, 0)
			scanner := bufio.NewScanner(file2)
			for scanner.Scan() {
				line := scanner.Text()
				expectedStrs[line]++
			}
			if !reflect.DeepEqual(expectedStrs, readedStrs) {
				t.Errorf("Expected: %v, got: %v", expectedStrs, readedStrs)
			}
		})
	}
}
