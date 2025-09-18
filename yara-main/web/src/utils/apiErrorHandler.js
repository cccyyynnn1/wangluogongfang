import toast from "./toast.js";

// API错误处理工具
export class ApiErrorHandler {
  static handle(error, context = "") {
    console.error(`API Error in ${context}:`, error);

    let message = "操作失败";

    if (error.response) {
      // 服务器响应了错误状态码
      const status = error.response.status;
      const data = error.response.data;

      switch (status) {
        case 400:
          message = data?.message || "请求参数错误";
          break;
        case 401:
          message = "未授权访问";
          break;
        case 403:
          message = "权限不足";
          break;
        case 404:
          message = data?.message || "资源不存在";
          break;
        case 500:
          message = data?.message || "服务器内部错误";
          break;
        case 502:
          message = "网关错误";
          break;
        case 503:
          message = "服务不可用";
          break;
        default:
          message = data?.message || `请求失败 (${status})`;
      }
    } else if (error.request) {
      // 请求已发出但没有收到响应
      message = "网络连接失败，请检查网络设置";
    } else {
      // 请求配置出错
      message = error.message || "请求配置错误";
    }

    toast.error(message);
    return message;
  }

  static async wrap(apiCall, context = "") {
    try {
      const response = await apiCall();
      return { success: true, data: response.data };
    } catch (error) {
      const message = this.handle(error, context);
      return { success: false, message, error };
    }
  }

  static validateResponse(response, context = "") {
    if (!response || !response.data) {
      const message = "响应数据格式错误";
      toast.error(message);
      return { success: false, message };
    }

    if (response.data.code !== 200) {
      const message = response.data.message || "操作失败";
      toast.error(message);
      return { success: false, message };
    }

    return { success: true, data: response.data };
  }
}

// 通用API调用包装器
export const apiWrapper = {
  async call(apiFunction, params = null, context = "") {
    try {
      const response = params ? await apiFunction(params) : await apiFunction();
      return ApiErrorHandler.validateResponse(response, context);
    } catch (error) {
      return { success: false, ...ApiErrorHandler.handle(error, context) };
    }
  },

  async get(apiFunction, context = "") {
    return this.call(apiFunction, null, context);
  },

  async post(apiFunction, data, context = "") {
    return this.call(apiFunction, data, context);
  },
};

export default ApiErrorHandler;
