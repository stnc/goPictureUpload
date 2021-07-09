package msgedit

import (
	"fmt"
	"net/url"
	"stncCms/app/domain/repository"
	"strings"
)

//oku
//https://faq.whatsapp.com/general/chats/how-to-use-click-to-chat/?lang=tr

type (
	// Inow : Custom Renderer for templates
	Msg struct{ Debug bool }
)

func (msgModul Msg) Wp(slug string) string {
	return onlyDate(slug)
}

func telReplace(tel string) string {
	return strings.Replace(tel, " ", "", -1)
}

func msgReplace(msg string, slug string, url string) string {
	return strings.Replace(msg, "[link]", url+"kurbanBilgi/"+slug, -1)
}
func onlyDate(slug string) string {
	db := repository.DB
	access := repository.KurbanRepositoryInit(db)
	// kurbanData, _ := access.GetKurbanOpenInfo(slug)

	if kurbanData, err := access.GetKurbanOpenInfo(slug); err == nil {

		tel := telReplace(kurbanData.Telefon)

		appOptions := repository.OptionRepositoryInit(db)
		msgg := msgReplace(appOptions.GetOption("whatsAppMsg"), slug, appOptions.GetOption("siteUrl"))
		fmt.Println(appOptions.GetOption("siteUrl"))
		return "https://wa.me/09" + tel + "?text=" + url.QueryEscape(msgg)
	} else {
		return ""
	}
}
