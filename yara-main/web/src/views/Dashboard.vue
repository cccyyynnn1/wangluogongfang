<template>
  <div class="space-y-8">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <div class="bg-white rounded-lg shadow-md p-6 card-hover">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-blue-100 text-blue-600">
            <i class="fas fa-file-alt text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">已扫描文件</p>
            <p class="text-2xl font-bold text-gray-900">
              {{ dashboardStats.scannedFiles }}
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6 card-hover">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-red-100 text-red-600">
            <i class="fas fa-exclamation-triangle text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">威胁检测</p>
            <p class="text-2xl font-bold text-gray-900">
              {{ dashboardStats.threatsDetected }}
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6 card-hover">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-green-100 text-green-600">
            <i class="fas fa-cogs text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">活跃进程</p>
            <p class="text-2xl font-bold text-gray-900">
              {{ dashboardStats.activeProcesses }}
            </p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6 card-hover">
        <div class="flex items-center">
          <div class="p-3 rounded-full bg-purple-100 text-purple-600">
            <i class="fas fa-network-wired text-xl"></i>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">网络连接</p>
            <p class="text-2xl font-bold text-gray-900">
              {{ dashboardStats.networkConnections }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- 系统健康状态 -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">
        <i class="fas fa-heartbeat mr-2"></i>系统健康状态
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="flex items-center space-x-3">
          <div
            class="w-3 h-3 rounded-full"
            :class="
              healthStatus.status === 'healthy' ? 'bg-green-500' : 'bg-red-500'
            "
          ></div>
          <span class="text-sm font-medium"
            >API服务:
            {{ healthStatus.status === "healthy" ? "正常" : "异常" }}</span
          >
        </div>
        <div class="text-sm text-gray-600">
          响应时间: {{ healthStatus.response_time || "N/A" }}ms
        </div>
        <div class="text-sm text-gray-600">
          最后检查: {{ healthStatus.last_check || "N/A" }}
        </div>
      </div>
      <div class="mt-4 flex gap-3">
        <button
          @click="checkHealth"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          <i class="fas fa-sync-alt mr-2"></i>检查健康状态
        </button>
        <button
          @click="getMetrics"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
        >
          <i class="fas fa-chart-line mr-2"></i>获取性能指标
        </button>
        <button
          @click="resetMetrics"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          <i class="fas fa-undo mr-2"></i>重置性能指标
        </button>
      </div>
    </div>

    <!-- 主要功能区域 -->
    <div class="grid grid-cols-1 xl:grid-cols-4 gap-6">
      <!-- 文件扫描 - 占据左侧2列 -->
      <div class="xl:col-span-2">
        <div class="bg-white rounded-lg shadow-md p-6 h-full flex flex-col">
          <h2 class="text-lg font-semibold text-gray-900 mb-6">
            <i class="fas fa-search mr-2"></i>文件扫描
          </h2>

          <div class="space-y-6 flex-grow">
            <div class="space-y-3">
              <label class="block text-sm font-medium text-gray-700"
                >扫描类型</label
              >
              <select
                v-model="scanType"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
              >
                <option value="file">单文件扫描</option>
                <option value="directory">目录扫描</option>
                <option value="buffer">内存缓冲区扫描</option>
              </select>
            </div>

            <div v-if="scanType === 'file'" class="space-y-3">
              <label class="block text-sm font-medium text-gray-700"
                >文件路径</label
              >
              <input
                v-model="filePath"
                type="text"
                placeholder="输入文件路径，如: C:\test.exe"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
              />
            </div>

            <div v-if="scanType === 'directory'" class="space-y-3">
              <label class="block text-sm font-medium text-gray-700"
                >目录路径</label
              >
              <input
                v-model="directoryPath"
                type="text"
                placeholder="输入目录路径，如: C:\test"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
              />
            </div>

            <div v-if="scanType === 'buffer'" class="space-y-3">
              <label class="block text-sm font-medium text-gray-700"
                >缓冲区内容</label
              >
              <textarea
                v-model="bufferContent"
                placeholder="输入要扫描的内容"
                rows="5"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors resize-none"
              ></textarea>
            </div>

            <div class="pt-4">
              <button
                @click="performScan"
                :disabled="isScanning"
                class="w-full bg-blue-600 text-white py-3 px-6 rounded-lg hover:bg-blue-700 transition-colors font-medium text-base disabled:opacity-50"
              >
                <i
                  class="fas fa-search mr-2"
                  :class="{ 'fa-spin': isScanning }"
                ></i>
                {{ isScanning ? "扫描中..." : "开始扫描" }}
              </button>
            </div>
          </div>

          <div v-if="scanResult" class="mt-6">
            <h3 class="text-md font-semibold text-gray-900 mb-3">扫描结果</h3>
            <div
              class="bg-gray-50 p-4 rounded-lg text-sm border border-gray-200"
            >
              <div v-if="scanResult.success" class="space-y-3">
                <div class="flex justify-between items-center">
                  <span class="font-medium">扫描状态:</span>
                  <span
                    class="px-2 py-1 rounded text-sm"
                    :class="
                      scanResult.data.is_infected
                        ? 'bg-red-100 text-red-800'
                        : 'bg-green-100 text-green-800'
                    "
                  >
                    {{ scanResult.data.is_infected ? "发现威胁" : "安全" }}
                  </span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="font-medium">威胁数量:</span>
                  <span class="text-lg font-bold">{{
                    scanResult.data.threats?.length || 0
                  }}</span>
                </div>
                <div class="flex justify-between items-center">
                  <span class="font-medium">扫描时间:</span>
                  <span>{{ scanResult.data.scan_time || "N/A" }}</span>
                </div>
                <div
                  v-if="
                    scanResult.data.threats &&
                    scanResult.data.threats.length > 0
                  "
                  class="mt-4"
                >
                  <h4 class="font-medium text-red-600 mb-2">检测到的威胁:</h4>
                  <div class="space-y-2">
                    <div
                      v-for="threat in scanResult.data.threats"
                      :key="threat.rule_name"
                      class="bg-red-50 p-3 rounded border-l-4 border-red-400"
                    >
                      <div class="font-medium text-red-800">
                        {{ threat.rule_name }}
                      </div>
                      <div class="text-sm text-red-600">
                        {{ threat.description }}
                      </div>
                      <div class="text-xs text-red-500">
                        严重性: {{ threat.severity }} | 类别:
                        {{ threat.category }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div v-else class="text-red-600">
                扫描失败: {{ scanResult.message }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧信息面板 - 占据右侧2列 -->
      <div class="xl:col-span-2 flex flex-col space-y-6">
        <!-- 系统状态 -->
        <div class="bg-white rounded-lg shadow-md p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">
            <i class="fas fa-heartbeat mr-2"></i>系统状态
          </h2>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">服务状态</span>
                <div class="flex items-center">
                  <span class="status-indicator status-healthy mr-2"></span>
                  <span class="text-sm font-medium text-green-600">正常</span>
                </div>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">规则引擎</span>
                <div class="flex items-center">
                  <span class="status-indicator status-healthy mr-2"></span>
                  <span class="text-sm font-medium text-green-600">已加载</span>
                </div>
              </div>
            </div>
            <div class="space-y-3">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">内存使用</span>
                <span class="text-sm font-medium text-gray-900"
                  >{{ systemStatus.memoryUsage }}%</span
                >
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">CPU使用</span>
                <span class="text-sm font-medium text-gray-900"
                  >{{ systemStatus.cpuUsage }}%</span
                >
              </div>
            </div>
          </div>
        </div>

        <!-- 快速操作 -->
        <div class="bg-white rounded-lg shadow-md p-6 flex-grow">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">
            <i class="fas fa-bolt mr-2"></i>快速操作
          </h2>
          <div class="grid grid-cols-2 gap-3">
            <router-link
              to="/processes"
              class="text-left p-3 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors"
            >
              <i class="fas fa-list mr-2"></i>进程列表
            </router-link>
            <router-link
              to="/network"
              class="text-left p-3 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors"
            >
              <i class="fas fa-network-wired mr-2"></i>网络连接
            </router-link>
            <router-link
              to="/security"
              class="text-left p-3 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors"
            >
              <i class="fas fa-shield-alt mr-2"></i>安全规则
            </router-link>
            <router-link
              to="/users"
              class="text-left p-3 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors"
            >
              <i class="fas fa-user mr-2"></i>用户信息
            </router-link>
            <router-link
              to="/registry"
              class="text-left p-3 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors"
            >
              <i class="fas fa-database mr-2"></i>注册表管理
            </router-link>
            <router-link
              to="/files"
              class="text-left p-3 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors"
            >
              <i class="fas fa-file-alt mr-2"></i>文件管理
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- 图表区域 -->
    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">扫描统计</h2>
        <canvas ref="scanChart" width="400" height="200"></canvas>
      </div>

      <div class="bg-white rounded-lg shadow-md p-6">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">威胁分布</h2>
        <canvas ref="threatChart" width="400" height="200"></canvas>
      </div>
    </div>

    <!-- 日志区域 -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-lg font-semibold text-gray-900">
          <i class="fas fa-clipboard-list mr-2"></i>系统日志
        </h2>
        <button
          @click="clearLogs"
          class="text-sm text-gray-500 hover:text-gray-700"
        >
          清空日志
        </button>
      </div>
      <div
        class="bg-gray-900 text-green-400 p-4 rounded-md h-64 overflow-y-auto font-mono text-sm"
      >
        <div
          v-for="log in recentLogs"
          :key="log.id"
          class="log-entry"
          :class="log.level"
        >
          [{{ log.timestamp }}] {{ log.message }}
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from "vue";
import { useAppStore } from "@/stores/app";
import { fileAPI } from "@/services/api";
import { metricsAPI } from "@/services/api";
import Chart from "chart.js/auto";

const appStore = useAppStore();

// 响应式数据
const scanType = ref("file");
const filePath = ref("");
const directoryPath = ref("");
const bufferContent = ref("");
const isScanning = ref(false);
const scanResult = ref(null);
const scanChart = ref(null);
const threatChart = ref(null);
const healthStatus = ref({
  status: "unknown",
  response_time: null,
  last_check: null,
});

// 计算属性
const dashboardStats = computed(() => appStore.dashboardStats);
const systemStatus = computed(() => appStore.systemStatus);
const recentLogs = computed(() => appStore.recentLogs);

// 方法
const performScan = async () => {
  if (isScanning.value) return;

  isScanning.value = true;
  scanResult.value = null;

  try {
    let response;
    switch (scanType.value) {
      case "file":
        if (!filePath.value) {
          throw new Error("请输入文件路径");
        }
        response = await fileAPI.scan({ file_path: filePath.value });
        break;
      case "directory":
        if (!directoryPath.value) {
          throw new Error("请输入目录路径");
        }
        response = await fileAPI.scanDirectory({
          path: directoryPath.value,
          recursive: true,
          max_depth: 10,
          include_patterns: [
            "*.exe",
            "*.dll",
            "*.sys",
            "*.bat",
            "*.cmd",
            "*.ps1",
            "*.vbs",
            "*.js",
            "*.jar",
            "*.msi",
            "*.scr",
            "*.com",
            "*.log",
            "*.txt",
          ],
          exclude_patterns: ["*.tmp", "*.bak"],
        });
        break;
      case "buffer":
        if (!bufferContent.value) {
          throw new Error("请输入缓冲区内容");
        }
        response = await fileAPI.scanBuffer({
          identifier: "web-buffer",
          data: bufferContent.value,
        });
        break;
    }

    scanResult.value = {
      success: true,
      data: response.data,
    };
    appStore.addLog(`扫描完成: ${scanType.value}`, "info");
  } catch (error) {
    scanResult.value = {
      success: false,
      message: error.message,
    };
    appStore.addLog(`扫描失败: ${error.message}`, "error");
  } finally {
    isScanning.value = false;
  }
};

const clearLogs = () => {
  appStore.clearLogs();
};

const checkHealth = async () => {
  try {
    const startTime = Date.now();
    const response = await fetch("/api/health");
    const endTime = Date.now();

    healthStatus.value = {
      status: response.ok ? "healthy" : "unhealthy",
      response_time: endTime - startTime,
      last_check: new Date().toLocaleString(),
    };

    if (response.ok) {
      appStore.addLog("健康检查通过", "info");
    } else {
      appStore.addLog("健康检查失败", "error");
    }
  } catch (error) {
    healthStatus.value = {
      status: "unhealthy",
      response_time: null,
      last_check: new Date().toLocaleString(),
    };
    appStore.addLog(`健康检查失败: ${error.message}`, "error");
  }
};

const getMetrics = async () => {
  try {
    const response = await metricsAPI.getMetrics();
    console.log("性能指标:", response.data);
    // 这里可以更新仪表板上的性能指标显示
  } catch (error) {
    console.error("获取性能指标失败:", error);
  }
};

const resetMetrics = async () => {
  if (!confirm("确定要重置性能指标吗？")) return;

  try {
    await metricsAPI.resetMetrics();
    appStore.addLog("性能指标已重置", "info");
  } catch (error) {
    console.error("重置性能指标失败:", error);
    appStore.addLog("重置性能指标失败", "error");
  }
};

const initCharts = () => {
  nextTick(() => {
    // 扫描统计图表
    if (scanChart.value) {
      new Chart(scanChart.value, {
        type: "line",
        data: {
          labels: ["00:00", "04:00", "08:00", "12:00", "16:00", "20:00"],
          datasets: [
            {
              label: "扫描文件数",
              data: [0, 0, 0, 0, 0, 0],
              borderColor: "rgb(59, 130, 246)",
              backgroundColor: "rgba(59, 130, 246, 0.1)",
              tension: 0.1,
            },
          ],
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              position: "top",
            },
          },
          scales: {
            y: {
              beginAtZero: true,
            },
          },
        },
      });
    }

    // 威胁分布图表
    if (threatChart.value) {
      new Chart(threatChart.value, {
        type: "doughnut",
        data: {
          labels: ["恶意软件", "特洛伊木马", "勒索软件", "后门程序", "其他"],
          datasets: [
            {
              data: [0, 0, 0, 0, 0],
              backgroundColor: [
                "rgba(255, 99, 132, 0.8)",
                "rgba(54, 162, 235, 0.8)",
                "rgba(255, 205, 86, 0.8)",
                "rgba(75, 192, 192, 0.8)",
                "rgba(153, 102, 255, 0.8)",
              ],
            },
          ],
        },
        options: {
          responsive: true,
          plugins: {
            legend: {
              position: "bottom",
            },
          },
        },
      });
    }
  });
};

// 生命周期
onMounted(() => {
  initCharts();
  checkHealth(); // 初始健康检查
});
</script>
