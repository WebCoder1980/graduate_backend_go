import { useState } from 'react'
import {Navigate, useNavigate} from 'react-router-dom'
import { useAuth } from '../context/useAuth'
import './pages.css'

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
    <div className="page-container">
      <div className="form-card">
        <h1>Регистрация</h1>
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Имя пользователя</label>
            <input
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
          </div>
          <div className="form-group">
            <label>Пароль</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
          {error && <p className="error-msg">{error}</p>}
          {success && <p className="success-msg">Регистрация успешна! Перенаправление на вход...</p>}
          <button className="submit-btn" type="submit">Зарегистрироваться</button>
        </form>
      </div>
    </div>
  )
}
