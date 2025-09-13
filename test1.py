import numpy as np
from paddleocr import PaddleOCR
import cv2

# 初始化（无 use_pipeline）
ocr = PaddleOCR(
    use_textline_orientation=True,
    lang='ch',
)

# 执行识别
result = ocr.predict(r"C:\Users\Admin\Desktop\20220429112505_e50c1.png")

# 处理结果
for item in result:
    if isinstance(item, dict):
        if item.get('text_type', 'rec_texts') == 'rec_texts':  # 默认按文本处理
            print(f"内容: {item['rec_texts']}, 置信度: {item.get('rec_scores', 0):.2f}")

# 可视化（需安装 opencv）
image = cv2.imread(r"C:\Users\Admin\Desktop\20220429112505_e50c1.png")
for item in result:
    if item['text_type'] != 'rec_texts':
        continue
    pts = np.array(item['rec_boxes'], np.int32).reshape((-1,1,2))
    cv2.polylines(image, [pts], True, (0,255,0), 2)
cv2.imwrite('result.jpg', image)