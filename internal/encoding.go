package internal

import (
	"bytes"
	"context"
	"sync"

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

type Encoder struct {
	buf  bytes.Buffer
	msgp *msgpack.Encoder
}

func NewEncoder() *Encoder {
	var enc Encoder
	enc.msgp = msgpack.NewEncoder(nil)
	return &enc
}

func (enc *Encoder) Encode(v interface{}) ([]byte, error) {
	enc.buf.Reset()
	enc.msgp.Reset(&enc.buf)
	enc.msgp.UseCompactInts(true)

	if err := enc.msgp.Encode(v); err != nil {
		return nil, err
	}
	return enc.buf.Bytes(), nil
}

func (enc *Encoder) EncodeZstd(v interface{}) ([]byte, error) {
	zw := getZstdWriter()
	defer putZstdWriter(zw)

	enc.buf.Reset()
	zw.Reset(&enc.buf)
	enc.msgp.Reset(zw)
	enc.msgp.UseCompactInts(true)

	if err := enc.msgp.Encode(v); err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return enc.buf.Bytes(), nil
}

var encPool = sync.Pool{
	New: func() interface{} {
		return NewEncoder()
	},
}

func GetEncoder() *Encoder {
	return encPool.Get().(*Encoder)
}

func PutEncoder(enc *Encoder) {
	enc.buf.Reset()
	encPool.Put(enc)
}

var zstdPool = sync.Pool{
	New: func() interface{} {
		zw, err := zstd.NewWriter(nil, zstd.WithEncoderConcurrency(1))
		if err != nil {
			panic(err)
		}
		return zw
	},
}

func getZstdWriter() *zstd.Encoder {
	return zstdPool.Get().(*zstd.Encoder)
}

func putZstdWriter(zw *zstd.Encoder) {
	zstdPool.Put(zw)
}
