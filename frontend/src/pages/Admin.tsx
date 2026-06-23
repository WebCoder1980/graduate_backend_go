import { useEffect, useState, forwardRef } from 'react'
import { useAuth } from '../context/useAuth'
import { Navigate } from "react-router-dom";
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';
import { Bar } from 'react-chartjs-2';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

const BarChart = forwardRef<any>(function BarChart(props, ref) {
    return <Bar ref={ref} {...props} />
})

const API_URL = '/api/v1/task/stats'

interface StatsData {
    Date: string;
    StatusId: number;
    StatusName: string;
    Count: number;
}

export default function Admin() {
    const { isAuthenticated, checkTokenExpired, logout, userRoles, tokens } = useAuth()
    const [statsData, setStatsData] = useState<StatsData[]>([])
    const [loading, setLoading] = useState(true)

    const tokenExpired = checkTokenExpired()

    useEffect(() => {
        if (tokenExpired) {
            logout()
        }
    }, [tokenExpired, logout])

    useEffect(() => {
        const fetchStats = async () => {
            if (!tokens) return;
            try {
                const response = await fetch(`${API_URL}?period_days=7`, {
                    headers: {
                        'Authorization': `Bearer ${tokens.access_token}`,
                        'Content-Type': 'application/json'
                    }
                });
                if (response.ok) {
                    const data: StatsData[] = await response.json();
                    setStatsData(data);
                }
            } catch (error) {
                console.error('Failed to fetch stats:', error);
            } finally {
                setLoading(false);
            }
        };
        fetchStats();
    }, [tokens]);

    if (tokenExpired || !isAuthenticated || !userRoles.includes('admin')) {
        return <Navigate to="/" replace />
    }

    if (loading) {
        return <main><h1>Панель администратора</h1></main>;
    }

    const uniqueDates = [...new Set(statsData.map(item => item.Date.split('T')[0]))].sort();
    const data = {
        labels: uniqueDates,
        datasets: [
            {
                label: 'Успех',
                data: uniqueDates.map(date => {
                    return statsData
                        .filter(item => item.Date.startsWith(date) && item.StatusName === 'Успех')
                        .reduce((sum, item) => sum + item.Count, 0);
                }),
                backgroundColor: 'rgba(75, 192, 192, 0.5)',
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 1
            },
            {
                label: 'Ошибка',
                data: uniqueDates.map(date => {
                    return statsData
                        .filter(item => item.Date.startsWith(date) && item.StatusName === 'Ошибка')
                        .reduce((sum, item) => sum + item.Count, 0);
                }),
                backgroundColor: 'rgba(255, 99, 132, 0.5)',
                borderColor: 'rgba(255, 99, 132, 1)',
                borderWidth: 1
            }
        ]
    };

    const options = {
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
            title: {
                display: true,
                text: 'Статистика задач по дням'
            }
        },
        scales: {
            x: {
                stacked: true
            },
            y: {
                stacked: true
            }
        }
    };

    return (
        <>
          <main>
            <h1>Панель администратора</h1>
            <div style={{ width: '100%', maxWidth: '800px', margin: '20px auto', padding: '20px', border: '1px solid #ccc', borderRadius: '8px' }}>
              <BarChart data={data} options={options} />
            </div>
          </main>
        </>
    )
}
