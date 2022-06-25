package main

import (
	"fmt"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xuri/excelize/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Strurktur DB
type (
	Data struct {
		gorm.Model
		Name  string
		Qty   int
		Price int
	}
	Datas []Data
)

// function for Connect to DB
func database() (*gorm.DB, error) {
	dsn := "root:password@tcp(127.0.0.1:3306)/testing_mysql?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return db, err
}

func main() {
	// Called function database() to Connect to database
	db, err := database()
	if err != nil {
		fmt.Println("Err : ", err)
	}

	// Migration to DB
	db.AutoMigrate(&Data{})

	// Open file Excel (Book1.xlsx)
	f, err := excelize.OpenFile("Book1.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Read Sheet1
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Goroutine
	go func(rows [][]string) {
		datas := Datas{}

		for i := 1; i < len(rows); i++ {
			qty, _ := strconv.Atoi(rows[i][2])
			price, _ := strconv.Atoi(rows[i][3])

			data := Data{}
			data.Name = rows[i][1]
			data.Qty = qty
			data.Price = price
			datas = append(datas, data)
		}

		db.Create(datas)
	}(rows)

	// message PASS
	fmt.Println("PASS!")

	time.Sleep(time.Second)

}
