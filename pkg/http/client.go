package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"time"
)

var nilResponseError = errors.New("response_empty")
var requestCreationError = errors.New("request_creation")
var makeTracedRequestStatsError = errors.New("make_traced_request_stats_error")

// ApiClient is the tool's http client
type ApiClient struct {
	Handle        http.Client
	TracedRequest *TracedRequest
	SentRequests  int
}

// TracedRequest holds the variables that each request uses to measure latency
type TracedRequest struct {
	dnsLookupBeginTime,
	dialBeginTime,
	dialDoneTime,
	obtainedConnTime,
	firstResponseByteTime,
	tlsHandshakeStartTime,
	tlsHandshakeDoneTime time.Time
	requestID string
	Request   *http.Request
}

// ApiResponse holds sufficient info from simple api client response
type ApiResponse struct {
	Body       []byte
	StatusCode int
	Err        error
	RequestId  string
}

func NewTracedGetRequestWithContext(ctx context.Context, baseUrl string, headers map[string]string, body []byte) (*TracedRequest, error) {
	var s TracedRequest
	trace := &httptrace.ClientTrace{
		DNSStart: func(_ httptrace.DNSStartInfo) { s.dnsLookupBeginTime = time.Now() },
		DNSDone:  func(_ httptrace.DNSDoneInfo) { s.dialBeginTime = time.Now() },
		ConnectStart: func(_, _ string) {
			if s.dialBeginTime.IsZero() { // connecting directly to IP
				s.dialBeginTime = time.Now()
			}
		},
		ConnectDone: func(net, addr string, err error) {
			if err != nil {
				return
			}
			s.dialDoneTime = time.Now()
		},
		GotConn:              func(_ httptrace.GotConnInfo) { s.obtainedConnTime = time.Now() },
		GotFirstResponseByte: func() { s.firstResponseByteTime = time.Now() },
		TLSHandshakeStart:    func() { s.tlsHandshakeStartTime = time.Now() },
		TLSHandshakeDone:     func(_ tls.ConnectionState, _ error) { s.tlsHandshakeDoneTime = time.Now() },
	}
	traceCtx := httptrace.WithClientTrace(ctx, trace)
	if body != nil {
		s.Request, _ = http.NewRequestWithContext(traceCtx, "POST", baseUrl, bytes.NewReader(body))
	} else {
		s.Request, _ = http.NewRequestWithContext(traceCtx, "GET", baseUrl, nil)
	}

	if s.Request == nil {
		return nil, requestCreationError
	}
	if headers != nil {
		for k, v := range headers {
			s.Request.Header.Add(k, v)
		}
	}
	return &s, nil
}

// ReadResponseBody cautiously empties the http.Response
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	b, err := ReadAllLimiter(resp.Body, 5000000)
	_ = resp.Body.Close()
	// ignoring error due to no appreciable impact
	// cf. https://stackoverflow.com/questions/47293975/should-i-error-check-close-on-a-response-body
	return b, err
}

// ReadAllLimiter limits reading of responses to 5MB
func ReadAllLimiter(closer io.ReadCloser, allowedResponseSize int64) ([]byte, error) {
	b, err := io.ReadAll(io.LimitReader(closer, allowedResponseSize))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// DoRequest sends a get request to the tool's configured url
func (c *ApiClient) DoRequest(tr *TracedRequest) (*ApiResponse, error) {

	var a ApiResponse

	r, err := c.Handle.Do(tr.Request)

	if err != nil {
		a.Err = err
		if r != nil {
			a.StatusCode = r.StatusCode
			a.Body, _ = ReadResponseBody(r)
		}
		return &a, err
	}

	a.Body, a.Err = ReadResponseBody(r)

	a.StatusCode = r.StatusCode

	return &a, nil
}

func (c *ApiClient) Query(
	address string,
	response interface{},
	requestHeaders map[string]string,
	insecure bool,
	timeoutSeconds int,
	body []byte,
) (*TracedRequestStats, error) {

	// default timeout will be 10 seconds
	timeout := time.Duration(1000000000 * 10)

	if timeoutSeconds != 0 {
		timeout = time.Duration(1000000000 * timeoutSeconds)
	}
	// control resource usage when running in goroutines
	if insecure {
		c.Handle = http.Client{Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        1,
				MaxConnsPerHost:     1,
				MaxIdleConnsPerHost: 1,
				DisableKeepAlives:   true,
				IdleConnTimeout:     time.Second * 1,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
		}
	} else {
		c.Handle = http.Client{Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        1,
				MaxConnsPerHost:     1,
				MaxIdleConnsPerHost: 1,
				DisableKeepAlives:   true,
				IdleConnTimeout:     time.Second * 1,
			},
		}
	}
	ctx, cancel := context.WithTimeout(context.TODO(), c.Handle.Timeout)
	defer cancel()
	tr, err := NewTracedGetRequestWithContext(ctx, address, requestHeaders, body)
	if err != nil {
		return nil, makeTracedRequestStatsError
	}

	apiResponse, err := c.DoRequest(tr)
	if err != nil {
		if apiResponse != nil && apiResponse.Body != nil {
			fmt.Println("httpclient error body", string(apiResponse.Body))
			return nil, err
		}
		fmt.Println("httpclient error", err.Error())
		return nil, err
	}

	if apiResponse.Body != nil {
		err := json.Unmarshal(apiResponse.Body, &response)
		if err != nil {
			return nil, errors.New("failed_to_unmarshal_response_body=" + string(apiResponse.Body))
		}
		trStats := tr.MakeTracedRequestStats(time.Now(), apiResponse, false)
		return &trStats, nil
	}

	return nil, nilResponseError
}

