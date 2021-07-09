package dto

//GruplarTableName table name

//Kurban dto
type Gruplar struct {
	ID                 uint64
	HayvanBilgisiID    uint64
	KesimSiraNo        int
	HissedarAdet       int
	KurbanFiyatiTipi   int
	ToplamKurbanFiyati float64
	Agirlik            int
	//Kurban             []Kurban

}

//Kurban dto
type KurbanUpdateRead struct {
	ID           uint64
	KurbanFiyati float64
	KasaBorcu    float64
	KalanUcret   float64
	BorcDurum    int
	Agirlik      int
}
