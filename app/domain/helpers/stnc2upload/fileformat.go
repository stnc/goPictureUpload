package stnc2upload

import (
	"path"
	"path/filepath"
	"stncCms/app/domain/helpers/stnchelper"
	"strings"

	"github.com/twinj/uuid"
)

func FormatFile(fn string) string {

	ext := path.Ext(fn)
	u := uuid.NewV4()

	newFileName := u.String() + ext

	return newFileName
}

func newFileNameFunc(filenameOrg string) string {

	filenameExtension := filepath.Ext(filenameOrg)

	realFilename := strings.Split(filenameOrg, ".")

	realFilenameSlug := stnchelper.GenericName(realFilename[0], 50)

	filename := realFilenameSlug + filenameExtension

	return filename
}
