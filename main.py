import sys
import os
import cv2
import numpy as np
from PyQt5.QtWidgets import QFileDialog, QApplication, QMainWindow, QDialog
from main_window import Ui_MainWindow
from child_window import Ui_Dialog
from paddleocr import PaddleOCR
from name import Name
from phonenum import PhoneNumber
from nation import Nation
from address import Address
from id import Id


class OCRProcessor:
    def __init__(self):
        self.ocr = PaddleOCR(
            device="gpu",
            precision='fp16',  # 保持FP16加速
            enable_mkldnn=False,  # 使用GPU时，MKLDNN关闭
            cpu_threads=4,  # 适当增加CPU线程数
            use_doc_orientation_classify=True,
            use_doc_unwarping=True,
            use_textline_orientation=True,
            text_recognition_model_name="PP-OCRv5_server_rec",  # 升级识别模型
            text_detection_model_name="PP-OCRv5_server_det",  # 升级检测模型
            #text_db_thresh=0.2,# 检测阈值
            #text_db_box_thresh=0.2, # 同样降低框阈值
            # rec_batch_num=6,                  # 如果需要批量处理，可调整
            # det_batch_num=6,                  # 如果需要批量处理，可调整
        )


    
    
    def process_and_extract(self, file_name):
        file_name = os.path.normpath(file_name)
        print(f"正在处理图片: {file_name}")

        try:
            # 尝试以彩色模式读取图片，以便进行后续处理
            image = cv2.imread(file_name, cv2.IMREAD_COLOR)
            print(f"image is None? {image is None}")
            if image is not None:
                print(f"image.shape: {image.shape}")
            # 增加鲁棒性检查：如果加载失败，cv2.imread会返回None
            if image is None:
                raise ValueError(f"无法加载图像，文件可能已损坏或格式不受支持。")

            # 新增鲁棒性检查：确保图像有有效的维度信息
            # 检查 image.shape 长度
            if not isinstance(image, np.ndarray) or len(image.shape) not in [2, 3]:
                raise ValueError(f"加载的图像对象无效或维度异常，shape={getattr(image, 'shape', None)}")
            if len(image.shape) == 2:
                image = cv2.cvtColor(image, cv2.COLOR_GRAY2BGR)
            elif len(image.shape) == 3 and image.shape[2] == 1:
                image = cv2.cvtColor(image, cv2.COLOR_GRAY2BGR)
            elif len(image.shape) == 3 and image.shape[2] == 3:
                pass  # 正常BGR
            else:
                raise ValueError(f"未知的图像通道数，shape={image.shape}")




            # 1. 重新缩放
            image_big = cv2.resize(image, dsize=(0, 0), fx=2, fy=2, interpolation=cv2.INTER_CUBIC)

            # 2. 灰度处理
            gray_image = cv2.cvtColor(image_big, cv2.COLOR_BGR2GRAY)


            # 3. 噪声消除
            denoised_image = cv2.medianBlur(gray_image, 3)

            # 4. 二值化
            #binary_image = cv2.adaptiveThreshold(denoised_image, 255, cv2.ADAPTIVE_THRESH_GAUSSIAN_C,cv2.THRESH_BINARY, 11, 2)
            _, binary_image = cv2.threshold(denoised_image, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)

            # 5. 形态学操作：膨胀(谨慎使用，可能会影响细节)
            #kernel = np.ones((1, 1), np.uint8)
            #dilated_image = cv2.dilate(binary_image, kernel, iterations=1)
            #dilated_image = cv2.morphologyEx(binary_image, cv2.MORPH_OPEN, kernel)
            # 6. 转换为三通道图像以适应OCR输入要求
            final_image = cv2.cvtColor(binary_image, cv2.COLOR_GRAY2BGR)

            print(f"最终输入模型的图像shape: {final_image}")
            # cv2.imwrite("processed_id_card.jpg", denoised_image)
            # 7. OCR识别
            try:
                # 使用预处理后的图像进行OCR识别
                result = self.ocr.predict(final_image,use_textline_orientation=True)
                # print("OCR result:", result)
                text = ""
                if result and isinstance(result[0], dict) and "rec_texts" in result[0]:
                    text = '\n'.join(result[0]["rec_texts"])
                return text
            except Exception as e:
                print("predict error:", e)
                # 出错时尝试使用原始图像进行识别
                try:
                    print("尝试使用原始图像进行OCR识别...")
                    result = self.ocr.predict(image)
                    if result and isinstance(result[0], dict) and "rec_texts" in result[0]:
                        text = '\n'.join(result[0]["rec_texts"])
                    return text
                except Exception as fallback_error:
                    print(f"备用方案也失败: {fallback_error}")
                    return ""

        except Exception as e:
            formatted_file_name = file_name.replace('\\', '/')
            print(f"处理文件 {formatted_file_name} 失败: {e}")
            return ""


