<template>
  <div class="space-y-6">
    <div class="bg-white rounded-lg shadow-md p-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">
        <i class="fas fa-cogs mr-2"></i>进程管理
      </h1>

      <!-- 统计信息 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-blue-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-blue-600">
            {{ processes.length }}
          </div>
          <div class="text-sm text-blue-600">总进程数</div>
        </div>
        <div class="bg-green-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-green-600">
            {{ runningProcesses }}
          </div>
          <div class="text-sm text-green-600">运行中</div>
        </div>
        <div class="bg-yellow-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-yellow-600">
            {{ suspendedProcesses }}
          </div>
          <div class="text-sm text-yellow-600">已暂停</div>
        </div>
        <div class="bg-red-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-red-600">
            {{ suspiciousProcesses }}
          </div>
          <div class="text-sm text-red-600">可疑进程</div>
        </div>
      </div>

      <!-- 高级操作面板 -->
      <div class="bg-gray-50 p-4 rounded-lg mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">高级操作</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >启动新进程</label
            >
            <div class="flex">
              <input
                v-model="newProcessCommand"
                type="text"
                placeholder="输入命令"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="startProcess"
                class="px-4 py-2 bg-green-600 text-white rounded-r-md hover:bg-green-700"
              >
                <i class="fas fa-play"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >检查进程状态</label
            >
            <div class="flex">
              <input
                v-model="checkProcessPID"
                type="number"
                placeholder="输入PID"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="checkProcessRunning"
                class="px-4 py-2 bg-blue-600 text-white rounded-r-md hover:bg-blue-700"
              >
                <i class="fas fa-check"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >系统模块管理</label
            >
            <div class="flex">
              <input
                v-model="moduleName"
                type="text"
                placeholder="模块名称"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="getModuleInfo"
                class="px-4 py-2 bg-purple-600 text-white rounded-r-md hover:bg-purple-700"
              >
                <i class="fas fa-info"></i>
              </button>
            </div>
            <div class="grid grid-cols-3 gap-2 mt-2">
              <button
                @click="suspendModule(moduleName)"
                class="px-3 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 text-sm"
              >
                <i class="fas fa-pause mr-1"></i>挂起
              </button>
              <button
                @click="resumeModule(moduleName)"
                class="px-3 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 text-sm"
              >
                <i class="fas fa-play mr-1"></i>恢复
              </button>
              <button
                @click="killModule(moduleName)"
                class="px-3 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 text-sm"
              >
                <i class="fas fa-times mr-1"></i>结束
              </button>
            </div>
            <button
              @click="getSystemModules"
              class="w-full mt-2 px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
            >
              <i class="fas fa-list mr-2"></i>系统模块列表
            </button>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >获取系统统计</label
            >
            <button
              @click="getProcessStatistics"
              class="w-full px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
            >
              <i class="fas fa-chart-bar mr-2"></i>统计信息
            </button>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >进程监控</label
            >
            <div class="grid grid-cols-2 gap-2">
              <button
                @click="enableProcessMonitoring"
                class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
              >
                <i class="fas fa-play mr-2"></i>启用监控
              </button>
              <button
                @click="disableProcessMonitoring"
                class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
              >
                <i class="fas fa-stop mr-2"></i>禁用监控
              </button>
            </div>
            <button
              @click="getMonitoredProcesses"
              class="w-full mt-2 px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700"
            >
              <i class="fas fa-list mr-2"></i>监控列表
            </button>
          </div>
        </div>
      </div>

      <!-- 搜索和过滤 -->
      <div class="flex flex-col md:flex-row gap-4 mb-6">
        <div class="flex-1">
          <input
            v-model="searchTerm"
            type="text"
            placeholder="搜索进程名称或PID..."
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div class="flex gap-2">
          <select
            v-model="statusFilter"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">全部状态</option>
            <option value="running">运行中</option>
            <option value="suspended">已暂停</option>
            <option value="stopped">已停止</option>
          </select>
          <button
            @click="refreshProcesses"
            :disabled="isLoading"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50"
          >
            <i
              class="fas fa-sync-alt mr-2"
              :class="{ 'fa-spin': isLoading }"
            ></i>
            刷新
          </button>
        </div>
      </div>

      <!-- 进程列表 -->
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                进程
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                PID
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                状态
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                CPU
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                内存
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                操作
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="process in filteredProcesses"
              :key="process.pid"
              class="hover:bg-gray-50"
            >
              <td class="px-4 py-3">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-8 w-8">
                    <div
                      class="h-8 w-8 rounded-full bg-gray-300 flex items-center justify-center"
                    >
                      <i class="fas fa-cog text-gray-600"></i>
                    </div>
                  </div>
                  <div class="ml-3">
                    <div class="text-sm font-medium text-gray-900">
                      {{ process.name || process.process_name }}
                    </div>
                    <div class="text-sm text-gray-500">
                      {{ process.exe || process.command }}
                    </div>
                  </div>
                </div>
              </td>
              <td class="px-4 py-3 text-sm text-gray-900">{{ process.pid }}</td>
              <td class="px-4 py-3">
                <span
                  class="px-2 py-1 text-xs rounded-full"
                  :class="getStatusClass(process.status)"
                >
                  {{ getStatusText(process.status) }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-900">
                {{ process.cpu_percent || 0 }}%
              </td>
              <td class="px-4 py-3 text-sm text-gray-900">
                {{
                  formatBytes(
                    process.memory_info?.rss || process.memory_usage || 0,
                  )
                }}
              </td>
              <td class="px-4 py-3 text-sm font-medium">
                <div class="flex space-x-2">
                  <button
                    @click="showProcessDetails(process)"
                    class="text-blue-600 hover:text-blue-900"
                    title="查看详情"
                  >
                    <i class="fas fa-eye"></i>
                  </button>
                  <button
                    @click="getProcessConnections(process.pid)"
                    class="text-green-600 hover:text-green-900"
                    title="网络连接"
                  >
                    <i class="fas fa-network-wired"></i>
                  </button>
                  <button
                    @click="getProcessMemory(process.pid)"
                    class="text-purple-600 hover:text-purple-900"
                    title="内存信息"
                  >
                    <i class="fas fa-memory"></i>
                  </button>
                  <button
                    @click="getProcessChildren(process.pid)"
                    class="text-orange-600 hover:text-orange-900"
                    title="子进程"
                  >
                    <i class="fas fa-sitemap"></i>
                  </button>
                  <button
                    @click="getProcessModules(process.pid)"
                    class="text-indigo-600 hover:text-indigo-900"
                    title="模块信息"
                  >
                    <i class="fas fa-puzzle-piece"></i>
                  </button>
                  <button
                    v-if="process.status === 'running'"
                    @click="suspendProcess(process.pid)"
                    class="text-yellow-600 hover:text-yellow-900"
                    title="暂停进程"
                  >
                    <i class="fas fa-pause"></i>
                  </button>
                  <button
                    v-if="process.status === 'suspended'"
                    @click="resumeProcess(process.pid)"
                    class="text-green-600 hover:text-green-900"
                    title="恢复进程"
                  >
                    <i class="fas fa-play"></i>
                  </button>
                  <button
                    @click="killProcess(process.pid)"
                    class="text-red-600 hover:text-red-900"
                    title="终止进程"
                  >
                    <i class="fas fa-times"></i>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- 进程详情模态框 -->
    <div
      v-if="selectedProcess"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 z-50"
    >
      <div class="flex items-center justify-center min-h-screen">
        <div
          class="bg-white rounded-lg shadow-xl max-w-4xl w-full mx-4 max-h-screen overflow-y-auto"
        >
          <div class="flex justify-between items-center p-6 border-b">
            <h3 class="text-lg font-semibold text-gray-900">进程详情</h3>
            <button
              @click="selectedProcess = null"
              class="text-gray-400 hover:text-gray-600"
            >
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="p-6">
            <div class="space-y-6">
              <!-- 基本信息 -->
              <div>
                <h4 class="text-md font-semibold text-gray-900 mb-3">
                  基本信息
                </h4>
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700"
                      >进程名称</label
                    >
                    <p class="mt-1 text-sm text-gray-900">
                      {{ selectedProcess.name || selectedProcess.process_name }}
                    </p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700"
                      >PID</label
                    >
                    <p class="mt-1 text-sm text-gray-900">
                      {{ selectedProcess.pid }}
                    </p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700"
                      >状态</label
                    >
                    <p class="mt-1 text-sm text-gray-900">
                      {{ getStatusText(selectedProcess.status) }}
                    </p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700"
                      >CPU使用率</label
                    >
                    <p class="mt-1 text-sm text-gray-900">
                      {{ selectedProcess.cpu_percent || 0 }}%
                    </p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700"
                      >内存使用</label
                    >
                    <p class="mt-1 text-sm text-gray-900">
                      {{
                        formatBytes(
                          selectedProcess.memory_info?.rss ||
                            selectedProcess.memory_usage ||
                            0,
                        )
                      }}
                    </p>
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700"
                      >创建时间</label
                    >
                    <p class="mt-1 text-sm text-gray-900">
                      {{ selectedProcess.create_time || "N/A" }}
                    </p>
                  </div>
                </div>
              </div>

              <!-- 详细信息 -->
              <div v-if="selectedProcess.exe">
                <h4 class="text-md font-semibold text-gray-900 mb-3">
                  可执行文件
                </h4>
                <p class="text-sm text-gray-900">{{ selectedProcess.exe }}</p>
              </div>

              <div v-if="selectedProcess.command">
                <h4 class="text-md font-semibold text-gray-900 mb-3">命令行</h4>
                <p
                  class="text-sm text-gray-900 font-mono bg-gray-100 p-2 rounded"
                >
                  {{ selectedProcess.command }}
                </p>
              </div>

              <!-- 实时信息 -->
              <div>
                <h4 class="text-md font-semibold text-gray-900 mb-3">
                  实时信息
                </h4>
                <div class="grid grid-cols-2 gap-4">
                  <div>
                    <button
                      @click="getProcessConnections(selectedProcess.pid)"
                      class="w-full px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
                    >
                      <i class="fas fa-network-wired mr-2"></i>网络连接
                    </button>
                  </div>
                  <div>
                    <button
                      @click="getProcessMemory(selectedProcess.pid)"
                      class="w-full px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700"
                    >
                      <i class="fas fa-memory mr-2"></i>内存详情
                    </button>
                  </div>
                  <div>
                    <button
                      @click="getProcessChildren(selectedProcess.pid)"
                      class="w-full px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700"
                    >
                      <i class="fas fa-sitemap mr-2"></i>子进程
                    </button>
                  </div>
                  <div>
                    <button
                      @click="getProcessModules(selectedProcess.pid)"
                      class="w-full px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
                    >
                      <i class="fas fa-puzzle-piece mr-2"></i>模块信息
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 信息显示模态框 -->
    <div
      v-if="infoModal.show"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 z-50"
    >
      <div class="flex items-center justify-center min-h-screen">
        <div
          class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-screen overflow-y-auto"
        >
          <div class="flex justify-between items-center p-6 border-b">
            <h3 class="text-lg font-semibold text-gray-900">
              {{ infoModal.title }}
            </h3>
            <button
              @click="infoModal.show = false"
              class="text-gray-400 hover:text-gray-600"
            >
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="p-6">
            <pre
              class="text-sm text-gray-900 bg-gray-100 p-4 rounded overflow-x-auto"
              >{{ infoModal.content }}</pre
            >
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { processAPI } from "@/services/api";

// 响应式数据
const processes = ref([]);
const isLoading = ref(false);
const searchTerm = ref("");
const statusFilter = ref("");
const selectedProcess = ref(null);

// 新增的响应式数据
const newProcessCommand = ref("");
const checkProcessPID = ref("");
const moduleName = ref("");
const infoModal = ref({
  show: false,
  title: "",
  content: "",
});

// 计算属性
const filteredProcesses = computed(() => {
  let filtered = processes.value;

  if (searchTerm.value) {
    const term = searchTerm.value.toLowerCase();
    filtered = filtered.filter(
      (process) =>
        (process.name && process.name.toLowerCase().includes(term)) ||
        (process.process_name &&
          process.process_name.toLowerCase().includes(term)) ||
        process.pid.toString().includes(term),
    );
  }

  if (statusFilter.value) {
    filtered = filtered.filter(
      (process) => process.status === statusFilter.value,
    );
  }

  return filtered;
});

const runningProcesses = computed(
  () => processes.value.filter((p) => p.status === "running").length,
);

const suspendedProcesses = computed(
  () => processes.value.filter((p) => p.status === "suspended").length,
);

const suspiciousProcesses = computed(
  () => processes.value.filter((p) => p.suspicious).length,
);

// 方法
const loadProcesses = async () => {
  isLoading.value = true;
  try {
    const response = await processAPI.getList();
    processes.value = response.data.processes || [];
  } catch (error) {
    console.error("加载进程列表失败:", error);
  } finally {
    isLoading.value = false;
  }
};

const refreshProcesses = () => {
  loadProcesses();
};

const showProcessDetails = (process) => {
  selectedProcess.value = process;
};

const suspendProcess = async (pid) => {
  try {
    await processAPI.suspend(pid);
    await loadProcesses();
  } catch (error) {
    console.error("暂停进程失败:", error);
  }
};

const resumeProcess = async (pid) => {
  try {
    await processAPI.resume(pid);
    await loadProcesses();
  } catch (error) {
    console.error("恢复进程失败:", error);
  }
};

const killProcess = async (pid) => {
  if (!confirm("确定要终止这个进程吗？")) return;

  try {
    await processAPI.kill(pid);
    await loadProcesses();
  } catch (error) {
    console.error("终止进程失败:", error);
  }
};

// 新增的方法
const startProcess = async () => {
  if (!newProcessCommand.value.trim()) {
    alert("请输入要启动的命令");
    return;
  }

  try {
    const response = await processAPI.start({
      command: newProcessCommand.value,
      args: [],
      working_dir: "",
    });
    console.log("进程启动成功:", response.data);
    newProcessCommand.value = "";
    await loadProcesses();
  } catch (error) {
    console.error("启动进程失败:", error);
  }
};

const checkProcessRunning = async () => {
  if (!checkProcessPID.value) {
    alert("请输入要检查的PID");
    return;
  }

  try {
    const response = await processAPI.getRunning(
      parseInt(checkProcessPID.value),
    );
    infoModal.value = {
      show: true,
      title: `进程 ${checkProcessPID.value} 状态检查`,
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("检查进程状态失败:", error);
  }
};

// 新增：进程监控功能
const enableProcessMonitoring = async () => {
  try {
    await processAPI.enableMonitoring();
    alert("进程监控已启用");
  } catch (error) {
    console.error("启用进程监控失败:", error);
  }
};

const disableProcessMonitoring = async () => {
  try {
    await processAPI.disableMonitoring();
    alert("进程监控已禁用");
  } catch (error) {
    console.error("禁用进程监控失败:", error);
  }
};

const getMonitoredProcesses = async () => {
  try {
    const response = await processAPI.getMonitoredProcesses();
    infoModal.value = {
      show: true,
      title: "监控进程列表",
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取监控进程列表失败:", error);
  }
};

const getProcessStatistics = async () => {
  try {
    const response = await processAPI.getStatistics();
    infoModal.value = {
      show: true,
      title: "进程统计信息",
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取进程统计信息失败:", error);
  }
};

const getModuleInfo = async () => {
  if (!moduleName.value.trim()) {
    alert("请输入模块名称");
    return;
  }

  try {
    const response = await processAPI.getModuleInfo(moduleName.value);
    infoModal.value = {
      show: true,
      title: `模块 ${moduleName.value} 信息`,
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取模块信息失败:", error);
  }
};

const getSystemModules = async () => {
  try {
    const response = await processAPI.getSystemModules();
    infoModal.value = {
      show: true,
      title: "系统模块列表",
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取系统模块失败:", error);
  }
};

const suspendModule = async (moduleName) => {
  try {
    await processAPI.suspendModule(moduleName);
    showToast("模块挂起成功", "success");
    await loadProcesses();
  } catch (error) {
    console.error("挂起模块失败:", error);
    showToast("挂起模块失败", "error");
  }
};

const resumeModule = async (moduleName) => {
  try {
    await processAPI.resumeModule(moduleName);
    showToast("模块恢复成功", "success");
    await loadProcesses();
  } catch (error) {
    console.error("恢复模块失败:", error);
    showToast("恢复模块失败", "error");
  }
};

const killModule = async (moduleName) => {
  try {
    await processAPI.killModule(moduleName);
    showToast("模块结束成功", "success");
    await loadProcesses();
  } catch (error) {
    console.error("结束模块失败:", error);
    showToast("结束模块失败", "error");
  }
};

const getProcessDetails = async (pid) => {
  try {
    const response = await processAPI.getById(pid);
    infoModal.value = {
      show: true,
      title: `进程 ${pid} 详细信息`,
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取进程详情失败:", error);
  }
};

const getProcessConnections = async (pid) => {
  try {
    const response = await processAPI.getConnections(pid);
    infoModal.value = {
      show: true,
      title: `进程 ${pid} 网络连接`,
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取进程连接失败:", error);
  }
};

const getProcessMemory = async (pid) => {
  try {
    const response = await processAPI.getMemory(pid);
    infoModal.value = {
      show: true,
      title: `进程 ${pid} 内存信息`,
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取进程内存失败:", error);
  }
};

const getProcessChildren = async (pid) => {
  try {
    const response = await processAPI.getChildren(pid);
    infoModal.value = {
      show: true,
      title: `进程 ${pid} 子进程`,
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取子进程失败:", error);
  }
};

const getProcessModules = async (pid) => {
  try {
    const response = await processAPI.getModules(pid);
    infoModal.value = {
      show: true,
      title: `进程 ${pid} 模块信息`,
      content: JSON.stringify(response.data, null, 2),
    };
  } catch (error) {
    console.error("获取进程模块失败:", error);
  }
};

const getStatusClass = (status) => {
  switch (status) {
    case "running":
      return "bg-green-100 text-green-800";
    case "suspended":
      return "bg-yellow-100 text-yellow-800";
    case "stopped":
      return "bg-red-100 text-red-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

const getStatusText = (status) => {
  switch (status) {
    case "running":
      return "运行中";
    case "suspended":
      return "已暂停";
    case "stopped":
      return "已停止";
    default:
      return "未知";
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
  loadProcesses();
});
</script>
