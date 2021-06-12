package jsonex

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	Marshal         = json.Marshal
	MarshalToString = json.MarshalToString
	MarshalIndent   = json.MarshalIndent

	Unmarshal           = json.Unmarshal
	UnmarshalFromString = json.UnmarshalFromString

	Get = json.Get

	NewEncoder = json.NewEncoder
	NewDecoder = json.NewDecoder

	Valid             = json.Valid
	RegisterExtension = json.RegisterExtension
	DecoderOf         = json.DecoderOf
	EncoderOf         = json.EncoderOf
)

func MustMarshalToString(data interface{}) string {
	str, err := MarshalToString(data)
	if err != nil {
		panic(err)
	}
	return str
}
