
import asyncio
import yaml
from src import BiliUser


async def main():
    initTasks = []
    startTasks = []
    with open('users.yaml', 'r', encoding='utf-8') as f:
        users = yaml.load(f, Loader=yaml.FullLoader)
    for user in users['USERS']:
        if user['access_key']:
            biliUser = BiliUser(user['access_key'], user['shared_uid'])
            initTasks.append(biliUser.init())
            startTasks.append(biliUser.start())
    await asyncio.gather(*initTasks)
    await asyncio.gather(*startTasks)
if __name__ == '__main__':
    loop = asyncio.get_event_loop()
    loop.run_until_complete(main())
