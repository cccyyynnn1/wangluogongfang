// Global variables
let selectedFile = null;
let lastResults = null;

// DOM elements
const elements = {
    // File upload
    fileUploadArea: document.getElementById('file-upload-area'),
    fileInput: document.getElementById('file-input'),
    fileSelected: document.getElementById('file-selected'),
    fileName: document.getElementById('file-name'),
    removeFile: document.getElementById('remove-file'),
    fileDetectBtn: document.getElementById('file-detect-btn'),
    
    // Image Preview
    imagePreviewSection: document.getElementById('image-preview-section'),
    imagePreview: document.getElementById('image-preview'),

    // Detection Options
    detectionOptions: document.querySelectorAll('.option-item input[type="checkbox"]'),
    
    // Status and Results
    loading: document.getElementById('loading'),
    resultsSection: document.getElementById('results-section'),
    resultsStats: document.getElementById('results-stats'),
    
    // Result Tabs
    tabBtns: document.querySelectorAll('.tab-btn'),
    visualResultContainer: document.getElementById('visual-result-container'),
    resultImage: document.getElementById('result-image'),
    ocrTextContainer: document.getElementById('ocr-text-container'),
    
    // Action Buttons
    exportBtn: document.getElementById('export-btn'),
    clearBtn: document.getElementById('clear-btn'),
    
    // Error Handling
    errorMessage: document.getElementById('error-message'),
    errorText: document.getElementById('error-text'),
    closeError: document.getElementById('close-error')
};

// Initialization
document.addEventListener('DOMContentLoaded', initializeEventListeners);

function initializeEventListeners() {
    // 确保只在图片页面执行初始化
    if (elements.imagePreviewSection) {
        // File upload events
        elements.fileUploadArea.addEventListener('click', () => elements.fileInput.click());
        elements.fileUploadArea.addEventListener('dragover', handleDragOver);
        elements.fileUploadArea.addEventListener('drop', handleFileDrop);
        elements.fileInput.addEventListener('change', handleFileSelect);
        elements.removeFile.addEventListener('click', removeSelectedFile);
        elements.fileDetectBtn.addEventListener('click', detectFromFile);
        
        // Tab switching
        elements.tabBtns.forEach(btn => {
            btn.addEventListener('click', () => switchTab(btn.dataset.tab));
        });

        // Action buttons
        elements.exportBtn.addEventListener('click', exportResults);
        elements.clearBtn.addEventListener('click', clearResults);
        elements.closeError.addEventListener('click', hideError);
    }
}

// --- File Handling ---

function handleDragOver(e) {
    e.preventDefault();
    elements.fileUploadArea.classList.add('dragover');
}

function handleFileDrop(e) {
    e.preventDefault();
    elements.fileUploadArea.classList.remove('dragover');
    const files = e.dataTransfer.files;
    if (files.length > 0) handleFileSelection(files[0]);
}

function handleFileSelect(e) {
    const files = e.target.files;
    if (files.length > 0) handleFileSelection(files[0]);
}

function handleFileSelection(file) {
    const allowedTypes = ['image/png', 'image/jpeg', 'image/webp'];
    if (!allowedTypes.includes(file.type)) {
        showError('只支持 .png, .jpg, 和 .webp 格式的图片');
        return;
    }
    
    if (file.size > 10 * 1024 * 1024) {
        showError('文件大小不能超过 10MB');
        return;
    }
    
    selectedFile = file;
    
    // Update UI
    elements.fileName.textContent = file.name;
    elements.fileUploadArea.querySelector('.upload-content').classList.add('hidden');
    elements.fileSelected.classList.remove('hidden');
    elements.fileDetectBtn.disabled = false;
    
    // Show image preview
    const reader = new FileReader();
    reader.onload = e => {
        elements.imagePreview.src = e.target.result;
        elements.imagePreviewSection.classList.remove('hidden');
    };
    reader.readAsDataURL(file);
    
    hideError();
    clearResults();
}

function removeSelectedFile() {
    selectedFile = null;
    elements.fileInput.value = '';
    elements.fileUploadArea.querySelector('.upload-content').classList.remove('hidden');
    elements.fileSelected.classList.add('hidden');
    elements.fileDetectBtn.disabled = true;
    elements.imagePreviewSection.classList.add('hidden');
    elements.imagePreview.src = '#';
    clearResults();
}

// --- Detection Logic ---

async function detectFromFile() {
    if (!selectedFile) {
        showError('请选择要检测的图片');
        return;
    }
    
    const types = Array.from(elements.detectionOptions)
        .filter(cb => cb.checked)
        .map(cb => cb.value);
        
    if (types.length === 0) {
        showError('请至少选择一种检测类型');
        return;
    }
    
    showLoading();
    hideError();

    // MOCK API CALL: Replace this with your actual API call
    // We simulate a delay and return a mock response
    setTimeout(() => {
        try {
            const mockResponse = generateMockResponse(selectedFile.name);
            if (mockResponse.success) {
                displayResults(mockResponse);
            } else {
                showError(mockResponse.message || '图片检测失败');
            }
        } catch (error) {
            console.error('Detection error:', error);
            showError('检测过程中发生未知错误');
        } finally {
            hideLoading();
        }
    }, 2000); // Simulate 2-second API call
}

