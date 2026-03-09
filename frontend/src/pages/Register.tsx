import { useState } from 'react'
import {Navigate, useNavigate} from 'react-router-dom'
import Header from '../components/Header'
import { useAuth } from '../context/useAuth'

export default function Register() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [success, setSuccess] = useState(false)
  const navigate = useNavigate()
  const { register, isAuthenticated } = useAuth()

  const handleSubmit = async (e: React.FormEvent) => {
    if (isAuthenticated) {
      return <Navigate to="/" replace />
    }

    e.preventDefault()
    setError('')
    setSuccess(false)
    
    try {
      await register(username, password)
      setSuccess(true)
      setTimeout(() => navigate('/login'), 1500)
    } catch {
      setError('Ошибка регистрации. Попробуйте другое имя пользователя.')
    }
  }

  return (
    <>
      <Header />
      <main>
        <h1>Регистрация</h1>
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
          {success && <p style={{ color: 'green' }}>Регистрация успешна! Перенаправление на вход...</p>}
          <button type="submit">Зарегистрироваться</button>
        </form>
      </main>
    </>
  )
}
