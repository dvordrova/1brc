package lib

import (
	"fmt"
	"io"
)

func Reader(reader io.ReaderAt, fileSize int, readerNum int, totalReaderCount int, outChan []chan []string, stepSize int) {
	outSize := len(outChan)
	batchCap := 1024
	var (
		stringBatch            = make([]string, 0, batchCap)
		prev             int   = 0
		readIndex        int64 = 0
		curOut           int   = 0
		length           int   = 0
		err              error
		skipFirstRead    bool = readerNum != 0
		continueReading  bool = true
		hadOneStartBlock bool = readerNum == 0
		// TODO: use this buffer
		// buffer            [][]byte = make([][]byte, outSize)
	)
	startPos := readerNum * fileSize / totalReaderCount
	endPos := (readerNum + 1) * fileSize / totalReaderCount
	fmt.Printf("total read should be from %d to %d\n", startPos, endPos)

	for continueReading {
		readBytes := make([]byte, stepSize)
		length, err = reader.ReadAt(readBytes, int64(startPos))
		readString := string(readBytes)

		// fmt.Printf("%d reader: read %d bytes from %d\n", readerNum, length, startPos)
		channelWritten := 0
		readIndex++
		i := 0
		prev = 0
		// skip first line
		if skipFirstRead {
			skipFirstRead = false
			for i < length && readString[i] != '\n' {
				i++
			}
			if i < length && startPos+i < endPos {
				hadOneStartBlock = true
			}
			i++
			// fmt.Printf("Skip first %d bytes\n", i)
			prev = i
		}
		// end line

		lastBiteWasNotWritten := false

		for ; i < length; i++ {
			for i < length && readString[i] != '\n' {
				i++
				// fmt.Println("i: ", i)
			}
			if startPos+i > endPos && !hadOneStartBlock {
				continueReading = false
				break
			}
			if i != length {
				// fmt.Println("i: ", i, "prev: ", prev, "read_bytes: ", string(read_bytes[prev:i]))
				// fmt.Printf("ready to write")
				stringBatch = append(stringBatch, readString[prev:i])
				channelWritten++
				// fmt.Printf("written to channel: %d\n", channelWritten)
				// clear lasted_read_bytes
				prev = i + 1
			} else {
				lastBiteWasNotWritten = true
				// fmt.Printf("lasted_read_bytes: %s prev: %d length: %d read_bytes: %s\n", lasted_read_bytes, prev, length, read_bytes[prev:length])
			}
			if len(stringBatch) == batchCap {
				outChan[curOut] <- stringBatch
				curOut = (curOut + 1) % outSize
				stringBatch = make([]string, 0, batchCap)
			}
			if startPos+i > endPos {
				continueReading = false
				// fmt.Printf("end of interval, read extra %d bytes\n", startPos+i-endPos)
				break
			}
		}

		if lastBiteWasNotWritten {
			startPos += prev
		} else {
			startPos += length
		}

		if err == io.EOF {
			continueReading = false
			continue
		} else {
			// fmt.Printf("%d reader: length: %d, err: %v\n", readerNum, length, err)
			Check("error while reading input file: ", err)
		}
		// break
	}

	if len(stringBatch) != 0 {
		outChan[curOut] <- stringBatch
	}
}
