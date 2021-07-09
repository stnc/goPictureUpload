package stnchelper

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func SendSms(msg string, phone string) {
	url := "http://api.tescom.com.tr:8080/api/smspost/v1"
	postData := `<sms>
		<username>Evma</username>
		<password>bosssss</password>
		<header>EMMs</header>
		<validity>2880</validity>
		<message>
		<gsm>
		<no></no>
		</gsm>
		<msg><![CDATA[-test-message-]]></msg>
		</message>
		</sms>`

	payload := strings.NewReader(postData)
	req, _ := http.NewRequest("POST", url, payload)

	// req.Header.Add("cache-control", "no-cache")
	// req.Header.Set("text/xml", "charset=UTF-8")

	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	stringSlice := strings.Split(string(body), " ")

	fmt.Println((stringSlice[0]))

}
