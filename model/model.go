package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Product struct {
	Id         int       `gorm:"primaryKey"; json : id`
	Name       string    `gorm:"unique; not null"; json : name`
	Price      float32   `gorm:"not null";json : price`
	Expiry     time.Time `gorm:"not null"; json : expiry`
	CategoryId int       `gorm:"not null"; json : categoryId`
}

type Datastore interface {
	Create(model *Product) (err error)
	Delete(model *Product, id int)
	Where(str string, id int) *gorm.DB
	Order(str string) *gorm.DB
	Find(model *Product) *gorm.DB
	Save(model *Product) (err error)
}
