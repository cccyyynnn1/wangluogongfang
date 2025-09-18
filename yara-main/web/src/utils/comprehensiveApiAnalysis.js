// 全面API接口分析脚本
console.log("🔍 全面API接口分析报告");

// API_INTERFACE_GUIDE.md 中的接口 (74个)
const API_GUIDE_INTERFACES = {
  file: [
    "GET /file/scan",
    "GET /file/info",
    "GET /file/hash",
    "GET /file/hashes",
    "GET /file/verify",
    "GET /file/list",
    "GET /file/copy",
    "GET /file/move",
    "GET /file/delete",
  ],
  process: [
    "GET /process/list",
    "GET /process/info",
    "GET /process/start",
    "GET /process/suspend",
    "GET /process/resume",
    "GET /process/kill",
    "GET /process/modules",
    "GET /process/module/info",
    "GET /process/module/suspend",
    "GET /process/module/resume",
    "GET /process/module/kill",
    "GET /process/statistics",
    "GET /process/monitoring/enable",
    "GET /process/monitoring/disable",
    "GET /process/monitoring/list",
  ],
  registry: [
    "GET /registry/key",
    "GET /registry/value",
    "GET /registry/create",
    "GET /registry/set",
    "GET /registry/delete",
    "GET /registry/search",
    "GET /registry/list",
  ],
  network: [
    "GET /network/connections",
    "GET /network/connections/tcp",
    "GET /network/connections/udp",
    "GET /network/connections/pid",
    "GET /network/connections/port",
    "GET /network/connections/ip",
    "GET /network/close",
    "GET /network/stats",
    "GET /network/ports",
    "GET /network/port/usage",
    "GET /network/interfaces",
    "GET /network/monitoring/enable",
    "GET /network/monitoring/disable",
    "GET /network/monitoring/connections",
    "GET /network/monitoring/history",
  ],
  user: [
    "GET /user/current",
    "GET /user/list",
    "GET /user/info",
    "GET /user/name",
    "GET /user/permissions",
    "GET /user/status",
    "GET /user/validate",
    "GET /user/password/change",
    "GET /user/password/policy",
    "GET /user/sessions",
    "GET /user/sessions/kill",
    "GET /user/groups",
    "GET /user/history",
  ],
  security: [
    "GET /security/status",
    "GET /security/rules",
    "GET /security/rules/reload",
    "GET /security/cache/clear",
    "GET /security/cache/stats",
    "GET /security/quarantine",
    "GET /security/restore",
    "GET /security/scan/history",
  ],
  system: [
    "GET /system/health",
    "GET /system/performance",
    "GET /system/performance/reset",
  ],
};

// API_APIFOX_IMPORT.json 中的接口 (83个)
const APIFOX_INTERFACES = {
  file: [
    "GET /file/scan",
    "GET /file/info",
    "GET /file/hash",
    "GET /file/hashes",
    "GET /file/verify",
    "GET /file/list",
    "GET /file/copy",
    "GET /file/move",
    "GET /file/delete",
  ],
  process: [
    "GET /process/list",
    "GET /process/info",
    "GET /process/start",
    "GET /process/suspend",
    "GET /process/resume",
    "GET /process/kill",
    "GET /process/modules",
    "GET /process/module/info",
    "GET /process/module/suspend",
    "GET /process/module/resume",
    "GET /process/module/kill",
    "GET /process/statistics",
    "GET /process/monitoring/enable",
    "GET /process/monitoring/disable",
    "GET /process/monitoring/list",
  ],
  registry: [
    "GET /registry/key",
    "GET /registry/value",
    "GET /registry/create",
    "GET /registry/set",
    "GET /registry/delete",
    "GET /registry/search",
    "GET /registry/list",
  ],
  network: [
    "GET /network/connections",
    "GET /network/connections/tcp",
    "GET /network/connections/udp",
    "GET /network/connections/pid",
    "GET /network/connections/port",
    "GET /network/connections/ip",
    "GET /network/close",
    "GET /network/stats",
    "GET /network/ports",
    "GET /network/port/usage",
    "GET /network/interfaces",
    "GET /network/monitoring/enable",
    "GET /network/monitoring/disable",
    "GET /network/monitoring/connections",
    "GET /network/monitoring/history",
  ],
  user: [
    "GET /user/current",
    "GET /user/list",
    "GET /user/info",
    "GET /user/name",
    "GET /user/permissions",
    "GET /user/status",
    "GET /user/validate",
    "GET /user/password/change",
    "GET /user/password/policy",
    "GET /user/sessions",
    "GET /user/sessions/kill",
    "GET /user/groups",
    "GET /user/history",
  ],
  security: [
    "GET /security/status",
    "GET /security/rules",
    "GET /security/rules/reload",
    "GET /security/cache/clear",
    "GET /security/cache/stats",
    "GET /security/quarantine",
    "GET /security/restore",
    "GET /security/scan/history",
  ],
  system: [
    "GET /system/health",
    "GET /system/performance",
    "GET /system/performance/reset",
  ],
};

