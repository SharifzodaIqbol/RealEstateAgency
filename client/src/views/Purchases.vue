<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';

const purchases = ref([]);
const roleID = parseInt(localStorage.getItem('role_id') || '3');

const loadPurchases = async () => {
    try {
        const endpoint = roleID === 1 ? '/purchases' : '/purchases/my';
        const res = await api.get(endpoint);
        purchases.value = res.data || [];
    } catch (e) {
        console.error("Ошибка загрузки покупок:", e);
        alert("Ошибка загрузки покупок: " + (e.response?.data || e.message));
    }
};

onMounted(loadPurchases);
</script>

<template>
    <div class="container">
        <h1>Список покупок</h1>
        <div class="card">
            <DataTable :value="purchases" emptyMessage="Нет данных о покупках.">
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
            </DataTable>
        </div>
    </div>
</template>