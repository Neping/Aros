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

func IntMax(a, b int) int {
	if a < b {
        return b
    }

    return a
}