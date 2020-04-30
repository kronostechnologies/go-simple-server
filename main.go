package main

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"
)

type Config struct {
	Listen       string
	Hosts        map[string]*ResponseConfig
	LogUnmatched bool
}

type ResponseConfig struct {
	Body     string
	Code     int
	Headers  map[string]string
	Redirect string
	Log      bool
}

type Response struct {
	Body     string
	Code     int
	Headers  map[string]string
	Redirect *template.Template
	Log      bool
}

type Redirection struct {
	Host       string
	RequestURI string
	URL        *url.URL
}

func AddResponse(host string, c *ResponseConfig) {
	var redirect *template.Template
	var code int

	if c.Redirect != "" {
		redirect = template.Must(template.New(host).Parse(c.Redirect))
		code = 301
	} else {
		code = 200
	}

	if c.Code > 0 {
		code = c.Code
	}

	r := &Response{
		Body:     c.Body,
		Code:     code,
		Headers:  c.Headers,
		Redirect: redirect,
	}

	responses[host] = r
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	host := strings.Split(r.Host, ":")[0]
	response, found := responses[host]

	requestString := host + r.RequestURI

	if !found {
		if logUnmatched {
			log.Printf("%s: %d", host, 404)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for k, v := range response.Headers {
		w.Header().Add(k, v)
	}

	if response.Redirect != nil {
		var redirectUrl bytes.Buffer
		re := response.Redirect.Execute(&redirectUrl, &Redirection{
			Host:       r.Host,
			RequestURI: r.RequestURI,
			URL:        r.URL,
		})
		if re != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		redirectString := redirectUrl.String()

		http.Redirect(w, r, redirectString, response.Code)

		if response.Log {
			log.Printf("%s: %d (%s)", requestString, response.Code, redirectString)
		}
	} else {
		w.WriteHeader(response.Code)

		_, we := io.WriteString(w, response.Body)
		if we != nil {
			log.Println(we)
		}

		if response.Log {
			log.Printf("%s: %d", requestString, response.Code)
		}
	}
}

var responses = make(map[string]*Response)
var logUnmatched bool

func main() {
	var config Config
	yamlFile, ye := ioutil.ReadFile("config.yaml")
	if ye != nil {
		panic(ye)
	}

	ue := yaml.Unmarshal(yamlFile, &config)
	if ue != nil {
		panic(ue)
	}

	logUnmatched = config.LogUnmatched

	AddResponse("status", &ResponseConfig{Body: "ok", Log: false})

	for h, r := range config.Hosts {
		log.Printf("configured %s\n", h)
		AddResponse(h, r)
	}

	log.Printf("listening on %s\n", config.Listen)
	if err := http.ListenAndServe(config.Listen, http.HandlerFunc(viewHandler)); err != nil {
		panic(err)
	}
}
