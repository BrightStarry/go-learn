package util

import (
	"crypto/aes"
	"crypto/cipher"
)

/**
aes 算法相关工具类
 */

//解密
func AESDecrypt(crypted,key []byte)[]byte{
	block,_ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block,key[:blockSize])
	origData := make([]byte,len(crypted))
	blockMode.CryptBlocks(origData,crypted)
	origData = PKCS7UnPadding(origData)
	return origData
}

//去补码
func PKCS7UnPadding(origData []byte)[]byte{
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}
