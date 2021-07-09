package dto

//burayı kişi ekleme ajax işlemleri kullanıyor

//KurbanTableName table name
var KisilerTableName string = "kisiler"

//Kurban dto
type Kisiler struct {
	ID      uint64
	AdSoyad string
	Telefon string
	Adres   string
}

type KisilerDataTableResult struct {
	Total    int       `json:"recordsTotal"`
	Filtered int       `json:"recordsFiltered"`
	Data     []Kisiler `json:"data"`
}
