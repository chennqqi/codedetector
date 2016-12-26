package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/chennqqi/codedetector/detector"
	"github.com/chennqqi/goutils/utils"
)

func main() {
	var fileName string
	var dirName string
	flag.StringVar(&dirName, "d", "", "set working in dir mode(<-d/-f> is must")
	flag.StringVar(&fileName, "f", "", "set target file path(<-d/-f> is must")
	flag.Parse()

	det, err := detector.LoadInteralRules()
	if err != nil {
		fmt.Println(err)
		return
	}

	baseName := filepath.Base(os.Args[0])
	if dirName == "" && fileName == "" {
		fmt.Println("usage: ", baseName, "<-f targetName/-d targetdir>")
		return
	}
	if fileName != "" {
		txt, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		r, err := det.Detect(txt)
		fmt.Println(fileName, r)
		return
	}

	var dcount, ndcount float64
	t := time.Now()

	utils.DoWalkDir(dirName, "", func(fileName string, isdir bool) error {
		if !isdir {
			txt, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Println(err)
			} else {
				r, _ := det.Detect(txt)
				fmt.Println("[OK]", fileName, r)
			}
		}
		return nil
	})
	fmt.Println("Rate", dcount/(dcount+ndcount), dcount+ndcount, "cost", time.Now().Sub(t))
}
