import { Link } from 'react-router-dom'
import './Header.css'

export default function Header() {
  return (
    <header className="header">
      <nav className="header-nav">
        <Link to="/">Главная</Link>
        <Link to="/login">Вход</Link>
        <Link to="/register">Регистрация</Link>
        <Link to="/upload">Загрузить</Link>
        <Link to="/gallery">Галерея</Link>
        <Link to="/admin">Админ</Link>
      </nav>
    </header>
  )
}
