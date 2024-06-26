package rycli

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

func createNewClient() *Client {
	return &Client{
		clientTimeout: 12 * time.Second,
		clientPool: &sync.Pool{
			New: func() interface{} {
				return new(fasthttp.Client)
			},
		},
		//customHeaders: getDefaultHeadersMap(),
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

/*
// SetCustomHeader setting custom header
func (cl *Client) SetCustomHeader(headerName, headerValue string) {

	cl.customHeaders.Store(headerName, headerValue)
	//cl.customHeaders[headerName] = headerValue
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
*/
/*
// DeleteAuthHeader clear basic auth header
func (cl *Client) DeleteAuthHeader() {
	cl.DeleteCustomHeader("Authorization")
}
*/

func (cl *Client) makeCallRequest(method string, args interface{}, headers map[string]string) ([]byte, int, error) {
	req := fasthttp.AcquireRequest()
	defer req.Reset()
	req.SetRequestURI(cl.BaseURL + method)

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", defaultContentType)

	name := strings.SplitN(method, "/", 5)
	if len(name) > 1 {

		req.Header.Set("func", name[1])

		for k, v := range headers {
			req.Header.Set(k, v)
		}
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

	cli := cl.clientPool.Get().(*fasthttp.Client)

	if cl.clientTimeout == 0 {
		if err := cli.Do(req, resp); err != nil {
			return nil, 0, err
		}
	} else {
		if err := cli.DoTimeout(req, resp, cl.clientTimeout); err != nil {
			return nil, 0, err
		}
	}

	statusCode := resp.StatusCode()
	if statusCode != 200 {
		err = fmt.Errorf("rpc call %s() status code: %d.", method, statusCode)
		return nil, 0, err
	}

	return resp.SwapBody(nil), statusCode, nil
}

// Call run remote procedure on JSON-RPC 2.0 API with parsing answer to provided structure or interface
func (cl *Client) Call(method string, headers map[string]string, args, result interface{}) error {

	resp, _, err := cl.makeCallRequest(method, args, headers)
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
