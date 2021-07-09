package dto

import (
	"time"
)

//HayvanSatisYerleriDataTable strcut
type HayvanSatisYerleriDataTable struct {
	ID         uint64    `json:"id"`
	YerAdi     string    `json:"yer_adi" `
	IlgiliKisi string    `json:"ilgili_kisi"`
	Adresi     string    `json:"adresi"`
	Telefon    string    `json:"telefon"`
	Durum      int       `json:"durum"`
	CreatedAt  time.Time `json:"created_at"`
}
type HayvanSatisYerleriDataTableResult struct {
	Total    int                           `json:"recordsTotal"`
	Filtered int                           `json:"recordsFiltered"`
	Data     []HayvanSatisYerleriDataTable `json:"data"`
}
