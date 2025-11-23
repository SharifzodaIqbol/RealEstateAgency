<script setup>
import { ref, onMounted } from 'vue';
import api from '../api/axios';
import { useRouter } from 'vue-router';

// Компоненты PrimeVue
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Select from 'primevue/select';

const properties = ref([]);
const router = useRouter();
const showDialog = ref(false);

// Модель для новой записи
const newProp = ref({
    address: '',
    type: '',
    price: null,
    status: 'active'
});

const types = ref(['Apartment', 'House', 'Studio', 'Office']);
const statuses = ref(['active', 'sold', 'reserved']);

// Загрузка данных
const loadData = async () => {
    try {
        const res = await api.get('/properties');
        properties.value = res.data || [];
    } catch (e) {
        console.error("Ошибка загрузки", e);
    }
};

// Создание записи
const createProperty = async () => {
    try {
        await api.post('/properties', {
            address: newProp.value.address,
            type: newProp.value.type,
            price: Number(newProp.value.price),
            status: newProp.value.status
        });
        showDialog.value = false;
        newProp.value = { address: '', type: '', price: null, status: 'active' };
        loadData();
    } catch (e) {
        alert("Ошибка при создании: " + (e.response?.data || e.message));
    }
};

const logout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('role_id');
    router.push('/login');
};

onMounted(loadData);
</script>

<template>
    <div class="container">
        <div class="header">
            <h1>Недвижимость</h1>
            <div>
                <Button label="Добавить" icon="pi pi-plus" @click="showDialog = true" class="mr-2" />
                <Button label="Выход" severity="secondary" icon="pi pi-sign-out" @click="logout" />
            </div>
        </div>

        <div class="card">
            <DataTable :value="properties" tableStyle="min-width: 50rem">
                <Column field="id" header="ID" sortable></Column>
                <Column field="address" header="Адрес" sortable></Column>
                <Column field="type" header="Тип" sortable></Column>
                <Column field="price" header="Цена ($)" sortable>
                    <template #body="slotProps">
                        {{ slotProps.data.price?.toLocaleString() }} $
                    </template>
                </Column>
                <Column field="status" header="Статус">
                    <template #body="slotProps">
                         <span :class="'status-badge status-' + slotProps.data.status">
                            {{ slotProps.data.status }}
                        </span>
                    </template>
                </Column>
            </DataTable>
        </div>

        <Dialog v-model:visible="showDialog" modal header="Добавить объект" :style="{ width: '25rem' }">
            <div class="field">
                <label for="addr" class="font-semibold w-6rem">Адрес</label>
                <InputText id="addr" v-model="newProp.address" class="flex-auto" autocomplete="off" />
            </div>
            <div class="field">
                <label for="type" class="font-semibold w-6rem">Тип</label>
                <Select v-model="newProp.type" :options="types" placeholder="Выберите тип" class="w-full" />
            </div>
            <div class="field">
                <label for="price" class="font-semibold w-6rem">Цена</label>
                <InputNumber id="price" v-model="newProp.price" inputId="currency-us" mode="currency" currency="USD" locale="en-US" class="w-full" />
            </div>
            <div class="field">
                <label for="status" class="font-semibold w-6rem">Статус</label>
                 <Select v-model="newProp.status" :options="statuses" class="w-full" />
            </div>
            
            <div class="flex justify-end gap-2 mt-4">
                <Button type="button" label="Отмена" severity="secondary" @click="showDialog = false"></Button>
                <Button type="button" label="Сохранить" @click="createProperty"></Button>
            </div>
        </Dialog>
    </div>
</template>

<style scoped>
.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
.mr-2 { margin-right: 0.5rem; }
.field { margin-bottom: 1rem; display: flex; flex-direction: column; gap: 5px; }
.w-full { width: 100%; }

.status-badge { padding: 4px 8px; border-radius: 4px; font-size: 0.9em; font-weight: bold; text-transform: uppercase; }
.status-active { background-color: #c8e6c9; color: #256029; }
.status-sold { background-color: #ffcdd2; color: #c63737; }
.status-reserved { background-color: #feedaf; color: #8a5340; }
</style>