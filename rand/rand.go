package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const SessionTokenBytes = 32

// Bytes 函数以给定的长度产生随机字节数组。
func Bytes(n int) ([]byte, error) {
	// 创建一个长度为 n 的字节数组。
	b := make([]byte, n)

	// 使用标准库的rand包的Read方法向字节数组中填充随机数据。
	nRead, err := rand.Read(b)
	if err != nil {
		// 生成随机字节时出现错误，返回已生成的字节数组和错误原因。
		return nil, fmt.Errorf("bytes:%w", err)
	}
	if nRead < n {
		// 若读取的随机字节不足，返回错误。
		return nil, fmt.Errorf("bytes:didnt read enough random bytes")
	}
	return b, nil
}

// String 函数以给定的长度产生随机字符串。
func String(n int) (string, error) {

	// 调用 Bytes 函数获得指定长度的随机字节数组。
	b, err := Bytes(n)
	if err != nil {
		return "", fmt.Errorf("string:%w", err)
	}

	// 对生成的字节数组进行base64编码，然后返回生成的字符串。
	return base64.URLEncoding.EncodeToString(b), nil
}

func SessionToken() (string, error) {
	return String(SessionTokenBytes)
}
