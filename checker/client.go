package checker

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type checker struct {
	responseCodes map[string]int
	urls          map[string]*url.URL
	m             sync.Mutex
}

const (
	timeout       = 10 * time.Second
	parallelLimit = 100
)

func Query(links []string) map[string]int {
	c := NewChecker(links)
	wg := sync.WaitGroup{}
	limiter := make(chan struct{}, parallelLimit)

	wg.Add(len(links))
	for _, link := range links {
		if !c.ValidLink(link) {
			continue
		}

		limiter <- struct{}{}

		go func() {
			c.fetchStatusCode(link)
			wg.Done()
			<-limiter
		}()
	}
	wg.Wait()

	return c.GetCodes()
}

func NewChecker(links []string) *checker {
	c := checker{
		responseCodes: make(map[string]int, len(links)),
		urls:          make(map[string]*url.URL),
	}
	c.setURLs(links)
	return &c
}

func (c *checker) setURLs(links []string) {
	for _, link := range links {
		c.responseCodes[link] = 0
		if u, err := url.Parse(link); err == nil && u.Scheme != "" && u.Host != "" {
			c.urls[link] = u
		} else {
			log.Printf("Invalid URL: %q", link)
		}
	}
}

func (c *checker) ValidLink(link string) bool {
	_, ok := c.urls[link]
	return ok
}

func (c *checker) GetCodes() map[string]int { return c.responseCodes }

func (c *checker) fetchStatusCode(link string) {
	request := buildRequest(link)
	for i := 1; i <= 5; i++ {
		response, err := (&http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:     true,
				IdleConnTimeout:       time.Second,
				ResponseHeaderTimeout: timeout,
			},
			Timeout: timeout,
		}).Do(request)
		if err != nil {
			log.Printf("Request error: %s", err)
			continue
		}
		c.m.Lock()
		c.responseCodes[link] = response.StatusCode
		c.m.Unlock()
		_ = response.Body.Close()
		break
	}
}

func buildRequest(link string) *http.Request {
	req, _ := http.NewRequest("GET", link, nil)
	return req
}
