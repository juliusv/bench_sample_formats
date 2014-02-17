package io_prometheus

import (
	"bytes"
	"encoding/gob"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"sync"
	"testing"
)

const numSamples = 5000

var once sync.Once

type Sample struct {
	Timestamp int64
	Value     float64
}

type SampleValueSeries []Sample

func prepareBuf(b *testing.B) *bytes.Buffer {
	s := make(SampleValueSeries, 0, numSamples)
	for i := 0; i < numSamples; i++ {
		s = append(s, Sample{
			Timestamp: rand.Int63(),
			Value:     rand.NormFloat64(),
		})
	}

	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(s); err != nil {
		b.Fatal(err)
	}
	return &buf
}

func BenchmarkUnmarshal(b *testing.B) {
	b.StopTimer()

	// BenchmarkUnmarshal is called multiple times.
	once.Do(func() {
		go http.ListenAndServe("localhost:9090", nil)
	})

	raw := prepareBuf(b).Bytes()

	for i := 0; i < b.N; i++ {
		buf := bytes.NewBuffer(raw)
		dec := gob.NewDecoder(buf)

		b.StartTimer()
		s := make(SampleValueSeries, 0, numSamples)
		if err := dec.Decode(&s); err != nil {
			b.Fatal(err)
		}
		b.StopTimer()

		if len(s) != numSamples {
			b.Fatal(len(s))
		}
	}
}
