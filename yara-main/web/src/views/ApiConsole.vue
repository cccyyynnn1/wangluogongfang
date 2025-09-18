<template>
  <div class="api-console container">
    <header class="header">
      <h1><i class="fas fa-network-wired"></i> API 控制台</h1>
      <p class="subtitle">在应用内直接测试后端 API</p>
    </header>

    <main class="main-content" style="display: flex; gap: 20px">
      <div style="width: 320px">
        <div class="input-method-selector" style="flex-wrap: wrap">
          <button
            v-for="panel in panels"
            :key="panel.id"
            class="method-btn"
            :class="{ active: activePanel === panel.id }"
            @click="activePanel = panel.id"
          >
            <i :class="panel.icon"></i> {{ panel.label }}
          </button>
          <div style="flex: 1"></div>
          <label
            style="
              display: flex;
              align-items: center;
              gap: 8px;
              color: #555;
              margin-top: 8px;
            "
          >
            API 基础 URL:
            <input
              v-model="baseUrl"
              style="
                margin-left: 8px;
                padding: 6px;
                border-radius: 6px;
                border: 1px solid #e0e0e0;
                width: 180px;
              "
            />
          </label>
        </div>

        <!-- 历史与过滤 -->
        <div style="margin-top: 18px">
          <h4 style="margin-bottom: 8px">请求历史</h4>
          <div style="display: flex; gap: 8px; margin-bottom: 8px">
            <input
              v-model="store.filterText"
              placeholder="搜索 URL / method / body"
              style="
                flex: 1;
                padding: 8px;
                border-radius: 6px;
                border: 1px solid #e0e0e0;
              "
            />
            <button class="action-btn secondary" @click="store.clearHistory">
              清空
            </button>
          </div>
          <div
            style="
              max-height: 420px;
              overflow: auto;
              border: 1px solid #f0f0f0;
              border-radius: 8px;
              padding: 8px;
              background: #fff;
            "
          >
            <div
              v-for="entry in store.paged"
              :key="entry.id"
              style="
                display: flex;
                justify-content: space-between;
                align-items: center;
                padding: 6px;
                border-bottom: 1px solid #f2f2f2;
              "
            >
              <div style="flex: 1; cursor: pointer" @click="selectEntry(entry)">
                <div style="font-weight: 600">
                  {{ entry.method }}
                  <span style="color: #666">{{ entry.url }}</span>
                </div>
                <div style="font-size: 12px; color: #888">
                  {{ entry.timestamp }} • {{ entry.status || "—" }}
                </div>
              </div>
              <div
                style="
                  display: flex;
                  flex-direction: column;
                  gap: 6px;
                  margin-left: 8px;
                "
              >
                <button
                  class="history-action-btn view"
                  @click.stop="selectEntry(entry)"
                >
                  查看
                </button>
                <button
                  class="history-action-btn delete"
                  @click.stop="store.removeEntry(entry.id)"
                >
                  删除
                </button>
              </div>
            </div>
          </div>
          <div
            style="
              display: flex;
              justify-content: center;
              align-items: center;
              gap: 8px;
              margin-top: 8px;
            "
          >
            <button
              class="action-btn"
              :disabled="store.page <= 1"
              @click="store.page--"
            >
              上页
            </button>
            <div>第 {{ store.page }} / {{ store.totalPages }} 页</div>
            <button
              class="action-btn"
              :disabled="store.page >= store.totalPages"
              @click="store.page++"
            >
              下页
            </button>
          </div>
        </div>
      </div>

      <!-- 主操作区域 -->
      <div style="flex: 1">
        <div class="input-method-selector" style="flex-wrap: wrap">
          <button
            v-for="panel in panels"
            :key="panel.id"
            class="method-btn"
            :class="{ active: activePanel === panel.id }"
            @click="activePanel = panel.id"
          >
            <i :class="panel.icon"></i> {{ panel.label }}
          </button>
        </div>

        <!-- Panels -->
        <section v-show="activePanel === 'file'">
          <h3>文件</h3>
          <div class="input-group">
            <label>选择文件</label>
            <input type="file" ref="fileScanInput" />
          </div>
          <div class="input-group">
            <button class="detect-btn" @click="fileScan">上传并扫描</button>
            <button class="action-btn" @click="getFileList">
              获取文件列表
            </button>
          </div>
          <div class="input-group">
            <label>路径查询</label>
            <input
              v-model="filePath"
              placeholder="例如 C:\\path\\file.exe"
              style="
                width: 100%;
                padding: 10px;
                border-radius: 8px;
                border: 1px solid #e0e0e0;
              "
            />
          </div>
          <div class="input-group">
            <button class="action-btn" @click="fileInfo">获取信息</button>
            <button class="action-btn" @click="fileHash">获取哈希</button>
            <button class="action-btn secondary" @click="fileDelete">
              删除
            </button>
          </div>
        </section>

        <section v-show="activePanel === 'network'" class="hidden">
          <h3>网络</h3>
          <div class="input-group">
            <button class="action-btn" @click="netConnections">获取连接</button>
            <button class="action-btn" @click="netStats">获取统计</button>
            <button class="action-btn secondary" @click="listeningPorts">
              监听端口
            </button>
          </div>
          <div class="input-group">
            <label>PID</label>
            <input
              v-model="netPid"
              placeholder="PID"
              style="
                padding: 8px;
                border-radius: 6px;
                border: 1px solid #e0e0e0;
              "
            />
            <button class="action-btn" @click="netByPid">查询</button>
          </div>
        </section>

        <section v-show="activePanel === 'process'" class="hidden">
          <h3>进程</h3>
          <div class="input-group">
            <button class="action-btn" @click="getProcesses">
              获取进程列表
            </button>
          </div>
          <div class="input-group">
            <label>启动命令</label>
            <input
              v-model="procCmd"
              placeholder="命令"
              style="
                width: 100%;
                padding: 8px;
                border-radius: 6px;
                border: 1px solid #e0e0e0;
              "
            />
            <input
              v-model="procArgs"
              placeholder="参数，逗号分隔"
              style="
                width: 100%;
                padding: 8px;
                margin-top: 8px;
                border-radius: 6px;
                border: 1px solid #e0e0e0;
              "
            />
            <button class="action-btn" @click="startProc">启动</button>
          </div>
        </section>

        <section v-show="activePanel === 'registry'" class="hidden">
          <h3>注册表</h3>
          <div class="input-group">
            <label>路径</label>
            <input
              v-model="regPath"
              placeholder="HKLM\\SOFTWARE\\..."
              style="
                width: 100%;
                padding: 8px;
                border-radius: 6px;
                border: 1px solid #e0e0e0;
              "
            />
          </div>
          <div class="input-group">
            <button class="action-btn" @click="listKeys">列出键</button>
            <button class="action-btn" @click="listValues">列出值</button>
            <button class="action-btn secondary" @click="searchRegistry">
              搜索
            </button>
          </div>
        </section>

        <section v-show="activePanel === 'security'" class="hidden">
          <h3>安全</h3>
          <div class="input-group">
            <label>选择文件</label>
            <input type="file" ref="secFileInput" />
          </div>
          <div class="input-group">
            <button class="detect-btn" @click="secScan">扫描</button>
            <button class="action-btn" @click="getRules">规则列表</button>
            <button class="action-btn secondary" @click="getScanHistory">
              扫描历史
            </button>
          </div>
        </section>

        <section v-show="activePanel === 'user'" class="hidden">
          <h3>用户</h3>
          <div class="input-group">
            <button class="action-btn" @click="getCurrentUser">当前用户</button>
            <button class="action-btn" @click="getAllUsers">所有用户</button>
          </div>
          <div class="input-group">
            <label>用户名</label>
            <input
              v-model="userName"
              placeholder="用户名"
              style="
                padding: 8px;
                border-radius: 6px;
                border: 1px solid #e0e0e0;
              "
            />
            <button class="action-btn" @click="getUserHistory">历史</button>
          </div>
        </section>

        <div style="display: flex; gap: 16px; margin-top: 18px">
          <div style="flex: 1">
            <div class="results-section" v-show="showResult">
              <div class="results-header">
                <h2><i class="fas fa-terminal"></i> 响应</h2>
              </div>
              <pre
                style="
                  background: #f8f9fa;
                  padding: 15px;
                  border-radius: 8px;
                  max-height: 360px;
                  overflow: auto;
                "
                >{{ prettyResp }}</pre
              >
              <div class="results-actions">
                <button class="action-btn" @click="copyResp">复制</button>
                <button class="action-btn secondary" @click="clearResp">
                  清空
                </button>
              </div>
            </div>
          </div>

          <!-- 请求详情面板 -->
          <div
            style="
              width: 420px;
              border: 1px solid #f0f0f0;
              border-radius: 8px;
              padding: 12px;
              background: #fff;
            "
          >
            <h4>请求详情</h4>
            <div v-if="selected" style="font-size: 13px">
              <div><strong>方法:</strong> {{ selected.method }}</div>
              <div><strong>URL:</strong> {{ selected.url }}</div>
              <div><strong>时间:</strong> {{ selected.timestamp }}</div>
              <div><strong>状态:</strong> {{ selected.status }}</div>
              <div style="margin-top: 8px">
                <strong>请求体:</strong>
                <pre
                  style="
                    max-height: 140px;
                    overflow: auto;
                    background: #f8f8f8;
                    padding: 8px;
                    border-radius: 6px;
                  "
                  >{{ format(selected.request) }}</pre
                >
              </div>
              <div style="margin-top: 8px">
                <strong>响应:</strong>
                <pre
                  style="
                    max-height: 140px;
                    overflow: auto;
                    background: #f8f8f8;
                    padding: 8px;
                    border-radius: 6px;
                  "
                  >{{ format(selected.response) }}</pre
                >
              </div>
            </div>
            <div v-else style="color: #888">选择左侧历史项查看详情</div>
          </div>
        </div>
      </div>
    </main>

    <footer class="footer"><p>&copy; 2025 Yara</p></footer>
  </div>
