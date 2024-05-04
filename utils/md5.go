package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

// 小写
func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	//h.Sum(nil) 会完成哈希计算过程并返回一个包含16字节（128位）MD5哈希结果的字节数组。
	tempStr := h.Sum(nil)
	//将二进制的哈希值转化为16进制的字符串
	return hex.EncodeToString(tempStr)
}

// 大写
func MD5Encode(data string) string {
	return strings.ToUpper(Md5Encode(data))
}

// 加密
func MakePassword(plainpwd, salt string) string {
	//plainpwd表示密码，salt表示随机数
	return Md5Encode(plainpwd + salt)
}

// 解密
func ValidPassword(plainpwd, salt, password string) bool {
	return Md5Encode(plainpwd+salt) == password
}
