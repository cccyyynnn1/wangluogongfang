<template>
  <div class="registry-container">
    <div class="header">
      <h1>注册表管理</h1>
      <div class="actions">
        <button
          @click="refreshData"
          class="btn btn-primary"
          :disabled="isLoading"
        >
          <i class="fas fa-sync" :class="{ 'fa-spin': isLoading }"></i> 刷新
        </button>
        <button @click="showCreateKeyModal = true" class="btn btn-success">
          <i class="fas fa-plus"></i> 新建键
        </button>
        <button @click="showSearchModal = true" class="btn btn-info">
          <i class="fas fa-search"></i> 搜索
        </button>
        <button @click="showBreadcrumbModal = true" class="btn btn-secondary">
          <i class="fas fa-route"></i> 路径导航
        </button>
      </div>
    </div>

    <div class="content">
      <div class="sidebar">
        <div class="registry-tree">
          <h3>注册表树</h3>
          <div
            class="tree-item"
            @click="navigateToPath('SOFTWARE')"
            :class="{ active: currentPath === 'SOFTWARE' }"
          >
            <i class="fas fa-folder"></i> SOFTWARE
          </div>
          <div
            class="tree-item"
            @click="navigateToPath('SYSTEM')"
            :class="{ active: currentPath === 'SYSTEM' }"
          >
            <i class="fas fa-folder"></i> SYSTEM
          </div>
          <div
            class="tree-item"
            @click="navigateToPath('HARDWARE')"
            :class="{ active: currentPath === 'HARDWARE' }"
          >
            <i class="fas fa-folder"></i> HARDWARE
          </div>
          <div
            class="tree-item"
            @click="navigateToPath('SAM')"
            :class="{ active: currentPath === 'SAM' }"
          >
            <i class="fas fa-folder"></i> SAM
          </div>
          <div
            class="tree-item"
            @click="navigateToPath('SECURITY')"
            :class="{ active: currentPath === 'SECURITY' }"
          >
            <i class="fas fa-folder"></i> SECURITY
          </div>

          <div class="tree-actions">
            <button @click="getKeyInfo" class="btn btn-sm btn-info">
              <i class="fas fa-info"></i> 键信息
            </button>
            <button @click="listKeys" class="btn btn-sm btn-secondary">
              <i class="fas fa-list"></i> 列出键
            </button>
            <button @click="listValues" class="btn btn-sm btn-warning">
              <i class="fas fa-list-alt"></i> 列出值
            </button>
          </div>
        </div>
      </div>

      <div class="main-content">
        <div class="path-navigator">
          <span class="path-label">当前路径:</span>
          <span class="path-value">{{ currentPath }}</span>
          <div class="path-actions">
            <button
              @click="goBack"
              class="btn btn-sm btn-secondary"
              :disabled="!canGoBack"
            >
              <i class="fas fa-arrow-left"></i> 返回
            </button>
            <button @click="copyPath" class="btn btn-sm btn-info">
              <i class="fas fa-copy"></i> 复制路径
            </button>
          </div>
        </div>

        <div v-if="isLoading" class="loading-overlay">
          <div class="loading-spinner">
            <i class="fas fa-spinner fa-spin"></i>
            <span>加载中...</span>
          </div>
        </div>

        <div v-else class="registry-content">
          <div class="keys-section">
            <h3>子键 ({{ subKeys.length }})</h3>
            <div class="keys-list">
              <div
                v-for="key in subKeys"
                :key="key"
                class="key-item"
                @click="navigateToSubKey(key)"
              >
                <i class="fas fa-folder"></i>
                <span>{{ key }}</span>
                <div class="key-actions">
                  <button
                    @click.stop="deleteKey(key)"
                    class="btn btn-sm btn-danger"
                  >
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </div>
              <div v-if="subKeys.length === 0" class="empty-state">
                <i class="fas fa-folder-open"></i>
                <span>没有子键</span>
              </div>
            </div>
          </div>

          <div class="values-section">
            <h3>值 ({{ values.length }})</h3>
            <div class="values-list">
              <div
                v-for="value in values"
                :key="value.name"
                class="value-item"
                @click="selectValue(value)"
              >
                <div class="value-name">{{ value.name || "(默认)" }}</div>
                <div class="value-type">{{ value.type }}</div>
                <div class="value-data">{{ value.value }}</div>
                <div class="value-actions">
                  <button
                    @click.stop="editValue(value)"
                    class="btn btn-sm btn-primary"
                  >
                    <i class="fas fa-edit"></i>
                  </button>
                  <button
                    @click.stop="deleteValue(value)"
                    class="btn btn-sm btn-danger"
                  >
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </div>
              <div v-if="values.length === 0" class="empty-state">
                <i class="fas fa-list"></i>
                <span>没有值</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建键模态框 -->
    <div v-if="showCreateKeyModal" class="modal">
      <div class="modal-content">
        <h3>创建注册表键</h3>
        <div class="form-group">
          <label>键路径:</label>
          <input v-model="newKeyPath" type="text" placeholder="输入键路径" />
        </div>
        <div class="modal-actions">
          <button
            @click="createKey"
            class="btn btn-success"
            :disabled="!newKeyPath.trim()"
          >
            <i class="fas fa-plus"></i> 创建
          </button>
          <button @click="showCreateKeyModal = false" class="btn btn-secondary">
            取消
          </button>
        </div>
      </div>
    </div>

    <!-- 搜索模态框 -->
    <div v-if="showSearchModal" class="modal">
      <div class="modal-content">
        <h3>搜索注册表</h3>
        <div class="form-group">
          <label>搜索关键词:</label>
          <input
            v-model="searchTerm"
            type="text"
            placeholder="输入搜索关键词"
          />
        </div>
        <div class="form-group">
          <label>搜索根路径:</label>
          <select v-model="searchRootPath">
            <option value="SOFTWARE">SOFTWARE</option>
            <option value="SYSTEM">SYSTEM</option>
            <option value="HARDWARE">HARDWARE</option>
          </select>
        </div>
        <div class="modal-actions">
          <button
            @click="searchRegistry"
            class="btn btn-primary"
            :disabled="!searchTerm.trim()"
          >
            <i class="fas fa-search"></i> 搜索
          </button>
          <button @click="showSearchModal = false" class="btn btn-secondary">
            取消
          </button>
        </div>
      </div>
    </div>

    <!-- 搜索结果显示模态框 -->
    <div v-if="showSearchResultsModal" class="modal">
      <div class="modal-content large">
        <h3>搜索结果 ({{ searchResults.length }} 个匹配项)</h3>
        <div class="search-results">
          <div
            v-for="result in searchResults"
            :key="result.path"
            class="search-result-item"
            @click="navigateToSearchResult(result)"
          >
            <div class="result-path">{{ result.path }}</div>
            <div class="result-type">{{ result.match_type }}</div>
            <div class="result-content">{{ result.match_content }}</div>
          </div>
        </div>
        <div class="modal-actions">
          <button
            @click="showSearchResultsModal = false"
            class="btn btn-secondary"
          >
            关闭
          </button>
        </div>
      </div>
    </div>

    <!-- 路径导航模态框 -->
    <div v-if="showBreadcrumbModal" class="modal">
      <div class="modal-content">
        <h3>路径导航</h3>
        <div class="form-group">
          <label>注册表路径:</label>
          <input
            v-model="navigatePath"
            type="text"
            placeholder="输入完整路径，如: SOFTWARE\\Microsoft\\Windows"
          />
        </div>
        <div class="modal-actions">
          <button
            @click="navigateToCustomPath"
            class="btn btn-primary"
            :disabled="!navigatePath.trim()"
          >
            <i class="fas fa-route"></i> 导航
          </button>
          <button
            @click="showBreadcrumbModal = false"
            class="btn btn-secondary"
          >
            取消
          </button>
        </div>
      </div>
    </div>

    <!-- 编辑值模态框 -->
    <div v-if="showEditValueModal" class="modal">
      <div class="modal-content">
        <h3>编辑注册表值</h3>
        <div class="form-group">
          <label>值名称:</label>
          <input v-model="editingValue.name" type="text" />
        </div>
        <div class="form-group">
          <label>值类型:</label>
          <select v-model="editingValue.type">
            <option value="REG_SZ">字符串 (REG_SZ)</option>
            <option value="REG_DWORD">DWORD (REG_DWORD)</option>
            <option value="REG_BINARY">二进制 (REG_BINARY)</option>
            <option value="REG_EXPAND_SZ">可扩展字符串 (REG_EXPAND_SZ)</option>
            <option value="REG_MULTI_SZ">多字符串 (REG_MULTI_SZ)</option>
          </select>
        </div>
        <div class="form-group">
          <label>值数据:</label>
          <textarea
            v-model="editingValue.value"
            rows="3"
            placeholder="输入值数据"
          ></textarea>
        </div>
        <div class="modal-actions">
          <button @click="saveValue" class="btn btn-success">
            <i class="fas fa-save"></i> 保存
          </button>
          <button @click="showEditValueModal = false" class="btn btn-secondary">
            取消
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { registryAPI } from "../services/api.js";
import { apiWrapper } from "../utils/apiErrorHandler.js";
import toast from "../utils/toast.js";

