// Copyright 2021 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

/*
Package monoflake is a highly scalable, single/multi node, human-readable,
predictable and incremental 64 bits (8bytes) unique id generator

# Time ordered

The `monoflake` package provides sequences based on the `monotonic` time which
represents the absolute elapsed wall-clock time since some arbitrary, fixed
point in the past. It isn't affected by changes in the system time-of-day clock.

# Epoch time

Epoch time value opens space for time value by subtracting the given value from the time sequence.

# Readable

It comes with `String()` method which encodes the ids into base62 as string and allows padded with zeros to 11 bytes.

# Ready to use bytes

`Bytes()` method allows converting the id directly into static 11 bytes.

# Multi Node Support

The `monoflake` package can be used on single/multiple nodes without the need for
machine coordination. It uses configured node identifier to generate ids by
attaching the node identifier to the end of the sequences.

# Zero allocation

Zero allocation magic with blazing fast results.
*/
package monoflake

import (
	"sync"
	"time"
)

type (
	MonoFlake struct {
		nodeID      int64
		maxSequence int64
		epoch       time.Time
		millisec    int64
		sequence    int64
		nodeBits    int64
		mu          *sync.Mutex
	}

	err string

	Option func(*MonoFlake) error
)

func (e err) Error() string {
	return string(e)
}

const (
	ErrEpochTooEarly          err = "epoch is too earlier than June 1st 2024 UTC"
	ErrNodeBitsLowerThanMin   err = "node bits must be greater than 8"
	ErrNodeBitsGreaterThanMax err = "node bits must be less than 13"

	minEpoch        int64 = 1717200000 // June 1st 2024 UTC
	minSequenceBits int64 = 10         // min 1024
	minNodeBits     int64 = 8          // min 256

	totalBits                   int64 = 64
	reservedSignBits            int64 = 1
	reservedEpochBits           int64 = 40 // 40 bits for milliseconds since epoch
	reservedAllocationBits      int64 = totalBits - reservedSignBits - reservedEpochBits
	defaultReservedNodeBits     int64 = 10
	defaultReservedSequenceBits int64 = reservedAllocationBits - defaultReservedNodeBits
	defaultMaxSequence          int64 = 2 << (defaultReservedSequenceBits - 1)
)

// WithMaxSequenceBits sets the maximum number of bits for node identifier and reserves the rest for sequence number
func WithNodeBits(bits int) Option {
	nodeBits := int64(bits)
	return func(t *MonoFlake) error {
		if nodeBits < minNodeBits {
			return ErrNodeBitsLowerThanMin
		}
		if nodeBits > reservedAllocationBits-minSequenceBits {
			return ErrNodeBitsGreaterThanMax
		}
		t.nodeBits = nodeBits
		t.maxSequence = 2 << (reservedAllocationBits - nodeBits - 1)
		return nil
	}
}

// WithEpoch sets the epoch time for the generator
func WithEpoch(epoch time.Time) Option {
	return func(t *MonoFlake) error {
		if epoch.Unix() < minEpoch {
			return ErrEpochTooEarly
		}
		t.epoch = epoch
		return nil
	}
}

// WithNodeID sets the node identifier for the generator
func withNodeID(nodeID int64) Option {
	return func(t *MonoFlake) error {
		t.nodeID = nodeID % (2 << (t.nodeBits - 1))
		return nil
	}
}

/*
	Default setup:
	| 1 bit (reserved) | 40 bits (since epoch) | 13 bits (sequencer) | 10 bits (node id) |
	| 0                | [0, 1099511627776)    | [0-8192)            | [0, 1024)         |
*/
// New creates a new MonoFlake generator
func New(nodeID uint16, opts ...Option) (*MonoFlake, error) {
	epoch := time.Unix(minEpoch, 0)
	tid := MonoFlake{
		epoch:       epoch,
		maxSequence: defaultMaxSequence,
		nodeBits:    defaultReservedNodeBits,
		mu:          &sync.Mutex{},
	}
	opts = append(opts, withNodeID(int64(nodeID)))
	for _, opt := range opts {
		if err := opt(&tid); err != nil {
			return nil, err
		}
	}
	return &tid, nil
}

// Next generates a new unique int64 ID
func (t *MonoFlake) Next() ID {
	t.mu.Lock()
	defer t.mu.Unlock()

	seq, ms := t.sequence, t.millisec
	since := time.Since(t.epoch).Milliseconds()

	if since < ms {
		since = ms
	} else if since > ms {
		seq = 0
	}
	nextMs := since

	nextSeq := seq + 1
	if nextSeq >= t.maxSequence {
		nextSeq = 0
		nextMs++
	}

	t.millisec = nextMs
	t.sequence = nextSeq

	return ID(since<<reservedAllocationBits | seq<<t.nodeBits | t.nodeID)
}
