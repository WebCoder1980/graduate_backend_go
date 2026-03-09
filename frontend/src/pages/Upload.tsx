import Header from '../components/Header'
import {useAuth} from "../context/useAuth.ts";
import {Navigate} from "react-router-dom";

export default function Upload() {
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated) {
      return <Navigate to="/" replace />
  }

  return (
    <>
      <Header />
      <main>
        <h1>Загрузка фото</h1>
      </main>
    </>
  )
}
