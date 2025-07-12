import { Transaction } from "@/lib/types"
import { getAccessToken } from "../lib/authToken"
import { refreshToken } from "./auth"

const API_BASE = import.meta.env.VITE_GATEWAY_API_URL

async function fetchWithAuth(input: string, init?: RequestInit): Promise<Response> {
  const token = getAccessToken()
  let res = await fetch(`${API_BASE}${input}`, {
    ...init,
    headers: {
      ...(init?.headers || {}),
      Authorization: `${token}`,
      "Content-Type": "application/json",
    },
  })

  if (res.status === 401) {
    await refreshToken()
    const newToken = getAccessToken()
    res = await fetch(`${API_BASE}${input}`, {
      ...init,
      headers: {
        ...(init?.headers || {}),
        Authorization: `${newToken}`,
        "Content-Type": "application/json",
      },
    })
  }

  return res
}

export async function fetchTransactions(
  type?: string,
  category_id?: string,
  start_date?: Date,
  end_date?: Date
): Promise<Transaction[]> {
  const params = new URLSearchParams()

  if (type) params.append("type", type)
  if (category_id) params.append("category_id", category_id)
  if (start_date) params.append("start_date", start_date.toISOString().slice(0, 10))
  if (end_date) params.append("end_date", end_date.toISOString().slice(0, 10))

  const res = await fetchWithAuth(`/transactions?${params.toString()}`)
  if (!res.ok) throw new Error("Failed to fetch transactions")

  return await res.json()
}

export async function createTransaction(data: Omit<Transaction, "id" | "category">): Promise<Transaction> {
  const res = await fetchWithAuth("/transactions", {
    method: "POST",
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    let message = "Failed to create transaction"
    try {
      const errorData = await res.json()
      if (errorData.error) {
        message = errorData.error
      }
    } catch {
      // JSON parse error fallback
    }
    throw new Error(message)
  }
  return res.json()
}

export async function updateTransaction(id: string, data: Omit<Transaction, "id" | "category">): Promise<Transaction> {
  const res = await fetchWithAuth(`/transactions/${id}`, {
    method: "PUT",
    body: JSON.stringify(data),
  })
  if (!res.ok) {
    let message = "Failed to create transaction"
    try {
      const errorData = await res.json()
      if (errorData.error) {
        message = errorData.error
      }
    } catch {
      // JSON parse error fallback
    }
    throw new Error(message)
  }
  return res.json()
}

export async function getTransactionById(id: string): Promise<Transaction> {
  const res = await fetchWithAuth(`/transactions/${id}`)
  if (!res.ok) throw new Error("Failed to fetch transaction detail")
  return res.json()
}

export async function deleteTransaction(id: string): Promise<void> {
  const res = await fetchWithAuth(`/transactions/${id}`, {
    method: "DELETE",
  })
  if (!res.ok) throw new Error("Failed to delete transaction")
}
