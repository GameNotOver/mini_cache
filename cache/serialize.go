package cache

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
)

type Serializer interface {
	Serialize(src interface{}) ([]byte, error)
	Deserialize(s []byte, ptr interface{}) error
}

type jsonSerializer struct{}

var JsonSerializer jsonSerializer

func (js jsonSerializer) Serialize(src interface{}) ([]byte, error) {
	return json.Marshal(src)
}

func (js jsonSerializer) Deserialize(s []byte, ptr interface{}) error {
	return json.Unmarshal(s, ptr)
}

type compressionSerializer struct{}

var CompressionSerializer compressionSerializer

func (m compressionSerializer) Serialize(src interface{}) ([]byte, error) {
	jsonBytes, err := jsoniter.Marshal(src)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBuffer(nil)
	writer := gzip.NewWriter(buff)
	if _, err = writer.Write(jsonBytes); err != nil {
		return nil, err
	}
	if err = writer.Close(); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil

}

func (m compressionSerializer) Deserialize(s []byte, ptr interface{}) error {
	reader, err := gzip.NewReader(bytes.NewBuffer(s))
	if err != nil {
		return err
	}
	if err = jsoniter.NewDecoder(reader).Decode(ptr); err != nil {
		return err
	}
	return nil
}
