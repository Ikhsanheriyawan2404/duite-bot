import { AuthResponse, User } from "@/lib/types"
import { saveTokens, getRefreshToken } from "../lib/authToken"

const API_BASE = "https://gateway-duite.brogrammer.id"

export async function magicLogin(token: string): Promise<void> {
  const res = await fetch(`${API_BASE}/auth/magic-login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ token }),
  })
  if (!res.ok) throw new Error("Login failed")
  const data: AuthResponse = await res.json()
  saveTokens(data.access_token, data.refresh_token)
}

export async function refreshToken(): Promise<void> {
  const refresh = getRefreshToken()
  if (!refresh) throw new Error("No refresh token")
  const res = await fetch(`${API_BASE}/auth/refresh`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ refresh_token: refresh }),
  })
  if (!res.ok) throw new Error("Refresh token failed")
  const data: AuthResponse = await res.json()
  saveTokens(data.access_token, data.refresh_token)
}

export async function getMe(): Promise<User> {
  const token = localStorage.getItem("access_token")
  const res = await fetch(`${API_BASE}/user/me`, {
    headers: {
      Authorization: `${token}`,
    },
  })

  if (res.status === 401) {
    await refreshToken()
    return getMe()
  }

  if (!res.ok) throw new Error("Failed to get user")
  return res.json()
}
