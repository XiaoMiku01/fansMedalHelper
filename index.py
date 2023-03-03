from main import run


# Tencent SCF
def main_handler(event, context):
    run()
    return


# Aliyun FC
def handler(event, context):
    run()
    return


if __name__ == '__main__':
    run()
