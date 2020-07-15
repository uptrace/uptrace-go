package internal

import (
	"bytes"
	"sync"

	"github.com/klauspost/compress/s2"
	"github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/kv/value"
)

type KVMap map[kv.Key]value.Value

func (m KVMap) EncodeMsgpack(enc *msgpack.Encoder) error {
	_ = enc.EncodeMapLen(len(m))
	for k, v := range m {
		EncodeKey(enc, k)
		EncodeValue(enc, v)
	}
	return nil
}

//------------------------------------------------------------------------------

type KVSlice []kv.KeyValue

func (slice KVSlice) EncodeMsgpack(enc *msgpack.Encoder) error {
	_ = enc.EncodeMapLen(len(slice))
	for _, el := range slice {
		EncodeKey(enc, el.Key)
		EncodeValue(enc, el.Value)
	}
	return nil
}

//------------------------------------------------------------------------------

func EncodeKey(enc *msgpack.Encoder, k kv.Key) {
	_ = enc.EncodeString(string(k))
}

func EncodeValue(enc *msgpack.Encoder, v value.Value) {
	switch v.Type() {
	case value.BOOL:
		_ = enc.EncodeBool(v.AsBool())
	case value.INT32:
		_ = enc.EncodeInt32(v.AsInt32())
	case value.INT64:
		_ = enc.EncodeInt64(v.AsInt64())
	case value.UINT32:
		_ = enc.EncodeUint32(v.AsUint32())
	case value.UINT64:
		_ = enc.EncodeUint64(v.AsUint64())
	case value.FLOAT32:
		_ = enc.EncodeFloat32(v.AsFloat32())
	case value.FLOAT64:
		_ = enc.EncodeFloat64(v.AsFloat64())
	case value.STRING:
		_ = enc.EncodeString(v.AsString())
	case value.ARRAY:
		_ = enc.Encode(v.AsArray())
	default:
		logrus.WithField("type", v.Type()).Error("unknown type")
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

func (enc *Encoder) Encode(v interface{}) (*bytes.Buffer, error) {
	enc.buf.Reset()
	enc.msgp.Reset(&enc.buf)
	enc.msgp.UseCompactInts(true)

	if err := enc.msgp.Encode(v); err != nil {
		return nil, err
	}
	return &enc.buf, nil
}

func (enc *Encoder) EncodeS2(v interface{}) (*bytes.Buffer, error) {
	s2w := getS2Writer()
	defer putS2Writer(s2w)

	enc.buf.Reset()
	s2w.Reset(&enc.buf)
	enc.msgp.Reset(s2w)
	enc.msgp.UseCompactInts(true)

	if err := enc.msgp.Encode(v); err != nil {
		return nil, err
	}
	if err := s2w.Close(); err != nil {
		return nil, err
	}
	return &enc.buf, nil
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

var s2Pool = sync.Pool{
	New: func() interface{} {
		return s2.NewWriter(nil, s2.WriterConcurrency(1))
	},
}

func getS2Writer() *s2.Writer {
	return s2Pool.Get().(*s2.Writer)
}

func putS2Writer(s2w *s2.Writer) {
	s2Pool.Put(s2w)
}
