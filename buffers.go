package next_utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	logging "github.com/ipfs/go-log/v2"
	"math/big"
)

var log = logging.Logger("utils")

func PackUint64LE(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}

func PackInt64BE(n int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(n))
	return b
}

func PackUint64BE(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func PackUint32LE(n uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)
	return b
}

func PackUint32BE(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

func PackInt32BE(n int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return b
}

func PackUint16LE(n uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, n)
	return b
}

func PackUint16BE(n uint16) []byte {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return b
}

func VarIntBytes(n uint64) []byte {
	if n < 0xFD {
		return []byte{byte(n)}
	}

	if n <= 0xFFFF {
		buff := make([]byte, 3)
		buff[0] = 0xFD
		binary.LittleEndian.PutUint16(buff[1:], uint16(n))
		return buff
	}

	if n <= 0xFFFFFFFF {
		buff := make([]byte, 5)
		buff[0] = 0xFE
		binary.LittleEndian.PutUint32(buff[1:], uint32(n))
		return buff
	}

	buff := make([]byte, 9)
	buff[0] = 0xFF
	binary.LittleEndian.PutUint64(buff[1:], uint64(n))
	return buff
}

func VarStringBytes(str string) []byte {
	bStr := []byte(str)
	return bytes.Join([][]byte{
		VarIntBytes(uint64(len(bStr))),
		bStr,
	}, nil)
}

func SerializeString(s string) []byte {
	if len(s) < 253 {
		return bytes.Join([][]byte{
			{byte(len(s))},
			[]byte(s),
		}, nil)
	} else if len(s) < 0x10000 {
		return bytes.Join([][]byte{
			{253},
			PackUint16LE(uint16(len(s))),
			[]byte(s),
		}, nil)
	} else if len(s) < 0x100000000 {
		return bytes.Join([][]byte{
			{254},
			PackUint32LE(uint32(len(s))),
			[]byte(s),
		}, nil)
	} else {
		return bytes.Join([][]byte{
			{255},
			PackUint64LE(uint64(len(s))),
			[]byte(s),
		}, nil)
	}
}

func SerializeNumber(n uint64) []byte {
	if n >= 1 && n <= 16 {
		return []byte{
			0x50 + byte(n),
		}
	}

	l := 1
	buff := make([]byte, 9)
	for n > 0x7f {
		buff[l] = byte(n & 0xff)
		l++
		n >>= 8
	}
	buff[0] = byte(l)
	buff[l] = byte(n)

	return buff[0 : l+1]
}

func Uint256BytesFromHash(h string) []byte {
	container := make([]byte, 32)
	fromHex, err := hex.DecodeString(h)
	if err != nil {
		log.Error(err)
	}

	copy(container, fromHex)

	return ReverseBytes(container)
}

func ReverseBytes(b []byte) []byte {
	_b := make([]byte, len(b))
	copy(_b, b)

	for i, j := 0, len(_b)-1; i < j; i, j = i+1, j-1 {
		_b[i], _b[j] = _b[j], _b[i]
	}
	return _b
}

func BigIntFromBitsHex(bits string) *big.Int {
	bBits, err := hex.DecodeString(bits)
	if err != nil {
		log.Panic(err)
	}
	return BigIntFromBitsBytes(bBits)
}

func BigIntFromBitsBytes(bits []byte) *big.Int {
	bytesNumber := bits[0]

	bigBits := new(big.Int).SetBytes(bits[1:])
	return new(big.Int).Mul(bigBits, new(big.Int).Exp(big.NewInt(2), big.NewInt(8*int64(bytesNumber-3)), nil))
}

// ReverseByteOrder LE <-> BE
func ReverseByteOrder(b []byte) []byte {
	_b := make([]byte, len(b))
	copy(_b, b)

	for i := 0; i < 8; i++ {
		binary.LittleEndian.PutUint32(_b[i*4:], binary.BigEndian.Uint32(_b[i*4:]))
	}
	return ReverseBytes(_b)
}

