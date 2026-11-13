package model

// migration applies the schema for the selected database.
func migration(driver string) {
	db := DB
	if driver == "mysql" {
		db = db.Set("gorm:table_options", "charset=utf8mb4")
	}
	if err := db.AutoMigrate(&User{}).AutoMigrate(&Video{}).Error; err != nil {
		panic(err)
	}
}
