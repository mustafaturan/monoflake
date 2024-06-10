package monoflake

import (
	"bytes"
	"strings"
	"testing"
)

func TestInt64(t *testing.T) {
	tests := []struct {
		id   ID
		want int64
	}{
		{ID(0), 0},
		{ID(1), 1},
		{ID(3263530505704195), 3263530505704195},
		{ID(3263530505704640), 3263530505704640},
	}

	msg := "Int64(%d) = %d, but returned %d"
	for _, test := range tests {
		got := test.id.Int64()
		if got != test.want {
			t.Errorf(msg, test.id, test.want, got)
		}
	}
}

func TestBytes(t *testing.T) {
	tid, _ := New(3843)
	tid1, tid2 := tid.Next(), tid.Next()

	t.Run("generates greater sequences on each call", func(t *testing.T) {
		t.Parallel()
		if bytes.Compare(tid1.Bytes()[:], tid2.Bytes()[:]) >= 0 {
			t.Errorf("Bytes(): %s >= Bytes(): %s", tid1, tid2)
		}
	})

	t.Run("generates 11 bytes sequences", func(t *testing.T) {
		t.Parallel()
		results := [][]byte{tid1.Bytes()[:], tid2.Bytes()[:]}
		for _, r := range results {
			if len(r) != 11 {
				t.Errorf("Bytes(): %s couldn't produce 11 bytes", r)
			}
		}
	})
}

func TestString(t *testing.T) {
	tid, _ := New(3843)
	tid1 := tid.Next()
	tid2 := tid.Next()

	t.Run("generates greater sequences on each call", func(t *testing.T) {
		t.Parallel()
		if strings.Compare(tid1.String(), tid2.String()) >= 0 {
			t.Errorf("String(): %s (%d) >= String(): %s (%d)", tid1, int64(tid1), tid2, int64(tid2))
		}
	})

	t.Run("generates 11 bytes string sequences", func(t *testing.T) {
		t.Parallel()
		results := []string{tid1.String(), tid2.String()}
		for _, r := range results {
			if len(r) != 11 {
				t.Errorf("String(): %s couldn't produce 11 bytes string", r)
			}
		}
	})
}

func TestSince(t *testing.T) {
	tests := []struct {
		id   ID
		want int64
	}{
		{ID(0), 0},
		{ID(1), 0},
		{ID(778086306<<reservedAllocationBits | 1<<defaultReservedNodeBits | 771), 778086306},
		{ID(101<<reservedAllocationBits | 1<<defaultReservedNodeBits | 771), 101},
	}

	msg := "Since(%d) = %d, but returned %d"
	for _, test := range tests {
		got := test.id.Since()
		if got != test.want {
			t.Errorf(msg, test.id, test.want, got)
		}
	}
}

func TestSequence(t *testing.T) {
	tests := []struct {
		id   ID
		want int64
	}{
		{ID(0), 0},
		{ID(1), 0},
		{ID(778086306<<reservedAllocationBits | 100<<defaultReservedNodeBits | 771), 100},
		{ID(778086306<<reservedAllocationBits | 101<<defaultReservedNodeBits | 771), 101},
	}

	msg := "Sequence(%d, %d) = %d, but returned %d"
	for _, test := range tests {
		got := test.id.Sequence(defaultReservedNodeBits)
		if got != test.want {
			t.Errorf(msg, test.id, defaultReservedNodeBits, test.want, got)
		}
	}
}

func TestNodeID(t *testing.T) {
	tests := []struct {
		id   ID
		want int64
	}{
		{ID(0), 0},
		{ID(1), 1},
		{ID(778086306<<reservedAllocationBits | 100<<defaultReservedNodeBits | 771), 771},
		{ID(778086306<<reservedAllocationBits | 101<<defaultReservedNodeBits | 192), 192},
	}

	msg := "NodeID(%d, %d) = %d, but returned %d"
	for _, test := range tests {
		got := test.id.NodeID(defaultReservedNodeBits)
		if got != test.want {
			t.Errorf(msg, test.id, defaultReservedNodeBits, test.want, got)
		}
	}
}

func TestToBase62WithPaddingZeros(t *testing.T) {
	tests := []struct {
		val     uint64
		padding int
		want    string
	}{
		{1, 11, "00000000001"},
		{63, 2, "11"},
		{124, 3, "020"},
		{125, 4, "0021"},
		{1<<64 - 1, 11, "LygHa16AHYF"},
	}

	msg := "toBase62WithPaddingZeros(%d, %d) = %v, but returned %v"
	for _, test := range tests {
		got := toBase62WithPaddingZeros(test.val, test.padding)
		if string(got) != test.want {
			t.Errorf(msg, test.val, test.padding, test.want, string(got))
		}
	}
}