class MyMainWindow(QMainWindow, Ui_MainWindow):
    IMAGE_EXTENSIONS = {
        ".jpg", ".jpeg", ".jpe", ".jfif", ".png", ".gif", ".webp",
        ".tiff", ".tif", ".bmp", ".ppm", ".pgm", ".pbm"
    }

    def __init__(self):
        super().__init__()
        self.setupUi(self)
        self.txt_files = []
        self.image_files = []
        self.name_num = 0
        self.id_num = 0
        self.phone_num = 0
        self.nation_num = 0
        self.address_num = 0
        self.ocr_processor = OCRProcessor()

    def button1_clicked(self):
        directory = QFileDialog.getExistingDirectory(self, "选取文件夹", "./")
        if directory:
            self.textEdit.setText(directory)
            self.txt_files = []
            self.image_files = []

    def button2_clicked(self):
        self.name_num = 0
        self.id_num = 0
        self.phone_num = 0
        self.nation_num = 0
        self.address_num = 0
        self.get_files()
        self.process_files()
        self.show_child()

    def check_box_1_clicked(self):
        checked = self.checkBox_1.isChecked()
        self.checkBox_2.setChecked(checked)
        self.checkBox_3.setChecked(checked)
        self.checkBox_4.setChecked(checked)
        self.checkBox_5.setChecked(checked)
        self.checkBox_6.setChecked(checked)

    def get_files(self):
        folder_path = self.textEdit.toPlainText()
        if not folder_path or not os.path.exists(folder_path):
            print("文件夹路径无效，请重新选择。")
            return

        self.txt_files = []
        self.image_files = []

        for file_name in os.listdir(folder_path):
            file_path = os.path.join(folder_path, file_name)
            if os.path.isfile(file_path):
                ext = os.path.splitext(file_name.lower())[1]
                if ext == ".txt":
                    self.txt_files.append(file_path)
                elif ext in self.IMAGE_EXTENSIONS:
                    self.image_files.append(file_path)

        if self.txt_files or self.image_files:
            output_dir = os.path.join(folder_path, 'out_put')
            if not os.path.exists(output_dir):
                os.makedirs(output_dir)

    def process_files(self):
        folder_path = self.textEdit.toPlainText()

        # 处理图片文件
        for file_path in self.image_files:
            txt_content = self.ocr_processor.process_and_extract(file_path)
            base_name = os.path.splitext(os.path.basename(file_path))[0]
            output_dir = os.path.join(folder_path, 'out_put', base_name)

            if not os.path.exists(output_dir):
                os.makedirs(output_dir)

            output_file_path = os.path.join(output_dir, f"{base_name}.txt")
            with open(output_file_path, 'w', encoding='utf-8') as output_file:
                output_file.write(txt_content)

            self.all_extract(output_file_path)

        # 处理已有的txt文件
        for file_path in self.txt_files:
            base_name = os.path.splitext(os.path.basename(file_path))[0]
            output_dir = os.path.join(folder_path, 'out_put', base_name)

            if not os.path.exists(output_dir):
                os.makedirs(output_dir)

            self.all_extract(file_path)

    def all_extract(self, file_path):
        if self.checkBox_2.isChecked():
            print('extract name')
            self.extract_name(file_path)
        if self.checkBox_3.isChecked():
            print('extract phone')
            self.extract_phone(file_path)
        if self.checkBox_4.isChecked():
            print('extract nation')
            self.extract_nation(file_path)
        if self.checkBox_6.isChecked():
            print('extract address')
            self.extract_address(file_path)
        if self.checkBox_5.isChecked():
            print('extract id')
            self.extract_id(file_path)

    def show_child(self):
        try:
            child_window = Child(str(self.name_num), str(self.phone_num), str(self.id_num), str(self.nation_num),
                                 str(self.address_num))
            child_window.exec_()
        except Exception as e:
            print(e)

    def extract_name(self, file_path):
        name = Name(file_path)
        self.name_num += name.main()

    def extract_phone(self, file_path):
        phone = PhoneNumber(file_path)
        self.phone_num += phone.main()

    def extract_nation(self, file_path):
        nation = Nation(file_path)
        self.nation_num += nation.main()

    def extract_address(self, file_path):
        address = Address(file_path)
        self.address_num += address.main()

    def extract_id(self, file_path):
        id = Id(file_path)
        self.id_num += id.main()


class Child(QDialog, Ui_Dialog):
    def __init__(self, name, phone, id, nation, address):
        super().__init__()
        self.setupUi(self)
        self.setWindowTitle("提取成功")
        self.label_3.setText(name)
        self.label_4.setText(phone)
        self.label_6.setText(id)
        self.label_8.setText(nation)
        self.label_10.setText(address)

    def close_windows(self):
        self.close()


if __name__ == '__main__':
    app = QApplication([])
    window = QMainWindow()
    main = MyMainWindow()
    main.show()
    sys.exit(app.exec_())