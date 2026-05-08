import { useState } from 'react'
import {Navigate, useNavigate} from 'react-router-dom'
import { useAuth } from '../context/useAuth'
import './pages.css'

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
    <div className="page-container">
      <div className="form-card">
        <h1>Вход</h1>
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
          <button className="submit-btn" type="submit">Войти</button>
        </form>
      </div>
    </div>
  )
}
