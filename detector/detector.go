package detector

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"regexp"

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
	//Get Rule language Index(language list)
	GetRuleIndex() []string

	//Detect a txt and return result, if not matched return ErrNotMatch
	Detect(txt []byte) (Result, error)
}

type ScriptCodeDetector struct {
	dict Rules
}

//Load yml rules
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

func doCalScore(rules []RuleItem, line []byte) int {
	var score int
	for _, rule := range rules {
		if verbose > 1 {
			fmt.Print("match: ", string(line), rule.Value)
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

func (sd *ScriptCodeDetector) Detect(txt []byte) (Result, error) {
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
		if line <= 10 {
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
			if verbose > 2 {
				fmt.Println("domatch:", lineidx, string(linesofCode[lineidx]))
			}
			if neartop(lineidx) {
				if verbose > 2 {
					fmt.Println("neartop:", lineidx)
				}
				t.Score += doCalScore(sd.dict[idx].NearTopRules, linesofCode[lineidx])
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
