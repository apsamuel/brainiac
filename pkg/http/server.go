package http

import (
	"encoding/json"
	"io"
	"net/http"
)

type HTTPServer interface {
	Serve() error
}

// ReadRequestBody cautiously empties the http.Response
func ReadRequestBody(req *http.Request) ([]byte, error) {
	b, err := ReadRequestBodyLimiter(req.Body, 5000000)
	_ = req.Body.Close()
	// ignoring error due to no appreciable impact
	// cf. https://stackoverflow.com/questions/47293975/should-i-error-check-close-on-a-response-body
	return b, err
}

// ReadRequestBodyLimiter limits reading of responses to 5MB
func ReadRequestBodyLimiter(closer io.ReadCloser, allowedRequestSize int64) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(closer, allowedRequestSize))
	if err != nil {
		return nil, err
	}
	return b, nil
}

func UnmarshalRequestBody(i interface{}, r *http.Request) ([]byte, error) {
	b, err := ReadRequestBody(r)
	if err != nil {
		return b, err
	}
	err = json.Unmarshal(b, &i)
	if err != nil {
		return b, err
	}
	return b, nil
}
