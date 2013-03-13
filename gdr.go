package main

import (
	"github.com/daviddengcn/gdr/gdrf"
	"github.com/daviddengcn/go-villa"
	"log"
	"os"
	"fmt"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Usage: gdr file.go [...]")
		os.Exit(1)
	}
	
	// Extract go files
	var files []villa.Path
	for i := 1; i < len(os.Args); i ++ {
		s := strings.ToLower(strings.TrimSpace(os.Args[i]))
		if strings.HasSuffix(s, ".go") {
			files = append(files, villa.Path(os.Args[i]).Clean())
		} else {
			break;
		}
	}
	
	if len(files) == 0 {
		fmt.Println("Usage: gdr file.go [...]")
		os.Exit(1)
	}

	// Create temporary directory
	tmpDir, err := villa.Path("").TempDir("gdr_")
	if err != nil {
		log.Fatal(err)
	}
	
	// Make file list in tmpDir
	var params villa.StringSlice
	for _, inF := range files {
		dst := tmpDir.Join(inF.Base())
		params.Add(dst)
		func() {
			dstF, err := dst.Create()
			if err != nil {
				log.Fatal(err)
			}
			defer dstF.Close()

			err = gdrf.FilterFile(inF, dstF)
			if err !=  nil {
				log.Fatal(err)
			}
		}()
	}
	
	// Run it
	params.InsertSlice(len(params), os.Args[len(files)+1:])
	params.Insert(0, "run")
	cmd := villa.Path("go").Command(params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}