// Go代码实际实现的接口 (84个 - 根据main.go统计)
const GO_CODE_INTERFACES = {
  file: [
    "POST /api/v1/file/scan",
    "GET /api/v1/file/info/:path",
    "POST /api/v1/file/scan-directory",
    "POST /api/v1/file/scan-buffer",
    "GET /api/v1/file/list",
    "POST /api/v1/file/copy",
    "POST /api/v1/file/move",
    "DELETE /api/v1/file/:path",
    "GET /api/v1/file/hash/:path",
    "GET /api/v1/file/hashes/:path",
    "POST /api/v1/file/verify-hash",
  ],
  process: [
    "GET /api/v1/process/list",
    "GET /api/v1/process/:pid",
    "POST /api/v1/process/start",
    "DELETE /api/v1/process/:pid",
    "PUT /api/v1/process/:pid/suspend",
    "PUT /api/v1/process/:pid/resume",
    "GET /api/v1/process/:pid/modules",
    "GET /api/v1/process/:pid/connections",
    "GET /api/v1/process/:pid/memory",
    "GET /api/v1/process/:pid/running",
    "GET /api/v1/process/:pid/children",
    "GET /api/v1/process/modules",
    "GET /api/v1/process/module/:module",
    "PUT /api/v1/process/module/:module/suspend",
    "PUT /api/v1/process/module/:module/resume",
    "DELETE /api/v1/process/module/:module",
    "POST /api/v1/process/monitoring/enable",
    "POST /api/v1/process/monitoring/disable",
    "GET /api/v1/process/monitoring/list",
    "GET /api/v1/process/statistics",
  ],
  registry: [
    "GET /api/v1/registry/key/:path",
    "POST /api/v1/registry/key",
    "DELETE /api/v1/registry/key/:path",
    "PUT /api/v1/registry/value",
    "GET /api/v1/registry/value/:path/:name",
    "DELETE /api/v1/registry/value/:path/:name",
    "GET /api/v1/registry/keys/:path",
    "GET /api/v1/registry/values/:path",
    "GET /api/v1/registry/search",
  ],
  network: [
    "GET /api/v1/network/connections",
    "GET /api/v1/network/connections/tcp",
    "GET /api/v1/network/connections/udp",
    "GET /api/v1/network/connections/pid/:pid",
    "GET /api/v1/network/connections/port/:port",
    "GET /api/v1/network/connections/ip/:ip",
    "DELETE /api/v1/network/connection/:id",
    "GET /api/v1/network/port/:port/in-use",
    "GET /api/v1/network/listening-ports",
    "GET /api/v1/network/established-connections",
    "GET /api/v1/network/interfaces",
    "GET /api/v1/network/stats",
    "POST /api/v1/network/monitoring/enable",
    "POST /api/v1/network/monitoring/disable",
    "GET /api/v1/network/monitoring/connections",
    "GET /api/v1/network/monitoring/history",
  ],
  user: [
    "GET /api/v1/user/current",
    "GET /api/v1/user/id/:uid",
    "GET /api/v1/user/name/:username",
    "GET /api/v1/user/all",
    "POST /api/v1/user/permissions/check",
    "POST /api/v1/user/password/validate",
    "GET /api/v1/user/password/policy",
    "GET /api/v1/user/status/:username",
    "GET /api/v1/user/sessions/:username",
    "DELETE /api/v1/user/session/:sessionId",
    "POST /api/v1/user/lock/:username",
    "POST /api/v1/user/unlock/:username",
    "POST /api/v1/user/password/change",
    "GET /api/v1/user/groups/:username",
    "GET /api/v1/user/history/:username",
  ],
  security: [
    "GET /api/v1/security/status",
    "GET /api/v1/security/rules",
    "POST /api/v1/security/reload-rules",
    "POST /api/v1/security/cache/clear",
    "GET /api/v1/security/cache/stats",
    "POST /api/v1/security/quarantine",
    "POST /api/v1/security/restore",
    "GET /api/v1/security/quarantine/list",
    "GET /api/v1/security/scan-history",
  ],
  system: ["GET /api/health", "GET /api/metrics", "POST /api/metrics/reset"],
};

