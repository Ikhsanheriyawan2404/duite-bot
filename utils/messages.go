package utils

const WelcomeText = `👋 Hai, selamat datang di Duite Bot  
Siap bantu kamu catat pengeluaran harian dengan cepat dan tanpa ribet.

Cukup kirim pesan seperti ini:
➡️ Makan siang 25k

Bot akan langsung menyimpan transaksi kamu secara otomatis! 🔥`

const HelpText = `📘 *Bantuan Finance Bot*

Berikut adalah command yang tersedia:

/start - Mulai untuk membuka menu
/harian - Melihat laporan transaksi hari ini
/bulanan - Melihat laporan transaksi bulan ini
/hapus - Menghapus data transaksi dengan cara /hapus (ID transaksi) -> berupa angka dan tanpa tanda pagar
/daftar - Daftarkan akun untuk lihat dashboard
/dashboard - Lihat Dashboard
/bantuan - Bantuan informasi
`

const PromptDefault = `Analisis deskripsi transaksi berikut:
Deskripsi: "%s"

Ekstrak informasi:
1. Tipe transaksi (INCOME/EXPENSE)
2. Nilai transaksi
3. Kategori transaksi (bahasa indonesia)

Hanya respon dengan format JSON:
{
	"type": "INCOME|EXPENSE",
	"amount": number,
	"category": "string"
}`
