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

func CreateTables() {
	// users, projects ve tasks tablolarını oluşturuyoruz
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		owner_id INTEGER REFERENCES users(id),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        description TEXT, 
        status TEXT DEFAULT 'pending',
        project_id INTEGER REFERENCES projects(id),
        assigned_to INTEGER REFERENCES users(id),
        due_date TIMESTAMP,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
    );`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Tablolar oluşturulurken hata: %v", err)
	}

	log.Println("📁 Database tables are ready!")
}
