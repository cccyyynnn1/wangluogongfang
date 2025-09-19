import re
import os


class Id:
    def __init__(self, file_path):
        self.file_path = file_path
        self.output_dir = os.path.dirname(file_path)
        self.output_file = os.path.join(self.output_dir, 'id.txt')
    def check(self, id):
        if int(id[6:10]) % 4 == 0 or (int(id[6:10]) % 100 == 0 and int(id[6:10]) % 4 == 0):
            birthday = re.compile(
                '[1-9][0-9]{5}(19[0-9]{2}|20[0-9]{2})((01|03|05|07|08|10|12)(0[1-9]|[1-2][0-9]|3[0-1])|(04|06|09|11)('
                '0[1-9]|[1-2][0-9]|30)|02(0[1-9]|[1-2][0-9]))[0-9]{3}[0-9Xx]$')
        else:
            # 出生日期平年时合法性正则表达式
            birthday = re.compile(
                '[1-9][0-9]{5}(19[0-9]{2}|20[0-9]{2})((01|03|05|07|08|10|12)(0[1-9]|[1-2][0-9]|3[0-1])|(04|06|09|11)('
                '0[1-9]|[1-2][0-9]|30)|02(0[1-9]|1[0-9]|2[0-8]))[0-9]{3}[0-9Xx]$')
        if not (re.match(birthday, id)):
            return False

        mod = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
        jym = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2']
        sum = 0
        for i in range(0, 17):
            sum += int(id[i]) * mod[i]
        sum %= 11
        if (jym[sum]) == id[17]:
            return True
        else:
            return False

    def main(self):
        input_file_name = self.file_path
        folder_path, file_name = os.path.split(self.file_path)
        output_file_name = os.path.join(folder_path, 'out_put', file_name.replace('.txt', ''), 'id.txt')
        print(input_file_name)
        print(output_file_name)
        i = 1
        num = 0
        with open(input_file_name, 'r', encoding='utf-8') as input_file:
            for line in input_file:
                print(line)
                id_list = re.findall(r'\d{7}\d{10}[\dxX]', line)
                print(id_list)
                if len(id_list) != 0:
                    for id in id_list:
                        if self.check(id):
                            with open(output_file_name, 'a') as output_file:
                                output_file.write(str(i) + ',' + id + '\n')
                                num = num + 1
                i = i + 1
        return num
