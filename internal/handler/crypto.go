package handler

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type CryptoHandler struct{}

func NewCryptoHandler() *CryptoHandler {
	return &CryptoHandler{}
}

// EncryptRequest 加密请求结构
type EncryptRequest struct {
	Mode string `json:"mode"` // aes128-cbc, aes256-cbc, aes128-ecb
	Key  string `json:"key"`
	Text string `json:"text"`
}

// DecryptRequest 解密请求结构
type DecryptRequest struct {
	Mode string `json:"mode"`
	Key  string `json:"key"`
	Text string `json:"text"`
}

// CryptoResponse 响应结构
type CryptoResponse struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
	Message string `json:"message"`
}

// parseMode 解析加密模式，返回 keySize 和是否是 ECB
func parseMode(mode string) (keySize int, isECB bool, err error) {
	switch mode {
	case "aes128-cbc", "aes128":
		return 16, false, nil
	case "aes256-cbc", "aes256":
		return 32, false, nil
	case "aes128-ecb":
		return 16, true, nil
	default:
		return 0, false, fmt.Errorf("不支持的加密模式: %s", mode)
	}
}

// Encrypt AES 加密
func (h *CryptoHandler) Encrypt(mode, key, text string) CryptoResponse {
	if text == "" {
		return CryptoResponse{Success: false, Message: "请输入要加密的内容"}
	}
	if key == "" {
		return CryptoResponse{Success: false, Message: "请输入加密密钥"}
	}

	keySize, isECB, err := parseMode(mode)
	if err != nil {
		return CryptoResponse{Success: false, Message: err.Error()}
	}

	keyBytes := adjustKey(key, keySize)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return CryptoResponse{Success: false, Message: fmt.Sprintf("加密失败: %v", err)}
	}

	// PKCS7 填充
	plaintext := pkcs7Padding([]byte(text), aes.BlockSize)

	if isECB {
		// ECB 模式加密：每个块独立加密，不需要 IV
		ciphertext := make([]byte, len(plaintext))
		ecbEncrypt(block, ciphertext, plaintext)
		result := base64.StdEncoding.EncodeToString(ciphertext)
		return CryptoResponse{Success: true, Result: result, Message: "加密成功"}
	}

	// CBC 模式加密
	iv := make([]byte, aes.BlockSize)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return CryptoResponse{Success: false, Message: fmt.Sprintf("加密失败: %v", err)}
	}

	ciphertext := make([]byte, len(plaintext))
	modeEnc := cipher.NewCBCEncrypter(block, iv)
	modeEnc.CryptBlocks(ciphertext, plaintext)

	// 将 IV 和密文拼接，然后 base64 编码
	result := base64.StdEncoding.EncodeToString(append(iv, ciphertext...))
	return CryptoResponse{Success: true, Result: result, Message: "加密成功"}
}

// Decrypt AES 解密
func (h *CryptoHandler) Decrypt(mode, key, text string) CryptoResponse {
	if text == "" {
		return CryptoResponse{Success: false, Message: "请输入要解密的内容"}
	}
	if key == "" {
		return CryptoResponse{Success: false, Message: "请输入解密密钥"}
	}

	keySize, isECB, err := parseMode(mode)
	if err != nil {
		return CryptoResponse{Success: false, Message: err.Error()}
	}

	keyBytes := adjustKey(key, keySize)
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return CryptoResponse{Success: false, Message: fmt.Sprintf("解密失败: %v", err)}
	}

	// Base64 解码
	data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return CryptoResponse{Success: false, Message: "密文格式错误，请输入有效的 Base64 编码"}
	}

	if isECB {
		// ECB 模式解密：直接解密密文（必须是块大小的整数倍）
		if len(data) == 0 || len(data)%aes.BlockSize != 0 {
			return CryptoResponse{Success: false, Message: "密文长度无效，不是块大小的整数倍"}
		}
		plaintext := make([]byte, len(data))
		ecbDecrypt(block, plaintext, data)
		plaintext, err = pkcs7Unpadding(plaintext)
		if err != nil {
			return CryptoResponse{Success: false, Message: "解密失败，请检查密钥是否正确"}
		}
		return CryptoResponse{Success: true, Result: string(plaintext), Message: "解密成功"}
	}

	// CBC 模式解密
	if len(data) < aes.BlockSize*2 {
		return CryptoResponse{Success: false, Message: "密文长度不足"}
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	modeDec := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	modeDec.CryptBlocks(plaintext, ciphertext)

	plaintext, err = pkcs7Unpadding(plaintext)
	if err != nil {
		return CryptoResponse{Success: false, Message: "解密失败，请检查密钥是否正确"}
	}

	return CryptoResponse{Success: true, Result: string(plaintext), Message: "解密成功"}
}

// ecbEncrypt ECB 模式加密（每个块独立加密）
func ecbEncrypt(block cipher.Block, dst, src []byte) {
	blockSize := block.BlockSize()
	for i := 0; i < len(src); i += blockSize {
		block.Encrypt(dst[i:i+blockSize], src[i:i+blockSize])
	}
}

// ecbDecrypt ECB 模式解密（每个块独立解密）
func ecbDecrypt(block cipher.Block, dst, src []byte) {
	blockSize := block.BlockSize()
	for i := 0; i < len(src); i += blockSize {
		block.Decrypt(dst[i:i+blockSize], src[i:i+blockSize])
	}
}

// adjustKey 调整密钥长度到指定大小
func adjustKey(key string, size int) []byte {
	keyBytes := []byte(key)
	keyLen := len(keyBytes)

	if keyLen >= size {
		return keyBytes[:size]
	}

	// 如果密钥太短，重复填充直到达到指定长度
	result := make([]byte, size)
	copy(result, keyBytes)

	for i := keyLen; i < size; i++ {
		result[i] = keyBytes[i%keyLen]
	}

	return result
}

// pkcs7Padding PKCS7 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

// pkcs7Unpadding PKCS7 去除填充
func pkcs7Unpadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("数据为空")
	}
	padding := int(data[length-1])
	if padding > length || padding > aes.BlockSize {
		return nil, fmt.Errorf("填充错误")
	}
	return data[:length-padding], nil
}
