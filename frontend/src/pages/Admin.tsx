import Header from '../components/Header'
import { useAuth } from '../context/useAuth'
import {Navigate} from "react-router-dom";

export default function Admin() {
    const { isAuthenticated } = useAuth()
    if (!isAuthenticated) {
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
