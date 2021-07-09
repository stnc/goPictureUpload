//  golang gin framework mvc and clean code project
//  Licensed under the Apache License 2.0
//  @author Selman TUNÇ <selmantunc@gmail.com>
//  @link: https://github.com/stnc/go-mvc-blog-clean-code
//  @license: Apache License 2.0
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/inancgumus/screen"
	"github.com/joho/godotenv"
	"github.com/scylladb/termtables"
)

//https://github.com/stnc-go/gobyexample/blob/master/pongo2render/render.go
func init() {
	//To load our environmental variables.

	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}

	/* //bu sunucuda çalışıyor
		    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	        if err != nil {
	            log.Fatal(err)
	        }
	        environmentPath := filepath.Join(dir, ".env")
	        err = godotenv.Load(environmentPath)
	        fatal(err)
	*/

}

func main() {

	err := filepath.Walk("res1/1",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// table := termtables.CreateTable()
			// table.AddHeaders("PATH", "FİLE")
			// table.AddRow(path, info.Name())
			// fmt.Println(table.Render())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	table2 := termtables.CreateTable()
	table2.AddHeaders("RESİM YÜKLEYİCİ")
	table2.AddRow("1. resim yukleme için 1 basınız")
	table2.AddRow("2. tüm resimleri silmek için 2 ye basınızz")
	// table2.AddRow("3. sistemdeki hayvan bilgisi için")
	fmt.Println(table2.Render())

	fmt.Print("-> ")

	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Println(err)
	}

	// print out the unicode value i.e. A -> 65, a -> 97
	fmt.Println(char)

	switch char {
	case '1':
		fmt.Println("A Key Pressed")
		screen.Clear()

		for {
			// Moves the cursor to the top left corner of the screen
			screen.MoveTopLeft()

			// fmt.Println(time.Now())
			// time.Sleep(time.Second)
		}

		break
	case '2':
		fmt.Println("2 ye bastın a Key Pressed")

		break
	}
}

/*
Data Entered
================================
Database: root_CRm2
Username: root_CRm2
Password: dtw7JV3M
*/
