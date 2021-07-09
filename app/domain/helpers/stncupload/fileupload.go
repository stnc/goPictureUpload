package stncupload

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"stncCms/app/domain/entity"
	"stncCms/app/domain/helpers/stnccollection"
	"stncCms/app/domain/helpers/stncsession"
	"stncCms/app/domain/repository"

	"github.com/disintegration/imaging"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

/*
direk erişimli kullanım
	filename, uploadError = stncupload.NewFileUpload().UploadFile(filenameForm, c.PostForm("Resim2"))

	_, filetype := stncupload.NewFileUpload().RealImage("public/upl/kiosk/" + filename)
*/

/*
kullanım 2
	upl := stncupload.FileUpload{UploadPath: "public/upl/kiosk2/"}
	filename, uploadError = upl.UploadFile(filenameForm, c.PostForm("Resim2"))
*/

//buradan eşiebilmesi için func (fu FileUpload) Uplo bolyle olmalı fonksiyon
//NewFileUpload  construct s
func NewFileUpload() *FileUpload {
	return &FileUpload{}
}

type FileUpload struct {
	UploadPath string
	UploadSize int
	Types      []string
	MaxFiles   int
}

//UploadFileInterface interface
type UploadFileInterface interface {
	UploadFileForMinio(file *multipart.FileHeader) (string, error)
	UploadFile(filest *multipart.FileHeader) (string, string)
	MultipleUploadFile(filest []*multipart.FileHeader, originalName string)
	RealFileType(fileName string) (bool, string)
	InitUploader(c *gin.Context, IDint int, filenameForm *multipart.FileHeader, accData map[string]string)
}

//So what is exposed is Uploader
var _ UploadFileInterface = &FileUpload{}

//TODO: https://github.com/gin-gonic/examples/tree/master/upload-file upload ornekleri var
//TODO: gerçek resim dosayasını tespit eden fonksiyon başka yere alınablir
//TODO: boyutlandırma https://github.com/disintegration/imaging
//https://socketloop.com/tutorials/golang-how-to-verify-uploaded-file-is-image-or-allowed-file-types
//https://www.golangprograms.com/how-to-get-dimensions-of-an-image-jpg-jpeg-png-or-gif.html
//UploadFile standart upload
func (fu FileUpload) UploadFile(filest *multipart.FileHeader) (filename string, errorReturn string) {
	var uploadFilePath string = fu.UploadPath

	// var deleteFilename string
	// var filename string
	// var errorReturn string

	if filest != nil {
		f, err := filest.Open()
		defer f.Close()
		if err != nil {
			errorReturn = err.Error()
		}

		if filest.Header != nil {

			size := filest.Size
			// var size2 = strconv.FormatUint(uint64(size), 10)
			if size > int64(1024000*fu.UploadSize) { // 1 MB
				uploadSizeStr := stnccollection.IntToString(fu.UploadSize)

				errorReturn = "HATA: Resim boyutu çok yüksek maximum " + uploadSizeStr + " MB olmalıdır" //+ size2
				filename = "false"
			}

			filename = newFileNameFunc(filest.Filename)
			// deleteFilename = filename

			fmt.Println(filename)

			out, err := os.Create(uploadFilePath + filename)

			defer out.Close()

			if err != nil {
				log.Fatal(err)
				errorReturn = err.Error()
				filename = "false"
			}

			_, err = io.Copy(out, f)

			if err != nil {
				log.Fatal(err)
				errorReturn = err.Error()
				filename = "false"
			}

		}
	}
	return filename, errorReturn
}

func (fu FileUpload) RealFileType(fileName string) (bool, string) {
	var uploadFilePath string = fu.UploadPath
	var typeControlError bool = false
	// open the uploaded file
	file, err := os.Open(uploadFilePath + fileName)
	defer file.Close()
	if err != nil {
		//TODO: buraya log koymak gerekiyor
		fmt.Println(err)
		err.Error()
		// os.Exit(1)
	}

	buff := make([]byte, 512) // why 512 bytes ? see http://golang.org/pkg/net/http/#DetectContentType
	_, err = file.Read(buff)

	if err != nil {
		fmt.Println(err)
		err.Error()
		// os.Exit(1)
	}

	filetype := http.DetectContentType(buff)

	fmt.Println(fu.Types)
	_, typeControlError = stnccollection.FindSlice(fu.Types, filetype)

	fmt.Println(filetype)
	fmt.Println("typeControlError")
	fmt.Println(typeControlError)
	return typeControlError, filetype
}

//MultipleUploadFile dropzone olduğu için henuz kullanımda değildir
func (fu FileUpload) MultipleUploadFile(files []*multipart.FileHeader, originalName string) {
	var uploadFilePath string = "public/upl/kiosk/"

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		fmt.Println(files[i].Filename)
		defer file.Close()
		if err != nil {
			// fmt.Fprintln(w, err)
			return
		}

		out, err := os.Create(uploadFilePath + files[i].Filename)

		defer out.Close()
		if err != nil {
			// fmt.Fprintf(w, "Unable to create the file for writing. Check your write access privilege")
			return
		}

		_, err = io.Copy(out, file) // file not files[i] !

		if err != nil {
			// fmt.Fprintln(w, err)
			return
		}

		fmt.Println("Files uploaded successfully : ")
		fmt.Println(files[i].Filename + "\n")

	}
}

