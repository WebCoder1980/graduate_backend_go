import { useEffect } from 'react'
import Header from '../components/Header'
import { useAuth } from '../context/useAuth'
import {Navigate} from "react-router-dom";

export default function Admin() {
    const { isAuthenticated, checkTokenExpired, logout } = useAuth()

    const tokenExpired = checkTokenExpired()

    useEffect(() => {
      if (tokenExpired) {
        logout()
      }
    }, [tokenExpired, logout])

    if (tokenExpired || !isAuthenticated) {
        return <Navigate to="/" replace />
    }

    return (
    <>
      <Header />
      <main>
        <h1>Панель администратора</h1>
      </main>
    </>
  )
}
