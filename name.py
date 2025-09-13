import os
import jieba
import jieba.posseg as pseg


class Name:
    def __init__(self, file_path):
        self.file_path = file_path
        # 在处理文件前，先加载自定义词典
        self.load_custom_name_dict('custom_names_dict.txt')
        # 在处理文件前，先加载黑名单
        self.load_black_list_word('blacklist.txt')

    def load_custom_name_dict(self, custom_dict_file):
        """
        加载自定义人名词典，确保特定姓名（如“吴乐彬”）能被正确识别。
        :param custom_dict_file: 自定义词典文件名
        """
        if not os.path.isfile(custom_dict_file):
            print(f"Warning: Custom dictionary file {custom_dict_file} not found.")
            return

        with open(custom_dict_file, 'r', encoding='utf-8') as f:
            for line in f:
                parts = line.strip().split()
                if not parts:
                    continue
                word = parts[0]
                # 如果提供了词频和词性，则使用；否则使用默认值
                freq = int(parts[1]) if len(parts) > 1 else 100  # 默认给一个较高的词频
                tag = parts[2] if len(parts) > 2 else 'nr'  # 默认词性为人名
                jieba.add_word(word, freq=freq, tag=tag)
                print(f"Added '{word}' to dictionary with freq={freq} and tag={tag}")  # 可选：打印日志确认

    def get_all_name(self, message_content):
        # 使用pseg.cut对消息内容进行分词和词性标注
        words = pseg.cut(message_content)  # words 是一个生成器，产生 pair 对象
        names = []
        for word_flag in words:  # 每个 word_flag 是一个 pair 对象
            # 通过 .word 和 .flag 属性访问词语和词性
            if word_flag.flag == 'nr' and len(word_flag.word) >= 2:  # 'nr'代表中文人名
                names.append(word_flag.word)  # 这里也要修正，应该是 word_flag.word 而不是 words.word
        return names

    def change_tag(self, word_list):
        """
        为黑名单中的词语添加自定义词性
        :param word_list: 包含黑名单词语的列表
        """
        for word in word_list:
            # 使用jieba.add_word添加词语，并指定词频和词性
            # 词频（freq）必须是一个数字，不能是None
            jieba.add_word(word, freq=10, tag='n')

    def load_black_list_word(self, black_list_file_name):
        """
        加载黑名单文件中的词语
        :param black_list_file_name: 黑名单文件名
        """
        if not os.path.isfile(black_list_file_name):
            # 如果黑名单文件不存在，则创建一个空文件
            with open(black_list_file_name, 'w', encoding='utf-8') as f:
                pass

        with open(black_list_file_name, 'r', encoding='utf-8') as f:
            black_list = []
            for line in f.readlines():
                # 去除每行首尾的空白字符和换行符
                word = line.strip()
                if word:  # 确保不是空行
                    black_list.append(word)

            # 调用change_tag函数来处理黑名单中的词语
            self.change_tag(black_list)

    def main(self):


        input_file_name = self.file_path
        folder_path, file_name = os.path.split(self.file_path)
        base_name = os.path.splitext(file_name)[0]
        output_dir = os.path.join(folder_path, 'out_put', base_name)

        if not os.path.exists(output_dir):
            os.makedirs(output_dir)

        output_file_name = os.path.join(output_dir, 'name.txt')

        i = 1
        num = 0
        with open(input_file_name, 'r', encoding='utf-8') as input_file, \
                open(output_file_name, 'a', encoding='utf-8') as output_file:
            for line in input_file:
                if not line.strip():  # 跳过空行，避免不必要的处理
                    continue

                name_list = self.get_all_name(line)
                for name in name_list:
                    output_file.write(str(i) + ',' + name + '\n')
                    num += 1
                i += 1
        return num