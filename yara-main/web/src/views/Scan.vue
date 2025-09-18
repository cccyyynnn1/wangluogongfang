<template>
  <div class="space-y-6">
    <div class="bg-white rounded-lg shadow-md p-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">
        <i class="fas fa-search mr-2"></i>文件扫描
      </h1>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
        <!-- 扫描配置 -->
        <div class="space-y-6">
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
            <div class="flex items-center space-x-4">
              <label class="flex items-center">
                <input v-model="recursive" type="checkbox" class="mr-2" />
                <span class="text-sm text-gray-700">递归扫描</span>
              </label>
              <div class="flex items-center space-x-2">
                <span class="text-sm text-gray-700">最大深度:</span>
                <input
                  v-model="maxDepth"
                  type="number"
                  min="1"
                  max="20"
                  class="w-20 px-2 py-1 border border-gray-300 rounded"
                />
              </div>
            </div>
          </div>

          <div v-if="scanType === 'buffer'" class="space-y-3">
            <label class="block text-sm font-medium text-gray-700"
              >缓冲区内容</label
            >
            <textarea
              v-model="bufferContent"
              placeholder="输入要扫描的内容"
              rows="8"
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

        <!-- 扫描结果 -->
        <div class="space-y-4">
          <h3 class="text-lg font-semibold text-gray-900">扫描结果</h3>

          <div v-if="!scanResult" class="text-gray-500 text-center py-8">
            暂无扫描结果
          </div>

          <div v-else-if="scanResult.success" class="space-y-4">
            <div class="bg-gray-50 p-4 rounded-lg border border-gray-200">
              <div class="flex justify-between items-center mb-3">
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
              <div class="flex justify-between items-center mb-3">
                <span class="font-medium">威胁数量:</span>
                <span class="text-lg font-bold">{{
                  scanResult.data.threats?.length || 0
                }}</span>
              </div>
              <div class="flex justify-between items-center mb-3">
                <span class="font-medium">扫描时间:</span>
                <span>{{ scanResult.data.scan_time || "N/A" }}</span>
              </div>
            </div>

            <div
              v-if="
                scanResult.data.threats && scanResult.data.threats.length > 0
              "
              class="space-y-3"
            >
              <h4 class="font-medium text-red-600">检测到的威胁:</h4>
              <div class="space-y-2 max-h-64 overflow-y-auto">
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
                    严重性: {{ threat.severity }} | 类别: {{ threat.category }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div v-else class="bg-red-50 p-4 rounded-lg border border-red-200">
            <div class="text-red-600">
              <i class="fas fa-exclamation-triangle mr-2"></i>
              扫描失败: {{ scanResult.message }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 扫描历史 -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">
        <i class="fas fa-history mr-2"></i>扫描历史
      </h2>
      <div class="space-y-3">
        <div
          v-for="history in scanHistory"
          :key="history.id"
          class="flex justify-between items-center p-3 bg-gray-50 rounded"
        >
          <div>
            <div class="font-medium">{{ history.type }}</div>
            <div class="text-sm text-gray-600">{{ history.path }}</div>
          </div>
          <div class="text-right">
            <div class="text-sm text-gray-600">{{ history.time }}</div>
            <div
              class="text-sm"
              :class="history.infected ? 'text-red-600' : 'text-green-600'"
            >
              {{ history.infected ? "发现威胁" : "安全" }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { fileAPI } from "@/services/api";

// 响应式数据
const scanType = ref("file");
const filePath = ref("");
const directoryPath = ref("");
const bufferContent = ref("");
const recursive = ref(true);
const maxDepth = ref(10);
const isScanning = ref(false);
const scanResult = ref(null);
const scanHistory = ref([]);

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
          recursive: recursive.value,
          max_depth: maxDepth.value,
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

    // 添加到历史记录
    addToHistory({
      type: scanType.value,
      path:
        scanType.value === "file"
          ? filePath.value
          : scanType.value === "directory"
            ? directoryPath.value
            : "缓冲区",
      time: new Date().toLocaleString(),
      infected: response.data.is_infected,
    });
  } catch (error) {
    scanResult.value = {
      success: false,
      message: error.message,
    };
  } finally {
    isScanning.value = false;
  }
};

const addToHistory = (record) => {
  scanHistory.value.unshift({
    id: Date.now(),
    ...record,
  });

  // 保持最多20条历史记录
  if (scanHistory.value.length > 20) {
    scanHistory.value = scanHistory.value.slice(0, 20);
  }
};

// 生命周期
onMounted(() => {
  // 加载扫描历史
  scanHistory.value = [
    {
      id: 1,
      type: "file",
      path: "C:\test.exe",
      time: "2024-01-01 10:00:00",
      infected: false,
    },
    {
      id: 2,
      type: "directory",
      path: "C:\temp",
      time: "2024-01-01 09:30:00",
      infected: true,
    },
  ];
});
</script>
