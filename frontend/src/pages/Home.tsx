import {Navigate} from "react-router-dom";
import {useAuth} from "../context/useAuth.ts";

export default function Home() {
  const { isAuthenticated } = useAuth()
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return <Navigate to="/upload" replace />
}
