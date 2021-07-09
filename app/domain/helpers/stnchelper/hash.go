package stnchelper

import "encoding/base64"

//BoxString eğer veriye hiç bişey Savemezse boşluk eklemeye yarar
func BoxString(x string) *string {
	return &x
}

func encode(src []byte) string {
	return base64.URLEncoding.EncodeToString(src)
}

func decode(src string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(src)
}

func bytesToString(data []byte) string {
	return string(data[:])
}
