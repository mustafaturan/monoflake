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

// BigEndianBytes returns 8 bytes array of the ID as big endian format
func (id ID) BigEndianBytes() []byte {
	v := fromInt64ToBytes(id.Int64())
	return v[:]
}

// Bytes returns 11 bytes array of the ID based on base62 encoding
func (id ID) Bytes() []byte {
	v := toBase62WithPaddingZeros(uint64(id))
	return v[:]
}

// String returns the Base62 encoded string of the ID
func (id ID) String() string {
	v := toBase62WithPaddingZeros(uint64(id))
	return string(v[:])
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

// IDFromBase62 is a helper fn that converts Base62 encoded string to ID
// returns -1 on overflow
func IDFromBase62(base62 string) ID {
	var id int64
	for _, r := range base62 {
		id = id*62 + fromBase62RuneToInt64(r)
		if id < 0 {
			return -1
		}
	}
	return ID(id)
}

// IDFromBigEndianBytes is a helper fn that converts big endian bytes to ID
// returns -1 on overflow
func IDFromBigEndianBytes(b []byte) ID {
	var id int64
	for _, r := range b {
		id = id*256 + int64(r)
		if id < 0 {
			return -1
		}
	}
	return ID(id)
}

// ToBase62WithPaddingZeros converts int types to Base62 encoded byte array
// with padding zeros
func toBase62WithPaddingZeros(u uint64) [11]byte {
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
	for i > 0 {
		i--
		a[i] = base62Mapping[0]
	}
	return a
}

func fromBase62RuneToInt64(char rune) int64 {
	if char >= '0' && char <= '9' {
		return int64(char - '0')
	}
	if char >= 'A' && char <= 'Z' {
		return int64(char - 'A' + 10)
	}
	if char >= 'a' && char <= 'z' {
		return int64(char - 'a' + 36)
	}
	return 0
}

func fromInt64ToBytes(n int64) [8]byte {
	const size = 8
	var b [size]byte
	for i := 0; i < size; i++ {
		b[size-1-i] = byte(n >> (i * 8))
	}
	return b
}
