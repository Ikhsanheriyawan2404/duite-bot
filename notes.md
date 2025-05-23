## ✅ Daftar Task Proyek Duite Bot

### 🔐 Authentication & Access Control

* [✓] Handle user yang tidak terdaftar ketika mengakses dashboard

---

### 🎨 Frontend Improvements

* \[✓] Formatting tanggal (frontend)
* \[✓] Slug untuk kategori (frontend)
* \[✓] ENUM untuk tipe transaksi (contoh: INCOME → Masuk, EXPENSE → Keluar)
* \[✓] Penanganan kondisi data kosong di chart (line & pie chart)
* [ ] Filter transaksi update rendering: card, pie chart, line chart, dan tabel transaksi

---

### 🐞 Bug & Error Handling

* [ ] Cek dan cari bug (data, chart, dll)
* [ ] Penanganan error & logging yang rapi
* \[✓] Bug transaksi harian dan bulanan (filter data)
* [✓] Tangani kondisi error aneh di chart (backend dan frontend)

---

### ⚙️ Backend Improvements

* [ ] Optimasi query transaksi
* [✓] Refactor ke pola clean architecture
* \[✓] Validasi dan ENUM tipe transaksi (EXPENSE/INCOME)
* [ ] Validasi input chat Telegram user sebelum hit API LLM
* [ ] Buat filter transaksi (by bulan, kategori, rentang tanggal)

---

### 🚀 CI/CD & Testing

* [ ] Implementasi CI/CD (build, lint, test, deploy)
* [ ] Implementasi unit testing

---

### 📝 Dokumentasi & Perencanaan

* [✓] Catatan daftar fitur yang sudah ada
* [✓] Rencana fitur selanjutnya (next feature)

---

### 🔬 Research & Development

* [ ] Pelajari tokenization & cara kerja LLM (temperature, cache hit, dll)
* [ ] Perbandingan model: ChatGPT (rate limit) vs DeepSeek (no rate limit, kompetitif)
* [ ] Efisiensi token input bahasa Inggris vs Bahasa Indonesia, output tetap Bahasa Indonesia

---

## ✅ Fitur yang Tersedia Saat Ini

### 🤖 Bot Telegram Commands

| Command     | Deskripsi                                                           |
| ----------- | ------------------------------------------------------------------- |
| `start`     | Menampilkan informasi penggunaan bot                                |
| `harian`    | Melihat daftar transaksi hari ini                                   |
| `bulanan`   | Melihat transaksi bulan ini                                         |
| `hapus`     | Mengelola transaksi: hapus transaksi tertentu                       |
| `daftar`    | Registrasi pengguna                                                 |
| `dashboard` | Menampilkan link atau informasi untuk mengakses dashboard (website) |

---

### 🔗 Integrasi LLM

* Model yang digunakan:

  * `ChatGPT-4.1 nano`
  * `DeepSeek-Chat`
* Fungsi:

  * Parsing input pengguna
  * Klasifikasi jenis transaksi
  * Ekstraksi informasi jumlah, tipe transaksi, dan kategori

---

### 📊 Dashboard Website

* **Ringkasan:**

  * Total pemasukan
  * Total pengeluaran
* **Line Chart:**

  * Tren transaksi tahunan
* **Pie Chart:**

  * Distribusi pengeluaran per kategori (top 5 + Others)
* **Tabel Transaksi:**

  * Daftar transaksi terakhir atau semua transaksi

---

## 🛠️ Rencana Update Fitur Selanjutnya

### 🎯 Filter Interaktif di Dashboard

* Filter berdasarkan:

  * Rentang tanggal (start–end)
  * Kategori transaksi

---

## 💡 Rencana Fitur Tambahan

### 🔔 Notifikasi & Reminder

* Pengingat mencatat transaksi harian/bulanan
* Notifikasi pengeluaran besar (threshold alert)

### 🧠 Analisis Keuangan Pintar (LLM)

* Rekomendasi penghematan
* Ringkasan keuangan otomatis
* Tanya jawab keuangan via bot (contoh: “Berapa total pengeluaran minggu lalu?”)

### 📅 Budgeting

* Target bulanan per kategori (makan, transportasi)
* Laporan jika melebihi budget

### 📥 Import Transaksi

* Upload CSV untuk auto-import
* Sinkronisasi dengan e-wallet/bank (jika memungkinkan)

### 👥 Multi-User (Keluarga/Tim)

* Satu grup mengelola keuangan bersama

### 📈 Insight Khusus

* Hari paling banyak transaksi
* Deteksi subscription besar dan aktif