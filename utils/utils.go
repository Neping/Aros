package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

func IntToHex(num int64) []byte  {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if nil != err {
		log.Panic(err)
	}

	return buff.Bytes()
}

func Uint64ToByte(number uint64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, number)
	if nil != err {
		log.Panic(err)
	}

	return buf.Bytes()
}

func IntMax(a, b int) int {
	if a < b {
        return b
    }

    return a
}