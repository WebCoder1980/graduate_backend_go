import { Outlet } from 'react-router-dom'
import Header from './Header'

export default function Layout() {
  return (
    <>
      <Header />
      <div style={{ padding: '20px' }}>
        <Outlet />
      </div>
    </>
  )
}
