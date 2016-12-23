package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/chennqqi/goutils/utils"
	"github.com/chennqqi/goutils/yamlconfig"
)

var (
	ErrNotMatch = errors.New("Not Matched")
)

type Result struct {
	Language string
	Score    float64
}

type CodeDetector interface {
	GetRuleIndex() []string
	Detect(txt []byte) (Result, error)
	DetectAll(txt []byte) ([]Result, error)
	DetectBest(txt []byte) (Result, error)
}

type ScriptCodeDetector struct {
	dict Rules
}

func LoadRules(rulename string) (CodeDetector, error) {
	var sd ScriptCodeDetector
	if err := yamlconfig.Load(&sd.dict, rulename); err != nil {
		return nil, err
	}
	for idx := 0; idx < len(sd.dict); idx++ {
		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
			sd.dict[idx].Rules[jdx].exp = regexp.MustCompile(
				sd.dict[idx].Rules[jdx].Value)
		}
	}
	return &sd, nil
}

func (sd *ScriptCodeDetector) GetRuleIndex() []string {
	indexs := make([]string, len(sd.dict))
	for idx := 0; idx < len(sd.dict); idx++ {
		indexs[idx] = sd.dict[idx].Language
	}
	return indexs
}

func (sd *ScriptCodeDetector) Detect(txt []byte) (Result, error) {
	for idx := 0; idx < len(sd.dict); idx++ {
		var r Result
		r.Language = sd.dict[idx].Language
		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
			if sd.dict[idx].Rules[jdx].exp.Match(txt) {
				r.Score += sd.dict[idx].Rules[jdx].Score
			}
			if r.Score >= 1.0 {
				return r, nil
			}
		}
	}
	return Result{}, ErrNotMatch
}

func (sd *ScriptCodeDetector) DetectAll(txt []byte) ([]Result, error) {
	var r []Result

	for idx := 0; idx < len(sd.dict); idx++ {
		var t Result
		t.Language = sd.dict[idx].Language
		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
			if sd.dict[idx].Rules[jdx].exp.Match(txt) {
				t.Score += sd.dict[idx].Rules[jdx].Score
			}
			if t.Score >= 1.0 {
				r = append(r, t)
				break
			} else if t.Score <= -1.0 {
				break
			}
		}
	}
	if len(r) > 0 {
		return r, nil
	}
	return r, ErrNotMatch
}

func (sd *ScriptCodeDetector) DetectBest(txt []byte) (Result, error) {
	var r []Result
	for idx := 0; idx < len(sd.dict); idx++ {
		var t Result
		t.Language = sd.dict[idx].Language
		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
			if sd.dict[idx].Rules[jdx].exp.Match(txt) {
				t.Score += sd.dict[idx].Rules[jdx].Score
			}
			if t.Score <= -1.0 {
				break
			}
		}
	}
	var maxIdx int
	for idx := 1; idx < len(r); idx++ {
		if r[idx].Score > r[maxIdx].Score {
			maxIdx = idx
		}
	}

	return r[maxIdx], nil
}

func (sd *ScriptCodeDetector) DetectBest2(txt []byte) (Result, error) {
	newTxt := bytes.Replace(txt, []byte("\r\n"), []byte("\n"))
	linesofCode := bytes.Split(newTxt, []byte("\n"))
	linenos := len(linesofCode)
	neartop := func(line int) {
		if linenos <= 10 {
			return true
		}
		return line < linenos/10
	}

	for idx := 0; idx < len(sd.dict); idx++ {
		var t Result
		t.Language = sd.dict[idx].Language
		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
			if sd.dict[idx].Rules[jdx].exp.Match(txt) {
				t.Score += sd.dict[idx].Rules[jdx].Score
			}
			if t.Score <= -1.0 {
				break
			}
		}
	}
	var maxIdx int
	for idx := 1; idx < len(r); idx++ {
		if r[idx].Score > r[maxIdx].Score {
			maxIdx = idx
		}
	}

	return r[maxIdx], nil
}

func main() {
	det, err := LoadRules("rule.yml")
	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("Indexs:", det.GetRuleIndex())
	var dcount, ndcount float64
	t := time.Now()
	utils.DoWalkDir("E:/centosshare/webshell/WebSHArk/WebSHArk", "", func(fileName string, isdir bool) error {
		if !isdir && strings.HasSuffix(fileName, "") {
			txt, err := ioutil.ReadFile(fileName)
			if err != nil {
				fmt.Println(err)
			} else {
				r, err := det.DetectAll(txt)
				if err != nil {
					fmt.Println(fileName, err)
					ndcount += 1
				} else {
					ndcount += 1
					if len(r) > 1 {
						fmt.Println(fileName, r)
					}
				}
			}
		}
		return nil
	})
	fmt.Println("识别率为", dcount/(dcount+ndcount), dcount+ndcount, "cost", time.Now().Sub(t))
}
