<template>
  <div class="files-container">
    <div class="header">
      <h1>文件管理</h1>
      <div class="actions">
        <button
          @click="refreshData"
          class="btn btn-primary"
          :disabled="isLoading"
        >
          <i class="fas fa-sync" :class="{ 'fa-spin': isLoading }"></i> 刷新
        </button>
        <button @click="showScanModal = true" class="btn btn-success">
          <i class="fas fa-search"></i> 扫描文件
        </button>
      </div>
    </div>

    <div class="content">
      <div class="sidebar">
        <div class="file-actions">
          <h3>文件操作</h3>
          <div class="action-group">
            <label>目录路径:</label>
            <input
              v-model="currentPath"
              type="text"
              placeholder="输入目录路径"
            />
            <button @click="loadFileList" class="btn btn-primary btn-sm">
              <i class="fas fa-folder-open"></i> 浏览
            </button>
          </div>

          <div class="action-group">
            <label>文件路径:</label>
            <input
              v-model="selectedFile"
              type="text"
              placeholder="选择要操作的文件"
            />
            <button @click="getFileInfo" class="btn btn-info btn-sm">
              <i class="fas fa-info"></i> 信息
            </button>
          </div>

          <div class="action-group">
            <label>哈希算法:</label>
            <select v-model="hashAlgorithm">
              <option value="md5">MD5</option>
              <option value="sha1">SHA1</option>
              <option value="sha256">SHA256</option>
            </select>
            <button @click="getFileHash" class="btn btn-secondary btn-sm">
              <i class="fas fa-hashtag"></i> 哈希
            </button>
            <button @click="getFileHashes" class="btn btn-primary btn-sm">
              <i class="fas fa-list"></i> 所有哈希
            </button>
          </div>

          <div class="action-group">
            <label>文件操作:</label>
            <div class="operation-buttons">
              <button @click="copyFile" class="btn btn-warning btn-sm">
                <i class="fas fa-copy"></i> 复制
              </button>
              <button @click="moveFile" class="btn btn-info btn-sm">
                <i class="fas fa-cut"></i> 移动
              </button>
              <button @click="deleteFile" class="btn btn-danger btn-sm">
                <i class="fas fa-trash"></i> 删除
              </button>
            </div>
          </div>

          <div class="action-group">
            <label>哈希验证:</label>
            <input
              v-model="verifyHash"
              type="text"
              placeholder="输入要验证的哈希值"
            />
            <button @click="verifyFileHash" class="btn btn-success btn-sm">
              <i class="fas fa-check"></i> 验证
            </button>
          </div>
        </div>
      </div>

      <div class="main-content">
        <div v-if="isLoading" class="loading-overlay">
          <div class="loading-spinner">
            <i class="fas fa-spinner fa-spin"></i>
            <span>加载中...</span>
          </div>
        </div>

        <div v-else class="file-content">
          <div class="file-list">
            <h3>文件列表 ({{ fileList.length }})</h3>
            <div class="files">
              <div
                v-for="file in fileList"
                :key="file.path"
                class="file-item"
                @click="selectFile(file)"
                :class="{ selected: selectedFile === file.path }"
              >
                <div class="file-icon">
                  <i class="fas fa-file"></i>
                </div>
                <div class="file-info">
                  <div class="file-name">{{ file.name }}</div>
                  <div class="file-path">{{ file.path }}</div>
                  <div class="file-size">{{ formatFileSize(file.size) }}</div>
                </div>
                <div class="file-actions">
                  <button
                    @click.stop="getFileInfo(file.path)"
                    class="btn btn-sm btn-info"
                  >
                    <i class="fas fa-info"></i>
                  </button>
                  <button
                    @click.stop="deleteFile(file.path)"
                    class="btn btn-sm btn-danger"
                  >
                    <i class="fas fa-trash"></i>
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="file-details">
            <h3>文件详情</h3>
            <div v-if="fileInfo" class="details-content">
              <div class="detail-item">
                <span class="label">路径:</span>
                <span class="value">{{ fileInfo.path }}</span>
              </div>
              <div class="detail-item">
                <span class="label">大小:</span>
                <span class="value">{{ formatFileSize(fileInfo.size) }}</span>
              </div>
              <div class="detail-item">
                <span class="label">创建时间:</span>
                <span class="value">{{ fileInfo.create_time }}</span>
              </div>
              <div class="detail-item">
                <span class="label">修改时间:</span>
                <span class="value">{{ fileInfo.modify_time }}</span>
              </div>
              <div class="detail-item">
                <span class="label">权限:</span>
                <span class="value">{{ fileInfo.permissions }}</span>
              </div>
              <div class="detail-item">
                <span class="label">所有者:</span>
                <span class="value">{{ fileInfo.owner }}</span>
              </div>
              <div class="detail-item">
                <span class="label">文件类型:</span>
                <span class="value">{{ fileInfo.file_type }}</span>
              </div>
              <div class="detail-item">
                <span class="label">威胁等级:</span>
                <span
                  class="value"
                  :class="getThreatClass(fileInfo.threat_level)"
                >
                  {{ fileInfo.threat_level }}
                </span>
              </div>
            </div>
            <div v-else class="no-selection">
              <i class="fas fa-file-alt"></i>
              <span>选择文件查看详情</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 文件扫描模态框 -->
    <div v-if="showScanModal" class="modal">
      <div class="modal-content">
        <h3>文件扫描</h3>
        <div class="form-group">
          <label>扫描类型:</label>
          <select v-model="scanType">
            <option value="file">单文件扫描</option>
            <option value="directory">目录扫描</option>
            <option value="buffer">内存缓冲区扫描</option>
          </select>
        </div>

        <div v-if="scanType === 'file'" class="form-group">
          <label>文件路径:</label>
          <input v-model="scanPath" type="text" placeholder="输入文件路径" />
        </div>

        <div v-if="scanType === 'directory'" class="form-group">
          <label>目录路径:</label>
          <input v-model="scanPath" type="text" placeholder="输入目录路径" />
          <div class="checkbox-group">
            <label>
              <input v-model="recursiveScan" type="checkbox" />
              递归扫描
            </label>
            <label>
              最大深度:
              <input v-model="maxDepth" type="number" min="1" max="10" />
            </label>
          </div>
        </div>

        <div v-if="scanType === 'buffer'" class="form-group">
          <label>缓冲区内容:</label>
          <textarea
            v-model="bufferContent"
            rows="5"
            placeholder="输入要扫描的内容"
          ></textarea>
        </div>

        <div class="modal-actions">
          <button
            @click="performScan"
            class="btn btn-success"
            :disabled="!canScan"
          >
            <i class="fas fa-search"></i> 开始扫描
          </button>
          <button @click="showScanModal = false" class="btn btn-secondary">
            取消
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { fileAPI } from "../services/api.js";
import { apiWrapper } from "../utils/apiErrorHandler.js";
import toast from "../utils/toast.js";

