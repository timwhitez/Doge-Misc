package main

import (
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
)

func main(){
	fmt.Println("sha1Encrypt")
	fmt.Println("sha1test")
	fmt.Println(sha1Encrypt("sha1test"))

	fmt.Println("str2sha1")
	fmt.Println("sha1test")
	fmt.Println(str2sha1("sha1test"))

	fmt.Println("sha256Encrypt")
	fmt.Println("sha256test")
	fmt.Println(sha256Encrypt("sha256test"))


	fmt.Println("Sha256Hex")
	fmt.Println("sha256test")
	fmt.Println(str2sha256("sha256test"))

}


// Encrypt encrypts any type of variable using SHA1 algorithms.
// It uses package gconv to convert `v` to its bytes type.
func sha1Encrypt(v interface{}) string {
	r := sha1.Sum(gconv.Bytes(v))
	return hex.EncodeToString(r[:])
}

// Encrypt encrypts any type of variable using SHA1 algorithms.
// It uses package gconv to convert `v` to its bytes type.
func sha256Encrypt(v interface{}) string {
	digest:=sha256.New()
	digest.Write(gconv.Bytes(v))
	r := digest.Sum(nil)
	return hex.EncodeToString(r[:])
}



func str2sha1(s string) string{
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}


func str2sha256(s string)string{
	digest:=sha256.New()
	digest.Write([]byte(s))
	return hex.EncodeToString(digest.Sum(nil))
}