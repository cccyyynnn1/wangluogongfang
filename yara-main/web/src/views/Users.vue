<template>
  <div class="space-y-6">
    <div class="bg-white rounded-lg shadow-md p-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-6">
        <i class="fas fa-user mr-2"></i>用户管理
      </h1>

      <!-- 当前用户信息 -->
      <div class="bg-blue-50 p-4 rounded-lg mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-3">当前用户</h3>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700"
              >用户名</label
            >
            <p class="mt-1 text-sm text-gray-900">
              {{ currentUser.username || currentUser.name || "N/A" }}
            </p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700"
              >用户ID</label
            >
            <p class="mt-1 text-sm text-gray-900">
              {{ currentUser.uid || currentUser.id || "N/A" }}
            </p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700"
              >主目录</label
            >
            <p class="mt-1 text-sm text-gray-900">
              {{ currentUser.home || currentUser.home_dir || "N/A" }}
            </p>
          </div>
        </div>
      </div>

      <!-- 用户列表 -->
      <div class="mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">系统用户</h3>
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead class="bg-gray-50">
              <tr>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  用户
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  UID
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  GID
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  主目录
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  Shell
                </th>
                <th
                  class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  状态
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
                v-for="user in users"
                :key="user.uid"
                class="hover:bg-gray-50"
              >
                <td class="px-4 py-3">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-8 w-8">
                      <div
                        class="h-8 w-8 rounded-full bg-gray-300 flex items-center justify-center"
                      >
                        <i class="fas fa-user text-gray-600"></i>
                      </div>
                    </div>
                    <div class="ml-3">
                      <div class="text-sm font-medium text-gray-900">
                        {{ user.username || user.name }}
                      </div>
                      <div class="text-sm text-gray-500">
                        {{ user.gecos || user.comment || "" }}
                      </div>
                    </div>
                  </div>
                </td>
                <td class="px-4 py-3 text-sm text-gray-900">
                  {{ user.uid || user.id }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-900">{{ user.gid }}</td>
                <td class="px-4 py-3 text-sm text-gray-900">
                  {{ user.home || user.home_dir }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-900">
                  {{ user.shell }}
                </td>
                <td class="px-4 py-3">
                  <span
                    class="px-2 py-1 text-xs rounded-full"
                    :class="getUserStatusClass(user.status)"
                  >
                    {{ getUserStatusText(user.status) }}
                  </span>
                </td>
                <td class="px-4 py-3 text-sm font-medium">
                  <div class="flex space-x-2">
                    <button
                      @click="showUserDetails(user)"
                      class="text-blue-600 hover:text-blue-900"
                    >
                      <i class="fas fa-eye"></i>
                    </button>
                    <button
                      v-if="user.status === 'active'"
                      @click="lockUser(user.username || user.name)"
                      class="text-yellow-600 hover:text-yellow-900"
                    >
                      <i class="fas fa-lock"></i>
                    </button>
                    <button
                      v-if="user.status === 'locked'"
                      @click="unlockUser(user.username || user.name)"
                      class="text-green-600 hover:text-green-900"
                    >
                      <i class="fas fa-unlock"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 用户操作面板 -->
      <div class="bg-gray-50 p-4 rounded-lg mb-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">用户操作</h3>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >用户ID查询</label
            >
            <div class="flex">
              <input
                v-model="searchUserId"
                type="number"
                placeholder="输入用户ID"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="getUserById"
                class="px-4 py-2 bg-blue-600 text-white rounded-r-md hover:bg-blue-700"
              >
                <i class="fas fa-search"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >用户名查询</label
            >
            <div class="flex">
              <input
                v-model="searchUsername"
                type="text"
                placeholder="输入用户名"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="getUserByName"
                class="px-4 py-2 bg-blue-600 text-white rounded-r-md hover:bg-blue-700"
              >
                <i class="fas fa-search"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >账户状态检查</label
            >
            <div class="flex">
              <input
                v-model="checkStatusUsername"
                type="text"
                placeholder="输入用户名"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="checkUserStatus"
                class="px-4 py-2 bg-yellow-600 text-white rounded-r-md hover:bg-yellow-700"
              >
                <i class="fas fa-check"></i>
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >密码验证</label
            >
            <div class="flex">
              <input
                v-model="validatePassword"
                type="password"
                placeholder="输入密码"
                class="flex-1 px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="validateUserPassword"
                class="px-4 py-2 bg-green-600 text-white rounded-r-md hover:bg-green-700"
              >
                <i class="fas fa-key"></i>
              </button>
            </div>
          </div>
        </div>

        <!-- 新增的高级操作 -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mt-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >权限检查</label
            >
            <div class="space-y-2">
              <input
                v-model="checkPermissionsUsername"
                type="text"
                placeholder="用户名"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <input
                v-model="checkPermissionsPermission"
                type="text"
                placeholder="权限名称"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="checkUserPermissions"
                class="w-full px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700"
              >
                <i class="fas fa-shield-alt mr-2"></i>检查权限
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >密码策略</label
            >
            <button
              @click="getPasswordPolicy"
              class="w-full px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700"
            >
              <i class="fas fa-cog mr-2"></i>获取策略
            </button>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >密码修改</label
            >
            <div class="space-y-2">
              <input
                v-model="changePasswordUsername"
                type="text"
                placeholder="用户名"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <input
                v-model="changePasswordNewPassword"
                type="password"
                placeholder="新密码"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
              <button
                @click="changeUserPassword"
                class="w-full px-4 py-2 bg-orange-600 text-white rounded-lg hover:bg-orange-700"
              >
                <i class="fas fa-edit mr-2"></i>修改密码
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2"
              >用户信息</label
            >
            <div class="space-y-2">
              <button
                @click="getUserGroups(currentUser.username || currentUser.name)"
                class="w-full px-4 py-2 bg-teal-600 text-white rounded-lg hover:bg-teal-700"
              >
                <i class="fas fa-users mr-2"></i>用户组
              </button>
              <button
                @click="
                  getUserHistory(currentUser.username || currentUser.name)
                "
                class="w-full px-4 py-2 bg-pink-600 text-white rounded-lg hover:bg-pink-700"
              >
                <i class="fas fa-history mr-2"></i>登录历史
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 用户会话 -->
      <div>
        <h3 class="text-lg font-semibold text-gray-900 mb-4">活跃会话</h3>
        <div class="space-y-3">
          <div
            v-for="session in userSessions"
            :key="session.id"
            class="border border-gray-200 rounded-lg p-4"
          >
            <div class="flex justify-between items-start">
              <div class="flex-1">
                <h4 class="font-medium text-gray-900">
                  {{ session.username }}
                </h4>
                <p class="text-sm text-gray-600 mt-1">
                  会话ID: {{ session.id }}
                </p>
                <div
                  class="flex items-center mt-2 space-x-4 text-xs text-gray-500"
                >
                  <span>登录时间: {{ session.login_time }}</span>
                  <span>终端: {{ session.terminal || "N/A" }}</span>
                  <span>主机: {{ session.host || "N/A" }}</span>
                </div>
              </div>
              <div class="flex items-center space-x-2">
                <button
                  @click="killSession(session.id)"
                  class="text-red-600 hover:text-red-900"
                >
                  <i class="fas fa-times"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 用户详情模态框 -->
    <div
      v-if="selectedUser"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 z-50"
    >
      <div class="flex items-center justify-center min-h-screen">
        <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full mx-4">
          <div class="flex justify-between items-center p-6 border-b">
            <h3 class="text-lg font-semibold text-gray-900">用户详情</h3>
            <button
              @click="selectedUser = null"
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
                    >用户名</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedUser.username || selectedUser.name }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >UID</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedUser.uid || selectedUser.id }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >GID</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedUser.gid }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >状态</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ getUserStatusText(selectedUser.status) }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >主目录</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedUser.home || selectedUser.home_dir }}
                  </p>
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700"
                    >Shell</label
                  >
                  <p class="mt-1 text-sm text-gray-900">
                    {{ selectedUser.shell }}
                  </p>
                </div>
              </div>

              <div v-if="selectedUser.groups">
                <label class="block text-sm font-medium text-gray-700"
                  >用户组</label
                >
                <p class="mt-1 text-sm text-gray-900">
                  {{ selectedUser.groups.join(", ") }}
                </p>
              </div>

              <div v-if="selectedUser.gecos">
                <label class="block text-sm font-medium text-gray-700"
                  >描述</label
                >
                <p class="mt-1 text-sm text-gray-900">
                  {{ selectedUser.gecos }}
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
import { ref, onMounted } from "vue";
import { userAPI } from "@/services/api";