func (t *TracedRequestStats) String() string {
	ts, e := json.Marshal(t)
	if e != nil {
		return "{\"err\":\"json.Marshal(TracedRequestStats)\"}"
	}
	return string(ts)
}

// TracedRequestStats are the measured outcome of a traced sent into a channel; the values in the TracedRequestStats struct are
// aggregated and descriptive stats are calculated and displayed
type TracedRequestStats struct {
	DnsLookup        float64 `json:"DnsLookup"`
	TcpConn          float64 `json:"TcpConn"`
	TlsHandshake     float64 `json:"TlsHandshake"`
	ServerProcessing float64 `json:"ServerProcessing"`
	ContentTransfer  float64 `json:"ContentTransfer"`
	TotalDuration    float64 `json:"TotalDuration"`
	StatusCode       int     `json:"StatusCode"`
	RequestIdMatch   bool    `json:"RequestIdMatch"`
	Success          bool    `json:"Success"`
}

// MakeTracedRequestStats calculates the durations and stats for the request
// When the goroutine that executes the traced request runs, we send tracedRequestStats, which contain duration statistics
// about the parts of the HTTP request on the output channel for later aggregation:
// https://www.w3.org/Protocols/HTTP-NG/http-prob.html
func (tr *TracedRequest) MakeTracedRequestStats(bodyReadTime time.Time, resp *ApiResponse, isTraceTLS bool) TracedRequestStats {
	if tr.dnsLookupBeginTime.IsZero() {
		tr.dnsLookupBeginTime = tr.dialBeginTime
	}
	out := TracedRequestStats{
		DnsLookup:        getDuration(tr.dnsLookupBeginTime, tr.dialBeginTime),
		ServerProcessing: getDuration(tr.obtainedConnTime, tr.firstResponseByteTime),
		ContentTransfer:  getDuration(tr.firstResponseByteTime, bodyReadTime),
		StatusCode:       resp.StatusCode,
		RequestIdMatch:   resp.RequestId == tr.requestID,
		Success:          isResponseDescriptionSuccess(resp),
	}
	if isTraceTLS {
		out.TcpConn = getDuration(tr.dialBeginTime, tr.dialDoneTime)
		out.TlsHandshake = getDuration(tr.tlsHandshakeStartTime, tr.tlsHandshakeDoneTime)
	} else {
		out.TcpConn = getDuration(tr.dialBeginTime, tr.obtainedConnTime)
	}
	out.TotalDuration = out.DnsLookup + out.TcpConn + out.TlsHandshake + out.ServerProcessing + out.ContentTransfer
	return out
}

// getDuration returns the difference between a start and finish time.Time argument in milliseconds.
func getDuration(start, finish time.Time) float64 {
	return float64(finish.Sub(start).Milliseconds())
}

// isResponseDescriptionSuccess can add custom stuff
func isResponseDescriptionSuccess(resp *ApiResponse) bool {
	if resp.Err == nil || resp.StatusCode == 200 {
		return true
	}
	return false
}

func NewApiClient(insecure bool, timeout time.Duration) *ApiClient {
	tr := &http.Transport{
		MaxIdleConns:        1,
		MaxConnsPerHost:     1,
		MaxIdleConnsPerHost: 1,
		DisableKeepAlives:   true,
		IdleConnTimeout:     1 * time.Second,
	}

	if insecure {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &ApiClient{
		Handle: http.Client{
			Timeout:   timeout,
			Transport: tr,
		},
	}
}
