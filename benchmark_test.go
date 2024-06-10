package monoflake_test

import (
	"testing"
	"time"

	"github.com/mustafaturan/monoflake"
)

func BenchmarkNext(b *testing.B) {
	b.ReportAllocs()

	mf, _ := monoflake.New(0)
	for n := 0; n < b.N; n++ {
		_ = mf.Next()
	}
}

func BenchmarkNextCompare(b *testing.B) {
	b.ReportAllocs()

	mf, err := monoflake.New(0)
	if err != nil {
		b.Fatal(err)
	}
	var id1, id2 monoflake.ID
	for n := 0; n < b.N; n++ {
		id1, id2 = mf.Next(), mf.Next()
		if id1 > id2 {
			b.Fatalf("Next(): %d >= Next(): %d", id1, id2)
		}
	}
}

func BenchmarkNextBase62(b *testing.B) {
	b.ReportAllocs()

	mf, _ := monoflake.New(0)
	for n := 0; n < b.N; n++ {
		_ = mf.Next().String()
	}
}

func BenchmarkNextBytes(b *testing.B) {
	b.ReportAllocs()

	mf, _ := monoflake.New(0, monoflake.WithEpoch(time.Now()))
	for n := 0; n < b.N; n++ {
		_ = mf.Next().Bytes()
	}
}
