package stnccollection

import (
	"strconv"
)

//FloatToString64 float 2 string
func FloatToString64(inputNum float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(inputNum, 'f', 10, 64)
}

//StringToFloat64  to convert a float number to a string
func StringToFloat64(str string) (returnData float64, err2 error) {
	returnData, err2 = strconv.ParseFloat(str, 64)
	return returnData, err2

}

/*
https://selmantunc.com.tr/post/635793192618442752/golang-numeric-conversions

Atoi (string to int)
i, err := strconv.Atoi(“-42”)

——————————-

Itoa (int to string).
s := strconv.Itoa(-42)

  ———————

int64 to string
str:= strconv.FormatInt(int64(165), 10)

——————————-

uint64 to string
lastID := strconv.FormatUint(uint64(5656556666), 10)

——————————–

string to  uint64
catID, _ := strconv.ParseUint(“string”, 10, 64)

interface return to string
session.Get(key).(string)
*/

//Uint64toString uint64 2 string
func Uint64toString(inputNum uint64) string {
	return strconv.FormatUint(uint64(inputNum), 10)
}

//StringtoUint64 string 2 uint64
func StringtoUint64(inputStr string) (uintInt uint64) {
	uintInt, _ = strconv.ParseUint(inputStr, 10, 64)
	return uintInt
}

//StringToint string 2 int
//TODO: error vermemek sorun olur mu ?
func StringToint(inputStr string) (IntType int) {
	IntType, _ = strconv.Atoi(inputStr)
	return IntType
}

//IntToString int 2 string
//TODO: error vermemek sorun olur mu ?
func IntToString(inputStr int) string {
	return strconv.Itoa(inputStr)
}
