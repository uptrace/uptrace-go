package internal

import (
	"bytes"
	"context"

	"github.com/klauspost/compress/zstd"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/attribute"
)

type KVMap map[attribute.Key]attribute.Value

func (m KVMap) EncodeMsgpack(enc *msgpack.Encoder) error {
	_ = enc.EncodeMapLen(len(m))
	for k, v := range m {
		EncodeKey(enc, k)
		EncodeValue(enc, v)
	}
	return nil
}

//------------------------------------------------------------------------------

type KeyValueSlice []attribute.KeyValue

var _ msgpack.CustomEncoder = (*KeyValueSlice)(nil)

func (slice KeyValueSlice) EncodeMsgpack(enc *msgpack.Encoder) error {
	if len(slice) == 0 {
		return enc.EncodeNil()
	}

	_ = enc.EncodeMapLen(len(slice))
	for _, el := range slice {
		EncodeKey(enc, el.Key)
		EncodeValue(enc, el.Value)
	}
	return nil
}

var _ msgpack.CustomDecoder = (*KeyValueSlice)(nil)

func (slice *KeyValueSlice) DecodeMsgpack(dec *msgpack.Decoder) error {
	n, err := dec.DecodeMapLen()
	if err != nil {
		return err
	}

	if n == -1 {
		*slice = nil
		return nil
	}

	*slice = make(KeyValueSlice, n)

	for i := 0; i < n; i++ {
		key, err := dec.DecodeString()
		if err != nil {
			return err
		}

		val, err := dec.DecodeInterface()
		if err != nil {
			return err
		}

		(*slice)[i] = attribute.Any(key, val)
	}

	return nil
}

//------------------------------------------------------------------------------

func EncodeKey(enc *msgpack.Encoder, k attribute.Key) {
	_ = enc.EncodeString(string(k))
}

func EncodeValue(enc *msgpack.Encoder, v attribute.Value) {
	switch v.Type() {
	case attribute.BOOL:
		_ = enc.EncodeBool(v.AsBool())
	case attribute.INT64:
		_ = enc.EncodeInt64(v.AsInt64())
	case attribute.FLOAT64:
		_ = enc.EncodeFloat64(v.AsFloat64())
	case attribute.STRING:
		_ = enc.EncodeString(v.AsString())
	case attribute.ARRAY:
		_ = enc.Encode(v.AsArray())
	default:
		Logger.Printf(context.TODO(), "unknown otel type: %s", v.Type())
		_ = enc.EncodeString("unknown otel type: " + v.Type().String())
	}
}

//------------------------------------------------------------------------------

func EncodeMsgpack(v interface{}) ([]byte, error) {
	enc := msgpack.GetEncoder()
	defer msgpack.PutEncoder(enc)

	var buf bytes.Buffer
	enc.Reset(&buf)

	if err := enc.Encode(v); err != nil {
		return nil, err
	}

	zenc, err := zstdEncoder()
	if err != nil {
		return nil, err
	}

	return zenc.EncodeAll(buf.Bytes(), nil), nil
}

var (
	zencOnce Once
	zenc     *zstd.Encoder
)

func zstdEncoder() (*zstd.Encoder, error) {
	if err := zencOnce.Do(func() error {
		var err error
		zenc, err = zstd.NewWriter(nil)
		return err
	}); err != nil {
		return nil, err
	}
	return zenc, nil
}
