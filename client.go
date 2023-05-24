package rycli

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

func getDefaultHeadersMap() map[string]string {
	headers := make(map[string]string)
	headers["User-Agent"] = userAgent
	headers["Content-Type"] = defaultContentType
	return headers
}

func createNewClient() *Client {
	return &Client{
		clientPool: &sync.Pool{
			New: func() interface{} {
				return new(fasthttp.Client)
			},
		},
		customHeaders: getDefaultHeadersMap(),
	}
}

// NewClient returns new configured Client to start work with JSON-RPC 2.0 protocol
func NewClient() *Client {
	return createNewClient()
}

// SetBaseURL setting basic url for API
func (cl *Client) SetBaseURL(baseURL string) {
	cl.BaseURL = baseURL
}

// DisableHeaderNamesNormalizing setting normalize headers or not
func (cl *Client) DisableHeaderNamesNormalizing(fix bool) {
	cl.disableHeaderNamesNormalizing = fix
}

// SetClientTimeout this method sets globally for client its timeout
func (cl *Client) SetClientTimeout(duration time.Duration) {
	cl.clientTimeout = duration
}

// SetCustomHeader setting custom header
func (cl *Client) SetCustomHeader(headerName string, headerValue string) {
	cl.customHeaders[headerName] = headerValue
}

// DeleteCustomHeader delete custom header
func (cl *Client) DeleteCustomHeader(headerName string) {
	delete(cl.customHeaders, headerName)
}

// SetBasicAuthHeader setting basic auth header
func (cl *Client) SetBasicAuthHeader(login string, password string) {
	cl.SetCustomAuthHeader("Basic", base64.StdEncoding.EncodeToString([]byte(login+":"+password)))
}

// SetCustomAuthHeader setting custom auth header with type of auth and auth data
func (cl *Client) SetCustomAuthHeader(authType string, authData string) {
	cl.SetCustomHeader("Authorization", authType+" "+authData)
}

// DeleteAuthHeader clear basic auth header
func (cl *Client) DeleteAuthHeader() {
	cl.DeleteCustomHeader("Authorization")
}

// SetUserAgent setting custom User Agent header
func (cl *Client) SetUserAgent(userAgent string) {
	cl.SetCustomHeader("User-Agent", userAgent)
}

func (cl *Client) makeCallRequest(method string, args interface{}) ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer req.Reset()
	req.SetRequestURI(cl.BaseURL + "/" + method)

	cl.SetCustomHeader("X-Func", method)
	for key, val := range cl.customHeaders {
		req.Header.Set(key, val)
	}

	req.Header.SetMethod("POST")
	byteBody, err := encodeClientRequest(method, args)
	if err != nil {
		return nil, 0, err
	}

	//debugLogging(cl, logrus.Fields{"headers": req.Header.String(), "request": byteBody}, "request prepared")

	req.SetBody(byteBody)
	resp := fasthttp.AcquireResponse()
	defer resp.Reset()

	client := cl.clientPool.Get().(*fasthttp.Client)

	client.DisableHeaderNamesNormalizing = cl.disableHeaderNamesNormalizing

	if cl.clientTimeout == 0 {
		if err := client.Do(req, resp); err != nil {
			return nil, 0, err
		}
	} else {
		if err := client.DoTimeout(req, resp, cl.clientTimeout); err != nil {
			return nil, 0, err
		}
	}

	cl.clientPool.Put(client)

	statusCode := resp.StatusCode()
	if statusCode != 200 {
		err = fmt.Errorf("rpc call %s() status code: %d. could not decode body to rpc response: %s", method, statusCode, err.Error())
		return nil, 0, err
	}

	return resp.SwapBody(nil), statusCode, nil
}

// Call run remote procedure on JSON-RPC 2.0 API with parsing answer to provided structure or interface
func (cl *Client) Call(method string, args, result interface{}) error {

	resp, _, err := cl.makeCallRequest(method, args)
	//fmt.Println("Call = ", string(resp))
	if err != nil {
		return err
	}
	err = decodeClientResponse(method, resp, result)
	return err
}

/*
// CallForMap run remote procedure on JSON-RPC 2.0 API with returning map[string]interface{}
func (cl *Client) CallForMap(urlPath string, method string, args interface{}) (map[string]interface{}, error) {
	resp, statusCode, err := cl.makeCallRequest(urlPath, method, args)
	if err != nil {
		return nil, err
	}
	dst := make(map[string]interface{})
	err = decodeClientResponse(resp, &dst)
	return dst, err
}
*/
/*
func (cl *Client) CallBatch(urlPath string, method string, args interface{}) {

}

func (cl *Client) AsyncCall(urlPath string, method string, args interface{}, ch chan<- interface{}) {
	var result interface{}
	ch <- result
}
*/
