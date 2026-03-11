# 🚀 TaskFlow Backend - Go Rest API

Bu proje, yüksek performanslı ve ölçeklenebilir bir görev yönetim sisteminin backend altyapısıdır.

## 🛠️ Teknolojiler
- **Dil:** Go (Golang)
- **Veritabanı:** PostgreSQL
- **Yetkilendirme:** JWT (JSON Web Token)
- **Dosya Yönetimi:** UUID tabanlı güvenli depolama sistemi

## ✨ Öne Çıkan Özellikler
- **İlişkisel Dosya Yönetimi:** Görevlere (Tasks) bağlı çoklu dosya yükleme desteği.
- **Güvenli Mimari:** Middleware tabanlı yetkilendirme ve dosya tipi doğrulama.
- **Clean Architecture:** Bakımı kolay, modüler kod yapısı.
- 
# 🏗️ Architecture & Folder Structure
Bu projede, sorumlulukların ayrıştırılması (Separation of Concerns) prensibi doğrultusunda modüler bir yapı tercih edilmiştir:

lib/
 ├── cmd/             # Uygulama giriş noktası (main.go)
 ├── internal/        # Dışarıya kapalı, çekirdek mantık
 │    ├── handlers/   # HTTP isteklerini karşılayan katman
 │    ├── middleware/ # Auth & CORS kontrolleri
 │    ├── models/     # Veritabanı şemaları (Structs)
 │    └── database/   # DB bağlantı ve konfigürasyonu
 └── uploads/         # Fiziksel dosya depolama alanı
 
## 🚀 Kurulum
1. `go mod download`
2. `.env` dosyasını yapılandırın.
3. `go run main.go`
