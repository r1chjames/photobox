package http

import (
	"bytes"
	"encoding/json"
	"gitlab.com/r1chjames/photobox/api/internal/appconfig"
	"log"
	"net/http"
	"net/http/httptest"
)

func PerformRequest(r http.Handler, method, path string, queryParams map[string]string, body interface{}) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		payload, err := json.Marshal(&body)
		if err != nil {
			log.Fatalf("unable to marshall JSON body payload, %s", err)
		}
		bodyPayload := bytes.NewBuffer(payload)
		req, _ = http.NewRequest(method, path, bodyPayload)
	} else {
		req, _ = http.NewRequest(method, path, nil)
	}

	if len(queryParams) > 0 {
		q := req.URL.Query()
		for param, value := range queryParams {
			q.Add(param, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAppConfig() appconfig.AppConfig {
	return appconfig.AppConfig{ApiBasePath: "test"}
}
