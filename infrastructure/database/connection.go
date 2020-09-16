package data

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	// registering database driver
	_ "github.com/lib/pq"
)

var (
	data *Data
	once sync.Once
)
// Data manages the connection to the database.
type Data struct {
	DB *sql.DB
}

// New returns a new instance of Data with the database connection ready.
func New() *Data {
	once.Do(initDB)
	return data
}

func NewTest() *Data {
	once.Do(initDBTest)
	return data
}

// initialize the data variable with the connection to the database.
func initDB() {

	db, err := GetConnection()
	if err != nil {
		fmt.Println("Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database")
	}

	err = MakeMigration(db)
	if err != nil {
		log.Fatal("This is the error:", err)
	}

	data = &Data{
		DB: db,
	}
}

func initDBTest() {

	db, err := GetConnectionTest()
	if err != nil {
		fmt.Println("Cannot connect to database test")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Println("We are connected to the database test")
	}

	err = MakeMigrationTest(db)
	if err != nil {
		log.Fatal("This is the error:", err)
	}

	data = &Data{
		DB: db,
	}
}

// Close closes the resources used by data.
func Close() error {
	if data == nil {
		return nil
	}

	return data.DB.Close()
}

func GetConnection() (*sql.DB, error) {
	DbHost := os.Getenv("DB_HOST")
	DbDriver := os.Getenv("DB_DRIVER")
	DbUser := os.Getenv("DB_USER")
	DbPassword := os.Getenv("DB_PASSWORD")
	DbName := os.Getenv("DB_NAME")
	DbPort := os.Getenv("DB_PORT")

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	return sql.Open(DbDriver, uri)
}

func GetConnectionTest() (*sql.DB, error) {
	DbHost := os.Getenv("TestDbHost")
	DbDriver := os.Getenv("TestDbDriver")
	DbUser := os.Getenv("TestDbUser")
	DbPassword := os.Getenv("TestDbPassword")
	DbName := os.Getenv("TestDbName")
	DbPort := os.Getenv("TestDbPort")

	uri := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	return sql.Open(DbDriver, uri)
}

// MakeMigration creates all the tables in the database
func MakeMigration(db *sql.DB) error {
	b, err := ioutil.ReadFile("./infrastructure/database/models.sql")
	if err != nil {
		return err
	}

	rows, err := db.Query(string(b))
	if err != nil {
		return err
	}

	return rows.Close()
}

// MakeMigrationTest creates all the tables in the database
func MakeMigrationTest(db *sql.DB) error {
	b, err := ioutil.ReadFile("../../infrastructure/database/models.sql")
	if err != nil {
		return err
	}

	rows, err := db.Query(string(b))
	if err != nil {
		return err
	}

	return rows.Close()
}
