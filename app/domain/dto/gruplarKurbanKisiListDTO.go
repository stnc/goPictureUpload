package dto

type GruplarExcelandIndex struct {
	ID                       uint64  `json:"id"`
	UserID                   uint64  `json:"userId"`
	HayvanBilgisiID          uint64  `json:"hayvanBilgisiID"  `
	GrupAdi                  string  `json:"grupAdi"`
	KesimSiraNo              int     `json:"kesimSiraNo"`
	HissedarAdet             int     `json:"hissedarAdet"`
	Aciklama                 string  `json:"aciklama"`
	SatisFiyatTuru           int     `json:"satisFiyatTuru"` //ayalardaki satış fiyatlraında hangisi kullanımış
	SatisFiyati              float64 `json:"satisFiyati" `   //ayarlardaki kilo başı satış fiyatının hangisi seçilmiş
	Siralama                 int     `json:"Siralama"`
	AgirlikTipi              int     `json:"agirlikTipi"` // yüksek ağırlık , orta ağırlık , düşük ağırlık gibi
	Durum                    int     `json:"durum" `
	Slug                     string  `json:"slug"`
	ToplamKurbanFiyati       float64 `json:"toplamKurbanFiyati"`
	KurbanBayramiYili        int     `json:"kurbanBayramiYili"` //otomaitk atılacak isitem tarafından ayarlradan ceksin
	Agirlik                  int     `json:"agirlik"`
	KupeNo                   string  `json:"kupeNo"`
	KisiBasiDusenHisseFiyati float64 `json:"kisiBasiDusenHisseFiyati"`
	ToplamOdemeler           float64 `json:"toplamOdemeler"`
	KalanBorclar             float64 `json:"kalanBorclar"`
	KasaBorcu                float64 `json:"kasaBorcu"`

	GrupIsoTopeTRname string `gorm:"-" `
	GrupIsoTopeName   string `gorm:"-" `
	GrupIsoTopeAlert  string `gorm:"-" `

	KurbanKisiList []KurbanListForGrouplar
}

type KurbanListForGrouplar struct {
	ID            uint64  `json:"id"`
	GrupID        uint64  `json:"grup_id"`
	KisiID        uint64  `json:"kisi_id"`
	VekaletDurumu int     `json:"vekalet"`
	KurbanTuru    int     `json:"kurban_turu"`
	Durum         int     `json:"durum"`
	BorcDurum     int     `json:"borc_durum"`
	GrupLideri    int     `json:"grup_lideri"`
	KurbanFiyati  float64 `json:"kurban_fiyati"`
	Borc          float64 `json:"borc"`
	Alacak        float64 `json:"alacak"`
	Bakiye        float64 `json:"bakiye"`
	RefKisiID     uint64  `json:"referans_kisi_id"`
	KisiAdSoyad   string  `json:"kisi_ad_soyad"` //join
	KisiTelefon   string  `json:"kisitelefon"`   //join
	Slug          string  `json:"slug"`

	ReferansID      uint64 `json:"referansID"`
	ReferansGrupID  uint64 `json:"referansGrupID"`
	ReferansAdSoyad string `json:"referansAdSoyad"`
	ReferansTelefon string `json:"referansTelefon"`
}
