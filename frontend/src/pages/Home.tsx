import { useEffect } from 'react'
import {Navigate} from "react-router-dom";
import {useAuth} from "../context/useAuth.ts";

export default function Home() {
  const { isAuthenticated, checkTokenExpired, logout, userRoles } = useAuth()

  const tokenExpired = checkTokenExpired()

  useEffect(() => {
    if (tokenExpired) {
      logout()
    }
  }, [tokenExpired, logout])

  if (tokenExpired || !isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  if (userRoles.includes('admin')) {
    return <Navigate to="/admin" replace />
  }

  return <Navigate to="/upload" replace />
}