// --- MOCK API RESPONSE GENERATOR ---
function generateMockResponse(fileName) {
    // This function creates fake data for demonstration purposes.
    // In a real application, this data would come from your backend server.
    return {
        success: true,
        stats: {
            total_detections: 5,
            processing_time_ms: 1850000,
            image_dimensions: { width: 1200, height: 800 }
        },
        ocr_text: `联系人: 张三\n电话: 13812345678\n地址是北京市朝阳区建国门外大街1号。\n请将款项汇入卡号 6222020100123456789。如有疑问，请发送邮件至 zhang.san@example.com`,
        results: [
            { type: 'PersonName', value: '张三', confidence: 0.95, bounding_box: { x: 100, y: 50, width: 80, height: 30 } },
            { type: 'PhoneNumber', value: '13812345678', confidence: 0.99, bounding_box: { x: 100, y: 120, width: 200, height: 30 } },
            { type: 'Address', value: '北京市朝阳区建国门外大街1号', confidence: 0.85, bounding_box: { x: 80, y: 190, width: 450, height: 35 } },
            { type: 'BankCard', value: '6222020100123456789', confidence: 0.98, bounding_box: { x: 250, y: 260, width: 380, height: 35 } },
            { type: 'Email', value: 'zhang.san@example.com', confidence: 0.92, bounding_box: { x: 400, y: 330, width: 320, height: 30 } }
        ]
    };
}


// --- Display Logic ---

function displayResults(result) {
    lastResults = result;
    
    // Display stats
    elements.resultsStats.innerHTML = `
        <div class="stat-item"><i class="fas fa-search"></i> 检测到 ${result.stats.total_detections} 项敏感信息</div>
        <div class="stat-item"><i class="fas fa-clock"></i> 处理时间 ${(result.stats.processing_time_ms / 1000000).toFixed(2)} ms</div>
        <div class="stat-item"><i class="fas fa-image"></i> 图片尺寸 ${result.stats.image_dimensions.width}x${result.stats.image_dimensions.height}</div>
    `;

    // Populate OCR Text Tab
    elements.ocrTextContainer.textContent = result.ocr_text;

    // Populate Visual Results Tab
    const imageUrl = URL.createObjectURL(selectedFile);
    elements.resultImage.src = imageUrl;
    
    // Clear previous bounding boxes
    elements.visualResultContainer.querySelectorAll('.bounding-box').forEach(box => box.remove());

    // Create new bounding boxes
    result.results.forEach(item => {
        const box = document.createElement('div');
        box.className = 'bounding-box';
        
        // Calculate relative position and size
        const relX = (item.bounding_box.x / result.stats.image_dimensions.width) * 100;
        const relY = (item.bounding_box.y / result.stats.image_dimensions.height) * 100;
        const relW = (item.bounding_box.width / result.stats.image_dimensions.width) * 100;
        const relH = (item.bounding_box.height / result.stats.image_dimensions.height) * 100;

        box.style.left = `${relX}%`;
        box.style.top = `${relY}%`;
        box.style.width = `${relW}%`;
        box.style.height = `${relH}%`;
        
        // Add a tooltip
        const tooltip = document.createElement('span');
        tooltip.className = 'bbox-tooltip';
        tooltip.textContent = item.type;
        box.appendChild(tooltip);

        elements.visualResultContainer.appendChild(box);
    });

    // Show results section and scroll to it
    elements.resultsSection.classList.remove('hidden');
    elements.resultsSection.scrollIntoView({ behavior: 'smooth' });
}

function switchTab(tabId) {
    elements.tabBtns.forEach(btn => btn.classList.toggle('active', btn.dataset.tab === tabId));
    document.querySelectorAll('.tab-content').forEach(content => {
        content.classList.toggle('active', content.id === `${tabId}-tab`);
    });
}

function exportResults() {
    if (!lastResults) {
        showError('没有可导出的结果');
        return;
    }
    
    const exportData = {
        timestamp: new Date().toISOString(),
        stats: lastResults.stats,
        ocr_text: lastResults.ocr_text,
        detections: lastResults.results
    };
    
    const dataStr = JSON.stringify(exportData, null, 2);
    const dataBlob = new Blob([dataStr], { type: 'application/json' });
    
    const link = document.createElement('a');
    link.href = URL.createObjectURL(dataBlob);
    link.download = `image_detection_results_${new Date().toISOString().slice(0, 19).replace(/:/g, '-')}.json`;
    link.click();
    
    URL.revokeObjectURL(link.href);
}

function clearResults() {
    elements.resultsSection.classList.add('hidden');
    elements.resultsStats.innerHTML = '';
    elements.ocrTextContainer.textContent = '';
    elements.visualResultContainer.querySelectorAll('.bounding-box').forEach(box => box.remove());
    lastResults = null;
    hideError();
}


// --- UI State Management ---

function showLoading() {
    elements.loading.classList.remove('hidden');
    elements.resultsSection.classList.add('hidden');
}

function hideLoading() {
    elements.loading.classList.add('hidden');
}

function showError(message) {
    elements.errorText.textContent = message;
    elements.errorMessage.classList.remove('hidden');
    elements.errorMessage.scrollIntoView({ behavior: 'smooth' });
}

function hideError() {
    elements.errorMessage.classList.add('hidden');
}