export default {
  name: "Registry",
  data() {
    return {
      currentPath: "SOFTWARE",
      pathHistory: ["SOFTWARE"],
      subKeys: [],
      values: [],
      isLoading: false,
      showCreateKeyModal: false,
      showSearchModal: false,
      showSearchResultsModal: false,
      showBreadcrumbModal: false,
      showEditValueModal: false,
      newKeyPath: "",
      searchTerm: "",
      searchRootPath: "SOFTWARE",
      navigatePath: "",
      editingValue: {},
      searchResults: [],
    };
  },
  computed: {
    canGoBack() {
      return this.pathHistory.length > 1;
    },
  },
  mounted() {
    this.loadRegistryData();
  },
  methods: {
    async loadRegistryData() {
      this.isLoading = true;
      try {
        const result = await apiWrapper.get(
          () => registryAPI.getKeyInfo(this.currentPath),
          "加载注册表数据",
        );

        if (result.success) {
          this.subKeys = result.data.data.sub_keys || [];
          this.values = result.data.data.values || [];
          toast.success("数据加载成功");
        }
      } catch (error) {
        console.error("加载注册表数据失败:", error);
      } finally {
        this.isLoading = false;
      }
    },

    async navigateToPath(path) {
      this.currentPath = path;
      this.pathHistory = [path];
      await this.loadRegistryData();
    },

    async navigateToSubKey(key) {
      const newPath = `${this.currentPath}\\${key}`;
      this.currentPath = newPath;
      this.pathHistory.push(newPath);
      await this.loadRegistryData();
    },

    async navigateToCustomPath() {
      if (this.navigatePath.trim()) {
        this.currentPath = this.navigatePath.trim();
        this.pathHistory = [this.currentPath];
        this.showBreadcrumbModal = false;
        await this.loadRegistryData();
      }
    },

    navigateToSearchResult(result) {
      this.currentPath = result.path;
      this.pathHistory = [result.path];
      this.showSearchResultsModal = false;
      this.loadRegistryData();
    },

    goBack() {
      if (this.pathHistory.length > 1) {
        this.pathHistory.pop();
        this.currentPath = this.pathHistory[this.pathHistory.length - 1];
        this.loadRegistryData();
      }
    },

    copyPath() {
      navigator.clipboard
        .writeText(this.currentPath)
        .then(() => {
          toast.success("路径已复制到剪贴板");
        })
        .catch(() => {
          toast.error("复制失败");
        });
    },

    async createKey() {
      if (!this.newKeyPath.trim()) {
        toast.warning("请输入键路径");
        return;
      }

      try {
        const result = await apiWrapper.post(
          () => registryAPI.createKey({ path: this.newKeyPath }),
          { path: this.newKeyPath },
          "创建注册表键",
        );

        if (result.success) {
          toast.success("注册表键创建成功");
          this.showCreateKeyModal = false;
          this.newKeyPath = "";
          await this.loadRegistryData();
        }
      } catch (error) {
        console.error("创建注册表键失败:", error);
      }
    },

    async searchRegistry() {
      if (!this.searchTerm.trim()) {
        toast.warning("请输入搜索关键词");
        return;
      }

      try {
        const result = await apiWrapper.get(
          () =>
            registryAPI.search({
              search_term: this.searchTerm,
              root_path: this.searchRootPath,
            }),
          "搜索注册表",
        );

        if (result.success) {
          this.searchResults = result.data.data.results || [];
          this.showSearchModal = false;
          this.showSearchResultsModal = true;
          toast.success(`找到 ${this.searchResults.length} 个匹配项`);
        }
      } catch (error) {
        console.error("搜索注册表失败:", error);
      }
    },

    selectValue(value) {
      this.editingValue = { ...value };
      this.showEditValueModal = true;
    },

    editValue(value) {
      this.editingValue = { ...value };
      this.showEditValueModal = true;
    },

    async saveValue() {
      try {
        const result = await apiWrapper.post(
          () =>
            registryAPI.setValue({
              path: this.currentPath,
              name: this.editingValue.name,
              type: this.editingValue.type,
              value: this.editingValue.value,
            }),
          {
            path: this.currentPath,
            name: this.editingValue.name,
            type: this.editingValue.type,
            value: this.editingValue.value,
          },
          "保存注册表值",
        );

        if (result.success) {
          toast.success("注册表值保存成功");
          this.showEditValueModal = false;
          await this.loadRegistryData();
        }
      } catch (error) {
        console.error("保存注册表值失败:", error);
      }
    },

    async deleteValue(value) {
      if (confirm(`确定要删除值 "${value.name}" 吗？`)) {
        try {
          const result = await apiWrapper.call(
            () => registryAPI.deleteValue(this.currentPath, value.name),
            null,
            "删除注册表值",
          );

          if (result.success) {
            toast.success("注册表值删除成功");
            await this.loadRegistryData();
          }
        } catch (error) {
          console.error("删除注册表值失败:", error);
        }
      }
    },

    async deleteKey(key) {
      if (confirm(`确定要删除键 "${key}" 吗？`)) {
        try {
          const result = await apiWrapper.call(
            () => registryAPI.deleteKey(`${this.currentPath}\\${key}`),
            null,
            "删除注册表键",
          );

          if (result.success) {
            toast.success("注册表键删除成功");
            await this.loadRegistryData();
          }
        } catch (error) {
          console.error("删除注册表键失败:", error);
        }
      }
    },

    async getKeyInfo() {
      if (!this.currentPath) {
        toast.warning("请选择注册表路径");
        return;
      }

      try {
        const result = await apiWrapper.get(
          () => registryAPI.getKeyInfo(this.currentPath),
          "获取注册表键信息",
        );

        if (result.success) {
          console.log("注册表键信息:", result.data);
          toast.success("获取注册表键信息成功");
        }
      } catch (error) {
        console.error("获取注册表键信息失败:", error);
      }
    },

    async listKeys() {
      if (!this.currentPath) {
        toast.warning("请选择注册表路径");
        return;
      }

      try {
        const result = await apiWrapper.get(
          () => registryAPI.listKeys(this.currentPath),
          "列出注册表键",
        );

        if (result.success) {
          console.log("注册表键列表:", result.data);
          toast.success("获取注册表键列表成功");
        }
      } catch (error) {
        console.error("获取注册表键列表失败:", error);
      }
    },

    async listValues() {
      if (!this.currentPath) {
        toast.warning("请选择注册表路径");
        return;
      }

      try {
        const result = await apiWrapper.get(
          () => registryAPI.listValues(this.currentPath),
          "列出注册表值",
        );

        if (result.success) {
          console.log("注册表值列表:", result.data);
          toast.success("获取注册表值列表成功");
        }
      } catch (error) {
        console.error("获取注册表值列表失败:", error);
      }
    },

    async refreshData() {
      await this.loadRegistryData();
    },
  },
};
</script>

