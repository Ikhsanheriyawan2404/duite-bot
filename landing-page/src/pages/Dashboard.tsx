"use client"

import { useState, useMemo, useCallback } from "react"
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
import { format, startOfMonth, endOfMonth } from "date-fns"
import { DateRange } from "react-day-picker"

// Mock initial data
const initialTransactions: Transaction[] = [
  {
    id: "1",
    date: new Date("2025-03-15"),
    description: "Old Salary",
    category: "Salary",
    amount: 4800,
    type: "income",
  },
  {
    id: "1a",
    date: new Date("2025-03-20"),
    description: "Old Groceries",
    category: "Groceries",
    amount: 120,
    type: "expense",
  },
  {
    id: "1b",
    date: new Date("2025-04-10"),
    description: "Consulting Gig",
    category: "Freelance",
    amount: 600,
    type: "income",
  },
  {
    id: "1c",
    date: new Date("2025-04-25"),
    description: "Utilities",
    category: "Utilities",
    amount: 90,
    type: "expense",
  },
  { id: "2", date: new Date("2025-05-01"), description: "Salary", category: "Salary", amount: 5000, type: "income" },
  {
    id: "3",
    date: new Date("2025-05-05"),
    description: "Groceries",
    category: "Groceries",
    amount: 150,
    type: "expense",
  },
  { id: "4", date: new Date("2025-05-10"), description: "Rent", category: "Rent", amount: 1200, type: "expense" },
  {
    id: "5",
    date: new Date("2025-06-01"),
    description: "Freelance Project",
    category: "Freelance",
    amount: 800,
    type: "income",
  },
  {
    id: "6",
    date: new Date("2025-06-03"),
    description: "Dinner Out",
    category: "Entertainment",
    amount: 75,
    type: "expense",
  },
  {
    id: "7",
    date: new Date("2025-06-15"),
    description: "New Phone",
    category: "Electronics",
    amount: 600,
    type: "expense",
  },
  { id: "8", date: new Date("2025-07-01"), description: "Salary", category: "Salary", amount: 5200, type: "income" },
  {
    id: "9",
    date: new Date("2025-07-05"),
    description: "Groceries",
    category: "Groceries",
    amount: 160,
    type: "expense",
  },
]

