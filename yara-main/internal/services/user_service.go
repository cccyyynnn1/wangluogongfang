package services

import (
	"fmt"
	"sync"
	"time"

	"yara-security-service/internal/models"
	"yara-security-service/pkg/user"

	"github.com/sirupsen/logrus"
)

// UserService 用户服务
type UserService struct {
	manager *user.Manager
	logger  *logrus.Logger
	mu      sync.RWMutex

	// 用户会话管理
	userSessions map[string]*models.UserSession
	sessionTTL   time.Duration
}

// NewUserService 创建用户服务
func NewUserService(manager *user.Manager, logger *logrus.Logger) *UserService {
	return &UserService{
		manager:      manager,
		logger:       logger,
		userSessions: make(map[string]*models.UserSession),
		sessionTTL:   30 * time.Minute,
	}
}

// GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser() (*models.UserInfo, error) {
	userInfo, err := s.manager.GetCurrentUser()
	if err != nil {
		return nil, fmt.Errorf("获取当前用户信息失败: %w", err)
	}

	s.logger.Debug("获取当前用户信息成功")
	return userInfo, nil
}

// GetUserByID 根据用户ID获取用户信息
func (s *UserService) GetUserByID(uid string) (*models.UserInfo, error) {
	userInfo, err := s.manager.GetUserByID(uid)
	if err != nil {
		return nil, fmt.Errorf("根据用户ID获取用户信息失败: %w", err)
	}

	s.logger.Debugf("根据用户ID获取用户信息成功: %s", uid)
	return userInfo, nil
}

// GetUserByName 根据用户名获取用户信息
func (s *UserService) GetUserByName(username string) (*models.UserInfo, error) {
	userInfo, err := s.manager.GetUserByName(username)
	if err != nil {
		return nil, fmt.Errorf("根据用户名获取用户信息失败: %w", err)
	}

	s.logger.Debugf("根据用户名获取用户信息成功: %s", username)
	return userInfo, nil
}

// GetAllUsers 获取所有用户列表
func (s *UserService) GetAllUsers() ([]*models.UserInfo, error) {
	users, err := s.manager.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("获取所有用户列表失败: %w", err)
	}

	s.logger.Debugf("获取所有用户列表成功: %d 个用户", len(users))
	return users, nil
}

// CheckUserPermissions 检查用户权限
func (s *UserService) CheckUserPermissions(username string, permissions []string) (*models.PermissionResult, error) {
	userPerms, err := s.manager.CheckUserPermissions(username, permissions)
	if err != nil {
		return nil, fmt.Errorf("检查用户权限失败: %w", err)
	}

	// 转换为PermissionResult
	result := &models.PermissionResult{
		Username:    userPerms.Username,
		Permissions: userPerms.Permissions,
		IsAdmin:     userPerms.IsAdmin,
		CheckTime:   userPerms.CheckTime,
	}

	s.logger.Debugf("检查用户权限成功: %s", username)
	return result, nil
}

// ValidatePassword 验证用户密码
func (s *UserService) ValidatePassword(username, password string) (*models.PasswordValidationResponse, error) {
	// 验证密码
	validationResult, err := s.manager.ValidatePassword(username, password)
	if err != nil {
		return nil, fmt.Errorf("密码验证失败: %w", err)
	}

	s.logger.Debugf("验证用户密码: %s, 结果: %v", username, validationResult.IsValid)
	return validationResult, nil
}

// GetPasswordPolicy 获取密码策略
func (s *UserService) GetPasswordPolicy() (*models.PasswordPolicy, error) {
	policy := s.manager.GetPasswordPolicy()
	s.logger.Debug("获取密码策略成功")
	return policy, nil
}

// CheckAccountStatus 检查账户状态
func (s *UserService) CheckAccountStatus(username string) (*models.AccountStatus, error) {
	status, err := s.manager.CheckAccountStatus(username)
	if err != nil {
		return nil, fmt.Errorf("检查账户状态失败: %w", err)
	}

	s.logger.Debugf("检查账户状态成功: %s", username)
	return status, nil
}

// GetUserSessions 获取用户会话
func (s *UserService) GetUserSessions(username string) ([]*models.UserSession, error) {
	sessions, err := s.manager.GetUserSessions(username)
	if err != nil {
		return nil, fmt.Errorf("获取用户会话失败: %w", err)
	}

	s.logger.Debugf("获取用户会话成功: %s, 会话数: %d", username, len(sessions))
	return sessions, nil
}

// KillUserSession 结束用户会话
func (s *UserService) KillUserSession(sessionID string) error {
	err := s.manager.KillUserSession(sessionID)
	if err != nil {
		return fmt.Errorf("结束用户会话失败: %w", err)
	}

	s.logger.Infof("结束用户会话成功: %s", sessionID)
	return nil
}

// LockUserAccount 锁定用户账户
func (s *UserService) LockUserAccount(username string) error {
	err := s.manager.LockUserAccount(username)
	if err != nil {
		return fmt.Errorf("锁定用户账户失败: %w", err)
	}

	s.logger.Infof("锁定用户账户成功: %s", username)
	return nil
}