<style scoped>
.registry-container {
  padding: 20px;
  height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h1 {
  margin: 0;
  color: #333;
}

.actions {
  display: flex;
  gap: 10px;
}

.content {
  display: flex;
  flex: 1;
  gap: 20px;
  position: relative;
}

.sidebar {
  width: 250px;
  background: #f5f5f5;
  border-radius: 8px;
  padding: 15px;
}

.registry-tree {
  height: 100%;
}

.tree-item {
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 4px;
  margin-bottom: 5px;
  transition: background-color 0.2s;
}

.tree-item:hover {
  background-color: #e0e0e0;
}

.tree-item.active {
  background-color: #007bff;
  color: white;
}

.tree-item i {
  margin-right: 8px;
  color: #666;
}

.tree-item.active i {
  color: white;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
}

.loading-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10;
}

.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: #007bff;
}

.path-navigator {
  background: #f8f9fa;
  padding: 10px 15px;
  border-radius: 8px;
  margin-bottom: 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.path-label {
  font-weight: bold;
  margin-right: 10px;
}

.path-value {
  font-family: monospace;
  color: #007bff;
  flex: 1;
}

.path-actions {
  display: flex;
  gap: 5px;
}

.registry-content {
  display: flex;
  flex: 1;
  gap: 20px;
}

.keys-section,
.values-section {
  flex: 1;
  background: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.keys-section h3,
.values-section h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 5px;
}

.key-item,
.value-item {
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.key-item:hover,
.value-item:hover {
  background-color: #f8f9fa;
}

.key-item i {
  margin-right: 8px;
  color: #007bff;
}

.value-name {
  font-weight: bold;
  flex: 1;
}

.value-type {
  color: #666;
  font-size: 0.9em;
  flex: 1;
}

.value-data {
  color: #333;
  flex: 2;
  font-family: monospace;
  word-break: break-all;
}

.value-actions,
.key-actions {
  display: flex;
  gap: 5px;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.empty-state i {
  font-size: 48px;
  margin-bottom: 10px;
  display: block;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  min-width: 400px;
  max-width: 600px;
  max-height: 80vh;
  overflow-y: auto;
}

.modal-content.large {
  min-width: 600px;
  max-width: 800px;
}

.modal-content h3 {
  margin-top: 0;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group textarea {
  resize: vertical;
  min-height: 80px;
}

.modal-actions {
  display: flex;
  gap: 10px;
  justify-content: flex-end;
  margin-top: 20px;
}

.search-results {
  max-height: 400px;
  overflow-y: auto;
}

.search-result-item {
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.search-result-item:hover {
  background-color: #f8f9fa;
}

.result-path {
  font-weight: bold;
  color: #007bff;
  margin-bottom: 5px;
}

.result-type {
  color: #666;
  font-size: 0.9em;
  margin-bottom: 5px;
}

.result-content {
  color: #333;
  font-family: monospace;
  word-break: break-all;
}

.btn {
  padding: 8px 16px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
  transition: background-color 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background-color: #007bff;
  color: white;
}

.btn-success {
  background-color: #28a745;
  color: white;
}

.btn-danger {
  background-color: #dc3545;
  color: white;
}

.btn-secondary {
  background-color: #6c757d;
  color: white;
}

.btn-info {
  background-color: #17a2b8;
  color: white;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.btn:hover:not(:disabled) {
  opacity: 0.8;
}
</style>
