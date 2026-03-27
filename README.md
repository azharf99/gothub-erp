# Gothub ERP - Core Authentication API 🚀

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-PostgreSQL-336791?style=flat&logo=postgresql)](https://gorm.io/)
[![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)](https://www.docker.com/)

Sebuah modul backend API berskala industri (Production-Grade) untuk sistem *Enterprise Resource Planning* (ERP). Proyek ini dibangun menggunakan arsitektur *Clean Architecture* untuk memastikan skalabilitas, kemudahan *testing*, dan pemisahan logika bisnis yang jelas.

Repositori ini juga merupakan bagian dari materi *masterclass* pengembangan web *Full Stack* dan *Cybersecurity* di kanal YouTube **Gothub**.

## ✨ Fitur Utama
- **Arsitektur Standar Industri**: Pemisahan lapisan *Handler, Usecase/Service, Repository, dan Model* (Clean Architecture).
- **Keamanan Lapis Ganda**: Autentikasi menggunakan standar JWT (*JSON Web Token*) dengan implementasi *Access Token* dan *Refresh Token*.
- **Proteksi Kata Sandi**: Hashing kredensial pengguna menggunakan `bcrypt`.
- **Database Terkelola**: Relasi dan migrasi otomatis menggunakan GORM dan PostgreSQL.
- **Containerized**: Siap dideploy ke VPS atau *cloud* mana pun menggunakan Docker & Docker Compose dalam satu perintah.

## 🛠️ Tech Stack
- **Bahasa**: Go (Golang)
- **Web Framework**: Gin HTTP Framework
- **ORM**: GORM
- **Database**: PostgreSQL
- **Keamanan**: JWT (golang-jwt/jwt/v5), Bcrypt
- **Infrastruktur**: Docker, Docker Compose (Alpine Linux & Multi-stage Build)

## 📁 Struktur Direktori
```text
.
├── cmd/api/main.go              # Entry point & Dependency Injection
├── internal/
│   ├── handler/                 # Layer HTTP, Request Validation & Response JSON
│   ├── middleware/              # Layer proteksi rute (JWT Auth)
│   ├── models/                  # Struct entitas DB dan Payload JSON
│   ├── repository/              # Layer interaksi langsung dengan PostgreSQL
│   ├── routes/                  # Pendaftaran Endpoint API
│   └── utils/                   # Fungsi bantuan (Generate JWT, dsb)
├── Dockerfile                   # Resep containerisasi (Multi-stage build)
├── docker-compose.yml           # Orkestrasi layanan API dan Database
└── .env.example                 # Contoh environment variables


🚀 Cara Menjalankan Aplikasi
Syarat Sistem
Pastikan kamu telah menginstal Docker dan Docker Compose di komputermu.

Langkah Instalasi
Clone repositori ini:

Bash
git clone [https://github.com/username-kamu/gothub-erp.git](https://github.com/username-kamu/gothub-erp.git)
cd gothub-erp
Siapkan Environment:
Duplikat file .env.example menjadi .env dan sesuaikan kredensial di dalamnya (Gunakan kata sandi yang kuat untuk JWT dan Database).

Jalankan dengan Docker Compose:

Bash
docker-compose up -d --build
Testing API:
Aplikasi akan berjalan di http://localhost:8080. Kamu bisa mengimpor file Gothub_ERP.postman_collection.json (jika kamu menyertakannya di repo) ke dalam aplikasi Postman untuk langsung menguji endpoint /register, /login, /refresh-token, dan /profile.

👨‍💻 Penulis
Azhar Faturohman Ahidin, S.Si.
Full Stack Web Developer | Creator of Gothub

Mari berjejaring! Jika kamu menemukan proyek ini bermanfaat, jangan ragu untuk memberikan ⭐ (Star) pada repositori ini.


---

### Langkah Terakhir (Push ke GitHub)

Agar proyekmu siap dikloning orang lain dengan aman, buat satu file lagi bernama `.env.example`. File ini adalah *template* kosong dari `.env` yang isinya tidak berisi rahasia aslimu, sekadar memberi tahu orang lain variabel apa saja yang dibutuhkan.

**Isi `.env.example`:**
```env
PORT=8080
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=masukkan_password_kamu_di_sini
DB_NAME=erp_db
DB_PORT=5432
JWT_SECRET=masukkan_kunci_rahasia_jwt_di_sini