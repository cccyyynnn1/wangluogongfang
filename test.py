from paddleocr import PaddleOCR
import cv2

# 初始化 PP-OCRv5 模型
ocr = PaddleOCR(
    ocr_version="PP-OCRv5",  # 明确指定版本[2](@ref)
    use_textline_orientation=False, # 禁用方向分类（简化输出）
    lang="ch"
)

# 执行预测
result = ocr.predict(r"C:\Users\Admin\Desktop\20220429112505_e50c1.png")

# 解析结果
if result and any(line[1][0] for line in result[0]):  # 检查非空
    for line in result[0]:
        text, score = line[1][0], line[1][1]
        print(f"文本: {text}, 置信度: {score:.2f}")
else:
    print("未识别到有效文本，请检查图像或模型配置")