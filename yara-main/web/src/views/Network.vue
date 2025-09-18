<template>
  <div class="space-y-6">
    <div class="bg-white rounded-lg shadow-md p-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">
        <i class="fas fa-network-wired mr-2"></i>网络管理
      </h1>

      <!-- 统计信息 -->
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <div class="bg-blue-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-blue-600">
            {{ connections.length }}
          </div>
          <div class="text-sm text-blue-600">总连接数</div>
        </div>
        <div class="bg-green-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-green-600">
            {{ tcpConnections }}
          </div>
          <div class="text-sm text-green-600">TCP连接</div>
        </div>
        <div class="bg-yellow-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-yellow-600">
            {{ udpConnections }}
          </div>
          <div class="text-sm text-yellow-600">UDP连接</div>
        </div>
        <div class="bg-red-50 p-4 rounded-lg">
          <div class="text-2xl font-bold text-red-600">
            {{ listeningPorts }}
          </div>
          <div class="text-sm text-red-600">监听端口</div>
        </div>
      </div>

      <!-- 过滤和搜索 -->
      <div class="flex flex-col md:flex-row gap-4 mb-6">
        <div class="flex-1">
          <input
            v-model="searchTerm"
            type="text"
            placeholder="搜索IP地址或端口..."
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        <div class="flex gap-2">
          <select
            v-model="protocolFilter"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">全部协议</option>
            <option value="tcp">TCP</option>
            <option value="udp">UDP</option>
          </select>
          <select
            v-model="statusFilter"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">全部状态</option>
            <option value="ESTABLISHED">已建立</option>
            <option value="LISTEN">监听中</option>
            <option value="TIME_WAIT">等待中</option>
            <option value="CLOSE_WAIT">关闭等待</option>
          </select>
          <button
            @click="refreshConnections"
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

      <!-- 网络监控控制 -->
      <div class="bg-gray-50 p-4 rounded-lg mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">网络监控</h3>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <button
            @click="enableNetworkMonitoring"
            class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
          >
            <i class="fas fa-play mr-2"></i>启用监控
          </button>
          <button
            @click="disableNetworkMonitoring"
            class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
          >
            <i class="fas fa-stop mr-2"></i>禁用监控
          </button>
          <button
            @click="getMonitoredConnections"
            class="px-4 py-2 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700"
          >
            <i class="fas fa-list mr-2"></i>监控列表
          </button>
          <button
            @click="getConnectionHistory"
            class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700"
          >
            <i class="fas fa-history mr-2"></i>连接历史
          </button>
        </div>
      </div>

      <!-- 网络查询工具 -->
      <div class="bg-blue-50 p-4 rounded-lg mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">网络查询工具</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >按PID查询连接</label
            >
            <div class="flex">
              <input
                v-model="queryPID"
                type="number"
                placeholder="输入PID"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="getConnectionsByPID(queryPID)"
                class="px-4 py-2 bg-blue-600 text-white rounded-r-md hover:bg-blue-700"
              >
                <i class="fas fa-search"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >按端口查询连接</label
            >
            <div class="flex">
              <input
                v-model="queryPort"
                type="number"
                placeholder="输入端口"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="getConnectionsByPort(queryPort)"
                class="px-4 py-2 bg-blue-600 text-white rounded-r-md hover:bg-blue-700"
              >
                <i class="fas fa-search"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >按IP查询连接</label
            >
            <div class="flex">
              <input
                v-model="queryIP"
                type="text"
                placeholder="输入IP地址"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="getConnectionsByIP(queryIP)"
                class="px-4 py-2 bg-blue-600 text-white rounded-r-md hover:bg-blue-700"
              >
                <i class="fas fa-search"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >检查端口占用</label
            >
            <div class="flex">
              <input
                v-model="checkPort"
                type="number"
                placeholder="输入端口"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="checkPortInUse(checkPort)"
                class="px-4 py-2 bg-yellow-600 text-white rounded-r-md hover:bg-yellow-700"
              >
                <i class="fas fa-check"></i>
              </button>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
          <button
            @click="getListeningPorts"
            class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700"
          >
            <i class="fas fa-list mr-2"></i>监听端口
          </button>
          <button
            @click="getEstablishedConnections"
            class="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700"
          >
            <i class="fas fa-link mr-2"></i>已建立连接
          </button>
          <button
            @click="getTCPConnections"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            <i class="fas fa-network-wired mr-2"></i>TCP连接
          </button>
          <button
            @click="getUDPConnections"
            class="px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700"
          >
            <i class="fas fa-broadcast-tower mr-2"></i>UDP连接
          </button>
        </div>
      </div>

      <!-- 连接列表 -->
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead class="bg-gray-50">
            <tr>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                协议
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                本地地址
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                远程地址
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                状态
              </th>
              <th
                class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                PID
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
              v-for="conn in filteredConnections"
              :key="conn.id"
              class="hover:bg-gray-50"
            >
              <td class="px-4 py-3">
                <span
                  class="px-2 py-1 text-xs rounded-full"
                  :class="getProtocolClass(conn.protocol)"
                >
                  {{ conn.protocol?.toUpperCase() || "TCP" }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-900">
                {{ formatAddress(conn.local_addr || conn.laddr) }}
              </td>
              <td class="px-4 py-3 text-sm text-gray-900">
                {{ formatAddress(conn.remote_addr || conn.raddr) }}
              </td>
              <td class="px-4 py-3">
                <span
                  class="px-2 py-1 text-xs rounded-full"
                  :class="getStatusClass(conn.status)"
                >
                  {{ conn.status || "UNKNOWN" }}
                </span>
              </td>
              <td class="px-4 py-3 text-sm text-gray-900">
                {{ conn.pid || "N/A" }}
              </td>
              <td class="px-4 py-3 text-sm font-medium">
                <div class="flex space-x-2">
                  <button
                    @click="showConnectionDetails(conn)"
                    class="text-blue-600 hover:text-blue-900"
                  >
                    <i class="fas fa-eye"></i>
                  </button>
                  <button
                    v-if="conn.status !== 'LISTEN'"
                    @click="closeConnection(conn.id)"
                    class="text-red-600 hover:text-red-900"
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

    <!-- 网络接口信息 -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-lg font-semibold text-gray-900 mb-4">
        <i class="fas fa-ethernet mr-2"></i>网络接口
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div
          v-for="iface in interfaces"
          :key="iface.name"
          class="border border-gray-200 rounded-lg p-4"
        >
          <div class="flex items-center justify-between mb-2">
            <h3 class="font-medium text-gray-900">{{ iface.name }}</h3>
            <span
              class="px-2 py-1 text-xs rounded-full"
              :class="
                iface.status === 'up'
                  ? 'bg-green-100 text-green-800'
                  : 'bg-red-100 text-red-800'
              "
            >
              {{ iface.status === "up" ? "启用" : "禁用" }}
            </span>
          </div>
          <div class="space-y-1 text-sm text-gray-600">
            <div>IP地址: {{ iface.addresses?.join(", ") || "N/A" }}</div>
            <div>MAC地址: {{ iface.mac || "N/A" }}</div>
            <div>MTU: {{ iface.mtu || "N/A" }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 连接详情模态框 -->
    <div
      v-if="selectedConnection"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 z-50"
    >
      <div class="flex items-center justify-center min-h-screen">
        <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4">
          <div class="flex justify-between items-center p-6 border-b">
            <h3 class="text-lg font-semibold text-gray-900">连接详情</h3>
            <button
              @click="selectedConnection = null"
              class="text-gray-400 hover:text-gray-600"
            >
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="p-6">
            <div class="space-y-4">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >协议</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedConnection.protocol?.toUpperCase() || "TCP" }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >状态</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedConnection.status || "UNKNOWN" }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >PID</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedConnection.pid || "N/A" }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >进程名称</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedConnection.process_name || "N/A" }}
                  </p>
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700"
                  >本地地址</label
                >
                <p class="mt-1 text-sm text-gray-900 font-mono">
                  {{
                    formatAddress(
                      selectedConnection.local_addr || selectedConnection.laddr,
                    )
                  }}
                </p>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700"
                  >远程地址</label
                >
                <p class="mt-1 text-sm text-gray-900 font-mono">
                  {{
                    formatAddress(
                      selectedConnection.remote_addr ||
                        selectedConnection.raddr,
                    )
                  }}
                </p>
              </div>

              <div v-if="selectedConnection.create_time">
                <label class="block text-sm font-medium text-gray-700"
                  >创建时间</label
                >
                <p class="mt-1 text-sm text-gray-900">
                  {{ selectedConnection.create_time }}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { networkAPI } from "@/services/api";

