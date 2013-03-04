package main

import (
	"github.com/daviddengcn/gdr/gdrf"
	"github.com/daviddengcn/go-villa"
	"log"
	"os"
)

func main() {
	//inFn := villa.Path("gdr.go")
	inFn := villa.Path("gdrf/filter.go")
	err := gdrf.FilterFile(inFn, os.Stdout)
	if err != nil {
		log.Println(err)
	}
}
