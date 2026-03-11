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

## 🏗️ Architecture & Folder Structure
Bu projede, **Separation of Concerns** (Sorumlulukların Ayrıştırılması) prensibi doğrultusunda modüler bir yapı tercih edilmiştir:

```text
├── cmd/
│   └── main.go              # Uygulama giriş noktası
├── internal/
│   ├── handlers/            # HTTP isteklerini (Request/Response) yöneten katman
│   ├── middleware/          # JWT Yetkilendirme ve CORS kontrolleri
│   ├── models/              # Veritabanı şemaları ve veri modelleri (Structs)
│   └── database/            # PostgreSQL bağlantı ve konfigürasyonu
├── uploads/                 # Kullanıcılar tarafından yüklenen fiziksel dosyalar
├── .env                     # Çevresel değişkenler (Gizli bilgiler)
└── go.mod                   # Bağımlılık yönetimi
```
## 🚀 Kurulum
1. `go mod download`
2. `.env` dosyasını yapılandırın.
3. `go run main.go`
