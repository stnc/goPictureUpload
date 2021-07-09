package stncdatetime

import (
	"fmt"
	"strings"
	"time"
)

//oku
//https://www.golangprograms.com/get-current-date-and-time-in-various-format-in-golang.html

type (
	// Inow : Custom Renderer for templates
	Inow struct{ Debug bool }
)

// var tarih stncdatetime.Inow
// fmt.Println("AY")
// fmt.Println(tarih.AylarListe("May"))
// fmt.Println("sql Date")
// fmt.Println(tarih.OnlyDate(posts.Odemeler[0].CreatedAt.String()))

//https://ednsquare.com/story/date-and-time-manipulation-golang-with-examples------cU1FjK
//https://yayprogramming.com/mysql-records-between-two-dates-in-go-language/  now ile yapmış

//      TODO buradan devam edelim
//https://yourbasic.org/golang/format-parse-string-time-date-example/#standard-time-and-date-formats tam formatlar
//Tarih test
func (testerModules Inow) Tarih() {

	//TODO: burası hata moduna bağlı olacak ama nasıl ????

	// if r.Debug {
	// 	fmt.Println("hata açık")
	// } else {
	// 	fmt.Println("hata kapalı")
	// }
	// The date we're trying to parse, work with and format
	myDateString := "2020-05-21 05:08"
	fmt.Println("My Starting Date:\t", myDateString)

	// Parse the date string into Go's time object
	// The 1st param specifies the format, 2nd is our date string
	myDate, _ := time.Parse("2006-01-02 15:04", myDateString)
	// Format uses the same formatting style as parse, or we can use a pre-made constant
	fmt.Println("My Date Reformatted:\t", myDate.Format(time.RFC822))
	fmt.Println("My RFC 1123:\t", myDate.Format(time.RFC1123))

	// In Y-m-d
	fmt.Println("Just The Date:\t\t", myDate.Format("02-01-2006 15:04"))
}

//TarihFullSQL for mysql
func (testerModules Inow) TarihFullSQL(myDateString string) string {

	//https://www.dotnetperls.com/substring-go
	/*
		myDateString := "2020-05-21 05:08"
		fmt.Println("My Starting Date:\t", myDateString)
	*/
	fmt.Println("My Starting Date:\t", myDateString)
	fmt.Println(len(myDateString))
	uzunluk := len(myDateString) - 3
	newDateString := myDateString[:uzunluk]
	fmt.Println(newDateString)

	s := strings.Split("Thu, 21 May 2020 05:08:00 UTC", ",")
	fmt.Println(s[0])

	s1 := strings.Split("Thu, 21 May 2020 05:08:00 UTC", " ")
	fmt.Println(s1[2])

	myDate, _ := time.Parse("2006-01-02 15:04", newDateString)
	fmt.Println("My RFC 1123:\t", myDate.Format(time.RFC1123))
	return myDate.String()
}

//https://www.php2golang.com/method/function.explode.html
func explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

/*
//FormatTarihForMysql tr
func (testerModules Inow) FormatTarihForMysql(tarih string) string {
	var gun, ay, yil string
	runes := []rune(tarih)
	fmt.Println("First 2", string(runes[0:2]))
	fmt.Println("Last 2", string(runes[1:]))
	gun = substr(tarih, 8, 2)
	ay = substr(tarih, 5, 2)
	yil = substr(tarih, 0, 4)
	return gun + "." + ay + ".".yil
}
*/

/*
2020-05-17 05:08:40
$tarih=date("Y-m-d H:i:s");
echo $bugun=format_tarih($tarih);
çalışan betiğin çıktısı 27.12.2010
*/
//AylarListe ayı türkçe verir
func (testerModules Inow) AylarListe(ay string) string {
	/*
		aylarKisa := make(map[string]string)
		aylarKisa["Jan"] = "Ocak"
		aylarKisa["Feb"] = "Şubat"
		aylarKisa["Mar"] = "Mart"
		aylarKisa["Apr"] = "Nisan"
		aylarKisa["May"] = "Mayıs"
		aylarKisa["Jun"] = "Haziran"
		aylarKisa["Jul"] = "Haziran"
		aylarKisa["Aug"] = "Ağustos"
		aylarKisa["Sep"] = "Eylül"
		aylarKisa["Oct"] = "Ekim"
		aylarKisa["Nov"] = "Kasım"
		aylarKisa["Dec"] = "Aralık"
		aylarUzun := make(map[string]string)
		aylarUzun["January"] = "Ocak"
		aylarUzun["February"] = "Şubat"
		aylarUzun["March"] = "Mart"
		aylarUzun["April"] = "Nisan"
		aylarUzun["May"] = "Mayıs"
		aylarUzun["June"] = "Haziran"
		aylarUzun["July"] = "Haziran"
		aylarUzun["August"] = "Ağustos"
		aylarUzun["September"] = "Eylül"
		aylarUzun["October"] = "Ekim"
		aylarUzun["November"] = "Kasım"
		aylarUzun["December"] = "Aralık"
	*/
	aylarFull := make(map[string]string)
	aylarFull["Jan"] = "Ocak"
	aylarFull["Feb"] = "Şubat"
	aylarFull["Mar"] = "Mart"
	aylarFull["Apr"] = "Nisan"
	aylarFull["May"] = "Mayıs"
	aylarFull["Jun"] = "Haziran"
	aylarFull["Jul"] = "Haziran"
	aylarFull["Aug"] = "Ağustos"
	aylarFull["Sep"] = "Eylül"
	aylarFull["Oct"] = "Ekim"
	aylarFull["Nov"] = "Kasım"
	aylarFull["Dec"] = "Aralık"
	aylarFull["January"] = "Ocak"
	aylarFull["February"] = "Şubat"
	aylarFull["March"] = "Mart"
	aylarFull["April"] = "Nisan"
	aylarFull["May"] = "Mayıs"
	aylarFull["June"] = "Haziran"
	aylarFull["July"] = "Haziran"
	aylarFull["August"] = "Ağustos"
	aylarFull["September"] = "Eylül"
	aylarFull["October"] = "Ekim"
	aylarFull["November"] = "Kasım"
	aylarFull["December"] = "Aralık"
	// ay := "Sat"
	return aylarFull[ay]
}

