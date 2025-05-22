const { default: makeWASocket, useMultiFileAuthState, DisconnectReason, fetchLatestBaileysVersion } = require('@whiskeysockets/baileys');
const { Boom } = require('@hapi/boom');
const qrcode = require('qrcode-terminal');
const path = require('path');

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
      // Tampilkan QR code di terminal
      qrcode.generate(qr, { small: true });
    }

    if (connection === 'close') {
      const shouldReconnect = (lastDisconnect.error = new Boom(lastDisconnect?.error))?.output?.statusCode !== DisconnectReason.loggedOut;
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

    console.log(`📩 Pesan dari ${from}: ${body}`);

    if (body.toLowerCase() === 'ping') {
      await sock.sendMessage(from, { text: 'pong!' });
    } else if (body.toLowerCase() === 'johar cupu') {
      const menu = `sudah pasti cupu, ga usah di tanya lagi. cupu tidak bisa ngoding`;
      await sock.sendMessage(from, { text: menu });
    } else if (body.toLowerCase() === 'menu') {
      const menu = `📋 *MENU BOT*\n\n1. ping - Tes bot\n2. info - Informasi bot`;
      await sock.sendMessage(from, { text: menu });
    } else if (body.toLowerCase() === 'info') {
      await sock.sendMessage(from, {
        text: '🤖 Bot ini dibuat dengan Baileys dan Node.js!',
      });
    }
  });
}

startBot();
