import os
import chardet


class Address:
    def __init__(self, file_path):
        self.file_path = file_path

    def detect_file_encoding(self):
        """自动检测文件的编码格式"""
        try:
            with open(self.file_path, 'rb') as f:
                # 读取文件的一部分内容来检测编码
                raw_data = f.read(10000)
                result = chardet.detect(raw_data)
                encoding = result['encoding']
                confidence = result['confidence']

                # 如果置信度较低或检测结果为None，则默认使用utf-8
                if encoding is None or confidence < 0.6:
                    return 'utf-8'
                return encoding
        except Exception as e:
            print(f"检测文件编码时出错: {e}")
            return 'utf-8'  # 默认回退到UTF-8

    def main(self):
        input_file_name = self.file_path
        folder_path, file_name = os.path.split(self.file_path)
        base_name = os.path.splitext(file_name)[0]
        output_dir = os.path.join(folder_path, 'out_put', base_name)

        # 确保输出目录存在
        if not os.path.exists(output_dir):
            os.makedirs(output_dir)

        output_file_name = os.path.join(output_dir, 'address.txt')

        # 检测文件编码
        file_encoding = self.detect_file_encoding()
        print(f"检测到文件编码: {file_encoding}")

        i = 1
        num = 0

        # 一次性打开输入和输出文件，避免在循环内重复打开关闭
        try:
            with open(input_file_name, 'r', encoding=file_encoding, errors='replace') as input_file, \
                    open(output_file_name, 'w', encoding='utf-8') as output_file:  # 输出统一使用UTF-8

                for line in input_file:
                    line = line.strip()  # 去除首尾空白字符
                    if not line:  # 跳过空行
                        i += 1
                        continue

                    # 分割句子和子句
                    for part1 in line.split('。'):
                        for part2 in part1.split('，'):
                            cleaned_part = part2.strip()
                            # 检查是否包含'省'和'市'
                            if cleaned_part and '省' in cleaned_part and '市' in cleaned_part:
                                output_file.write(f"{i},{cleaned_part}\n")
                                num += 1
                    i += 1

        except FileNotFoundError:
            print(f"错误：输入文件 '{input_file_name}' 不存在。")
            return 0
        except UnicodeDecodeError as e:
            print(f"解码错误：{e}。尝试使用不同的编码。")
            return 0
        except Exception as e:
            print(f"处理文件时发生未知错误：{e}")
            return 0

        print(f"处理完成，共找到 {num} 个地址。")
        return num


# 使用示例
if __name__ == "__main__":
    # 实例化Address类，传入要处理的文件路径
    address_extractor = Address("your_input_file.txt")  # 请替换为你的文件路径
    # 执行处理并获取找到的地址数量
    count = address_extractor.main()
    print(f"提取的地址数量: {count}")