export default function DashboardPage() {
  const [transactions, setTransactions] = useState<Transaction[]>(initialTransactions)
  const [isAddDialogOpen, setIsAddDialogOpen] = useState(false)
  const [isEditDialogOpen, setIsEditDialogOpen] = useState(false)
  const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false)
  const [currentTransaction, setCurrentTransaction] = useState<Transaction | undefined>(undefined)
  const [transactionToDeleteId, setTransactionToDeleteId] = useState<string | null>(null)

  // const [dateRange, setDateRange] = useState<{ from: Date; to?: Date }>({from: new Date()})
  const [dateRange, setDateRange] = useState<DateRange>({
    from: undefined,
    to: undefined,
  });
  const [selectedType, setSelectedType] = useState<string>("all")
  const [selectedCategory, setSelectedCategory] = useState<string>("all")

  const uniqueCategories = useMemo(() => {
    const categories = new Set(transactions.map((t) => t.category))
    return ["all", ...Array.from(categories).sort()]
  }, [transactions])

  const filteredTransactions = useMemo(() => {
    return transactions
      .filter((transaction) => {
        const transactionDate = new Date(transaction.date)
        const { from, to } = dateRange
        if (from && transactionDate < startOfMonth(from)) return false
        if (to) {
          const toEndOfDay = endOfMonth(to) // Ensure 'to' date includes the whole month if filtering by month
          toEndOfDay.setHours(23, 59, 59, 999)
          if (transactionDate > toEndOfDay) return false
        }
        if (selectedType !== "all" && transaction.type !== selectedType) return false
        if (selectedCategory !== "all" && transaction.category !== selectedCategory) return false
        return true
      })
      .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()) // Sort for balance trend
  }, [transactions, dateRange, selectedType, selectedCategory])

  const summaryData = useMemo(() => {
    // Summary should reflect ALL transactions, not filtered ones
    const totalIncome = transactions.filter((t) => t.type === "income").reduce((sum, t) => sum + t.amount, 0)
    const totalExpenses = transactions.filter((t) => t.type === "expense").reduce((sum, t) => sum + t.amount, 0)
    const balance = totalIncome - totalExpenses
    return { totalIncome, totalExpenses, balance }
  }, [transactions])

  // Data for Monthly Income/Expense Bar Chart (FinancialChart)
  const monthlyIncomeExpenseData = useMemo((): MonthlyData[] => {
    const monthlyMap: Record<string, { income: number; expense: number }> = {}
    filteredTransactions.forEach((transaction) => {
      const date = new Date(transaction.date)
      const monthYear = format(date, "MMM yyyy")
      if (!monthlyMap[monthYear]) {
        monthlyMap[monthYear] = { income: 0, expense: 0 }
      }
      if (transaction.type === "income") {
        monthlyMap[monthYear].income += transaction.amount
      } else {
        monthlyMap[monthYear].expense += transaction.amount
      }
    })
    return Object.entries(monthlyMap)
      .map(([month, data]) => ({ month, ...data }))
      .sort((a, b) => new Date(a.month).getTime() - new Date(b.month).getTime())
  }, [filteredTransactions])

  // Data for Expense Category Pie Chart
  const expenseCategoryData = useMemo((): CategoryExpenseData[] => {
    const categoryMap: Record<string, number> = {}
    filteredTransactions
      .filter((t) => t.type === "expense")
      .forEach((transaction) => {
        categoryMap[transaction.category] = (categoryMap[transaction.category] || 0) + transaction.amount
      })
    return Object.entries(categoryMap)
      .map(([category, totalExpense]) => ({ category, totalExpense }))
      .sort((a, b) => b.totalExpense - a.totalExpense) // Sort for consistent pie chart segment order
  }, [filteredTransactions])

  // Data for Balance Trend Line Chart
  const balanceTrendData = useMemo((): BalanceTrendData[] => {
    if (filteredTransactions.length === 0) return []

    const trend: BalanceTrendData[] = []
    let currentBalance = 0

    // Calculate initial balance before the filtered period if filters are applied
    // This gives a more accurate starting point for the trend within the filtered view.
    // For simplicity in this example, we'll start the trend from the first filtered transaction.
    // A more robust solution would calculate balance from all transactions up to the start of the filtered period.

    const transactionsByMonth: Record<string, Transaction[]> = {}
    filteredTransactions.forEach((t) => {
      const monthYear = format(new Date(t.date), "MMM yyyy")
      if (!transactionsByMonth[monthYear]) {
        transactionsByMonth[monthYear] = []
      }
      transactionsByMonth[monthYear].push(t)
    })

    const sortedMonths = Object.keys(transactionsByMonth).sort((a, b) => new Date(a).getTime() - new Date(b).getTime())

    // Find overall starting balance before any filtered transactions
    // This is complex if we want the true historical balance.
    // For this example, let's assume the balance trend starts from 0 before the first transaction in the *original* dataset,
    // and then we calculate the balance up to the start of the filtered period.
    let balanceBeforeFilteredPeriod = 0
    const firstFilteredDate = filteredTransactions.length > 0 ? new Date(filteredTransactions[0].date) : new Date()

    initialTransactions.forEach((t) => {
      if (new Date(t.date) < firstFilteredDate) {
        balanceBeforeFilteredPeriod += t.type === "income" ? t.amount : -t.amount
      }
    })
    currentBalance = balanceBeforeFilteredPeriod

    sortedMonths.forEach((month) => {
      transactionsByMonth[month].forEach((t) => {
        currentBalance += t.type === "income" ? t.amount : -t.amount
      })
      trend.push({ date: month, balance: Number.parseFloat(currentBalance.toFixed(2)) })
    })

    // If no transactions in a month within the range, we might want to carry forward the balance.
    // This current logic shows balance at month-end where transactions occurred.
    return trend
  }, [filteredTransactions, initialTransactions])

  const handleAddTransaction = useCallback((values: Omit<Transaction, "id">) => {
    setTransactions((prev) =>
      [...prev, { ...values, id: crypto.randomUUID() }].sort(
        (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime(), // Sort by date asc for consistency
      ),
    )
    setIsAddDialogOpen(false)
  }, [])

  const handleEditTransaction = useCallback(
    (values: Omit<Transaction, "id">) => {
      if (!currentTransaction) return
      setTransactions((prev) =>
        prev
          .map((t) => (t.id === currentTransaction.id ? { ...values, id: t.id } : t))
          .sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime()),
      )
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
    setDateRange({from: new Date()})
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

      <section className="grid gap-4 md:grid-cols-3">
        <SummaryCard title="Total Pemasukan" value={`$${summaryData.totalIncome.toFixed(2)}`} icon={TrendingUp} />
        <SummaryCard title="Total Pengeluaran" value={`$${summaryData.totalExpenses.toFixed(2)}`} icon={TrendingDown} />
        <SummaryCard
          title="Saldo"
          value={`$${summaryData.balance.toFixed(2)}`}
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
        <div className="h-full min-h-[300px] md:col-span-2">
          <BalanceTrendChart data={balanceTrendData} />
        </div>
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
          transactions={filteredTransactions.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime())}
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
              Confirm Deletion
            </DialogTitle>
            <DialogDescription>
              Are you sure you want to delete this transaction? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>
          <div className="flex justify-end space-x-2 pt-4">
            <Button variant="outline" onClick={() => setIsDeleteDialogOpen(false)}>
              Cancel
            </Button>
            <Button variant="destructive" onClick={handleDeleteTransaction}>
              Delete
            </Button>
          </div>
        </DialogContent>
      </Dialog>
    </div>
  )
}
