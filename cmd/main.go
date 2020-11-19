package main

import (
	"rest/api"
	//"database/sql"
	//"strings"
	//"encoding/json"
	//"io/ioutil"
	"fmt"
	//"time"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // switch dialects to change between dbs
)

/*type Product struct{
	Id int `gorm:"primaryKey"; json : id`
	Name string  `gorm:"unique; not null"; json : name`
	Price float32 `gorm:"not null";json : price`
	Expiry time.Time `gorm:"not null"; json : expiry`
	CategoryId int `gorm:"not null"; json : categoryId`
}*/

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "go_inventory"
)

/*func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := &Product{}
	err = json.Unmarshal(jsn,data)
	if err != nil {
		panic(err)
	}
	if data.Name==""{
		msg["error"]="name is missing"
		json.NewEncoder(w).Encode(msg)
	}else if data.Price <= 0{
		msg["error"]="price is either missing or invalid"
		json.NewEncoder(w).Encode(msg)
	}else if data.CategoryId == 0{
		msg["error"]="category is missing"
		json.NewEncoder(w).Encode(msg)
	}else{
		err = db.Create(data).Error
		fmt.Printf("%T,%v",err,err)
		if err != nil{ // to check if create causes an error
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint"){ // to check if create causes an integrity error
				msg["error"]="name already exists"
				json.NewEncoder(w).Encode(msg)
			}
		}else{
			msg["error"]="created successfully"
			json.NewEncoder(w).Encode(msg)
		}
	}
}

func delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	id := mux.Vars(r)["id"]
	data := &Product{}
	db.Delete(data, id)
	msg["msg"]="deleted successfully"
	json.NewEncoder(w).Encode(msg)
}
func list(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var prod []Product
	params :=r.URL.Query()
	id, ok := params["categoryId"]
	if ok{
		db = db.Where("category_id = ?", id[0])
	}
	sort, ok := params["sort"]
	if ok{
		order, ok := params["order"]
		if ok{
			if sort[0] == "price"{
				if order[0] == "desc"{
					db = db.Order("price desc")
				}else{
					db = db.Order("price")
				}
			}else{
				if order[0] == "desc"{
					db = db.Order("expiry desc")
				}else{
					db = db.Order("expiry")
				}
			}
		}
	}
	db = db.Find(&prod)
	if len(prod) == 0 {
		msg["error"]="invalid category"
		json.NewEncoder(w).Encode(msg)
	}else{
		json.NewEncoder(w).Encode(prod)
	}
}

func update(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	msg:=make(map[string]string)
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	id := mux.Vars(r)["id"]
	data := &Product{}
	db=db.Where("id = ?", id).Find(data)
	if db.Error != nil{
		msg["error"]="product is not available"
		json.NewEncoder(w).Encode(msg)	
	}else{
		jsn, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(jsn,data)
		if err != nil {
			panic(err)
		}
		if data.Price < 0{
			msg["error"]="price is invalid"
			json.NewEncoder(w).Encode(msg)
		}else{
			err = db.Save(data).Error
			if err != nil{ // to check if create causes an error
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint"){ // to check if create causes an integrity error
					msg["error"]="name already exists"
					json.NewEncoder(w).Encode(msg)
				}
			}else{
				msg["error"]="updated successfully"
				json.NewEncoder(w).Encode(msg)
			}
		}
	}
}

func OpenConnection() *gorm.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	database :=db.DB()
	err = database.Ping()

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&Product{}) // creates a table based on the constraints in postgres if not created yet.

	return db
}

func handleRequests(){
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/delete/{id}",delete).Methods("DELETE")
	myRouter.HandleFunc("/get",list).Methods("GET")
	myRouter.HandleFunc("/create",create).Methods("POST")
	myRouter.HandleFunc("/update/{id}",update).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080",myRouter))
}*/

func main(){
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	datastore := NewProductDataStore(db)
	ctrl := NewController(datastore)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/delete/{id}",ctrl.DeleteProd).Methods("DELETE")
	myRouter.HandleFunc("/get",ctrl.ListProd).Methods("GET")
	myRouter.HandleFunc("/create",ctrl.CreateProd).Methods("POST")
	myRouter.HandleFunc("/update/{id}",ctrl.UpdateProd).Methods("PUT")
	log.Fatal(http.ListenAndServe(":8080",myRouter))
}
