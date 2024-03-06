package lib

import (
	"fmt"
	"log"
)

func Check(explanation string, e error, print ...bool) {
	if len(print) > 0 {
		fmt.Printf("checking '%s' expression", explanation)
	}
	if e != nil {
		log.Fatal(explanation, e)
	}
}
