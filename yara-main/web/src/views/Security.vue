<template>
  <div class="space-y-6">
    <div class="bg-white rounded-lg shadow-md p-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">
        <i class="fas fa-shield-alt mr-2"></i>安全管理
      </h1>

      <!-- 安全状态概览 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-blue-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-blue-600">
            {{ securityStatus.rules_count || 0 }}
          </div>
          <div class="text-sm text-blue-600">已加载规则</div>
        </div>
        <div class="bg-green-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-green-600">
            {{ securityStatus.scanned_files || 0 }}
          </div>
          <div class="text-sm text-green-600">已扫描文件</div>
        </div>
        <div class="bg-red-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-red-600">
            {{ securityStatus.threat_count || 0 }}
          </div>
          <div class="text-sm text-red-600">检测威胁</div>
        </div>
        <div class="bg-yellow-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-yellow-600">
            {{ quarantineList.length }}
          </div>
          <div class="text-sm text-yellow-600">隔离文件</div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex flex-wrap gap-3 mb-6">
        <button
          @click="reloadRules"
          :disabled="isLoading"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
        >
          <i class="fas fa-sync-alt mr-2" :class="{ 'fa-spin': isLoading }"></i>
          重新加载规则
        </button>
        <button
          @click="clearCache"
          :disabled="isLoading"
          class="px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 disabled:opacity-50"
        >
          <i class="fas fa-trash mr-2"></i>
          清空缓存
        </button>
        <button
          @click="getCacheStats"
          :disabled="isLoading"
          class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 disabled:opacity-50"
        >
          <i class="fas fa-chart-bar mr-2"></i>
          缓存统计
        </button>
        <button
          @click="getScanHistory"
          :disabled="isLoading"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50"
        >
          <i class="fas fa-history mr-2"></i>
          扫描历史
        </button>
        <button
          @click="loadSecurityData"
          :disabled="isLoading"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50"
        >
          <i class="fas fa-sync-alt mr-2" :class="{ 'fa-spin': isLoading }"></i>
          刷新数据
        </button>
      </div>

      <!-- 文件隔离操作 -->
      <div class="bg-red-50 p-4 rounded-lg mb-6">
        <h3 class="text-lg font-semibold text-red-900 mb-3">文件隔离</h3>
        <div class="flex gap-3">
          <input
            v-model="quarantineFilePath"
            type="text"
            placeholder="输入要隔离的文件路径"
            class="flex-1 px-4 py-2 border border-red-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-red-500"
          />
          <button
            @click="quarantineFile(quarantineFilePath)"
            :disabled="!quarantineFilePath || isLoading"
            class="px-6 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50"
          >
            <i class="fas fa-shield-alt mr-2"></i>
            隔离文件
          </button>
        </div>
      </div>

      <!-- 规则引擎状态 -->
      <div class="bg-gray-50 p-4 rounded-lg mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-3">规则引擎状态</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <span class="text-sm text-gray-600">状态:</span>
            <span
              class="ml-2 px-2 py-1 text-xs rounded-full bg-green-100 text-green-800"
            >
              {{ securityStatus.status || "unknown" }}
            </span>
          </div>
          <div>
            <span class="text-sm text-gray-600">规则目录:</span>
            <span class="ml-2 text-sm text-gray-900">{{
              securityStatus.rules_dir || "N/A"
            }}</span>
          </div>
          <div>
            <span class="text-sm text-gray-600">缓存状态:</span>
            <span
              class="ml-2 px-2 py-1 text-xs rounded-full bg-blue-100 text-blue-800"
            >
              {{ cacheStats.enabled ? "启用" : "禁用" }}
            </span>
          </div>
        </div>
      </div>

      <!-- 安全规则列表 -->
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">安全规则</h3>
        <div class="space-y-3 max-h-64 overflow-y-auto">
          <div
            v-for="rule in securityRules"
            :key="rule.name"
            class="border border-gray-200 rounded-lg p-4"
          >
            <div class="flex justify-between items-start">
              <div class="flex-1">
                <h4 class="font-medium text-gray-900">
                  {{ rule.name || rule.rule_name }}
                </h4>
                <p class="text-sm text-gray-600 mt-1">{{ rule.description }}</p>
                <div
                  class="flex items-center mt-2 space-x-4 text-xs text-gray-500"
                >
                  <span>严重性: {{ rule.severity }}</span>
                  <span>类别: {{ rule.category }}</span>
                  <span>作者: {{ rule.author || "Unknown" }}</span>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <span
                  class="px-2 py-1 text-xs rounded-full"
                  :class="getSeverityClass(rule.severity)"
                >
                  {{ rule.severity }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 隔离文件列表 -->
      <div>
        <h3 class="text-lg font-semibold text-gray-900 mb-4">隔离文件</h3>
        <div class="space-y-3">
          <div
            v-for="file in quarantineList"
            :key="file.path"
            class="border border-red-200 rounded-lg p-4 bg-red-50"
          >
            <div class="flex justify-between items-start">
              <div class="flex-1">
                <h4 class="font-medium text-red-900">{{ file.path }}</h4>
                <p class="text-sm text-red-600 mt-1">
                  隔离时间: {{ file.quarantine_time }}
                </p>
                <div
                  class="flex items-center mt-2 space-x-4 text-xs text-red-500"
                >
                  <span>威胁数: {{ file.threats?.length || 0 }}</span>
                  <span>文件大小: {{ formatBytes(file.size || 0) }}</span>
                </div>
                <div
                  v-if="file.threats && file.threats.length > 0"
                  class="mt-2"
                >
                  <div class="text-xs text-red-600">检测到的威胁:</div>
                  <div class="space-y-1 mt-1">
                    <div
                      v-for="threat in file.threats"
                      :key="threat.rule_name"
                      class="text-xs text-red-500"
                    >
                      • {{ threat.rule_name }}: {{ threat.description }}
                    </div>
                  </div>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <button
                  @click="restoreFile(file.path)"
                  class="text-green-600 hover:text-green-900"
                >
                  <i class="fas fa-undo"></i>
                </button>
                <button
                  @click="deleteFile(file.path)"
                  class="text-red-600 hover:text-red-900"
                >
                  <i class="fas fa-trash"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 缓存统计 -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">
        <i class="fas fa-database mr-2"></i>缓存统计
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div class="bg-blue-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-blue-600">
            {{ cacheStats.hit_count || 0 }}
          </div>
          <div class="text-sm text-blue-600">缓存命中</div>
        </div>
        <div class="bg-green-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-green-600">
            {{ cacheStats.miss_count || 0 }}
          </div>
          <div class="text-sm text-green-600">缓存未命中</div>
        </div>
        <div class="bg-yellow-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-yellow-600">
            {{ cacheStats.size || 0 }}
          </div>
          <div class="text-sm text-yellow-600">缓存大小</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { securityAPI } from "@/services/api";

// 响应式数据
const securityStatus = ref({});
const securityRules = ref([]);
const quarantineList = ref([]);
const cacheStats = ref({});
const isLoading = ref(false);
const quarantineFilePath = ref("");

// 方法
const loadSecurityData = async () => {
  isLoading.value = true;
  try {
    const [statusResponse, rulesResponse, quarantineResponse, cacheResponse] =
      await Promise.all([
        securityAPI.getStatus(),
        securityAPI.getRules(),
        securityAPI.getQuarantineList(),
        securityAPI.getCacheStats(),
      ]);

    securityStatus.value = statusResponse.data;
    securityRules.value = rulesResponse.data.rules || [];
    quarantineList.value = quarantineResponse.data.files || [];
    cacheStats.value = cacheResponse.data;
  } catch (error) {
    console.error("加载安全数据失败:", error);
  } finally {
    isLoading.value = false;
  }
};

const reloadRules = async () => {
  try {
    await securityAPI.reloadRules();
    await loadSecurityData();
  } catch (error) {
    console.error("重新加载规则失败:", error);
  }
};

const clearCache = async () => {
  if (!confirm("确定要清空缓存吗？")) return;

  try {
    await securityAPI.clearCache();
    await loadSecurityData();
  } catch (error) {
    console.error("清空缓存失败:", error);
  }
};

const getCacheStats = async () => {
  try {
    const response = await securityAPI.getCacheStats();
    console.log("缓存统计:", response.data);
    alert(`缓存统计信息: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("获取缓存统计失败:", error);
    alert("获取缓存统计失败");
  }
};

const getScanHistory = async () => {
  try {
    const response = await securityAPI.getScanHistory();
    console.log("扫描历史:", response.data);
    alert(`扫描历史: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("获取扫描历史失败:", error);
    alert("获取扫描历史失败");
  }
};

const restoreFile = async (filePath) => {
  if (!confirm("确定要恢复这个文件吗？")) return;

  try {
    await securityAPI.restore({ file_path: filePath });
    await loadSecurityData();
  } catch (error) {
    console.error("恢复文件失败:", error);
  }
};

const deleteFile = async (filePath) => {
  if (!confirm("确定要删除这个文件吗？此操作不可恢复！")) return;

  try {
    // 这里需要实现删除隔离文件的API
    console.log("删除文件:", filePath);
    await loadSecurityData();
  } catch (error) {
    console.error("删除文件失败:", error);
  }
};

const quarantineFile = async (filePath) => {
  if (!confirm("确定要隔离这个文件吗？")) return;

  try {
    await securityAPI.quarantine({ file_path: filePath });
    await loadSecurityData();
  } catch (error) {
    console.error("隔离文件失败:", error);
  }
};

const getSeverityClass = (severity) => {
  switch (severity?.toLowerCase()) {
    case "high":
      return "bg-red-100 text-red-800";
    case "medium":
      return "bg-yellow-100 text-yellow-800";
    case "low":
      return "bg-green-100 text-green-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

const formatBytes = (bytes) => {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
};

// 生命周期
onMounted(() => {
  loadSecurityData();
});
</script>
