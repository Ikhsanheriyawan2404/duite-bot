const { default: makeWASocket, useMultiFileAuthState, DisconnectReason, fetchLatestBaileysVersion } = require('@whiskeysockets/baileys');
const { Boom } = require('@hapi/boom');
const qrcode = require('qrcode-terminal');
const path = require('path');
const { config } = require('./config') 
const { WelcomeText, HelpText, ErrorMessage, WrongCommandText, PromptText } = require('./static');
const { getDailyTransaction, registerUser, checkUser, getMonthlyTransaction, hitAiClassifyTransaction, saveTransaction } = require('./api-client');
const { formatDailyReport, formatRupiah, slugify, translateTransactionType, formatMonthlyReport } = require('./utils');

async function startBot() {
  const { state, saveCreds } = await useMultiFileAuthState(path.resolve('./sessions'));
  const { version, isLatest } = await fetchLatestBaileysVersion();
  console.log(`Menggunakan WhatsApp versi ${version.join('.')}, terbaru: ${isLatest}`);

  const sock = makeWASocket({
    version,
    auth: state,
  });

  sock.ev.on('creds.update', saveCreds);

  sock.ev.on('connection.update', (update) => {
    const { connection, lastDisconnect, qr } = update;

    if (qr) {
      qrcode.generate(qr, { small: true });
    }

    if (connection === 'close') {
      const shouldReconnect = new Boom(lastDisconnect?.error)?.output?.statusCode !== DisconnectReason.loggedOut;
      console.log('Koneksi terputus, mencoba menyambung kembali:', shouldReconnect);
      if (shouldReconnect) {
        startBot();
      }
    } else if (connection === 'open') {
      console.log('🟢 Bot terhubung!');
    }
  });

  sock.ev.on('messages.upsert', async ({ messages, type }) => {
    if (type !== 'notify') return;
    const msg = messages[0];
    if (!msg.message || msg.key.fromMe) return;

    const from = msg.key.remoteJid;
    const body = msg.message.conversation || msg.message.extendedTextMessage?.text || '';
    const replaceTextPhoneNumber = from.replace('@s.whatsapp.net', '');
    const senderNumber = Number(replaceTextPhoneNumber)

    console.log(`📩 Pesan dari ${from}: ${body}`);

    try {
      const command = body.trim().toLowerCase();
      const [mainCommand, ...args] = command.split(' ');
    
      switch (mainCommand) {
        case 'start':
          await sock.sendMessage(from, { text: WelcomeText });
          break;

        case 'bantuan':
          await sock.sendMessage(from, { text: HelpText });
          break;
    
        case 'harian':
          const dailyTx = await getDailyTransaction(senderNumber);
          if (dailyTx.length === 0) {
            await sock.sendMessage(from, { text: '📅 Hari ini dompet kamu anteng banget~ Ga ada transaksi masuk/keluar nih 😎' });
            break;
          }

          await sock.sendMessage(from, { text: formatDailyReport(dailyTx) });
          break;
    
        case 'bulanan':
          const monthlyTx = await getMonthlyTransaction(senderNumber);
          if (monthlyTx.length === 0) {
            await sock.sendMessage(from, { text: '📅 No transaksi bulan ini... lagi mode hemat atau lupa nyatet nih? 🤔📉' });
            break;
          }

          await sock.sendMessage(from, { text: formatMonthlyReport(monthlyTx) });
          break;
    
        case 'hapus':
          const id = args[0];
          if (id) {
            await sock.sendMessage(from, { text: `🗑️ Menghapus transaksi dengan ID ${id}...` });
            // TODO: Hapus transaksi
          } else {
            await sock.sendMessage(from, { text: '⚠️ Format salah. Contoh: hapus 123' });
          }
          break;
    
        case 'daftar':
          const fullName = args.join(' ').trim();

          if (!fullName) {
            await sock.sendMessage(from, {
              text: `Waduuhh, kamu belum isi nama nih 😅\nCoba ketik kayak gini ya:\n👉 *daftar Udin Andria*`,
            });
            break;
          }
          const resultRegister = await registerUser(senderNumber, fullName);
          if(resultRegister.error) {
            await sock.sendMessage(from, { text: 'Eh, btw kamu udah daftar sebelumnya, hehe'});
          } else {
            await sock.sendMessage(from, {
              text: `Hai ${resultRegister.name}, mau aku bantu lihat dashboard?\nklik disini ya ${config.DASHBOARD_URL}`
            });
          }
          break;
    
        case 'dashboard':
          const checkedUser = await checkUser(senderNumber);
          if (checkedUser.exist) {
            await sock.sendMessage(from, {
              text: `📊 Dashboard klik disini ya 👉 ${config.DASHBOARD_URL}`,
            });
          } else {
            await sock.sendMessage(from, {
              text: 'Yuk, daftar dulu biar bisa lanjut!\nGampang kok, tinggal ketik: daftar NamaKamu\nContoh: daftar Budi Corcodillo',
            });
          }
          break;
    
        default:
          const fullPrompt = PromptText.replace('%s', body);
          const result = await hitAiClassifyTransaction(senderNumber, fullPrompt);
          const { type, amount, category, date } = result.result;

          // const typeText = translateTransactionType(type);

          let transactionDate = new Date(); // default: sekarang

          if (date && typeof date === 'string' && date.trim() !== '') {
            const parsedDate = new Date(date); // format harus "YYYY-MM-DD"
            if (!isNaN(parsedDate)) {
              transactionDate = parsedDate;
            }
          }
          
          const response = await saveTransaction({
            chat_id: senderNumber,
            original_text: body,
            transaction_type: type,
            amount,
            category: slugify(category),
            transaction_date: transactionDate,
          });

          if (response.error) {
            await sock.sendMessage(from, { text: '❌ Gagal menyimpan transaksi. Silakan coba lagi.' });
            break;
          }

          let replyLines = [
            '✅ Siap, transaksi kamu udah ke-record! 🎉',
            '',
            `📂 Tipe     : ${type === 'INCOME' ? '🟢 Pemasukan' : '🔴 Pengeluaran'}`,
            `💰 Nominal  : ${formatRupiah(amount)}`,
            `🏷️ Kategori : ${category}`
          ];
          
          if (date) {
            replyLines.push(`🗓️ Tanggal  : ${date}`);
          }
          
          const reply = replyLines.join('\n');
          await sock.sendMessage(from, { text: reply });
      }
    } catch (err) {
      console.error('❌ Error handling message:', err);
      await sock.sendMessage(from, { text: ErrorMessage });
    }
  });
}

startBot();