func (testerModules Inow) Gunler(gun string) string {
	gunler := make(map[string]string)

	gunler["Mon"] = "Pazartesi"
	gunler["Tues"] = "Salı"
	gunler["Wed"] = "Çarşamba"
	gunler["Thu"] = "Perşembe"
	gunler["Fri"] = "Cuma"
	gunler["Sat"] = "Cumartesi"
	gunler["Sun"] = "Pazar"

	gunler["Monday"] = "Pazartesi"
	gunler["Tuesday"] = "Salı"
	gunler["Wednesday"] = "Çarşamba"
	gunler["Thursday"] = "Perşembe"
	gunler["Friday"] = "Cuma"
	gunler["Saturday"] = "Cumartesi"
	gunler["Sunday"] = "Pazar"

	return gunler[gun]
}

//TarihCons test
func (testerModules Inow) TarihCons(myDateString string) {

	//bu ksımda saati ayırdık
	saat := strings.Split(myDateString, " ")
	saatFull := saat[1]
	fmt.Println(saatFull)

	saatWithoutMS := saatFull[:5]
	fmt.Println(saatWithoutMS)
	myDateString = saat[0] + " " + saatWithoutMS

	fmt.Println("Orginal Date:\t", myDateString)

	myDate, _ := time.Parse("2006-01-02 15:04", myDateString)
	trFormat := myDate.Format("02 Jan 2006 15:04, Mon")
	fmt.Println("Tr Format:\t\t", trFormat)
	fmt.Println("ING date:\t\t", myDate.Format("Mon, 02 Jan 2006 15:04"))

	this := Inow{}

	gun := strings.Split(trFormat, ",")
	gunFullEng := strings.TrimSpace(gun[1])
	gunFull := this.Gunler(strings.TrimSpace(gun[1]))
	fmt.Println(gunFull)

	ay := strings.Split(trFormat, " ")
	ayFullEng := ay[1]
	ayFull := this.AylarListe(ay[1])
	fmt.Println(ayFull)
	var result string
	result = strings.Replace(trFormat, ayFullEng, ayFull, -1)

	result = strings.Replace(result, gunFullEng, gunFull, -1)
	fmt.Println(result)

}

//tarihFinalInterface bal bla   	tarih.TarihCons("2020-05-17 05:08:40")
func tarihFinalInterface(myDateString string) string {

	//bu ksımda saati ayırdık
	saat := strings.Split(myDateString, " ")
	saatFull := saat[1]

	saatWithoutMS := saatFull[:5]

	myDateString = saat[0] + " " + saatWithoutMS

	myDate, _ := time.Parse("2006-01-02 15:04", myDateString)
	trFormat := myDate.Format("02 Jan 2006 15:04, Mon")

	this := Inow{}

	gun := strings.Split(trFormat, ",")
	gunFullEng := strings.TrimSpace(gun[1])
	gunFull := this.Gunler(strings.TrimSpace(gun[1]))

	ay := strings.Split(trFormat, " ")
	ayFullEng := ay[1]
	ayFull := this.AylarListe(ay[1])

	var result string
	result = strings.Replace(trFormat, ayFullEng, ayFull, -1)
	result = strings.Replace(result, gunFullEng, gunFull, -1)
	return result
}

func onlyDate(myDateString string) string {
	//bu ksımda saati ayırdık
	saat := strings.Split(myDateString, " ")
	saatFull := saat[1]

	saatWithoutMS := saatFull[:5]

	myDateString = saat[0] + " " + saatWithoutMS

	myDate, _ := time.Parse("2006-01-02 15:04", myDateString)
	trFormat := myDate.Format("02 Jan 2006")

	this := Inow{}

	ay := strings.Split(trFormat, " ")
	ayFullEng := ay[1]
	ayFull := this.AylarListe(ay[1])

	var result string
	result = strings.Replace(trFormat, ayFullEng, ayFull, -1)

	return result
}

//TarihFinalPointer d
/*
example
		var CreatedAt string
		 CreatedAt = smsList[0].CreatedAt.String()
	tarih.TarihFinalPointer(&CreatedAt)
*/
func (testerModules Inow) TarihFinalPointer(myDateString *string) {
	result := tarihFinalInterface(*myDateString)
	*myDateString = result

}

//TarihFinal cons
func (testerModules Inow) TarihFinal(myDateString string) string {
	return tarihFinalInterface(myDateString)
}

func (testerModules Inow) OnlyDate(myDateString string) string {
	return onlyDate(myDateString)
}

/*
"github.com/jinzhu/now"
func DateReturn() string {
	myConfig := &now.Config{
		WeekStartDay: time.Monday,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}
	dateim, _ := myConfig.Parse("2020-12-05 21:13:48")
	return dateim.String()
}
*/
/* ilk örnek
func DateReturn() string {
	input := "2020-12-05 21:13:48"
	layout := "2006-01-02"
	t, _ := time.Parse(layout, input)
	return t.Format("02-Jan-2006")
}
*/
