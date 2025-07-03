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
//========================================== first prompt template
// 1. Transaction type classification
const PromptTypeClassification = `
Transaction:
"%s"

Reply *only* in JSON:
{"type":"INCOME"|"EXPENSE"}
`
//========================================== first prompt template

//========================================== second prompt template
// 2. Full extraction: type, amount, category, date
const PromptFullClassification = `
Current Date: "%s"
Transaction: "%s"
Type (INCOME/EXPENSE): "%s"

Interpret any relative time expressions based on the current date above.

Reply *only* in JSON:
{"type": "%s","amount": number,"category_id": number,"date": "YYYY-MM-DD" or null}

Use one of these categories:
%s
`
//========================================== second prompt template