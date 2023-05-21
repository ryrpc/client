package rycli

import (
	"fmt"
	"math"

	"lukechampine.com/frand"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
)

// func printObject(v interface{}) string {
// 	res2B, _ := ffjson.Marshal(v)
// 	return string(res2B)
// }

// encodeClientRequest encodes parameters for a JSON-RPC client request.
func encodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		ID:      frand.Uint64n(math.MaxUint64),
		Version: "2.0",
		Method:  method,
		Params:  args,
	}
	return sonic.Marshal(c)
}

// decodeClientResponse decodes the response body of a client request into the interface reply.
func decodeClientResponse(method string, r []byte, statusCode int) (SrvResponse, error) {

	var res SrvResponse

	decoder := decoder.NewDecoder(string(r))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	err := decoder.Decode(&res)

	if err != nil {
		err := fmt.Errorf("rpc call %s() on could not decode body to rpc response: %s", method, err.Error())
		return res, err
	}

	return res, nil
}
