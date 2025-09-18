// API 控制台脚本
const E = {
  panels: document.querySelectorAll(".panel"),
  methodBtns: document.querySelectorAll(".method-btn"),
  baseUrlInput: document.getElementById("base-url"),
  // file panel
  fileScanInput: document.getElementById("file-scan-input"),
  btnFileScan: document.getElementById("btn-file-scan"),
  btnGetFileList: document.getElementById("btn-get-file-list"),
  filePathInput: document.getElementById("file-path"),
  btnFileInfo: document.getElementById("btn-file-info"),
  btnFileHash: document.getElementById("btn-file-hash"),
  btnFileDelete: document.getElementById("btn-file-delete"),
  // network
  btnNetConns: document.getElementById("btn-net-conns"),
  btnNetStats: document.getElementById("btn-net-stats"),
  btnListeningPorts: document.getElementById("btn-listening-ports"),
  netPid: document.getElementById("net-pid"),
  btnNetByPid: document.getElementById("btn-net-by-pid"),
  // process
  btnProcessList: document.getElementById("btn-process-list"),
  procCmd: document.getElementById("proc-cmd"),
  procArgs: document.getElementById("proc-args"),
  btnProcStart: document.getElementById("btn-proc-start"),
  // registry
  regPath: document.getElementById("reg-path"),
  btnRegKeys: document.getElementById("btn-reg-keys"),
  btnRegValues: document.getElementById("btn-reg-values"),
  btnRegSearch: document.getElementById("btn-reg-search"),
  // security
  secFileInput: document.getElementById("sec-file-input"),
  btnSecScan: document.getElementById("btn-sec-scan"),
  btnGetRules: document.getElementById("btn-get-rules"),
  btnScanHistory: document.getElementById("btn-scan-history"),
  // user
  btnUserCurrent: document.getElementById("btn-user-current"),
  btnUserAll: document.getElementById("btn-user-all"),
  userName: document.getElementById("user-name"),
  btnUserHistory: document.getElementById("btn-user-history"),
  // results
  resultsSection: document.getElementById("api-results"),
  responsePre: document.getElementById("api-response"),
  statsDiv: document.getElementById("api-stats"),
  btnCopy: document.getElementById("btn-copy-response"),
  btnClear: document.getElementById("btn-clear-response"),
  // error
  apiError: document.getElementById("api-error"),
  apiErrorText: document.getElementById("api-error-text"),
  apiErrorClose: document.getElementById("api-error-close"),
};

// 辅助函数
function baseUrl() {
  return (E.baseUrlInput && E.baseUrlInput.value) || "/api/v1";
}
function showResult(obj) {
  E.responsePre.textContent = JSON.stringify(obj, null, 2);
  E.resultsSection.classList.remove("hidden");
}
function showError(msg) {
  E.apiErrorText.textContent = msg;
  E.apiError.classList.remove("hidden");
}
function hideError() {
  E.apiError.classList.add("hidden");
}

// 面板切换
E.methodBtns.forEach((btn) => {
  btn.addEventListener("click", () => {
    E.methodBtns.forEach((b) => b.classList.remove("active"));
    btn.classList.add("active");
    const panel = btn.dataset.panel;
    E.panels.forEach((p) => p.classList.add("hidden"));
    document.getElementById(panel).classList.remove("hidden");
    hideError();
    E.resultsSection.classList.add("hidden");
  });
});

E.apiErrorClose.addEventListener("click", hideError);
E.btnCopy.addEventListener("click", () => {
  navigator.clipboard
    .writeText(E.responsePre.textContent)
    .then(() => alert("已复制到剪贴板"));
});
E.btnClear.addEventListener("click", () => {
  E.responsePre.textContent = "";
  E.resultsSection.classList.add("hidden");
});

// 网络请求工具
async function doGet(path) {
  hideError();
  try {
    const res = await fetch(baseUrl() + path, { method: "GET" });
    const json = await res.json();
    showResult(json);
    return json;
  } catch (err) {
    showError(err.message || err);
    console.error(err);
  }
}
async function doPostJson(path, body) {
  hideError();
  try {
    const res = await fetch(baseUrl() + path, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(body),
    });
    const json = await res.json();
    showResult(json);
    return json;
  } catch (err) {
    showError(err.message || err);
    console.error(err);
  }
}
async function doPostFormData(path, formData) {
  hideError();
  try {
    const res = await fetch(baseUrl() + path, {
      method: "POST",
      body: formData,
    });
    const json = await res.json();
    showResult(json);
    return json;
  } catch (err) {
    showError(err.message || err);
    console.error(err);
  }
}
async function doDelete(path) {
  hideError();
  try {
    const res = await fetch(baseUrl() + path, { method: "DELETE" });
    const json = await res.json();
    showResult(json);
    return json;
  } catch (err) {
    showError(err.message || err);
    console.error(err);
  }
}

