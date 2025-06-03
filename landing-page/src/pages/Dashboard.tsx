"use client"

import { useState, useMemo, useCallback, useEffect } from "react"
import { Button } from "@/components/ui/button"
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  DialogDescription,
} from "@/components/ui/dialog"
import { TransactionFilters } from "@/components/dashboard/transaction-filters"
import { SummaryCard } from "@/components/dashboard/summary-card"
import { FinancialChart } from "@/components/dashboard/financial-chart" // Monthly Income/Expense Bar Chart
import { ExpenseCategoryPieChart } from "@/components/dashboard/expense-category-pie-chart" // New
import { BalanceTrendChart } from "@/components/dashboard/balance-trend-chart" // New
import { TransactionsTable } from "@/components/dashboard/transactions-table"
import { TransactionForm } from "@/components/dashboard/transaction-form"
import type { Transaction, MonthlyData, CategoryExpenseData, BalanceTrendData } from "@/lib/types"
import { TrendingUp, TrendingDown, PlusCircle, AlertTriangle, DollarSign } from "lucide-react"
import { format } from "date-fns"
import { DateRange } from "react-day-picker"
import { getMe, magicLogin } from "@/api/auth"
import { fetchTransactions } from "@/api/transactions"
import { formatDate } from "@/lib/utils"

const initialTransactions: Transaction[] = []