</template>

<script>
import { ref, reactive, computed } from "vue";
import api, {
  fileAPI,
  processAPI,
  networkAPI,
  registryAPI,
  securityAPI,
  userAPI,
} from "@/services/api";
import { useApiConsoleStore } from "@/stores/apiConsole";

export default {
  name: "ApiConsole",
  setup() {
    const store = useApiConsoleStore();
    const panels = [
      { id: "file", label: "文件", icon: "fas fa-file" },
      { id: "network", label: "网络", icon: "fas fa-network-wired" },
      { id: "process", label: "进程", icon: "fas fa-tasks" },
      { id: "registry", label: "注册表", icon: "fas fa-database" },
      { id: "security", label: "安全", icon: "fas fa-shield-alt" },
      { id: "user", label: "用户", icon: "fas fa-user" },
    ];

    const activePanel = ref("file");
    const baseUrl = ref("/api/v1");

    // file
    const filePath = ref("");
    const fileScanInput = ref(null);

    // network
    const netPid = ref("");

    // process
    const procCmd = ref("");
    const procArgs = ref("");

    // registry
    const regPath = ref("");

    // security
    const secFileInput = ref(null);

    // user
    const userName = ref("");

    const resp = ref(null);
    const showResult = ref(false);
    const error = ref(null);
    const selected = ref(null);

    const prettyResp = computed(() =>
      resp.value ? JSON.stringify(resp.value, null, 2) : "",
    );

    function handleError(e) {
      try {
        error.value =
          e.response?.data?.message || e.message || JSON.stringify(e);
      } catch {
        error.value = String(e);
      }
    }

    function genId() {
      return (
        Date.now().toString() + "-" + Math.random().toString(36).slice(2, 8)
      );
    }
    function nowStr() {
      return new Date().toLocaleString();
    }
    function format(obj) {
      try {
        return JSON.stringify(obj, null, 2);
      } catch {
        return String(obj);
      }
    }
    function selectEntry(entry) {
      selected.value = entry;
    }
    function recordEntry({
      method = "GET",
      url = "-",
      request = null,
      response = null,
      status = null,
    }) {
      try {
        store.addEntry({
          id: genId(),
          method,
          url,
          timestamp: nowStr(),
          status:
            status?.toString?.() ||
            (response && response.status) ||
            status ||
            (response && response.code) ||
            null,
          request,
          response,
        });
      } catch (e) {
        console.error("记录历史失败", e);
      }
    }

    // 文件相关
    async function fileScan() {
      error.value = null;
      const f = fileScanInput.value?.files?.[0];
      if (!f) {
        error.value = "请选择文件";
        return;
      }
      const fd = new FormData();
      fd.append("file", f);
      try {
        // use securityAPI.scan for multipart uploads (backend expects security/scan)
        const r = await securityAPI.scan(fd);
        resp.value = r.data;
        showResult.value = true;
      } catch (e) {
        handleError(e);
      }
    }
    // instrumented
    async function fileScan_instrumented() {
      error.value = null;
      const f = fileScanInput.value?.files?.[0];
      if (!f) {
        error.value = "请选择文件";
        return;
      }
      const fd = new FormData();
      fd.append("file", f);
      const url = baseUrl.value + "/security/scan";
      try {
        // call security API scan endpoint for uploaded file content
        const r = await securityAPI.scan(fd);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "POST",
          url,
          request: { fileName: f.name },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "POST",
          url,
          request: { fileName: f.name },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function getFileList() {
      error.value = null;
      const url = baseUrl.value + "/file/list";
      try {
        const r = await fileAPI.getList();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function fileInfo() {
      error.value = null;
      if (!filePath.value) {
        error.value = "请输入路径";
        return;
      }
      const url =
        baseUrl.value + "/file/info/" + encodeURIComponent(filePath.value);
      try {
        const r = await fileAPI.getInfo(filePath.value);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: { path: filePath.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: { path: filePath.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function fileHash() {
      error.value = null;
      if (!filePath.value) {
        error.value = "请输入路径";
        return;
      }
      const url =
        baseUrl.value + "/file/hash/" + encodeURIComponent(filePath.value);
      try {
        const r = await fileAPI.getHash(filePath.value);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: { path: filePath.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: { path: filePath.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function fileDelete() {
      error.value = null;
      if (!filePath.value) {
        error.value = "请输入路径";
        return;
      }
      if (!confirm("确认删除？")) return;
      const url =
        baseUrl.value + "/file/delete/" + encodeURIComponent(filePath.value);
      try {
        const r = await fileAPI.delete(filePath.value);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "DELETE",
          url,
          request: { path: filePath.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "DELETE",
          url,
          request: { path: filePath.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }

    // 网络
    async function netConnections() {
      error.value = null;
      const url = baseUrl.value + "/network/connections";
      try {
        const r = await networkAPI.getConnections();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function netStats() {
      error.value = null;
      const url = baseUrl.value + "/network/stats";
      try {
        const r = await networkAPI.getStats();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function listeningPorts() {
      error.value = null;
      const url = baseUrl.value + "/network/listening-ports";
      try {
        const r = await networkAPI.getListeningPorts();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function netByPid() {
      error.value = null;
      if (!netPid.value) {
        error.value = "请输入PID";
        return;
      }
      const url =
        baseUrl.value +
        "/network/connections?pid=" +
        encodeURIComponent(netPid.value);
      try {
        const r = await networkAPI.getConnectionsByPID(netPid.value);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: { pid: netPid.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: { pid: netPid.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }

    // 进程
    async function getProcesses() {
      error.value = null;
      const url = baseUrl.value + "/process/list";
      try {
        const r = await processAPI.getList();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function startProc() {
      error.value = null;
      if (!procCmd.value) {
        error.value = "请输入命令";
        return;
      }
      const args = procArgs.value
        ? procArgs.value
            .split(",")
            .map((s) => s.trim())
            .filter(Boolean)
        : [];
      const url = baseUrl.value + "/process/start";
      try {
        const r = await processAPI.start({ command: procCmd.value, args });
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "POST",
          url,
          request: { command: procCmd.value, args },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "POST",
          url,
          request: { command: procCmd.value, args },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }

    // 注册表
    async function listKeys() {
      error.value = null;
      if (!regPath.value) {
        error.value = "请输入路径";
        return;
      }
      const url =
        baseUrl.value + "/registry/keys/" + encodeURIComponent(regPath.value);
      try {
        const r = await registryAPI.listKeys(regPath.value);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: { path: regPath.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: { path: regPath.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function listValues() {
      error.value = null;
      if (!regPath.value) {
        error.value = "请输入路径";
        return;
      }
      const url =
        baseUrl.value + "/registry/values/" + encodeURIComponent(regPath.value);
      try {
        const r = await registryAPI.listValues(regPath.value);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: { path: regPath.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: { path: regPath.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function searchRegistry() {
      error.value = null;
      if (!regPath.value) {
        error.value = "请输入搜索关键词";
        return;
      }
      const url = baseUrl.value + "/registry/search";
      try {
        const r = await api.post("/registry/search", {
          pattern: regPath.value,
        });
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "POST",
          url,
          request: { pattern: regPath.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "POST",
          url,
          request: { pattern: regPath.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }

    // 安全
    async function secScan() {
      error.value = null;
      const f = secFileInput.value?.files?.[0];
      if (!f) {
        error.value = "请选择文件";
        return;
      }
      const fd = new FormData();
      fd.append("file", f);
      const url = baseUrl.value + "/security/scan";
      try {
        const r = await securityAPI.scan(fd);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "POST",
          url,
          request: { fileName: f.name },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "POST",
          url,
          request: { fileName: f.name },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function getRules() {
      error.value = null;
      const url = baseUrl.value + "/security/rules";
      try {
        const r = await securityAPI.getRules();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function getScanHistory() {
      error.value = null;
      const url = baseUrl.value + "/security/scan-history";
      try {
        const r = await securityAPI.getScanHistory();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }

    // user
    async function getCurrentUser() {
      error.value = null;
      const url = baseUrl.value + "/user/current";
      try {
        const r = await userAPI.getCurrent();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function getAllUsers() {
      error.value = null;
      const url = baseUrl.value + "/user/all";
      try {
        const r = await userAPI.getAll();
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: null,
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }
    async function getUserHistory() {
      error.value = null;
      if (!userName.value) {
        error.value = "请输入用户名";
        return;
      }
      const url =
        baseUrl.value + "/user/history/" + encodeURIComponent(userName.value);
      try {
        const r = await userAPI.getHistory(userName.value);
        resp.value = r.data;
        showResult.value = true;
        recordEntry({
          method: "GET",
          url,
          request: { username: userName.value },
          response: r.data,
          status: r.status,
        });
      } catch (e) {
        handleError(e);
        recordEntry({
          method: "GET",
          url,
          request: { username: userName.value },
          response: e.response?.data || { error: e.message },
          status: e.response?.status || "ERR",
        });
      }
    }

    function copyResp() {
      if (resp.value)
        navigator.clipboard.writeText(JSON.stringify(resp.value, null, 2));
    }
    function clearResp() {
      resp.value = null;
      showResult.value = false;
    }

    return {
      panels,
      activePanel,
      baseUrl,
      // file
      filePath,
      fileScanInput,
      // network
      netPid,
      // process
      procCmd,
      procArgs,
      // registry
      regPath,
      // security
      secFileInput,
      // user
      userName,
      // helpers
      resp,
      prettyResp,
      showResult,
      error,
      store,
      selected,
      selectEntry,
      format,
      // methods
      fileScan: fileScan_instrumented,
      getFileList,
      fileInfo,
      fileHash,
      fileDelete,
      netConnections,
      netStats,
      listeningPorts,
      netByPid,
      getProcesses,
      startProc,
      listKeys,
      listValues,
      searchRegistry,
      secScan,
      getRules,
      getScanHistory,
      getCurrentUser,
      getAllUsers,
      getUserHistory,
      copyResp,
      clearResp,
    };
  },
};
</script>

<style scoped>
/* 可添加组件特定的样式（复用全局 styles.css） */
</style>
