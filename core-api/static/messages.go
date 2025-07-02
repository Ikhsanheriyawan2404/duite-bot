package static

/** VERSI 0.5.0
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

/** VERSI 1.0.0 peningkatan di grounding category & pengetahuan waktu saat ini
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
*/

// Versi 1.5.0 peningkatan grounding category & dinamis category
const PromptDefault = `
Current date is: "%s

Analyze the following transaction description:
"%s

Interpret any relative time expressions based on the current date above.

Extract and reply in JSON format only:
{
	"type": "INCOME|EXPENSE",
	"amount": number,
	"category_id": number,
	"date": "YYYY-MM-DD" or null
}

Use one of these categories:

#1 Gaji & Pendapatan Tetap
#2 Pendapatan Lain / Usaha
#3 Investasi & Bunga
#4 Hadiah & Lain-lain
#5 Kebutuhan Harian
#6 Transportasi
#7 Tagihan & Cicilan
#8 Kesehatan
#9 Pendidikan
#10 Hiburan & Sosial
#11 Tabungan & Investasi
#12 Lain-lain Pengeluaran
#13 Donasi
#14 Belanja Online

Output must be in Bahasa Indonesia.
`