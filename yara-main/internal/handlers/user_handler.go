package handlers

import (
	"strconv"
	"yara-security-service/internal/models"
	"yara-security-service/internal/services"
	"yara-security-service/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *services.UserService
	logger      *logrus.Logger
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *services.UserService, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetCurrentUser 获取当前用户信息
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userInfo, err := h.userService.GetCurrentUser()
	if err != nil {
		h.logger.Errorf("获取当前用户信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取当前用户信息失败")
		return
	}

	utils.InfoSuccessResponse(c, userInfo)
}

// GetUserByID 根据用户ID获取用户信息
func (h *UserHandler) GetUserByID(c *gin.Context) {
	uid := c.Param("uid")
	if uid == "" {
		utils.BadRequestResponse(c, "用户ID不能为空")
		return
	}

	userInfo, err := h.userService.GetUserByID(uid)
	if err != nil {
		h.logger.Errorf("获取用户信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取用户信息失败")
		return
	}

	utils.InfoSuccessResponse(c, userInfo)
}

// GetUserByName 根据用户名获取用户信息
func (h *UserHandler) GetUserByName(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	userInfo, err := h.userService.GetUserByName(username)
	if err != nil {
		h.logger.Errorf("获取用户信息失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取用户信息失败")
		return
	}

	utils.InfoSuccessResponse(c, userInfo)
}

// GetAllUsers 获取所有用户
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.logger.Errorf("获取所有用户失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取所有用户失败")
		return
	}

	utils.ListSuccessResponse(c, users, len(users))
}

// CheckUserPermissions 检查用户权限
func (h *UserHandler) CheckUserPermissions(c *gin.Context) {
	var request models.CheckPermissionsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	permissions, err := h.userService.CheckUserPermissions(request.Username, request.Permissions)
	if err != nil {
		h.logger.Errorf("检查用户权限失败: %v", err)
		utils.InternalServerErrorResponse(c, "检查用户权限失败")
		return
	}

	utils.InfoSuccessResponse(c, permissions)
}

// ValidatePassword 验证密码
func (h *UserHandler) ValidatePassword(c *gin.Context) {
	var request models.ValidatePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	isValid, err := h.userService.ValidatePassword(request.Username, request.Password)
	if err != nil {
		h.logger.Errorf("验证密码失败: %v", err)
		utils.InternalServerErrorResponse(c, "验证密码失败")
		return
	}

	utils.InfoSuccessResponse(c, gin.H{"is_valid": isValid})
}

// GetPasswordPolicy 获取密码策略
func (h *UserHandler) GetPasswordPolicy(c *gin.Context) {
	policy, err := h.userService.GetPasswordPolicy()
	if err != nil {
		h.logger.Errorf("获取密码策略失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取密码策略失败")
		return
	}

	utils.InfoSuccessResponse(c, policy)
}

// CheckAccountStatus 检查账户状态
func (h *UserHandler) CheckAccountStatus(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	status, err := h.userService.CheckAccountStatus(username)
	if err != nil {
		h.logger.Errorf("检查账户状态失败: %v", err)
		utils.InternalServerErrorResponse(c, "检查账户状态失败")
		return
	}

	utils.StatusSuccessResponse(c, status)
}

// GetUserSessions 获取用户会话
func (h *UserHandler) GetUserSessions(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	sessions, err := h.userService.GetUserSessions(username)
	if err != nil {
		h.logger.Errorf("获取用户会话失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取用户会话失败")
		return
	}

	utils.ListSuccessResponse(c, sessions, len(sessions))
}

// KillUserSession 结束用户会话
func (h *UserHandler) KillUserSession(c *gin.Context) {
	sessionID := c.Param("sessionId")
	if sessionID == "" {
		utils.BadRequestResponse(c, "会话ID不能为空")
		return
	}

	err := h.userService.KillUserSession(sessionID)
	if err != nil {
		h.logger.Errorf("结束用户会话失败: %v", err)
		utils.OperationFailedResponse(c, "结束用户会话", err)
		return
	}

	utils.OperationSuccessResponse(c, "结束用户会话")
}

// LockUserAccount 锁定用户账户
func (h *UserHandler) LockUserAccount(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	err := h.userService.LockUserAccount(username)
	if err != nil {
		h.logger.Errorf("锁定用户账户失败: %v", err)
		utils.OperationFailedResponse(c, "锁定用户账户", err)
		return
	}

	utils.OperationSuccessResponse(c, "锁定用户账户")
}

// UnlockUserAccount 解锁用户账户
func (h *UserHandler) UnlockUserAccount(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	err := h.userService.UnlockUserAccount(username)
	if err != nil {
		h.logger.Errorf("解锁用户账户失败: %v", err)
		utils.OperationFailedResponse(c, "解锁用户账户", err)
		return
	}

	utils.OperationSuccessResponse(c, "解锁用户账户")
}

// ChangeUserPassword 修改用户密码
func (h *UserHandler) ChangeUserPassword(c *gin.Context) {
	var request models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	err := h.userService.ChangeUserPassword(request.Username, request.NewPassword)
	if err != nil {
		h.logger.Errorf("修改用户密码失败: %v", err)
		utils.OperationFailedResponse(c, "修改用户密码", err)
		return
	}

	utils.OperationSuccessResponse(c, "修改用户密码")
}

// GetUserGroups 获取用户组
func (h *UserHandler) GetUserGroups(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	groups, err := h.userService.GetUserGroups(username)
	if err != nil {
		h.logger.Errorf("获取用户组失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取用户组失败")
		return
	}

	utils.ListSuccessResponse(c, groups, len(groups))
}

// GetUserLoginHistory 获取用户登录历史
func (h *UserHandler) GetUserLoginHistory(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	// 获取limit参数，默认为100
	limitStr := c.DefaultQuery("limit", "100")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 100 // 默认值
	}

	history, err := h.userService.GetUserLoginHistory(username, limit)
	if err != nil {
		h.logger.Errorf("获取用户登录历史失败: %v", err)
		utils.InternalServerErrorResponse(c, "获取用户登录历史失败")
		return
	}

	utils.ListSuccessResponse(c, history, len(history))
}

// ClearAllCache 清除所有缓存
func (h *UserHandler) ClearAllCache(c *gin.Context) {
	err := h.userService.ClearAllCache()
	if err != nil {
		h.logger.Errorf("清除所有缓存失败: %v", err)
		utils.InternalServerErrorResponse(c, "清除所有缓存失败")
		return
	}

	utils.OperationSuccessResponse(c, "所有缓存已清除")
}

// ClearUserCache 清除指定用户的缓存
func (h *UserHandler) ClearUserCache(c *gin.Context) {
	username := c.Param("username")
	if username == "" {
		utils.BadRequestResponse(c, "用户名不能为空")
		return
	}

	err := h.userService.ClearUserCache(username)
	if err != nil {
		h.logger.Errorf("清除用户缓存失败: %v", err)
		utils.InternalServerErrorResponse(c, "清除用户缓存失败")
		return
	}

	utils.OperationSuccessResponse(c, "用户缓存已清除")
}

// GetCacheStatus 获取缓存状态
func (h *UserHandler) GetCacheStatus(c *gin.Context) {
	status := h.userService.GetCacheStatus()
	utils.InfoSuccessResponse(c, status)
}
