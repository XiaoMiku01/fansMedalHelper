import os
import json
import yaml

if __name__ == '__main__':
    users = os.environ.get('USERS')
    if users:
        with open("users.yaml", "w") as fw:
            yaml.dump(json.loads(users), fw)
        with open("users.yaml", "r") as fr:
            users = yaml.load(fr, Loader=yaml.FullLoader)
    os.system("git pull")
    os.system("python main.py")
