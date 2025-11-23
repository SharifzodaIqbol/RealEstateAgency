<template>
    <div class="container">
        <div class="header">
            <h1>🛡️ Панель Администратора</h1>
            <Button label="Выход" severity="secondary" icon="pi pi-sign-out" @click="logout" />
        </div>
        
        <!-- Пользователи -->
        <div class="card">
            <h2>👥 Управление пользователями</h2>
            <DataTable :value="users" :loading="loading" emptyMessage="Нет данных о пользователях.">
                <Column field="id" header="ID"></Column>
                <Column field="username" header="Имя пользователя"></Column>
                <Column field="email" header="Email"></Column>
                <Column field="role_id" header="Роль">
                    <template #body="slotProps">
                        <span :class="'role-badge role-' + slotProps.data.role_id">
                            {{ getRoleName(slotProps.data.role_id) }}
                        </span>
                    </template>
                </Column>
                <Column header="Действия">
                    <template #body="slotProps">
                        <Button icon="pi pi-trash" severity="danger" 
                                @click="deleteUser(slotProps.data.id)" 
                                :disabled="slotProps.data.role_id === 1" />
                    </template>
                </Column>
            </DataTable>
        </div>

        <!-- Покупки -->
        <div class="card">
            <h2>💰 Все покупки</h2>
            <DataTable :value="purchases" :loading="loading" emptyMessage="Нет данных о покупках.">
                <Column field="id" header="ID"></Column>
                <Column field="property_id" header="ID Недвижимости"></Column>
                <Column field="initial_price" header="Начальная цена">
                    <template #body="slotProps">
                        {{ slotProps.data.initial_price?.toLocaleString() }} $
                    </template>
                </Column>
                <Column field="purchase_date" header="Дата покупки">
                    <template #body="slotProps">
                        {{ new Date(slotProps.data.purchase_date).toLocaleDateString() }}
                    </template>
                </Column>
                <Column header="Действия">
                    <template #body="slotProps">
                        <Button icon="pi pi-trash" severity="danger" 
                                @click="deletePurchase(slotProps.data.id)" />
                    </template>
                </Column>
            </DataTable>
        </div>

        <!-- Продажи -->
        <div class="card">
            <h2>🏪 Все продажи</h2>
            <DataTable :value="sales" :loading="loading" emptyMessage="Нет данных о продажах.">
                <Column field="id" header="ID"></Column>
                <Column field="property_id" header="ID Недвижимости"></Column>
                <Column field="final_price" header="Финальная цена">
                    <template #body="slotProps">
                        {{ slotProps.data.final_price?.toLocaleString() }} $
                    </template>
                </Column>
                <Column field="sale_date" header="Дата продажи">
                    <template #body="slotProps">
                        {{ new Date(slotProps.data.sale_date).toLocaleDateString() }}
                    </template>
                </Column>
                <Column header="Действия">
                    <template #body="slotProps">
                        <Button icon="pi pi-trash" severity="danger" 
                                @click="deleteSale(slotProps.data.id)" />
                    </template>
                </Column>
            </DataTable>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';
import { useRouter } from 'vue-router';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';

const users = ref([]);
const purchases = ref([]);
const sales = ref([]);
const loading = ref(false);
const router = useRouter();

const loadUsers = async () => {
    try {
        const res = await api.get('/admin/users');
        users.value = res.data || [];
    } catch (e) {
        console.error("Ошибка загрузки пользователей:", e);
        alert("Ошибка загрузки пользователей: " + (e.response?.data || e.message));
    }
};

const loadPurchases = async () => {
    try {
        const res = await api.get('/purchases');
        purchases.value = res.data || [];
    } catch (e) {
        console.error("Ошибка загрузки покупок:", e);
    }
};

const loadSales = async () => {
    try {
        const res = await api.get('/sales');
        sales.value = res.data || [];
    } catch (e) {
        console.error("Ошибка загрузки продаж:", e);
    }
};

const getRoleName = (roleId) => {
    const roles = { 1: 'Админ', 2: 'Агент', 3: 'Пользователь' };
    return roles[roleId] || 'Неизвестно';
};

const deleteUser = async (id) => {
    if (confirm('Вы уверены, что хотите удалить пользователя?')) {
        try {
            await api.delete(`/admin/users/${id}`);
            loadUsers();
        } catch (e) {
            alert("Ошибка удаления пользователя: " + (e.response?.data || e.message));
        }
    }
};

const deletePurchase = async (id) => {
    if (confirm('Вы уверены, что хотите удалить покупку?')) {
        try {
            await api.delete(`/admin/purchases/${id}`);
            loadPurchases();
        } catch (e) {
            alert("Ошибка удаления покупки: " + (e.response?.data || e.message));
        }
    }
};

const deleteSale = async (id) => {
    if (confirm('Вы уверены, что хотите удалить продажу?')) {
        try {
            await api.delete(`/admin/sales/${id}`);
            loadSales();
        } catch (e) {
            alert("Ошибка удаления продажи: " + (e.response?.data || e.message));
        }
    }
};

const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('role_id');
    router.push('/login');
};

onMounted(() => {
    loadUsers();
    loadPurchases();
    loadSales();
});
</script>

<style scoped>
.container { max-width: 1200px; margin: 0 auto; padding: 20px; }
.card { margin-top: 20px; padding: 1.5rem; background: var(--p-surface-card); border-radius: 12px; }
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }

.role-badge { padding: 4px 8px; border-radius: 4px; font-size: 0.9em; font-weight: bold; }
.role-1 { background-color: #ffcdd2; color: #c63737; }
.role-2 { background-color: #ffd8b2; color: #805b36; }
.role-3 { background-color: #c8e6c9; color: #256029; }
</style>