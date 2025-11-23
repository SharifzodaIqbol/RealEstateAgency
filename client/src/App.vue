<template>
  <div id="app">
    <!-- Навигационное меню ТОЛЬКО для авторизованных пользователей -->
    <nav v-if="showNavigation" class="main-nav">
      <div class="nav-container">
        <div class="nav-links">
          <router-link to="/properties" class="nav-link">Недвижимость</router-link>
          <router-link 
            v-if="userRole <= 2" 
            to="/purchases" 
            class="nav-link"
          >
            Покупки
          </router-link>
          <router-link 
            v-if="userRole <= 2" 
            to="/sales" 
            class="nav-link"
          >
            Продажи
          </router-link>
          <router-link 
            v-if="userRole === 1" 
            to="/admin" 
            class="nav-link"
          >
            Админ панель
          </router-link>
        </div>
        <div class="nav-user">
          <span class="user-role">Роль: {{ roleName }}</span>
          <button @click="logout" class="logout-btn">Выйти</button>
        </div>
      </div>
    </nav>

    <!-- Основной контент -->
    <main :class="{ 'with-nav': showNavigation }">
      <router-view />
    </main>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';

export default {
  name: 'App',
  setup() {
    const router = useRouter();
    const route = useRoute();
    
    const userRole = ref(parseInt(localStorage.getItem('role_id')) || 3);
    
    // Вычисляемое свойство для проверки авторизации
    const isAuthenticated = computed(() => {
      return !!localStorage.getItem('token');
    });
    
    // Вычисляемое свойство для показа навигации
    const showNavigation = computed(() => {
      // Не показывать навигацию на страницах логина и регистрации
      const publicPages = ['/login', '/register'];
      return isAuthenticated.value && !publicPages.includes(route.path);
    });
    
    const roleName = computed(() => {
      const roles = { 1: 'Админ', 2: 'Агент', 3: 'Пользователь' };
      return roles[userRole.value] || 'Пользователь';
    });
    
    const logout = () => {
      localStorage.removeItem('token');
      localStorage.removeItem('role_id');
      userRole.value = 3;
      router.push('/login');
    };
    
    // Следим за изменением маршрута и обновляем роль
    watch(route, () => {
      userRole.value = parseInt(localStorage.getItem('role_id')) || 3;
    });
    
    onMounted(() => {
      userRole.value = parseInt(localStorage.getItem('role_id')) || 3;
    });
    
    return {
      userRole,
      isAuthenticated,
      showNavigation,
      roleName,
      logout
    };
  }
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.main-nav {
  background: #2c3e50;
  color: white;
  padding: 0;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.nav-container {
  max-width: 1200px;
  margin: 0 auto;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
}

.nav-links {
  display: flex;
  gap: 2rem;
}

.nav-link {
  color: white;
  text-decoration: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  transition: background-color 0.3s;
}

.nav-link:hover {
  background-color: #34495e;
}

.nav-link.router-link-active {
  background-color: #3498db;
}

.nav-user {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-role {
  color: #bdc3c7;
  font-size: 0.9rem;
}

.logout-btn {
  background: #e74c3c;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 4px;
  cursor: pointer;
}

.logout-btn:hover {
  background: #c0392b;
}

main.with-nav {
  padding-top: 0;
}
</style>