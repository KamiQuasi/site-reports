package pagespeed

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Site struct
type Site struct {
	URL        string
	Protocol   string
	Repository string
	Analysis   []Analyzer
	Scores     []Result
}

func (p *Site) analyze() error {
	request := Request{}
	request.Strategy = "mobile"
	p.Scores = append(p.Scores, request.get(p)[0])

	return nil
}

// Analyzer interface
type Analyzer interface {
	analyze()
}

// Request struct
type Request struct {
	URL                       string
	FilterThirdPartyResources bool
	Locale                    string
	Rule                      string
	Screenshot                bool
	Strategy                  string
}

func (ps *Request) get(sites ...*Site) []Result {
	var results []Result

	for _, site := range sites {
		resp, err := http.Get(fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v2/runPagespeed?url=%s&strategy=%s", site.URL, ps.Strategy))
		if err != nil {

		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {

		}
		result := Result{}
		json.Unmarshal(body, &result)
		results = append(results, result)
	}

	return results
}

// Result struct
type Result struct {
	Kind             string           `json:"kind"`
	ID               string           `json:"id"`
	ResponseCode     int              `json:"responseCode"`
	Title            string           `json:"title"`
	RuleGroups       TestGroups       `json:"ruleGroups,omitempty"`
	PageStats        PageStats        `json:"pageStats"`
	Version          Version          `json:"version"`
	FormattedResults FormattedResults `json:"formattedResults"`
}

// TestGroups struct
type TestGroups struct {
	Speed     RuleSet `json:"SPEED,omitempty"`
	Usability RuleSet `json:"USABILITY,omitempty"`
}

// RuleSet struct
type RuleSet struct {
	Score int `json:"score,omitempty"`
}

// PageStats struct
type PageStats struct {
	NumberResources         int `json:"numberResources"`
	NumberHosts             int `json:"numberHosts"`
	TotalRequestBytes       int `json:"totalRequestBytes,string"`
	NumberStaticResources   int `json:"numberStaticResources,string"`
	HTMLResponseBytes       int `json:"htmlResponseBytes,string"`
	CSSResponseBytes        int `json:"cssResponseBytes,string"`
	ImageResponseBytes      int `json:"imageResponseBytes,string"`
	JavaScriptResponseBytes int `json:"javascriptResponseBytes,string"`
	OtherResponseBytes      int `json:"otherResponseBytes,string"`
	NumberJsResources       int `json:"numberJsResources"`
	NumberCSSResorces       int `json:"numberCssResources"`
}

// Version struct
type Version struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

// FormattedResults struct
type FormattedResults struct {
	Locale      string      `json:"locale"`
	RuleResults RuleResults `json:"ruleResults"`
}

// RuleResults struct
type RuleResults struct {
	AvoidLandingPageRedirects       Rule `json:"AvoidLandingPageRedirects"`
	EnableGzipCompression           Rule `json:"EnableGzipCompression"`
	LeverageBrowserCaching          Rule `json:"LeverageBrowserCaching"`
	MainResourceServerResponseTime  Rule `json:"MainResourceServerResponseTime"`
	MinifyCSS                       Rule `json:"MinifyCss"`
	MinifyHTML                      Rule `json:"MinifyHTML"`
	MinifyJavaScript                Rule `json:"MinifyJavaScript"`
	MinimizeRenderBlockingResources Rule `json:"MinimizeRenderBlockingResources"`
	OptimizeImages                  Rule `json:"OptimizeImages"`
	PrioritizeVisibleContent        Rule `json:"PrioritizeVisibleContent"`
}

// Rule struct
type Rule struct {
	LocalizedRuleName string     `json:"localizedRuleName"`
	RuleImpact        float64    `json:"ruleImpact"`
	Groups            []string   `json:"groups"`
	Summary           Formatter  `json:"summary"`
	URLBlocks         []URLBlock `json:"urlBlocks"`
}

// URLBlock struct
type URLBlock struct {
	Header Formatter   `json:"header"`
	URLs   []URLResult `json:"urls"`
}

// URLResult struct
type URLResult struct {
	Result Formatter `json:"result"`
}

// Formatter struct
type Formatter struct {
	Format string     `json:"format"`
	Args   []Argument `json:"args"`
}

// Argument struct
type Argument struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