// 响应式数据
const connections = ref([]);
const interfaces = ref([]);
const isLoading = ref(false);
const searchTerm = ref("");
const protocolFilter = ref("");
const statusFilter = ref("");
const selectedConnection = ref(null);

// 查询工具数据
const queryPID = ref("");
const queryPort = ref("");
const queryIP = ref("");
const checkPort = ref("");

// 计算属性
const filteredConnections = computed(() => {
  let filtered = connections.value;

  if (searchTerm.value) {
    const term = searchTerm.value.toLowerCase();
    filtered = filtered.filter(
      (conn) =>
        (conn.local_addr && conn.local_addr.toLowerCase().includes(term)) ||
        (conn.remote_addr && conn.remote_addr.toLowerCase().includes(term)) ||
        (conn.laddr &&
          conn.laddr.ip &&
          conn.laddr.ip.toLowerCase().includes(term)) ||
        (conn.raddr &&
          conn.raddr.ip &&
          conn.raddr.ip.toLowerCase().includes(term)),
    );
  }

  if (protocolFilter.value) {
    filtered = filtered.filter(
      (conn) => conn.protocol === protocolFilter.value,
    );
  }

  if (statusFilter.value) {
    filtered = filtered.filter((conn) => conn.status === statusFilter.value);
  }

  return filtered;
});

