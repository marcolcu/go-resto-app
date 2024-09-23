package database

import (
    "fmt"
    "os"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    var err error

    // Ambil konfigurasi dari environment variables yang disediakan oleh Railway
    user := os.Getenv("MYSQLUSER")
    password := os.Getenv("MYSQLPASSWORD")
    host := os.Getenv("MYSQLHOST")
    port := os.Getenv("MYSQLPORT")
    dbname := os.Getenv("MYSQLDATABASE")

    // Format DSN menggunakan environment variables dari Railway
    DSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        user, password, host, port, dbname)

    // Membuka koneksi ke database
    DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})
    if err != nil {
        panic("Can't connect to database")
    }

    fmt.Println("Connected to database")
}
