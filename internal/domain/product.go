package domain

type Product struct {
	Id        int32   `gorm:"column:id;primaryKey"`
	Name      string  `gorm:"column:name"`
	Code      string  `gorm:"column:code"`
	Price     float64 `gorm:"column:price"`
	CreatedBy int     `gorm:"column:created_by"`
	UpdatedBy int32   `gorm:"column:updated_by"`
}

func (Product) TableName() string {
	return "product"
}
