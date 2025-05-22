Berikut adalah rangkuman fitur yang sudah tersedia dan rencana pengembangan untuk **bot Telegram catatan keuangan Duite Bot**

---

## ✅ **Fitur yang Tersedia Saat Ini**

### 📌 Bot Telegram Commands

| Command     | Deskripsi                                                           |
| ----------- | ------------------------------------------------------------------- |
| `start`    | Menampilkan informasi penggunaan bot                                |
| `harian`    | Melihat daftar transaksi hari ini                                   |
| `bulanan`   | Melihat transaksi bulan ini                                         |
| `hapus`     | Mengelola transaksi: hapus transaksi tertentu                       |
| `daftar`    | Registrasi pengguna                                                 |
| `dashboard` | Menampilkan link atau informasi untuk mengakses dashboard (website) |

---

### 🔗 **Integrasi dengan LLM**

* Menggunakan model:

  * `ChatGPT-4.1 nano`
  * `DeepSeek-Chat`
* Untuk:

  * Membantu parsing input pengguna
  * Klasifikasi jenis transaksi dari input bebas
  * Ekstraksi informasi jumlah, jenis (masuk/keluar), dan kategori

---

### 📊 **Dashboard Website**

Fitur visualisasi dan tampilan informasi keuangan user:

* **Card Ringkasan:**

  * Total pemasukan
  * Total pengeluaran

* **Line Chart:**

  * Visualisasi tren transaksi (masuk/keluar) sepanjang tahun

* **Pie Chart:**

  * Distribusi pengeluaran berdasarkan 5 kategori teratas + "Others"

* **Tabel Transaksi:**

  * Daftar transaksi terakhir atau semua transaksi lengkap

---

## 🛠️ **Rencana Update Fitur Selanjutnya**

### 🎯 Filter Interaktif di Dashboard

* Filter transaksi berdasarkan:

  * **Rentang tanggal** (start - end)
  * **Kategori transaksi**

### 💡 Rencana Pengembangan Tambahan (Saran)

Berikut beberapa ide fitur tambahan yang bisa kamu pertimbangkan:

#### 🔔 Notifikasi & Reminder

* Pengingat harian/bulanan untuk mencatat pengeluaran
* Notifikasi jika ada pengeluaran besar (threshold alert)

#### 🧠 Analisis Keuangan Pintar (LLM-enhanced)

* Rekomendasi penghematan
* Ringkasan keuangan bulanan secara otomatis
* Tanya jawab tentang keuangan pribadi via bot (contoh: “Berapa total pengeluaran saya minggu lalu?”)

#### 📅 Fitur Budgeting

* Set target bulanan untuk kategori tertentu (makanan, transportasi)
* Laporan jika melebihi batas budget

#### 📥 Import Transaksi

* Upload CSV atau file rekening untuk auto-import
* Sinkronisasi dengan e-wallet atau bank (jika memungkinkan ke depan)

#### 👥 Multi-User dalam Satu Keluarga/Tim

* Satu grup bisa kelola keuangan bersama (misalnya rumah tangga)

#### 📈 Insight Khusus

* Transaksi terbanyak terjadi pada hari apa?
* Kategori langganan (subscription) yang aktif dan besar

---
