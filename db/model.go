package db

import (
	"gorm.io/gorm"
	"time"
)

// region Base Model Define

// BaseModel Define
type BaseModel struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

//
//// GetID Get ID
//func (m *BaseModel) GetID() int64 {
//	return m.ID
//}
//
//// SetID Set ID
//func (m *BaseModel) SetID(id int64) {
//	m.ID = id
//}
//func (m *BaseModel) IsNil() bool {
//	if m == nil || m.ID == 0 {
//		return true
//	}
//	return false
//}

// endregion Base Model Define
