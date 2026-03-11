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
│   └── main.go              # Uygulamanın başlatıldığı ana dosya
├── internal/
│   ├── handlers/            # HTTP Request/Response mantığının yönetildiği katman
│   ├── middleware/          # JWT Auth, Logging ve CORS güvenlik katmanları
│   ├── models/              # Veritabanı tablolarının Go karşılığı olan yapılar (Structs)
│   ├── database/            # PostgreSQL bağlantı havuzu yönetimi
│   ├── config/              # .env ve çevresel değişkenlerin yüklendiği konfigürasyon katmanı
│   └── router/              # Uygulamanın tüm API uç noktalarının (Routes) tanımlandığı yer
├── uploads/                 # Sistem üzerinden yüklenen fiziksel dökümanlar
├── .env                     # Hassas veritabanı ve JWT anahtarı bilgileri
```
## 🚀 Kurulum
1. `go mod download`
2. `.env` dosyasını yapılandırın.
3. `go run main.go`