// 文件事件
E.btnGetFileList.addEventListener("click", () => doGet("/file/list"));
E.btnFileScan.addEventListener("click", async () => {
  const f = E.fileScanInput.files[0];
  if (!f) {
    showError("请先选择要扫描的文件");
    return;
  }
  const fd = new FormData();
  fd.append("file", f);
  await doPostFormData("/file/scan", fd);
});
E.btnFileInfo.addEventListener("click", () => {
  const p = E.filePathInput.value.trim();
  if (!p) {
    showError("请输入路径");
    return;
  }
  // encode path as query param fullPath 如果路径包含反斜杠，使用fullPath查询
  if (p.includes("\\") || p.includes("/")) {
    doGet(
      `/file/info/${encodeURIComponent(p)}?fullPath=${encodeURIComponent(p)}`,
    );
  } else {
    doGet(`/file/info/${encodeURIComponent(p)}`);
  }
});
E.btnFileHash.addEventListener("click", () => {
  const p = E.filePathInput.value.trim();
  if (!p) {
    showError("请输入路径");
    return;
  }
  if (p.includes("\\") || p.includes("/")) {
    doGet(
      `/file/hash/${encodeURIComponent(p)}?fullPath=${encodeURIComponent(p)}`,
    );
  } else {
    doGet(`/file/hash/${encodeURIComponent(p)}`);
  }
});
E.btnFileDelete.addEventListener("click", () => {
  const p = E.filePathInput.value.trim();
  if (!p) {
    showError("请输入路径");
    return;
  }
  if (confirm("确认删除文件？该操作不可撤销。")) {
    if (p.includes("\\") || p.includes("/")) {
      doDelete(
        `/file/${encodeURIComponent(p)}?fullPath=${encodeURIComponent(p)}`,
      );
    } else {
      doDelete(`/file/${encodeURIComponent(p)}`);
    }
  }
});

// 网络事件
E.btnNetConns.addEventListener("click", () => doGet("/network/connections"));
E.btnNetStats.addEventListener("click", () => doGet("/network/stats"));
E.btnListeningPorts.addEventListener("click", () =>
  doGet("/network/listening-ports"),
);
E.btnNetByPid.addEventListener("click", () => {
  const pid = E.netPid.value.trim();
  if (!pid) {
    showError("请输入 PID");
    return;
  }
  doGet(`/network/connections/pid/${encodeURIComponent(pid)}`);
});

// 进程事件
E.btnProcessList.addEventListener("click", () => doGet("/process/list"));
E.btnProcStart.addEventListener("click", () => {
  const cmd = E.procCmd.value.trim();
  if (!cmd) {
    showError("请输入要启动的命令");
    return;
  }
  const args = E.procArgs.value.trim();
  const body = {
    command: cmd,
    args: args
      ? args
          .split(",")
          .map((s) => s.trim())
          .filter(Boolean)
      : [],
  };
  doPostJson("/process/start", body);
});

// 注册表事件
E.btnRegKeys.addEventListener("click", () => {
  const path = E.regPath.value.trim();
  if (!path) {
    showError("请输入注册表路径");
    return;
  }
  doGet(`/registry/keys/${encodeURIComponent(path)}`);
});
E.btnRegValues.addEventListener("click", () => {
  const path = E.regPath.value.trim();
  if (!path) {
    showError("请输入注册表路径");
    return;
  }
  doGet(`/registry/values/${encodeURIComponent(path)}`);
});
E.btnRegSearch.addEventListener("click", () => {
  const path = E.regPath.value.trim();
  if (!path) {
    showError("请输入搜索关键词或路径");
    return;
  }
  // POST 搜索给后端更灵活
  doPostJson("/registry/search", { pattern: path });
});

// 安全事件
E.btnSecScan.addEventListener("click", async () => {
  const f = E.secFileInput.files[0];
  if (!f) {
    showError("请选择要扫描的文件");
    return;
  }
  const fd = new FormData();
  fd.append("file", f);
  await doPostFormData("/security/scan", fd);
});
E.btnGetRules.addEventListener("click", () => doGet("/security/rules"));
E.btnScanHistory.addEventListener("click", () =>
  doGet("/security/scan-history"),
);

// 用户事件
E.btnUserCurrent.addEventListener("click", () => doGet("/user/current"));
E.btnUserAll.addEventListener("click", () => doGet("/user/all"));
E.btnUserHistory.addEventListener("click", () => {
  const name = E.userName.value.trim();
  if (!name) {
    showError("请输入用户名");
    return;
  }
  doGet(`/user/history/${encodeURIComponent(name)}`);
});

// 初始化：隐藏所有 panel 除第一个
document.addEventListener("DOMContentLoaded", () => {
  E.panels.forEach((p, i) => {
    if (i !== 0) p.classList.add("hidden");
  });
});

// 处理跨域或非200错误时的显示（fetch 已在上层 catch 捕获网络错误，但HTTP错误需要额外处理）
(async function patchFetch() {
  const _fetch = window.fetch;
  window.fetch = async function (input, init) {
    const res = await _fetch(input, init);
    const ct = res.headers.get("content-type") || "";
    if (!res.ok) {
      let text = await res.text();
      try {
        const json = JSON.parse(text);
        showError(json.message || text);
        showResult(json);
      } catch (e) {
        showError(text);
      }
      return res;
    }
    if (ct.includes("application/json")) return res;
    return res;
  };
})();
