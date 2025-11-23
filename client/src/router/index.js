// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Properties from '../views/Properties.vue';
import Purchases from '../views/Purchases.vue';
import Sales from '../views/Sales.vue';
import AdminDashboard from '../views/AdminDashboard.vue';

const ROLE_ADMIN = 1;
const ROLE_AGENT = 2;
const ROLE_USER = 3;

const routes = [
    { path: '/login', component: Login },
    { path: '/register', component: Register },
    
    // Основные маршруты
    { path: '/properties', component: Properties, name: 'properties' },
    { path: '/', redirect: '/properties' }, 
    
    // Маршруты для агентов и админов
    { 
        path: '/purchases', 
        component: Purchases, 
        name: 'purchases',
        meta: { minRole: ROLE_AGENT } 
    },
    { 
        path: '/sales', 
        component: Sales, 
        name: 'sales',
        meta: { minRole: ROLE_AGENT } 
    },

    // Маршрут только для админов
    { 
        path: '/admin', 
        component: AdminDashboard, 
        name: 'admin',
        meta: { minRole: ROLE_ADMIN } 
    },
];

const router = createRouter({
    history: createWebHistory(),
    routes
});

router.beforeEach((to, from, next) => {
    const publicPages = ['/login', '/register'];
    const authRequired = !publicPages.includes(to.path);
    const token = localStorage.getItem('token');
    
    // Получаем роль из localStorage
    const roleID = parseInt(localStorage.getItem('role_id')) || ROLE_USER;

    // Проверка авторизации
    if (authRequired && !token) {
        return next('/login');
    }

    // Проверка прав доступа
    if (to.meta?.minRole && roleID > to.meta.minRole) {
        alert('Недостаточно прав для доступа к этой странице');
        return next('/properties');
    }

    next();
});

export default router;