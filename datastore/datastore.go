package datastore

import (
	//"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"rest/model"
)

type ProductDataStore struct {
	db *gorm.DB
}

func NewProductDataStore(db *gorm.DB) ProductDataStore {
	return ProductDataStore{
		db: db,
	}
}

func (pd ProductDataStore) Create(model *model.Product) (err error) {
	return pd.db.Create(model).Error
}

func (pd ProductDataStore) Delete(model *model.Product, id int) {
	pd.db.Delete(model, id)
}

func (pd ProductDataStore) Where(str string, id int) *gorm.DB {
	return pd.db.Where(str, id)
}
func (pd ProductDataStore) Order(str string) *gorm.DB {
	return pd.db.Order(str)
}
func (pd ProductDataStore) Find(model *model.Product) *gorm.DB {
	return pd.db.Find(model)
}
func (pd ProductDataStore) Save(model *model.Product) (err error) {
	return pd.db.Save(model).Error
}
