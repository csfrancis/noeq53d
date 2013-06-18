package main

import (
	"bytes"
	"testing"
	"sync"
)

const (
	numIds = uint8(32)
)

type nopWriter struct{}

func (nw *nopWriter) Write(b []byte) (n int, err error) {
	return
}

func BenchmarkIdGeneration(b *testing.B) {
	var x int64
	l := new(sync.Mutex)
	for n := 0; n < b.N; n++ {
		nextId(&x, l)
	}
}

func BenchmarkServe01(b *testing.B) {
	for n := 0; n < (b.N / int(numIds)); n++ {
		for m := uint8(0); m < numIds; m++ {
			b.StopTimer()
			i, o := bytes.NewBuffer([]byte{1, byte(m)}), new(nopWriter)
			b.StartTimer()
			serve(i, o)
		}
	}
}

func BenchmarkServe02(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{2, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}

func BenchmarkServe03(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{3, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}

func BenchmarkServe05(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{5, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}

func BenchmarkServe08(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{8, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}

func BenchmarkServe13(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{13, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}

func BenchmarkServe21(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{21, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}

func BenchmarkServe34(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{34, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}

func BenchmarkServe55(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		i, o := bytes.NewBuffer([]byte{55, 0}), new(nopWriter)
		b.StartTimer()
		serve(i, o)
	}
}
