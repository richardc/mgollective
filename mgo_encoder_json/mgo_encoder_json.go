package mgo_encoder_json

import (
	"encoding/json"
	"github.com/richardc/mgollective/mgollective"
)

type JsonEncoder struct {
}

func (e JsonEncoder) Name() string {
	return "json"
}

func (e JsonEncoder) Order() int {
	return 5
}

func (e JsonEncoder) EncodeRequest(message mgollective.RequestMessage) []byte {
	bytes, _ := json.Marshal(message)
	return bytes
}

func (e JsonEncoder) DecodeRequest(bytes []byte) mgollective.RequestMessage {
	message := &mgollective.RequestMessage{}
	json.Unmarshal(bytes, message)
	return *message
}

func makeJsonEncoder(app *mgollective.Mgollective) mgollective.Encoder {
	return &JsonEncoder{}
}

func init() {
	mgollective.RegisterEncoder("json", makeJsonEncoder)
}
