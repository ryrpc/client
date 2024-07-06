package rycli

import (
	"errors"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

// func printObject(v interface{}) string {
// 	res2B, _ := ffjson.Marshal(v)
// 	return string(res2B)
// }

// encodeClientRequest encodes parameters for a JSON-RPC client request.
func encodeClientRequest(method string, args interface{}) ([]byte, error) {

	if val, ok := args.(string); ok {
		return []byte(val), nil
	} else if val, ok := args.([]byte); ok {
		return val, nil
	} else {
		b, err := cbor.Marshal(args)
		if err != nil {
			return b, err
		}
		return b, nil
	}
}

// decodeClientResponse decodes the response body of a client request into the interface reply.
func decodeClientResponse(method string, r []byte, result interface{}) error {

	arg := &Base{}

	_, err := arg.Unmarshal(r)
	if err != nil {
		return err
	}

	if len(arg.Err) > 0 {
		//err1 := fmt.Errorf("rpc call %s on rpc error: %s", method, arg.GetErr())
		return errors.New(arg.Err)
	}
	/*
		if !arg.Has("result") {
			err1 := fmt.Errorf("rpc call %s on could not decode body to rpc response: not found", method)
			return err1
		}
	*/
	if vv, ok := result.(*string); ok {
		*vv = string(arg.Data)
	} else if vv, ok := result.(*[]byte); ok {
		*vv = arg.Data
	} else {
		err = cbor.Unmarshal(arg.Data, result)
		if err != nil {
			err1 := fmt.Errorf("rpc call %s() on could not decode body to rpc Decode: %s", method, err.Error())
			return err1
		}
	}

	return nil
}
