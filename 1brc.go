package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"sync"

	"golang.org/x/exp/mmap"

	"github.com/pkg/profile"

	"github.com/dvordrova/1brc/lib"
)

const numParsers = 12
const numReaders = 6
const readSpeed = 10 * 1024 * 1024

var fileName = flag.String("input", "", "path to the input file to evaluate")
var profilingType = flag.String("profiling", "", "type of profiling to use")

func main() {
	flag.Parse()
	if *fileName == "" {
		log.Fatal("input file is required")
	}
	if profilingType != nil {
		switch *profilingType {
		case "cpu":
			defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profiles")).Stop()
		case "mem":
			defer profile.Start(profile.MemProfile, profile.ProfilePath("./profiles")).Stop()
		case "trace":
			defer profile.Start(profile.TraceProfile, profile.ProfilePath("./profiles")).Stop()
		default:
			log.Fatal("unknown profiling type")
		}
	}
	fmt.Println(evaluate(*fileName))
	// print mapNameMeasurment
	// for k, v := range mapNameMeasurment {
	// 	fmt.Printf("key: %d\n", k)
	// 	for _, m := range v {
	// 		fmt.Printf("%s: max: %d, min: %d, avg: %f\n", m.name, m.max, m.min, float64(m.sum)/float64(m.count))
	// 	}
	// }
}

func evaluate(inputFile string) string {
	file, err := mmap.Open(inputFile)
	lib.Check("failed to open file for mmap: ", err)
	defer file.Close()
	fileSize := file.Len()

	var bytesChunkChannel = make([]chan []string, 0, numParsers)
	for i := 0; i < numParsers; i++ {
		bytesChunkChannel = append(bytesChunkChannel, make(chan []string, 20000))
	}
	outChan := make(chan (map[string]*lib.Measurement))
	var wg sync.WaitGroup
	wg.Add(numReaders)

	for i := 0; i < numParsers; i++ {
		go lib.Parser(bytesChunkChannel[i], outChan)
	}
	for i := 0; i < numReaders; i++ {
		var (
			startChannel int
			endChannel   int
		)

		if numParsers > numReaders {
			if numParsers%numReaders != 0 {
				log.Fatal("numParsers should be divisible by numReaders or wise versa")
			}
			startChannel = i * numParsers / numReaders
			endChannel = (i + 1) * numParsers / numReaders
		} else {
			if numReaders%numParsers != 0 {
				log.Fatal("numParsers should be divisible by numReaders or wise versa")
			}
			howMuchWriters := numReaders / numParsers
			startChannel = i / howMuchWriters
			endChannel = startChannel + 1
		}

		fmt.Printf("Channels are from %d to %d\n", startChannel, endChannel)

		go func() {
			defer wg.Done()
			lib.Reader(file, fileSize, i, numReaders, bytesChunkChannel[startChannel:endChannel], readSpeed)
		}()
	}

	wg.Wait()
	for i := 0; i < numParsers; i++ {
		close(bytesChunkChannel[i])
	}

	resultMap := make(map[string]*lib.Measurement)
	allKeys := make([]string, 0, 10000)
	for i := 0; i < numParsers; i++ {
		mapNameMeasurment := <-outChan
		for k, v := range mapNameMeasurment {
			if _, ok := resultMap[k]; !ok {
				allKeys = append(allKeys, k)
				resultMap[k] = v
			} else {
				resultMap[k].Max = max(resultMap[k].Max, v.Max)
				resultMap[k].Min = min(resultMap[k].Min, v.Min)
				resultMap[k].Sum += v.Sum
				resultMap[k].Count += v.Count
			}
		}
	}
	sort.Strings(allKeys)

	var result string
	for i, k := range allKeys {
		if i != 0 {
			result += ", "
		}
		result += fmt.Sprintf("%s=%s/%s/%s",
			k,
			lib.FormattedNumber(resultMap[k].Min),
			lib.FormattedAvg(resultMap[k].Sum, resultMap[k].Count),
			lib.FormattedNumber(resultMap[k].Max),
		)
	}
	return result
}
