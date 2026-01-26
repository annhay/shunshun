package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"shunshun/internal/pkg/global"
)

// 加密、解密工具

// Md5 加密
//
// 参数:
//   - str: 要加密的字符串
//
// 返回值:
//   - string: MD5加密后的字符串
func Md5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}

// 高级加密标准（Advanced Encryption Standard ,AES）

// GetPwdKey 获取AES加密密钥
//
// 返回值:
//   - []byte: AES加密密钥
func GetPwdKey() []byte {
	if global.AppConf != nil && global.AppConf.AES.SecretKey != "" {
		return []byte(global.AppConf.AES.SecretKey)
	}
	// 默认密钥（仅用于开发环境，生产环境必须配置）
	return []byte("default-secure-key-32bytes")
}

// PwdKey 16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
// key不能泄露
var PwdKey = GetPwdKey()

// PKCS7Padding PKCS7 填充模式
//
// 参数:
//   - ciphertext: 要填充的数据
//   - blockSize: 块大小
//
// 返回值:
//   - []byte: 填充后的数据
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	//Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 填充的反向操作，删除填充字符串
//
// 参数:
//   - origData: 要去填充的数据
//
// 返回值:
//   - []byte: 去填充后的数据
//   - error: 错误信息
func PKCS7UnPadding(origData []byte) ([]byte, error) {
	//获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	} else {
		//获取填充字符串长度
		unpadding := int(origData[length-1])
		//截取切片，删除填充字节，并且返回明文
		return origData[:(length - unpadding)], nil
	}
}

// AesEcrypt 实现AES加密
//
// 参数:
//   - origData: 要加密的原始数据
//   - key: 加密密钥
//
// 返回值:
//   - []byte: 加密后的数据
//   - error: 错误信息
func AesEcrypt(origData []byte, key []byte) []byte {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//对数据进行填充，让数据长度满足需求
	origData = PKCS7Padding(origData, blockSize)
	//采用 AES 加密方法中 CBC 加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	//执行加密
	blocMode.CryptBlocks(crypted, origData)
	return crypted
}

// AesDeCrypt 实现AES解密
//
// 参数:
//   - cypted: 要解密的数据
//   - key: 解密密钥
//
// 返回值:
//   - []byte: 解密后的数据
//   - error: 错误信息
func AesDeCrypt(cypted []byte, key []byte) []byte {
	//创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil
	}
	//获取块大小
	blockSize := block.BlockSize()
	//创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))
	//执行解密
	blockMode.CryptBlocks(origData, cypted)
	//去除填充字符串
	origData, err = PKCS7UnPadding(origData)
	if err != nil {
		return nil
	}
	return origData
}

// EnPwdCode 加密并进行base64编码
//
// 参数:
//   - pwd: 要加密的数据
//
// 返回值:
//   - string: 加密并base64编码后的字符串
//   - error: 错误信息
func EnPwdCode(pwd []byte) string {
	result := AesEcrypt(pwd, PwdKey)
	return base64.StdEncoding.EncodeToString(result)
}

// DePwdCode 解密base64编码的字符串
//
// 参数:
//   - pwd: 要解密的base64编码字符串
//
// 返回值:
//   - []byte: 解密后的数据
//   - error: 错误信息
func DePwdCode(pwd string) []byte {
	//解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return nil
	}
	//执行 AES 解密
	return AesDeCrypt(pwdByte, PwdKey)

}