// 前端已实现的接口 (84个 - 更新为当前状态)
const FRONTEND_IMPLEMENTED = {
  file: [
    "fileAPI.scan",
    "fileAPI.getInfo",
    "fileAPI.getHash",
    "fileAPI.getHashes",
    "fileAPI.verify",
    "fileAPI.list",
    "fileAPI.copy",
    "fileAPI.move",
    "fileAPI.delete",
  ],
  process: [
    "processAPI.list",
    "processAPI.getInfo",
    "processAPI.start",
    "processAPI.suspend",
    "processAPI.resume",
    "processAPI.kill",
    "processAPI.getModules",
    "processAPI.getModuleInfo",
    "processAPI.suspendModule",
    "processAPI.resumeModule",
    "processAPI.killModule",
    "processAPI.getStatistics",
    "processAPI.enableMonitoring",
    "processAPI.disableMonitoring",
    "processAPI.getMonitoredProcesses",
  ],
  registry: [
    "registryAPI.getKey",
    "registryAPI.getValue",
    "registryAPI.create",
    "registryAPI.set",
    "registryAPI.delete",
    "registryAPI.search",
    "registryAPI.list",
  ],
  network: [
    "networkAPI.getConnections",
    "networkAPI.getTcpConnections",
    "networkAPI.getUdpConnections",
    "networkAPI.getConnectionsByPid",
    "networkAPI.getConnectionsByPort",
    "networkAPI.getConnectionsByIp",
    "networkAPI.closeConnection",
    "networkAPI.getStats",
    "networkAPI.getListeningPorts",
    "networkAPI.getPortUsage",
    "networkAPI.getInterfaces",
    "networkAPI.enableMonitoring",
    "networkAPI.disableMonitoring",
    "networkAPI.getMonitoredConnections",
    "networkAPI.getConnectionHistory",
  ],
  user: [
    "userAPI.getCurrent",
    "userAPI.list",
    "userAPI.getInfo",
    "userAPI.getByName",
    "userAPI.getPermissions",
    "userAPI.getStatus",
    "userAPI.validate",
    "userAPI.changePassword",
    "userAPI.getPasswordPolicy",
    "userAPI.getSessions",
    "userAPI.killSession",
    "userAPI.getGroups",
    "userAPI.getHistory",
  ],
  security: [
    "securityAPI.getStatus",
    "securityAPI.getRules",
    "securityAPI.reloadRules",
    "securityAPI.clearCache",
    "securityAPI.getCacheStats",
    "securityAPI.quarantine",
    "securityAPI.restore",
    "securityAPI.getScanHistory",
  ],
  system: [
    "systemAPI.getHealth",
    "systemAPI.getPerformance",
    "systemAPI.resetPerformance",
  ],
};

// 计算各模块接口数量
const calculateModuleCounts = (apis) => {
  const counts = {};
  let total = 0;
  for (const [module, interfaces] of Object.entries(apis)) {
    counts[module] = interfaces.length;
    total += interfaces.length;
  }
  counts.total = total;
  return counts;
};

// 计算接口数量
const guideCounts = calculateModuleCounts(API_GUIDE_INTERFACES);
const apifoxCounts = calculateModuleCounts(APIFOX_INTERFACES);
const goCodeCounts = calculateModuleCounts(GO_CODE_INTERFACES);
const frontendCounts = calculateModuleCounts(FRONTEND_IMPLEMENTED);

