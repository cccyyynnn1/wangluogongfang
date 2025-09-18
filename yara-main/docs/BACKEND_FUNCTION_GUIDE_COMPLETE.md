# Yara 安全服务后端功能函数调用说明文档

## 📋 目录索引

### 🛡️ 安全扫描模块

- [type QuarantineFileRequest](#type-quarantinefilerequest)
- [type RestoreFileRequest](#type-restorefilerequest)
- [type SecurityStatus](#type-securitystatus)
- [type RulesInfo](#type-rulesinfo)
- [type RuleInfo](#type-ruleinfo)
- [type SecurityHandler](#type-securityhandler)
- [type SecurityService](#type-securityservice)
- [type Scanner](#type-scanner)
  - [func NewSecurityHandler(securityService *services.SecurityService, logger *logrus.Logger) \*SecurityHandler](#func-newsecurityhandler)
  - [func (h *SecurityHandler) GetSecurityStatus(c *gin.Context)](#func-h-securityhandler-getsecuritystatus)
  - [func (h *SecurityHandler) GetRulesInfo(c *gin.Context)](#func-h-securityhandler-getrulesinfo)
  - [func (h *SecurityHandler) ReloadRules(c *gin.Context)](#func-h-securityhandler-reloadrules)
  - [func (h *SecurityHandler) ClearCache(c *gin.Context)](#func-h-securityhandler-clearcache)
  - [func (h *SecurityHandler) GetCacheStats(c *gin.Context)](#func-h-securityhandler-getcachestats)
  - [func (h *SecurityHandler) QuarantineFile(c *gin.Context)](#func-h-securityhandler-quarantinefile)
  - [func (h *SecurityHandler) RestoreFile(c *gin.Context)](#func-h-securityhandler-restorefile)
  - [func (h *SecurityHandler) GetQuarantineList(c *gin.Context)](#func-h-securityhandler-getquarantinelist)
  - [func (h *SecurityHandler) GetScanHistory(c *gin.Context)](#func-h-securityhandler-getscanhistory)

### 📁 文件管理模块

- [type DeepFileAnalysis](#type-deepfileanalysis)
- [type BasicFileInfo](#type-basicfileinfo)
- [type ContentAnalysis](#type-contentanalysis)
- [type TextAnalysis](#type-textanalysis)
- [type BinaryAnalysis](#type-binaryanalysis)
- [type StructureAnalysis](#type-structureanalysis)
- [type PEAnalysis](#type-peanalysis)
- [type ThreatAnalysis](#type-threatanalysis)
- [type FileScanResult](#type-filescanresult)
- [type ArchiveAnalysis](#type-archiveanalysis)
- [type DocumentAnalysis](#type-documentanalysis)
- [type BehaviorAnalysis](#type-behavioranalysis)
- [type SignatureAnalysis](#type-signatureanalysis)
- [type EntropyAnalysis](#type-entropyanalysis)
- [type EncodingAnalysis](#type-encodinganalysis)
- [type ScanFileRequest](#type-scanfilerequest)
- [type ScanResult](#type-scanresult)
- [type ThreatInfo](#type-threatinfo)
- [type FileInfo](#type-fileinfo)
- [type BufferScanRequest](#type-bufferscanrequest)
- [type FileListResponse](#type-filelistresponse)
- [type DirectoryScanRequest](#type-directoryscanrequest)
- [type FileHandler](#type-filehandler)
- [type FileManager](#type-filemanager)
  - [func NewFileHandler(securityService *services.SecurityService, scanner *security.Scanner, logger *logrus.Logger) *FileHandler](#func-newfilehandler)
  - [func (h *FileHandler) ScanFile(c *gin.Context)](#func-h-filehandler-scanfile)
  - [func (h *FileHandler) GetFileInfo(c *gin.Context)](#func-h-filehandler-getfileinfo)
  - [func (h \*FileHandler) calculateFileHashes(filePath string) (map[string]string, error)](#func-h-filehandler-calculatefilehashes)
  - [func (h *FileHandler) ScanDirectory(c *gin.Context)](#func-h-filehandler-scandirectory)
  - [func (h *FileHandler) ScanBuffer(c *gin.Context)](#func-h-filehandler-scanbuffer)
  - [func (h *FileHandler) GetFileList(c *gin.Context)](#func-h-filehandler-getfilelist)
  - [func (h *FileHandler) getFileList(directory string, recursive bool) ([]*models.FileInfo, error)](#func-h-filehandler-getfilelist)
  - [func (m \*FileManager) CopyFile(source, dest string) error](#func-m-filemanager-copyfile)
  - [func (m \*FileManager) MoveFile(source, dest string) error](#func-m-filemanager-movefile)
  - [func (m \*FileManager) DeleteFile(filePath string) error](#func-m-filemanager-deletefile)
  - [func (m \*FileManager) GetFileHash(filePath, algorithm string) (string, error)](#func-m-filemanager-getfilehash)
  - [func (m \*FileManager) GetFileHashes(filePath string) (map[string]string, error)](#func-m-filemanager-getfilehashes)
  - [func (m \*FileManager) VerifyFileHash(filePath, algorithm, expectedHash string) (bool, error)](#func-m-filemanager-verifyfilehash)

### 🔍 进程管理模块

- [type ProcessEvent](#type-processevent)
- [type DeepProcessAnalysis](#type-deepprocessanalysis)
- [type BasicProcessInfo](#type-basicprocessinfo)
- [type ProcessInfo](#type-processinfo)
- [type ProcessThreatAnalysis](#type-processthreatanalysis)
- [type AddrInfo](#type-addrinfo)
- [type ProcessMemoryInfo](#type-processmemoryinfo)
- [type ProcessModuleInfo](#type-processmoduleinfo)
- [type ProcessListResponse](#type-processlistresponse)
- [type ProcessDetailResponse](#type-processdetailresponse)
- [type StartProcessRequest](#type-startprocessrequest)
- [type StartProcessResponse](#type-startprocessresponse)
- [type ProcessStatusResponse](#type-processstatusresponse)
- [type ProcessModulesResponse](#type-processmodulesresponse)
- [type ProcessConnectionInfo](#type-processconnectioninfo)
- [type ProcessConnectionsResponse](#type-processconnectionsresponse)
- [type ProcessMemoryResponse](#type-processmemoryresponse)
- [type ProcessRunningStatus](#type-processrunningstatus)
- [type ProcessRunningResponse](#type-processrunningresponse)
- [type ProcessChildrenResponse](#type-processchildrenresponse)
- [type SystemModuleInfo](#type-systemmoduleinfo)
- [type SystemModulesResponse](#type-systemmodulesresponse)
- [type SystemModuleDetailResponse](#type-systemmoduledetailresponse)
- [type ModuleStatusResponse](#type-modulestatusresponse)
- [type MonitoringStatusResponse](#type-monitoringstatusresponse)
- [type MonitoredProcessInfo](#type-monitoredprocessinfo)
- [type MonitoredProcessListResponse](#type-monitoredprocesslistresponse)
- [type ProcessStatistics](#type-processstatistics)
- [type TopProcessInfo](#type-topprocessinfo)
- [type ProcessHandler](#type-processhandler)
  - [func NewProcessHandler(processService *services.ProcessService, logger *logrus.Logger) \*ProcessHandler](#func-newprocesshandler)
  - [func (h *ProcessHandler) GetProcesses(c *gin.Context)](#func-h-processhandler-getprocesses)
  - [func (h *ProcessHandler) GetProcessByPID(c *gin.Context)](#func-h-processhandler-getprocessbypid)
  - [func (h *ProcessHandler) StartProcess(c *gin.Context)](#func-h-processhandler-startprocess)
  - [func (h *ProcessHandler) KillProcess(c *gin.Context)](#func-h-processhandler-killprocess)
  - [func (h *ProcessHandler) SuspendProcess(c *gin.Context)](#func-h-processhandler-suspendprocess)
  - [func (h *ProcessHandler) ResumeProcess(c *gin.Context)](#func-h-processhandler-resumeprocess)
  - [func (h *ProcessHandler) GetProcessModules(c *gin.Context)](#func-h-processhandler-getprocessmodules)
  - [func (h *ProcessHandler) GetProcessConnections(c *gin.Context)](#func-h-processhandler-getprocessconnections)
  - [func (h *ProcessHandler) GetProcessMemoryInfo(c *gin.Context)](#func-h-processhandler-getprocessmemoryinfo)
  - [func (h *ProcessHandler) IsProcessRunning(c *gin.Context)](#func-h-processhandler-isprocessrunning)
  - [func (h *ProcessHandler) GetProcessChildren(c *gin.Context)](#func-h-processhandler-getprocesschildren)
  - [func (h *ProcessHandler) GetSystemModules(c *gin.Context)](#func-h-processhandler-getsystemmodules)
  - [func (h *ProcessHandler) GetModuleInfo(c *gin.Context)](#func-h-processhandler-getmoduleinfo)
  - [func (h *ProcessHandler) SuspendModule(c *gin.Context)](#func-h-processhandler-suspendmodule)
  - [func (h *ProcessHandler) ResumeModule(c *gin.Context)](#func-h-processhandler-resumemodule)
  - [func (h *ProcessHandler) KillModule(c *gin.Context)](#func-h-processhandler-killmodule)
  - [func EnableProcessMonitoring(c \*gin.Context)](#func-enableprocessmonitoring)
  - [func DisableProcessMonitoring(c \*gin.Context)](#func-disableprocessmonitoring)
  - [func GetMonitoredProcesses(c \*gin.Context)](#func-getmonitoredprocesses)
  - [func GetProcessStatistics(c \*gin.Context)](#func-getprocessstatistics)

### 🔧 注册表管理模块

- [type RegistryKey](#type-registrykey)
- [type CreateRegistryKeyRequest](#type-createregistrykeyrequest)
- [type SetRegistryValueRequest](#type-setregistryvaluerequest)
- [type SearchRegistryRequest](#type-searchregistryrequest)
- [type RegistryHandler](#type-registryhandler)
- [type RegistryKeyResponse](#type-registrykeyresponse)
- [type RegistryKeyInfo](#type-registrykeyinfo)
- [type RegistryKeysResponse](#type-registrykeysresponse)
- [type RegistrySubKey](#type-registrysubkey)
- [type RegistryValuesResponse](#type-registryvaluesresponse)
- [type RegistryValue](#type-registryvalue)
- [type CreateRegistryKeyResponse](#type-createregistrykeyresponse)
- [type DeleteRegistryKeyResponse](#type-deleteregistrykeyresponse)
- [type RegistryValueResponse](#type-registryvalueresponse)
- [type RegistryValueDetail](#type-registryvaluedetail)
- [type SetRegistryValueResponse](#type-setregistryvalueresponse)
- [type DeleteRegistryValueResponse](#type-deleteregistryvalueresponse)
- [type RegistrySearchResponse](#type-registrysearchresponse)
- [type RegistrySearchResult](#type-registrysearchresult)
  - [func NewRegistryHandler(registryService *services.RegistryService, logger *logrus.Logger) \*RegistryHandler](#func-newregistryhandler)
  - [func (h *RegistryHandler) GetRegistryKey(c *gin.Context)](#func-h-registryhandler-getregistrykey)
  - [func (h *RegistryHandler) CreateRegistryKey(c *gin.Context)](#func-h-registryhandler-createregistrykey)
  - [func (h *RegistryHandler) DeleteRegistryKey(c *gin.Context)](#func-h-registryhandler-deleteregistrykey)
  - [func (h *RegistryHandler) SetRegistryValue(c *gin.Context)](#func-h-registryhandler-setregistryvalue)
  - [func (h *RegistryHandler) GetRegistryValue(c *gin.Context)](#func-h-registryhandler-getregistryvalue)
  - [func (h *RegistryHandler) DeleteRegistryValue(c *gin.Context)](#func-h-registryhandler-deleteregistryvalue)
  - [func (h *RegistryHandler) ListRegistryKeys(c *gin.Context)](#func-h-registryhandler-listregistrykeys)
  - [func (h *RegistryHandler) ListRegistryValues(c *gin.Context)](#func-h-registryhandler-listregistryvalues)
  - [func (h *RegistryHandler) SearchRegistry(c *gin.Context)](#func-h-registryhandler-searchregistry)

### 🌐 网络管理模块

- [type NetworkConnection](#type-networkconnection)
- [type NetworkPacket](#type-networkpacket)
- [type NetworkInterface](#type-networkinterface)
- [type NetworkThreat](#type-networkthreat)
- [type NetworkHandler](#type-networkhandler)
- [type NetworkConnectionsResponse](#type-networkconnectionsresponse)
- [type TCPConnectionsResponse](#type-tcpconnectionsresponse)
- [type UDPConnectionsResponse](#type-udpconnectionsresponse)
- [type PortConnectionsResponse](#type-portconnectionsresponse)
- [type IPConnectionsResponse](#type-ipconnectionsresponse)
- [type ConnectionDetailResponse](#type-connectiondetailresponse)
- [type ConnectionDetail](#type-connectiondetail)
- [type CloseConnectionRequest](#type-closeconnectionrequest)
- [type CloseConnectionResponse](#type-closeconnectionresponse)
- [type PortsInUseResponse](#type-portsinuseresponse)
- [type PortInfo](#type-portinfo)
- [type PortStatusResponse](#type-portstatusresponse)
- [type PortStatus](#type-portstatus)
- [type ListeningPortsResponse](#type-listeningportsresponse)
- [type ListeningPort](#type-listeningport)
- [type EstablishedConnectionsResponse](#type-establishedconnectionsresponse)
- [type NetworkInterfacesResponse](#type-networkinterfacesresponse)
- [type NetworkStatsResponse](#type-networkstatsresponse)
- [type EnableMonitoringResponse](#type-enablemonitoringresponse)
- [type DisableMonitoringResponse](#type-disablemonitoringresponse)
- [type MonitoringConnectionsResponse](#type-monitoringconnectionsresponse)
- [type ConnectionHistoryResponse](#type-connectionhistoryresponse)
- [type ConnectionHistory](#type-connectionhistory)
  - [func NewNetworkHandler(networkService *services.NetworkService, logger *logrus.Logger) \*NetworkHandler](#func-newnetworkhandler)
  - [func (h *NetworkHandler) GetConnections(c *gin.Context)](#func-h-networkhandler-getconnections)
  - [func (h *NetworkHandler) GetTCPConnections(c *gin.Context)](#func-h-networkhandler-gettcpconnections)
  - [func (h *NetworkHandler) GetUDPConnections(c *gin.Context)](#func-h-networkhandler-getudpconnections)
  - [func (h *NetworkHandler) GetConnectionsByPID(c *gin.Context)](#func-h-networkhandler-getconnectionsbypid)
  - [func (h *NetworkHandler) CloseConnection(c *gin.Context)](#func-h-networkhandler-closeconnection)
  - [func (h *NetworkHandler) GetNetworkInterfaces(c *gin.Context)](#func-h-networkhandler-getnetworkinterfaces)
  - [func (h *NetworkHandler) GetNetworkStats(c *gin.Context)](#func-h-networkhandler-getnetworkstats)
  - [func (h *NetworkHandler) GetConnectionsByPort(c *gin.Context)](#func-h-networkhandler-getconnectionsbyport)
  - [func (h *NetworkHandler) GetConnectionsByIP(c *gin.Context)](#func-h-networkhandler-getconnectionsbyip)
  - [func (h *NetworkHandler) IsPortInUse(c *gin.Context)](#func-h-networkhandler-isportinuse)
  - [func (h *NetworkHandler) GetListeningPorts(c *gin.Context)](#func-h-networkhandler-getlisteningports)
  - [func (h *NetworkHandler) GetEstablishedConnections(c *gin.Context)](#func-h-networkhandler-getestablishedconnections)
  - [func EnableNetworkMonitoring(c \*gin.Context)](#func-enablenetworkmonitoring)
  - [func DisableNetworkMonitoring(c \*gin.Context)](#func-disablenetworkmonitoring)
  - [func GetMonitoredConnections(c \*gin.Context)](#func-getmonitoredconnections)
  - [func GetConnectionHistory(c \*gin.Context)](#func-getconnectionhistory)
  - [func (n *network.Manager) GetEstablishedConnections() ([]*models.NetworkConnection, error)](#func-n-networkmanager-getestablishedconnections)
  - [func (n *network.Manager) GetNetworkInterfaces() ([]*models.NetworkInterface, error)](#func-n-networkmanager-getnetworkinterfaces)
  - [func (n \*network.Service) EnableMonitoring() error](#func-n-networkservice-enablemonitoring)
  - [func (n \*network.Service) DisableMonitoring() error](#func-n-networkservice-disablemonitoring)
  - [func (n *network.Service) GetMonitoredConnections() ([]*models.MonitoredNetworkConnection, error)](#func-n-networkservice-getmonitoredconnections)
  - [func (n *network.Service) GetConnectionHistory(limit, offset int, startDate, endDate string) ([]*models.NetworkConnectionHistory, error)](#func-n-networkservice-getconnectionhistory)

### 👤 用户管理模块

- [type UserHandler](#type-userhandler)
- [type CurrentUserResponse](#type-currentuserresponse)
- [type CurrentUser](#type-currentuser)
- [type AllUsersResponse](#type-allusersresponse)
- [type UserInfo](#type-userinfo)
- [type UserByIdResponse](#type-userbyidresponse)
- [type UserDetail](#type-userdetail)
- [type UserByNameResponse](#type-userbynameresponse)
- [type UserGroupsResponse](#type-usergroupsresponse)
- [type UserGroup](#type-usergroup)
- [type CheckPermissionRequest](#type-checkpermissionrequest)
- [type CheckPermissionResponse](#type-checkpermissionresponse)
- [type ValidatePasswordRequest](#type-validatepasswordrequest)
- [type ValidatePasswordResponse](#type-validatepasswordresponse)
- [type PasswordPolicyResponse](#type-passwordpolicyresponse)
- [type PasswordPolicy](#type-passwordpolicy)
- [type UserStatusResponse](#type-userstatusresponse)
- [type UserStatus](#type-userstatus)
- [type UserSessionsResponse](#type-usersessionsresponse)
- [type UserSession](#type-usersession)
- [type EndSessionResponse](#type-endsessionresponse)
- [type LockUserResponse](#type-lockuserresponse)
- [type UnlockUserResponse](#type-unlockuserresponse)
- [type ChangePasswordRequest](#type-changepasswordrequest)
- [type ChangePasswordResponse](#type-changepasswordresponse)
- [type UserHistoryResponse](#type-userhistoryresponse)
- [type LoginHistory](#type-loginhistory)
  - [func NewUserHandler(manager *user.Manager, logger *logrus.Logger) \*UserHandler](#func-newuserhandler)
  - [func (h *UserHandler) GetCurrentUser(c *gin.Context)](#func-h-userhandler-getcurrentuser)
  - [func (h *UserHandler) GetUserByID(c *gin.Context)](#func-h-userhandler-getuserbyid)
  - [func (h *UserHandler) GetUserByName(c *gin.Context)](#func-h-userhandler-getuserbyname)
  - [func (h *UserHandler) GetAllUsers(c *gin.Context)](#func-h-userhandler-getallusers)
  - [func (h *UserHandler) CheckUserPermissions(c *gin.Context)](#func-h-userhandler-checkuserpermissions)
  - [func (h *UserHandler) ValidatePassword(c *gin.Context)](#func-h-userhandler-validatepassword)
  - [func (h *UserHandler) GetPasswordPolicy(c *gin.Context)](#func-h-userhandler-getpasswordpolicy)
  - [func (h *UserHandler) CheckAccountStatus(c *gin.Context)](#func-h-userhandler-checkaccountstatus)
  - [func (h *UserHandler) GetUserSessions(c *gin.Context)](#func-h-userhandler-getusersessions)
  - [func (h *UserHandler) KillUserSession(c *gin.Context)](#func-h-userhandler-killusersession)
  - [func (h *UserHandler) LockUserAccount(c *gin.Context)](#func-h-userhandler-lockuseraccount)
  - [func (h *UserHandler) UnlockUserAccount(c *gin.Context)](#func-h-userhandler-unlockuseraccount)
  - [func (h *UserHandler) ChangeUserPassword(c *gin.Context)](#func-h-userhandler-changeuserpassword)
  - [func (h *UserHandler) GetUserGroups(c *gin.Context)](#func-h-userhandler-getusergroups)
  - [func (h *UserHandler) GetUserLoginHistory(c *gin.Context)](#func-h-userhandler-getuserloginhistory)

### 🔧 中间件模块

- [type AuthConfig](#type-authconfig)
  - [func DefaultAuthConfig() \*AuthConfig](#func-defaultauthconfig)
  - [func Auth(config *AuthConfig, logger *logrus.Logger) gin.HandlerFunc](#func-auth)
- [type RateLimitConfig](#type-ratelimitconfig)
  - [func DefaultRateLimitConfig() \*RateLimitConfig](#func-defaultratelimitconfig)
  - [func RateLimit(config \*RateLimitConfig) gin.HandlerFunc](#func-ratelimit)
- [type MetricsCollector](#type-metricscollector)
  - [func NewMetricsCollector(logger *logrus.Logger) *MetricsCollector](#func-newmetricscollector)
  - [func (mc \*MetricsCollector) Collect() gin.HandlerFunc](#func-mc-metricscollector-collect)
  - [func (mc *MetricsCollector) GetMetrics() *Metrics](#func-mc-metricscollector-getmetrics)
  - [func (mc \*MetricsCollector) ResetMetrics()](#func-mc-metricscollector-resetmetrics)
  - [func (mc \*MetricsCollector) ReportMetrics()](#func-mc-metricscollector-reportmetrics)
- [type TimeoutConfig](#type-timeoutconfig)
  - [func DefaultTimeoutConfig() \*TimeoutConfig](#func-defaulttimeoutconfig)
  - [func Timeout(config *TimeoutConfig, logger *logrus.Logger) gin.HandlerFunc](#func-timeout)
- [type ErrorHandler](#type-errorhandler)
  - [func NewErrorHandler(logger *logrus.Logger) *ErrorHandler](#func-newerrorhandler)
  - [func (h \*ErrorHandler) Recovery() gin.HandlerFunc](#func-h-errorhandler-recovery)
  - [func (h \*ErrorHandler) ErrorResponse() gin.HandlerFunc](#func-h-errorhandler-errorresponse)
  - [func (h \*ErrorHandler) Timeout(timeout time.Duration) gin.HandlerFunc](#func-h-errorhandler-timeout)
  - [func (h \*ErrorHandler) RateLimit(requests int, window time.Duration) gin.HandlerFunc](#func-h-errorhandler-ratelimit)
  - [func (h \*ErrorHandler) Validation() gin.HandlerFunc](#func-h-errorhandler-validation)
  - [func (h \*ErrorHandler) Logging() gin.HandlerFunc](#func-h-errorhandler-logging)
  - [func (h \*ErrorHandler) Security() gin.HandlerFunc](#func-h-errorhandler-security)
  - [func (h \*ErrorHandler) HealthCheck() gin.HandlerFunc](#func-h-errorhandler-healthcheck)
  - [func (h *ErrorHandler) HandlePanic(c *gin.Context)](#func-h-errorhandler-handlepanic)
- [type LoggingConfig](#type-loggingconfig)
  - [func DefaultLoggingConfig() \*LoggingConfig](#func-defaultloggingconfig)
  - [func Logging(config *LoggingConfig, logger *logrus.Logger) gin.HandlerFunc](#func-logging)
  - [func Recovery(logger \*logrus.Logger) gin.HandlerFunc](#func-recovery)
  - [func RequestID() gin.HandlerFunc](#func-requestid)
- [type CORSConfig](#type-corsconfig)
  - [func DefaultCORSConfig() \*CORSConfig](#func-defaultcorsconfig)
  - [func CORS(config \*CORSConfig) gin.HandlerFunc](#func-cors)

### 🛠️ 工具模块

- [type CryptoUtils](#type-cryptoutils)
  - [func NewCryptoUtils() \*CryptoUtils](#func-newcryptoutils)
  - [func (c \*CryptoUtils) GenerateRandomBytes(length int) ([]byte, error)](#func-c-cryptoutils-generaterandombytes)
  - [func (c \*CryptoUtils) GenerateRandomString(length int) (string, error)](#func-c-cryptoutils-generaterandomstring)
  - [func (c \*CryptoUtils) HashSHA256(data []byte) string](#func-c-cryptoutils-hashsha256)
  - [func (c \*CryptoUtils) HashSHA512(data []byte) string](#func-c-cryptoutils-hashsha512)
  - [func (c \*CryptoUtils) HashWithSalt(data []byte, salt []byte) string](#func-c-cryptoutils-hashwithsalt)
  - [func (c \*CryptoUtils) GenerateHMAC(data []byte, key []byte) string](#func-c-cryptoutils-generatehmac)
  - [func (c \*CryptoUtils) VerifyHMAC(data []byte, key []byte, expectedHMAC string) bool](#func-c-cryptoutils-verifyhmac)
  - [func (c \*CryptoUtils) EncryptAES(data []byte, key []byte) ([]byte, error)](#func-c-cryptoutils-encryptaes)
  - [func (c \*CryptoUtils) DecryptAES(ciphertext []byte, key []byte) ([]byte, error)](#func-c-cryptoutils-decryptaes)
  - [func (c \*CryptoUtils) GenerateKey(keySize int) ([]byte, error)](#func-c-cryptoutils-generatekey)
  - [func (c \*CryptoUtils) HashFile(data []byte, algorithm string) (string, error)](#func-c-cryptoutils-hashfile)
  - [func (c \*CryptoUtils) Base64Encode(data []byte) string](#func-c-cryptoutils-base64encode)
  - [func (c \*CryptoUtils) Base64Decode(encoded string) ([]byte, error)](#func-c-cryptoutils-base64decode)
  - [func (c \*CryptoUtils) HexEncode(data []byte) string](#func-c-cryptoutils-hexencode)
  - [func (c \*CryptoUtils) HexDecode(encoded string) ([]byte, error)](#func-c-cryptoutils-hexdecode)
  - [func (c \*CryptoUtils) GeneratePasswordHash(password string, salt []byte) string](#func-c-cryptoutils-generatepasswordhash)
  - [func (c \*CryptoUtils) VerifyPassword(password string, salt []byte, expectedHash string) bool](#func-c-cryptoutils-verifypassword)
  - [func (c \*CryptoUtils) GenerateAPIKey(length int) (string, error)](#func-c-cryptoutils-generateapikey)
  - [func (c \*CryptoUtils) EncryptSensitiveData(data string, key []byte) (string, error)](#func-c-cryptoutils-encryptsensitivedata)
  - [func (c \*CryptoUtils) DecryptSensitiveData(encryptedData string, key []byte) (string, error)](#func-c-cryptoutils-decryptsensitivedata)
- [type Validator](#type-validator)
- [func NewValidator() \*Validator](#func-newvalidator)
- [func (v \*Validator) ValidateFilePath(filePath string) error](#func-v-validator-validatefilepath)
- [func (v \*Validator) ValidateDirectoryPath(dirPath string) error](#func-v-validator-validatedirectorypath)
- [func (v \*Validator) ValidatePID(pid int32) error](#func-v-validator-validatepid)
- [func (v \*Validator) ValidatePort(port int) error](#func-v-validator-validateport)
- [func (v \*Validator) ValidateIPAddress(ip string) error](#func-v-validator-validateipaddress)
- [func (v \*Validator) ValidateIPRange(ipRange string) error](#func-v-validator-validateiprange)
- [func (v \*Validator) ValidateEmail(email string) error](#func-v-validator-validateemail)
- [func (v \*Validator) ValidateURL(url string) error](#func-v-validator-validateurl)
- [func (v \*Validator) ValidateFileName(fileName string) error](#func-v-validator-validatefilename)
- [func (v \*Validator) ValidateFileSize(size int64, maxSize int64) error](#func-v-validator-validatefilesize)
- [func (v \*Validator) ValidateStringLength(str string, minLength, maxLength int) error](#func-v-validator-validatestringlength)
- [func (v \*Validator) ValidateInteger(value int, min, max int) error](#func-v-validator-validateinteger)
- [func (v \*Validator) ValidateFloat(value float64, min, max float64) error](#func-v-validator-validatefloat)
- [func (v \*Validator) ValidateRegex(pattern string) error](#func-v-validator-validateregex)
- [func (v \*Validator) ValidateFileExtension(fileName string, allowedExtensions []string) error](#func-v-validator-validatefileextension)
- [func (v \*Validator) ValidatePathDepth(path string, maxDepth int) error](#func-v-validator-validatepathdepth)
- [func (v \*Validator) ValidateRegistryPath(regPath string) error](#func-v-validator-validateregistrypath)
- [func (v \*Validator) ValidateProcessName(processName string) error](#func-v-validator-validateprocessname)
- [func (v \*Validator) ValidateScanDepth(depth int) error](#func-v-validator-validatescandepth)
- [func (v \*Validator) ValidateTimeout(timeout int) error](#func-v-validator-validatetimeout)
- [func (v \*Validator) ValidateAPIKey(apiKey string) error](#func-v-validator-validateapikey)
- [func (v \*Validator) ValidateHash(hash string, algorithm string) error](#func-v-validator-validatehash)
- [func (v \*Validator) ValidateNumericString(str string) error](#func-v-validator-validatenumericstring)
- [func (v \*Validator) ValidateHexString(str string) error](#func-v-validator-validatehexstring)

---

## 🛡️ 安全扫描模块

### type QuarantineFileRequest

```go
type QuarantineFileRequest struct {
    FilePath     string `json:"file_path" binding:"required"`
    ThreatType   string `json:"threat_type"`
    Description  string `json:"description"`
}
```

QuarantineFileRequest 表示隔离文件请求结构。

**字段说明:**

- `FilePath` (string, 必需): 要隔离的文件路径
- `ThreatType` (string): 威胁类型
- `Description` (string): 威胁描述

### type RestoreFileRequest

```go
type RestoreFileRequest struct {
    QuarantinePath string `json:"quarantine_path" binding:"required"`
    RestorePath    string `json:"restore_path" binding:"required"`
}
```

RestoreFileRequest 表示恢复文件请求结构。

**字段说明:**

- `QuarantinePath` (string, 必需): 隔离区文件路径
- `RestorePath` (string, 必需): 恢复目标路径

### type SecurityStatus

```go
type SecurityStatus struct {
    ProtectionEnabled    bool      `json:"protection_enabled"`
    RealTimeProtection   bool      `json:"real_time_protection"`
    ThreatCount          int       `json:"threat_count"`
    QuarantinedCount     int       `json:"quarantined_count"`
    LastScanTime         time.Time `json:"last_scan_time"`
    LastUpdateTime       time.Time `json:"last_update_time"`
    EngineVersion        string    `json:"engine_version"`
    DatabaseVersion      string    `json:"database_version"`
}
```

SecurityStatus 表示安全状态结构。

**字段说明:**

- `ProtectionEnabled` (bool): 防护是否启用
- `RealTimeProtection` (bool): 实时防护是否启用
- `ThreatCount` (int): 威胁数量
- `QuarantinedCount` (int): 隔离文件数量
- `LastScanTime` (time.Time): 最后扫描时间
- `LastUpdateTime` (time.Time): 最后更新时间
- `EngineVersion` (string): 引擎版本
- `DatabaseVersion` (string): 数据库版本

### type RulesInfo

```go
type RulesInfo struct {
    TotalRules    int                    `json:"total_rules"`
    ActiveRules   int                    `json:"active_rules"`
    RuleList      []RuleInfo             `json:"rule_list"`
    LastUpdated   time.Time              `json:"last_updated"`
    Categories    map[string]int         `json:"categories"`
    SeverityLevel map[string]int         `json:"severity_level"`
}
```

RulesInfo 表示规则信息结构。

**字段说明:**

- `TotalRules` (int): 总规则数
- `ActiveRules` (int): 活跃规则数
- `RuleList` ([]RuleInfo): 规则列表
- `LastUpdated` (time.Time): 最后更新时间
- `Categories` (map[string]int): 分类统计
- `SeverityLevel` (map[string]int): 严重级别统计

### type RuleInfo

```go
type RuleInfo struct {
    Name        string    `json:"name"`
    Category    string    `json:"category"`
    Severity    string    `json:"severity"`
    Description string    `json:"description"`
    Version     string    `json:"version"`
    Enabled     bool      `json:"enabled"`
    LastUsed    time.Time `json:"last_used"`
}
```

RuleInfo 表示规则信息结构。

**字段说明:**

- `Name` (string): 规则名称
- `Category` (string): 规则分类
- `Severity` (string): 严重级别
- `Description` (string): 规则描述
- `Version` (string): 规则版本
- `Enabled` (bool): 是否启用
- `LastUsed` (time.Time): 最后使用时间

### type SecurityHandler

```go
type SecurityHandler struct {
    securityService *services.SecurityService
    logger          *logrus.Logger
}
```

SecurityHandler 表示一个安全功能处理器，提供系统安全状态监控、规则管理和威胁处理功能。

### type SecurityService

```go
type SecurityService struct {
    scanner *security.Scanner
    logger  *logrus.Logger
    mu      sync.RWMutex

    // 优化的缓存机制
    scanCache     map[string]*models.ScanResult
    cacheTTL      time.Duration
    cacheStats    *CacheStats
    cacheMutex    sync.RWMutex
    maxCacheSize  int
    evictionQueue []string // LRU队列

    // 统计
    scanStats *ScanStatistics

    // 并发控制
    scanWorkers   int
    scanSemaphore chan struct{}
    scanTimeout   time.Duration
}
```

SecurityService 表示安全服务，提供文件扫描、缓存管理和统计功能。

**字段说明:**

- `scanner` (\*security.Scanner): 扫描器实例
- `logger` (\*logrus.Logger): 日志记录器
- `mu` (sync.RWMutex): 读写互斥锁
- `scanCache` (map[string]\*models.ScanResult): 扫描缓存
- `cacheTTL` (time.Duration): 缓存生存时间
- `cacheStats` (\*CacheStats): 缓存统计
- `cacheMutex` (sync.RWMutex): 缓存互斥锁
- `maxCacheSize` (int): 最大缓存大小
- `evictionQueue` ([]string): LRU 队列
- `scanStats` (\*ScanStatistics): 扫描统计
- `scanWorkers` (int): 扫描工作线程数
- `scanSemaphore` (chan struct{}): 扫描信号量
- `scanTimeout` (time.Duration): 扫描超时时间

### type Scanner

```go
type Scanner struct {
    yaraEngine *YaraEngine
    logger     *logrus.Logger
    rules      []*yara.Rule
    mu         sync.RWMutex
}
```

Scanner 表示扫描器，提供核心的 Yara 规则扫描功能。

**字段说明:**

- `yaraEngine` (\*YaraEngine): Yara 引擎
- `logger` (\*logrus.Logger): 日志记录器
- `rules` ([]\*yara.Rule): Yara 规则列表
- `mu` (sync.RWMutex): 读写互斥锁

### func NewSecurityHandler

```go
func NewSecurityHandler(securityService *services.SecurityService, logger *logrus.Logger) *SecurityHandler
```

NewSecurityHandler 使用指定的安全服务和日志记录器创建 SecurityHandler 的新实例。

### func (h \*SecurityHandler) GetSecurityStatus

```go
func (h *SecurityHandler) GetSecurityStatus(c *gin.Context)
```

GetSecurityStatus 返回系统安全状态概览，包括防护状态、实时保护状态、威胁数量等。

**功能说明:**

1. 检查系统防护状态
2. 获取实时保护状态
3. 统计威胁数量和隔离文件数量
4. 返回安全状态概览

### func (h \*SecurityHandler) GetRulesInfo

```go
func (h *SecurityHandler) GetRulesInfo(c *gin.Context)
```

GetRulesInfo 获取已加载的安全规则信息，包括规则数量、规则名称、严重级别等。

**功能说明:**

1. 获取已加载的安全规则总数
2. 获取规则列表和详细信息
3. 获取最后更新时间
4. 返回规则信息

### func (h \*SecurityHandler) ReloadRules

```go
func (h *SecurityHandler) ReloadRules(c *gin.Context)
```

ReloadRules 重新加载安全规则文件，更新检测规则库，确保使用最新的威胁检测规则。

**功能说明:**

1. 重新加载安全规则文件
2. 更新检测规则库
3. 验证规则加载状态
4. 返回重新加载结果

### func (h \*SecurityHandler) ClearCache

```go
func (h *SecurityHandler) ClearCache(c *gin.Context)
```

ClearCache 清空扫描结果缓存，释放内存资源，强制下次扫描重新计算。

**功能说明:**

1. 清空扫描结果缓存
2. 释放内存资源
3. 记录清理时间
4. 返回清理结果

### func (h \*SecurityHandler) GetCacheStats

```go
func (h *SecurityHandler) GetCacheStats(c *gin.Context)
```

GetCacheStats 获取缓存使用统计，包括缓存命中率、缓存大小、最后清理时间等。

**功能说明:**

1. 获取缓存大小和命中率
2. 统计总请求数和缓存命中数
3. 获取最后清理时间
4. 返回缓存统计信息

### func (h \*SecurityHandler) QuarantineFile

```go
func (h *SecurityHandler) QuarantineFile(c *gin.Context)
```

QuarantineFile 将检测到的威胁文件移动到隔离区，防止文件继续造成危害。

**功能说明:**

1. 从请求体中解析文件路径和威胁类型
2. 验证文件路径的有效性
3. 将文件移动到隔离区
4. 返回隔离结果

### func (h \*SecurityHandler) RestoreFile

```go
func (h *SecurityHandler) RestoreFile(c *gin.Context)
```

RestoreFile 将隔离区中的文件恢复到原始位置，适用于误报文件的恢复。

**功能说明:**

1. 从请求体中解析隔离路径和恢复路径
2. 验证文件路径的有效性
3. 将文件从隔离区恢复到原始位置
4. 返回恢复结果

### func (h \*SecurityHandler) GetQuarantineList

```go
func (h *SecurityHandler) GetQuarantineList(c *gin.Context)
```

GetQuarantineList 获取所有被隔离的文件信息，包括文件路径、隔离原因、隔离时间等。

**功能说明:**

1. 获取隔离区中的所有文件信息
2. 统计隔离文件总数
3. 获取每个文件的详细信息
4. 返回隔离文件列表

### func (h \*SecurityHandler) GetScanHistory

```go
func (h *SecurityHandler) GetScanHistory(c *gin.Context)
```

GetScanHistory 获取系统扫描历史记录，包括扫描时间、扫描文件数、检测到的威胁等。

**功能说明:**

1. 获取系统扫描历史记录
2. 统计扫描文件数和检测到的威胁数
3. 获取扫描时间和结果
4. 返回扫描历史记录

---

## 📁 文件管理模块

### type DeepFileAnalysis

```go
type DeepFileAnalysis struct {
    FilePath          string             `json:"file_path"`
    Timestamp         time.Time          `json:"timestamp"`
    IsSuspicious      bool               `json:"is_suspicious"`
    ThreatLevel       string             `json:"threat_level"`
    RiskScore         int                `json:"risk_score"`
    RiskFactors       []string           `json:"risk_factors"`
    BasicInfo         *BasicFileInfo     `json:"basic_info"`
    ContentAnalysis   *ContentAnalysis   `json:"content_analysis"`
    StructureAnalysis *StructureAnalysis `json:"structure_analysis"`
    ThreatAnalysis    *ThreatAnalysis    `json:"threat_analysis"`
}
```

DeepFileAnalysis 表示深度文件分析结果结构。

**字段说明:**

- `FilePath` (string): 文件路径
- `Timestamp` (time.Time): 分析时间戳
- `IsSuspicious` (bool): 是否可疑
- `ThreatLevel` (string): 威胁等级
- `RiskScore` (int): 风险评分
- `RiskFactors` ([]string): 风险因素列表
- `BasicInfo` (\*BasicFileInfo): 基础文件信息
- `ContentAnalysis` (\*ContentAnalysis): 内容分析
- `StructureAnalysis` (\*StructureAnalysis): 结构分析
- `ThreatAnalysis` (\*ThreatAnalysis): 威胁分析

### type BasicFileInfo

```go
type BasicFileInfo struct {
    Name         string    `json:"name"`
    Size         int64     `json:"size"`
    SizeCategory string    `json:"size_category"`
    Extension    string    `json:"extension"`
    FileType     string    `json:"file_type"`
    ModTime      time.Time `json:"mod_time"`
    CreateTime   time.Time `json:"create_time"`
    AccessTime   time.Time `json:"access_time"`
    Permissions  string    `json:"permissions"`
    MD5Hash      string    `json:"md5_hash"`
    SHA256Hash   string    `json:"sha256_hash"`
    IsHidden     bool      `json:"is_hidden"`
    IsSystem     bool      `json:"is_system"`
}
```

BasicFileInfo 表示基础文件信息结构。

**字段说明:**

- `Name` (string): 文件名
- `Size` (int64): 文件大小
- `SizeCategory` (string): 文件大小类别
- `Extension` (string): 文件扩展名
- `FileType` (string): 文件类型
- `ModTime` (time.Time): 修改时间
- `CreateTime` (time.Time): 创建时间
- `AccessTime` (time.Time): 访问时间
- `Permissions` (string): 权限信息
- `MD5Hash` (string): MD5 哈希值
- `SHA256Hash` (string): SHA256 哈希值
- `IsHidden` (bool): 是否隐藏文件
- `IsSystem` (bool): 是否系统文件

### type ContentAnalysis

```go
type ContentAnalysis struct {
    Timestamp          time.Time         `json:"timestamp"`
    TextAnalysis       *TextAnalysis     `json:"text_analysis"`
    BinaryAnalysis     *BinaryAnalysis   `json:"binary_analysis"`
    EncodingAnalysis   *EncodingAnalysis `json:"encoding_analysis"`
    SuspiciousPatterns []string          `json:"suspicious_patterns"`
}
```

ContentAnalysis 表示内容分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `TextAnalysis` (\*TextAnalysis): 文本分析
- `BinaryAnalysis` (\*BinaryAnalysis): 二进制分析
- `EncodingAnalysis` (\*EncodingAnalysis): 编码分析
- `SuspiciousPatterns` ([]string): 可疑模式列表

### type TextAnalysis

```go
type TextAnalysis struct {
    Timestamp         time.Time `json:"timestamp"`
    Encoding          string    `json:"encoding"`
    LineCount         int       `json:"line_count"`
    WordCount         int       `json:"word_count"`
    CharacterCount    int       `json:"character_count"`
    Language          string    `json:"language"`
    SpecialCharacters []string  `json:"special_characters"`
}
```

TextAnalysis 表示文本分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `Encoding` (string): 编码格式
- `LineCount` (int): 行数
- `WordCount` (int): 单词数
- `CharacterCount` (int): 字符数
- `Language` (string): 语言
- `SpecialCharacters` ([]string): 特殊字符列表

### type BinaryAnalysis

```go
type BinaryAnalysis struct {
    Timestamp   time.Time `json:"timestamp"`
    Entropy     float64   `json:"entropy"`
    HighEntropy bool      `json:"high_entropy"`
    FileHeader  string    `json:"file_header"`
    Strings     []string  `json:"strings"`
    HexPatterns []string  `json:"hex_patterns"`
}
```

BinaryAnalysis 表示二进制分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `Entropy` (float64): 熵值
- `HighEntropy` (bool): 是否高熵值
- `FileHeader` (string): 文件头
- `Strings` ([]string): 字符串列表
- `HexPatterns` ([]string): 十六进制模式列表

### type StructureAnalysis

```go
type StructureAnalysis struct {
    Timestamp        time.Time         `json:"timestamp"`
    PEAnalysis       *PEAnalysis       `json:"pe_analysis"`
    ArchiveAnalysis  *ArchiveAnalysis  `json:"archive_analysis"`
    DocumentAnalysis *DocumentAnalysis `json:"document_analysis"`
}
```

StructureAnalysis 表示结构分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `PEAnalysis` (\*PEAnalysis): PE 文件分析
- `ArchiveAnalysis` (\*ArchiveAnalysis): 压缩文件分析
- `DocumentAnalysis` (\*DocumentAnalysis): 文档分析

### type PEAnalysis

```go
type PEAnalysis struct {
    Timestamp       time.Time `json:"timestamp"`
    IsValidPE       bool      `json:"is_valid_pe"`
    EntryPoint      string    `json:"entry_point"`
    ImageBase       string    `json:"image_base"`
    Subsystem       string    `json:"subsystem"`
    Characteristics string    `json:"characteristics"`
    Imports         []string  `json:"imports"`
    Exports         []string  `json:"exports"`
    Sections        []string  `json:"sections"`
}
```

PEAnalysis 表示 PE 文件分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `IsValidPE` (bool): 是否为有效 PE 文件
- `EntryPoint` (string): 入口点
- `ImageBase` (string): 镜像基址
- `Subsystem` (string): 子系统
- `Characteristics` (string): 特征
- `Imports` ([]string): 导入列表
- `Exports` ([]string): 导出列表
- `Sections` ([]string): 节列表

### type ThreatAnalysis

```go
type ThreatAnalysis struct {
    Timestamp         time.Time          `json:"timestamp"`
    ScanResult        *ScanResult        `json:"scan_result"`
    IsInfected        bool               `json:"is_infected"`
    Threats           []ThreatInfo       `json:"threats"`
    BehaviorAnalysis  *BehaviorAnalysis  `json:"behavior_analysis"`
    SignatureAnalysis *SignatureAnalysis `json:"signature_analysis"`
    EntropyAnalysis   *EntropyAnalysis   `json:"entropy_analysis"`
}
```

ThreatAnalysis 表示威胁分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `ScanResult` (\*ScanResult): 扫描结果
- `IsInfected` (bool): 是否被感染
- `Threats` ([]ThreatInfo): 威胁列表
- `BehaviorAnalysis` (\*BehaviorAnalysis): 行为分析
- `SignatureAnalysis` (\*SignatureAnalysis): 签名分析
- `EntropyAnalysis` (\*EntropyAnalysis): 熵值分析

### type FileScanResult

```go
type FileScanResult struct {
    FilePath     string            `json:"file_path"`
    FileType     string            `json:"file_type"`
    IsSuspicious bool              `json:"is_suspicious"`
    Threats      []ThreatInfo      `json:"threats"`
    ScanTime     time.Time         `json:"scan_time"`
    ScanDuration time.Duration     `json:"scan_duration"`
    ThreatLevel  string            `json:"threat_level"`
    Category     string            `json:"category"`
    Metadata     map[string]string `json:"metadata,omitempty"`
}
```

FileScanResult 表示文件扫描结果结构。

**字段说明:**

- `FilePath` (string): 文件路径
- `FileType` (string): 文件类型
- `IsSuspicious` (bool): 是否可疑
- `Threats` ([]ThreatInfo): 威胁列表
- `ScanTime` (time.Time): 扫描时间
- `ScanDuration` (time.Duration): 扫描耗时
- `ThreatLevel` (string): 威胁等级
- `Category` (string): 类别
- `Metadata` (map[string]string): 元数据

### type ArchiveAnalysis

```go
type ArchiveAnalysis struct {
    Timestamp        time.Time `json:"timestamp"`
    IsValidArchive   bool      `json:"is_valid_archive"`
    ArchiveType      string    `json:"archive_type"`
    FileCount        int       `json:"file_count"`
    CompressedSize   int64     `json:"compressed_size"`
    UncompressedSize int64     `json:"uncompressed_size"`
    Files            []string  `json:"files"`
}
```

ArchiveAnalysis 表示压缩文件分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `IsValidArchive` (bool): 是否为有效压缩文件
- `ArchiveType` (string): 压缩文件类型
- `FileCount` (int): 文件数量
- `CompressedSize` (int64): 压缩后大小
- `UncompressedSize` (int64): 解压后大小
- `Files` ([]string): 文件列表

### type DocumentAnalysis

```go
type DocumentAnalysis struct {
    Timestamp       time.Time         `json:"timestamp"`
    IsValidDocument bool              `json:"is_valid_document"`
    DocumentType    string            `json:"document_type"`
    HasMacros       bool              `json:"has_macros"`
    MacroCount      int               `json:"macro_count"`
    EmbeddedObjects []string          `json:"embedded_objects"`
    Metadata        map[string]string `json:"metadata"`
}
```

DocumentAnalysis 表示文档分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `IsValidDocument` (bool): 是否为有效文档
- `DocumentType` (string): 文档类型
- `HasMacros` (bool): 是否包含宏
- `MacroCount` (int): 宏数量
- `EmbeddedObjects` ([]string): 嵌入对象列表
- `Metadata` (map[string]string): 元数据

### type BehaviorAnalysis

```go
type BehaviorAnalysis struct {
    Timestamp       time.Time `json:"timestamp"`
    NetworkActivity bool      `json:"network_activity"`
    FileOperations  bool      `json:"file_operations"`
    RegistryAccess  bool      `json:"registry_access"`
    ProcessCreation bool      `json:"process_creation"`
    APIHooking      bool      `json:"api_hooking"`
    AntiDebug       bool      `json:"anti_debug"`
    VMDetection     bool      `json:"vm_detection"`
}
```

BehaviorAnalysis 表示行为分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `NetworkActivity` (bool): 网络活动
- `FileOperations` (bool): 文件操作
- `RegistryAccess` (bool): 注册表访问
- `ProcessCreation` (bool): 进程创建
- `APIHooking` (bool): API 钩子
- `AntiDebug` (bool): 反调试
- `VMDetection` (bool): 虚拟机检测

### type SignatureAnalysis

```go
type SignatureAnalysis struct {
    Timestamp        time.Time         `json:"timestamp"`
    IsSigned         bool              `json:"is_signed"`
    SignerName       string            `json:"signer_name"`
    CertificateValid bool              `json:"certificate_valid"`
    SignatureValid   bool              `json:"signature_valid"`
    CertificateInfo  map[string]string `json:"certificate_info"`
}
```

SignatureAnalysis 表示签名分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `IsSigned` (bool): 是否已签名
- `SignerName` (string): 签名者名称
- `CertificateValid` (bool): 证书是否有效
- `SignatureValid` (bool): 签名是否有效
- `CertificateInfo` (map[string]string): 证书信息

### type EntropyAnalysis

```go
type EntropyAnalysis struct {
    Timestamp           time.Time          `json:"timestamp"`
    OverallEntropy      float64            `json:"overall_entropy"`
    SectionEntropies    map[string]float64 `json:"section_entropies"`
    HighEntropySections []string           `json:"high_entropy_sections"`
    EntropyThreshold    float64            `json:"entropy_threshold"`
}
```

EntropyAnalysis 表示熵值分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `OverallEntropy` (float64): 整体熵值
- `SectionEntropies` (map[string]float64): 各节熵值
- `HighEntropySections` ([]string): 高熵值节列表
- `EntropyThreshold` (float64): 熵值阈值

### type EncodingAnalysis

```go
type EncodingAnalysis struct {
    Timestamp         time.Time `json:"timestamp"`
    EncodingType      string    `json:"encoding_type"`
    IsEncoded         bool      `json:"is_encoded"`
    EncodingAlgorithm string    `json:"encoding_algorithm"`
}
```

EncodingAnalysis 表示编码分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `EncodingType` (string): 编码类型
- `IsEncoded` (bool): 是否已编码
- `EncodingAlgorithm` (string): 编码算法

### type ScanFileRequest

```go
type ScanFileRequest struct {
    Path string `json:"path" binding:"required"`
}
```

ScanFileRequest 表示扫描文件的请求结构。

**字段说明:**

- `Path` (string, 必需): 要扫描的文件路径

### type ScanResult

```go
type ScanResult struct {
    FilePath     string            `json:"file_path"`
    IsInfected   bool              `json:"is_infected"`
    Threats      []ThreatInfo      `json:"threats"`
    ScanTime     time.Time         `json:"scan_time"`
    ScanDuration time.Duration     `json:"scan_duration"`
    FileInfo     *FileInfo         `json:"file_info"`
    Metadata     map[string]string `json:"metadata,omitempty"`
}
```

ScanResult 表示文件扫描的结果结构。

**字段说明:**

- `FilePath` (string): 扫描的文件路径
- `IsInfected` (bool): 是否被感染
- `Threats` ([]ThreatInfo): 威胁信息列表
- `ScanTime` (time.Time): 扫描时间
- `ScanDuration` (time.Duration): 扫描耗时
- `FileInfo` (\*FileInfo): 文件信息
- `Metadata` (map[string]string): 元数据信息

### type ThreatInfo

```go
type ThreatInfo struct {
    RuleName    string `json:"rule_name"`
    Description string `json:"description"`
    Severity    string `json:"severity"`
    Category    string `json:"category"`
    Tags        string `json:"tags"`
}
```

ThreatInfo 表示威胁信息结构。

**字段说明:**

- `RuleName` (string): 规则名称
- `Description` (string): 威胁描述
- `Severity` (string): 严重程度
- `Category` (string): 威胁类别
- `Tags` (string): 标签信息

### type FileInfo

```go
type FileInfo struct {
    Path         string    `json:"path"`
    Name         string    `json:"name"`
    Size         int64     `json:"size"`
    IsDir        bool      `json:"is_dir"`
    ModTime      time.Time `json:"mod_time"`
    CreateTime   time.Time `json:"create_time"`
    AccessTime   time.Time `json:"access_time"`
    Permissions  string    `json:"permissions"`
    Owner        string    `json:"owner"`
    Group        string    `json:"group"`
    MD5          string    `json:"md5,omitempty"`
    SHA256       string    `json:"sha256,omitempty"`
    IsSuspicious bool      `json:"is_suspicious,omitempty"`
    ThreatLevel  string    `json:"threat_level,omitempty"`
}
```

FileInfo 表示文件信息结构。

**字段说明:**

- `Path` (string): 文件路径
- `Name` (string): 文件名
- `Size` (int64): 文件大小
- `IsDir` (bool): 是否为目录
- `ModTime` (time.Time): 修改时间
- `CreateTime` (time.Time): 创建时间
- `AccessTime` (time.Time): 访问时间
- `Permissions` (string): 权限信息
- `Owner` (string): 所有者
- `Group` (string): 用户组
- `MD5` (string): MD5 哈希值
- `SHA256` (string): SHA256 哈希值
- `IsSuspicious` (bool): 是否可疑
- `ThreatLevel` (string): 威胁等级

### type BufferScanRequest

```go
type BufferScanRequest struct {
    Data       string `json:"data"`
    Identifier string `json:"identifier"`
}
```

BufferScanRequest 表示缓冲区扫描请求结构。

**字段说明:**

- `Data` (string): 要扫描的数据内容
- `Identifier` (string): 数据标识符

### type FileListResponse

```go
type FileListResponse struct {
    Directory string     `json:"directory"`
    Files     []FileInfo `json:"files"`
    Count     int        `json:"count"`
    Error     string     `json:"error,omitempty"`
}
```

FileListResponse 表示文件列表响应结构。

**字段说明:**

- `Directory` (string): 目录路径
- `Files` ([]FileInfo): 文件列表
- `Count` (int): 文件数量
- `Error` (string): 错误信息

### type DirectoryScanRequest

```go
type DirectoryScanRequest struct {
    Path            string   `json:"path"`
    Recursive       bool     `json:"recursive"`
    MaxDepth        int      `json:"max_depth"`
    IncludePatterns []string `json:"include_patterns,omitempty"`
    ExcludePatterns []string `json:"exclude_patterns,omitempty"`
}
```

DirectoryScanRequest 表示目录扫描的请求结构。

**字段说明:**

- `Path` (string): 要扫描的目录路径
- `Recursive` (bool): 是否递归扫描子目录
- `MaxDepth` (int): 最大扫描深度
- `IncludePatterns` ([]string): 包含的文件模式
- `ExcludePatterns` ([]string): 排除的文件模式

### type FileHandler

```go
type FileHandler struct {
    securityService *services.SecurityService
    scanner         *security.Scanner
    logger          *logrus.Logger
}
```

FileHandler 表示一个文件处理处理器，提供文件扫描和分析功能。

**字段说明:**

- `securityService` (\*services.SecurityService): 安全服务实例
- `scanner` (\*security.Scanner): 扫描器实例
- `logger` (\*logrus.Logger): 日志记录器

### func NewFileHandler

```go
func NewFileHandler(securityService *services.SecurityService, scanner *security.Scanner, logger *logrus.Logger) *FileHandler
```

NewFileHandler 使用指定的安全服务、扫描器和日志记录器创建 FileHandler 的新实例。

**参数:**

- `securityService` (\*services.SecurityService): 安全服务实例
- `scanner` (\*security.Scanner): 扫描器实例
- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `*FileHandler`: 新创建的文件处理器实例

### func (h \*FileHandler) ScanFile

```go
func (h *FileHandler) ScanFile(c *gin.Context)
```

ScanFile 使用 Yara 规则对单个文件进行安全扫描，检测病毒、木马和其他威胁。返回扫描结果和威胁信息。

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**功能说明:**

1. 从请求体中解析文件路径
2. 验证文件路径的有效性
3. 执行 Yara 规则扫描
4. 计算文件哈希值
5. 返回扫描结果

### func (h \*FileHandler) GetFileInfo

```go
func (h *FileHandler) GetFileInfo(c *gin.Context)
```

GetFileInfo 获取指定文件的详细信息，包括文件大小、创建时间、修改时间、权限、所有者、哈希值和其他详细属性。

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**功能说明:**

1. 从路径参数中获取文件路径
2. 检查文件是否存在
3. 获取文件详细信息
4. 计算文件哈希值
5. 返回文件信息

### func (h \*FileHandler) calculateFileHashes

```go
func (h *FileHandler) calculateFileHashes(filePath string) (map[string]string, error)
```

calculateFileHashes 计算文件的 MD5 和 SHA256 哈希值。

**参数:**

- `filePath` (string): 文件路径

**返回值:**

- `map[string]string`: 哈希值映射
- `error`: 错误信息

**功能说明:**

1. 读取文件内容
2. 计算 MD5 哈希值
3. 计算 SHA256 哈希值
4. 返回哈希值结果

### func (h \*FileHandler) ScanDirectory

```go
func (h *FileHandler) ScanDirectory(c *gin.Context)
```

ScanDirectory 递归扫描指定目录中的所有文件，支持扫描深度配置和文件过滤，用于批量威胁检测。

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**功能说明:**

1. 从请求体中解析扫描配置
2. 验证目录路径的有效性
3. 递归扫描目录中的文件
4. 应用文件过滤规则
5. 执行 Yara 规则扫描
6. 返回扫描结果汇总

### func (h \*FileHandler) ScanBuffer

```go
func (h *FileHandler) ScanBuffer(c *gin.Context)
```

ScanBuffer 对内存缓冲区中的数据进行安全扫描，适用于检测网络载荷和内存中的恶意代码。

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**功能说明:**

1. 从请求体中解析数据内容
2. 验证数据有效性
3. 执行 Yara 规则扫描
4. 检测恶意字符和数据结构
5. 返回扫描结果

### func (h \*FileHandler) GetFileList

```go
func (h *FileHandler) GetFileList(c *gin.Context)
```

GetFileList 获取指定目录中的文件列表，提供基本的文件信息，如名称、大小和类型。

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**功能说明:**

1. 从查询参数中获取目录路径
2. 验证目录路径的有效性
3. 扫描目录中的文件
4. 应用递归和深度限制
5. 返回文件列表

### func (h \*FileHandler) getFileList

```go
func (h *FileHandler) getFileList(directory string, recursive bool) ([]*models.FileInfo, error)
```

getFileList 内部方法，获取指定目录的文件列表。

**参数:**

- `dirPath` (string): 目录路径
- `recursive` (bool): 是否递归

**返回值:**

- `[]*models.FileInfo`: 文件信息列表
- `error`: 错误信息

**功能说明:**

1. 读取目录内容
2. 过滤文件和目录
3. 获取文件详细信息
4. 返回文件列表

### type FileManager

```go
type FileManager struct {
    logger *logrus.Logger
}
```

FileManager 表示文件管理器，提供文件操作功能。

**字段说明:**

- `logger` (\*logrus.Logger): 日志记录器

### func (m \*FileManager) CopyFile

```go
func (m *FileManager) CopyFile(source, dest string) error
```

CopyFile 复制文件从源路径到目标路径，支持文件覆盖和权限保持。

**参数:**

- `source` (string): 源文件路径
- `dest` (string): 目标文件路径

**返回值:**

- `error`: 错误信息

### func (m \*FileManager) MoveFile

```go
func (m *FileManager) MoveFile(source, dest string) error
```

MoveFile 移动文件从源路径到目标路径，支持文件重命名。

**参数:**

- `source` (string): 源文件路径
- `dest` (string): 目标文件路径

**返回值:**

- `error`: 错误信息

### func (m \*FileManager) DeleteFile

```go
func (m *FileManager) DeleteFile(filePath string) error
```

DeleteFile 删除指定路径的文件，支持强制删除。

**参数:**

- `filePath` (string): 文件路径

**返回值:**

- `error`: 错误信息

### func (m \*FileManager) GetFileHash

```go
func (m *FileManager) GetFileHash(filePath, algorithm string) (string, error)
```

GetFileHash 获取指定文件的哈希值，支持多种哈希算法。

**参数:**

- `filePath` (string): 文件路径
- `algorithm` (string): 哈希算法

**返回值:**

- `string`: 哈希值
- `error`: 错误信息

### func (m \*FileManager) GetFileHashes

```go
func (m *FileManager) GetFileHashes(filePath string) (map[string]string, error)
```

GetFileHashes 获取指定文件的所有哈希值（MD5、SHA1、SHA256、SHA512）。

**参数:**

- `filePath` (string): 文件路径

**返回值:**

- `map[string]string`: 哈希值映射
- `error`: 错误信息

### func (m \*FileManager) VerifyFileHash

```go
func (m *FileManager) VerifyFileHash(filePath, algorithm, expectedHash string) (bool, error)
```

VerifyFileHash 验证文件的哈希值是否与预期值匹配。

**参数:**

- `filePath` (string): 文件路径
- `algorithm` (string): 哈希算法
- `expectedHash` (string): 期望的哈希值

**返回值:**

- `bool`: 验证结果
- `error`: 错误信息

---

## 🔍 进程管理模块

### type ProcessEvent

```go
type ProcessEvent struct {
    PID         int32                  `json:"pid"`
    ProcessName string                 `json:"process_name"`
    EventType   string                 `json:"event_type"`
    Timestamp   time.Time              `json:"timestamp"`
    Details     map[string]interface{} `json:"details"`
}
```

ProcessEvent 表示进程事件结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `ProcessName` (string): 进程名称
- `EventType` (string): 事件类型
- `Timestamp` (time.Time): 时间戳
- `Details` (map[string]interface{}): 详细信息

### type DeepProcessAnalysis

```go
type DeepProcessAnalysis struct {
    PID              int32                    `json:"pid"`
    Timestamp        time.Time                `json:"timestamp"`
    IsSuspicious     bool                     `json:"is_suspicious"`
    ThreatLevel      string                   `json:"threat_level"`
    RiskScore        int                      `json:"risk_score"`
    RiskFactors      []string                 `json:"risk_factors"`
    BasicInfo        *BasicProcessInfo        `json:"basic_info"`
    MemoryAnalysis   *MemoryAnalysis          `json:"memory_analysis"`
    NetworkAnalysis  *NetworkAnalysis         `json:"network_analysis"`
    FileAnalysis     *FileAnalysis            `json:"file_analysis"`
    BehaviorAnalysis *ProcessBehaviorAnalysis `json:"behavior_analysis"`
    ThreatAnalysis   *ProcessThreatAnalysis   `json:"threat_analysis"`
}
```

DeepProcessAnalysis 表示深度进程分析结果结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `Timestamp` (time.Time): 分析时间戳
- `IsSuspicious` (bool): 是否可疑
- `ThreatLevel` (string): 威胁等级
- `RiskScore` (int): 风险评分
- `RiskFactors` ([]string): 风险因素列表
- `BasicInfo` (\*BasicProcessInfo): 基础进程信息
- `MemoryAnalysis` (\*MemoryAnalysis): 内存分析
- `NetworkAnalysis` (\*NetworkAnalysis): 网络分析
- `FileAnalysis` (\*FileAnalysis): 文件分析
- `BehaviorAnalysis` (\*ProcessBehaviorAnalysis): 行为分析
- `ThreatAnalysis` (\*ProcessThreatAnalysis): 威胁分析

### type BasicProcessInfo

```go
type BasicProcessInfo struct {
    Timestamp        time.Time `json:"timestamp"`
    Name             string    `json:"name"`
    CommandLine      string    `json:"command_line"`
    Executable       string    `json:"executable"`
    WorkingDirectory string    `json:"working_directory"`
    Status           string    `json:"status"`
    CreateTime       time.Time `json:"create_time"`
    Username         string    `json:"username"`
    ParentPID        int32     `json:"parent_pid"`
    CPUPercent       float64   `json:"cpu_percent"`
    MemoryPercent    float32   `json:"memory_percent"`
    NumThreads       int32     `json:"num_threads"`
    Priority         int32     `json:"priority"`
}
```

BasicProcessInfo 表示基础进程信息结构。

**字段说明:**

- `Timestamp` (time.Time): 时间戳
- `Name` (string): 进程名称
- `CommandLine` (string): 命令行
- `Executable` (string): 可执行文件路径
- `WorkingDirectory` (string): 工作目录
- `Status` (string): 进程状态
- `CreateTime` (time.Time): 创建时间
- `Username` (string): 用户名
- `ParentPID` (int32): 父进程 ID
- `CPUPercent` (float64): CPU 使用率
- `MemoryPercent` (float32): 内存使用率
- `NumThreads` (int32): 线程数
- `Priority` (int32): 优先级

### type ProcessInfo

```go
type ProcessInfo struct {
    PID           int32     `json:"pid"`
    Name          string    `json:"name"`
    Executable    string    `json:"executable"`
    CommandLine   string    `json:"command_line"`
    CreateTime    time.Time `json:"create_time"`
    Status        string    `json:"status"`
    Username      string    `json:"username"`
    ParentPID     int32     `json:"parent_pid"`
    WorkingDir    string    `json:"working_directory"`
    CPUPercent    float64   `json:"cpu_percent"`
    MemoryPercent float64   `json:"memory_percent"`
    NumThreads    int32     `json:"num_threads"`
    Priority      int32     `json:"priority"`
}
```

ProcessInfo 表示进程信息结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `Name` (string): 进程名称
- `Executable` (string): 可执行文件路径
- `CommandLine` (string): 命令行参数
- `CreateTime` (time.Time): 创建时间
- `Status` (string): 进程状态
- `Username` (string): 用户名
- `ParentPID` (int32): 父进程 ID
- `WorkingDir` (string): 工作目录
- `CPUPercent` (float64): CPU 使用率
- `MemoryPercent` (float64): 内存使用率
- `NumThreads` (int32): 线程数
- `Priority` (int32): 优先级

### type ProcessThreatAnalysis

```go
type ProcessThreatAnalysis struct {
    Timestamp         time.Time          `json:"timestamp"`
    IsSuspicious      bool               `json:"is_suspicious"`
    ThreatLevel       string             `json:"threat_level"`
    RiskScore         int                `json:"risk_score"`
    RiskFactors       []string           `json:"risk_factors"`
    SuspiciousBehaviors []string         `json:"suspicious_behaviors"`
    NetworkThreats    []string           `json:"network_threats"`
    FileThreats       []string           `json:"file_threats"`
}
```

ProcessThreatAnalysis 表示进程威胁分析结构。

**字段说明:**

- `Timestamp` (time.Time): 分析时间戳
- `IsSuspicious` (bool): 是否可疑
- `ThreatLevel` (string): 威胁等级
- `RiskScore` (int): 风险评分
- `RiskFactors` ([]string): 风险因素列表
- `SuspiciousBehaviors` ([]string): 可疑行为列表
- `NetworkThreats` ([]string): 网络威胁列表
- `FileThreats` ([]string): 文件威胁列表

### type AddrInfo

```go
type AddrInfo struct {
    IP   string `json:"ip"`
    Port int    `json:"port"`
}
```

AddrInfo 表示地址信息结构。

**字段说明:**

- `IP` (string): IP 地址
- `Port` (int): 端口号

### type ProcessMemoryInfo

```go
type ProcessMemoryInfo struct {
    PID       int32   `json:"pid"`
    RSS       uint64  `json:"rss"`
    VMS       uint64  `json:"vms"`
    Percent   float64 `json:"percent"`
    Available uint64  `json:"available"`
    Used      uint64  `json:"used"`
    Free      uint64  `json:"free"`
    Total     uint64  `json:"total"`
}
```

ProcessMemoryInfo 表示进程内存信息结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `RSS` (uint64): 物理内存使用量
- `VMS` (uint64): 虚拟内存使用量
- `Percent` (float64): 内存使用百分比
- `Available` (uint64): 可用内存
- `Used` (uint64): 已使用内存
- `Free` (uint64): 空闲内存
- `Total` (uint64): 总内存

### type ProcessConnectionsResponse

```go
type ProcessConnectionsResponse struct {
    PID              int                    `json:"pid"`
    ProcessName      string                 `json:"process_name"`
    TotalConnections int                    `json:"total_connections"`
    Connections      []ProcessConnectionInfo `json:"connections"`
}
```

ProcessConnectionsResponse 表示进程连接列表响应结构。

**字段说明:**

- `PID` (int): 进程 ID
- `ProcessName` (string): 进程名称
- `TotalConnections` (int): 总连接数量
- `Connections` ([]ProcessConnectionInfo): 连接信息列表

### type ProcessMemoryResponse

```go
type ProcessMemoryResponse struct {
    PID        int32              `json:"pid"`
    MemoryInfo ProcessMemoryInfo  `json:"memory_info"`
}
```

ProcessMemoryResponse 表示进程内存响应结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `MemoryInfo` (ProcessMemoryInfo): 内存信息

### type ProcessModuleInfo

```go
type ProcessModuleInfo struct {
    Name        string `json:"name"`
    Path        string `json:"path"`
    BaseAddress uint64 `json:"base_address"`
    Size        uint64 `json:"size"`
    Version     string `json:"version"`
    Description string `json:"description"`
}
```

ProcessModuleInfo 表示进程模块信息结构。

**字段说明:**

- `Name` (string): 模块名称
- `Path` (string): 模块路径
- `BaseAddress` (uint64): 基地址
- `Size` (uint64): 模块大小
- `Version` (string): 版本信息
- `Description` (string): 描述信息

### type ProcessListResponse

```go
type ProcessListResponse struct {
    TotalCount int           `json:"total_count"`
    Processes  []ProcessInfo `json:"processes"`
}
```

ProcessListResponse 表示进程列表响应结构。

**字段说明:**

- `TotalCount` (int): 总进程数量
- `Processes` ([]ProcessInfo): 进程信息列表

### type ProcessDetailResponse

```go
type ProcessDetailResponse struct {
    BasicInfo    ProcessInfo      `json:"basic_info"`
    MemoryInfo   ProcessMemoryInfo `json:"memory_info"`
    CPUInfo      CPUInfo          `json:"cpu_info"`
    FileInfo     FileInfo         `json:"file_info"`
    NetworkInfo  NetworkInfo      `json:"network_info"`
    ModuleInfo   ModuleInfo       `json:"module_info"`
    ThreadInfo   ThreadInfo       `json:"thread_info"`
    HandleInfo   HandleInfo       `json:"handle_info"`
}
```

ProcessDetailResponse 表示进程详细信息响应结构。

**字段说明:**

- `BasicInfo` (ProcessInfo): 基本进程信息
- `MemoryInfo` (ProcessMemoryInfo): 内存信息
- `CPUInfo` (CPUInfo): CPU 信息
- `FileInfo` (FileInfo): 文件信息
- `NetworkInfo` (NetworkInfo): 网络信息
- `ModuleInfo` (ModuleInfo): 模块信息
- `ThreadInfo` (ThreadInfo): 线程信息
- `HandleInfo` (HandleInfo): 句柄信息

### type StartProcessRequest

```go
type StartProcessRequest struct {
    Command   string   `json:"command" binding:"required"`
    Args      []string `json:"args"`
    WorkingDir string  `json:"working_dir"`
}
```

StartProcessRequest 表示启动进程请求结构。

**字段说明:**

- `Command` (string, 必需): 要执行的命令
- `Args` ([]string): 命令行参数列表
- `WorkingDir` (string): 工作目录

### type StartProcessResponse

```go
type StartProcessResponse struct {
    PID       int32     `json:"pid"`
    Command   string    `json:"command"`
    StartTime time.Time `json:"start_time"`
}
```

StartProcessResponse 表示启动进程响应结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `Command` (string): 执行的命令
- `StartTime` (time.Time): 启动时间

### type ProcessStatusResponse

```go
type ProcessStatusResponse struct {
    PID           int32     `json:"pid"`
    Status        string    `json:"status"`
    OperationTime time.Time `json:"operation_time"`
}
```

ProcessStatusResponse 表示进程状态响应结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `Status` (string): 进程状态
- `OperationTime` (time.Time): 操作时间

### type ProcessModulesResponse

```go
type ProcessModulesResponse struct {
    PID           int32              `json:"pid"`
    ProcessName   string             `json:"process_name"`
    TotalModules  int                `json:"total_modules"`
    Modules       []ProcessModuleInfo `json:"modules"`
}
```

ProcessModulesResponse 表示进程模块列表响应结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `ProcessName` (string): 进程名称
- `TotalModules` (int): 总模块数量
- `Modules` ([]ProcessModuleInfo): 模块信息列表

### type ProcessConnectionInfo

```go
type ProcessConnectionInfo struct {
    ID           string `json:"id"`
    LocalAddr    string `json:"local_addr"`
    RemoteAddr   string `json:"remote_addr"`
    LocalPort    int    `json:"local_port"`
    RemotePort   int    `json:"remote_port"`
    Status       string `json:"status"`
    Type         string `json:"type"`
}
```

ProcessConnectionInfo 表示进程连接信息结构。

**字段说明:**

- `ID` (string): 连接 ID
- `LocalAddr` (string): 本地地址
- `RemoteAddr` (string): 远程地址
- `LocalPort` (int): 本地端口
- `RemotePort` (int): 远程端口
- `Status` (string): 连接状态
- `Type` (string): 连接类型

### type ProcessRunningStatus

```go
type ProcessRunningStatus struct {
    PID      int32  `json:"pid"`
    IsRunning bool  `json:"is_running"`
    Status   string `json:"status"`
}
```

ProcessRunningStatus 表示进程运行状态结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `IsRunning` (bool): 是否正在运行
- `Status` (string): 进程状态

### type ProcessRunningResponse

```go
type ProcessRunningResponse struct {
    PID      int32  `json:"pid"`
    IsRunning bool  `json:"is_running"`
    Status   string `json:"status"`
}
```

ProcessRunningResponse 表示进程运行状态响应结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `IsRunning` (bool): 是否正在运行
- `Status` (string): 进程状态

### type ProcessChildrenResponse

```go
type ProcessChildrenResponse struct {
    PID       int32         `json:"pid"`
    Children  []ProcessInfo `json:"children"`
    Count     int           `json:"count"`
}
```

ProcessChildrenResponse 表示进程子进程响应结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `Children` ([]ProcessInfo): 子进程列表
- `Count` (int): 子进程数量

### type SystemModuleInfo

```go
type SystemModuleInfo struct {
    Name           string `json:"name"`
    ProcessCount   int    `json:"process_count"`
    TotalInstances int    `json:"total_instances"`
    IsSuspicious   bool   `json:"is_suspicious"`
}
```

SystemModuleInfo 表示系统模块信息结构。

**字段说明:**

- `Name` (string): 模块名称
- `ProcessCount` (int): 使用该模块的进程数量
- `TotalInstances` (int): 总实例数
- `IsSuspicious` (bool): 是否可疑

### type SystemModulesResponse

```go
type SystemModulesResponse struct {
    TotalModules int                `json:"total_modules"`
    Modules      []SystemModuleInfo `json:"modules"`
}
```

SystemModulesResponse 表示系统模块列表响应结构。

**字段说明:**

- `TotalModules` (int): 总模块数量
- `Modules` ([]SystemModuleInfo): 模块信息列表

### type SystemModuleDetailResponse

```go
type SystemModuleDetailResponse struct {
    Name           string         `json:"name"`
    Path           string         `json:"path"`
    Size           uint64         `json:"size"`
    ProcessCount   int            `json:"process_count"`
    AffectedProcesses []ProcessInfo `json:"affected_processes"`
}
```

SystemModuleDetailResponse 表示系统模块详细信息响应结构。

**字段说明:**

- `Name` (string): 模块名称
- `Path` (string): 模块路径
- `Size` (uint64): 模块大小
- `ProcessCount` (int): 使用该模块的进程数量
- `AffectedProcesses` ([]ProcessInfo): 受影响的进程列表

### type ModuleStatusResponse

```go
type ModuleStatusResponse struct {
    ModuleName string `json:"module_name"`
    Status     string `json:"status"`
    Message    string `json:"message"`
}
```

ModuleStatusResponse 表示模块状态响应结构。

**字段说明:**

- `ModuleName` (string): 模块名称
- `Status` (string): 操作状态
- `Message` (string): 状态消息

### type MonitoringStatusResponse

```go
type MonitoringStatusResponse struct {
    IsEnabled    bool      `json:"is_enabled"`
    StartTime    time.Time `json:"start_time"`
    ProcessCount int       `json:"process_count"`
}
```

MonitoringStatusResponse 表示监控状态响应结构。

**字段说明:**

- `IsEnabled` (bool): 是否启用监控
- `StartTime` (time.Time): 开始监控时间
- `ProcessCount` (int): 监控的进程数量

### type MonitoredProcessInfo

```go
type MonitoredProcessInfo struct {
    PID         int32     `json:"pid"`
    Name        string    `json:"name"`
    StartTime   time.Time `json:"start_time"`
    LastSeen    time.Time `json:"last_seen"`
    EventCount  int       `json:"event_count"`
}
```

MonitoredProcessInfo 表示被监控进程信息结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `Name` (string): 进程名称
- `StartTime` (time.Time): 开始监控时间
- `LastSeen` (time.Time): 最后看到时间
- `EventCount` (int): 事件数量

### type MonitoredProcessListResponse

```go
type MonitoredProcessListResponse struct {
    TotalCount int                    `json:"total_count"`
    Processes  []MonitoredProcessInfo `json:"processes"`
}
```

MonitoredProcessListResponse 表示被监控进程列表响应结构。

**字段说明:**

- `TotalCount` (int): 总进程数量
- `Processes` ([]MonitoredProcessInfo): 被监控进程列表

### type ProcessStatistics

```go
type ProcessStatistics struct {
    TotalProcesses    int     `json:"total_processes"`
    RunningProcesses  int     `json:"running_processes"`
    SuspendedProcesses int    `json:"suspended_processes"`
    TotalMemoryUsage  uint64  `json:"total_memory_usage"`
    TotalCPUUsage     float64 `json:"total_cpu_usage"`
}
```

ProcessStatistics 表示进程统计信息结构。

**字段说明:**

- `TotalProcesses` (int): 总进程数
- `RunningProcesses` (int): 运行中进程数
- `SuspendedProcesses` (int): 暂停进程数
- `TotalMemoryUsage` (uint64): 总内存使用量
- `TotalCPUUsage` (float64): 总 CPU 使用率

### type TopProcessInfo

```go
type TopProcessInfo struct {
    PID           int32   `json:"pid"`
    Name          string  `json:"name"`
    CPUPercent    float64 `json:"cpu_percent"`
    MemoryPercent float64 `json:"memory_percent"`
    Rank          int     `json:"rank"`
}
```

TopProcessInfo 表示顶级进程信息结构。

**字段说明:**

- `PID` (int32): 进程 ID
- `Name` (string): 进程名称
- `CPUPercent` (float64): CPU 使用率
- `MemoryPercent` (float64): 内存使用率
- `Rank` (int): 排名

### type ProcessHandler

```go
type ProcessHandler struct {
    manager *process.Manager
    logger  *logrus.Logger
}
```

ProcessHandler 表示一个进程管理处理器，提供进程监控和控制功能。

**字段说明:**

- `manager` (\*process.Manager): 进程管理器实例
- `logger` (\*logrus.Logger): 日志记录器

### func NewProcessHandler

```go
func NewProcessHandler(processService *services.ProcessService, logger *logrus.Logger) *ProcessHandler
```

NewProcessHandler 使用指定的进程服务和日志记录器创建 ProcessHandler 的新实例。

**参数:**

- `processService` (\*services.ProcessService): 进程服务实例
- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `*ProcessHandler`: 新创建的进程处理器实例

### func (h \*ProcessHandler) GetProcesses

```go
func (h *ProcessHandler) GetProcesses(c *gin.Context)
```

GetProcesses 获取系统中所有进程的基本信息，包括 PID、进程名称、CPU 使用率、内存使用率和其他系统指标。

**功能说明:**

1. 调用进程服务获取进程列表
2. 处理响应数据
3. 返回进程列表和总数信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含进程列表的 JSON 响应，包括总数和进程详细信息

### func (h \*ProcessHandler) GetProcessByPID

```go
func (h *ProcessHandler) GetProcessByPID(c *gin.Context)
```

GetProcessByPID 通过 PID 获取特定进程的详细信息，包括进程状态、启动时间、命令行参数、内存使用情况和其他详细指标。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性
3. 调用进程管理器获取进程详细信息
4. 返回进程的基本信息和内存信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含进程详细信息的 JSON 响应，包括基本信息和内存信息

### func (h \*ProcessHandler) StartProcess

```go
func (h *ProcessHandler) StartProcess(c *gin.Context)
```

StartProcess 使用指定的命令和参数启动新进程，支持工作目录和命令行参数配置。

**功能说明:**

1. 解析请求体中的命令、参数和工作目录
2. 验证命令的有效性和安全性
3. 调用进程管理器启动新进程
4. 返回新进程的 PID 和状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含新进程信息的 JSON 响应，包括 PID、命令和状态

### func (h \*ProcessHandler) KillProcess

```go
func (h *ProcessHandler) KillProcess(c *gin.Context)
```

KillProcess 通过 PID 强制终止指定进程，适用于终止恶意进程或失控进程。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性和进程存在性
3. 调用进程管理器强制终止进程
4. 返回终止结果和状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含终止结果的 JSON 响应，包括 PID 和状态

### func (h \*ProcessHandler) SuspendProcess

```go
func (h *ProcessHandler) SuspendProcess(c *gin.Context)
```

SuspendProcess 暂停指定进程的执行，进程保持在内存中但不再占用 CPU 资源。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性和进程存在性
3. 调用进程管理器暂停进程执行
4. 返回暂停结果和状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含暂停结果的 JSON 响应，包括 PID 和状态

### func (h \*ProcessHandler) ResumeProcess

```go
func (h *ProcessHandler) ResumeProcess(c *gin.Context)
```

ResumeProcess 恢复之前被挂起的进程，使其重新开始执行。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性和进程存在性
3. 调用进程管理器恢复进程执行
4. 返回恢复结果和状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含恢复结果的 JSON 响应，包括 PID 和状态

### func (h \*ProcessHandler) GetProcessModules

```go
func (h *ProcessHandler) GetProcessModules(c *gin.Context)
```

GetProcessModules 获取特定进程加载的所有 DLL 模块信息，包括模块路径、基地址、大小和其他模块详情。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性和进程存在性
3. 调用进程管理器获取进程模块列表
4. 返回模块列表和数量信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含模块列表的 JSON 响应，包括 PID、模块信息和数量

### func (h \*ProcessHandler) GetProcessConnections

```go
func (h *ProcessHandler) GetProcessConnections(c *gin.Context)
```

GetProcessConnections 获取特定进程的所有网络连接，包括 TCP/UDP 连接、本地地址、远程地址和连接状态。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性和进程存在性
3. 调用网络管理器获取进程连接信息
4. 返回连接列表和数量信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含连接列表的 JSON 响应，包括 PID、连接信息和数量

### func (h \*ProcessHandler) GetProcessMemoryInfo

```go
func (h *ProcessHandler) GetProcessMemoryInfo(c *gin.Context)
```

GetProcessMemoryInfo 获取特定进程的详细内存使用信息，包括物理内存、虚拟内存、内存峰值和其他内存指标。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性和进程存在性
3. 调用进程管理器获取内存使用信息
4. 返回内存详细信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含内存信息的 JSON 响应，包括 PID 和内存详细信息

### func (h \*ProcessHandler) IsProcessRunning

```go
func (h *ProcessHandler) IsProcessRunning(c *gin.Context)
```

IsProcessRunning 检查指定 PID 的进程是否仍在运行，返回运行状态。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性
3. 调用进程管理器检查进程运行状态
4. 返回进程运行状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含进程运行状态的 JSON 响应，包括 PID 和运行状态

### func (h \*ProcessHandler) GetProcessChildren

```go
func (h *ProcessHandler) GetProcessChildren(c *gin.Context)
```

GetProcessChildren 获取指定进程的所有子进程，用于分析进程树结构。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性和进程存在性
3. 调用进程管理器获取子进程列表
4. 返回子进程列表和数量信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含子进程列表的 JSON 响应，包括 PID、子进程信息和数量

### func (h \*ProcessHandler) GetSystemModules

```go
func (h *ProcessHandler) GetSystemModules(c *gin.Context)
```

GetSystemModules 获取系统中进程加载的所有模块的统计信息，包括模块名称和受影响的进程数量。

**功能说明:**

1. 调用进程管理器获取系统模块统计信息
2. 统计每个模块被多少进程使用
3. 返回模块列表和受影响进程信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含系统模块统计信息的 JSON 响应，包括总数、模块信息和受影响进程

### func (h \*ProcessHandler) GetModuleInfo

```go
func (h *ProcessHandler) GetModuleInfo(c *gin.Context)
```

GetModuleInfo 获取特定模块的详细信息，包括模块路径、大小和受影响进程列表。

**功能说明:**

1. 从路径参数中获取模块名称
2. 验证模块名称的有效性
3. 调用进程管理器获取模块详细信息
4. 返回模块信息和受影响进程列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `module_name`: 模块名称 (字符串)

**返回值:**

- 返回包含模块详细信息的 JSON 响应，包括名称、路径、大小和受影响进程

### func (h \*ProcessHandler) SuspendModule

```go
func (h *ProcessHandler) SuspendModule(c *gin.Context)
```

SuspendModule 暂停所有进程中指定模块的加载，用于隔离可疑模块。

**功能说明:**

1. 从路径参数中获取模块名称
2. 验证模块名称的有效性
3. 调用进程管理器暂停模块在所有进程中的加载
4. 返回暂停结果和受影响进程信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `module_name`: 模块名称 (字符串)

**返回值:**

- 返回包含暂停结果的 JSON 响应，包括模块名称、状态和受影响进程

### func (h \*ProcessHandler) ResumeModule

```go
func (h *ProcessHandler) ResumeModule(c *gin.Context)
```

ResumeModule 恢复之前被挂起的模块，使其在所有相关进程中重新可用。

**功能说明:**

1. 从路径参数中获取模块名称
2. 验证模块名称的有效性
3. 调用进程管理器恢复模块在所有进程中的加载
4. 返回恢复结果和状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `module_name`: 模块名称 (字符串)

**返回值:**

- 返回包含恢复结果的 JSON 响应，包括模块名称和状态

### func (h \*ProcessHandler) KillModule

```go
func (h *ProcessHandler) KillModule(c *gin.Context)
```

KillModule 强制从所有已加载它的进程中卸载指定模块。

**功能说明:**

1. 从路径参数中获取模块名称
2. 验证模块名称的有效性
3. 调用进程管理器强制卸载模块
4. 返回卸载结果和受影响进程信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `module_name`: 模块名称 (字符串)

**返回值:**

- 返回包含卸载结果的 JSON 响应，包括模块名称、状态和受影响进程

---

## 🔍 进程监控扩展功能

### func EnableProcessMonitoring

```go
func EnableProcessMonitoring(c *gin.Context)
```

EnableProcessMonitoring 启用进程监控功能，开始实时监控系统进程。

**功能说明:**

1. 检查当前监控状态
2. 启动进程监控服务
3. 记录监控开始时间
4. 返回监控启用结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含监控状态的 JSON 响应，包括状态和开始时间

### func DisableProcessMonitoring

```go
func DisableProcessMonitoring(c *gin.Context)
```

DisableProcessMonitoring 禁用进程监控功能，停止实时监控。

**功能说明:**

1. 检查当前监控状态
2. 停止进程监控服务
3. 记录监控停止时间
4. 返回监控禁用结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含监控状态的 JSON 响应，包括状态和停止时间

### func GetMonitoredProcesses

```go
func GetMonitoredProcesses(c *gin.Context)
```

GetMonitoredProcesses 获取当前被监控的进程列表。

**功能说明:**

1. 检查监控服务状态
2. 获取当前被监控的进程列表
3. 返回监控进程的详细信息
4. 统计监控进程总数

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含监控进程列表的 JSON 响应，包括进程信息和总数

### func GetProcessStatistics

```go
func GetProcessStatistics(c *gin.Context)
```

GetProcessStatistics 获取进程统计信息，包括总数、运行中数量、内存使用等。

**功能说明:**

1. 获取系统进程总数统计
2. 统计运行中的进程数量
3. 计算总内存和 CPU 使用情况
4. 统计被监控的进程数量
5. 返回综合统计信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含进程统计信息的 JSON 响应，包括总数、运行中数量、内存使用等

---

## 🔧 注册表管理模块

### type RegistryKey

```go
type RegistryKey struct {
    Path         string            `json:"path"`
    Name         string            `json:"name"`
    Type         string            `json:"type"`
    Value        interface{}       `json:"value"`
    SubKeys      []string          `json:"sub_keys,omitempty"`
    Values       map[string]string `json:"values,omitempty"`
    LastModified time.Time         `json:"last_modified"`
}
```

RegistryKey 表示注册表键信息结构。

**字段说明:**

- `Path` (string): 注册表路径
- `Name` (string): 键名称
- `Type` (string): 键类型
- `Value` (interface{}): 键值
- `SubKeys` ([]string): 子键列表
- `Values` (map[string]string): 值映射
- `LastModified` (time.Time): 最后修改时间

### type CreateRegistryKeyRequest

```go
type CreateRegistryKeyRequest struct {
    Path string `json:"path" binding:"required"`
}
```

CreateRegistryKeyRequest 表示创建注册表键请求结构。

**字段说明:**

- `Path` (string): 要创建的注册表键路径，必需字段

### type SetRegistryValueRequest

```go
type SetRegistryValueRequest struct {
    Path  string      `json:"path" binding:"required"`
    Name  string      `json:"name" binding:"required"`
    Value interface{} `json:"value" binding:"required"`
    Type  string      `json:"type" binding:"required"`
}
```

SetRegistryValueRequest 表示设置注册表值请求结构。

**字段说明:**

- `Path` (string): 注册表键路径，必需字段
- `Name` (string): 值名称，必需字段
- `Value` (interface{}): 值内容，必需字段
- `Type` (string): 值类型，必需字段

### type SearchRegistryRequest

```go
type SearchRegistryRequest struct {
    Root         string `json:"root" binding:"required"`
    Pattern      string `json:"pattern"`
    ValuePattern string `json:"value_pattern"`
}
```

SearchRegistryRequest 表示搜索注册表请求结构。

**字段说明:**

- `Root` (string): 搜索根路径，必需字段
- `Pattern` (string): 搜索模式
- `ValuePattern` (string): 值搜索模式

### type RegistryKeyResponse

```go
type RegistryKeyResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    RegistryKey `json:"data"`
}
```

RegistryKeyResponse 表示注册表键响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (RegistryKey): 注册表键数据

### type RegistryKeyInfo

```go
type RegistryKeyInfo struct {
    Path         string            `json:"path"`
    Name         string            `json:"name"`
    Type         string            `json:"type"`
    Value        interface{}       `json:"value"`
    SubKeys      []string          `json:"sub_keys,omitempty"`
    Values       map[string]string `json:"values,omitempty"`
    LastModified time.Time         `json:"last_modified"`
}
```

RegistryKeyInfo 表示注册表键信息结构。

**字段说明:**

- `Path` (string): 注册表路径
- `Name` (string): 键名称
- `Type` (string): 键类型
- `Value` (interface{}): 键值
- `SubKeys` ([]string): 子键列表
- `Values` (map[string]string): 值映射
- `LastModified` (time.Time): 最后修改时间

### type RegistryKeysResponse

```go
type RegistryKeysResponse struct {
    Code    int           `json:"code"`
    Message string        `json:"message"`
    Data    []RegistryKey `json:"data"`
}
```

RegistryKeysResponse 表示注册表键列表响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]RegistryKey): 注册表键列表

### type RegistrySubKey

```go
type RegistrySubKey struct {
    Path         string    `json:"path"`
    Name         string    `json:"name"`
    Type         string    `json:"type"`
    LastModified time.Time `json:"last_modified"`
}
```

RegistrySubKey 表示注册表子键结构。

**字段说明:**

- `Path` (string): 注册表路径
- `Name` (string): 键名称
- `Type` (string): 键类型
- `LastModified` (time.Time): 最后修改时间

### type RegistryValuesResponse

```go
type RegistryValuesResponse struct {
    Code    int            `json:"code"`
    Message string         `json:"message"`
    Data    []RegistryValue `json:"data"`
}
```

RegistryValuesResponse 表示注册表值列表响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]RegistryValue): 注册表值列表

### type RegistryValue

```go
type RegistryValue struct {
    Path  string      `json:"path"`
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
    Type  string      `json:"type"`
}
```

RegistryValue 表示注册表值结构。

**字段说明:**

- `Path` (string): 注册表路径
- `Name` (string): 值名称
- `Value` (interface{}): 值内容
- `Type` (string): 值类型

### type CreateRegistryKeyResponse

```go
type CreateRegistryKeyResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

CreateRegistryKeyResponse 表示创建注册表键响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type DeleteRegistryKeyResponse

```go
type DeleteRegistryKeyResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

DeleteRegistryKeyResponse 表示删除注册表键响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type RegistryValueResponse

```go
type RegistryValueResponse struct {
    Code    int          `json:"code"`
    Message string       `json:"message"`
    Data    RegistryValue `json:"data"`
}
```

RegistryValueResponse 表示注册表值响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (RegistryValue): 注册表值数据

### type RegistryValueDetail

```go
type RegistryValueDetail struct {
    Path  string      `json:"path"`
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
    Type  string      `json:"type"`
}
```

RegistryValueDetail 表示注册表值详细信息结构。

**字段说明:**

- `Path` (string): 注册表路径
- `Name` (string): 值名称
- `Value` (interface{}): 值内容
- `Type` (string): 值类型

### type SetRegistryValueResponse

```go
type SetRegistryValueResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

SetRegistryValueResponse 表示设置注册表值响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type DeleteRegistryValueResponse

```go
type DeleteRegistryValueResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

DeleteRegistryValueResponse 表示删除注册表值响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type RegistrySearchResponse

```go
type RegistrySearchResponse struct {
    Code    int                  `json:"code"`
    Message string               `json:"message"`
    Data    []RegistrySearchResult `json:"data"`
}
```

RegistrySearchResponse 表示注册表搜索响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]RegistrySearchResult): 搜索结果列表

### type RegistrySearchResult

```go
type RegistrySearchResult struct {
    Path  string      `json:"path"`
    Name  string      `json:"name"`
    Value interface{} `json:"value"`
    Type  string      `json:"type"`
}
```

RegistrySearchResult 表示注册表搜索结果结构。

**字段说明:**

- `Path` (string): 注册表路径
- `Name` (string): 键名称
- `Value` (interface{}): 键值
- `Type` (string): 键类型

### type RegistryHandler

```go
type RegistryHandler struct {
    registryService *services.RegistryService
    logger          *logrus.Logger
}
```

RegistryHandler 表示一个注册表操作处理器，提供注册表键值的管理和查询功能。

### func NewRegistryHandler

```go
func NewRegistryHandler(registryService *services.RegistryService, logger *logrus.Logger) *RegistryHandler
```

NewRegistryHandler 使用指定的注册表服务和日志记录器创建 RegistryHandler 的新实例。

### func (h \*RegistryHandler) GetRegistryKey

```go
func (h *RegistryHandler) GetRegistryKey(c *gin.Context)
```

GetRegistryKey 获取指定注册表路径的键信息，包括子键列表、值列表、最后修改时间等详细信息。

**功能说明:**

1. 从路径参数中获取注册表路径
2. 验证注册表路径的有效性
3. 获取注册表键的详细信息
4. 返回键信息、子键列表和值列表

### func (h \*RegistryHandler) CreateRegistryKey

```go
func (h *RegistryHandler) CreateRegistryKey(c *gin.Context)
```

CreateRegistryKey 在指定路径下创建新的注册表键，支持创建多级键结构。

**功能说明:**

1. 从请求体中解析注册表路径和描述信息
2. 验证注册表路径的有效性
3. 创建新的注册表键
4. 返回创建结果

### func (h \*RegistryHandler) DeleteRegistryKey

```go
func (h *RegistryHandler) DeleteRegistryKey(c *gin.Context)
```

DeleteRegistryKey 删除指定路径的注册表键及其所有子键和值。

**功能说明:**

1. 从路径参数中获取注册表路径
2. 验证注册表路径的有效性
3. 删除注册表键及其所有子键和值
4. 返回删除结果

### func (h \*RegistryHandler) SetRegistryValue

```go
func (h *RegistryHandler) SetRegistryValue(c *gin.Context)
```

SetRegistryValue 在指定注册表键下设置值，支持多种数据类型（字符串、DWORD、二进制等）。

**功能说明:**

1. 从请求体中解析注册表路径、值名称、类型和数据
2. 验证注册表路径和值名称的有效性
3. 设置注册表值
4. 返回设置结果

### func (h \*RegistryHandler) GetRegistryValue

```go
func (h *RegistryHandler) GetRegistryValue(c *gin.Context)
```

GetRegistryValue 获取指定注册表键下特定值的详细信息，包括值类型、数据内容等。

**功能说明:**

1. 从路径参数中获取注册表路径和值名称
2. 验证注册表路径和值名称的有效性
3. 获取注册表值的详细信息
4. 返回值类型、数据内容和大小

### func (h \*RegistryHandler) DeleteRegistryValue

```go
func (h *RegistryHandler) DeleteRegistryValue(c *gin.Context)
```

DeleteRegistryValue 删除指定注册表键下的特定值，保留键结构。

**功能说明:**

1. 从路径参数中获取注册表路径和值名称
2. 验证注册表路径和值名称的有效性
3. 删除指定的注册表值
4. 返回删除结果

### func (h \*RegistryHandler) ListRegistryKeys

```go
func (h *RegistryHandler) ListRegistryKeys(c *gin.Context)
```

ListRegistryKeys 获取指定注册表路径下的所有子键名称列表。

**功能说明:**

1. 从路径参数中获取注册表路径
2. 验证注册表路径的有效性
3. 获取所有子键名称列表
4. 返回子键列表和数量

### func (h \*RegistryHandler) ListRegistryValues

```go
func (h *RegistryHandler) ListRegistryValues(c *gin.Context)
```

ListRegistryValues 获取指定注册表路径下的所有值，包括值名称、类型、数据等。

**功能说明:**

1. 从路径参数中获取注册表路径
2. 验证注册表路径的有效性
3. 获取所有值的详细信息
4. 返回值列表和数量

### func (h \*RegistryHandler) SearchRegistry

```go
func (h *RegistryHandler) SearchRegistry(c *gin.Context)
```

SearchRegistry 在指定根路径下搜索包含特定关键词的注册表键或值，支持模糊搜索。

**功能说明:**

1. 从请求体中解析搜索参数（根路径、关键词、搜索选项）
2. 验证搜索参数的有效性
3. 执行注册表搜索
4. 返回搜索结果和匹配数量

---

## 🌐 网络管理模块

### type NetworkConnection

```go
type NetworkConnection struct {
    ID           string    `json:"id"`
    PID          int32     `json:"pid"`
    ProcessName  string    `json:"process_name"`
    LocalAddr    string    `json:"local_addr"`
    LocalPort    int       `json:"local_port"`
    RemoteAddr   string    `json:"remote_addr"`
    RemotePort   int       `json:"remote_port"`
    Status       string    `json:"status"`
    Type         string    `json:"type"`
    CreateTime   time.Time `json:"create_time"`
    LastActivity time.Time `json:"last_activity"`
}
```

NetworkConnection 表示网络连接信息结构。

**字段说明:**

- `ID` (string): 连接唯一标识符
- `PID` (int32): 进程 ID
- `ProcessName` (string): 进程名称
- `LocalAddr` (string): 本地地址
- `LocalPort` (int): 本地端口
- `RemoteAddr` (string): 远程地址
- `RemotePort` (int): 远程端口
- `Status` (string): 连接状态
- `Type` (string): 连接类型
- `CreateTime` (time.Time): 创建时间
- `LastActivity` (time.Time): 最后活动时间

### type NetworkPacket

```go
type NetworkPacket struct {
    ID           string    `json:"id"`
    Timestamp    time.Time `json:"timestamp"`
    SourceIP     string    `json:"source_ip"`
    SourcePort   int       `json:"source_port"`
    DestIP       string    `json:"dest_ip"`
    DestPort     int       `json:"dest_port"`
    Protocol     string    `json:"protocol"`
    Size         int       `json:"size"`
    Payload      []byte    `json:"payload"`
    IsSuspicious bool      `json:"is_suspicious"`
}
```

NetworkPacket 表示网络数据包结构。

**字段说明:**

- `ID` (string): 数据包唯一标识符
- `Timestamp` (time.Time): 时间戳
- `SourceIP` (string): 源 IP 地址
- `SourcePort` (int): 源端口
- `DestIP` (string): 目标 IP 地址
- `DestPort` (int): 目标端口
- `Protocol` (string): 协议类型
- `Size` (int): 数据包大小
- `Payload` ([]byte): 数据载荷
- `IsSuspicious` (bool): 是否可疑

### type NetworkInterface

```go
type NetworkInterface struct {
    Name         string   `json:"name"`
    Index        int      `json:"index"`
    MTU          int      `json:"mtu"`
    Flags        uint32   `json:"flags"`
    HardwareAddr string   `json:"hardware_addr"`
    Addresses    []string `json:"addresses"`
    IsUp         bool     `json:"is_up"`
    IsLoopback   bool     `json:"is_loopback"`
}
```

NetworkInterface 表示网络接口信息结构。

**字段说明:**

- `Name` (string): 接口名称
- `Index` (int): 接口索引
- `MTU` (int): 最大传输单元
- `Flags` (uint32): 接口标志
- `HardwareAddr` (string): 硬件地址
- `Addresses` ([]string): IP 地址列表
- `IsUp` (bool): 是否启用
- `IsLoopback` (bool): 是否回环接口

### type NetworkThreat

```go
type NetworkThreat struct {
    ID           string    `json:"id"`
    Timestamp    time.Time `json:"timestamp"`
    ThreatType   string    `json:"threat_type"`
    SourceIP     string    `json:"source_ip"`
    DestIP       string    `json:"dest_ip"`
    Port         int       `json:"port"`
    Protocol     string    `json:"protocol"`
    Description  string    `json:"description"`
    Severity     string    `json:"severity"`
    IsBlocked    bool      `json:"is_blocked"`
}
```

NetworkThreat 表示网络威胁信息结构。

**字段说明:**

- `ID` (string): 威胁唯一标识符
- `Timestamp` (time.Time): 检测时间
- `ThreatType` (string): 威胁类型
- `SourceIP` (string): 源 IP 地址
- `DestIP` (string): 目标 IP 地址
- `Port` (int): 端口号
- `Protocol` (string): 协议类型
- `Description` (string): 威胁描述
- `Severity` (string): 严重程度
- `IsBlocked` (bool): 是否已阻止

### type NetworkHandler

```go
type NetworkHandler struct {
    networkService *services.NetworkService
    logger         *logrus.Logger
}
```

NetworkHandler 表示一个网络管理处理器，提供网络连接监控和管理功能。

**字段说明:**

- `networkService` (\*services.NetworkService): 网络服务实例
- `logger` (\*logrus.Logger): 日志记录器

### func NewNetworkHandler

```go
func NewNetworkHandler(networkService *services.NetworkService, logger *logrus.Logger) *NetworkHandler
```

NewNetworkHandler 使用指定的网络服务和日志记录器创建 NetworkHandler 的新实例。

**参数:**

- `networkService` (\*services.NetworkService): 网络服务实例
- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `*NetworkHandler`: 新创建的网络处理器实例

### func (h \*NetworkHandler) GetConnections

```go
func (h *NetworkHandler) GetConnections(c *gin.Context)
```

GetConnections 获取系统中所有活跃的网络连接信息，包括 TCP 和 UDP 连接。

**功能说明:**

1. 调用网络服务获取所有连接
2. 过滤活跃连接
3. 返回连接列表和统计信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含网络连接列表的 JSON 响应

### func (h \*NetworkHandler) GetTCPConnections

```go
func (h *NetworkHandler) GetTCPConnections(c *gin.Context)
```

GetTCPConnections 获取系统中所有 TCP 连接信息。

**功能说明:**

1. 调用网络服务获取 TCP 连接
2. 过滤 TCP 连接
3. 返回 TCP 连接列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含 TCP 连接列表的 JSON 响应

### func (h \*NetworkHandler) GetUDPConnections

```go
func (h *NetworkHandler) GetUDPConnections(c *gin.Context)
```

GetUDPConnections 获取系统中所有 UDP 连接信息。

**功能说明:**

1. 调用网络服务获取 UDP 连接
2. 过滤 UDP 连接
3. 返回 UDP 连接列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含 UDP 连接列表的 JSON 响应

### func (h \*NetworkHandler) GetConnectionsByPID

```go
func (h *NetworkHandler) GetConnectionsByPID(c *gin.Context)
```

GetConnectionsByPID 获取指定进程 ID 的所有网络连接。

**功能说明:**

1. 从路径参数中获取 PID
2. 验证 PID 的有效性
3. 获取该进程的所有网络连接
4. 返回连接列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `pid`: 进程 ID (整数)

**返回值:**

- 返回包含指定进程网络连接的 JSON 响应

### func (h \*NetworkHandler) CloseConnection

```go
func (h *NetworkHandler) CloseConnection(c *gin.Context)
```

CloseConnection 强制关闭指定的网络连接。

**功能说明:**

1. 从路径参数中获取连接 ID
2. 验证连接 ID 的有效性
3. 强制关闭网络连接
4. 返回关闭结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `connection_id`: 连接 ID (字符串)

**返回值:**

- 返回包含关闭结果的 JSON 响应

### func (h \*NetworkHandler) GetNetworkInterfaces

```go
func (h *NetworkHandler) GetNetworkInterfaces(c *gin.Context)
```

GetNetworkInterfaces 获取系统中所有网络接口的详细信息。

**功能说明:**

1. 调用网络服务获取网络接口列表
2. 获取每个接口的详细信息
3. 返回接口列表和统计信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含网络接口列表的 JSON 响应

### func (h \*NetworkHandler) GetNetworkStats

```go
func (h *NetworkHandler) GetNetworkStats(c *gin.Context)
```

GetNetworkStats 获取网络统计信息，包括流量统计、连接统计等。

**功能说明:**

1. 调用网络服务获取统计信息
2. 计算网络流量统计
3. 统计连接数量
4. 返回网络统计信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含网络统计信息的 JSON 响应

### func (h \*NetworkHandler) GetConnectionsByPort

```go
func (h *NetworkHandler) GetConnectionsByPort(c *gin.Context)
```

GetConnectionsByPort 获取指定端口的所有网络连接。

**功能说明:**

1. 从路径参数中获取端口号
2. 验证端口号的有效性
3. 获取该端口的所有连接
4. 返回连接列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `port`: 端口号 (整数)

**返回值:**

- 返回包含指定端口连接的 JSON 响应

### func (h \*NetworkHandler) GetConnectionsByIP

```go
func (h *NetworkHandler) GetConnectionsByIP(c *gin.Context)
```

GetConnectionsByIP 获取指定 IP 地址的所有网络连接。

**功能说明:**

1. 从路径参数中获取 IP 地址
2. 验证 IP 地址的有效性
3. 获取该 IP 的所有连接
4. 返回连接列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `ip`: IP 地址 (字符串)

**返回值:**

- 返回包含指定 IP 连接的 JSON 响应

### func (h \*NetworkHandler) IsPortInUse

```go
func (h *NetworkHandler) IsPortInUse(c *gin.Context)
```

IsPortInUse 检查指定端口是否被占用。

**功能说明:**

1. 从路径参数中获取端口号
2. 验证端口号的有效性
3. 检查端口占用状态
4. 返回端口状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `port`: 端口号 (整数)

**返回值:**

- 返回包含端口占用状态的 JSON 响应

### func (h \*NetworkHandler) GetListeningPorts

```go
func (h *NetworkHandler) GetListeningPorts(c *gin.Context)
```

GetListeningPorts 获取系统中所有监听端口的信息。

**功能说明:**

1. 调用网络服务获取监听端口
2. 获取每个端口的详细信息
3. 返回监听端口列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含监听端口列表的 JSON 响应

### func (h \*NetworkHandler) GetEstablishedConnections

```go
func (h *NetworkHandler) GetEstablishedConnections(c *gin.Context)
```

GetEstablishedConnections 获取系统中所有已建立的网络连接。

**功能说明:**

1. 调用网络服务获取已建立连接
2. 过滤已建立状态的连接
3. 返回已建立连接列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含已建立连接列表的 JSON 响应

---

## 🌐 网络监控扩展功能

### func EnableNetworkMonitoring

```go
func EnableNetworkMonitoring(c *gin.Context)
```

EnableNetworkMonitoring 启用网络监控功能，开始实时监控网络连接。

**功能说明:**

1. 检查当前监控状态
2. 启动网络监控服务
3. 记录监控开始时间
4. 返回监控启用结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含监控状态的 JSON 响应，包括状态和开始时间

### func DisableNetworkMonitoring

```go
func DisableNetworkMonitoring(c *gin.Context)
```

DisableNetworkMonitoring 禁用网络监控功能，停止实时监控。

**功能说明:**

1. 检查当前监控状态
2. 停止网络监控服务
3. 记录监控停止时间
4. 返回监控禁用结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含监控状态的 JSON 响应，包括状态和停止时间

### func GetMonitoredConnections

```go
func GetMonitoredConnections(c *gin.Context)
```

GetMonitoredConnections 获取当前被监控的网络连接列表。

**功能说明:**

1. 检查监控服务状态
2. 获取当前被监控的连接列表
3. 返回监控连接的详细信息
4. 统计监控连接总数

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含监控连接列表的 JSON 响应，包括连接信息和总数

### func GetConnectionHistory

```go
func GetConnectionHistory(c *gin.Context)
```

GetConnectionHistory 获取网络连接历史记录，包括已关闭的连接。

**功能说明:**

1. 从查询参数中获取时间范围和限制
2. 获取指定时间范围内的连接历史
3. 返回连接历史记录列表
4. 统计历史连接总数

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**查询参数:**

- `start_date`: 开始日期 (可选)
- `end_date`: 结束日期 (可选)
- `limit`: 限制数量 (可选)
- `offset`: 偏移量 (可选)

**返回值:**

- 返回包含连接历史记录的 JSON 响应，包括记录列表和总数

---

## 🌐 网络管理器扩展功能

### func (n \*network.Manager) GetEstablishedConnections

```go
func (n *network.Manager) GetEstablishedConnections() ([]*models.NetworkConnection, error)
```

GetEstablishedConnections 获取所有已建立的网络连接。

**返回值:**

- `[]*models.NetworkConnection`: 已建立连接列表
- `error`: 错误信息

### func (n \*network.Manager) GetNetworkInterfaces

```go
func (n *network.Manager) GetNetworkInterfaces() ([]*models.NetworkInterface, error)
```

GetNetworkInterfaces 获取所有网络接口信息。

**返回值:**

- `[]*models.NetworkInterface`: 网络接口列表
- `error`: 错误信息

---

## 🌐 网络服务扩展功能

### func (n \*network.Service) EnableMonitoring

```go
func (n *network.Service) EnableMonitoring() error
```

EnableMonitoring 启用网络监控服务。

**返回值:**

- `error`: 错误信息

### func (n \*network.Service) DisableMonitoring

```go
func (n *network.Service) DisableMonitoring() error
```

DisableMonitoring 禁用网络监控服务。

**返回值:**

- `error`: 错误信息

### func (n \*network.Service) GetMonitoredConnections

```go
func (n *network.Service) GetMonitoredConnections() ([]*models.MonitoredNetworkConnection, error)
```

GetMonitoredConnections 获取被监控的网络连接列表。

**返回值:**

- `[]*models.MonitoredNetworkConnection`: 被监控连接列表
- `error`: 错误信息

### func (n \*network.Service) GetConnectionHistory

```go
func (n *network.Service) GetConnectionHistory(limit, offset int, startDate, endDate string) ([]*models.NetworkConnectionHistory, error)
```

GetConnectionHistory 获取网络连接历史记录。

**参数:**

- `limit` (int): 限制数量
- `offset` (int): 偏移量
- `startDate` (string): 开始日期
- `endDate` (string): 结束日期

**返回值:**

- `[]*models.NetworkConnectionHistory`: 连接历史记录列表
- `error`: 错误信息

---

## 📊 网络管理响应类型定义

### type NetworkConnectionsResponse

```go
type NetworkConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

NetworkConnectionsResponse 表示网络连接列表响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): 网络连接列表

### type TCPConnectionsResponse

```go
type TCPConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

TCPConnectionsResponse 表示TCP连接列表响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): TCP连接列表

### type UDPConnectionsResponse

```go
type UDPConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

UDPConnectionsResponse 表示UDP连接列表响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): UDP连接列表

### type ProcessConnectionsResponse

```go
type ProcessConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

ProcessConnectionsResponse 表示进程网络连接响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): 网络连接列表

### type PortConnectionsResponse

```go
type PortConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

PortConnectionsResponse 表示端口连接响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): 网络连接列表

### type IPConnectionsResponse

```go
type IPConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

IPConnectionsResponse 表示IP连接响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): 网络连接列表

### type ConnectionDetailResponse

```go
type ConnectionDetailResponse struct {
    Code    int             `json:"code"`
    Message string          `json:"message"`
    Data    ConnectionDetail `json:"data"`
}
```

ConnectionDetailResponse 表示连接详情响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (ConnectionDetail): 连接详情数据

### type ConnectionDetail

```go
type ConnectionDetail struct {
    LocalAddress  string `json:"local_address"`
    LocalPort     int    `json:"local_port"`
    RemoteAddress string `json:"remote_address"`
    RemotePort    int    `json:"remote_port"`
    Protocol      string `json:"protocol"`
    State         string `json:"state"`
    PID           int    `json:"pid"`
    ProcessName   string `json:"process_name"`
    CreateTime    string `json:"create_time"`
    LastActivity  string `json:"last_activity"`
}
```

ConnectionDetail 表示连接详细信息结构。

**字段说明:**

- `LocalAddress` (string): 本地地址
- `LocalPort` (int): 本地端口
- `RemoteAddress` (string): 远程地址
- `RemotePort` (int): 远程端口
- `Protocol` (string): 协议类型
- `State` (string): 连接状态
- `PID` (int): 进程ID
- `ProcessName` (string): 进程名称
- `CreateTime` (string): 创建时间
- `LastActivity` (string): 最后活动时间

### type CloseConnectionRequest

```go
type CloseConnectionRequest struct {
    LocalAddress  string `json:"local_address" binding:"required"`
    LocalPort     int    `json:"local_port" binding:"required"`
    RemoteAddress string `json:"remote_address" binding:"required"`
    RemotePort    int    `json:"remote_port" binding:"required"`
    Protocol      string `json:"protocol" binding:"required"`
}
```

CloseConnectionRequest 表示关闭连接请求结构。

**字段说明:**

- `LocalAddress` (string): 本地地址，必需字段
- `LocalPort` (int): 本地端口，必需字段
- `RemoteAddress` (string): 远程地址，必需字段
- `RemotePort` (int): 远程端口，必需字段
- `Protocol` (string): 协议类型，必需字段

### type CloseConnectionResponse

```go
type CloseConnectionResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

CloseConnectionResponse 表示关闭连接响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type PortsInUseResponse

```go
type PortsInUseResponse struct {
    Code    int        `json:"code"`
    Message string     `json:"message"`
    Data    []PortInfo `json:"data"`
}
```

PortsInUseResponse 表示正在使用的端口响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]PortInfo): 端口信息列表

### type PortInfo

```go
type PortInfo struct {
    Port        int    `json:"port"`
    Protocol    string `json:"protocol"`
    State       string `json:"state"`
    PID         int    `json:"pid"`
    ProcessName string `json:"process_name"`
}
```

PortInfo 表示端口信息结构。

**字段说明:**

- `Port` (int): 端口号
- `Protocol` (string): 协议类型
- `State` (string): 端口状态
- `PID` (int): 进程ID
- `ProcessName` (string): 进程名称

### type PortStatusResponse

```go
type PortStatusResponse struct {
    Code    int        `json:"code"`
    Message string     `json:"message"`
    Data    PortStatus `json:"data"`
}
```

PortStatusResponse 表示端口状态响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (PortStatus): 端口状态数据

### type PortStatus

```go
type PortStatus struct {
    Port        int    `json:"port"`
    Protocol    string `json:"protocol"`
    State       string `json:"state"`
    PID         int    `json:"pid"`
    ProcessName string `json:"process_name"`
    IsInUse     bool   `json:"is_in_use"`
}
```

PortStatus 表示端口状态结构。

**字段说明:**

- `Port` (int): 端口号
- `Protocol` (string): 协议类型
- `State` (string): 端口状态
- `PID` (int): 进程ID
- `ProcessName` (string): 进程名称
- `IsInUse` (bool): 是否正在使用

### type ListeningPortsResponse

```go
type ListeningPortsResponse struct {
    Code    int             `json:"code"`
    Message string          `json:"message"`
    Data    []ListeningPort `json:"data"`
}
```

ListeningPortsResponse 表示监听端口响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]ListeningPort): 监听端口列表

### type ListeningPort

```go
type ListeningPort struct {
    Port        int    `json:"port"`
    Protocol    string `json:"protocol"`
    PID         int    `json:"pid"`
    ProcessName string `json:"process_name"`
}
```

ListeningPort 表示监听端口结构。

**字段说明:**

- `Port` (int): 端口号
- `Protocol` (string): 协议类型
- `PID` (int): 进程ID
- `ProcessName` (string): 进程名称

### type EstablishedConnectionsResponse

```go
type EstablishedConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

EstablishedConnectionsResponse 表示已建立连接响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): 网络连接列表

### type NetworkInterfacesResponse

```go
type NetworkInterfacesResponse struct {
    Code    int                `json:"code"`
    Message string             `json:"message"`
    Data    []NetworkInterface `json:"data"`
}
```

NetworkInterfacesResponse 表示网络接口响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkInterface): 网络接口列表

### type NetworkStatsResponse

```go
type NetworkStatsResponse struct {
    Code    int          `json:"code"`
    Message string       `json:"message"`
    Data    NetworkStats `json:"data"`
}
```

NetworkStatsResponse 表示网络统计信息响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (NetworkStats): 网络统计信息数据

### type EnableMonitoringResponse

```go
type EnableMonitoringResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

EnableMonitoringResponse 表示启用监控响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type DisableMonitoringResponse

```go
type DisableMonitoringResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

DisableMonitoringResponse 表示禁用监控响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type MonitoringConnectionsResponse

```go
type MonitoringConnectionsResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []NetworkConnection `json:"data"`
}
```

MonitoringConnectionsResponse 表示监控连接响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]NetworkConnection): 网络连接列表

### type ConnectionHistoryResponse

```go
type ConnectionHistoryResponse struct {
    Code    int                 `json:"code"`
    Message string              `json:"message"`
    Data    []ConnectionHistory `json:"data"`
}
```

ConnectionHistoryResponse 表示连接历史响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]ConnectionHistory): 连接历史列表

### type ConnectionHistory

```go
type ConnectionHistory struct {
    LocalAddress  string `json:"local_address"`
    LocalPort     int    `json:"local_port"`
    RemoteAddress string `json:"remote_address"`
    RemotePort    int    `json:"remote_port"`
    Protocol      string `json:"protocol"`
    PID           int    `json:"pid"`
    ProcessName   string `json:"process_name"`
    ConnectTime   string `json:"connect_time"`
    DisconnectTime string `json:"disconnect_time"`
}
```

ConnectionHistory 表示连接历史结构。

**字段说明:**

- `LocalAddress` (string): 本地地址
- `LocalPort` (int): 本地端口
- `RemoteAddress` (string): 远程地址
- `RemotePort` (int): 远程端口
- `Protocol` (string): 协议类型
- `PID` (int): 进程ID
- `ProcessName` (string): 进程名称
- `ConnectTime` (string): 连接时间
- `DisconnectTime` (string): 断开时间

---

## 👤 用户管理模块

### type UserHandler

```go
type UserHandler struct {
    manager *user.Manager
    logger  *logrus.Logger
}
```

UserHandler 表示一个用户管理处理器，提供用户信息查询和权限管理功能。

**字段说明:**

- `manager` (\*user.Manager): 用户管理器实例
- `logger` (\*logrus.Logger): 日志记录器

### func NewUserHandler

```go
func NewUserHandler(manager *user.Manager, logger *logrus.Logger) *UserHandler
```

NewUserHandler 使用指定的用户管理器和日志记录器创建 UserHandler 的新实例。

**参数:**

- `manager` (\*user.Manager): 用户管理器实例
- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `*UserHandler`: 新创建的用户处理器实例

### func (h \*UserHandler) GetCurrentUser

```go
func (h *UserHandler) GetCurrentUser(c *gin.Context)
```

GetCurrentUser 获取当前登录用户的详细信息。

**功能说明:**

1. 从认证上下文获取当前用户信息
2. 获取用户详细信息
3. 返回用户信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含当前用户信息的 JSON 响应

### func (h \*UserHandler) GetUserByID

```go
func (h *UserHandler) GetUserByID(c *gin.Context)
```

GetUserByID 通过用户 ID 获取用户详细信息。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 获取用户详细信息
4. 返回用户信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**返回值:**

- 返回包含用户详细信息的 JSON 响应

### func (h \*UserHandler) GetUserByName

```go
func (h *UserHandler) GetUserByName(c *gin.Context)
```

GetUserByName 通过用户名获取用户详细信息。

**功能说明:**

1. 从路径参数中获取用户名
2. 验证用户名的有效性
3. 获取用户详细信息
4. 返回用户信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `username`: 用户名 (字符串)

**返回值:**

- 返回包含用户详细信息的 JSON 响应

### func (h \*UserHandler) GetAllUsers

```go
func (h *UserHandler) GetAllUsers(c *gin.Context)
```

GetAllUsers 获取系统中所有用户的列表信息。

**功能说明:**

1. 调用用户管理器获取所有用户
2. 过滤敏感信息
3. 返回用户列表和总数

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含所有用户列表的 JSON 响应

### func (h \*UserHandler) CheckUserPermissions

```go
func (h *UserHandler) CheckUserPermissions(c *gin.Context)
```

CheckUserPermissions 检查指定用户的权限信息。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 获取用户权限信息
4. 返回权限列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**返回值:**

- 返回包含用户权限信息的 JSON 响应

### func (h \*UserHandler) ValidatePassword

```go
func (h *UserHandler) ValidatePassword(c *gin.Context)
```

ValidatePassword 验证用户密码的有效性。

**功能说明:**

1. 从请求体中解析用户 ID 和密码
2. 验证用户 ID 的有效性
3. 验证密码格式和强度
4. 返回验证结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**请求体:**

- `user_id`: 用户 ID (整数)
- `password`: 密码 (字符串)

**返回值:**

- 返回包含密码验证结果的 JSON 响应

### func (h \*UserHandler) GetPasswordPolicy

```go
func (h *UserHandler) GetPasswordPolicy(c *gin.Context)
```

GetPasswordPolicy 获取系统密码策略配置。

**功能说明:**

1. 获取系统密码策略配置
2. 返回密码复杂度要求
3. 返回密码有效期设置
4. 返回密码历史限制

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**返回值:**

- 返回包含密码策略的 JSON 响应

### func (h \*UserHandler) CheckAccountStatus

```go
func (h *UserHandler) CheckAccountStatus(c *gin.Context)
```

CheckAccountStatus 检查用户账户状态。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 检查账户锁定状态
4. 检查账户过期状态
5. 返回账户状态信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**返回值:**

- 返回包含账户状态的 JSON 响应

### func (h \*UserHandler) GetUserSessions

```go
func (h *UserHandler) GetUserSessions(c *gin.Context)
```

GetUserSessions 获取用户的所有活跃会话信息。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 获取用户活跃会话列表
4. 返回会话信息

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**返回值:**

- 返回包含用户会话列表的 JSON 响应

### func (h \*UserHandler) KillUserSession

```go
func (h *UserHandler) KillUserSession(c *gin.Context)
```

KillUserSession 强制终止指定的用户会话。

**功能说明:**

1. 从路径参数中获取会话 ID
2. 验证会话 ID 的有效性
3. 强制终止用户会话
4. 返回终止结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `session_id`: 会话 ID (字符串)

**返回值:**

- 返回包含会话终止结果的 JSON 响应

### func (h \*UserHandler) LockUserAccount

```go
func (h *UserHandler) LockUserAccount(c *gin.Context)
```

LockUserAccount 锁定指定的用户账户。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 锁定用户账户
4. 记录锁定原因和时间
5. 返回锁定结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**请求体:**

- `reason`: 锁定原因 (字符串)

**返回值:**

- 返回包含账户锁定结果的 JSON 响应

### func (h \*UserHandler) UnlockUserAccount

```go
func (h *UserHandler) UnlockUserAccount(c *gin.Context)
```

UnlockUserAccount 解锁指定的用户账户。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 解锁用户账户
4. 记录解锁时间和操作者
5. 返回解锁结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**返回值:**

- 返回包含账户解锁结果的 JSON 响应

### func (h \*UserHandler) ChangeUserPassword

```go
func (h *UserHandler) ChangeUserPassword(c *gin.Context)
```

ChangeUserPassword 更改指定用户的密码。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 验证新密码的复杂度
4. 更新用户密码
5. 记录密码更改时间
6. 返回更改结果

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**请求体:**

- `new_password`: 新密码 (字符串)
- `confirm_password`: 确认密码 (字符串)

**返回值:**

- 返回包含密码更改结果的 JSON 响应

### func (h \*UserHandler) GetUserGroups

```go
func (h *UserHandler) GetUserGroups(c *gin.Context)
```

GetUserGroups 获取指定用户所属的用户组列表。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 获取用户所属的用户组
4. 返回用户组列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**返回值:**

- 返回包含用户组列表的 JSON 响应

### func (h \*UserHandler) GetUserLoginHistory

```go
func (h *UserHandler) GetUserLoginHistory(c *gin.Context)
```

GetUserLoginHistory 获取指定用户的登录历史记录。

**功能说明:**

1. 从路径参数中获取用户 ID
2. 验证用户 ID 的有效性
3. 获取用户登录历史记录
4. 返回登录历史列表

**参数:**

- `c *gin.Context`: 包含请求的 Gin HTTP 上下文

**路径参数:**

- `user_id`: 用户 ID (整数)

**查询参数:**

- `start_date`: 开始日期 (可选)
- `end_date`: 结束日期 (可选)
- `limit`: 限制数量 (可选)
- `offset`: 偏移量 (可选)

**返回值:**

- 返回包含用户登录历史的 JSON 响应

---

## 📊 用户管理响应类型定义

### type CurrentUserResponse

```go
type CurrentUserResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    CurrentUser `json:"data"`
}
```

CurrentUserResponse 表示当前用户响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (CurrentUser): 当前用户数据

### type CurrentUser

```go
type CurrentUser struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Status   string `json:"status"`
}
```

CurrentUser 表示当前用户结构。

**字段说明:**

- `ID` (int): 用户ID
- `Username` (string): 用户名
- `Email` (string): 邮箱
- `Role` (string): 角色
- `Status` (string): 状态

### type AllUsersResponse

```go
type AllUsersResponse struct {
    Code    int        `json:"code"`
    Message string     `json:"message"`
    Data    []UserInfo `json:"data"`
}
```

AllUsersResponse 表示所有用户响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]UserInfo): 用户信息列表

### type UserInfo

```go
type UserInfo struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Status   string `json:"status"`
}
```

UserInfo 表示用户信息结构。

**字段说明:**

- `ID` (int): 用户ID
- `Username` (string): 用户名
- `Email` (string): 邮箱
- `Role` (string): 角色
- `Status` (string): 状态

### type UserByIdResponse

```go
type UserByIdResponse struct {
    Code    int        `json:"code"`
    Message string     `json:"message"`
    Data    UserDetail `json:"data"`
}
```

UserByIdResponse 表示根据ID获取用户响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (UserDetail): 用户详细信息

### type UserDetail

```go
type UserDetail struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    Status   string `json:"status"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

UserDetail 表示用户详细信息结构。

**字段说明:**

- `ID` (int): 用户ID
- `Username` (string): 用户名
- `Email` (string): 邮箱
- `Role` (string): 角色
- `Status` (string): 状态
- `CreatedAt` (string): 创建时间
- `UpdatedAt` (string): 更新时间

### type UserByNameResponse

```go
type UserByNameResponse struct {
    Code    int        `json:"code"`
    Message string     `json:"message"`
    Data    UserDetail `json:"data"`
}
```

UserByNameResponse 表示根据用户名获取用户响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (UserDetail): 用户详细信息

### type UserGroupsResponse

```go
type UserGroupsResponse struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    []UserGroup `json:"data"`
}
```

UserGroupsResponse 表示用户组响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]UserGroup): 用户组列表

### type UserGroup

```go
type UserGroup struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
    Members     int    `json:"members"`
}
```

UserGroup 表示用户组结构。

**字段说明:**

- `ID` (int): 组ID
- `Name` (string): 组名
- `Description` (string): 组描述
- `Members` (int): 成员数量

### type CheckPermissionRequest

```go
type CheckPermissionRequest struct {
    UserID       int    `json:"user_id" binding:"required"`
    Permission   string `json:"permission" binding:"required"`
    Resource     string `json:"resource" binding:"required"`
}
```

CheckPermissionRequest 表示检查权限请求结构。

**字段说明:**

- `UserID` (int): 用户ID，必需字段
- `Permission` (string): 权限名称，必需字段
- `Resource` (string): 资源名称，必需字段

### type CheckPermissionResponse

```go
type CheckPermissionResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    bool   `json:"data"`
}
```

CheckPermissionResponse 表示检查权限响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (bool): 是否有权限

### type ValidatePasswordRequest

```go
type ValidatePasswordRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}
```

ValidatePasswordRequest 表示验证密码请求结构。

**字段说明:**

- `Username` (string): 用户名，必需字段
- `Password` (string): 密码，必需字段

### type ValidatePasswordResponse

```go
type ValidatePasswordResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    bool   `json:"data"`
}
```

ValidatePasswordResponse 表示验证密码响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (bool): 密码是否有效

### type PasswordPolicyResponse

```go
type PasswordPolicyResponse struct {
    Code    int            `json:"code"`
    Message string         `json:"message"`
    Data    PasswordPolicy `json:"data"`
}
```

PasswordPolicyResponse 表示密码策略响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (PasswordPolicy): 密码策略数据

### type PasswordPolicy

```go
type PasswordPolicy struct {
    MinLength     int  `json:"min_length"`
    RequireUpper  bool `json:"require_upper"`
    RequireLower  bool `json:"require_lower"`
    RequireDigit  bool `json:"require_digit"`
    RequireSpecial bool `json:"require_special"`
    MaxAge        int  `json:"max_age"`
}
```

PasswordPolicy 表示密码策略结构。

**字段说明:**

- `MinLength` (int): 最小长度
- `RequireUpper` (bool): 是否需要大写字母
- `RequireLower` (bool): 是否需要小写字母
- `RequireDigit` (bool): 是否需要数字
- `RequireSpecial` (bool): 是否需要特殊字符
- `MaxAge` (int): 最大有效期

### type UserStatusResponse

```go
type UserStatusResponse struct {
    Code    int        `json:"code"`
    Message string     `json:"message"`
    Data    UserStatus `json:"data"`
}
```

UserStatusResponse 表示用户状态响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` (UserStatus): 用户状态数据

### type UserStatus

```go
type UserStatus struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Status   string `json:"status"`
    Locked   bool   `json:"locked"`
    LockTime string `json:"lock_time"`
}
```

UserStatus 表示用户状态结构。

**字段说明:**

- `ID` (int): 用户ID
- `Username` (string): 用户名
- `Status` (string): 状态
- `Locked` (bool): 是否锁定
- `LockTime` (string): 锁定时间

### type UserSessionsResponse

```go
type UserSessionsResponse struct {
    Code    int           `json:"code"`
    Message string        `json:"message"`
    Data    []UserSession `json:"data"`
}
```

UserSessionsResponse 表示用户会话响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]UserSession): 用户会话列表

### type UserSession

```go
type UserSession struct {
    ID        string `json:"id"`
    UserID    int    `json:"user_id"`
    Username  string `json:"username"`
    LoginTime string `json:"login_time"`
    LastActivity string `json:"last_activity"`
    IPAddress string `json:"ip_address"`
}
```

UserSession 表示用户会话结构。

**字段说明:**

- `ID` (string): 会话ID
- `UserID` (int): 用户ID
- `Username` (string): 用户名
- `LoginTime` (string): 登录时间
- `LastActivity` (string): 最后活动时间
- `IPAddress` (string): IP地址

### type EndSessionResponse

```go
type EndSessionResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

EndSessionResponse 表示结束会话响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type LockUserResponse

```go
type LockUserResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

LockUserResponse 表示锁定用户响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type UnlockUserResponse

```go
type UnlockUserResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

UnlockUserResponse 表示解锁用户响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type ChangePasswordRequest

```go
type ChangePasswordRequest struct {
    UserID      int    `json:"user_id" binding:"required"`
    OldPassword string `json:"old_password" binding:"required"`
    NewPassword string `json:"new_password" binding:"required"`
}
```

ChangePasswordRequest 表示修改密码请求结构。

**字段说明:**

- `UserID` (int): 用户ID，必需字段
- `OldPassword` (string): 旧密码，必需字段
- `NewPassword` (string): 新密码，必需字段

### type ChangePasswordResponse

```go
type ChangePasswordResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

ChangePasswordResponse 表示修改密码响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息

### type UserHistoryResponse

```go
type UserHistoryResponse struct {
    Code    int            `json:"code"`
    Message string         `json:"message"`
    Data    []LoginHistory `json:"data"`
}
```

UserHistoryResponse 表示用户历史响应结构。

**字段说明:**

- `Code` (int): 响应状态码
- `Message` (string): 响应消息
- `Data` ([]LoginHistory): 登录历史列表

### type LoginHistory

```go
type LoginHistory struct {
    ID        int    `json:"id"`
    UserID    int    `json:"user_id"`
    Username  string `json:"username"`
    LoginTime string `json:"login_time"`
    LogoutTime string `json:"logout_time"`
    IPAddress string `json:"ip_address"`
    Success   bool   `json:"success"`
}
```

LoginHistory 表示登录历史结构。

**字段说明:**

- `ID` (int): 记录ID
- `UserID` (int): 用户ID
- `Username` (string): 用户名
- `LoginTime` (string): 登录时间
- `LogoutTime` (string): 登出时间
- `IPAddress` (string): IP地址
- `Success` (bool): 是否成功

---

## 🔧 中间件模块

### type AuthConfig

```go
type AuthConfig struct {
    SecretKey     string        `json:"secret_key"`
    TokenExpiry   time.Duration `json:"token_expiry"`
    RefreshExpiry time.Duration `json:"refresh_expiry"`
    Issuer        string        `json:"issuer"`
    Audience      string        `json:"audience"`
}
```

AuthConfig 表示认证中间件配置结构。

**字段说明:**

- `SecretKey` (string): JWT 签名密钥
- `TokenExpiry` (time.Duration): 访问令牌过期时间
- `RefreshExpiry` (time.Duration): 刷新令牌过期时间
- `Issuer` (string): 令牌发行者
- `Audience` (string): 令牌受众

### func DefaultAuthConfig

```go
func DefaultAuthConfig() *AuthConfig
```

DefaultAuthConfig 返回默认的认证配置。

**返回值:**

- `*AuthConfig`: 默认认证配置实例

### func Auth

```go
func Auth(config *AuthConfig, logger *logrus.Logger) gin.HandlerFunc
```

Auth 创建认证中间件，用于验证 JWT 令牌和用户身份。

**参数:**

- `config` (\*AuthConfig): 认证配置
- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `gin.HandlerFunc`: 认证中间件函数

---

### type RateLimitConfig

```go
type RateLimitConfig struct {
    Requests    int           `json:"requests"`
    Window      time.Duration `json:"window"`
    Burst       int           `json:"burst"`
    BlockTime   time.Duration `json:"block_time"`
    StorageType string        `json:"storage_type"`
}
```

RateLimitConfig 表示速率限制中间件配置结构。

**字段说明:**

- `Requests` (int): 允许的请求数量
- `Window` (time.Duration): 时间窗口
- `Burst` (int): 突发请求数量
- `BlockTime` (time.Duration): 阻止时间
- `StorageType` (string): 存储类型

### func DefaultRateLimitConfig

```go
func DefaultRateLimitConfig() *RateLimitConfig
```

DefaultRateLimitConfig 返回默认的速率限制配置。

**返回值:**

- `*RateLimitConfig`: 默认速率限制配置实例

### func RateLimit

```go
func RateLimit(config *RateLimitConfig) gin.HandlerFunc
```

RateLimit 创建速率限制中间件，用于控制 API 请求频率。

**参数:**

- `config` (\*RateLimitConfig): 速率限制配置

**返回值:**

- `gin.HandlerFunc`: 速率限制中间件函数

---

### type MetricsCollector

```go
type MetricsCollector struct {
    logger *logrus.Logger
    mu     sync.RWMutex

    // 请求统计
    totalRequests    int64
    successfulRequests int64
    failedRequests   int64

    // 响应时间统计
    responseTimes    []time.Duration
    avgResponseTime  time.Duration

    // 错误统计
    errorCounts      map[string]int64

    // 活跃连接数
    activeConnections int64
}
```

MetricsCollector 表示指标收集器，用于收集和统计系统性能指标。

**字段说明:**

- `logger` (\*logrus.Logger): 日志记录器
- `mu` (sync.RWMutex): 读写互斥锁
- `totalRequests` (int64): 总请求数
- `successfulRequests` (int64): 成功请求数
- `failedRequests` (int64): 失败请求数
- `responseTimes` ([]time.Duration): 响应时间列表
- `avgResponseTime` (time.Duration): 平均响应时间
- `errorCounts` (map[string]int64): 错误计数映射
- `activeConnections` (int64): 活跃连接数

### func NewMetricsCollector

```go
func NewMetricsCollector(logger *logrus.Logger) *MetricsCollector
```

NewMetricsCollector 创建新的指标收集器实例。

**参数:**

- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `*MetricsCollector`: 新的指标收集器实例

### func (mc \*MetricsCollector) Collect

```go
func (mc *MetricsCollector) Collect() gin.HandlerFunc
```

Collect 创建指标收集中间件，用于收集请求和响应指标。

**返回值:**

- `gin.HandlerFunc`: 指标收集中间件函数

### func (mc \*MetricsCollector) GetMetrics

```go
func (mc *MetricsCollector) GetMetrics() *Metrics
```

GetMetrics 获取当前收集的指标数据。

**返回值:**

- `*Metrics`: 指标数据

### func (mc \*MetricsCollector) ResetMetrics

```go
func (mc *MetricsCollector) ResetMetrics()
```

ResetMetrics 重置所有指标数据。

### func (mc \*MetricsCollector) ReportMetrics

```go
func (mc *MetricsCollector) ReportMetrics()
```

ReportMetrics 报告当前指标数据到日志系统。

---

### type TimeoutConfig

```go
type TimeoutConfig struct {
    ReadTimeout  time.Duration `json:"read_timeout"`
    WriteTimeout time.Duration `json:"write_timeout"`
    IdleTimeout  time.Duration `json:"idle_timeout"`
}
```

TimeoutConfig 表示超时中间件配置结构。

**字段说明:**

- `ReadTimeout` (time.Duration): 读取超时时间
- `WriteTimeout` (time.Duration): 写入超时时间
- `IdleTimeout` (time.Duration): 空闲超时时间

### func DefaultTimeoutConfig

```go
func DefaultTimeoutConfig() *TimeoutConfig
```

DefaultTimeoutConfig 返回默认的超时配置。

**返回值:**

- `*TimeoutConfig`: 默认超时配置实例

### func Timeout

```go
func Timeout(config *TimeoutConfig, logger *logrus.Logger) gin.HandlerFunc
```

Timeout 创建超时中间件，用于控制请求处理超时。

**参数:**

- `config` (\*TimeoutConfig): 超时配置
- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `gin.HandlerFunc`: 超时中间件函数

---

### type ErrorHandler

```go
type ErrorHandler struct {
    logger *logrus.Logger
}
```

ErrorHandler 表示错误处理中间件，用于统一处理系统错误。

**字段说明:**

- `logger` (\*logrus.Logger): 日志记录器

### func NewErrorHandler

```go
func NewErrorHandler(logger *logrus.Logger) *ErrorHandler
```

NewErrorHandler 创建新的错误处理器实例。

**参数:**

- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `*ErrorHandler`: 新的错误处理器实例

### func (h \*ErrorHandler) Recovery

```go
func (h *ErrorHandler) Recovery() gin.HandlerFunc
```

Recovery 创建恢复中间件，用于捕获和处理 panic。

**返回值:**

- `gin.HandlerFunc`: 恢复中间件函数

### func (h \*ErrorHandler) ErrorResponse

```go
func (h *ErrorHandler) ErrorResponse() gin.HandlerFunc
```

ErrorResponse 创建错误响应中间件，用于统一错误响应格式。

**返回值:**

- `gin.HandlerFunc`: 错误响应中间件函数

### func (h \*ErrorHandler) Timeout

```go
func (h *ErrorHandler) Timeout(timeout time.Duration) gin.HandlerFunc
```

Timeout 创建超时中间件，用于控制请求超时。

**参数:**

- `timeout` (time.Duration): 超时时间

**返回值:**

- `gin.HandlerFunc`: 超时中间件函数

### func (h \*ErrorHandler) RateLimit

```go
func (h *ErrorHandler) RateLimit(requests int, window time.Duration) gin.HandlerFunc
```

RateLimit 创建速率限制中间件。

**参数:**

- `requests` (int): 允许的请求数量
- `window` (time.Duration): 时间窗口

**返回值:**

- `gin.HandlerFunc`: 速率限制中间件函数

### func (h \*ErrorHandler) Validation

```go
func (h *ErrorHandler) Validation() gin.HandlerFunc
```

Validation 创建验证中间件，用于验证请求参数。

**返回值:**

- `gin.HandlerFunc`: 验证中间件函数

### func (h \*ErrorHandler) Logging

```go
func (h *ErrorHandler) Logging() gin.HandlerFunc
```

Logging 创建日志中间件，用于记录请求和响应日志。

**返回值:**

- `gin.HandlerFunc`: 日志中间件函数

### func (h \*ErrorHandler) Security

```go
func (h *ErrorHandler) Security() gin.HandlerFunc
```

Security 创建安全中间件，用于安全检查。

**返回值:**

- `gin.HandlerFunc`: 安全中间件函数

### func (h \*ErrorHandler) HealthCheck

```go
func (h *ErrorHandler) HealthCheck() gin.HandlerFunc
```

HealthCheck 创建健康检查中间件。

**返回值:**

- `gin.HandlerFunc`: 健康检查中间件函数

### func (h \*ErrorHandler) HandlePanic

```go
func (h *ErrorHandler) HandlePanic(c *gin.Context)
```

HandlePanic 处理 panic 异常。

**参数:**

- `c *gin.Context`: Gin HTTP 上下文

---

### type LoggingConfig

```go
type LoggingConfig struct {
    Level      string `json:"level"`
    Format     string `json:"format"`
    Output     string `json:"output"`
    MaxSize    int    `json:"max_size"`
    MaxBackups int    `json:"max_backups"`
    MaxAge     int    `json:"max_age"`
    Compress   bool   `json:"compress"`
}
```

LoggingConfig 表示日志中间件配置结构。

**字段说明:**

- `Level` (string): 日志级别
- `Format` (string): 日志格式
- `Output` (string): 输出目标
- `MaxSize` (int): 最大文件大小
- `MaxBackups` (int): 最大备份数量
- `MaxAge` (int): 最大保存天数
- `Compress` (bool): 是否压缩

### func DefaultLoggingConfig

```go
func DefaultLoggingConfig() *LoggingConfig
```

DefaultLoggingConfig 返回默认的日志配置。

**返回值:**

- `*LoggingConfig`: 默认日志配置实例

### func Logging

```go
func Logging(config *LoggingConfig, logger *logrus.Logger) gin.HandlerFunc
```

Logging 创建日志中间件，用于记录请求和响应日志。

**参数:**

- `config` (\*LoggingConfig): 日志配置
- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `gin.HandlerFunc`: 日志中间件函数

### func Recovery

```go
func Recovery(logger *logrus.Logger) gin.HandlerFunc
```

Recovery 创建恢复中间件，用于捕获和处理 panic。

**参数:**

- `logger` (\*logrus.Logger): 日志记录器

**返回值:**

- `gin.HandlerFunc`: 恢复中间件函数

### func RequestID

```go
func RequestID() gin.HandlerFunc
```

RequestID 创建请求 ID 中间件，为每个请求生成唯一标识符。

**返回值:**

- `gin.HandlerFunc`: 请求 ID 中间件函数

---

### type CORSConfig

```go
type CORSConfig struct {
    AllowedOrigins   []string `json:"allowed_origins"`
    AllowedMethods   []string `json:"allowed_methods"`
    AllowedHeaders   []string `json:"allowed_headers"`
    ExposedHeaders   []string `json:"exposed_headers"`
    AllowCredentials bool     `json:"allow_credentials"`
    MaxAge           int      `json:"max_age"`
}
```

CORSConfig 表示 CORS 中间件配置结构。

**字段说明:**

- `AllowedOrigins` ([]string): 允许的源
- `AllowedMethods` ([]string): 允许的方法
- `AllowedHeaders` ([]string): 允许的头部
- `ExposedHeaders` ([]string): 暴露的头部
- `AllowCredentials` (bool): 是否允许凭据
- `MaxAge` (int): 预检请求缓存时间

### func DefaultCORSConfig

```go
func DefaultCORSConfig() *CORSConfig
```

DefaultCORSConfig 返回默认的 CORS 配置。

**返回值:**

- `*CORSConfig`: 默认 CORS 配置实例

### func CORS

```go
func CORS(config *CORSConfig) gin.HandlerFunc
```

CORS 创建 CORS 中间件，用于处理跨域请求。

**参数:**

- `config` (\*CORSConfig): CORS 配置

**返回值:**

- `gin.HandlerFunc`: CORS 中间件函数

---

## 🛠️ 工具模块

### type CryptoUtils

```go
type CryptoUtils struct {
    secretKey []byte
    logger    *logrus.Logger
}
```

CryptoUtils 表示加密工具类，提供各种加密、哈希和编码功能。

**字段说明:**

- `secretKey` ([]byte): 加密密钥
- `logger` (\*logrus.Logger): 日志记录器

### func NewCryptoUtils

```go
func NewCryptoUtils() *CryptoUtils
```

NewCryptoUtils 创建新的加密工具实例。

**返回值:**

- `*CryptoUtils`: 新的加密工具实例

### func (c \*CryptoUtils) GenerateRandomBytes

```go
func (c *CryptoUtils) GenerateRandomBytes(length int) ([]byte, error)
```

GenerateRandomBytes 生成指定长度的随机字节数组。

**参数:**

- `length` (int): 字节长度

**返回值:**

- `[]byte`: 随机字节数组
- `error`: 错误信息

### func (c \*CryptoUtils) GenerateRandomString

```go
func (c *CryptoUtils) GenerateRandomString(length int) (string, error)
```

GenerateRandomString 生成指定长度的随机字符串。

**参数:**

- `length` (int): 字符串长度

**返回值:**

- `string`: 随机字符串
- `error`: 错误信息

### func (c \*CryptoUtils) HashSHA256

```go
func (c *CryptoUtils) HashSHA256(data []byte) string
```

HashSHA256 计算数据的 SHA256 哈希值。

**参数:**

- `data` ([]byte): 要哈希的数据

**返回值:**

- `string`: SHA256 哈希值

### func (c \*CryptoUtils) HashSHA512

```go
func (c *CryptoUtils) HashSHA512(data []byte) string
```

HashSHA512 计算数据的 SHA512 哈希值。

**参数:**

- `data` ([]byte): 要哈希的数据

**返回值:**

- `string`: SHA512 哈希值

### func (c \*CryptoUtils) HashWithSalt

```go
func (c *CryptoUtils) HashWithSalt(data []byte, salt []byte) string
```

HashWithSalt 使用盐值计算数据的哈希值。

**参数:**

- `data` ([]byte): 要哈希的数据
- `salt` ([]byte): 盐值

**返回值:**

- `string`: 带盐值的哈希值

### func (c \*CryptoUtils) GenerateHMAC

```go
func (c *CryptoUtils) GenerateHMAC(data []byte, key []byte) string
```

GenerateHMAC 生成数据的 HMAC 值。

**参数:**

- `data` ([]byte): 要计算 HMAC 的数据
- `key` ([]byte): HMAC 密钥

**返回值:**

- `string`: HMAC 值

### func (c \*CryptoUtils) VerifyHMAC

```go
func (c *CryptoUtils) VerifyHMAC(data []byte, key []byte, expectedHMAC string) bool
```

VerifyHMAC 验证数据的 HMAC 值。

**参数:**

- `data` ([]byte): 要验证的数据
- `key` ([]byte): HMAC 密钥
- `expectedHMAC` (string): 期望的 HMAC 值

**返回值:**

- `bool`: 验证结果

### func (c \*CryptoUtils) EncryptAES

```go
func (c *CryptoUtils) EncryptAES(data []byte, key []byte) ([]byte, error)
```

EncryptAES 使用 AES 算法加密数据。

**参数:**

- `data` ([]byte): 要加密的数据
- `key` ([]byte): 加密密钥

**返回值:**

- `[]byte`: 加密后的数据
- `error`: 错误信息

### func (c \*CryptoUtils) DecryptAES

```go
func (c *CryptoUtils) DecryptAES(ciphertext []byte, key []byte) ([]byte, error)
```

DecryptAES 使用 AES 算法解密数据。

**参数:**

- `ciphertext` ([]byte): 要解密的数据
- `key` ([]byte): 解密密钥

**返回值:**

- `[]byte`: 解密后的数据
- `error`: 错误信息

### func (c \*CryptoUtils) GenerateKey

```go
func (c *CryptoUtils) GenerateKey(keySize int) ([]byte, error)
```

GenerateKey 生成指定长度的加密密钥。

**参数:**

- `keySize` (int): 密钥长度（位）

**返回值:**

- `[]byte`: 生成的密钥
- `error`: 错误信息

### func (c \*CryptoUtils) HashFile

```go
func (c *CryptoUtils) HashFile(data []byte, algorithm string) (string, error)
```

HashFile 计算文件的哈希值。

**参数:**

- `data` ([]byte): 文件数据
- `algorithm` (string): 哈希算法

**返回值:**

- `string`: 文件哈希值
- `error`: 错误信息

### func (c \*CryptoUtils) Base64Encode

```go
func (c *CryptoUtils) Base64Encode(data []byte) string
```

Base64Encode 将数据编码为 Base64 格式。

**参数:**

- `data` ([]byte): 要编码的数据

**返回值:**

- `string`: Base64 编码后的字符串

### func (c \*CryptoUtils) Base64Decode

```go
func (c *CryptoUtils) Base64Decode(encoded string) ([]byte, error)
```

Base64Decode 将 Base64 格式的字符串解码为原始数据。

**参数:**

- `encoded` (string): Base64 编码的字符串

**返回值:**

- `[]byte`: 解码后的数据
- `error`: 错误信息

### func (c \*CryptoUtils) HexEncode

```go
func (c *CryptoUtils) HexEncode(data []byte) string
```

HexEncode 将数据编码为十六进制格式。

**参数:**

- `data` ([]byte): 要编码的数据

**返回值:**

- `string`: 十六进制编码后的字符串

### func (c \*CryptoUtils) HexDecode

```go
func (c *CryptoUtils) HexDecode(encoded string) ([]byte, error)
```

HexDecode 将十六进制格式的字符串解码为原始数据。

**参数:**

- `encoded` (string): 十六进制编码的字符串

**返回值:**

- `[]byte`: 解码后的数据
- `error`: 错误信息

### func (c \*CryptoUtils) GeneratePasswordHash

```go
func (c *CryptoUtils) GeneratePasswordHash(password string, salt []byte) string
```

GeneratePasswordHash 生成密码哈希值。

**参数:**

- `password` (string): 密码
- `salt` ([]byte): 盐值

**返回值:**

- `string`: 密码哈希值

### func (c \*CryptoUtils) VerifyPassword

```go
func (c *CryptoUtils) VerifyPassword(password string, salt []byte, expectedHash string) bool
```

VerifyPassword 验证密码是否正确。

**参数:**

- `password` (string): 密码
- `salt` ([]byte): 盐值
- `expectedHash` (string): 期望的哈希值

**返回值:**

- `bool`: 验证结果

### func (c \*CryptoUtils) GenerateAPIKey

```go
func (c *CryptoUtils) GenerateAPIKey(length int) (string, error)
```

GenerateAPIKey 生成 API 密钥。

**参数:**

- `length` (int): 密钥长度

**返回值:**

- `string`: 生成的 API 密钥
- `error`: 错误信息

### func (c \*CryptoUtils) EncryptSensitiveData

```go
func (c *CryptoUtils) EncryptSensitiveData(data string, key []byte) (string, error)
```

EncryptSensitiveData 加密敏感数据。

**参数:**

- `data` (string): 敏感数据
- `key` ([]byte): 加密密钥

**返回值:**

- `string`: 加密后的数据
- `error`: 错误信息

### func (c \*CryptoUtils) DecryptSensitiveData

```go
func (c *CryptoUtils) DecryptSensitiveData(encryptedData string, key []byte) (string, error)
```

DecryptSensitiveData 解密敏感数据。

**参数:**

- `encryptedData` (string): 加密后的数据
- `key` ([]byte): 解密密钥

**返回值:**

- `string`: 解密后的数据
- `error`: 错误信息

---

### type Validator

```go
type Validator struct {
    logger *logrus.Logger
}
```

Validator 表示验证器，提供各种数据验证功能。

**字段说明:**

- `logger` (\*logrus.Logger): 日志记录器

### func NewValidator

```go
func NewValidator() *Validator
```

NewValidator 创建新的验证器实例。

**返回值:**

- `*Validator`: 新的验证器实例

### func (v \*Validator) ValidateFilePath

```go
func (v *Validator) ValidateFilePath(filePath string) error
```

ValidateFilePath 验证文件路径的有效性。

**参数:**

- `filePath` (string): 文件路径

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateDirectoryPath

```go
func (v *Validator) ValidateDirectoryPath(dirPath string) error
```

ValidateDirectoryPath 验证目录路径的有效性。

**参数:**

- `dirPath` (string): 目录路径

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidatePID

```go
func (v *Validator) ValidatePID(pid int32) error
```

ValidatePID 验证进程 ID 的有效性。

**参数:**

- `pid` (int32): 进程 ID

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidatePort

```go
func (v *Validator) ValidatePort(port int) error
```

ValidatePort 验证端口号的有效性。

**参数:**

- `port` (int): 端口号

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateIPAddress

```go
func (v *Validator) ValidateIPAddress(ip string) error
```

ValidateIPAddress 验证 IP 地址的有效性。

**参数:**

- `ip` (string): IP 地址

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateIPRange

```go
func (v *Validator) ValidateIPRange(ipRange string) error
```

ValidateIPRange 验证 IP 地址范围的有效性。

**参数:**

- `ipRange` (string): IP 地址范围

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateEmail

```go
func (v *Validator) ValidateEmail(email string) error
```

ValidateEmail 验证邮箱地址的有效性。

**参数:**

- `email` (string): 邮箱地址

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateURL

```go
func (v *Validator) ValidateURL(url string) error
```

ValidateURL 验证 URL 的有效性。

**参数:**

- `url` (string): URL

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateFileName

```go
func (v *Validator) ValidateFileName(fileName string) error
```

ValidateFileName 验证文件名的有效性。

**参数:**

- `fileName` (string): 文件名

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateFileSize

```go
func (v *Validator) ValidateFileSize(size int64, maxSize int64) error
```

ValidateFileSize 验证文件大小的有效性。

**参数:**

- `size` (int64): 文件大小
- `maxSize` (int64): 最大文件大小

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateStringLength

```go
func (v *Validator) ValidateStringLength(str string, minLength, maxLength int) error
```

ValidateStringLength 验证字符串长度的有效性。

**参数:**

- `str` (string): 字符串
- `minLength` (int): 最小长度
- `maxLength` (int): 最大长度

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateInteger

```go
func (v *Validator) ValidateInteger(value int, min, max int) error
```

ValidateInteger 验证整数值的有效性。

**参数:**

- `value` (int): 整数值
- `min` (int): 最小值
- `max` (int): 最大值

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateFloat

```go
func (v *Validator) ValidateFloat(value float64, min, max float64) error
```

ValidateFloat 验证浮点数值的有效性。

**参数:**

- `value` (float64): 浮点数值
- `min` (float64): 最小值
- `max` (float64): 最大值

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateRegex

```go
func (v *Validator) ValidateRegex(pattern string) error
```

ValidateRegex 验证正则表达式的有效性。

**参数:**

- `pattern` (string): 正则表达式

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateFileExtension

```go
func (v *Validator) ValidateFileExtension(fileName string, allowedExtensions []string) error
```

ValidateFileExtension 验证文件扩展名的有效性。

**参数:**

- `fileName` (string): 文件名
- `allowedExtensions` ([]string): 允许的扩展名列表

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidatePathDepth

```go
func (v *Validator) ValidatePathDepth(path string, maxDepth int) error
```

ValidatePathDepth 验证路径深度的有效性。

**参数:**

- `path` (string): 路径
- `maxDepth` (int): 最大深度

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateRegistryPath

```go
func (v *Validator) ValidateRegistryPath(regPath string) error
```

ValidateRegistryPath 验证注册表路径的有效性。

**参数:**

- `regPath` (string): 注册表路径

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateProcessName

```go
func (v *Validator) ValidateProcessName(processName string) error
```

ValidateProcessName 验证进程名的有效性。

**参数:**

- `processName` (string): 进程名

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateScanDepth

```go
func (v *Validator) ValidateScanDepth(depth int) error
```

ValidateScanDepth 验证扫描深度的有效性。

**参数:**

- `depth` (int): 扫描深度

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateTimeout

```go
func (v *Validator) ValidateTimeout(timeout int) error
```

ValidateTimeout 验证超时时间的有效性。

**参数:**

- `timeout` (int): 超时时间

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateAPIKey

```go
func (v *Validator) ValidateAPIKey(apiKey string) error
```

ValidateAPIKey 验证 API 密钥的有效性。

**参数:**

- `apiKey` (string): API 密钥

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateHash

```go
func (v *Validator) ValidateHash(hash string, algorithm string) error
```

ValidateHash 验证哈希值的有效性。

**参数:**

- `hash` (string): 哈希值
- `algorithm` (string): 哈希算法

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateNumericString

```go
func (v *Validator) ValidateNumericString(str string) error
```

ValidateNumericString 验证字符串是否为数字。

**参数:**

- `str` (string): 字符串

**返回值:**

- `error`: 错误信息

### func (v \*Validator) ValidateHexString

```go
func (v *Validator) ValidateHexString(str string) error
```

ValidateHexString 验证字符串是否为十六进制格式。

**参数:**

- `str` (string): 字符串

**返回值:**

- `error`: 错误信息

---

**文档版本** 3.0  
**最后更新** 2025-08-10  
**维护人员** LYS