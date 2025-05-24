package static

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