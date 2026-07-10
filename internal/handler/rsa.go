package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

type RsaHandler struct{}

func NewRsaHandler() *RsaHandler {
	return &RsaHandler{}
}

// RSAKeyPair RSA密钥对响应
 type RSAKeyPair struct {
	Success   bool   `json:"success"`
	PublicKey string `json:"publicKey"`
	PrivateKey string `json:"privateKey"`
	Message   string `json:"message"`
}

// GenerateKeyPair 生成RSA密钥对
 func (h *RsaHandler) GenerateKeyPair(keySize int, keyFormat, password string) RSAKeyPair {
	// 生成RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return RSAKeyPair{
			Success: false,
			Message: fmt.Sprintf("生成密钥失败: %v", err),
		}
	}

	// 编码私钥
	var privateKeyPEM []byte
	if keyFormat == "PKCS#8" {
		// PKCS#8 格式
		privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return RSAKeyPair{
				Success: false,
				Message: fmt.Sprintf("编码私钥失败: %v", err),
			}
		}
		privateKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		})
	} else {
		// PKCS#1 格式 (默认)
		privateKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
		})
	}

	// 编码公钥（格式与私钥保持一致）
	var publicKeyPEM []byte
	if keyFormat == "PKCS#8" {
		// PKCS#8/PKIX 格式
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return RSAKeyPair{
				Success: false,
				Message: fmt.Sprintf("编码公钥失败: %v", err),
			}
		}
		publicKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		})
	} else {
		// PKCS#1 格式 (默认)
		publicKeyBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
		publicKeyPEM = pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		})
	}

	return RSAKeyPair{
		Success:    true,
		PublicKey:  string(publicKeyPEM),
		PrivateKey: string(privateKeyPEM),
		Message:    "密钥生成成功",
	}
}