func (fu FileUpload) InitUploader(c *gin.Context, IDint int, filenameForm *multipart.FileHeader, accData map[string]string) {
	var uploadError string
	var filename string

	uploadPath := accData["uploadPath"]

	db := repository.DB
	access := repository.MediaRepositoryInit(db)

	maxFiles := stnccollection.StringToint(accData["maxFiles"])

	_, err := os.Stat(uploadPath)
	if os.IsNotExist(err) {
		uploadError = uploadPath + " klasörü yok "
		//Create a folder/directory at a full qualified path
		err = os.Mkdir(uploadPath, 0755)
		if err != nil {
			// log.Fatal(err)
			log.Default()
		}

		err = os.Mkdir(uploadPath+"big/", 0755)
		if err != nil {
			// log.Fatal(err)
			log.Default()
		}

		err = os.Mkdir(uploadPath+"thumb/", 0755)
		if err != nil {
			// log.Fatal(err)
			log.Default()
		}
	}

	var total int

	mediaModulID := stnccollection.StringToint(accData["modulID"])

	access.Count(mediaModulID, IDint, &total)

	if total >= maxFiles {
		uploadError = "Yükleme sayısını aştınız maksimum " + stnccollection.IntToString(maxFiles) + " dosya yükleyebilirsiniz"
		c.JSON(http.StatusBadGateway, uploadError)
		return
	}

	filename, uploadError = fu.UploadFile(filenameForm)

	// countTypes := len(upl.Types)

	var typeControlError bool = true

	var filetype string

	typeControlError, filetype = fu.RealFileType(filename)

	switch filetype {
	case "image/jpeg", "image/jpg":
		fu.ImageResize(filename, uploadPath, accData)
		FileDelete(uploadPath, filename)
	case "image/gif":
		fu.ImageResize(filename, uploadPath, accData)
		FileDelete(uploadPath, filename)
	case "image/png":
		fu.ImageResize(filename, uploadPath, accData)
		FileDelete(uploadPath, filename)
		// default:
		// 	returnData = false
	}

	var mediaData = entity.Media{}
	mediaData.MediaName = filename
	mediaData.ModulID = stnccollection.StringtoUint64(accData["modulID"])
	mediaData.ContentID = uint64(IDint)
	mediaData.UserID = stncsession.GetUserID2(c)
	mediaData.MimeType = filetype

	saveData, saveErr := access.Save(&mediaData)

	if saveErr != nil {
		uploadError = "veritabanı hatası"
	}

	var returnID int
	returnID = saveData.ID
	if typeControlError == false {
		// stringSlice := strings.Split(filetype, "/") stringSlice[1] //yok açma
		access.Delete(uint64(returnID)) //aç
		FileDelete(uploadPath, filename)
		uploadError = "HATA: Yükleyeceğiniz dosya tipi uygun değildir, bunlar biri olmalıdır mp4, webm, jpeg, jpg, gif, png "
		c.JSON(http.StatusBadGateway, uploadError)
		return
	}

	// https://play.golang.org/p/UKZbcuJUPP

	if uploadError == "" {
		// c.JSON(http.StatusOK, "Başarı ile yuklendi") //TODO: v2 de acılması gerekiyor
		c.JSON(http.StatusBadGateway, "Başarı ile yuklendi, sayfayı yenilediğinizde resimler gelecektir")

	} else {
		c.JSON(http.StatusBadGateway, uploadError)
	}

}

//imageResize image resize
func (fu FileUpload) ImageResize(filename string, path string, accData map[string]string) {

	fmt.Println("girer")
	src, err := imaging.Open(path + filename)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	// src = imaging.Resize(src, 1815, 2510, imaging.Center)
	// var srcBig image.Image

	widthBig := stnccollection.StringToint(accData["bigImageWidth"])
	// heightHeight := stnccollection.StringToint(accData["bigImageHeight"])

	weightThumb := stnccollection.StringToint(accData["thumbImageWidth"])
	heightThumb := stnccollection.StringToint(accData["thumbImageHeight"])

	var srcBig = imaging.Resize(src, widthBig, 0, imaging.Lanczos)

	errBig := imaging.Save(srcBig, path+"big/"+filename)
	if errBig != nil {
		log.Fatalf("failed to save image: %v", errBig)
	}

	var srcThumb = imaging.Resize(src, weightThumb, heightThumb, imaging.Lanczos)

	errth := imaging.Save(srcThumb, path+"thumb/"+filename)
	if errth != nil {
		log.Fatalf("failed to save image: %v", errth)
	}
}

//MediaDelete Delete file data
func (fu FileUpload) MediaDelete(c *gin.Context, ID uint64, uploadPath string) {

	db := repository.DB
	access := repository.MediaRepositoryInit(db)

	if mediaData, err := access.GetByID(ID); err == nil {
		mediaName := mediaData.MediaName
		FileDelete(uploadPath, mediaName)
		FileDelete(uploadPath+"big/", mediaName)
		FileDelete(uploadPath+"thumb/", mediaName)
		fmt.Println(mediaData)
	}

	access.Delete(ID)

	viewData := pongo2.Context{
		"status": "ok",
		"msg":    "Kayıt Başarı ile Silindi",
	}
	fmt.Println("girer")
	c.JSON(http.StatusOK, viewData)

	// c.JSON(http.StatusBadGateway, uploadError)

	return

}
