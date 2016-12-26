package detector

import (
	"regexp"
)

type RuleItem struct {
	Value string         `json:"value" yaml:"value"`
	Score int            `json:"score" yaml:"score"`
	exp   *regexp.Regexp `-`
}

type Rule struct {
	Language     string     `json:"language" yaml:"language"`
	TopRules     []RuleItem `json:"top" yaml:"top"`
	NearTopRules []RuleItem `json:"neartop" yaml:"neartop"`
	Rules        []RuleItem `json:"rules" yaml:"rules"`
}

type Rules []Rule
