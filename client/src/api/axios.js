import axios from 'axios';
import router from '../router';

const api = axios.create({
    baseURL: 'http://localhost:3000', // Адрес вашего Go сервера
    headers: {
        'Content-Type': 'application/json'
    }
});

// Перехватчик: добавляем токен к каждому запросу
api.interceptors.request.use(config => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

// Перехватчик: если токен протух (401), выкидываем на логин
api.interceptors.response.use(response => response, error => {
    if (error.response && error.response.status === 401) {
        localStorage.removeItem('token');
        router.push('/login');
    }
    return Promise.reject(error);
});

export default api;