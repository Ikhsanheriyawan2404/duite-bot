# 📘 REST API Documentation

## 🖥️ Base URL

```
https://duite-bot.brogrammer.id
```

---

## ✅ Health Check

### `GET /health`

**Deskripsi**: Memeriksa apakah server API hidup.

**Response**

```json
"ok"
```

---

## 👤 Users

### `POST /users/register`

**Deskripsi**: Mendaftarkan pengguna baru.

**Request Body**

```json
{
  "chatId": "123456789",
  "name": "John Doe"
}
```

**Successful Response — 201 Created**

```json
{
  "id": 2,
  "uuid": "fe9b06dd-b689-47f8-9840-aa2afae397e4",
  "chat_id": 0,
  "name": "IkhsanJH",
  "is_paid": null,
  "created_at": "2025-05-23T09:49:29.381641626Z",
  "updated_at": "2025-05-23T09:49:29.381641716Z",
  "Transactions": null
}
```

**400 Bad Request — User already exists**

```json
{
  "error": "user sudah terdaftar"
}
```

**Skenario**:

* ✅ Register user baru → sukses.
* ❌ Register dengan chatId yang sudah ada → error `400 Bad Request` atau pesan duplikat.

---

### `GET /users/:chatId/exists`

**Deskripsi**: Memeriksa apakah user dengan `chatId` tertentu sudah terdaftar.

**200 Response (jika ada)**

```json
{ "exists": true }
```

**404 Response (jika tidak ada)**

```json
{ "exists": false }
```

---

### `GET /users/:chatId/transactions/daily`

**Deskripsi**: Mengambil ringkasan transaksi harian user berdasarkan `chatId`.

**Response**

```json
[
  {
    "id": 2,
    "transaction_type": "INCOME",
    "amount": 1000000,
    "category": "gaji-bulanan",
    "description": "",
    "transaction_date": "2025-05-23T04:25:50.109Z",
    "chat_id": 6282117088123,
    "original_text": "masuk nih 1jt gaji bulanan",
    "created_at": "0001-01-01T00:00:00Z",
    "updated_at": "0001-01-01T00:00:00Z"
  },
]
```

**Skenario**:

* ✅ Ada transaksi hari ini → tampilkan ringkasan.
* ❌ Belum ada transaksi → `[]` (array kosong).

---

### `GET /users/:chatId/transactions/monthly`

**Deskripsi**: Mengambil ringkasan transaksi bulanan user.

**Response**

```json
{
  "month": "2025-07",
  "total": 1500000,
  "transactions": [ ... ]
}
```

---

### `DELETE /users/:chatId/transactions/:transactionId`

**Deskripsi**: Menghapus satu transaksi berdasarkan `chatId` dan `transactionId`.

**Response**

```json
{
  "message": "Transaction deleted"
}
```

**Skenario**:

* ✅ Transaksi ditemukan dan dihapus → sukses.
* ❌ Transaksi tidak ditemukan → `404 Not Found`.

---

### `POST /users/:chatId/transactions/ai-classify`

**Deskripsi**: Mengklasifikasikan transaksi berdasarkan deskripsi menggunakan AI.

**Request Body**

```json
{
  "prompt": "Analisa deskripsi transaksi berikut: beli nasi goreng 15k kemaren

Ekstrak:
1. type INCOME atau EXPENSE
2. amount angka rupiah
3. category transaksi ringkas (Indonesia)
4. date tanggal (format YYYY-MM-DD), atau null

Balas hanya JSON:
{
	"type": "INCOME|EXPENSE",
	"amount": number,
	"category": string
	"date": "YYYY-MM-DD" | null
}"
}
```

**Response**

```json
{
  "category": "Makanan",
  "confidence": 0.95
}
```

**Skenario**:

* ✅ Input valid → balikan kategori dengan confidence.
* ❌ Deskripsi kosong → `400 Bad Request`.

---

## 💰 Transactions

### `POST /transactions/`

**Deskripsi**: Menambahkan transaksi baru.

**Request Body**

```json
{
  "chatId": "123456789",
  "description": "Beli kopi",
  "amount": 15000,
  "category": "Minuman"
}
```

**Response**

```json
{
  "message": "Transaction created",
  "transaction": {
    "id": 101,
    "chatId": "123456789",
    "description": "Beli kopi",
    "amount": 15000,
    "category": "Minuman"
  }
}
```

**Skenario**:

* ✅ Input valid → transaksi berhasil dibuat.
* ❌ Tidak ada `chatId` → error.
* ❌ Invalid `amount` (misal 0 atau negatif) → error validasi.

---

## 🧪 Testing & Headers

* Semua endpoint menerima dan merespons dalam format JSON.
* Gunakan header berikut saat testing:

```
Content-Type: application/json
```

---

Jika kamu ingin dokumentasi ini dalam format `README.md` atau di-render dalam Swagger / Postman, saya juga bisa bantu!
