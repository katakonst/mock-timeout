package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {

	conf := InitConfigs()
	logger := NewLogger("info")
	rule, err := ParseRule(conf.ruleFile)
	if err != nil {
		logger.Fatalf("Reading rule file: %v", err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Errorf("Error while reading req body: %v", err)
			return
		}

		if err = rule.ProcessRequest(r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Errorf("Error while processing req : %v", err)
			return
		}

		logger.Infof("Request for uri %s for domain %s", r.RequestURI,
			rule.Domain)

		url := rule.Domain + "/" + r.RequestURI
		proxyReq, err := http.NewRequest(r.Method, url, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Errorf("Error while creating request: %v", err)
			return
		}

		proxyReq.Header = r.Header
		client := http.Client{}
		resp, err := client.Do(proxyReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			logger.Errorf("Error while doing request: %v", err)
			return
		}
		defer resp.Body.Close()

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("Error while reading response body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for name, values := range resp.Header {
			w.Header().Set(name, strings.Join(values, ","))
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	})
	logger.Infof("Started server on %s", rule.Host)
	log.Fatal(http.ListenAndServe(rule.Host, nil))
}
