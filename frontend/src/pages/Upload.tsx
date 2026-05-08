import {useAuth} from "../context/useAuth.ts";
import {Navigate} from "react-router-dom";
import { useState, useEffect } from 'react'
import './pages.css'

const API_URL = '/api/v1/task'

export default function Upload() {
  const { isAuthenticated, tokens, checkTokenExpired, logout } = useAuth()
    const [files, setFiles] = useState<File[]>([])
    const [width, setWidth] = useState('')
    const [height, setHeight] = useState('')
    const [quality, setQuality] = useState('')
    const [format, setFormat] = useState('')
    const [uploading, setUploading] = useState(false)
    const [error, setError] = useState('')
    const [success, setSuccess] = useState('')

    const tokenExpired = checkTokenExpired()

    useEffect(() => {
      if (tokenExpired) {
        logout()
      }
    }, [tokenExpired, logout])

    if (tokenExpired || !isAuthenticated) {
        return <Navigate to="/" replace />
    }

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            const newFiles = Array.from(e.target.files).filter(f =>
                /\.(jpg|jpeg|png|webp)$/i.test(f.name)
            )
            setFiles(prev => [...prev, ...newFiles])
        }
    }

    const removeFile = (index: number) => {
        setFiles(prev => prev.filter((_, i) => i !== index))
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        if (files.length === 0) return

        setUploading(true)
        setError('')
        setSuccess('')

        const params = new URLSearchParams()
        if (width) params.set('width', width)
        if (height) params.set('height', height)
        if (quality) params.set('quality', quality)
        if (format) params.set('format', format)

        const url = `${API_URL}${params.toString() ? '?' + params.toString() : ''}`

        const formData = new FormData()
        files.forEach(file => formData.append('file', file))

        try {
            const res = await fetch(url, {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${tokens?.access_token}` },
                body: formData
            })
            if (!res.ok) throw new Error('Upload failed')
            const taskId = await res.text()
            setSuccess(`Files uploaded successfully. Task ID: ${taskId}`)
            setFiles([])
        } catch {
            setError('Upload failed')
        } finally {
            setUploading(false)
        }
    }

    return (
        <div className="page-container">
            <div className="form-card">
                <h1>Загрузка фото</h1>
                <form onSubmit={handleSubmit}>
                    <div className="form-group">
                        <label>Файлы</label>
                        <div className="file-input-wrapper">
                            <div className="file-input-btn">+ Выбрать файлы</div>
                            <input type="file" multiple accept=".jpg,.jpeg,.png,.webp" onChange={handleFileChange} />
                        </div>
                        <div className="file-list">
                            {files.map((f, i) => (
                                <div className="file-item" key={i}>
                                    <span>{f.name}</span>
                                    <button type="button" onClick={() => removeFile(i)}>×</button>
                                </div>
                            ))}
                        </div>
                    </div>
                    <div className="row">
                        <div className="form-group">
                            <label>Ширина</label>
                            <input type="number" value={width} onChange={e => setWidth(e.target.value)} />
                        </div>
                        <div className="form-group">
                            <label>Высота</label>
                            <input type="number" value={height} onChange={e => setHeight(e.target.value)} />
                        </div>
                    </div>
                    <div className="row">
                        <div className="form-group">
                            <label>Качество (0-1)</label>
                            <input type="number" step="0.1" min="0" max="1" value={quality} onChange={e => setQuality(e.target.value)} />
                        </div>
                        <div className="form-group">
                            <label>Формат</label>
                            <select value={format} onChange={e => setFormat(e.target.value)}>
                                <option value="">Original</option>
                                <option value="jpg">JPG</option>
                                <option value="png">PNG</option>
                                <option value="webp">WEBP</option>
                            </select>
                        </div>
                    </div>
                    <button className="submit-btn" type="submit" disabled={uploading || files.length === 0}>
                        {uploading ? 'Загрузка...' : 'Загрузить'}
                    </button>
                </form>
                {error && <p className="error-msg">{error}</p>}
                {success && <p className="success-msg">{success}</p>}
            </div>
        </div>
    )
}
