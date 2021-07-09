package repository

import (
	"fmt"
	"os"
	"stncCms/app/domain/entity"
	"stncCms/app/services"

	"github.com/hypnoglow/gormzap"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	_ "github.com/lib/pq" // here
	_ "gorm.io/driver/mysql"
	_ "gorm.io/driver/postgres"
)

var DB *gorm.DB

//Repositories strcut
type Repositories struct {
	User               services.UserAppInterface
	Post               services.PostAppInterface
	Cat                services.CatAppInterface
	CatPost            services.CatPostAppInterface
	Lang               services.LanguageAppInterface
	HayvanSatisYerleri services.HayvanSatisYerleriAppInterface
	HayvanBilgisi      services.HayvanBilgisiAppInterface
	Kurban             services.KurbanAppInterface
	Kodemeler          services.OdemelerAppInterface
	Gruplar            services.GruplarAppInterface
	Kisiler            services.KisilerAppInterface
	WebArchive         services.WebArchiveAppInterface
	WebArchiveLink     services.WebArchiveLinksAppInterface
	Options            services.OptionsAppInterface
	Media              services.MediaAppInterface
	DB                 *gorm.DB
}

//DbConnect initial
/*TODO: burada db verisi pointer olarak işaretlenecek oyle gidecek veri*/
func DbConnect() *gorm.DB {
	dbdriver := os.Getenv("DB_DRIVER")
	dbHost := os.Getenv("DB_HOST")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	gormAdvancedLogger := os.Getenv("GORM_ZAP_LOGGER")
	debug := os.Getenv("MODE")
	//	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPassword) //bu postresql

	//DBURL := "root:sel123C#@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local" //mysql
	var DBURL string

	if dbdriver == "mysql" {
		DBURL = dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local"
	} else if dbdriver == "postgres" {
		DBURL = fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbUser, dbName, dbPassword) //Build connection string
	}
	db, err := gorm.Open(dbdriver, DBURL)

	// }
	// db, err := gorm.Open(dbdriver, DBURL)
	//nunlar gorm 2 versionunda prfexi falan var
	// db, err := gorm.Open(mysql.Open(DBURL), &gorm.Config{
	// 	NamingStrategy: schema.NamingStrategy{
	// 		TablePrefix:   "krbn_", // table name prefix, table for `User` would be `t_users`
	// 		SingularTable: true,    // use singular table name, table for `User` would be `user` with this option enabled
	// 	},
	// 	// Logger: gorm_logrus.New(),
	// })

	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}

	if debug == "DEBUG" && gormAdvancedLogger == "ENABLE" {
		db.LogMode(true)
		log := zap.NewExample()
		db.SetLogger(gormzap.New(log, gormzap.WithLevel(zap.DebugLevel)))
	} else if debug == "DEBUG" || debug == "TEST" && gormAdvancedLogger == "ENABLE" {
		db.LogMode(true)
	} else if debug == "RELEASE" {
		fmt.Println(debug)
		db.LogMode(false)
	}
	DB = db

	db.SingularTable(true)

	return db
}

//https://techinscribed.com/different-approaches-to-pass-database-connection-into-controllers-in-golang/

//RepositoriesInit initial
func RepositoriesInit(db *gorm.DB) (*Repositories, error) {

	return &Repositories{
		User:               UserRepositoryInit(db),
		Post:               PostRepositoryInit(db),
		Cat:                CatRepositoryInit(db),
		CatPost:            CatPostRepositoryInit(db),
		Lang:               LanguageRepositoryInit(db),
		Kurban:             KurbanRepositoryInit(db),
		Kodemeler:          OdemelerRepositoryInit(db),
		Gruplar:            GruplarRepositoryInit(db),
		Kisiler:            KisilerRepositoryInit(db),
		HayvanSatisYerleri: HayvanSatisYerleriRepositoryInit(db),
		HayvanBilgisi:      HayvanBilgisiRepositoryInit(db),
		WebArchive:         WebArchiveRepositoryInit(db),
		WebArchiveLink:     WebArchiveLinksRepositoryInit(db),
		Options:            OptionRepositoryInit(db),
		Media:              MediaRepositoryInit(db),

		DB: db,
	}, nil
}

//Close closes the  database connection
// func (s *Repositories) Close() error {
// 	return s.db.Close()
// }

//Automigrate This migrate all tables
func (s *Repositories) Automigrate() error {
	s.DB.AutoMigrate(&entity.Kurban{}, &entity.Odemeler{}, &entity.Gruplar{}, &entity.Kisiler{}, &entity.User{},
		&entity.HayvanSatisYerleri{}, &entity.Languages{}, &entity.Modules{}, &entity.Notes{},
		&entity.HayvanBilgisi{}, &entity.Options{}, &entity.Media{})

	// &entity.Post{}, &entity.Categories{}, &entity.CategoryPosts{}, &entity.WebArchive{}, &entity.WebArchiveLinks{}
	s.DB.Model(&entity.HayvanBilgisi{}).AddForeignKey("hayvan_satis_yerleri_id", "hayvan_satis_yerleri(id)", "CASCADE", "CASCADE") // one to one (one=hayvan_satis_yerleri) (one=HayvanBilgisi)

	s.DB.Model(&entity.Odemeler{}).AddForeignKey("kurban_id", "kurbanlar(id)", "CASCADE", "CASCADE")        // one to many (one=kurbanlar) (many=odemeler)
	s.DB.Model(&entity.Kurban{}).AddForeignKey("grup_id", "gruplar(id)", "CASCADE", "CASCADE")              // one to many (one=gruplar) (many=kurbanlar)
	return s.DB.Model(&entity.Kurban{}).AddForeignKey("kisi_id", "kisiler(id)", "CASCADE", "CASCADE").Error // one to many (one=kisiler) (many=kurbanlar)
	// s.DB.Model(&entity.WebArchiveLinks{}).AddForeignKey("web_archive_id", "web_archive(id)", "CASCADE", "CASCADE").Error // one to many (one=web_archives) (many=WebArchiveLinks)

}
