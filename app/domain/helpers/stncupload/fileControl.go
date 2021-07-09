package stncupload

import (
	"fmt"
	"log"
	"os"
)

func FileDelete(uploadPath string, filename string) {

	//https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	// if _, err := os.Stat(uploadPath + filename); os.IsExist(err) {
	// 	// uploadError = uploadPath + " dosya var sile "
	// 	fmt.Println("uploadPath + filename")
	// 	fmt.Println(uploadPath + filename)
	// 	e := os.Remove(uploadPath + filename)
	// 	if e != nil {
	// 		log.Default()
	// 	}
	// }

	//https://stackoverflow.com/questions/12518876/how-to-check-if-a-file-exists-in-go
	if _, err := os.Stat(uploadPath + filename); err == nil {
		fmt.Println("uploadPath + filename")
		fmt.Println(uploadPath + filename)
		e := os.Remove(uploadPath + filename)
		if e != nil {
			log.Default()
		}
	}
	// else if os.IsNotExist(err) {
	//     //dosya yok demek
	//   }

}