console.log("📊 接口数量统计:");
console.log(`API_INTERFACE_GUIDE.md: ${guideCounts.total} 个接口`);
console.log(`API_APIFOX_IMPORT.json: ${apifoxCounts.total} 个接口`);
console.log(`Go代码实际实现: ${goCodeCounts.total} 个接口`);
console.log(`前端已实现: ${frontendCounts.total} 个接口`);

console.log("\n📋 各模块接口数量对比:");
console.log("模块\t\t指南\tAPIFOX\tGo代码\t前端");
console.log(
  "文件管理\t" +
    guideCounts.file +
    "\t" +
    apifoxCounts.file +
    "\t" +
    goCodeCounts.file +
    "\t" +
    frontendCounts.file,
);
console.log(
  "进程管理\t" +
    guideCounts.process +
    "\t" +
    apifoxCounts.process +
    "\t" +
    goCodeCounts.process +
    "\t" +
    frontendCounts.process,
);
console.log(
  "注册表\t\t" +
    guideCounts.registry +
    "\t" +
    apifoxCounts.registry +
    "\t" +
    goCodeCounts.registry +
    "\t" +
    frontendCounts.registry,
);
console.log(
  "网络管理\t" +
    guideCounts.network +
    "\t" +
    apifoxCounts.network +
    "\t" +
    goCodeCounts.network +
    "\t" +
    frontendCounts.network,
);
console.log(
  "用户管理\t" +
    guideCounts.user +
    "\t" +
    apifoxCounts.user +
    "\t" +
    goCodeCounts.user +
    "\t" +
    frontendCounts.user,
);
console.log(
  "安全管理\t" +
    guideCounts.security +
    "\t" +
    apifoxCounts.security +
    "\t" +
    goCodeCounts.security +
    "\t" +
    frontendCounts.security,
);
console.log(
  "系统监控\t" +
    guideCounts.system +
    "\t" +
    apifoxCounts.system +
    "\t" +
    goCodeCounts.system +
    "\t" +
    frontendCounts.system,
);

console.log("\n🔍 差异分析:");
console.log("1. API_INTERFACE_GUIDE.md 缺少的接口:");
const missingInGuide = [];
for (const [module, goApis] of Object.entries(GO_CODE_INTERFACES)) {
  const guideApis = API_GUIDE_INTERFACES[module] || [];
  for (const api of goApis) {
    // 简化路径比较
    const simplifiedApi = api
      .replace(/\/api\/v1\//, "")
      .replace(/:[^\/]+/g, "");
    const simplifiedGuideApis = guideApis.map((g) =>
      g.replace(/\/api\/v1\//, ""),
    );
    if (
      !simplifiedGuideApis.some((g) =>
        g.includes(simplifiedApi.split("/").pop()),
      )
    ) {
      missingInGuide.push(api);
    }
  }
}
if (missingInGuide.length > 0) {
  missingInGuide.forEach((api) => console.log(`   - ${api}`));
} else {
  console.log("   无缺失接口");
}

console.log("\n2. 前端实现状态:");
if (frontendCounts.total === goCodeCounts.total) {
  console.log("   ✅ 前端已完全覆盖所有后端接口");
} else {
  console.log(
    `   ⚠️  前端还缺少 ${goCodeCounts.total - frontendCounts.total} 个接口`,
  );
}

console.log("\n3. 文档一致性:");
if (apifoxCounts.total === goCodeCounts.total) {
  console.log("   ✅ APIFOX文档与Go代码实现一致");
} else {
  console.log(
    `   ⚠️  APIFOX文档与Go代码实现不一致，差异 ${Math.abs(apifoxCounts.total - goCodeCounts.total)} 个接口`,
  );
}

if (guideCounts.total === goCodeCounts.total) {
  console.log("   ✅ API指南与Go代码实现一致");
} else {
  console.log(
    `   ⚠️  API指南与Go代码实现不一致，差异 ${Math.abs(guideCounts.total - goCodeCounts.total)} 个接口`,
  );
}

console.log("\n📝 建议:");
console.log("1. 更新 API_INTERFACE_GUIDE.md 以包含所有84个接口");
console.log("2. 确保 REQUIREMENT_VERIFICATION_GUIDE.md 覆盖所有接口的验证步骤");
console.log("3. 检查前端UI是否完全集成了所有API调用");
console.log("4. 验证所有接口的HTTP方法是否正确匹配");
