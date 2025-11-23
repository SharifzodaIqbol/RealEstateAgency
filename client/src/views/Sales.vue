<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';

const sales = ref([]);
const roleID = parseInt(localStorage.getItem('role_id') || '3');

const loadSales = async () => {
    try {
        const endpoint = roleID === 1 ? '/sales' : '/sales/my';
        const res = await api.get(endpoint);
        sales.value = res.data || [];
    } catch (e) {
        console.error("Ошибка загрузки продаж:", e);
        alert("Ошибка загрузки продаж: " + (e.response?.data || e.message));
    }
};

onMounted(loadSales);
</script>

<template>
    <div class="container">
        <h1>Список продаж</h1>
        <div class="card">
            <DataTable :value="sales" emptyMessage="Нет данных о продажах.">
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
            </DataTable>
        </div>
    </div>
</template>