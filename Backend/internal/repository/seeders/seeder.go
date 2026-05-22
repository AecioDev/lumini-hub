package seeders

import "gorm.io/gorm"

func RunAll(db *gorm.DB) {
	SeedUserAdm(db)

	SeedCountries(db)
	SeedStates(db)
	SeedCities(db)

	SeedRolesPermissions(db)
	//SeedMeasurementUnit(db)
	//SeedPaymentMethod(db)
	//SeedProductCategory(db)
}