// UnlockUserAccount 解锁用户账户
func (s *UserService) UnlockUserAccount(username string) error {
	err := s.manager.UnlockUserAccount(username)
	if err != nil {
		return fmt.Errorf("解锁用户账户失败: %w", err)
	}

	s.logger.Infof("解锁用户账户成功: %s", username)
	return nil
}

// ChangeUserPassword 修改用户密码
func (s *UserService) ChangeUserPassword(username, newPassword string) error {
	err := s.manager.ChangeUserPassword(username, newPassword)
	if err != nil {
		return fmt.Errorf("修改用户密码失败: %w", err)
	}

	s.logger.Infof("修改用户密码成功: %s", username)
	return nil
}

// GetUserGroups 获取用户组信息
func (s *UserService) GetUserGroups(username string) ([]string, error) {
	groups, err := s.manager.GetUserGroups(username)
	if err != nil {
		return nil, fmt.Errorf("获取用户组信息失败: %w", err)
	}

	s.logger.Debugf("获取用户组信息: %s, 组数: %d", username, len(groups))
	return groups, nil
}

// GetUserLoginHistory 获取用户登录历史
func (s *UserService) GetUserLoginHistory(username string, limit int) ([]*models.LoginRecord, error) {
	history, err := s.manager.GetUserLoginHistory(username, limit)
	if err != nil {
		return nil, fmt.Errorf("获取用户登录历史失败: %w", err)
	}

	s.logger.Debugf("获取用户登录历史成功: %s, 记录数: %d", username, len(history))
	return history, nil
}

// CreateUserSession 创建用户会话
func (s *UserService) CreateUserSession(username string) (*models.UserSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	session := &models.UserSession{
		SessionID:  generateSessionID(),
		Username:   username,
		LoginTime:  time.Now(),
		LastActive: time.Now(),
		IPAddress:  "", // 需要从请求中获取
		UserAgent:  "", // 需要从请求中获取
		IsActive:   true,
	}

	s.userSessions[session.SessionID] = session

	s.logger.Infof("创建用户会话成功: %s, 会话ID: %s", username, session.SessionID)
	return session, nil
}

// UpdateSessionActivity 更新会话活动
func (s *UserService) UpdateSessionActivity(sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session, exists := s.userSessions[sessionID]
	if !exists {
		return fmt.Errorf("会话不存在: %s", sessionID)
	}

	session.LastActive = time.Now()
	s.logger.Debugf("更新会话活动: %s", sessionID)
	return nil
}

// ValidateSession 验证会话
func (s *UserService) ValidateSession(sessionID string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, exists := s.userSessions[sessionID]
	if !exists {
		return false, nil
	}

	// 检查会话是否过期
	if time.Since(session.LastActive) > s.sessionTTL {
		delete(s.userSessions, sessionID)
		return false, nil
	}

	return session.IsActive, nil
}

// GetActiveSessions 获取活跃会话
func (s *UserService) GetActiveSessions() []*models.UserSession {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var activeSessions []*models.UserSession
	now := time.Now()

	for _, session := range s.userSessions {
		if session.IsActive && now.Sub(session.LastActive) <= s.sessionTTL {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions
}

// CleanupExpiredSessions 清理过期会话
func (s *UserService) CleanupExpiredSessions() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	cleanedCount := 0
	now := time.Now()

	for sessionID, session := range s.userSessions {
		if now.Sub(session.LastActive) > s.sessionTTL {
			delete(s.userSessions, sessionID)
			cleanedCount++
		}
	}

	if cleanedCount > 0 {
		s.logger.Infof("清理过期会话: %d 个", cleanedCount)
	}

	return cleanedCount
}

// GetUserStatistics 获取用户统计信息
func (s *UserService) GetUserStatistics() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["active_sessions"] = len(s.userSessions)
	stats["session_ttl"] = s.sessionTTL.String()
	stats["cache_time"] = time.Now().Format(time.RFC3339)

	return stats
}

// ClearAllCache 清除所有缓存
func (s *UserService) ClearAllCache() error {
	s.manager.ClearAllCache()
	s.logger.Info("所有缓存已清除")
	return nil
}

// ClearUserCache 清除指定用户的缓存
func (s *UserService) ClearUserCache(username string) error {
	s.manager.ClearUserAdminCache(username)
	s.manager.ClearUserGroupsCache(username)
	s.logger.Infof("用户 %s 的缓存已清除", username)
	return nil
}

// GetCacheStatus 获取缓存状态
func (s *UserService) GetCacheStatus() map[string]interface{} {
	// 这里可以添加更详细的缓存状态信息
	// 目前返回基本信息
	return map[string]interface{}{
		"cache_cleared": true,
		"timestamp":     time.Now().Format(time.RFC3339),
	}
}

// generateSessionID 生成会话ID
func generateSessionID() string {
	// 这里可以使用UUID或其他方式生成会话ID
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}
