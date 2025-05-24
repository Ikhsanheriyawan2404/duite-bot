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

### POST `/users/:chatId/transactions/ai-classify`

**Deskripsi**: Endpoint ini menerima deskripsi transaksi dalam parameter `prompt` dan menggunakan AI (LLM) untuk mengklasifikasikan transaksi menjadi tipe pemasukan (`INCOME`) atau pengeluaran (`EXPENSE`), beserta nominal, kategori, dan tanggal transaksi. Jika valid, transaksi disimpan ke database.

### Request

* **Path Parameter**

| Nama   | Tipe  | Keterangan          |
| ------ | ----- | ------------------- |
| chatId | int64 | ID chat/user tujuan |

* **Body JSON**

```json
{
  "prompt": "string"
}
```

* **prompt**: Deskripsi transaksi yang mengandung informasi nominal dan deskripsi lain.

---

### Response

#### 1. Success (201 Created)

```json
{
  "message": "Transaksi berhasil disimpan",
  "usage": {
    "prompt_tokens": 45,
    "completion_tokens": 20,
    "total_tokens": 65
  },
  "data": {
    "chatId": 123456,
    "originalText": "beli nasi goreng 15k kemaren",
    "transactionType": "EXPENSE",
    "amount": 15000,
    "category": "makanan",
    "transactionDate": "2025-05-24T00:00:00Z"
  }
}
```

---

#### 2. Error Response

| HTTP Status               | Kondisi                                                        | Response JSON                                                                                                             |
| ------------------------- | -------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------- |
| 400 Bad Request           | `prompt` kosong atau tidak dapat diparse                       | `{"error": "Prompt is required"}`                                                                                         |
| 400 Bad Request           | `prompt` tidak mengandung nominal (angka/k)                    | `{"error": "Prompt harus mengandung nominal (angka atau k)", "help": "Contoh: 'Makan siang 25k' atau 'gaji masuk 1000'"}` |
| 500 Internal Server Error | Gagal komunikasi dengan LLM (AI)                               | `{"error": "Gagal mengklasifikasi transaksi"}`                                                                            |
| 400 Bad Request           | Hasil klasifikasi tipe transaksi bukan `INCOME` atau `EXPENSE` | `{"error": "Tipe transaksi tidak valid (harus INCOME atau EXPENSE)"}`                                                     |
| 500 Internal Server Error | Gagal menyimpan transaksi ke database                          | `{"error": "Gagal menyimpan transaksi"}`                                                                                  |

---

### Skema Proses Logika

1. Terima `prompt` dari body JSON.
2. Validasi `prompt` tidak kosong.
3. Validasi `prompt` mengandung nominal (angka atau karakter `k`/`K`).
4. Kirim `prompt` ke fungsi klasifikasi AI (`hitChatGpt`).
5. Terima hasil klasifikasi dengan atribut:

   * `type` (`INCOME`/`EXPENSE`)
   * `amount` (float64)
   * `category` (string)
   * `date` (string dalam format RFC3339 atau kosong)
6. Validasi tipe transaksi (hanya `INCOME` atau `EXPENSE` diterima).
7. Parse tanggal transaksi jika ada, fallback ke `time.Now()` jika kosong atau parsing gagal.
8. Simpan transaksi ke database via `transactionService.CreateTransaction`.
9. Berikan response sukses dengan data transaksi dan usage AI.

---

### Contoh Request

```bash
curl -X POST "http://yourapi.com/users/123456/transactions/ai-classify" \
-H "Content-Type: application/json" \
-d '{"prompt":"Beli kopi 12k hari ini"}'
```

---

### Contoh Response Sukses

```json
{
  "message": "Transaksi berhasil disimpan",
  "usage": {
    "prompt_tokens": 33,
    "completion_tokens": 15,
    "total_tokens": 48
  },
  "data": {
    "chatId": 123456,
    "originalText": "Beli kopi 12k hari ini",
    "transactionType": "EXPENSE",
    "amount": 12000,
    "category": "minuman",
    "transactionDate": "2025-05-24T00:00:00Z"
  }
}
```

---

### Contoh Response Error karena prompt kosong

```json
{
  "error": "Prompt is required"
}
```

---

### Contoh Response Error karena nominal tidak ditemukan

```json
{
  "error": "Prompt harus mengandung nominal (angka atau k)",
  "help": "Contoh: 'Makan siang 25k' atau 'gaji masuk 1000'"
}
```

---

Kalau kamu pakai Swagger atau OpenAPI, saya bisa bantu generate format YAML/JSON-nya juga. Butuh?


**Skenario**:

* ✅ Input valid → balikan kategori dengan confidence.
* ❌ Deskripsi kosong → `400 Bad Request`.

---

## 🧪 Testing & Headers

* Semua endpoint menerima dan merespons dalam format JSON.
* Gunakan header berikut saat testing:

```
Content-Type: application/json
```

---

Jika kamu ingin dokumentasi ini dalam format `README.md` atau di-render dalam Swagger / Postman, saya juga bisa bantu!
