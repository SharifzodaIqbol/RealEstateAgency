<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api/axios';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Message from 'primevue/message';

const username = ref('');
const email = ref('');
const password = ref('');
const router = useRouter();
// Используем объект для хранения сообщения и его типа (severity)
const message = ref({ text: '', severity: '' }); 

const handleRegister = async () => {
    // Очищаем предыдущее сообщение
    message.value = { text: '', severity: '' }; 
    
    // Простая валидация на фронтенде
    if (!username.value || !email.value || !password.value) {
        message.value = { text: 'Пожалуйста, заполните все поля.', severity: 'warn' };
        return;
    }
    
    try {
        // Отправляем запрос на ваш Go handler Register (POST /register)
        const response = await api.post('/register', { 
            username: username.value, // Требуется вашим бэкендом
            email: email.value, 
            password: password.value 
        });
        
        // Успешная регистрация
        const successMessage = response.data?.message || 'Вы успешно зарегистрировались!';
        message.value = { text: successMessage, severity: 'success' };
        
        // Очистка полей и перенаправление на логин через 2 секунды
        username.value = '';
        email.value = '';
        password.value = '';
        
        setTimeout(() => router.push('/login'), 2000);

    } catch (e) {
        // Обработка ошибки
        const errorMessage = e.response?.data || 'Ошибка регистрации. Возможно, пользователь уже существует.';
        message.value = { text: errorMessage, severity: 'error' }; 
    }
};
</script>

<template>
    <div class="flex-center">
        <div class="card login-box">
            <h2 style="text-align: center;">Регистрация пользователя</h2>
            
            <Message v-if="message.text" :severity="message.severity" class="mb-2">{{ message.text }}</Message>
            
            <div class="field"> 
                <label>Имя пользователя</label>
                <InputText v-model="username" class="w-full" />
            </div>
            
            <div class="field"> 
                <label>Email</label>
                <InputText v-model="email" class="w-full" />
            </div>
            
            <div class="field"> 
                <label>Пароль</label>
                <InputText v-model="password" type="password" class="w-full" />
            </div>

            <Button label="Зарегистрироваться" @click="handleRegister" class="w-full mt-2" />
            
            <p class="mt-3 text-center">
                Уже есть аккаунт? <router-link to="/login">Войти</router-link>
            </p>
        </div>
    </div>
</template>

<style scoped>
/* Стиль для центрирования блока формы на весь экран */
.flex-center {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh; 
}
/* Стиль для ограничения ширины формы */
.login-box {
    width: 100%;
    max-width: 400px; /* Фиксированная максимальная ширина */
}
/* Стандартный стиль карточки */
.card {
    background: var(--p-surface-card, #2a2a2a); 
    padding: 2rem;
    border-radius: 12px;
    box-shadow: 0 4px 10px rgba(0, 0, 0, 0.5); 
}
/* Стиль для каждого поля ввода */
.field { 
    margin-bottom: 1rem; 
    display: flex; 
    flex-direction: column; 
    gap: 5px; 
}

/* Вспомогательные классы */
.w-full { width: 100%; }
.mt-2 { margin-top: 1rem; }
.mt-3 { margin-top: 1.5rem; }
.mb-2 { margin-bottom: 1rem; }
.text-center { text-align: center; }
</style>