package stnccollection

import (
	"fmt"
	"math"
	"strings"
)

//https://golangcode.com/check-if-row-exists-in-slice/
//FindSlice elemnt
func FindSlice(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func FindSliceTypes(slice []string, val string) bool {

	for _, n := range slice {
		if val == n {
			return true
		}
	}
	return false

}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

//link https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision
//ToFixedDecimal decimal format
func ToFixedDecimal(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

//SayiYuvarla sayı ayarlar
func SayiYuvarla(value64 float64) float64 {
	var orginalValue, resultData float64
	orginalValue = value64
	// fmt.Println(orginalValue)
	var value string
	value = FloatToString64(value64)

	stringSlice := strings.Split(value, ".")

	value = stringSlice[0]
	total := len(value)

	total1Eksi := len(value) - 1

	sonSayi := value[total1Eksi:total]

	var sonSayiKontrol float64
	sonSayiKontrol, _ = StringToFloat64(sonSayi)
	if sonSayiKontrol > 5 && sonSayiKontrol <= 9 {
		resultfloat64, _ := StringToFloat64(sonSayi)
		// fmt.Println(resultfloat64)
		resultData = (10 - resultfloat64)
		// fmt.Println(resultData)
		resultData = resultData + orginalValue
		// fmt.Println("resultData ust ")
		fmt.Println(resultData)

	} else {
		resultData = orginalValue - sonSayiKontrol
		// fmt.Println("resultData alt ")
		fmt.Println(resultData)

	}
	return resultData
}
