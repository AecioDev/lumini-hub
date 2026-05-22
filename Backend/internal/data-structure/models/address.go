package models

import "gorm.io/gorm"

type Country struct {
	gorm.Model

	Name      string `gorm:"not null;unique" json:"name"`
	BacenCode string `gorm:"not null;unique" json:"bacen_code"` // Código Bacen do país
	PhoneCode string `gorm:"not null" json:"phone_code"`        // Código telefônico internacional
}

// TableName especifica o nome da tabela
func (Country) TableName() string {
	return "country"
}

type State struct {
	gorm.Model

	Name      string  `gorm:"not null;unique" json:"name"`
	UF        string  `gorm:"size:2;not null;unique" json:"uf"`    // Sigla do Estado (ex: SP, RJ)
	IBGECode  string  `gorm:"not null;unique" json:"ibge_code"`    // Código IBGE da UF
	CountryID uint    `gorm:"not null" json:"country_id"`          // Relacionamento com Country
	Country   Country `gorm:"foreignKey:CountryID" json:"country"` // Associação com o país
}

// TableName especifica o nome da tabela
func (State) TableName() string {
	return "state"
}

type City struct {
	gorm.Model

	Name     string `gorm:"not null" json:"name"`
	IBGECode string `gorm:"not null;unique" json:"ibge_code"` // Código IBGE da cidade
	StateID  uint   `gorm:"not null" json:"state_id"`         // Relacionamento com State
	State    State  `gorm:"foreignKey:StateID" json:"state"`  // Associação com o estado
}

// TableName especifica o nome da tabela
func (City) TableName() string {
	return "city"
}

type Address struct {
	gorm.Model

	Street       string `gorm:"not null" json:"street"`
	Number       string `json:"number"`
	Neighborhood string `gorm:"not null" json:"neighborhood"`
	ZipCode      string `gorm:"not null;size:8" json:"zip_code"` // CEP com 8 dígitos
	CityID       uint   `gorm:"not null" json:"city_id"`         // Relacionamento com City
	City         City   `gorm:"foreignKey:CityID" json:"city"`   // Associação com a cidade

	CustomerID *uint     `gorm:"index" json:"customer_id,omitempty"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"-"`

	SupplierID *uint     `gorm:"index" json:"supplier_id,omitempty"`
	Supplier   *Supplier `gorm:"foreignKey:SupplierID" json:"-"`
}

// TableName especifica o nome da tabela
func (Address) TableName() string {
	return "address"
}
