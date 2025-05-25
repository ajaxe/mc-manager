package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

var httpClient = &http.Client{
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func buildApiURL(b, p string) string {
	return fmt.Sprintf("%s/api/%s", strings.TrimSuffix(b, "/"), strings.TrimPrefix(p, "/"))
}

func appBaseURL() string {
	b := app.Window().URL()
	b.Path = ""
	return b.String()
}

func httpGet(u string, v interface{}) error {
	res, err := http.Get(u)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error code: %v", res.StatusCode)
	}

	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &v)

	return err
}
func httpPost(u string, payload, response interface{}) error {
	return httpWithPayload(http.MethodPost, u, payload, response)
}

func httpPut(u string, payload, response interface{}) error {
	return httpWithPayload(http.MethodPut, u, payload, response)
}
func httpDelete(u string, response interface{}) error {
	return httpWithPayload(http.MethodDelete, u, nil, response)
}
func httpWithPayload(method, u string, payload, response interface{}) (err error) {
	res, err := httpWithPayloadInternal(method, u, payload)
	if err != nil {
		return
	}
	defer res.Body.Close()

	b, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(b, &response)

	return
}
func httpWithPayloadInternal(method, u string, payload interface{}) (res *http.Response, err error) {
	buf := bytes.NewBuffer([]byte{})

	if payload != nil {
		v, e := json.Marshal(payload)
		buf = bytes.NewBuffer(v)
		if e != nil {
			err = e
			return
		}
	}

	req, err := http.NewRequest(method, u, buf)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return
	}
	res, err = httpClient.Do(req)
	return
}

func LoginCheck() (redirect string) {
	res, e := httpWithPayloadInternal(http.MethodPost, buildApiURL(appBaseURL(), "/login/check"), struct{}{})
	if h := res.Header.Get("Location"); res.StatusCode == http.StatusFound && h != "" {
		redirect = h
	}
	app.Logf("headers: %v, status: %v, URL: %v", res.Header, res.StatusCode, res.Request.URL)
	app.Logf("error: %v", e)
	return
}
