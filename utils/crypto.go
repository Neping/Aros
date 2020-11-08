package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

/**
 * 末尾填充字节
 */
func paddingText(src []byte, blockSize int) []byte {
	padding  := blockSize - len(src)%blockSize
	text := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(src, text...)
}

/**
 * 删除末尾填充字节
 */
func unpaddingText(src []byte) []byte {
	len := len(src)
	number := int(src[len -1])

	return src[:len - number]
}

func EncryptAES(src, key []byte) []byte {
	block, err := aes.NewCipher(key)
	if nil != err {
		panic(err)
	}
	src = paddingText(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)

	return src
}

func DecryptAES(src, key []byte) []byte {
	block, err  := aes.NewCipher(key)
	if nil != err {
		panic(err)
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)

	return unpaddingText(src)
}


func TestAES()  {
	src := []byte("this is a aes test")
	key := []byte("1234567890abcdef")
	str := EncryptAES(src, key)
	fmt.Println("EncryptAES: "+ string(str))
	str = DecryptAES(str, key)
	fmt.Println("DecryptAES: "+ string(str))
}
