package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
)

func Servinit() (*rsa.PrivateKey, []byte) {

	log.Println("服务端生成RSA密钥对并存储到变量中")
	// 1. 生成RSA密钥对并存储到变量中
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Error generating RSA key pair:", err)
		return nil, nil
	}
	publicKey := privateKey.PublicKey
	pubASN1, err := x509.MarshalPKIXPublicKey(&publicKey)

	if err != nil {
		fmt.Println("Error marshalling public key:", err)
		return nil, nil
	}
	publicKeyBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	return privateKey, publicKeyBytes

}

func ServStp1(innerKey string, publicKeyBytes []byte) []byte {
	log.Println("服务端通过预制key AES的方式加密RSA的私钥")
	// 2. 使用设置好的一个内置key，通过AES的方式加密rsa的私钥
	ciphertext, err := encrypt_AES([]byte(innerKey), publicKeyBytes)
	if err != nil {
		fmt.Println("Error encrypting public key:", err)
		return nil
	}
	return ciphertext
}

func ClientStp1(innerKey string, ciphertext []byte) *rsa.PublicKey {
	log.Println("客户端通过预制key AES的方式解密RSA公钥")
	// 3. 使用内置key解密刚才加密的公钥
	decryptedPubBytes, err := decrypt_AES([]byte(innerKey), ciphertext)
	if err != nil {
		fmt.Println("Error decrypting encrypted public key:", err)
		return nil
	}

	// 把pem格式解密出来
	block, _ := pem.Decode(decryptedPubBytes)
	decryptedPubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("Error parsing decrypted public key:", err)
		return nil
	}
	rsaPubKey := decryptedPubKey.(*rsa.PublicKey)
	return rsaPubKey
}

func ClientStp2(rsaPubKey *rsa.PublicKey) []byte {
	log.Println("客户端生成一个随机的key1 并且用RSA公钥加密")
	// 4. 客户端生成一个随机的key1 并且用公钥加密
	key1 := make([]byte, 32)
	_, err := io.ReadFull(rand.Reader, key1)
	if err != nil {
		fmt.Println("Error generating random key1:", err)
		return nil
	}
	fmt.Printf("原始Key1：%x\n", key1)
	// 5. 公钥，rsa加密key1
	key1Encrypted, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, key1, []byte{})
	if err != nil {
		fmt.Println("Error encrypting key1 with private key:", err)
		return nil
	}
	return key1Encrypted
}

func ServStp2(privateKey *rsa.PrivateKey, key1Encrypted []byte) []byte {
	log.Println("服务端使用私钥解密key1")
	// 6. 使用私钥解密key1并且打印出来
	key1Decrypted, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, key1Encrypted, []byte{})
	if err != nil {
		fmt.Println("Error decrypting key1 with public key:", err)
		return nil
	}

	return key1Decrypted

}

func main() {
	var PrivateKey *rsa.PrivateKey
	var PublicKeyBytes []byte
	var InnerKey = "abcdefghijklmnop" // 预制的对称加密key
	log.Printf("客户端与服务端预制的key为: %s\n", InnerKey)

	PrivateKey, PublicKeyBytes = Servinit()

	pubEncBytes := ServStp1(InnerKey, PublicKeyBytes)

	tmpPubKey := ClientStp1(InnerKey, pubEncBytes)
	key1Enc := ClientStp2(tmpPubKey)

	realKey1 := ServStp2(PrivateKey, key1Enc)

	fmt.Printf("解密后的Key1：%x\n", realKey1)
}

func encrypt_AES(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	return gcm.Seal(nonce, nonce, data, nil), nil
}

func decrypt_AES(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