export default {
  name: "Files",
  data() {
    return {
      currentPath: "C:\\",
      selectedFile: "",
      fileList: [],
      fileInfo: null,
      isLoading: false,
      hashAlgorithm: "sha256",

      // 扫描相关
      showScanModal: false,
      scanType: "file",
      scanPath: "",
      recursiveScan: true,
      maxDepth: 3,
      bufferContent: "",

      // 文件操作相关
      verifyHash: "",
      copyTargetPath: "",
      moveTargetPath: "",
    };
  },
  computed: {
    canScan() {
      if (this.scanType === "file" || this.scanType === "directory") {
        return this.scanPath.trim();
      }
      return this.bufferContent.trim();
    },
  },
  mounted() {
    this.loadFileList();
  },
  methods: {
    async loadFileList() {
      if (!this.currentPath.trim()) {
        toast.warning("请输入目录路径");
        return;
      }

      this.isLoading = true;
      try {
        const result = await apiWrapper.get(
          () => fileAPI.getList(this.currentPath),
          "加载文件列表",
        );

        if (result.success) {
          this.fileList = result.data.data.files || [];
          toast.success(`加载了 ${this.fileList.length} 个文件`);
        }
      } catch (error) {
        console.error("加载文件列表失败:", error);
      } finally {
        this.isLoading = false;
      }
    },

    selectFile(file) {
      this.selectedFile = file.path;
      this.getFileInfo(file.path);
    },

    async getFileInfo(path) {
      if (!path) {
        toast.warning("请选择文件");
        return;
      }

      try {
        const result = await apiWrapper.get(
          () => fileAPI.getInfo(path),
          "获取文件信息",
        );

        if (result.success) {
          this.fileInfo = result.data.data;
          toast.success("文件信息获取成功");
        }
      } catch (error) {
        console.error("获取文件信息失败:", error);
      }
    },

    async getFileHash() {
      if (!this.selectedFile) {
        toast.warning("请选择文件");
        return;
      }

      try {
        const result = await apiWrapper.get(
          () => fileAPI.getHash(this.selectedFile, this.hashAlgorithm),
          "获取文件哈希",
        );

        if (result.success) {
          const hash = result.data.data.hash;
          toast.success(`${this.hashAlgorithm.toUpperCase()} 哈希: ${hash}`);

          // 复制到剪贴板
          navigator.clipboard.writeText(hash).then(() => {
            toast.info("哈希值已复制到剪贴板");
          });
        }
      } catch (error) {
        console.error("获取文件哈希失败:", error);
      }
    },

    // 新增：获取文件所有哈希值
    async getFileHashes() {
      if (!this.selectedFile) {
        toast.warning("请选择文件");
        return;
      }

      try {
        const result = await apiWrapper.get(
          () => fileAPI.getHashes(this.selectedFile),
          "获取文件所有哈希值",
        );

        if (result.success) {
          const hashes = result.data.data.hashes;
          const hashInfo = `MD5: ${hashes.md5}\nSHA1: ${hashes.sha1}\nSHA256: ${hashes.sha256}`;
          alert(`文件哈希值:\n${hashInfo}`);
          toast.success("获取所有哈希值成功");
        }
      } catch (error) {
        console.error("获取文件所有哈希值失败:", error);
      }
    },

    async deleteFile(path) {
      if (confirm(`确定要删除文件 "${path}" 吗？`)) {
        try {
          const result = await apiWrapper.call(
            () => fileAPI.delete(path),
            null,
            "删除文件",
          );

          if (result.success) {
            toast.success("文件删除成功");
            this.loadFileList();
          }
        } catch (error) {
          console.error("删除文件失败:", error);
        }
      }
    },

    async copyFile() {
      if (!this.selectedFile) {
        toast.warning("请选择要复制的文件");
        return;
      }

      const targetPath = prompt("请输入目标路径:");
      if (!targetPath) return;

      try {
        const result = await apiWrapper.post(
          () =>
            fileAPI.copy({
              source: this.selectedFile,
              destination: targetPath,
            }),
          {
            source: this.selectedFile,
            destination: targetPath,
          },
          "复制文件",
        );

        if (result.success) {
          toast.success("文件复制成功");
          this.loadFileList();
        }
      } catch (error) {
        console.error("复制文件失败:", error);
      }
    },

    async moveFile() {
      if (!this.selectedFile) {
        toast.warning("请选择要移动的文件");
        return;
      }

      const targetPath = prompt("请输入目标路径:");
      if (!targetPath) return;

      try {
        const result = await apiWrapper.post(
          () =>
            fileAPI.move({
              source: this.selectedFile,
              destination: targetPath,
            }),
          {
            source: this.selectedFile,
            destination: targetPath,
          },
          "移动文件",
        );

        if (result.success) {
          toast.success("文件移动成功");
          this.loadFileList();
        }
      } catch (error) {
        console.error("移动文件失败:", error);
      }
    },

    async verifyFileHash() {
      if (!this.selectedFile || !this.verifyHash) {
        toast.warning("请选择文件并输入要验证的哈希值");
        return;
      }

      try {
        const result = await apiWrapper.post(
          () =>
            fileAPI.verifyHash({
              path: this.selectedFile,
              hash: this.verifyHash,
              algorithm: this.hashAlgorithm,
            }),
          {
            path: this.selectedFile,
            hash: this.verifyHash,
            algorithm: this.hashAlgorithm,
          },
          "验证文件哈希",
        );

        if (result.success) {
          const isValid = result.data.data.valid;
          if (isValid) {
            toast.success("哈希验证成功，文件完整性正常");
          } else {
            toast.error("哈希验证失败，文件可能已损坏或被篡改");
          }
        }
      } catch (error) {
        console.error("验证文件哈希失败:", error);
      }
    },

    async performScan() {
      this.isLoading = true;
      try {
        let result;
        switch (this.scanType) {
          case "file":
            result = await apiWrapper.post(
              () => fileAPI.scan({ path: this.scanPath }),
              { path: this.scanPath },
              "扫描文件",
            );
            break;
          case "directory":
            result = await apiWrapper.post(
              () =>
                fileAPI.scanDirectory({
                  path: this.scanPath,
                  recursive: this.recursiveScan,
                  max_depth: this.maxDepth,
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
                }),
              {
                path: this.scanPath,
                recursive: this.recursiveScan,
                max_depth: this.maxDepth,
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
              },
              "扫描目录",
            );
            break;
          case "buffer":
            result = await apiWrapper.post(
              () =>
                fileAPI.scanBuffer({
                  identifier: "web-buffer",
                  data: this.bufferContent,
                }),
              {
                identifier: "web-buffer",
                data: this.bufferContent,
              },
              "扫描缓冲区",
            );
            break;
        }

        if (result.success) {
          toast.success("扫描完成");
          this.showScanModal = false;
        }
      } catch (error) {
        console.error("扫描失败:", error);
      } finally {
        this.isLoading = false;
      }
    },

    refreshData() {
      this.loadFileList();
    },

    formatFileSize(bytes) {
      if (bytes === 0) return "0 B";
      const k = 1024;
      const sizes = ["B", "KB", "MB", "GB", "TB"];
      const i = Math.floor(Math.log(bytes) / Math.log(k));
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    },

    getThreatClass(level) {
      switch (level) {
        case "high":
          return "threat-high";
        case "medium":
          return "threat-medium";
        case "low":
          return "threat-low";
        default:
          return "threat-safe";
      }
    },
  },
};
</script>

