package dto

import "time"

//OdemelerTableName table name
var OdemelerTableName string = "odemeler"

//Odemeler dto layer
type Odemeler struct {
	ID             uint64    `json:"id"`
	KurbanID       uint64    `json:"kurbanId"`
	UserID         uint64    `json:"userId"`
	Aciklama       string    `json:"aciklama"`
	Makbuz         string    `json:"makbuz" `
	BorcDurum      int       `validate:"required"`
	Durum          int       `validate:"required"`
	KurbanFiyati   float64   `validate:"numeric" `
	VerilenUcret   float64   `validate:"required,numeric" `
	Alacak         float64   `validate:"numeric" `
	Borc           float64   `validate:"numeric" `
	Bakiye         float64   `validate:"numeric" `
	VerildigiTarih time.Time `json:"verildigiTarih"`
}

//OdemelerSonFiyat dto layer
type OdemelerSonFiyat struct {
	ID           uint64
	VerilenUcret float64 `validate:"required,numeric"`
	Alacak       float64 `validate:"required,numeric"`
	Borc         float64 `validate:"required,numeric"`
}

//TableName override
func (gk *Odemeler) TableName() string {
	return OdemelerTableName
}

type OdemeMakbuzu struct {
	OdemeID                uint64    `json:"odemeID"`
	Aciklama               string    `json:"Aciklama"`
	Makbuz                 string    `json:"Makbuz"`
	OdemeDurum             int       `json:"OdemeDurum"`
	OdemelerBorcDurum      int       `json:"OdemelerBorcDurum"`
	OdemelerKurbanFiyati   float64   `json:"OdemelerKurbanFiyati"`
	OdemelerVerilenUcret   float64   `json:"OdemelerVerilenUcret"`
	OdemelerVerildigiTarih time.Time `json:"OdemelerVerildigiTarih"`
	KurbanGrupId           uint64    `json:"KurbanGrupId"`
	KurbanID               uint64    `json:"KurbanID"`
	KurbanDurum            int       `json:"KurbanDurum"`
	KurbanBorcDurum        int       `json:"KurbanBorcDurum"`
	GrupAdi                string    `json:"GrupAdi"`
	KesimNo                int       `json:"KesimNo"`
	HissedarAdet           int       `json:"hissedarAdet"`
	AdSoyad                string    `json:"AdSoyad"`
	Telefon                string    `json:"Telefon"`
	KisiId                 uint64    `json:"KisiId"`
}
