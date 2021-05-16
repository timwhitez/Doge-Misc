package main
 
import (
	"encoding/base64"
	"fmt"
)
 
/*通用编码表，编码解码*/
func raw_64Encode(s string)string {
	//编码
	s64_std := base64.StdEncoding.EncodeToString([]byte(s))
	return s64_std
}

/*通用编码表，编码解码*/
func raw_64Decode(s64_std string) string {
	//解码
	decodeBytes, err := base64.StdEncoding.DecodeString(s64_std)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(decodeBytes)
}


/*
自定义编码表，编码，解码（自定义码表 可随意变换字母排列顺序，然后会自动生成解密表）
Base64只能算是一个编码算法
*/
func self_64Encode(s string)string {
	encodeStd := "rstuKLYZlmn9+/TUopqOGHIhijkABCwxyz0VWXMNgD345QRSPv6abcEF812Jdef7"
	s64 := base64.NewEncoding(encodeStd).EncodeToString([]byte(s))
	return s64
}



func self_64Decode(s64 string) string {
	encodeStd := "rstuKLYZlmn9+/TUopqOGHIhijkABCwxyz0VWXMNgD345QRSPv6abcEF812Jdef7"
	//解码
	decodeBytes, err := base64.NewEncoding(encodeStd).DecodeString(s64)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(decodeBytes)
}

 
func main() {
	fmt.Println(raw_64Encode("Hello World!!!"))
	fmt.Println(raw_64Decode(raw_64Encode("Hello World!!!")))

	fmt.Println(self_64Encode("Hello World!!!"))
	fmt.Println(self_64Decode(self_64Encode("Hello World!!!")))
}