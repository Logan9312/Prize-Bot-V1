package database

import (
	"fmt"
	"reflect"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var sourceDB, targetDB *gorm.DB

func NewDBConnect(password, host string) {
	fmt.Println("Connecting to Database for migration")

	targetDB = NewDB(password, host)
	sourceDB = DB

	err := targetDB.AutoMigrate(AuctionSetup{}, Auction{}, AuctionQueue{}, GiveawaySetup{}, Giveaway{}, ClaimSetup{}, CurrencySetup{}, Claim{}, DevSetup{}, UserProfile{}, ShopSetup{}, WhiteLabels{})
	if err != nil {
		fmt.Println(err)
	}
}

func NewDB(password, host string) *gorm.DB {
	dbuser := "postgres"
	port := "5527"
	dbname := "railway"

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", host, port, dbuser, dbname, password)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func CopyTablesToNewDB() {
	// Define a list with all of your tables
	tables := []interface{}{
		&AuctionSetup{},
		&Auction{},
		&AuctionQueue{},
		&GiveawaySetup{},
		&Giveaway{},
		&ClaimSetup{},
		&CurrencySetup{},
		&Claim{},
		&DevSetup{},
		&UserProfile{},
		&ShopSetup{},
		&WhiteLabels{},
	}

	// Loop over all the tables
	for _, table := range tables {
		// Using reflect we can create a new slice to hold our data
		slicePtr := reflect.New(reflect.SliceOf(reflect.TypeOf(table).Elem()))
		slice := reflect.Indirect(slicePtr)

		// Find all records in the source database and put them in the slice
		if err := sourceDB.Find(slice.Addr().Interface()).Error; err != nil {
			fmt.Printf("Could not copy table %s: %s\n", reflect.TypeOf(table).Elem().Name(), err)
			continue
		}

		// Now create the table in the target database
		if err := targetDB.Migrator().CreateTable(table); err != nil {
			fmt.Printf("Could not create table %s: %s\n", reflect.TypeOf(table).Elem().Name(), err)
			continue
		}

		// Copy all records to the target database
		if err := targetDB.Create(slice.Interface()).Error; err != nil {
			fmt.Printf("Could not copy data to table %s: %s\n", reflect.TypeOf(table).Elem().Name(), err)
		}
	}

	DB = targetDB

	fmt.Println("Finished copying tables to new database")
}
