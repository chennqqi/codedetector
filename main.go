package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/chennqqi/codedetector/detector"
)

func main() {
	var ext string
	var name string
	flag.StringVar(&ext, "e", "php", "set specific languange(default:php)")
	flag.StringVar(&name, "name", "", "set specific languange(default:php)")
	flag.IntVar(&verbose, "v", 0, "set verbose(default:0)")
	flag.Parse()

	if name == "" {
		name = ext
	}

	det, err := detector.LoadInteralRules("rule.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("Indexs:", det.GetRuleIndex())
	var dcount, ndcount float64
	t := time.Now()

	language := ext
	utils.DoWalkDir("E:/centosshare/webshell/WebSHArk/WebSHArk", "", func(fileName string, isdir bool) error {
		if !isdir && strings.HasSuffix(fileName, name) {
			txt, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Println(err)
			} else {
				r, _ := det.DetectBest(txt)
				if r.Language != language {
					ndcount += 1
					fmt.Println(fileName, r)
				} else {
					dcount += 1
					//fmt.Println("[OK]", fileName, r)
				}
			}
		}
		return nil
	})
	fmt.Println("识别率为", dcount/(dcount+ndcount), dcount+ndcount, "cost", time.Now().Sub(t))
}
