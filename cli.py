import sys
import json
import os
from address import Address
from id import Id
from name import Name
from nation import Nation
from phonenum import PhoneNumber


def extract_data(file_path):
    """
    根据文件路径提取所有信息，并返回一个字典。
    """
    if not os.path.exists(file_path):
        return {"error": f"File not found: {file_path}"}

    # 实例化所有提取器
    address_extractor = Address(file_path)
    id_extractor = Id(file_path)
    name_extractor = Name(file_path)
    nation_extractor = Nation(file_path)
    phone_extractor = PhoneNumber(file_path)

    results = {
        "name_count": name_extractor.main(),
        "id_count": id_extractor.main(),
        "phone_count": phone_extractor.main(),
        "nation_count": nation_extractor.main(),
        "address_count": address_extractor.main()
    }

    # 为了简化，这里只返回提取到的数量，您可以修改以返回详细列表
    # 例如：
    # "names": name_extractor.get_all_names()
    # "ids": id_extractor.get_all_ids()

    return results


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "No file path provided"}))
        sys.exit(1)

    input_file = sys.argv[1]

    try:
        data = extract_data(input_file)
        print(json.dumps(data, ensure_ascii=False))
    except Exception as e:
        print(json.dumps({"error": str(e)}))
        sys.exit(1)