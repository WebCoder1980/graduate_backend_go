import { useState } from 'react'
import {Navigate, useNavigate} from 'react-router-dom'
import Header from '../components/Header'
import { useAuth } from '../context/useAuth'

export default function Login() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login, isAuthenticated } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    if (isAuthenticated) {
      return <Navigate to="/" replace />
    }

    e.preventDefault()
    setError('')
    
    try {
      await login(username, password)
      navigate('/')
    } catch (error) {
      if (error instanceof Error) {
        setError(error.message)
      } else {
        setError(String(error))
      }
    }
  }

  return (
    <>
      <Header />
      <main>
        <h1>Вход</h1>
        <form onSubmit={handleSubmit}>
          <div>
            <label>
              Имя пользователя:
              <input
                type="text"
                value={username}
                onChange={(e) => setUsername(e.target.value)}
                required
              />
            </label>
          </div>
          <div>
            <label>
              Пароль:
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </label>
          </div>
          {error && <p style={{ color: 'red' }}>{error}</p>}
          <button type="submit">Войти</button>
        </form>
      </main>
    </>
  )
}
