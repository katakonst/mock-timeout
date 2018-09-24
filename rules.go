package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

type Rule struct {
	Host   string
	Domain string
	Routes []Route
}

type Route struct {
	Pattern string
	Method  string
	Timeout string
}

func (rule *Rule) ProcessRequest(r *http.Request) error {
	for _, route := range rule.Routes {
		routePat, err := regexp.Compile(route.Pattern)
		if err != nil {
			return fmt.Errorf("computing regex %s: %v", route.Pattern, err)
		}

		if routePat.MatchString(r.RequestURI) == true &&
			r.Method == route.Method {
			if d, err := time.ParseDuration(route.Timeout); err == nil {
				time.Sleep(d)
			} else {
				return fmt.Errorf("Error at timeout %v", err)
			}
			return nil
		}
	}
	return nil
}

func ParseRule(filename string) (*Rule, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file: %v", err)
	}
	var rule Rule

	if err = json.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("unmarshaling file: %v", err)
	}
	return &rule, nil
}
