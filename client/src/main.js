import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import PrimeVue from 'primevue/config';
import Aura from '@primevue/themes/aura'; // Красивая современная тема
import 'primeicons/primeicons.css'; // Иконки
import './style.css' 

const app = createApp(App);

app.use(router);
app.use(PrimeVue, {
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: 'system', // Авто-тема (светлая/темная)
        }
    }
});

app.mount('#app');