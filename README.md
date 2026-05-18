# Gothub ERP - Core Authentication & Course Management API 🚀

[![Go Version](https://img.shields.io/badge/Go-1.26-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat)](https://github.com/gin-gonic/gin)
[![GORM](https://img.shields.io/badge/GORM-PostgreSQL-336791?style=flat&logo=postgresql)](https://gorm.io/)
[![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Gothub ERP adalah backend API modular berskala industri yang dirancang untuk sistem *Enterprise Resource Planning*. Proyek ini mengimplementasikan **Clean Architecture** (Handler, Service, Repository, Model) untuk memastikan kode yang mudah dikelola, diuji, dan dikembangkan.

Proyek ini mencakup sistem autentikasi yang kuat dengan **Role-Based Access Control (RBAC)** dan manajemen kursus/pengguna yang aman.

## 🌟 Fitur Utama

- **Clean Architecture**: Pemisahan tanggung jawab yang jelas antar lapisan.
- **Autentikasi Lanjut**: Implementasi JWT dengan *Access Token* dan *Refresh Token*.
- **Role-Based Access Control (RBAC)**:
  - **Admin**: Akses penuh ke manajemen user dan sistem.
  - **Guru**: Manajemen kursus dan nilai.
  - **Siswa**: Akses terbatas ke kursus dan jadwal.
- **Keamanan Ketat**:
  - Password Hashing dengan `bcrypt`.
  - Rate Limiting untuk mencegah brute force.
  - Security Headers & CORS terkonfigurasi.
  - Global Error Handling.
- **Database Otomatis**: Migrasi otomatis (Auto-migrate) dan Seeder untuk Super Admin.
- **Siap Docker**: Dukungan penuh containerization dengan Docker & Docker Compose.

## 🛠️ Tech Stack

- **Backend**: Go (Golang) v1.26
- **Web Framework**: [Gin Gonic](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: PostgreSQL
- **Security**: JWT (golang-jwt), Bcrypt, Security Headers
- **Infrastruktur**: Docker, Docker Compose

## 📁 Struktur Proyek

```text
.
├── cmd/api/main.go              # Titik masuk aplikasi (Entry Point)
├── internal/
│   ├── database/                # Seeder & Konfigurasi DB
│   ├── handler/                 # Controller / HTTP Interface
│   ├── middleware/              # Auth, Error, Security, Rate Limit
│   ├── models/                  # Struct entitas & Database
│   ├── repository/              # Data Access Layer
│   ├── routes/                  # Definisi Endpoint API
│   ├── service/                 # Business Logic Layer (Otak aplikasi)
│   └── utils/                   # Helper (JWT, Response, Error)
├── .env.example                 # Template konfigurasi environment
├── Dockerfile                   # Konfigurasi Docker (Multi-stage build)
└── docker-compose.yml           # Orkestrasi API & Database
```

## 🚀 Memulai (Getting Started)

### Prasyarat
- Docker & Docker Compose terinstal di mesin Anda.

### Langkah Instalasi

1.  **Clone Repositori**:
    ```bash
    git clone https://github.com/azharf99/gothub-erp.git
    cd gothub-erp
    ```

2.  **Konfigurasi Environment**:
    Salin file `.env.example` menjadi `.env` dan sesuaikan nilainya:
    ```bash
    cp .env.example .env
    ```

3.  **Jalankan dengan Docker**:
    ```bash
    docker-compose up -d --build
    ```

4.  **Akses API**:
    Server akan berjalan di `http://localhost:8080`.

## 📌 Endpoint API Utama (v1)

### Publik
- `POST /api/v1/register` - Pendaftaran akun baru
- `POST /api/v1/login` - Login untuk mendapatkan token
- `POST /api/v1/refresh-token` - Memperbarui access token

### Terproteksi (Wajib Login)
- `GET /api/v1/profile` - Melihat profil saya
- `POST /api/v1/logout` - Logout sistem
- `GET /api/v1/courses` - Melihat semua kursus

### Manajemen Kursus (Guru & Admin)
- `POST /api/v1/courses` - Membuat kursus baru
- `PUT /api/v1/courses/:id` - Memperbarui kursus
- `DELETE /api/v1/courses/:id` - Menghapus kursus

### Manajemen User (Admin Saja)
- `GET /api/v1/users` - List semua user
- `POST /api/v1/users` - Membuat user baru
- `PUT /api/v1/users/:id` - Edit user
- `DELETE /api/v1/users/:id` - Hapus user

## 📄 Lisensi & Atribusi

Proyek ini dilisensikan di bawah **Apache License 2.0**.

Sesuai dengan ketentuan lisensi, siapa pun yang menggunakan, mendistribusikan, atau memodifikasi kode ini **WAJIB mencantumkan nama penulis asli**.

**Penulis Asli**:
**Azhar Faturohman Ahidin**
[GitHub @azharf99](https://github.com/azharf99)

Lihat file [LICENSE](LICENSE) dan [NOTICE](NOTICE) untuk detail lebih lanjut.

---
Dibuat dengan ❤️ oleh [Azhar Faturohman Ahidin](https://github.com/azharf99). Jika Anda menyukai proyek ini, silakan berikan ⭐!
