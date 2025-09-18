package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
)

// CryptoUtils 加密工具
type CryptoUtils struct{}

// NewCryptoUtils 创建加密工具实例
func NewCryptoUtils() *CryptoUtils {
	return &CryptoUtils{}
}

// GenerateRandomBytes 生成随机字节
func (c *CryptoUtils) GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, fmt.Errorf("生成随机字节失败: %w", err)
	}
	return bytes, nil
}

// GenerateRandomString 生成随机字符串
func (c *CryptoUtils) GenerateRandomString(length int) (string, error) {
	bytes, err := c.GenerateRandomBytes(length)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}

// HashSHA256 计算SHA256哈希
func (c *CryptoUtils) HashSHA256(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// HashSHA512 计算SHA512哈希
func (c *CryptoUtils) HashSHA512(data []byte) string {
	hash := sha512.Sum512(data)
	return hex.EncodeToString(hash[:])
}

// HashWithSalt 使用盐值计算哈希
func (c *CryptoUtils) HashWithSalt(data []byte, salt []byte) string {
	h := sha256.New()
	h.Write(data)
	h.Write(salt)
	return hex.EncodeToString(h.Sum(nil))
}

// GenerateHMAC 生成HMAC
func (c *CryptoUtils) GenerateHMAC(data []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// VerifyHMAC 验证HMAC
func (c *CryptoUtils) VerifyHMAC(data []byte, key []byte, expectedHMAC string) bool {
	actualHMAC := c.GenerateHMAC(data, key)
	return actualHMAC == expectedHMAC
}

// EncryptAES AES加密
func (c *CryptoUtils) EncryptAES(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES密码失败: %w", err)
	}

	// 生成随机IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("生成IV失败: %w", err)
	}

	// 填充数据
	paddedData := c.pkcs7Pad(data, aes.BlockSize)

	// 加密
	ciphertext := make([]byte, len(paddedData))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedData)

	// 返回IV + 密文
	result := make([]byte, 0, len(iv)+len(ciphertext))
	result = append(result, iv...)
	result = append(result, ciphertext...)

	return result, nil
}

// DecryptAES AES解密
func (c *CryptoUtils) DecryptAES(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES密码失败: %w", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("密文长度不足")
	}

	// 分离IV和密文
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// 解密
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)

	// 去除填充
	unpaddedData, err := c.pkcs7Unpad(plaintext)
	if err != nil {
		return nil, fmt.Errorf("去除填充失败: %w", err)
	}

	return unpaddedData, nil
}

// pkcs7Pad PKCS7填充
func (c *CryptoUtils) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// pkcs7Unpad PKCS7去除填充
func (c *CryptoUtils) pkcs7Unpad(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, fmt.Errorf("数据为空")
	}

	padding := int(data[length-1])
	if padding > length {
		return nil, fmt.Errorf("填充无效")
	}

	// 验证填充
	for i := length - padding; i < length; i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("填充无效")
		}
	}

	return data[:length-padding], nil
}

// GenerateKey 生成密钥
func (c *CryptoUtils) GenerateKey(keySize int) ([]byte, error) {
	return c.GenerateRandomBytes(keySize)
}

// HashFile 计算文件哈希
func (c *CryptoUtils) HashFile(data []byte, algorithm string) (string, error) {
	var h hash.Hash

	switch algorithm {
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return "", fmt.Errorf("不支持的哈希算法: %s", algorithm)
	}

	h.Write(data)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Base64Encode Base64编码
func (c *CryptoUtils) Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// Base64Decode Base64解码
func (c *CryptoUtils) Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

// HexEncode 十六进制编码
func (c *CryptoUtils) HexEncode(data []byte) string {
	return hex.EncodeToString(data)
}

// HexDecode 十六进制解码
func (c *CryptoUtils) HexDecode(encoded string) ([]byte, error) {
	return hex.DecodeString(encoded)
}

// GeneratePasswordHash 生成密码哈希
func (c *CryptoUtils) GeneratePasswordHash(password string, salt []byte) string {
	return c.HashWithSalt([]byte(password), salt)
}

// VerifyPassword 验证密码
func (c *CryptoUtils) VerifyPassword(password string, salt []byte, expectedHash string) bool {
	actualHash := c.GeneratePasswordHash(password, salt)
	return actualHash == expectedHash
}

// GenerateAPIKey 生成API密钥
func (c *CryptoUtils) GenerateAPIKey(length int) (string, error) {
	return c.GenerateRandomString(length)
}

// EncryptSensitiveData 加密敏感数据
func (c *CryptoUtils) EncryptSensitiveData(data string, key []byte) (string, error) {
	encrypted, err := c.EncryptAES([]byte(data), key)
	if err != nil {
		return "", err
	}
	return c.Base64Encode(encrypted), nil
}

// DecryptSensitiveData 解密敏感数据
func (c *CryptoUtils) DecryptSensitiveData(encryptedData string, key []byte) (string, error) {
	encrypted, err := c.Base64Decode(encryptedData)
	if err != nil {
		return "", err
	}

	decrypted, err := c.DecryptAES(encrypted, key)
	if err != nil {
		return "", err
	}

	return string(decrypted), nil
}
