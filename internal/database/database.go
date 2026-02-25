package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL sürücüsü
)

// DB nesnesini diğer paketlerden erişilebilir yapıyoruz
var DB *sql.DB

func Connect(dbURL string) {
	var err error
	// Veritabanına bağlantı açıyoruz
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Veritabanı sürücüsü hatası: %v", err)
	}

	// Bağlantıyı test ediyoruz (Ping)
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	log.Println("✅ Successfully connected to database!")
}
