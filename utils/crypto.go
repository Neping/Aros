package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"
	"os"
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

// src 明文
// key 16位密钥
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

// src 密文
// key 16位密钥
func DecryptAES(src, key []byte) []byte {
	block, err  := aes.NewCipher(key)
	if nil != err {
		panic(err)
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)

	return unpaddingText(src)
}

// RSA Public Key & Private Key Create
func RsaGenKey(bits int) error {
	/// 生成私钥
	priKey, err := rsa.GenerateKey(rand.Reader, bits)
	if nil != err {
		return err
	}
	// x509标准
	priStream := x509.MarshalPKCS1PrivateKey(priKey)
	block := pem.Block{
		Type: "RSA Private Key",
		Bytes: priStream,
	}
	// 通过pem将设置好的数据进行编码，并写入磁盘文件
	priFile, err := os.Create("private.pem")
	if nil != err {
		return err
	}
	defer priFile.Close()
	err  = pem.Encode(priFile, &block)
	if nil != err {
		return err
	}
	
	/// 生成公钥
	pubKey := priKey.PublicKey
	pubStream, err := x509.MarshalPKIXPublicKey(&pubKey)
	if nil != err {
		return err
	}
	block = pem.Block{
		Type: "RSA Public Key",
		Bytes: pubStream,
	}
	pubFile, err := os.Create("public.pem")
	if nil != err {
		return err
	}
	defer pubFile.Close()
	err = pem.Encode(pubFile, &block)
	if nil != err {
		return err
	}

	return nil
}

// RSA 公钥加密
// src 待加密数据
// pathName 公钥文件路径
func EncryptPubRSA(src []byte, pathName string) ([]byte, error) {
	ret := []byte("")
	fp, err := os.Open(pathName)
	if nil != err {
		return ret, err
	}
	// 得到文件属性[文件大小]
	info, err := fp.Stat()
	if nil != err {
		return ret, err
	}
	buf := make([]byte, info.Size())
	_, err = fp.Read(buf)
	defer fp.Close()
	if nil != err {
		return ret, err
	}
	// 字符串解码
	block, _ := pem.Decode(buf)
	// 按x509规范解析公钥
	pubParse, err := x509.ParsePKIXPublicKey(block.Bytes)
	if nil != err {
		return ret, err
	}
	pubKey := pubParse.(*rsa.PublicKey)
	// 加密数据
	ret, err = rsa.EncryptPKCS1v15(rand.Reader, pubKey, src)
	if nil != err {
		return ret, err
	}

	return ret, nil
}

// RSA 私钥解密
// src 待解密数据
// pathName 私钥文件路径
func DecryptPriRSA(src []byte, pathName string) ([]byte, error) {
	ret := []byte("")
	fp, err := os.Open(pathName)
	if nil != err {
		return ret, err
	}
	// 得到文件属性[文件大小]
	info, err := fp.Stat()
	if nil != err {
		return ret, err
	}
	buf := make([]byte, info.Size())
	_, err = fp.Read(buf)
	defer fp.Close()
	if nil != err {
		return ret, err
	}
	// 字符串解码
	block, _ := pem.Decode(buf)
	// 按x509规范还原私钥数据
	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if nil != err {
		return ret, err
	}
	// 通过私钥对数据解码
	ret, err = rsa.DecryptPKCS1v15(rand.Reader, priKey, src)
	if nil != err {
		return ret, err
	}

	return ret, nil
}

func TestAES()  {
	src := []byte("this is a aes test")
	key := []byte("1234567890abcdef")
	str := EncryptAES(src, key)
	fmt.Println("EncryptAES: "+ string(str))
	str = DecryptAES(str, key)
	fmt.Println("DecryptAES: "+ string(str))
}

func TestRsa()  {
	err := RsaGenKey(4096)
	fmt.Println("key create error info: ", err)
	// 公钥加密
	src := []byte("this is a rsa test!")
	fmt.Println("data: ", string(src))

	encr, err := EncryptPubRSA(src, "public.pem")
	if nil != err {
		panic(err)
	}
	fmt.Println("public encrypt data: ", string(encr))

	// 私钥解密
	decr, err :=  DecryptPriRSA(encr, "private.pem")
	if nil != err {
		panic(err)
	}
	fmt.Println("private encrypt data: ", string(decr))
}

// md5
// fmt.Println(utils.MD5([]byte("123456")))
func MD5(src []byte) string {
	res := md5.Sum(src)
	//ret := hex.EncodeToString(res[:])
	ret := fmt.Sprintf("%x", res)

	return ret
}

// 可多重追加md5
// fmt.Println(utils.MD5Muilt([]byte("123456")))
func MD5Muilt(src []byte) string {
	hash := md5.New()
	// 添加数据
	// 方法1
	io.WriteString(hash, string(src))
	// 方法2
	// res.Write(src)
	res := hash.Sum(nil)

	return hex.EncodeToString(res)
}