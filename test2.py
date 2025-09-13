
from paddleocr import PaddleOCR
from PIL import Image, ImageOps
import numpy as np

img_path = r"C:\Users\Admin\Desktop\gogogogo\742110b3-0e35-4927-97f8-be13c1733549.png"  # 这里输入你的图片路径
img = Image.open(img_path)

# 缩放图片（宽度固定为700px，高度按比例调整）
width = 700
height = int(img.height * (width / img.width))
img_resized = img.resize((width, height), Image.LANCZOS)

# 转换为灰度图 + 二值化
img_gray = ImageOps.grayscale(img_resized)
img_binary = img_gray.point(lambda x: 255 if x > 170 else 0)

# 将二值图像转换为3通道RGB格式
img = np.stack((img_binary,) * 3, axis=-1)

# 实例化OCR
ocr = PaddleOCR(

    # 基础设置
    # device="gpu",
    # use_tensorrt=True,
    precision='fp16',
    # enable_hpi=True,
    enable_mkldnn=True,
    cpu_threads=2,

    use_doc_orientation_classify=False,
    use_doc_unwarping=False,
    use_textline_orientation=False,

    # 输入优化
    # text_det_limit_type="min",
    # text_det_limit_side_len=640,
    # text_recognition_batch_size=1,
    # text_det_thresh = 0.4,
    # text_det_box_thresh = 0.7,
    # text_det_unclip_ratio = 1.2,

    text_recognition_model_name="PP-OCRv4_mobile_rec",
    text_detection_model_name="PP-OCRv5_mobile_det",
    #text_recognition_model_dir="D:/ProgramData/Paddle_models/PP-OCRv4_mobile_rec",  # 文本识别模型
    #text_detection_model_dir="D:/ProgramData/Paddle_models/PP-OCRv5_mobile_det",  # 文本检测模型
    # textline_orientation_model_dir="D:/ProgramData/Paddle_models/PP-LCNet_x1_0_doc_ori"  # 方向分类模型

)

# 使用predict方法
result = ocr.predict(img)
print(result[0].keys())
for line in result[0]['rec_texts']:
    print(line)