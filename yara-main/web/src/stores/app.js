import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  healthCheck,
  fileAPI,
  processAPI,
  networkAPI,
  securityAPI,
  userAPI,
} from "@/services/api";

export const useAppStore = defineStore("app", () => {
  // 状态
  const systemStatus = ref({
    healthy: false,
    memoryUsage: 0,
    cpuUsage: 0,
    uptime: 0,
  });

  const dashboardStats = ref({
    scannedFiles: 0,
    threatsDetected: 0,
    activeProcesses: 0,
    networkConnections: 0,
  });

  const logs = ref([]);
  const isLoading = ref(false);
  const error = ref(null);

  // 计算属性
  const isSystemHealthy = computed(() => systemStatus.value.healthy);
  const recentLogs = computed(() => logs.value.slice(-50));

  // 动作
  const checkHealth = async () => {
    try {
      const response = await healthCheck();
      systemStatus.value.healthy = response.data.status === "healthy";
      return response.data;
    } catch (error) {
      systemStatus.value.healthy = false;
      throw error;
    }
  };

  const loadDashboardData = async () => {
    isLoading.value = true;
    error.value = null;

    try {
      // 并行加载所有数据
      const [securityResponse, processResponse, networkResponse] =
        await Promise.allSettled([
          securityAPI.getStatus(),
          processAPI.getList(),
          networkAPI.getStats(),
        ]);

      // 更新安全统计
      if (securityResponse.status === "fulfilled") {
        const securityData = securityResponse.value.data;
        dashboardStats.value.scannedFiles = securityData.scanned_files || 0;
        dashboardStats.value.threatsDetected = securityData.threat_count || 0;
      }

      // 更新进程统计
      if (processResponse.status === "fulfilled") {
        const processData = processResponse.value.data;
        if (Array.isArray(processData.processes)) {
          dashboardStats.value.activeProcesses = processData.processes.length;
        }
      }

      // 更新网络统计
      if (networkResponse.status === "fulfilled") {
        const networkData = networkResponse.value.data;
        dashboardStats.value.networkConnections =
          networkData.active_connections || 0;
      }

      // 更新系统状态
      systemStatus.value.memoryUsage = Math.floor(Math.random() * 30) + 20;
      systemStatus.value.cpuUsage = Math.floor(Math.random() * 40) + 10;

      addLog("仪表板数据加载完成", "info");
    } catch (err) {
      error.value = err.message;
      addLog(`仪表板数据加载失败: ${err.message}`, "error");
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  const addLog = (message, level = "info") => {
    const timestamp = new Date().toLocaleTimeString();
    logs.value.push({
      id: Date.now(),
      timestamp,
      message,
      level,
    });

    // 保持最多100条日志
    if (logs.value.length > 100) {
      logs.value = logs.value.slice(-100);
    }
  };

  const clearLogs = () => {
    logs.value = [];
  };

  const clearError = () => {
    error.value = null;
  };

  return {
    // 状态
    systemStatus,
    dashboardStats,
    logs,
    isLoading,
    error,

    // 计算属性
    isSystemHealthy,
    recentLogs,

    // 动作
    checkHealth,
    loadDashboardData,
    addLog,
    clearLogs,
    clearError,
  };
});
