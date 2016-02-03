package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Fetch(imageUrl string) ([]byte, error) {
	url, err := url.Parse(imageUrl)
	if err != nil {
		return nil, fmt.Errorf("Invalid image URL: (url=%s)", url.RequestURI())
	}
	return fetchImage(url)
}

func fetchImage(url *url.URL) ([]byte, error) {
	req := createRequest(url)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error downloading image: %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Error downloading image: (status=%d) (url=%s)", res.StatusCode, req.URL.RequestURI())
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Unable to create image from response body: %s (url=%s)", req.URL.RequestURI(), err)
	}
	return buf, nil
}

func createRequest(url *url.URL) *http.Request {
	req, _ := http.NewRequest("GET", url.RequestURI(), nil)
	req.Header.Set("User-Agent", "resizr "+Version)
	req.URL = url
	return req
}
