import {useAuth} from "../context/useAuth.ts";
import {Navigate} from "react-router-dom";
import { useState, useEffect } from 'react'

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
        <>
            <main>
                <h1>Загрузка фото</h1>
                <form onSubmit={handleSubmit}>
                    <div>
                        <label>Files:</label>
                        <input type="file" multiple accept=".jpg,.jpeg,.png,.webp" onChange={handleFileChange} />
                        {files.map((f, i) => (
                            <div key={i}>
                                {f.name}
                                <button type="button" onClick={() => removeFile(i)}>×</button>
                            </div>
                        ))}
                    </div>
                    <div>
                        <label>Width:</label>
                        <input type="number" value={width} onChange={e => setWidth(e.target.value)} />
                    </div>
                    <div>
                        <label>Height:</label>
                        <input type="number" value={height} onChange={e => setHeight(e.target.value)} />
                    </div>
                    <div>
                        <label>Quality (0-1):</label>
                        <input type="number" step="0.1" min="0" max="1" value={quality} onChange={e => setQuality(e.target.value)} />
                    </div>
                    <div>
                        <label>Format:</label>
                        <select value={format} onChange={e => setFormat(e.target.value)}>
                            <option value="">Original</option>
                            <option value="jpg">JPG</option>
                            <option value="png">PNG</option>
                            <option value="webp">WEBP</option>
                        </select>
                    </div>
                    <button type="submit" disabled={uploading || files.length === 0}>
                        {uploading ? 'Uploading...' : 'Upload'}
                    </button>
                </form>
                {error && <p style={{color: 'red'}}>{error}</p>}
                {success && <p style={{color: 'green'}}>{success}</p>}
            </main>
        </>
    )
}
