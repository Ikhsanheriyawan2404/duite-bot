package static

/** VERSI 0.5
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
*/

const PromptDefault = `
Current date is: "%s

Analyze the following transaction description:
"%s

Interpret any relative time expressions based on the current date above.

Extract and reply in JSON format only:
{
	"type": "INCOME|EXPENSE",
	"amount": number,
	"category": string,
	"date": "YYYY-MM-DD" or null
}

Use one of these categories:
- Makanan & Minuman
- Transportasi
- Belanja
- Tagihan & Utilitas
- Pendidikan
- Kesehatan
- Hiburan
- Donasi & Sosial
- Perawatan Diri
- Cicilan / Kredit
- Pajak & Retribusi
- Investasi & Tabungan
- Biaya Rumah Tangga
- Biaya Anak
- Hewan Peliharaan
- Lain-lain
- Gaji / Pendapatan Tetap
- Bisnis / Dagang
- Investasi
- Hadiah / Hibah
- Penjualan Aset
- Pengembalian Uang
- Lain-lain

Output must be in Bahasa Indonesia.
`