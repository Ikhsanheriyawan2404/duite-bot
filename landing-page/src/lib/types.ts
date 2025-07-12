export type TransactionType = "INCOME" | "EXPENSE"

export interface Transaction {
  id: string
  transaction_date: string
  original_text: string
  category_id: number
  category: string
  amount: number
  transaction_type: TransactionType
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

export interface AuthResponse {
  access_token: string
  refresh_token: string
}

export interface User {
  id: number
  name: string
  chat_id: number
  uuid: string,
}

export interface Category {
  id: number
  name: string
  type: string
  created_at: string
  updated_at: string
}
