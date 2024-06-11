// Copyright 2024 Mustafa Turan. All rights reserved.
// Use of this source code is governed by a Apache License 2.0 license that can
// be found in the LICENSE file.

package monoflake

type ID int64

const (
	maxBase62     = uint64(62)
	base62Mapping = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

// Int64 returns the int64 value of the ID
func (id ID) Int64() int64 {
	return int64(id)
}

// Bytes returns the byte array of the ID
func (id ID) Bytes() []byte {
	return toBase62WithPaddingZeros(uint64(id), 11)
}

// String returns the Base62 encoded string of the ID
func (id ID) String() string {
	return string(toBase62WithPaddingZeros(uint64(id), 11))
}

// Since returns the milliseconds since epoch of the ID
func (id ID) Since() int64 {
	return int64(id) >> reservedAllocationBits
}

// Sequence returns the sequence number of the ID
func (id ID) Sequence(nodeBits int64) int64 {
	shiftedID := id.Int64() >> nodeBits
	mask := (int64(1) << (reservedAllocationBits - nodeBits)) - 1
	return shiftedID & mask
}

// NodeID returns the node identifier of the ID
func (id ID) NodeID(nodeBits int64) int64 {
	mask := int64((1 << nodeBits) - 1)
	return id.Int64() & mask
}

// ToBase62WithPaddingZeros converts int types to Base62 encoded byte array
// with padding zeros
func toBase62WithPaddingZeros(u uint64, length int) []byte {
	const size = 11 // largest uint64 in base62 occupies 11 bytes
	var a [size]byte
	i := size
	for u >= maxBase62 {
		i--
		// Avoid using r = a%b in addition to q = a/maxBase62
		// since 64bit division and modulo operations
		// are calculated by runtime functions on 32bit machines.
		q := u / maxBase62
		a[i] = base62Mapping[u-q*maxBase62]
		u = q
	}
	// when u < maxBase62
	i--
	a[i] = base62Mapping[u]
	for i > size-length {
		i--
		a[i] = base62Mapping[0]
	}
	return a[i:]
}
