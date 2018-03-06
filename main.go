package main

import (
	"appealimport/importdb"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
)

var config = struct {
	DB struct {
		Host     string
		User     string
		Password string `required:"true"`
		Database string
	}

	SFTP struct {
		Host     string
		Port     int
		User     string
		Password string `required:"true"`
		Folder   string `required:"true"`
	}
}{}

// DBConf (Public) -
type DBConf struct {
	Host     string
	User     string
	Password string
	Database string
}

func main() {

	if _, err := os.Stat("./config.yml"); err != nil {
		panic("no configuration file found")
	}
	configor.Load(&config, "config.yml")
	/*
		conf := importdb.SFTPConf{
			Address:  fmt.Sprintf("%v:%v", config.SFTP.Host, config.SFTP.Port),
			User:     config.SFTP.User,
			Password: config.SFTP.Password,
			Folder:   config.SFTP.Folder,
		}
		importdb.AppealsExtract(conf) /*/

	// Loading Database configurations
	dbconf := DBConf{
		Host:     config.DB.Host,
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Database,
	}
	// Step1: Create the Database instance and takes the boolean flag to either populate it with data or not
	db := seedDB(dbconf, true) // True does the TL in ETL
	// Defer closing the DB object until after all the commands in Main run below
	defer db.Close()
	// Optional Step: Insert Foreign keys into database
	// insertForeignKeys(db) // Add foreign keys for better performance???
	// Step2: Create appeal object that will have data injected into it.
	//a := importdb.Appeal{}
	appeals := importdb.AppealsTransform()
	for _, appeal := range appeals {
		db.Debug().Save(&appeal)
	}
	a := []importdb.Appeal{}
	//db.Debug().Preload("Appellant").Preload("Appellant.Rep").Preload("Appellant.Rep.Address").
	//	Preload("Appellant.Address").Preload("Reason").Preload("IDR").Preload("IDR.Issue").
	//	Preload("Action").Preload("Action.Hearing").Preload("Action.Accessibility").Find(&a, "ID = 71")

	// Step3: Run query to collect data and insert into the Appeal object created above.
	//db.Debug().Preload("Outreach").Find(&a, "ID = 71")
	// Step4: Either return the data or print it out in some meaningful way
	db.Order("ID desc").Limit("10").Preload("Appellant").Find(&a)
	for _, app := range a {
		fmt.Println(app.Appellant.LastName)
	}
	//*/
}

func seedDB(c DBConf, flag bool) *gorm.DB {
	db, err := gorm.Open("mssql", fmt.Sprintf("sqlserver://%v:%v@%v:1433?database=%v", c.User, c.Password, c.Host, c.Database))

	if err != nil {
		panic(err.Error())
	}

	if flag {
		begin := time.Now().UnixNano()
		var appealsInserted int
		db.DropTableIfExists(&importdb.Appeal{}, &importdb.Person{}, &importdb.Rep{},
			&importdb.Address{}, &importdb.Reason{}, &importdb.IDR{}, &importdb.Issue{},
			&importdb.Outreach{}, &importdb.Action{}, &importdb.Hearing{}, &importdb.Accessibility{},
			&importdb.Session{})
		db.AutoMigrate(&importdb.Appeal{}, &importdb.Person{}, &importdb.Rep{},
			&importdb.Address{}, &importdb.Reason{}, &importdb.IDR{}, &importdb.Issue{},
			&importdb.Outreach{}, &importdb.Action{}, &importdb.Hearing{}, &importdb.Accessibility{},
			&importdb.Session{})

		appeals := importdb.LoadData()
		db.Exec("SET IDENTITY_INSERT [dbo].[appeals] ON")
		for _, appeal := range appeals {
			db.Save(&appeal)
			appealsInserted++
		}
		db.Exec("SET IDENTITY_INSERT [dbo].[appeals] OFF")
		outreaches := importdb.LoadOutreach()
		db.Exec("SET IDENTITY_INSERT [dbo].[outreaches] ON")
		for _, outreach := range outreaches {
			db.Save(&outreach)
		}
		db.Exec("SET IDENTITY_INSERT [dbo].[outreaches] OFF")
		end := time.Now().UnixNano()
		duration := end - begin
		fmt.Printf("Time to import: %v:%v\n", duration/int64(time.Minute), (duration/int64(time.Second))%60)
		fmt.Printf("Appeals inserted: %v\n", appealsInserted)
	}

	return db
}

func insertForeignKeys(db *gorm.DB) {
	// Foreign key's, Cannot delete database without removing these constraints first or Dropping database
	// Think there might be a bug somewhere here.
	db.Debug().Model(&importdb.Person{}).AddForeignKey("appeal_id", "appeals(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.Rep{}).AddForeignKey("person_id", "people(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.Reason{}).AddForeignKey("appeal_id", "people(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.IDR{}).AddForeignKey("appeal_id", "appeals(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.Issue{}).AddForeignKey("idr_id", "idr(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.Outreach{}).AddForeignKey("appeal_id", "appeals(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.Action{}).AddForeignKey("appeal_id", "appeals(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.Hearing{}).AddForeignKey("action_id", "actions(ID)", "CASCADE", "CASCADE")
	db.Debug().Model(&importdb.Accessibility{}).AddForeignKey("action_id", "actions(ID)", "CASCADE", "CASCADE")
}
