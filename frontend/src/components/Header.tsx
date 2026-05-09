import { Link } from 'react-router-dom'
import { useAuth } from '../context/useAuth'
import './Header.css'

export default function Header() {
  const { isAuthenticated, logout, userRoles } = useAuth()

  return (
    <header className="header">
      <nav className="header-nav">
        {!isAuthenticated && <Link to="/login">Вход</Link>}
        {!isAuthenticated && <Link to="/register">Регистрация</Link>}
        {isAuthenticated && userRoles.includes('user') && <Link to="/upload">Загрузить</Link>}
        {isAuthenticated && userRoles.includes('user') &&  <Link to="/gallery">Галерея</Link>}
        {isAuthenticated && userRoles.includes('admin') && <Link to="/admin">Админ</Link>}
        {isAuthenticated && <button onClick={logout}>Выход</button>}
      </nav>
    </header>
  )
}
