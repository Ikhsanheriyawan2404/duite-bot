export type TransactionType = "income" | "expense"

export interface Transaction {
  id: string
  date: Date
  description: string
  category: string
  amount: number
  type: TransactionType
}

// For the monthly income/expense bar chart
export interface MonthlyData {
  month: string // e.g., "Jan 2024"
  income: number
  expense: number
}

// For expense category pie chart
export interface CategoryExpenseData {
  category: string
  totalExpense: number
  fill?: string // Optional: for chart color
}

// For balance trend line chart
export interface BalanceTrendData {
  date: string // e.g., "Jan 2024" or a specific date
  balance: number
}
