package db

import (
	"gorm.io/gorm"
	"time"
)

// region Base Model Define

type ID interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | string
}

// IModel ModelT Interface
type IModel[T ID] interface {
	GetID() T
	IsNil() bool
}

// ModelT Define

type ModelT[T ID] struct {
	ID        T              `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// GetID Get ID
func (m ModelT[T]) GetID() T {
	return m.ID
}

func (m ModelT[T]) IsNil() bool {
	var t T
	if m.ID == t {
		return true
	}
	return false
}

type ModelI64 struct {
	ID        int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (m ModelI64) GetID() int64 {
	return m.ID
}

func (m ModelI64) IsNil() bool {
	return m.ID == 0
}

// endregion Base Model Define
