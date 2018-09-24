package main

import "flag"

type Config struct {
	ruleFile string
	Host     string
}

func InitConfigs() Config {
	ruleFile := flag.String("ruleFile", "rule.json", "rule file")
	return Config{
		ruleFile: *ruleFile,
	}
}
