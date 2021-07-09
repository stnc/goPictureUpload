package controller

import (
	"fmt"
	"net/http"
	"regexp"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/services"

	"strconv"

	"github.com/astaxie/beego/utils/pagination"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

//WebArchive constructor
type WebArchive struct {
	webArchiveApp     services.WebArchiveAppInterface
	webArchiveLinkApp services.WebArchiveLinksAppInterface

	userApp services.UserAppInterface
}

const webArchivePath = "admin/WEB/webarchive/"

//InitWebArchive webarchive controller constructor
func InitWebArchive(wApp services.WebArchiveAppInterface, wLinkApp services.WebArchiveLinksAppInterface, uApp services.UserAppInterface) *WebArchive {
	return &WebArchive{
		webArchiveApp:     wApp,
		webArchiveLinkApp: wLinkApp,

		userApp: uApp,
	}
}

//Index list
func (access *WebArchive) Index(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	paginator = &pagination.Paginator{}
	var total int64
	access.webArchiveApp.Count(&total)
	PerPage := 8
	paginator := pagination.NewPaginator(c.Request, PerPage, total)
	offset := paginator.Offset()
	posts, _ := access.webArchiveApp.GetAllP(PerPage, offset)
	viewData := pongo2.Context{
		"paginator": paginator,
		"title":     "İçerik Ekleme",
		"posts":     posts,
		"csrf":      csrf.GetToken(c),
	}
	c.HTML(
		http.StatusOK,
		webArchivePath+"index.html",
		viewData,
	)
}

// Create webarchive
// Create webarchive
func (access *WebArchive) Create(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	viewData := pongo2.Context{
		"title": "WebArchive Add",
		"csrf":  csrf.GetToken(c),
	}
	/*
			str1 := `
			https://www.haberturk.com/kayseri-haberleri/82582582-erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi
		https://www.sabah.com.tr/kayseri/2020/11/23/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi
		https://www.sanayigazetesi.com.tr/ar-ge/icindeki-girisimciyi-fark-et-h24655.html
		https://www.talasexpresshaber.com/haber/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-344941
		https://www.ekovitrin.com/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h148866.html
		https://www.detayhaberler.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi/25073/
		https://beyazgazete.com/haber/2020/11/23/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-5848756.html
		https://www.menemenhaber.com.tr/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi/22505/
		https://www.hatayinternettv.com/haber/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-279789
		https://www.sondakika.com/haber/haber-erciyes-teknopark-icindeki-girisimciyi-fark-et-13753622/
		http://www.kayseriulkergazetesi.com/haber/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-24414.html
		http://www.kayserihaber.com.tr/haber/erciyes_teknopark_icindeki_girisimciyi_fark_et_etkinligi_duzenledi-50086.html
		http://www.kentturktv.com/2020/11/23/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi/
		https://www.kayseriyerelhaber.com/gundem/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h32216.html
		https://www.efetv.com.tr/erciyes-teknopark--icindeki-girisimciyi-fark-et--etkinligi-duzenledi-102718.html
		https://www.nehaber24.com/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h470913.html
		https://www.haberofisi.com.tr/tr/haberler/sabah/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-cs1bSHrAHlWA
		http://www.malatyaguncel.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-1654583h.htm
		http://www.kayserihaber.com.tr/haber/erciyes_teknopark_icindeki_girisimciyi_fark_et_etkinligi_duzenledi-50086.html
		https://www.haberekspresi.net/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-6942h.htm
		https://www.kayseriolay.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h68861.htm
		https://www.samsungazetesi.com/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h1403542.html
		https://www.iha.com.tr/kayseri-haberleri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-2752234/
		https://www.bursatv.com.tr/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h833047.html
		https://www.ayakligaste.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi/1896/
		http://www.haberalanya.com.tr/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h251071.html
		https://www.samsungazetesi.com/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h1403542.html
		https://kayserigurhaber.com/haber/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-47617.html
		https://www.inegolonline.com/guncel-olaylar/haber/erciyes-teknopark-icindeki-girisimciyi-fark-et-etk-1709864/
		https://www.haber16.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi/247099/
		https://www.kamu3.com/global-girisimcilik-haftasi-dolayisiyla-erciyes-teknopark-duzenledigi-basa/212943/
		https://www.kayseristarhaber.com.tr/ekonomi/icindeki-girisimciyi-fark-et-etkinligi-h10307.html
		http://www.haberlerantalya.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi/
		http://www.kayserikent.com/site/page.asp?dsy_id=120925&t=Erciyes-Teknopark-‘Icindeki-Girisimciyi-Fark-Et’-Etkinligi-Duzenledi
		http://yenidoganhaber.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi/316746/
		https://www.kayseriviphaber.com/icindeki-girisimciyi-fark-et/1706/
		http://www.ertehaber.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-148597h.htm
		http://www.bursahaber.com/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h1918470.html
		https://kayserihakimiyet2000.com/haber/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-10867.html
		https://www.haberler.com/erciyes-teknopark-icindeki-girisimciyi-fark-et-13753622-haberi
		https://www.gazeterize.com/kayseri/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-h14237386.html
		https://www.haberekspresi.net/erciyes-teknopark-icindeki-girisimciyi-fark-et-etkinligi-duzenledi-6942h.htm
			`
	*/

	//     var str = `https://www.geeksforgeeks.org/lru-cache-implementation/#hastest?name=234ksldjfkl&sjkldf=
	// https://developer.ibm.com/series/kubernetes-learning-path/
	// https://www.google.com/search?q=cannot+find+package+%22context%22+in+any+of%3A&oq=cannot+find+package+%22context%22+in+any+of%3A&aqs=chrome..69i57j0l5.205752j0j4&sourceid=chrome&ie=UTF-8#hello world

	// https://myaccount.google.com/name?utm_source=google-account&utm_medium=web&pli=1&rapt=AEjHL4PXc8zLqgFQOXxYQRbdKJI8\nyNq3u&name=rjl@1239015423#Kkn1HwudmOHSj7MpYAguZoOmrznrl0PF7hyLkZ17xFdrXToG5UDkMhO2bM1a6i5xw#path=234&sp=\n`
	// var re = regexp.MustCompile(`(?m)(?P<origin>(?P<protocol>http[s]?:)?\/\/(?P<host>[a-z0-9A-Z-_.]+))(?P<port>:\d+)?(?P<path>[\/a-zA-Z0-9-\.]+)?(?P<search>\?[^#\n]+)?(?P<hash>#.*)?`)
	// for i, match := range re.FindAllString(str1, -1) {
	// 	fmt.Println(match, "found at index", i)
	// }

	// windows := strings.Replace(str1, "\n", "\r\n", -1)
	// fmt.Printf("%+q\n", windows)
	// 	fmt.Printf("%+q\n", strings.Split(strings.Replace(windows, "\r\n", "", -1), "\t"))
	// kvs := strings.Split(strings.Replace(windows, "\r\n", "", -1), "\t")
	// for _, url := range kvs {

	// 	fmt.Println(url)
	// }
	c.HTML(
		http.StatusOK,
		webArchivePath+"create.html",
		viewData,
	)
}

//Store save method
func (access *WebArchive) Store(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	var webarchive, _, _ = webArchiveModel(c)
	fmt.Println(webarchive)
	var saveError = make(map[string]string)
	saveError = webarchive.Validate()

	if len(saveError) == 0 {
		saveData, saveErr := access.webArchiveApp.Save(&webarchive)
		if saveErr != nil {
			saveError = saveErr
		}
		lastID := strconv.FormatUint(uint64(saveData.ID), 10)
		linksAll := c.PostForm("LinksAll")
		var re = regexp.MustCompile(`(?m)(?P<origin>(?P<protocol>http[s]?:)?\/\/(?P<host>[a-z0-9A-Z-_.]+))(?P<port>:\d+)?(?P<path>[\/a-zA-Z0-9-\.]+)?(?P<search>\?[^#\n]+)?(?P<hash>#.*)?`)
		for i, match := range re.FindAllString(linksAll, -1) {
			fmt.Println(match, "found at index", i)
			// webarchivelink entity.WebArchiveLinks{}
			var webarchivelink = entity.WebArchiveLinks{}
			webarchivelink.Link = match
			webarchivelink.WebArchiveID = saveData.ID
			access.webArchiveLinkApp.Save(&webarchivelink)
		}

		stncsession.SetFlashMessage("Record successfull added", "success", c)
		c.Redirect(http.StatusMovedPermanently, "/admin/webarchive/edit/"+lastID)
		return
	}

	viewData := pongo2.Context{
		"title":      "webarchive create",
		"csrf":       csrf.GetToken(c),
		"err":        saveError,
		"post":       webarchive,
		"webarchive": webarchive,
	}
	c.HTML(
		http.StatusOK,
		webArchivePath+"create.html",
		viewData,
	)
}

//Edit edit data
func (access *WebArchive) Edit(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	flashMsg := stncsession.GetFlashMessage(c)
	if ID, err := strconv.ParseUint(c.Param("ID"), 10, 64); err == nil {
		if posts, err := access.webArchiveApp.GetByID(ID); err == nil {
			viewData := pongo2.Context{
				"title":    "webarchive edit",
				"post":     posts,
				"csrf":     csrf.GetToken(c),
				"flashMsg": flashMsg,
			}
			c.HTML(
				http.StatusOK,
				webArchivePath+"edit.html",
				viewData,
			)
		} else {
			c.AbortWithError(http.StatusNotFound, err)
		}
	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

//Update data
func (access *WebArchive) Update(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)
	var webarchive, _, id = webArchiveModel(c)

	var savePostError = make(map[string]string)
	savePostError = webarchive.Validate()
	if len(savePostError) == 0 {
		_, saveErr := access.webArchiveApp.Update(&webarchive)
		if saveErr != nil {
			savePostError = saveErr
		}
		stncsession.SetFlashMessage("Record successfully edit", "success", c)
		c.Redirect(http.StatusMovedPermanently, "/admin/webarchive/edit/"+id)
		return
	}

	viewData := pongo2.Context{
		"title": "içerik ekleme",
		"err":   savePostError,
		"csrf":  csrf.GetToken(c),
		"post":  webarchive,
	}
	c.HTML(
		http.StatusOK,
		webArchivePath+"edit.html",
		viewData,
	)
}

//Delete data
func (access *WebArchive) Delete(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	if ID, err := strconv.ParseUint(c.Param("ID"), 10, 64); err == nil {

		access.webArchiveApp.Delete(ID)
		stncsession.SetFlashMessage("Success Delete", "success", c)

		c.Redirect(http.StatusTemporaryRedirect, "/admin/webarchive")
		return

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

func (access *WebArchive) LinkPdfRun(c *gin.Context) {
	stncsession.IsLoggedInRedirect(c)

	if ID, err := strconv.ParseUint(c.Param("ID"), 10, 64); err == nil {

		access.webArchiveApp.Delete(ID)
		stncsession.SetFlashMessage("Success RUN", "success", c)

		c.Redirect(http.StatusMovedPermanently, "/admin/webarchive/edit/"+c.Param("ID"))
		return

	} else {
		c.AbortWithStatus(http.StatusNotFound)
	}
}

//form webarchive model
func webArchiveModel(c *gin.Context) (webarchive entity.WebArchive, idD uint64, idStr string) {
	id := c.PostForm("ID")
	title := c.PostForm("Title")
	linksAll := c.PostForm("LinksAll")
	excerpt := c.PostForm("Excerpt")
	idInt, _ := strconv.Atoi(id)
	var idN uint64
	idN = uint64(idInt)
	//	var webarchive = entity.WebArchive{}
	webarchive.ID = idN
	webarchive.UserID = 1
	webarchive.Title = title
	webarchive.Status = 1
	webarchive.LinksAll = linksAll
	webarchive.Excerpt = excerpt
	return webarchive, idN, id
}
