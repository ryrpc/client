package rycli

import (
	"fmt"
	"math"

	"github.com/fxamacker/cbor/v2"
	"github.com/golang/protobuf/proto"
	"github.com/valyala/fasthttp"
	"lukechampine.com/frand"
)

// func printObject(v interface{}) string {
// 	res2B, _ := ffjson.Marshal(v)
// 	return string(res2B)
// }

// encodeClientRequest encodes parameters for a JSON-RPC client request.
func encodeClientRequest(method string, args interface{}) ([]byte, error) {

	arg := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(arg)

	arg.Add("version", "2.0")
	arg.Add("method", method)
	arg.Add("id", fmt.Sprintf("%d", frand.Uint64n(math.MaxUint64)))

	if val, ok := args.(string); ok {
		arg.Add("params", val)
	} else {
		b, err := cbor.Marshal(args)
		if err != nil {
			return b, err
		}
		arg.AddBytesV("params", b)
	}

	qs := arg.QueryString()

	return qs, nil
}

// decodeClientResponse decodes the response body of a client request into the interface reply.
func decodeClientResponse(method string, r []byte, result interface{}) error {

	arg := &Base{}

	err := proto.Unmarshal(r, arg)
	if err != nil {
		return err
	}

	if len(arg.GetErr()) > 0 {
		err1 := fmt.Errorf("rpc call %s on could not decode body to rpc error: %s", method, arg.GetErr())
		return err1
	}
	/*
		if !arg.Has("result") {
			err1 := fmt.Errorf("rpc call %s on could not decode body to rpc response: not found", method)
			return err1
		}
	*/
	if _, ok := result.(string); ok {
		result = string(arg.GetData())
	} else {
		err = cbor.Unmarshal(arg.GetData(), result)
		if err != nil {
			err1 := fmt.Errorf("rpc call %s() on could not decode body to rpc Decode: %s", method, err.Error())
			return err1
		}
	}

	return nil
}
