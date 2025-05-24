const formatRupiah = (amount) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(amount);
}

const formatDailyReport = (transactions) => {
  let report = "📊 *Laporan Hari Ini*\n";
  let totalIn = 0;
  let totalOut = 0;

  transactions.forEach((tx) => {
    let transactionType = "";
    if (tx.transactionType === "EXPENSE") {
      transactionType = "🔴";
      totalOut += tx.amount;
    } else if (tx.transactionType === "INCOME") {
      transactionType = "🟢";
      totalIn += tx.amount;
    }

    const formattedAmount = formatRupiah(tx.amount);
    report += `#${tx.id} ${transactionType} ${formattedAmount} ${tx.original_text}\n`;
  });

  report += "\n";
  report += `🟢 Total Pemasukan: ${formatRupiah(totalIn)}\n`;
  report += `🔴 Total Pengeluaran: ${formatRupiah(totalOut)}\n`;

  return report;
}

const formatMonthlyReport = (transactions) => {
  let report = "📊 *Laporan Bulan Ini*\n";
  let totalIn = 0;
  let totalOut = 0;

  transactions.forEach((tx) => {
    let transactionType = "";
    if (tx.transactionType === "EXPENSE") {
      transactionType = "🔴";
      totalOut += tx.amount;
    } else if (tx.transactionType === "INCOME") {
      transactionType = "🟢";
      totalIn += tx.amount;
    }

    const formattedAmount = formatRupiah(tx.amount);
    report += `#${tx.id} ${transactionType} ${formattedAmount} ${tx.original_text}\n`;
  });

  report += "\n";
  report += `🟢 Total Pemasukan: ${formatRupiah(totalIn)}\n`;
  report += `🔴 Total Pengeluaran: ${formatRupiah(totalOut)}\n`;

  return report;
}

const slugify = (text) => {
  return text
    .toString()                    // pastikan input berupa string
    .normalize('NFD')              // pisahkan aksen (é -> e + ́)
    .replace(/[\u0300-\u036f]/g, '') // hapus aksen
    .toLowerCase()                 // jadi huruf kecil semua
    .trim()                        // hapus spasi di awal/akhir
    .replace(/[^a-z0-9 -]/g, '')   // hapus karakter yang bukan huruf/angka/spasi
    .replace(/\s+/g, '-')          // ganti spasi dengan tanda hubung
    .replace(/-+/g, '-');          // hapus tanda hubung dobel
}

function encodeChatID(chatID) {
  // Konversi chatID ke string
  const chatIDStr = chatID.toString();
  // Encode ke base64 dengan standar URL-safe
  return Buffer.from(chatIDStr, 'utf-8')
    .toString('base64')
    .replace(/\+/g, '-')
    .replace(/\//g, '_');
}

function translateTransactionType(type) {
  if (type === 'INCOME') return 'Pemasukan';
  if (type === 'EXPENSE') return 'Pengeluaran';
  return 'Tidak Dikenal'; // fallback
}


const formatDate = (date) => {
  const day = String(date.getDate()).padStart(2, '0');
  const month = String(date.getMonth() + 1).padStart(2, '0'); // getMonth() is zero-based
  const year = String(date.getFullYear()).slice(-2); // Get last two digits of year

  return `${day}/${month}/${year}`;
}

module.exports = {
  formatRupiah,
  formatDailyReport,
  formatMonthlyReport,
  slugify,
  translateTransactionType,
  formatDate,
  encodeChatID
}