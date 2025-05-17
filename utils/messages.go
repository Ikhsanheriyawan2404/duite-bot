package utils

const HelpText = `📘 *Bantuan Finance Bot*
Berikut adalah command yang tersedia:

/start - Mulai care dengan cashflow
/harian - Melihat laporan transaksi hari ini
/bulanan - Melihat laporan transaksi bulan ini
/hapus - Menghapus data transaksi dengan cara /hapus (ID transaksi) -> berupa angka dan tanpa tanda pagar
/daftar - Daftarkan akun untuk lihat dashboard
/dashboard - Lihat Dashboard
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
