import axios from "axios";

const API_BASE_URL = "/api/v1";

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    "Content-Type": "application/json",
  },
});

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    console.log(`API Request: ${config.method?.toUpperCase()} ${config.url}`);
    return config;
  },
  (error) => {
    console.error("API Request Error:", error);
    return Promise.reject(error);
  },
);

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    console.log(`API Response: ${response.status} ${response.config.url}`);
    return response;
  },
  (error) => {
    console.error("API Response Error:", error);
    return Promise.reject(error);
  },
);

// 健康检查
export const healthCheck = () => {
  return axios.get("/api/health");
};

// 文件相关API
export const fileAPI = {
  // 对于上传文件，需要明确使用 multipart/form-data，否则全局 JSON header 可能导致问题
  scan: (formData) =>
    api.post("/file/scan", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    }),
  scanDirectory: (data) => api.post("/file/scan-directory", data),
  scanBuffer: (data) => api.post("/file/scan-buffer", data),
  getInfo: (path) => api.get(`/file/info/${encodeURIComponent(path)}`),
  getHash: (path, algorithm = "sha256") =>
    api.get(`/file/hash/${encodeURIComponent(path)}?algorithm=${algorithm}`),
  getHashes: (path) => api.get(`/file/hashes/${encodeURIComponent(path)}`),
  verifyHash: (data) => api.post("/file/verify-hash", data),
  getList: (path) => api.get("/file/list", { params: { path } }),
  copy: (data) => api.post("/file/copy", data),
  move: (data) => api.post("/file/move", data),
  delete: (path) => api.delete(`/file/${encodeURIComponent(path)}`),
};

// 进程相关API
export const processAPI = {
  getList: (params) => api.get("/process/list", { params }),
  getById: (pid) => api.get(`/process/${pid}`),
  getConnections: (pid) => api.get(`/process/${pid}/connections`),
  getMemory: (pid) => api.get(`/process/${pid}/memory`),
  getRunning: (pid) => api.get(`/process/${pid}/running`),
  getChildren: (pid) => api.get(`/process/${pid}/children`),
  getModules: (pid) => api.get(`/process/${pid}/modules`),
  start: (data) => api.post("/process/start", data),
  kill: (pid) => api.delete(`/process/${pid}`),
  suspend: (pid) => api.put(`/process/${pid}/suspend`),
  resume: (pid) => api.put(`/process/${pid}/resume`),
  getStatistics: () => api.get("/process/statistics"),
  // 系统模块相关接口
  getSystemModules: () => api.get("/process/modules"),
  getModuleInfo: (moduleName) => api.get(`/process/module/${moduleName}`),
  suspendModule: (moduleName) =>
    api.put(`/process/module/${moduleName}/suspend`),
  resumeModule: (moduleName) => api.put(`/process/module/${moduleName}/resume`),
  killModule: (moduleName) => api.delete(`/process/module/${moduleName}`),
  // 监控相关接口
  enableMonitoring: () => api.post("/process/monitoring/enable"),
  disableMonitoring: () => api.post("/process/monitoring/disable"),
  getMonitoredProcesses: () => api.get("/process/monitoring/list"),
};

// 网络相关API
export const networkAPI = {
  getConnections: () => api.get("/network/connections"),
  getTCPConnections: () => api.get("/network/connections/tcp"),
  getUDPConnections: () => api.get("/network/connections/udp"),
  getConnectionsByPID: (pid) => api.get(`/network/connections/pid/${pid}`),
  getConnectionsByPort: (port) => api.get(`/network/connections/port/${port}`),
  getConnectionsByIP: (ip) => api.get(`/network/connections/ip/${ip}`),
  closeConnection: (id) => api.delete(`/network/connection/${id}`),
  isPortInUse: (port) => api.get(`/network/port/${port}/in-use`),
  getListeningPorts: () => api.get("/network/listening-ports"),
  getEstablishedConnections: () => api.get("/network/established-connections"),
  getInterfaces: () => api.get("/network/interfaces"),
  getStats: () => api.get("/network/stats"),
  // 监控相关接口
  enableMonitoring: () => api.post("/network/monitoring/enable"),
  disableMonitoring: () => api.post("/network/monitoring/disable"),
  getMonitoredConnections: () => api.get("/network/monitoring/connections"),
  getConnectionHistory: () => api.get("/network/monitoring/history"),
};

// 注册表相关API
export const registryAPI = {
  getKeyInfo: (path) => api.get(`/registry/key/${encodeURIComponent(path)}`),
  getValue: (path, valueName) =>
    api.get(
      `/registry/value/${encodeURIComponent(path)}/${encodeURIComponent(valueName)}`,
    ),
  createKey: (data) => api.post("/registry/key", data),
  setValue: (data) => api.put("/registry/value", data),
  deleteValue: (path, valueName) =>
    api.delete(
      `/registry/value/${encodeURIComponent(path)}/${encodeURIComponent(valueName)}`,
    ),
  deleteKey: (path) => api.delete(`/registry/key/${encodeURIComponent(path)}`),
  search: (params) => api.get("/registry/search", { params }),
  listKeys: (path) => api.get(`/registry/keys/${encodeURIComponent(path)}`),
  listValues: (path) => api.get(`/registry/values/${encodeURIComponent(path)}`),
};

// 安全相关API
export const securityAPI = {
  // 支持文件上传扫描
  scan: (formData) =>
    api.post("/security/scan", formData, {
      headers: { "Content-Type": "multipart/form-data" },
    }),
  getStatus: () => api.get("/security/status"),
  getRules: () => api.get("/security/rules"),
  reloadRules: () => api.post("/security/reload-rules"),
  clearCache: () => api.post("/security/cache/clear"),
  getCacheStats: () => api.get("/security/cache/stats"),
  quarantine: (data) => api.post("/security/quarantine", data),
  restore: (data) => api.post("/security/restore", data),
  getQuarantineList: () => api.get("/security/quarantine/list"),
  getScanHistory: () => api.get("/security/scan-history"),
};

// 用户相关API
export const userAPI = {
  getCurrent: () => api.get("/user/current"),
  getById: (uid) => api.get(`/user/id/${uid}`),
  getByName: (username) => api.get(`/user/name/${username}`),
  getAll: () => api.get("/user/all"),
  checkPermissions: (data) => api.post("/user/permissions/check", data),
  validatePassword: (data) => api.post("/user/password/validate", data),
  getPasswordPolicy: () => api.get("/user/password/policy"),
  checkStatus: (username) => api.get(`/user/status/${username}`),
  getSessions: (username) => api.get(`/user/sessions/${username}`),
  killSession: (sessionId) => api.delete(`/user/session/${sessionId}`),
  lockAccount: (username) => api.post(`/user/lock/${username}`),
  unlockAccount: (username) => api.post(`/user/unlock/${username}`),
  changePassword: (data) => api.post("/user/password/change", data),
  getGroups: (username) => api.get(`/user/groups/${username}`),
  getHistory: (username) => api.get(`/user/history/${username}`),
};

// 指标相关API
export const metricsAPI = {
  getMetrics: () => axios.get("/api/metrics"),
  resetMetrics: () => axios.post("/api/metrics/reset"),
};

export default api;
