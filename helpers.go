package rycli

import (
	"fmt"
	"math"

	"lukechampine.com/frand"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/decoder"
	"github.com/valyala/fastjson"
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
func decodeClientResponse(method string, r []byte, result interface{}) error {

	val, err := fastjson.ParseBytes(r)
	if err != nil {
		err1 := fmt.Errorf("rpc call %s on could not decode body to rpc ParseBytes: %s", method, err.Error())
		return err1
	}

	if val.Exists("error") {
		err1 := fmt.Errorf("rpc call %s on could not decode body to rpc error: %s", method, string(val.GetStringBytes("error")))
		return err1
	}

	if !val.Exists("result") {
		err1 := fmt.Errorf("rpc call %s on could not decode body to rpc response: not found", method)
		return err1
	}

	ss := string(val.GetStringBytes("result"))

	decoder := decoder.NewDecoder(ss)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	err = decoder.Decode(result)

	if err != nil {
		err1 := fmt.Errorf("rpc call %s() on could not decode body to rpc Decode: %s", method, err.Error())
		return err1
	}

	return nil
}
