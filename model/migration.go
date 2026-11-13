package model

// migration applies the schema for the selected database.
func migration(driver string) {
	db := DB
	if driver == "mysql" {
		db = db.Set("gorm:table_options", "charset=utf8mb4")
	}
	db.AutoMigrate(&User{}).AutoMigrate(&Video{})
}
