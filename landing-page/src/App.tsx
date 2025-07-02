import { Button } from "@/components/ui/button"
import { Card, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { ThemeToggle } from "@/components/theme-toggle"
import {
  ArrowRight,
  BarChart3,
  MessageCircle,
  PieChart,
  Smartphone,
  TrendingUp,
  Zap,
  Shield,
  Clock,
  Users,
  Check,
} from "lucide-react"

export default function App() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-50 to-blue-50 dark:from-slate-900 dark:to-slate-800 transition-colors">
      {/* Header */}
      <header className="border-b bg-white/80 dark:bg-slate-900/80 backdrop-blur-sm sticky top-0 z-50 transition-colors">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <img src="/Duite Bot.png" alt="Logo Duite Bot" className="h-8 w-8" />
            <span className="text-xl font-bold text-gray-900 dark:text-white">Duite Bot</span>
          </div>
          <nav className="hidden md:flex items-center space-x-6">
            <a href="#features" className="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">
              Fitur
            </a>
            <a href="#dashboard" className="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">
              Dasbor
            </a>
            <a href="#pricing" className="text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors">
              Harga
            </a>
            <ThemeToggle />
            <Button size="sm">Join Sekarang</Button>
          </nav>
        </div>
      </header>

      {/* Hero Section */}
      <section className="py-20 px-4">
        <div className="container mx-auto text-center max-w-4xl">
          <Badge variant="secondary" className="mb-4">
            <Zap className="h-3 w-3 mr-1" />
            AI-Powered Financial Tracking
          </Badge>
          <h1 className="text-5xl md:text-6xl font-bold text-gray-900 dark:text-white mb-6 leading-tight transition-colors">
            Catat Keuangan Anda dengan
            <span className="text-blue-600 dark:text-blue-400 block">Perintah Chat yang Sederhana</span>
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-300 mb-8 max-w-2xl mx-auto transition-colors">
            Catat pengeluaran, pemasukan, dan transaksi keuangan Anda dengan mudah melalui Telegram dan WhatsApp. Dapatkan wawasan instan melalui dasbor analitik berbasis AI kami.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center mb-12">
            <Button
              asChild
              size="lg"
              className="text-lg px-8"
            >
              <a
                href="#cta"
              >
                Coba Gratis
                <ArrowRight className="ml-2 h-5 w-5" />
              </a>
            </Button>
          </div>

          {/* Platform Badges */}
          <div className="flex justify-center items-center gap-6 mb-12">
            <div className="flex items-center gap-2 bg-white dark:bg-slate-800 rounded-full px-4 py-2 shadow-sm transition-colors">
              <MessageCircle className="h-5 w-5 text-blue-500" />
              <span className="font-medium dark:text-white">Telegram</span>
            </div>
            <div className="flex items-center gap-2 bg-white dark:bg-slate-800 rounded-full px-4 py-2 shadow-sm transition-colors">
              <MessageCircle className="h-5 w-5 text-green-500" />
              <span className="font-medium dark:text-white">WhatsApp</span>
            </div>
          </div>

          {/* Hero Image */}
          <div className="relative max-w-4xl mx-auto mb-8">
            <div className="bg-white dark:bg-slate-800 rounded-2xl shadow-2xl p-2 border dark:border-slate-700 transition-colors">
              <img
                src="./dashboard.png"
                alt="Pratinjau Dasbor Duite Bot"
                width={800}
                height={400}
                className="w-full rounded-lg"
                loading="lazy"
              />
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section id="features" className="py-20 px-4 bg-white dark:bg-slate-900 transition-colors">
        <div className="container mx-auto max-w-6xl">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 dark:text-white mb-4 transition-colors">Semua yang Anda Butuhkan untuk Mengelola Keuangan</h2>
            <p className="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto transition-colors">
              Fitur canggih yang dirancang untuk membuat pelacakan keuangan semudah mengirim pesan
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
            <Card className="border-0 shadow-lg hover:shadow-xl transition-shadow dark:bg-slate-800 dark:hover:bg-slate-700">
              <CardHeader>
                <MessageCircle className="h-12 w-12 text-blue-600 mb-4" />
                <CardTitle className="dark:text-white">Natural Language Processing</CardTitle>
                <CardDescription className="dark:text-gray-300">
                  Cukup ketik "Jajan cilok 25k", dan AI kami akan langsung memahaminya dan mengkategorikannya.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-0 shadow-lg hover:shadow-xl transition-shadow dark:bg-slate-800 dark:hover:bg-slate-700">
              <CardHeader>
                <BarChart3 className="h-12 w-12 text-green-600 mb-4" />
                <CardTitle className="dark:text-white">Real-time Analytics</CardTitle>
                <CardDescription className="dark:text-gray-300">
                  Dapatkan wawasan instan tentang pola pengeluaran, tren pendapatan, dan kesehatan finansial Anda melalui visualisasi menarik.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-0 shadow-lg hover:shadow-xl transition-shadow dark:bg-slate-800 dark:hover:bg-slate-700">
              <CardHeader>
                <Smartphone className="h-12 w-12 text-purple-600 mb-4" />
                <CardTitle className="dark:text-white">Multi-Platform Support</CardTitle>
                <CardDescription className="dark:text-gray-300">
                  Akses data keuangan Anda secara lancar melalui Telegram, WhatsApp, dan dasbor web/mobile app kami.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-0 shadow-lg hover:shadow-xl transition-shadow dark:bg-slate-800 dark:hover:bg-slate-700">
              <CardHeader>
                <Shield className="h-12 w-12 text-red-600 mb-4" />
                <CardTitle className="dark:text-white">User Data Protection</CardTitle>
                <CardDescription className="dark:text-gray-300">
                  Data transaksi Anda dilindungi dari akses tidak sah dan hanya digunakan untuk keperluan layanan.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-0 shadow-lg hover:shadow-xl transition-shadow dark:bg-slate-800 dark:hover:bg-slate-700">
              <CardHeader>
                <Clock className="h-12 w-12 text-orange-600 mb-4" />
                <CardTitle className="dark:text-white">Instant Tracking</CardTitle>
                <CardDescription className="dark:text-gray-300">
                  Catat transaksi dalam hitungan detik. Tanpa formulir, tanpa aplikasi – cukup kirim pesan.
                </CardDescription>
              </CardHeader>
            </Card>

            <Card className="border-0 shadow-lg hover:shadow-xl transition-shadow dark:bg-slate-800 dark:hover:bg-slate-700">
              <CardHeader>
                <TrendingUp className="h-12 w-12 text-indigo-600 mb-4" />
                <CardTitle className="dark:text-white">Smart Insights</CardTitle>
                <CardDescription className="dark:text-gray-300">
                  Dapatkan saran keuangan dan rekomendasi pengeluaran berdasarkan riwayat transaksi Anda.
                </CardDescription>
              </CardHeader>
            </Card>
          </div>
        </div>
      </section>

      {/* Dashboard Preview */}
      <section id="dashboard" className="py-20 px-4 bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-slate-800 dark:to-slate-700 transition-colors">
        <div className="container mx-auto max-w-6xl">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 dark:text-white mb-4 transition-colors">Powerful Analytics Dashboard</h2>
            <p className="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto transition-colors">
              Visualisasikan data keuangan Anda dengan grafik interaktif dan laporan lengkap
            </p>
          </div>

          <div className="grid lg:grid-cols-2 gap-8 items-center">
            <div className="space-y-6">
              <Card className="p-6 dark:bg-slate-800 dark:border-slate-700">
                <div className="flex items-center gap-4 mb-4">
                  <PieChart className="h-8 w-8 text-blue-600" />
                  <div>
                    <h3 className="font-semibold text-lg dark:text-white">Expense Categories</h3>
                    <p className="text-gray-600 dark:text-gray-300">Track spending by category</p>
                  </div>
                </div>
                <div className="h-32 bg-gradient-to-r from-blue-100 to-purple-100 dark:from-blue-900/20 dark:to-purple-900/20 rounded-lg flex items-center justify-center">
                  <span className="text-gray-500 dark:text-gray-400">Interactive Pie Chart</span>
                </div>
              </Card>

              <Card className="p-6 dark:bg-slate-800 dark:border-slate-700">
                <div className="flex items-center gap-4 mb-4">
                  <BarChart3 className="h-8 w-8 text-green-600" />
                  <div>
                    <h3 className="font-semibold text-lg dark:text-white">Monthly Trends</h3>
                    <p className="text-gray-600 dark:text-gray-300">Income vs expenses over time</p>
                  </div>
                </div>
                <div className="h-32 bg-gradient-to-r from-green-100 to-blue-100 dark:from-green-900/20 dark:to-blue-900/20 rounded-lg flex items-center justify-center">
                  <span className="text-gray-500 dark:text-gray-400">Interactive Bar Chart</span>
                </div>
              </Card>
            </div>

            <div className="bg-white dark:bg-slate-800 rounded-2xl shadow-2xl p-4 border dark:border-slate-700 max-w-md mx-auto transition-colors">
              <img
                src="./example.gif"
                width={300}
                height={250}
                className="w-full h-auto rounded-lg"
                loading="lazy"
                alt="App preview"
              />
            </div>

          </div>
        </div>
      </section>

      {/* Pricing Section */}
      <section id="pricing" className="py-20 px-4 bg-slate-50 dark:bg-slate-900 transition-colors">
        <div className="container mx-auto max-w-6xl">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 dark:text-white mb-4 transition-colors">Pilih Mode Finansial Kamu 🔥</h2>
            <p className="text-xl text-gray-600 dark:text-gray-300 max-w-2xl mx-auto transition-colors">
              Mulai dari yang santuy sampai yang full komit — sesuaikan sama gaya kamu. No tipu-tipu, harga transparan!
            </p>
          </div>

          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8 items-stretch">
            {/* Mode Santuy */}
            <Card className="flex flex-col border-gray-300 dark:border-slate-600 shadow-lg dark:bg-slate-800">
              <CardHeader className="pb-6">
                <CardTitle className="text-2xl font-semibold mb-2 dark:text-white">Mode Santuy</CardTitle>
                <CardDescription className="text-gray-600 dark:text-gray-300 h-12">
                  Buat kamu yang baru mulai nyatet keuangan.
                </CardDescription>
                <p className="text-4xl font-bold text-gray-900 dark:text-white mt-4">
                  Rp0<span className="text-xl font-normal text-gray-500 dark:text-gray-400">/bulan</span>
                </p>
              </CardHeader>
              <div className="p-6 flex-grow">
                <ul className="space-y-3 text-gray-700 dark:text-gray-300">
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>50 transaksi/bulan</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Analisa sederhana</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Integrasi Telegram & WhatsApp</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Dukungan komunitas</span>
                  </li>
                </ul>
              </div>
              <div className="p-6 mt-auto">
                <a href="#cta" className="w-full">
                  <Button variant="outline" className="w-full text-lg py-6">
                    Gaskeun Gratis
                  </Button>
                </a>
              </div>
            </Card>

            {/* Mode Serius */}
            <Card className="flex flex-col border-blue-500 ring-2 ring-blue-500 shadow-2xl relative dark:bg-slate-800">
              <Badge variant="default" className="absolute -top-3 left-1/2 -translate-x-1/2 bg-blue-600 text-white">
                Paling Laris
              </Badge>
              <CardHeader className="pb-6">
                <CardTitle className="text-2xl font-semibold mb-2 text-blue-600 dark:text-blue-400">Mode Serius</CardTitle>
                <CardDescription className="text-gray-600 dark:text-gray-300 h-12">
                  Buat yang mulai peduli dan pengen insight mendalam.
                </CardDescription>
                <p className="text-4xl font-bold text-gray-900 dark:text-white mt-4">
                  Rp15.000<span className="text-xl font-normal text-gray-500 dark:text-gray-400">/bulan</span>
                </p>
                <p className="text-sm text-red-500 line-through mt-1">
                  Rp20.000
                </p>
                <p className="text-sm text-green-600 font-semibold">
                  Hemat Rp5.000 (11%)
                </p>
              </CardHeader>
              <div className="p-6 flex-grow">
                <ul className="space-y-3 text-gray-700 dark:text-gray-300">
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Transaksi tanpa batas</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Analitik dan laporan mendalam</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Insight & rekomendasi pintar</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Dukungan prioritas</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Ekspor data (CSV, EXCEL, PDF)</span>
                  </li>
                </ul>
              </div>
              <div className="p-6 mt-auto">
                <a href="#cta" className="w-full">
                  <Button className="w-full text-lg py-6 bg-blue-600 hover:bg-blue-700">
                    Upgrade Sekarang
                  </Button>
                </a>
              </div>
            </Card>

            {/* Mode Komitmen */}
            <Card className="flex flex-col border-gray-300 dark:border-slate-600 shadow-lg dark:bg-slate-800">
              <CardHeader className="pb-6">
                <CardTitle className="text-2xl font-semibold mb-2 dark:text-white">Mode Komitmen</CardTitle>
                <CardDescription className="text-gray-600 dark:text-gray-300 h-12">
                  Buat tim atau bisnis yang udah serius dan butuh kontrol penuh.
                </CardDescription>
                <p className="text-4xl font-bold text-gray-900 dark:text-white mt-4">
                  Rp40.000<span className="text-xl font-normal text-gray-500 dark:text-gray-400">/3 bulan</span>
                </p>
                <p className="text-sm text-red-500 line-through mt-1">
                  Rp60.000
                </p>
                <p className="text-sm text-green-600 font-semibold">
                  Hemat Rp20.000 (20%)
                </p>
              </CardHeader>
              <div className="p-6 flex-grow">
                <ul className="space-y-3 text-gray-700 dark:text-gray-300">
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Semua fitur di Mode Serius</span>
                  </li>
                  <li className="flex items-center">
                    <Check className="h-5 w-5 text-green-500 mr-2 shrink-0" />
                    <span>Branding kustom</span>
                  </li>
                </ul>
              </div>
              <div className="p-6 mt-auto">
                <a href="#cta" className="w-full">
                  <Button variant="outline" className="w-full text-lg py-6">
                    Hubungi Tim Kami
                  </Button>
                </a>
              </div>
            </Card>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section className="py-20 px-4 bg-white dark:bg-slate-900 transition-colors">
        <div className="container mx-auto max-w-4xl">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 dark:text-white mb-4 transition-colors">Cara Kerjanya</h2>
            <p className="text-xl text-gray-600 dark:text-gray-300 transition-colors">Mulai dalam tiga langkah mudah</p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            <div className="text-center">
              <div className="bg-blue-100 dark:bg-blue-900/20 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                <span className="text-2xl font-bold text-blue-600 dark:text-blue-400">1</span>
              </div>
              <h3 className="text-xl font-semibold mb-2 dark:text-white">Hubungkan Platform Anda</h3>
              <p className="text-gray-600 dark:text-gray-300">Tambahkan bot kami ke Telegram atau WhatsApp Anda</p>
            </div>

            <div className="text-center">
              <div className="bg-green-100 dark:bg-green-900/20 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                <span className="text-2xl font-bold text-green-600 dark:text-green-400">2</span>
              </div>
              <h3 className="text-xl font-semibold mb-2 dark:text-white">Mulai Percakapan</h3>
              <p className="text-gray-600 dark:text-gray-300">Kirimkan pesan seperti biasa untuk mencatat transaksi Anda secara otomatis</p>
            </div>

            <div className="text-center">
              <div className="bg-purple-100 dark:bg-purple-900/20 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                <span className="text-2xl font-bold text-purple-600 dark:text-purple-400">3</span>
              </div>
              <h3 className="text-xl font-semibold mb-2 dark:text-white">Lihat Ringkasan & Analisis</h3>
              <p className="text-gray-600 dark:text-gray-300">Akses dasbor Anda untuk melihat laporan keuangan yang lengkap dan mudah dipahami.</p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section id="cta" className="py-20 px-4 bg-gradient-to-r from-blue-600 to-purple-600 text-white">
        <div className="container mx-auto text-center max-w-4xl">
          <h2 className="text-4xl font-bold mb-4">Siap Mengubah Cara Anda Melacak Keuangan?</h2>
          <p className="text-xl mb-8 opacity-90">
            Bergabunglah bersama ratusan pengguna lainnya yang telah mengelola keuangan mereka dengan lebih mudah melalui Duite Bot.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center flex-wrap">
            {/* Tombol Uji Coba */}
            <Button size="lg" variant="secondary" className="text-lg px-8">
              Mulai Uji Coba Gratis
              <ArrowRight className="ml-2 h-5 w-5" />
            </Button>

            {/* Tombol Gabung Telegram */}
            <Button
              asChild
              size="lg"
              variant="outline"
              className="text-lg px-8 border-white text-black hover:bg-white hover:text-blue-600"
            >
              <a
                href="https://t.me/grup_duite_bot"
                target="_blank"
                rel="noopener noreferrer"
              >
              <MessageCircle className="h-5 w-5 text-blue-500" /> Telegram
              </a>
            </Button>

            {/* Tombol Gabung WhatsApp */}
            <Button
              asChild
              size="lg"
              variant="outline"
              className="text-lg px-8 border-white text-black hover:bg-white hover:text-blue-600"
            >
              <a
                href="https://chat.whatsapp.com/LqK5R0OCUIOCqCwwE0dWDA"
                target="_blank"
                rel="noopener noreferrer"
              >
                <MessageCircle className="h-5 w-5 text-green-500" /> WhatsApp
              </a>
            </Button>
          </div>
          <div className="mt-8 flex items-center justify-center gap-8 text-sm opacity-75">
            <div className="flex items-center gap-2">
              <Users className="h-4 w-4" />
              <span>Lebih dari 100 Pengguna Aktif</span>
            </div>
            <div className="flex items-center gap-2">
              <Shield className="h-4 w-4" />
              <span>User Data Protection</span>
            </div>
            <div className="flex items-center gap-2">
              <Clock className="h-4 w-4" />
              <span>24/7 Support</span>
            </div>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="bg-gray-900 dark:bg-slate-950 text-white py-12 px-4 transition-colors">
        <div className="container mx-auto max-w-6xl">
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center space-x-2 mb-4">
                <img src="/Duite Bot.png" alt="Logo Duite Bot" className="h-8 w-8" />
                <span className="text-lg font-bold">Duite Bot</span>
              </div>
              <p className="text-gray-400 dark:text-gray-500">Membantu Anda mengelola keuangan secara efisien melalui chat pintar berbasis AI.</p>
            </div>
            <div>
              <h3 className="font-semibold mb-4">Product</h3>
              <ul className="space-y-2 text-gray-400 dark:text-gray-500">
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Fitur
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Dasbor
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Harga
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    API
                  </a>
                </li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold mb-4">Support</h3>
              <ul className="space-y-2 text-gray-400 dark:text-gray-500">
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Help Center
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Contact Us
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Privacy Policy
                  </a>
                </li>
                <li>
                  <a href="#" className="hover:text-white transition-colors">
                    Terms of Service
                  </a>
                </li>
              </ul>
            </div>
            <div>
              <h3 className="font-semibold mb-4">Connect</h3>
              <ul className="space-y-2 text-gray-400 dark:text-gray-500">
                <li>
                  <a href="https://t.me/grup_duite_bot" className="hover:text-white transition-colors">
                    Telegram Bot
                  </a>
                </li>
                <li>
                  <a href="https://chat.whatsapp.com/LqK5R0OCUIOCqCwwE0dWDA" className="hover:text-white transition-colors">
                    WhatsApp Bot
                  </a>
                </li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 dark:border-slate-700 mt-8 pt-8 text-center text-gray-400 dark:text-gray-500">
            <p>&copy; 2025-{new Date().getFullYear()} Duite Bot. All rights reserved. </p>
          </div>
        </div>
      </footer>
    </div>
  )
};