package monoflake

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("default node bits", func(t *testing.T) {
		tests := []struct {
			nodeID uint16
			epoch  int64
			want   error
		}{
			{0, minEpoch, nil},
			{1024, minEpoch + 100, nil},
			{1025, minEpoch, nil},
			{1024, minEpoch - 1, ErrEpochTooEarly},
		}

		msg := "New(%d, %d) = %v, but returned %v"
		for _, test := range tests {
			epoch := time.Unix(test.epoch, 0)
			_, got := New(test.nodeID, WithEpoch(epoch))
			if got != test.want {
				t.Errorf(msg, test.nodeID, test.epoch, test.want, got)
			}
		}
	})

	t.Run("assigned node bits", func(t *testing.T) {
		tests := []struct {
			nodeID   uint16
			epoch    int64
			nodeBits int
			want     error
		}{
			{0, minEpoch, 8, nil},
			{1024, minEpoch + 100, 10, nil},
			{1025, minEpoch, 12, nil},
			{1024, minEpoch - 1, 8, ErrEpochTooEarly},
			{1024, minEpoch, 7, ErrNodeBitsLowerThanMin},
			{1024, minEpoch, 14, ErrNodeBitsGreaterThanMax},
		}

		msg := "New(%d, %d) = %v, but returned %v"
		for _, test := range tests {
			epoch := time.Unix(test.epoch, 0)
			_, got := New(test.nodeID, WithEpoch(epoch), WithNodeBits(test.nodeBits))
			if got != test.want {
				t.Errorf(msg, test.nodeID, test.epoch, test.want, got)
			}
		}
	})

}

func TestNext(t *testing.T) {
	// test for all possible node bits
	for i := 8; i < 14; i++ {
		mf, _ := New(100, WithNodeBits(i))
		t.Run("generates greater sequences on each call until max sequence", func(t *testing.T) {
			t.Parallel()
			for i := 0; i < int(mf.maxSequence)*250; i++ {
				id1, id2 := mf.Next(), mf.Next()
				if id1 >= id2 {
					debug := map[string]int64{
						"since":     time.Since(mf.epoch).Milliseconds(),
						"max_seq":   mf.maxSequence,
						"time1":     id1.Since(),
						"time2":     id2.Since(),
						"sequence1": id1.Sequence(mf.nodeBits),
						"sequence2": id2.Sequence(mf.nodeBits),
						"node_id1":  id1.NodeID(mf.nodeBits),
						"node_id2":  id2.NodeID(mf.nodeBits),
						"i":         int64(i),
					}
					t.Errorf("Next(): %d >= Next(): %d %v", id1, id2, debug)
				}
			}
		})
	}

}

func TestErrorString(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{ErrEpochTooEarly, "epoch is too earlier than June 1st 2024 UTC"},
		{ErrNodeBitsLowerThanMin, "node bits must be greater than 8"},
		{ErrNodeBitsGreaterThanMax, "node bits must be less than 13"},
	}

	for _, test := range tests {
		if got := test.err.Error(); got != test.want {
			t.Errorf("ErrorString(%v) = %q, but returned %q", test.err, test.want, got)
		}
	}
}
