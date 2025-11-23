<script setup>
import { useRouter } from 'vue-router';
import Menubar from 'primevue/menubar';

const router = useRouter();
// Получаем роль, если ее нет, считаем ее Пользователем (3), чтобы не показывать лишнего.
const roleID = parseInt(localStorage.getItem('role_id')) || 3; 

const menuItems = [
    { label: 'Недвижимость', icon: 'pi pi-home', to: '/properties' },
    
    // Покупки/Продажи доступны Агентам (2) и Админам (1)
    ...(roleID <= 2 ? [{ label: 'Покупки', icon: 'pi pi-wallet', to: '/purchases' }] : []),
    ...(roleID <= 2 ? [{ label: 'Продажи', icon: 'pi pi-money-bill', to: '/sales' }] : []),
    
    // Панель администратора доступна только Админу (Роль 1)
    ...(roleID === 1 ? [{ label: 'Админ-панель', icon: 'pi pi-shield', to: '/admin' }] : []),
];

const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('role_id');
    router.push('/login');
};

// Добавляем кнопку выхода в конец
const model = [
    ...menuItems,
    { 
        label: 'Выход', 
        icon: 'pi pi-sign-out', 
        command: logout,
        class: 'ml-auto' // Сдвигает кнопку вправо
    }
];
</script>

<template>
    <Menubar :model="model" />
    <div class="container">
        <router-view />
    </div>
</template>

<style scoped>
.container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
}
</style>