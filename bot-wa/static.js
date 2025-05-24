// Static messages
const WelcomeText = `👋 Hai, selamat datang di Duite Bot  
Aku siap bantu kamu catat pengeluaran harian dengan cepat dan tanpa ribet.

Cukup kirim pesan seperti ini:
➡️ Makan siang 25k 
atau seperti ini
➡️ gaji masuk tgl 20 mei 2025

Setelah itu, aku akan langsung menyimpan transaksi kamu secara otomatis! 🔥`;

const HelpText = `📘 *Bantuan dari Duite Bot*

Aku bisa bantu kamu dengan perintah-perintah berikut:

• start - Aku akan membukakan menu untukmu
• harian - Aku akan tampilkan laporan transaksi hari ini
• bulanan - Aku akan tunjukkan laporan transaksi bulan ini
• hapus - Aku bisa hapus data transaksimu, cukup kirim: hapus (ID transaksi), contoh: hapus 123
• daftar - Aku akan daftarkan akunmu supaya bisa lihat dashboard
• dashboard - Aku akan tampilkan Dashboard kamu
• bantuan - Aku akan kirimkan informasi bantuan seperti ini`;

const ErrorMessage = '⚠️ Terjadi kesalahan saat menghubungi server.';
const CloseMenuText = '❌ Menu ditutup. Gunakan *start* untuk membuka kembali menu.';
const WrongCommandText = '❓ Perintah tidak dikenali. Ketik *bantuan* untuk melihat daftar perintah.';

module.exports = { HelpText, WelcomeText, ErrorMessage, CloseMenuText, WrongCommandText };