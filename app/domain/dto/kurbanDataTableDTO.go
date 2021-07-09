package dto

import (
	"time"
)

//KurbanTable strcut listelemde kullanılıyor
type KurbanTable struct {
	ID            uint64     `json:"id"`
	GrupID        uint64     `json:"grup_id"`
	KisiID        uint64     `json:"kisi_id"`
	VekaletDurumu int        `json:"vekalet"`
	KurbanTuru    int        `json:"kurban_turu"`
	Durum         int        `json:"durum"`
	BorcDurum     int        `json:"borc_durum"`
	KurbanFiyati  float64    `json:"kurban_fiyati"`
	Borc          float64    `json:"borc"`
	Alacak        float64    `json:"alacak"`
	Bakiye        float64    `json:"bakiye"`
	Tarih         *time.Time `json:"tarih"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
	KisiGrupID    uint64     `json:"kisi_grup_id"`
	KisiAdSoyad   string     `json:"kisi_ad_soyad"`
	KisiTelefon   string     `json:"telefon"`
}
type KurbanBilgisiDataTableResult struct {
	Total    int           `json:"recordsTotal"`
	Filtered int           `json:"recordsFiltered"`
	Data     []KurbanTable `json:"data"`
}

type KurbanOpenInfoData struct {
	KurbanGrupId    uint64  `json:"KurbanGrupId"`
	KurbanID        uint64  `json:"KurbanID"`
	KurbanDurum     int     `json:"KurbanDurum"`
	KurbanBorcDurum int     `json:"KurbanBorcDurum"`
	KurbanBorc      float64 `json:"KurbanBorc"`
	KurbanAlacak    float64 `json:"KurbanAlacak"`
	KurbanBakiye    float64 `json:"kurbanBakiye"`
	KurbanFiyati    float64 `json:"kurbanFiyati"`
	GrupAdi         string  `json:"GrupAdi"`
	KesimNo         int     `json:"KesimNo"`
	HissedarAdet    int     `json:"hissedarAdet"`
	AdSoyad         string  `json:"AdSoyad"`
	Telefon         string  `json:"Telefon"`
	KisiId          uint64  `json:"KisiId"`
	HayvanAgirlik   uint64  `json:"hayvanAgirlik"`
	HayvanBilgisiId uint64  `json:"HayvanBilgisiId"`
}
