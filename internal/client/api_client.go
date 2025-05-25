package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const (
	httpMethodGet    = "get"
	httpMethodPost   = "post"
	httpMethodPut    = "put"
	httpMethodDelete = "delete"
)

func buildApiURL(b, p string) string {
	return fmt.Sprintf("%s/api/%s", strings.TrimSuffix(b, "/"), strings.TrimPrefix(p, "/"))
}

func appBaseURL() string {
	b := app.Window().URL()
	b.Path = ""
	return b.String()
}

func httpGet(u string, v interface{}) error {
	resp, code, err := httpCall("", u, nil)

	if !successful(code) {
		return fmt.Errorf("error code: %v", code)
	}

	b, _ := io.ReadAll(strings.NewReader(*resp))
	err = json.Unmarshal(b, &v)

	return err
}
func httpPost(u string, payload, response interface{}) error {
	return httpWithPayload(httpMethodPost, u, payload, response)
}

func httpPut(u string, payload, response interface{}) error {
	return httpWithPayload(httpMethodPut, u, payload, response)
}
func httpDelete(u string, response interface{}) error {
	return httpWithPayload(httpMethodDelete, u, nil, response)
}
func httpWithPayload(method, u string, payload, response interface{}) (err error) {
	resp, _, err := httpCall(method, u, payload)
	if err != nil {
		return
	}

	b, _ := io.ReadAll(strings.NewReader(*resp))

	err = json.Unmarshal(b, &response)

	return
}
func httpCall(method, u string, payload interface{}) (resp *string, code int, err error) {
	p := map[string]any{}
	if payload != nil {
		buf := bytes.NewBuffer([]byte{})
		v, e := json.Marshal(payload)
		buf = bytes.NewBuffer(v)
		if e != nil {
			err = e
			return
		}
		p["body"] = string(buf.Bytes())
	}

	if method == "" {
		method = "GET"
	}
	p["method"] = method
	p["headers"] = map[string]any{
		"Content-Type": "application/json",
	}

	resp, code, err = fetch(u, &p)

	return
}

func LoginCheck() {
	_, _, e := httpCall("post", buildApiURL(appBaseURL(), "/login/check"), struct{}{})

	app.Logf("error: %v", e)
	return
}
func successful(code int) bool {
	return 100 <= code && code <= 399
}
