/* eslint-disable react-refresh/only-export-components */
import {jwtDecode, type JwtPayload} from 'jwt-decode'
import { createContext, useState, type ReactNode } from 'react'

interface AuthTokens {
  access_token: string
  id_token: string
  expires_in: number
  refresh_expires_in: number
  refresh_token: string
  token_type: string
  'not-before-policy': number
  session_state: string
  scope: string
  roles: string[]
}

interface AuthContextType {
  tokens: AuthTokens | null
  login: (username: string, password: string) => Promise<void>
  register: (username: string, password: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
  checkTokenExpired: () => boolean
}

interface KeycloakPayload {
  realm_access?: {
    roles?: string[]
  }
  resource_access?: {
    [key: string]: {
      roles?: string[]
    }
  }
}

export const AuthContext = createContext<AuthContextType | null>(null)

const API_URL = '/api/v1/user'

function getStoredTokens(): AuthTokens | null {
  const stored = localStorage.getItem('auth_tokens')
  if (!stored) return null
  try {
    const tokens: AuthTokens = JSON.parse(stored)
    const payload = jwtDecode(tokens.access_token) as JwtPayload
    if (payload.exp && payload.exp * 1000 < Date.now()) {
      localStorage.removeItem('auth_tokens')
      return null
    }
    return tokens
  } catch {
    localStorage.removeItem('auth_tokens')
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [tokens, setTokens] = useState<AuthTokens | null>(() => getStoredTokens())

  const login = async (username: string, password: string) => {
    const response = await fetch(`${API_URL}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    })
    
    if (!response.ok) {
      throw new Error('Login failed')
    }
    
    const data: AuthTokens = await response.json()
    const payload = jwtDecode(data.access_token) as JwtPayload & KeycloakPayload
    data.roles = payload.realm_access?.roles || []
    console.log(data.roles)

    setTokens(data)
    localStorage.setItem('auth_tokens', JSON.stringify(data))
  }

  const register = async (username: string, password: string) => {
    const response = await fetch(`${API_URL}/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    })
    
    if (!response.ok) {
      throw new Error('Registration failed')
    }
  }

  const logout = () => {
    setTokens(null)
    localStorage.removeItem('auth_tokens')
  }

  const checkTokenExpired = () => {
    if (!tokens) return false
    try {
      const { exp } = jwtDecode(tokens.access_token) as JwtPayload
      return exp ? exp * 1000 < Date.now() : false
    } catch {
      return true
    }
  }

  return (
    <AuthContext.Provider value={{ tokens, login, register, logout, isAuthenticated: !!tokens, checkTokenExpired }}>
      {children}
    </AuthContext.Provider>
  )
}
