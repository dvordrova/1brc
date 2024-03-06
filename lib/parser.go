package lib

import (
	// "math/rand"
	"strings"
)

func Parser(byteChan chan []string, outputChan chan (map[string]*Measurement)) {
	// parserRandName := rand.Int63()
	// runtime.Gosched()
	result := make(map[string]*Measurement, 10000)
	value := int64(0)
	// runtime.chanrecv2(byteChan)
	for lines := range byteChan {
		// if length > 0 {
		// 	continue
		// } else {
		// 	fmt.Println("Empty input")
		// 	continue
		// }
		for _, line := range lines {
			// fmt.Printf("parserRandName=%d: parser work with line: %s\n", parserRandName, line)
			sepInd := strings.IndexByte(line, ';')
			name := string(line[:sepInd])
			value = 0
			for _, c := range line[sepInd+1:] {
				if c != '-' && c != '.' {
					value = value*10 + int64(c-'0')
				}
			}
			if line[sepInd+1] == '-' {
				value *= -1
			}

			measurment, ok := result[name]
			if !ok {
				result[name] = &Measurement{value, value, value, 1}
			} else {
				measurment.Max = max(measurment.Max, value)
				measurment.Min = min(measurment.Min, value)
				measurment.Sum += value
				measurment.Count++
			}
		}
	}
	// debug what return
	// fmt.Printf("parserRandName=%d: parser return result: %v\n", parserRandName, result)
	outputChan <- result
}
