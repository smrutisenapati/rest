package rest

import(
	"github.com/jinzhu/gorm"
	"rest/model"
	"strconv"

	//"rest/datastore"
	"strings"
	"encoding/json"
	"io/ioutil"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Controller struct {
    datastore model.Datastore
}

func NewController(datastore model.Datastore) Controller{
    return Controller {
        datastore: datastore,
    }
}

func (ctrl Controller) CreateProd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()*/
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	data := &model.Product{}
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
		err := ctrl.datastore.Create(data)
		fmt.Printf("%T,%v",err,err)
		if err != nil { // to check if create causes an error
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
func (ctrl Controller) DeleteProd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()*/
	id := mux.Vars(r)["id"]
	data := &model.Product{}
	ctrl.datastore.Delete(data, id)
	msg["msg"]="deleted successfully"
	json.NewEncoder(w).Encode(msg)
}

func (ctrl Controller) ListProd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") // to send json response
	msg:=make(map[string]string)
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()*/
	var prod []model.Product
	params :=r.URL.Query()
	id, ok := params["categoryId"]
	var db *gorm.DB
	if ok{
		id, err := strconv.Atoi(id[0])
		if err != nil {

		}
		db = ctrl.datastore.Where("category_id = ?", id)
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
	db = db.Find(prod)
	if len(prod) == 0 {
		msg["error"]="invalid category"
		json.NewEncoder(w).Encode(msg)
	}else{
		json.NewEncoder(w).Encode(prod)
	}
}

func (ctrl Controller) UpdateProd(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	msg:=make(map[string]string)
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",host, port, user, password, dbname)
	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()*/
	id := mux.Vars(r)["id"]
	data := &model.Product{}
	db=ctrl.datastore.Where("id = ?", id).Find(data)
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
