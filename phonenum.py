
import re
import os
from phone import Phone

class PhoneNumber:
    def __init__(self, file_path):
        self.file_path = file_path

    def check(self, phone_number):
        p = Phone()
        data = p.find(phone_number)
        return data is not None

    def main(self):
        input_file_name = self.file_path
        folder_path, file_name = os.path.split(self.file_path)
        base_name = os.path.splitext(file_name)[0]
        output_dir = os.path.join(folder_path, 'out_put', base_name)
        if not os.path.exists(output_dir):
            os.makedirs(output_dir)
        output_file_name = os.path.join(output_dir, 'phone.txt')
        i = 1
        num = 0
        with open(input_file_name, 'r', encoding='utf-8') as input_file:
            for line in input_file:
                phone_list = re.findall(r'1\d{10}', line)
                if phone_list:
                    for phone in phone_list:
                        if self.check(phone):
                            with open(output_file_name, 'a') as output_file:
                                output_file.write(str(i) + ',' + phone + '\n')
                                num += 1
                i += 1
        return num