const tcpConnections = computed(
  () => connections.value.filter((c) => c.protocol === "tcp").length,
);

const udpConnections = computed(
  () => connections.value.filter((c) => c.protocol === "udp").length,
);

const listeningPorts = computed(
  () => connections.value.filter((c) => c.status === "LISTEN").length,
);

// 方法
const loadConnections = async () => {
  isLoading.value = true;
  try {
    const [connectionsResponse, interfacesResponse] = await Promise.all([
      networkAPI.getConnections(),
      networkAPI.getInterfaces(),
    ]);
    connections.value = connectionsResponse.data.connections || [];
    interfaces.value = interfacesResponse.data.interfaces || [];
  } catch (error) {
    console.error("加载网络连接失败:", error);
  } finally {
    isLoading.value = false;
  }
};

const refreshConnections = () => {
  loadConnections();
};

const showConnectionDetails = (connection) => {
  selectedConnection.value = connection;
};

const closeConnection = async (connectionId) => {
  if (!confirm("确定要关闭这个连接吗？")) return;

  try {
    await networkAPI.closeConnection(connectionId);
    await loadConnections();
  } catch (error) {
    console.error("关闭连接失败:", error);
  }
};

const getNetworkStats = async () => {
  try {
    const response = await networkAPI.getStats();
    console.log("网络统计:", response.data);
  } catch (error) {
    console.error("获取网络统计失败:", error);
  }
};

// 新增：网络监控功能
const enableNetworkMonitoring = async () => {
  try {
    await networkAPI.enableMonitoring();
    alert("网络监控已启用");
  } catch (error) {
    console.error("启用网络监控失败:", error);
  }
};

