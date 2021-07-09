package stnchelper

import (
	"fmt"
	"os"
)

func FileControlV1(uploadPath string) (uploadError string) {
	if _, err := os.Stat(uploadPath); err == nil {
		fmt.Println(err)
		uploadError = "kalsor var "
		// c.JSON(http.StatusOK, uploadError)

	} else if os.IsNotExist(err) {
		fmt.Println(err)
		uploadError = "kalsor yok "
		// c.JSON(http.StatusBadGateway, uploadError)
	}
	return uploadError
}

// uses
// func main() {
//     if fileExists("example.txt") {
//         fmt.Println("Example file exists")
//     } else {
//         fmt.Println("Example file does not exist (or is a directory)")
//     }
// }
// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
