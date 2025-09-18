<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- 导航栏 -->
    <nav class="gradient-bg text-white shadow-lg">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center py-4">
          <div class="flex items-center">
            <i class="fas fa-shield-alt text-2xl mr-3"></i>
            <h1 class="text-xl font-bold">Yara安全服务管理系统</h1>
          </div>
          <div class="flex items-center space-x-4">
            <span
              class="status-indicator"
              :class="isSystemHealthy ? 'status-healthy' : 'status-error'"
            ></span>
            <span class="text-sm">{{
              isSystemHealthy ? "系统运行中" : "系统异常"
            }}</span>
            <button
              @click="refreshData"
              :disabled="isLoading"
              class="bg-white bg-opacity-20 px-3 py-1 rounded-lg hover:bg-opacity-30 transition-all disabled:opacity-50"
            >
              <i
                class="fas fa-sync-alt mr-1"
                :class="{ 'fa-spin': isLoading }"
              ></i>
              {{ isLoading ? "刷新中..." : "刷新" }}
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- 主要内容 -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <router-view />
    </main>

    <!-- 刷新成功提示 -->
    <div v-if="showRefreshSuccess" class="refresh-success show">
      <i class="fas fa-check-circle"></i>
      <span>刷新成功！</span>
    </div>

    <!-- 错误提示 -->
    <div
      v-if="error"
      class="fixed top-4 right-4 bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg z-50"
    >
      <div class="flex items-center">
        <i class="fas fa-exclamation-triangle mr-2"></i>
        <span>{{ error }}</span>
        <button @click="clearError" class="ml-2">
          <i class="fas fa-times"></i>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useAppStore } from "@/stores/app";

const appStore = useAppStore();

const showRefreshSuccess = ref(false);

// 计算属性
const isLoading = computed(() => appStore.isLoading);
const isSystemHealthy = computed(() => appStore.isSystemHealthy);
const error = computed(() => appStore.error);

// 方法
const refreshData = async () => {
  try {
    await appStore.loadDashboardData();
    showRefreshSuccess.value = true;
    setTimeout(() => {
      showRefreshSuccess.value = false;
    }, 2000);
  } catch (err) {
    console.error("刷新数据失败:", err);
  }
};

const clearError = () => {
  appStore.clearError();
};

// 生命周期
onMounted(async () => {
  try {
    await appStore.checkHealth();
    await appStore.loadDashboardData();
    appStore.addLog("系统启动完成", "info");
  } catch (err) {
    console.error("初始化失败:", err);
    appStore.addLog("系统初始化失败", "error");
  }
});
</script>

<style>
/* 引入Font Awesome */
@import url("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.0.0/css/all.min.css");
</style>
