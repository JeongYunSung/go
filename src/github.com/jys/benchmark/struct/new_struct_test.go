package _struct

import (
	"fmt"
	"os"
	"runtime/trace"
	"testing"
)

type S struct {
	a, b, c int64
	d, e, f string
}

func byCopy() S {
	return S{
		a: 1, b: 2, c: 3,
		d: "Jeong", e: "Yun", f: "Sung",
	}
}

func byPointer() *S {
	return &S{
		a: 1, b: 2, c: 3,
		d: "Jeong", e: "Yun", f: "Sung",
	}
}

func (s S) stack(s1 S) {}

func (s *S) heap(s1 *S) {}

func BenchmarkMemoryStack2(b *testing.B) {
	var s S
	var s1 S

	s = byCopy()
	s1 = byCopy()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100000; i++ {
			s.stack(s1)
		}
	}
}

func BenchmarkMemoryHeap2(b *testing.B) {
	var s *S
	var s1 *S

	s = byPointer()
	s1 = byPointer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 100000; i++ {
			s.heap(s1)
		}
	}
}

func BenchmarkMemoryStack1(b *testing.B) {
	var s S

	f, err := os.Create("stack.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		s = byCopy()
	}

	trace.Stop()

	b.StopTimer()

	_ = fmt.Sprintf("%v", s.a)
}

func BenchmarkMemoryHeap1(b *testing.B) {
	var s *S

	f, err := os.Create("heap.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}

	for i := 0; i < b.N; i++ {
		s = byPointer()
	}

	trace.Stop()

	b.StopTimer()

	_ = fmt.Sprintf("%v", s.a)
}
