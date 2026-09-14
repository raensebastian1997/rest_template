# REST API Golang

API menggunakan arsitektur berlapis:

```text
HTTP Controller -> Service Interface -> Service Implementation
                -> Repository Interface -> GORM Repository -> MySQL
```

Controller hanya menangani HTTP. Validasi dan aturan bisnis berada di `internal/service`, sedangkan seluruh query database berada di `internal/repository`.

## Menjalankan aplikasi

Salin `.env.example` menjadi `.env`, lalu ganti `JWT_SECRET` dan `ADMIN_PASSWORD` dengan nilai yang aman. Setelah itu jalankan hot reload:

```powershell
air
```

Role `admin` dan `user` dibuat otomatis. Akun admin awal hanya dibuat jika `ADMIN_PASSWORD` diisi. Registrasi publik selalu menghasilkan akun dengan role `user`.

## Authentication

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
```

Kirim access token pada route terproteksi:

```text
Authorization: Bearer <access_token>
```

Role `user` dapat membaca product dan category. Role `admin` dapat mengelola product, category, user, dan melihat role.

## Query Index

Semua endpoint `Index` menerima:

```text
page=1
per_page=10
search=keyboard
sort_by=price
order=desc
category_id=1   # khusus product
```

Contoh:

```text
GET /api/v1/product/?page=1&per_page=20&search=keyboard&sort_by=price&order=desc&category_id=1
```

Response berisi `data` dan metadata `page`, `per_page`, `total`, `total_pages`, `sort_by`, dan `order`. Data product memuat object relasi `category` melalui GORM preload.
