package security

import (
	"crypto/sha1"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//TODO: buradaki şifreleme tutarsız davranıyor
/*
//
	hashPassword, _ := security.Hash(pass)
	fmt.Println("ŞİFELENMİŞ HALİ" +hashPassword )
	$2a$10$QPiWAgMpwHBkDjBL5pPd2.HBlfdniuGOvZd5kh.ILLjKFo67rvfsO
*/
//Hash is code hash yaparak verir
func HashOLD(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPasswordOLD(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

//VerifyPassword doğrulama
func VerifyPasswordApINEOLCAK(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Hash(Txt string) string {
	h := sha1.New()
	h.Write([]byte(Txt + "uhOx." + Txt))
	bs := h.Sum(nil)
	sh := string(fmt.Sprintf("%x", bs))
	return sh
}

func PassVer(InputPassword string) string {

	var userLoginHashPasword string

	userLoginHashPasword = Hash(InputPassword)

	return string(userLoginHashPasword)

}

func VerifyPassword(hashedPassword, InputPassword string) bool {

	var userLoginHashPasword string

	userLoginHashPasword = Hash(InputPassword)

	userLoginHashPasword = string(userLoginHashPasword)

	if hashedPassword == userLoginHashPasword {
		return true
	} else {
		return false
	}
}

func VerifyPasswordApi(hashedPassword, InputPassword string) bool {

	var userLoginHashPasword string

	userLoginHashPasword = Hash(InputPassword)

	userLoginHashPasword = string(userLoginHashPasword)

	if hashedPassword == userLoginHashPasword {
		return true
	} else {
		return false
	}
}
