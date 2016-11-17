package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

var classes = []interface{}{
	Applicant{},
	ApplicantData{},
	Oblast{},
	User{},
	ExamReference{},
}

func InitTestDb(db *gorm.DB) *gorm.DB {

	// Migrate the schema
	db.AutoMigrate(classes...)

	// initialize if db empty
	var nOblasts, nappl, nxref int
	if db.Model(&Oblast{}).Count(&nOblasts); nOblasts < 1 {
		fmt.Printf("No oblasts in db, creating new\n")
		for _, oblast := range InitialOblasts {
			db.Create(&oblast)
		}
	}
	// retrieve oblasts from db
	var oblasts []Oblast
	db.Find(&oblasts)

	if db.Model(&ExamReference{}).Count(&nxref); nxref < 1 {
		xref := ExamReference{
			Year:    time.Now().Year(),
			Results: [NQESTION]int{5, 5, 2, 3, 10, 15, 10, 10, 0, 0},
		}
		db.Create(&xref)
	}

	if db.Find(&Applicant{}).Count(&nappl); nappl < 1 {
		data := ApplicantData{
			LastName:    "Gans",
			FirstName:   "Gisbert",
			FathersName: "Gisbertovich",
			Phone:       "040441777",
			Home:        "Ducksburg",
			School:      "Nr. 14",
			Email:       "Giga@goosemail.com",
			Oblast:      oblasts[7],
			OrtSum:      111,
			OrtMath:     66,
			OrtPhys:     33,
			Results:     [NQESTION]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		goose := Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Gans",
			FirstName:   "Franz",
			FathersName: "Fietjevich",
			Phone:       "046774417",
			Home:        "Franzhausen",
			School:      "Nr. 66",
			Email:       "FraGa@goosemail.com",
			Oblast:      oblasts[6],
			OrtSum:      100,
			OrtMath:     66,
			OrtPhys:     33,
			Results:     [NQESTION]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Gans",
			FirstName:   "Gertrude",
			FathersName: "Gaggova",
			Phone:       "0467734567",
			Home:        "Ducksburg",
			School:      "Nr. 35",
			Email:       "Gerti@goosemail.com",
			Oblast:      oblasts[6],
			OrtSum:      178,
			OrtMath:     99,
			OrtPhys:     36,
			Results:     [NQESTION]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Duck",
			FirstName:   "Daisy",
			FathersName: "Waltova",
			Phone:       "08980080099",
			Home:        "München",
			School:      "Nr. 5",
			Email:       "Daisy@disney.com",
			Oblast:      oblasts[2],
			OrtSum:      230,
			OrtMath:     100,
			OrtPhys:     66,
			Results:     [NQESTION]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
		data = ApplicantData{
			LastName:    "Duck",
			FirstName:   "Donald",
			FathersName: "Disneyvich",
			Phone:       "08955443322",
			Home:        "Oberpfaffenhofen",
			School:      "Nr. 19",
			Email:       "Donald@disney.com",
			Oblast:      oblasts[2],
			OrtSum:      101,
			OrtMath:     50,
			Results:     [NQESTION]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
		goose = Applicant{Data: data}
		// create a new object with its dependents
		db.Create(&goose)
	}
	var nuser int
	if db.Model(&User{}).Count(&nuser); nuser < 1 {
		users := InitialUsers()
		for _, user := range users {
			db.Save(&user)
		}
	}
	return db
}

func ClearTestDb(db *gorm.DB) {

	for _, class := range classes {
		db.Delete(class)
		db.DropTableIfExists(class)
	}

}
