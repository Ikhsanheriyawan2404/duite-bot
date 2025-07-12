import { Category } from "@/lib/types"
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

export async function fetchCategories(): Promise<Category[]> {
  const res = await fetchWithAuth(`/categories`)
  if (!res.ok) throw new Error("Failed to fetch transactions")

  return await res.json()
}
