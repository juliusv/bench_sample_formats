package io_prometheus

import (
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"testing"

	"code.google.com/p/goprotobuf/proto"
)

const numSamples = 5000

var once sync.Once

func prepareBuf(b *testing.B) *proto.Buffer {
	buf := proto.NewBuffer(make([]byte, 0, 4096))
	v := &SampleValueSeries{Value: make([]*SampleValueSeries_Value, 0, numSamples)}
	for i := 0; i < numSamples; i++ {
		v.Value = append(v.Value, &SampleValueSeries_Value{
			Timestamp: proto.Int64(rand.Int63()),
			Value:     proto.Float64(rand.NormFloat64()),
		})
	}
	if err := buf.Marshal(v); err != nil {
		b.Fatal(err)
	}
	return buf
}

func BenchmarkUnmarshal(b *testing.B) {
	b.StopTimer()

	// BenchmarkUnmarshal is called multiple times.
	once.Do(func() {
		go http.ListenAndServe("localhost:9090", nil)
	})

	raw := prepareBuf(b).Bytes()
	buf := proto.NewBuffer(make([]byte, 0, 4096))

	for i := 0; i < b.N; i++ {
		buf.SetBuf(raw)

		b.StartTimer()
		v := &SampleValueSeries{Value: make([]*SampleValueSeries_Value, 0, numSamples)}
		if err := buf.Unmarshal(v); err != nil {
			b.Fatal(err)
		}
		b.StopTimer()

		if len(v.Value) != numSamples {
			b.Fatal(len(v.Value))
		}
	}
}
