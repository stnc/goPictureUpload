package stnchelper

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

//TestModeCode hepsini testi için
func TestModeCode() {
	// str := "SAADET SELMAN"
	str := "ÇEKİÇSADETTİN BÜYÜKÇEKic"

	// Convert string to rune slice before taking substrings.
	// ... This will handle Unicode characters correctly.
	//     Not needed for ASCII strings.
	runes := []rune(str)

	stringSlice := strings.Split(str, " ")
	fmt.Printf("%v\n", stringSlice)
	ilkKisim := stringSlice[0]
	ikiKisim := stringSlice[1]

	fmt.Println("ilk kısım ", ilkKisim)

	fmt.Println("ikinci kısım ", ikiKisim)

	toplamKarakter := utf8.RuneCountInString(str)
	fmt.Println("toplam karakter ", toplamKarakter)

	bassiz := string(runes[2:])
	fmt.Println(" basi_yok:", bassiz)
	fmt.Println(" bastan 2:", string(runes[:2]))
	fmt.Println(" sondan 2:", string(runes[toplamKarakter-2:]))
	return
	yapianHali := string(bassiz[0 : toplamKarakter-2])
	fmt.Println(" yapilan:", yapianHali)
	fmt.Println(strings.Replace(yapianHali, " ", "#", -1))
	toplamKarakter2 := len(yapianHali)
	fmt.Println(" yapilan:", toplamKarakter2)

	var text string

	for i := 1; i < toplamKarakter2; i++ {
		text += "*"
	}
	fmt.Println(" yapilan:", text)

}

func LetterPrefix(str string) string {
	if len(str) > 0 {
		runes := []rune(str)
		toplamKarakter := utf8.RuneCountInString(str)
		newStr := truncString(string(runes[:2]), toplamKarakter)
		fmt.Println(" prefix 2:", newStr)
		return newStr
	} else {
		return ""
	}

}

func LetterSuffix(str string) string {
	toplamKarakter := utf8.RuneCountInString(str)
	// str1 := truncString(str, toplamKarakter)
	fmt.Println("toplam str ", toplamKarakter)
	if toplamKarakter > 0 {
		runes := []rune(str)
		fmt.Println("oku ", string(runes[toplamKarakter-2:]))
		newStr := string(runes[toplamKarakter-2:])
		// fmt.Println(" suffix 2:", newStr)
		// fmt.Println("son karakter", newStr)
		return string(newStr)

	} else {
		return ""
	}
}

func LetterMiddle(str string) string {
	if utf8.RuneCountInString(str) > 0 {
		runes := []rune(str)
		chr := string(runes[2:])
		return chr
	} else {
		return ""
	}
}

//TODO: tr karakterleir öçevir
func LetterMiddleStars(str string) string {
	if utf8.RuneCountInString(str) > 0 {
		stringSlice := strings.Split(str, " ")
		fmt.Printf("%v\n", stringSlice)
		ilkKisim := stringSlice[0]
		ikiKisim := stringSlice[1]
		fmt.Println("ilk kısım ", ilkKisim)
		fmt.Println("ikinci kısım ", ikiKisim)
		var ilkKisimtxt string
		var sonKisimtxt string
		ilkKisimToplamKarakter := len(ilkKisim)
		sonKisimToplamKarakter := len(ikiKisim)
		for i := 1; i < ilkKisimToplamKarakter; i++ {
			ilkKisimtxt += "*"
		}

		for i2 := 1; i2 < sonKisimToplamKarakter; i2++ {
			sonKisimtxt += "*"
		}

		return ilkKisimtxt + " " + sonKisimtxt
	} else {
		return ""
	}
}

func truncString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	start := (maxLen + 1) / 2
	for start >= 0 && s[start]>>6 == 0b10 {
		start--
	}
	end := utf8.RuneCountInString(s) - (maxLen - start)
	for end < utf8.RuneCountInString(s) && s[end]>>6 == 0b10 {
		end++
	}
	return s[:start] + s[end:]
}