// 响应式数据
const currentUser = ref({});
const users = ref([]);
const userSessions = ref([]);
const selectedUser = ref(null);

// 用户操作相关
const searchUserId = ref("");
const searchUsername = ref("");
const checkStatusUsername = ref("");
const validatePassword = ref("");

// 权限检查相关
const checkPermissionsUsername = ref("");
const checkPermissionsPermission = ref("");

// 密码策略相关
const getPasswordPolicyUsername = ref("");

// 密码修改相关
const changePasswordUsername = ref("");
const changePasswordNewPassword = ref("");

// 方法
const loadUserData = async () => {
  try {
    const [currentResponse, usersResponse, sessionsResponse] =
      await Promise.all([
        userAPI.getCurrent(),
        userAPI.getAll(),
        userAPI.getSessions("current"),
      ]);

    currentUser.value = currentResponse.data;
    users.value = usersResponse.data.users || [];
    userSessions.value = sessionsResponse.data.sessions || [];
  } catch (error) {
    console.error("加载用户数据失败:", error);
  }
};

const showUserDetails = (user) => {
  selectedUser.value = user;
};

const lockUser = async (username) => {
  if (!confirm(`确定要锁定用户 ${username} 吗？`)) return;

  try {
    await userAPI.lockAccount(username);
    await loadUserData();
  } catch (error) {
    console.error("锁定用户失败:", error);
  }
};

