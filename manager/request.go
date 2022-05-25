package manager

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"MedalHelper/util/log"
)

func Get(rawUrl, cookie string, params url.Values) ([]byte, error) {
	structUrl, err := url.Parse(rawUrl)
	if err != nil {
		log.Error("url %s Parse error: %v", rawUrl, err)
		return nil, err
	}
	structUrl.RawQuery = params.Encode()
	req, err := http.NewRequest("GET", structUrl.String(), nil)
	if err != nil {
		log.Error("NewRequest error: %v", err)
		return nil, err
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		log.Error("Client.Do error: %v, req: %v", err, req)
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func PostForm(url string, data url.Values) error {
	resp, err := http.PostForm(url, data)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		log.Error("URL %s POST error: %v, data: %v", url, err, data)
		return err
	}
	resp.Body.Close()
	return nil
}

func PostFormWithCookie(url, cookie string, data url.Values) ([]byte, error) {
	req, err := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	if err != nil {
		log.Error("NewRequest error: %v", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil || resp.StatusCode != 200 {
		log.Error("Client.Do error: %v, req: %v", err, req)
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
