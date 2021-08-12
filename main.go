package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"net/http"
)

type Driver struct {
	gorm.Model
	Name string
	License string
	Cars []Car
}

type Car struct {
	gorm.Model
	Year int
	Make string
	ModelName string
	DriverID int
}

var db *gorm.DB

var err error

var (
	drivers = []Driver{
		{Name: "Jimmy Johnson", License: "ABC123"},

		{Name: "Howard Hills",  License: "XYZ789"},

		{Name: "Craig Colbin",  License: "DEF333"}, 
	}

	cars =[]Car{
		{Year: 2000, Make: "Toyota", ModelName: "Tundra", DriverID: 1},

		{Year: 2001, Make: "Honda",  ModelName: "Accord", DriverID: 1},

		{Year: 2002, Make: "Nissan", ModelName: "Sentra", DriverID: 2},

		{Year: 2003, Make: "Ford",   ModelName: "F-150",  DriverID: 3},
	}
)

func main(){
	router := mux.NewRouter()

	dsn := "host=localhost user=postgres password=gans356pb786 dbname=dbtest sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn),&gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true} })
	if err != nil{
		panic("failed to connect database")
	}

	db.AutoMigrate(&Driver{})
	db.AutoMigrate(&Car{})

	for index := range drivers{
		db.Create(&drivers[index])
	}

	router.HandleFunc("/cars", GetCars).Methods("GET")
	router.HandleFunc("/cars/{id}", GetCar).Methods("GET")
	router.HandleFunc("/drivers/{id}", GetDriver).Methods("GET")
	router.HandleFunc("/cars/{id}", DeleteCar).Methods("DELETE")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080",handler))
}


func GetCars(w http.ResponseWriter, r *http.Request){
	var cars []Car

	db.Find(&cars)
	json.NewEncoder(w).Encode(&cars)
}

func GetCar(w http.ResponseWriter, r*http.Request){
	params := mux.Vars(r)

	var car Car

	db.First(&car, params["id"])
	json.NewEncoder(w).Encode(&car)
}

func GetDriver(w http.ResponseWriter, r *http.Request){
	params:= mux.Vars(r)
	var driver Driver

	db.First(&driver, params["id"])
	db.Model(&driver).Preload("Cars")

	json.NewEncoder(w).Encode(&driver)

}

func DeleteCar(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	var car Car

	db.First(&car, params["id"])
	db.Delete(&car)

	var cars []Car

	db.Find(&cars)
	json.NewEncoder(w).Encode(&cars)
}
