import { useState, useEffect, useCallback } from 'react'
import { Navigate } from 'react-router-dom'
import { useAuth } from "../context/useAuth.ts"

interface Task {
  id: number
  created_dt: string
  width: string | null
  height: string | null
  format: string | null
  quality: number | null
  user_uuid: string
  common_status_id?: number
  images?: Image[]
}

interface Image {
  id: number
  name: string
  format: string
  task_id: number
  position: number
  status_id: number
  end_dt: string
}

const API_URL = '/api/v1'

export default function Gallery() {
  const { isAuthenticated, tokens, checkTokenExpired, logout, userRoles } = useAuth()
  const [tasks, setTasks] = useState<Task[]>([])
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [imageUrls, setImageUrls] = useState<Record<number, string>>({})
  const [selectedImage, setSelectedImage] = useState<Image | null>(null)

  const getAuthHeaders = () => ({
    'Authorization': `Bearer ${tokens?.access_token}`
  })

  const fetchImage = useCallback(async (imageId: number) => {
    if (imageUrls[imageId]) return
    try {
      const res = await fetch(`${API_URL}/image-processor/${imageId}`, {
        headers: getAuthHeaders()
      })
      if (!res.ok) throw new Error('Failed to fetch image')
      const blob = await res.blob()
      const url = URL.createObjectURL(blob)
      setImageUrls(prev => ({ ...prev, [imageId]: url }))
    } catch (e) {
      console.error('Error loading image:', e)
    }
  }, [tokens, imageUrls])

  const fetchTasks = useCallback(async () => {
    setLoading(true)
    setError('')
    try {
      const res = await fetch(`${API_URL}/task`, {
        headers: getAuthHeaders()
      })
      if (!res.ok) throw new Error('Failed to fetch tasks')
      const data: Task[] = await res.json()
      setTasks(data)
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }, [tokens])

  const fetchTaskDetails = useCallback(async (taskId: number) => {
    setLoading(true)
    setError('')
    try {
      const res = await fetch(`${API_URL}/task/${taskId}`, {
        headers: getAuthHeaders()
      })
      if (!res.ok) throw new Error('Failed to fetch task details')
      const data: Task = await res.json()
      setSelectedTask(data)
      
      if (data.images) {
        data.images.forEach(img => fetchImage(img.id))
      }
    } catch (e) {
      setError(e instanceof Error ? e.message : 'Unknown error')
    } finally {
      setLoading(false)
    }
  }, [tokens, fetchImage])

  useEffect(() => {
    return () => {
      Object.values(imageUrls).forEach(url => URL.revokeObjectURL(url))
    }
  }, [])

  useEffect(() => {
    if (isAuthenticated && tokens) {
      fetchTasks()
    }
  }, [isAuthenticated, tokens, fetchTasks])

  const tokenExpired = checkTokenExpired()

  useEffect(() => {
    if (tokenExpired) {
      logout()
    }
  }, [tokenExpired, logout])

  if (tokenExpired || !isAuthenticated || !userRoles.includes('user')) {
    return <Navigate to="/" replace />
  }

  return (
    <>
      <main>
        <h1>Галерея</h1>
        {error && <p style={{ color: 'red' }}>{error}</p>}
        {loading && <p>Загрузка...</p>}
        <div style={{ display: 'flex', gap: '20px' }}>
          <div style={{ flex: 1 }}>
            <h2>Задачи</h2>
            <ul>
              {tasks.map(task => (
                <li key={task.id} onClick={() => fetchTaskDetails(task.id)} style={{ cursor: 'pointer', padding: '8px', borderBottom: '1px solid #ccc' }}>
                  Задача #{task.id} - {new Date(task.created_dt).toLocaleString()}
                </li>
              ))}
            </ul>
          </div>
          {selectedTask && (
            <div style={{ flex: 2 }}>
              <h2>Детали задачи #{selectedTask.id}</h2>
              {(() => {
                return (
                  <>
                    <p>Размеры: {selectedTask.width ?? 'исходное'} x {selectedTask.height ?? 'исходное'}</p>
                    <p>Формат: {selectedTask.format ?? 'исходный'}</p>
                    <p>Качество: {selectedTask.quality ?? 'исходное'}</p>
                  </>
                )

              })(                )}
                {selectedTask.images && selectedTask.images.length > 0 && (
                  <button
                    onClick={() => {
                      selectedTask.images?.forEach(img => {
                        if (imageUrls[img.id]) {
                          const a = document.createElement('a')
                          a.href = imageUrls[img.id]
                          a.download = `${img.name}.${img.format}`
                          document.body.appendChild(a)
                          a.click()
                          document.body.removeChild(a)
                        }
                      })
                    }}
                    style={{ marginBottom: '15px', padding: '8px 16px', cursor: 'pointer', backgroundColor: '#4CAF50', color: 'white', border: 'none', borderRadius: '4px' }}
                  >
                    Скачать все изображения
                  </button>
                )}
                <div>
                  <h3>Изображения</h3>
                  {selectedTask.images?.map(img => (
                    <div key={img.id} style={{ marginBottom: '10px' }}>
                      <p>{img.name}.{img.format} (позиция: {img.position})</p>
                      {imageUrls[img.id] ? (
                        <>
                          <img 
                            src={imageUrls[img.id]}
                            alt={img.name}
                            style={{ maxWidth: '300px', maxHeight: '300px', cursor: 'pointer', display: 'block' }}
                            onClick={() => setSelectedImage(img)}
                          />
                          <button
                            onClick={() => {
                              const a = document.createElement('a')
                              a.href = imageUrls[img.id]
                              a.download = `${img.name}.${img.format}`
                              a.click()
                            }}
                            style={{ marginTop: '5px', padding: '5px 10px', cursor: 'pointer' }}
                          >
                            Скачать
                          </button>
                        </>
                      ) : (
                        <p>Загрузка изображения...</p>
                      )}
                    </div>
                  ))}
                </div>
            </div>
          )}
        </div>
      </main>
      {selectedImage && imageUrls[selectedImage.id] && (
        <div 
          style={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            backgroundColor: 'rgba(0,0,0,0.8)',
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center',
            zIndex: 1000
          }}
          onClick={() => setSelectedImage(null)}
        >
          <div style={{ position: 'relative' }}>
            <button 
              onClick={() => setSelectedImage(null)}
              style={{
                position: 'absolute',
                top: '-30px',
                right: '0',
                background: 'white',
                border: 'none',
                borderRadius: '50%',
                width: '30px',
                height: '30px',
                cursor: 'pointer'
              }}
            >
              ✕
            </button>
            <img 
              src={imageUrls[selectedImage.id]}
              alt={selectedImage.name}
              style={{ maxWidth: '90vw', maxHeight: '90vh' }}
            />
          </div>
        </div>
      )}
    </>
  )
}