<style scoped>
.files-container {
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
}

.sidebar {
  width: 300px;
  background: #f5f5f5;
  border-radius: 8px;
  padding: 15px;
}

.file-actions h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 5px;
}

.action-group {
  margin-bottom: 15px;
}

.action-group label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
  color: #333;
}

.action-group input,
.action-group select {
  width: 100%;
  padding: 8px 12px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
  margin-bottom: 5px;
}

.checkbox-group {
  display: flex;
  gap: 15px;
  align-items: center;
  margin-top: 5px;
}

.checkbox-group label {
  display: flex;
  align-items: center;
  gap: 5px;
  font-weight: normal;
  margin: 0;
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

.file-content {
  display: flex;
  flex: 1;
  gap: 20px;
}

.file-list {
  flex: 2;
  background: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.file-list h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 5px;
}

.files {
  max-height: 500px;
  overflow-y: auto;
}

.file-item {
  display: flex;
  align-items: center;
  padding: 10px;
  border: 1px solid #e0e0e0;
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.file-item:hover {
  background-color: #f8f9fa;
}

.file-item.selected {
  background-color: #e3f2fd;
  border-color: #007bff;
}

.file-icon {
  margin-right: 10px;
  color: #007bff;
  font-size: 18px;
}

.file-info {
  flex: 1;
}

.file-name {
  font-weight: bold;
  color: #333;
  margin-bottom: 2px;
}

.file-path {
  font-size: 0.9em;
  color: #666;
  margin-bottom: 2px;
}

.file-size {
  font-size: 0.8em;
  color: #999;
}

.file-actions {
  display: flex;
  gap: 5px;
}

.file-details {
  flex: 1;
  background: white;
  border-radius: 8px;
  padding: 15px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.file-details h3 {
  margin-top: 0;
  margin-bottom: 15px;
  color: #333;
  border-bottom: 2px solid #007bff;
  padding-bottom: 5px;
}

.details-content {
  space-y: 10px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.detail-item:last-child {
  border-bottom: none;
}

.detail-item .label {
  font-weight: bold;
  color: #666;
}

.detail-item .value {
  color: #333;
  font-family: monospace;
}

.threat-high {
  color: #dc3545;
  font-weight: bold;
}

.threat-medium {
  color: #ffc107;
  font-weight: bold;
}

.threat-low {
  color: #28a745;
  font-weight: bold;
}

.threat-safe {
  color: #28a745;
  font-weight: bold;
}

.no-selection {
  text-align: center;
  padding: 40px 20px;
  color: #999;
}

.no-selection i {
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

.btn-warning {
  background-color: #ffc107;
  color: #212529;
}

.btn-success {
  background-color: #28a745;
  color: white;
}

.btn-sm {
  padding: 4px 8px;
  font-size: 12px;
}

.operation-buttons {
  display: flex;
  gap: 5px;
  flex-wrap: wrap;
}

.btn:hover:not(:disabled) {
  opacity: 0.8;
}
</style>
