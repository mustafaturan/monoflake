package monoflake

import (
	"bytes"
	"encoding/binary"
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

func TestBigEndianBytes(t *testing.T) {
	tests := []struct {
		id   ID
		want []byte
	}{
		{ID(0), []byte{0, 0, 0, 0, 0, 0, 0, 0}},
		{ID(1), []byte{0, 0, 0, 0, 0, 0, 0, 1}},
		{ID(3263530505704195), []byte{0, 11, 152, 41, 232, 129, 147, 3}},
		{ID(3263530505704640), []byte{0, 11, 152, 41, 232, 129, 148, 192}},
	}

	msg := "BigEndianBytes(%d) = %v, but returned %v"
	for _, test := range tests {
		got := test.id.BigEndianBytes()
		if !bytes.Equal(got, test.want) {
			t.Errorf(msg, test.id, test.want, got)
		}
	}
}

func TestBytes(t *testing.T) {
	mf, _ := New(3843)
	id1, id2 := mf.Next(), mf.Next()

	t.Run("generates greater sequences on each call", func(t *testing.T) {
		t.Parallel()
		if bytes.Compare(id1.Bytes()[:], id2.Bytes()[:]) >= 0 {
			t.Errorf("Bytes(): %s >= Bytes(): %s", id1, id2)
		}
	})

	t.Run("generates 11 bytes sequences", func(t *testing.T) {
		t.Parallel()
		results := [][]byte{id1.Bytes()[:], id2.Bytes()[:]}
		for _, r := range results {
			if len(r) != 11 {
				t.Errorf("Bytes(): %s couldn't produce 11 bytes", r)
			}
		}
	})
}

func TestString(t *testing.T) {
	mf, _ := New(3843)
	id1 := mf.Next()
	id2 := mf.Next()

	t.Run("generates greater sequences on each call", func(t *testing.T) {
		t.Parallel()
		if strings.Compare(id1.String(), id2.String()) >= 0 {
			t.Errorf("String(): %s (%d) >= String(): %s (%d)", id1, int64(id1), id2, int64(id2))
		}
	})

	t.Run("generates 11 bytes string sequences", func(t *testing.T) {
		t.Parallel()
		results := []string{id1.String(), id2.String()}
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

func TestIDFromBase62(t *testing.T) {
	tests := []struct {
		base62 string
		want   int64
	}{
		{"00000000001", 1},
		{"11", 63},
		{"020", 124},
		{"0021", 125},
		{"AzL8n0Y58m7", 1<<63 - 1},
		{"ZZZZZZZZZZZ", -1},
		{"ZZZZZZZZZZZZ", -1},
	}

	msg := "IDFromBase62(%s) = %d, but returned %d"
	for _, test := range tests {
		got := IDFromBase62(test.base62).Int64()
		if got != test.want {
			t.Errorf(msg, test.base62, test.want, got)
		}
	}
}

func TestIDFromBigEndianBytes(t *testing.T) {
	tests := []struct {
		bigEndian []byte
		want      int64
	}{
		{[]byte{0, 0, 0, 0, 0, 0, 0, 1}, 1},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 63}, 63},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 124}, 124},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 125}, 125},
		{[]byte{127, 255, 255, 255, 255, 255, 255, 255}, 1<<63 - 1},
		{[]byte{255, 255, 255, 255, 255, 255, 255, 255}, -1},
		{[]byte{255, 255, 255, 255, 255, 255, 255, 255, 255}, -1},
	}

	msg := "IDFromBigEndianBytes(%s) = %d, but returned %d"
	for _, test := range tests {
		got := IDFromBigEndianBytes(test.bigEndian).Int64()
		if got != test.want {
			t.Errorf(msg, test.bigEndian, test.want, got)
		}
	}
}

func TestToBase62WithPaddingZeros(t *testing.T) {
	tests := []struct {
		val  uint64
		want string
	}{
		{1, "00000000001"},
		{63, "00000000011"},
		{124, "00000000020"},
		{125, "00000000021"},
		{1<<63 - 1, "AzL8n0Y58m7"},
		{1<<64 - 1, "LygHa16AHYF"},
	}

	msg := "toBase62WithPaddingZeros(%d, %d) = %v, but returned %v"
	for _, test := range tests {
		got := toBase62WithPaddingZeros(test.val)
		if string(got[:]) != test.want {
			t.Errorf(msg, test.val, test.want, string(got[:]))
		}
	}
}

func TestFromBase62RuneToInt64(t *testing.T) {
	var want, got int64
	msg := "fromBase62RuneToInt64(%s) = %d, but returned %d"
	for _, r := range base62Mapping {
		got = fromBase62RuneToInt64(r)
		want = int64(strings.IndexRune(base62Mapping, r))
		if got != want {
			t.Errorf(msg, string(r), want, got)
		}
	}

	got = fromBase62RuneToInt64('!')
	if got != want {
		t.Errorf(msg, string('!'), want, got)
	}
}

func TestFromInt64ToBytes(t *testing.T) {
	tests := []struct {
		want []byte
		val  int64
	}{
		{[]byte{0, 0, 0, 0, 0, 0, 0, 1}, 1},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 63}, 63},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 124}, 124},
		{[]byte{0, 0, 0, 0, 0, 0, 0, 125}, 125},
		{[]byte{127, 255, 255, 255, 255, 255, 255, 255}, 1<<63 - 1},
		{[]byte{255, 255, 255, 255, 255, 255, 255, 255}, -1},
	}

	msg := "IDFromBigEndianBytes(%s) = %d, but returned %d"
	for _, test := range tests {
		got := fromInt64ToBytes(test.val)
		if !bytes.Equal(got[:], test.want) {
			t.Errorf(msg, test.val, test.want, got)
		}

		buf := make([]byte, 0, 8)
		comparable := binary.BigEndian.AppendUint64(buf[0:], uint64(test.val))
		if !bytes.Equal(got[:], comparable) {
			t.Errorf(msg, test.val, comparable, got)
		}
	}
}