const unlockUser = async (username) => {
  if (!confirm(`确定要解锁用户 ${username} 吗？`)) return;

  try {
    await userAPI.unlockAccount(username);
    await loadUserData();
  } catch (error) {
    console.error("解锁用户失败:", error);
  }
};

const killSession = async (sessionId) => {
  if (!confirm("确定要终止这个会话吗？")) return;

  try {
    await userAPI.killSession(sessionId);
    await loadUserData();
  } catch (error) {
    console.error("终止会话失败:", error);
  }
};

const getUserById = async () => {
  if (!searchUserId.value) {
    alert("请输入用户ID");
    return;
  }

  try {
    const response = await userAPI.getById(searchUserId.value);
    console.log("用户信息:", response.data);
    alert(`用户信息: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("获取用户信息失败:", error);
    alert("获取用户信息失败");
  }
};

const getUserByName = async () => {
  if (!searchUsername.value) {
    alert("请输入用户名");
    return;
  }

  try {
    const response = await userAPI.getByName(searchUsername.value);
    console.log("用户信息:", response.data);
    alert(`用户信息: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("获取用户信息失败:", error);
    alert("获取用户信息失败");
  }
};

const checkUserStatus = async () => {
  if (!checkStatusUsername.value) {
    alert("请输入用户名");
    return;
  }

  try {
    const response = await userAPI.checkStatus(checkStatusUsername.value);
    console.log("用户状态:", response.data);
    alert(`用户状态: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("检查用户状态失败:", error);
    alert("检查用户状态失败");
  }
};

const validateUserPassword = async () => {
  if (!validatePassword.value) {
    alert("请输入密码");
    return;
  }

  try {
    const response = await userAPI.validatePassword({
      username: currentUser.value.username || currentUser.value.name,
      password: validatePassword.value,
    });
    console.log("密码验证结果:", response.data);
    alert(`密码验证结果: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("密码验证失败:", error);
    alert("密码验证失败");
  }
};

const checkUserPermissions = async () => {
  if (!checkPermissionsUsername.value || !checkPermissionsPermission.value) {
    alert("请输入用户名和权限");
    return;
  }

  try {
    const response = await userAPI.checkPermissions({
      username: checkPermissionsUsername.value,
      permission: checkPermissionsPermission.value,
    });
    console.log("权限检查结果:", response.data);
    alert(`权限检查结果: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("权限检查失败:", error);
    alert("权限检查失败");
  }
};

const getPasswordPolicy = async () => {
  try {
    const response = await userAPI.getPasswordPolicy();
    console.log("密码策略:", response.data);
    alert(`密码策略: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("获取密码策略失败:", error);
    alert("获取密码策略失败");
  }
};

const changeUserPassword = async () => {
  if (!changePasswordUsername.value || !changePasswordNewPassword.value) {
    alert("请输入用户名和新密码");
    return;
  }

  try {
    const response = await userAPI.changePassword({
      username: changePasswordUsername.value,
      new_password: changePasswordNewPassword.value,
    });
    console.log("密码修改结果:", response.data);
    alert(`密码修改结果: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("密码修改失败:", error);
    alert("密码修改失败");
  }
};

const getUserGroups = async (username) => {
  try {
    const response = await userAPI.getGroups(username);
    console.log("用户组信息:", response.data);
    alert(`用户组信息: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("获取用户组失败:", error);
    alert("获取用户组失败");
  }
};

const getUserHistory = async (username) => {
  try {
    const response = await userAPI.getHistory(username);
    console.log("用户历史记录:", response.data);
    alert(`用户历史记录: ${JSON.stringify(response.data, null, 2)}`);
  } catch (error) {
    console.error("获取用户历史失败:", error);
    alert("获取用户历史失败");
  }
};

const getUserStatusClass = (status) => {
  switch (status) {
    case "active":
      return "bg-green-100 text-green-800";
    case "locked":
      return "bg-red-100 text-red-800";
    case "inactive":
      return "bg-gray-100 text-gray-800";
    default:
      return "bg-gray-100 text-gray-800";
  }
};

const getUserStatusText = (status) => {
  switch (status) {
    case "active":
      return "活跃";
    case "locked":
      return "已锁定";
    case "inactive":
      return "非活跃";
    default:
      return "未知";
  }
};

// 生命周期
onMounted(() => {
  loadUserData();
});
</script>
