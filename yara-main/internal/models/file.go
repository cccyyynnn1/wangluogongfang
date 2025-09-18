package models

import (
	"time"
)

// DeepFileAnalysis 深度文件分析结果
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

// BasicFileInfo 基础文件信息
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

// ContentAnalysis 内容分析
type ContentAnalysis struct {
	Timestamp          time.Time         `json:"timestamp"`
	TextAnalysis       *TextAnalysis     `json:"text_analysis"`
	BinaryAnalysis     *BinaryAnalysis   `json:"binary_analysis"`
	EncodingAnalysis   *EncodingAnalysis `json:"encoding_analysis"`
	SuspiciousPatterns []string          `json:"suspicious_patterns"`
}

// TextAnalysis 文本分析
type TextAnalysis struct {
	Timestamp         time.Time `json:"timestamp"`
	Encoding          string    `json:"encoding"`
	LineCount         int       `json:"line_count"`
	WordCount         int       `json:"word_count"`
	CharacterCount    int       `json:"character_count"`
	Language          string    `json:"language"`
	SpecialCharacters []string  `json:"special_characters"`
}

// BinaryAnalysis 二进制分析
type BinaryAnalysis struct {
	Timestamp   time.Time `json:"timestamp"`
	Entropy     float64   `json:"entropy"`
	HighEntropy bool      `json:"high_entropy"`
	FileHeader  string    `json:"file_header"`
	Strings     []string  `json:"strings"`
	HexPatterns []string  `json:"hex_patterns"`
}

// EncodingAnalysis 编码分析
type EncodingAnalysis struct {
	Timestamp         time.Time `json:"timestamp"`
	EncodingType      string    `json:"encoding_type"`
	IsEncoded         bool      `json:"is_encoded"`
	EncodingAlgorithm string    `json:"encoding_algorithm"`
}

// StructureAnalysis 结构分析
type StructureAnalysis struct {
	Timestamp        time.Time         `json:"timestamp"`
	PEAnalysis       *PEAnalysis       `json:"pe_analysis"`
	ArchiveAnalysis  *ArchiveAnalysis  `json:"archive_analysis"`
	DocumentAnalysis *DocumentAnalysis `json:"document_analysis"`
}

// PEAnalysis PE文件分析
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

// ArchiveAnalysis 压缩文件分析
type ArchiveAnalysis struct {
	Timestamp        time.Time `json:"timestamp"`
	IsValidArchive   bool      `json:"is_valid_archive"`
	ArchiveType      string    `json:"archive_type"`
	FileCount        int       `json:"file_count"`
	CompressedSize   int64     `json:"compressed_size"`
	UncompressedSize int64     `json:"uncompressed_size"`
	Files            []string  `json:"files"`
}

// DocumentAnalysis 文档文件分析
type DocumentAnalysis struct {
	Timestamp       time.Time         `json:"timestamp"`
	IsValidDocument bool              `json:"is_valid_document"`
	DocumentType    string            `json:"document_type"`
	HasMacros       bool              `json:"has_macros"`
	MacroCount      int               `json:"macro_count"`
	EmbeddedObjects []string          `json:"embedded_objects"`
	Metadata        map[string]string `json:"metadata"`
}

// ThreatAnalysis 威胁分析
type ThreatAnalysis struct {
	Timestamp         time.Time          `json:"timestamp"`
	ScanResult        *ScanResult        `json:"scan_result"`
	IsInfected        bool               `json:"is_infected"`
	Threats           []ThreatInfo       `json:"threats"`
	BehaviorAnalysis  *BehaviorAnalysis  `json:"behavior_analysis"`
	SignatureAnalysis *SignatureAnalysis `json:"signature_analysis"`
	EntropyAnalysis   *EntropyAnalysis   `json:"entropy_analysis"`
}

// BehaviorAnalysis 行为分析
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

// SignatureAnalysis 签名分析
type SignatureAnalysis struct {
	Timestamp        time.Time         `json:"timestamp"`
	IsSigned         bool              `json:"is_signed"`
	SignerName       string            `json:"signer_name"`
	CertificateValid bool              `json:"certificate_valid"`
	SignatureValid   bool              `json:"signature_valid"`
	CertificateInfo  map[string]string `json:"certificate_info"`
}

// EntropyAnalysis 熵值分析
type EntropyAnalysis struct {
	Timestamp           time.Time          `json:"timestamp"`
	OverallEntropy      float64            `json:"overall_entropy"`
	SectionEntropies    map[string]float64 `json:"section_entropies"`
	HighEntropySections []string           `json:"high_entropy_sections"`
	EntropyThreshold    float64            `json:"entropy_threshold"`
}

// FileScanResult 文件扫描结果
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
