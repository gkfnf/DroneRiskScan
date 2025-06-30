<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <router-view />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSystemStore } from '@/stores/system'

const systemStore = useSystemStore()

// 初始化系统
systemStore.initialize()

onMounted(() => {
  // Additional logic can be added here if needed
})
</script>

<style>
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

/* 自定义样式 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from {
  transform: translateX(-100%);
}

.slide-leave-to {
  transform: translateX(100%);
}
</style>
