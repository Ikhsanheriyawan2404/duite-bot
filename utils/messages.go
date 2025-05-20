package utils

const WelcomeText = `👋 Hai, selamat datang di Duite Bot  
Siap bantu kamu catat pengeluaran harian dengan cepat dan tanpa ribet.

Cukup kirim pesan seperti ini:
➡️ Makan siang 25k atau beli motor 20 mei 2025

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

const PromptDefault = `Analisa deskripsi transaksi berikut: "%s

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
}`