const disableNetworkMonitoring = async () => {
  try {
    await networkAPI.disableMonitoring();
    alert("网络监控已禁用");
  } catch (error) {
    console.error("禁用网络监控失败:", error);
  }
};

const getMonitoredConnections = async () => {
  try {
    const response = await networkAPI.getMonitoredConnections();
    console.log("监控连接:", response.data);
  } catch (error) {
    console.error("获取监控连接失败:", error);
  }
};

const getConnectionHistory = async () => {
  try {
    const response = await networkAPI.getConnectionHistory();
    console.log("连接历史:", response.data);
  } catch (error) {
    console.error("获取连接历史失败:", error);
  }
};

const getListeningPorts = async () => {
  try {
    const response = await networkAPI.getListeningPorts();
    console.log("监听端口:", response.data);
    // 可以在这里显示监听端口信息
  } catch (error) {
    console.error("获取监听端口失败:", error);
  }
};

const checkPortInUse = async (port) => {
  try {
    const response = await networkAPI.isPortInUse(port);
    console.log(`端口 ${port} 使用情况:`, response.data);
    // 可以在这里显示端口使用情况
  } catch (error) {
    console.error("检查端口使用情况失败:", error);
  }
};

const getEstablishedConnections = async () => {
  try {
    const response = await networkAPI.getEstablishedConnections();
    console.log("已建立连接:", response.data);
    // 可以在这里显示已建立连接信息
  } catch (error) {
    console.error("获取已建立连接失败:", error);
  }
};

const getTCPConnections = async () => {
  try {
    const response = await networkAPI.getTCPConnections();
    console.log("TCP连接:", response.data);
    // 可以在这里显示TCP连接信息
  } catch (error) {
    console.error("获取TCP连接失败:", error);
  }
};

const getUDPConnections = async () => {
  try {
    const response = await networkAPI.getUDPConnections();
    console.log("UDP连接:", response.data);
    // 可以在这里显示UDP连接信息
  } catch (error) {
    console.error("获取UDP连接失败:", error);
  }
};

const getConnectionsByPID = async (pid) => {
  try {
    const response = await networkAPI.getConnectionsByPID(pid);
    console.log(`PID ${pid} 的连接:`, response.data);
    // 可以在这里显示特定PID的连接
  } catch (error) {
    console.error(`获取PID ${pid} 连接失败:`, error);
  }
};

const getConnectionsByPort = async (port) => {
  try {
    const response = await networkAPI.getConnectionsByPort(port);
    console.log(`端口 ${port} 的连接:`, response.data);
    // 可以在这里显示特定端口的连接
  } catch (error) {
    console.error(`获取端口 ${port} 连接失败:`, error);
  }
};

const getConnectionsByIP = async (ip) => {
  try {
    const response = await networkAPI.getConnectionsByIP(ip);
    console.log(`IP ${ip} 的连接:`, response.data);
    // 可以在这里显示特定IP的连接
  } catch (error) {
    console.error(`获取IP ${ip} 连接失败:`, error);
  }
};

const getProtocolClass = (protocol) => {
  switch (protocol?.toLowerCase()) {
    case "tcp":
      return "bg-blue-100 text-blue-800";
    case "udp":
      return "bg-green-100 text-green-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

const getStatusClass = (status) => {
  switch (status) {
    case "ESTABLISHED":
      return "bg-green-100 text-green-800";
    case "LISTEN":
      return "bg-blue-100 text-blue-800";
    case "TIME_WAIT":
      return "bg-yellow-100 text-yellow-800";
    case "CLOSE_WAIT":
      return "bg-red-100 text-red-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

const formatAddress = (addr) => {
  if (!addr) return "N/A";

  if (typeof addr === "string") {
    return addr;
  }

  if (addr.ip && addr.port) {
    return `${addr.ip}:${addr.port}`;
  }

  return "N/A";
};

// 生命周期
onMounted(() => {
  loadConnections();
});
</script>
