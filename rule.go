package main

import (
	"regexp"
)

type RuleItem struct {
	Value string         `json:"value" yaml:"value"`
	Score float64        `json:"score" yaml:"score"`
	exp   *regexp.Regexp `-`
}

type Rule struct {
	Language string     `json:"language" yaml:"language"`
	Rules    []RuleItem `json:"rules" yaml:"rules"`
}

type Rules []Rule
