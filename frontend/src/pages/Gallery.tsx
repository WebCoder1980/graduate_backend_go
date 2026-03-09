import { Navigate } from 'react-router-dom'
import Header from '../components/Header'
import {useAuth} from "../context/useAuth.ts";

export default function Gallery() {
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated) {
    return <Navigate to="/" replace />
  }

  return (
    <>
      <Header />
      <main>
        <h1>Галерея</h1>
      </main>
    </>
  )
}