export default function DashboardPage() {
  
  const loadFilterTransactions = async () => {
    try {
      const typeParam = selectedType === "all" ? undefined : selectedType;
      const categoryParam = selectedCategory === "all" ? undefined : selectedCategory;

      const fetched = await fetchTransactions(
        typeParam,
        categoryParam,
        dateRange.from,
        dateRange.to
      );
      console.log({fetched})
      setTransactions(fetched);
    } catch (error) {
      console.error("Failed to fetch transactions:", error);
    }
  };

  const [transactions, setTransactions] = useState<Transaction[]>(initialTransactions)
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
  const [currentTransaction, setCurrentTransaction] = useState<Transaction | undefined>(undefined)
  const [transactionToDeleteId, setTransactionToDeleteId] = useState<string | null>(null)

  const now = new Date();
  const awalBulan = new Date(now.getFullYear(), now.getMonth(), 1);
  const [dateRange, setDateRange] = useState<DateRange>({
    from: awalBulan,
    to: now,
  });
  const [selectedType, setSelectedType] = useState<string>("all")
  const [selectedCategory, setSelectedCategory] = useState<string>("all")

  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search)
    const token = urlParams.get("ref")

    const doAuthAndFetch = async () => {
      try {
        if (token) {
          // hanya panggil magic login jika ada token
          await magicLogin(token)
          // hilangkan ref token dari URL agar tidak terus login ulang
          const cleanUrl = window.location.origin + window.location.pathname
          window.history.replaceState({}, document.title, cleanUrl)
        }

        // lanjutkan login menggunakan token yang sudah disimpan
        await getMe()
        await loadFilterTransactions();
      } catch (err) {
        console.error("Login or fetch error:", err)
        // optional: redirect ke halaman login
      }
    }

    doAuthAndFetch()
  }, [])

  useEffect(() => {
    loadFilterTransactions();
  }, [selectedType, selectedCategory, dateRange]);

  const uniqueCategories = useMemo(() => {
    const categories = new Set(transactions.map((t) => t.category))
    return ["all", ...Array.from(categories).sort()]
  }, [transactions])

  const summaryData = useMemo(() => {
    // Summary should reflect ALL transactions, not filtered ones
    const totalIncome = transactions.filter((t) => t.transaction_type === "INCOME").reduce((sum, t) => sum + t.amount, 0)
    const totalExpenses = transactions.filter((t) => t.transaction_type === "EXPENSE").reduce((sum, t) => sum + t.amount, 0)
    const balance = totalIncome - totalExpenses
    return { totalIncome, totalExpenses, balance }
  }, [transactions])

  // Data for Monthly Income/Expense Bar Chart (FinancialChart)
  const monthlyIncomeExpenseData = useMemo((): MonthlyData[] => {
    const monthlyMap: Record<string, { income: number; expense: number }> = {}
    transactions.forEach((transaction) => {
      const date = new Date(transaction.transaction_date)
      const monthYear = format(date, "MMM yyyy")
      if (!monthlyMap[monthYear]) {
        monthlyMap[monthYear] = { income: 0, expense: 0 }
      }
      if (transaction.transaction_type === "INCOME") {
        monthlyMap[monthYear].income += transaction.amount
      } else if (transaction.transaction_type === "EXPENSE") {
        monthlyMap[monthYear].expense += transaction.amount
      }
    })
    return Object.entries(monthlyMap)
      .map(([month, data]) => ({ month, ...data }))
      .sort((a, b) => new Date(a.month).getTime() - new Date(b.month).getTime())
  }, [transactions])

  // Data for Expense Category Pie Chart
  const expenseCategoryData = useMemo((): CategoryExpenseData[] => {
    const categoryMap: Record<string, number> = {}
    transactions
      .filter((t) => t.transaction_type === "EXPENSE")
      .forEach((transaction) => {
        categoryMap[transaction.category] = (categoryMap[transaction.category] || 0) + transaction.amount
      })
    return Object.entries(categoryMap)
      .map(([category, totalExpense]) => ({ category, totalExpense }))
      .sort((a, b) => b.totalExpense - a.totalExpense) // Sort for consistent pie chart segment order
  }, [transactions])

  // Data for Balance Trend Line Chart
  // const balanceTrendData = useMemo((): BalanceTrendData[] => {
  //   if (transactions.length === 0) return []

  //   const trend: BalanceTrendData[] = []
  //   let currentBalance = 0

  //   // Calculate initial balance before the filtered period if filters are applied
  //   // This gives a more accurate starting point for the trend within the filtered view.
  //   // For simplicity in this example, we'll start the trend from the first filtered transaction.
  //   // A more robust solution would calculate balance from all transactions up to the start of the filtered period.

  //   const transactionsByMonth: Record<string, Transaction[]> = {}
  //   transactions.forEach((t) => {
  //     const monthYear = format(new Date(t.transaction_date), "MMM yyyy")
  //     if (!transactionsByMonth[monthYear]) {
  //       transactionsByMonth[monthYear] = []
  //     }
  //     transactionsByMonth[monthYear].push(t)
  //   })

  //   const sortedMonths = Object.keys(transactionsByMonth).sort((a, b) => new Date(a).getTime() - new Date(b).getTime())

  //   // Find overall starting balance before any filtered transactions
  //   // This is complex if we want the true historical balance.
  //   // For this example, let's assume the balance trend starts from 0 before the first transaction in the *original* dataset,
  //   // and then we calculate the balance up to the start of the filtered period.
  //   let balanceBeforeFilteredPeriod = 0
  //   const firstFilteredDate = transactions.length > 0 ? new Date(transactions[0].transaction_date) : new Date()

  //   initialTransactions.forEach((t) => {
  //     if (new Date(t.transaction_date) < firstFilteredDate) {
  //       balanceBeforeFilteredPeriod += t.transaction_type === "income" ? t.amount : -t.amount
  //     }
  //   })
  //   currentBalance = balanceBeforeFilteredPeriod

  //   sortedMonths.forEach((month) => {
  //     transactionsByMonth[month].forEach((t) => {
  //       currentBalance += t.transaction_type === "income" ? t.amount : -t.amount
  //     })
  //     trend.push({ date: month, balance: Number.parseFloat(currentBalance.toFixed(2)) })
  //   })

  //   // If no transactions in a month within the range, we might want to carry forward the balance.
  //   // This current logic shows balance at month-end where transactions occurred.
  //   return trend
  // }, [transactions, initialTransactions])

  const handleAddTransaction = useCallback((values: Omit<Transaction, "id">) => {
    // setTransactions((prev) =>
    //   [...prev, { ...values, id: crypto.randomUUID() }].sort(
    //     (a, b) => new Date(a.transaction_date).getTime() - new Date(b.transaction_date).getTime(), // Sort by date asc for consistency
    //   ),
    // )
    setIsAddDialogOpen(false)
  }, [])

  const handleEditTransaction = useCallback(
    (values: Omit<Transaction, "id">) => {
      if (!currentTransaction) return
      // setTransactions((prev) =>
      //   prev
      //     .map((t) => (t.id === currentTransaction.id ? { ...values, id: t.id } : t))
      //     .sort((a, b) => new Date(a.transaction_date).getTime() - new Date(b.transaction_date).getTime()),
      // )
      setIsEditDialogOpen(false)
      setCurrentTransaction(undefined)
    },
    [currentTransaction],
  )

  const handleDeleteTransaction = useCallback(() => {
    if (!transactionToDeleteId) return
    setTransactions((prev) => prev.filter((t) => t.id !== transactionToDeleteId))
    setIsDeleteDialogOpen(false)
    setTransactionToDeleteId(null)
  }, [transactionToDeleteId])

  const openEditDialog = (transaction: Transaction) => {
    setCurrentTransaction(transaction)
    setIsEditDialogOpen(true)
  }

  const openDeleteDialog = (transactionId: string) => {
    setTransactionToDeleteId(transactionId)
    setIsDeleteDialogOpen(true)
  }

  const resetFilters = () => {
    const now = new Date();
    const awalBulan = new Date(now.getFullYear(), now.getMonth(), 1);

    setDateRange({ from: awalBulan, to: now });
    setSelectedType("all")
    setSelectedCategory("all")
  }

  return (
    <div className="container mx-auto p-4 md:p-8 space-y-8 bg-background text-foreground min-h-screen">
      <header className="flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0">
        <h1 className="text-3xl font-bold tracking-tight">Duite Dashboard</h1>
        <Dialog open={isAddDialogOpen} onOpenChange={setIsAddDialogOpen}>
          <DialogTrigger asChild>
            <Button>
              <PlusCircle className="mr-2 h-4 w-4" /> Buat Transaksi
            </Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Tambah Transaksi Baru</DialogTitle>
              <DialogDescription>Lengkapi rincian transaksi baru Anda.</DialogDescription>
            </DialogHeader>
            <TransactionForm onSubmit={handleAddTransaction} onCancel={() => setIsAddDialogOpen(false)} />
          </DialogContent>
        </Dialog>
      </header>

      <div className="flex items-center justify-between">
        <p className="text-muted-foreground text-sm">
          Menampilkan data untuk: 
          <strong>
          {dateRange.from && dateRange.to
            ? `${formatDate(dateRange.from)} - ${formatDate(dateRange.to)}`
            : "-"}
          </strong>
        </p>
      </div>

      <section className="grid gap-4 md:grid-cols-3">
        <SummaryCard title="Total Pemasukan" value={`Rp${summaryData.totalIncome.toLocaleString('id')}`} icon={TrendingUp} />
        <SummaryCard title="Total Pengeluaran" value={`Rp${summaryData.totalExpenses.toLocaleString('id')}`} icon={TrendingDown} />
        <SummaryCard
          title="Saldo"
          value={`Rp${summaryData.balance.toLocaleString('id')}`}
          icon={DollarSign}
          description={summaryData.balance >= 0 ? "Keuangan sehat!" : "Perlu perhatian!"}
          descriptionVariant={summaryData.balance >= 0 ? 'success' : 'danger'}
          />
      </section>

      {/* Charts Section */}
      <section className="grid grid-cols-1 gap-4 md:grid-cols-2 md:gap-6">
        {/* Chart 1 - Akan sejajar dengan Chart 2 di desktop */}
        <div className="h-full min-h-[300px]">
          <FinancialChart data={monthlyIncomeExpenseData} />
        </div>
        
        {/* Chart 2 - Akan sejajar dengan Chart 1 di desktop */}
        <div className="h-full min-h-[300px]">
          <ExpenseCategoryPieChart data={expenseCategoryData} />
        </div>
        
        {/* Chart 3 - Selalu full width */}
        {/* <div className="h-full min-h-[300px] md:col-span-2">
          <BalanceTrendChart data={balanceTrendData} />
        </div> */}
      </section>

      <section>
        <h2 className="text-2xl font-semibold mb-4">Transactions</h2>
        <TransactionFilters
          dateRange={dateRange}
          onDateRangeChange={(range) => setDateRange(range ?? { from: undefined })}
          selectedType={selectedType}
          onSelectedTypeChange={setSelectedType}
          selectedCategory={selectedCategory}
          onSelectedCategoryChange={setSelectedCategory}
          categories={uniqueCategories}
          onResetFilters={resetFilters}
        />
        <TransactionsTable
          transactions={transactions}
          onEdit={openEditDialog}
          onDelete={openDeleteDialog}
        />
      </section>

      {/* Edit Transaction Dialog */}
      <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
        <DialogContent className="sm:max-w-[425px]">
          <DialogHeader>
            <DialogTitle>Edit Transaksi</DialogTitle>
            <DialogDescription>Perbarui detail transaksi Anda.</DialogDescription>
          </DialogHeader>
          {currentTransaction && (
            <TransactionForm
              onSubmit={handleEditTransaction}
              initialData={currentTransaction}
              onCancel={() => setIsEditDialogOpen(false)}
            />
          )}
        </DialogContent>
      </Dialog>

      {/* Delete Confirmation Dialog */}
      <Dialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle className="flex items-center">
              <AlertTriangle className="mr-2 h-6 w-6 text-red-500" />
              Konfirmasi Penghapusan
            </DialogTitle>
            <DialogDescription>
              Apakah Anda yakin ingin menghapus transaksi ini? Tindakan ini tidak dapat dibatalkan.
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-end space-x-2 pt-4">
            <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>
              Batal
            </Button>
            <Button variant="destructive" onClick={handleDeleteTransaction}>
              Hapus
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
