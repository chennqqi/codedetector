package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"regexp"
	"strings"
	"time"

	"github.com/chennqqi/goutils/utils"
	"github.com/chennqqi/goutils/yamlconfig"
)

var (
	ErrNotMatch = errors.New("Not Matched")
	verbose     = 0
)

type Result struct {
	Language string
	Score    int
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
		for jdx := 0; jdx < len(sd.dict[idx].TopRules); jdx++ {
			if verbose > 2 {
				fmt.Println(sd.dict[idx].Language, "TopRules:", sd.dict[idx].TopRules[jdx].Value)
			}
			sd.dict[idx].TopRules[jdx].exp = regexp.MustCompile(
				sd.dict[idx].TopRules[jdx].Value)
		}
		for jdx := 0; jdx < len(sd.dict[idx].NearTopRules); jdx++ {
			sd.dict[idx].NearTopRules[jdx].exp = regexp.MustCompile(
				sd.dict[idx].NearTopRules[jdx].Value)
			if verbose > 2 {
				fmt.Println(sd.dict[idx].Language, "NearTopRules:", sd.dict[idx].NearTopRules[jdx].Value)
			}
		}

		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
			sd.dict[idx].Rules[jdx].exp = regexp.MustCompile(
				sd.dict[idx].Rules[jdx].Value)
			if verbose > 2 {
				fmt.Println(sd.dict[idx].Language, "Rules:", sd.dict[idx].Rules[jdx].Value)
			}
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
	//	for idx := 0; idx < len(sd.dict); idx++ {
	//		var r Result
	//		r.Language = sd.dict[idx].Language
	//		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
	//			if sd.dict[idx].Rules[jdx].exp.Match(txt) {
	//				r.Score += sd.dict[idx].Rules[jdx].Score
	//			}
	//			if r.Score >= 1.0 {
	//				return r, nil
	//			}
	//		}
	//	}
	return Result{}, ErrNotMatch
}

func (sd *ScriptCodeDetector) DetectAll(txt []byte) ([]Result, error) {
	var r []Result

	//	for idx := 0; idx < len(sd.dict); idx++ {
	//		var t Result
	//		t.Language = sd.dict[idx].Language
	//		for jdx := 0; jdx < len(sd.dict[idx].Rules); jdx++ {
	//			if sd.dict[idx].Rules[jdx].exp.Match(txt) {
	//				t.Score += sd.dict[idx].Rules[jdx].Score
	//			}
	//			if t.Score >= 1.0 {
	//				r = append(r, t)
	//				break
	//			} else if t.Score <= -1.0 {
	//				break
	//			}
	//		}
	//	}
	//	if len(r) > 0 {
	//		return r, nil
	//	}
	return r, ErrNotMatch
}

func doCalScore(rules []RuleItem, line []byte) int {
	var score int
	for _, rule := range rules {
		if verbose > 1 {
			fmt.Print("match:", string(line), rule.Value)
		}
		if rule.exp.Match(line) {
			score += rule.Score
			if verbose > 1 {
				fmt.Println(" OK")
			}
		} else {
			if verbose > 1 {
				fmt.Println(" NOT")
			}
		}
	}
	return score
}

func (sd *ScriptCodeDetector) DetectBest(txt []byte) (Result, error) {
	var idxStart int
	//skip UTF-8 EFBBBF
	if len(txt) >= 3 && txt[0] == 0xEF && txt[1] == 0xBB && txt[2] == 0xBF {
		txt = txt[3:]
	}

	for idx := 0; idx < len(txt); idx++ {
		switch txt[idx] {
		case '\t', ' ', '\r', '\n':
			//donothing

		default:
			idxStart = idx
			goto BREAK_FOR
		}
	}
BREAK_FOR:

	newTxt := bytes.TrimFunc(txt[idxStart:], func(b rune) bool {
		return uint32(b) == '\r'
	})

	linesofCode := bytes.Split(newTxt, []byte("\n"))
	linecount := len(linesofCode)
	neartop := func(line int) bool {
		if linecount <= 10 {
			return true
		}
		return line < linecount/10
	}

	valididxMap := make(map[int]int)
	for idx := 0; idx < linecount; idx++ {
		if neartop(idx) {
			valididxMap[idx] = 1
		} else if idx%int(math.Ceil(float64(linecount)/500.0)) == 0 {
			valididxMap[idx] = 1
		}
	}
	if verbose > 0 {
		fmt.Println("validMap", valididxMap)
	}

	//match Top
	if linecount > 0 {
		for idx := 0; idx < len(sd.dict); idx++ {
			var t Result
			t.Language = sd.dict[idx].Language
			t.Score = doCalScore(sd.dict[idx].TopRules, linesofCode[0])
			if t.Score >= 100 {
				return t, nil
			}
		}
	}

	//match
	rs := make([]Result, len(sd.dict))
	for idx := 0; idx < len(sd.dict); idx++ {
		var t Result
		t.Language = sd.dict[idx].Language
		for lineidx, _ := range valididxMap {
			if neartop(lineidx) {
				t.Score += doCalScore(sd.dict[idx].NearTopRules, linesofCode[lineidx])
				if verbose > 1 {
					fmt.Println(sd.dict[idx].NearTopRules)
				}
			}
			t.Score += doCalScore(sd.dict[idx].Rules, linesofCode[lineidx])
		}
		rs[idx] = t
	}

	var maxIdx int
	for idx := 1; idx < len(rs); idx++ {
		if rs[idx].Score > rs[maxIdx].Score {
			maxIdx = idx
		}
	}
	if verbose > 0 {
		fmt.Println(rs)
	}

	if rs[maxIdx].Score < 2 {
		rs[maxIdx].Language = "unknown"
		return rs[maxIdx], ErrNotMatch
	}
	return rs[maxIdx], nil
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

	language := "jsp"
	utils.DoWalkDir("E:/centosshare/webshell/WebSHArk/WebSHArk", "", func(fileName string, isdir bool) error {
		if !isdir && strings.HasSuffix(fileName, language) {
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
