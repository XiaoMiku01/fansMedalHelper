
import asyncio
import yaml
from src import BiliUser

from apscheduler.schedulers.blocking import BlockingScheduler
from apscheduler.triggers.cron import CronTrigger


async def main():
    initTasks = []
    startTasks = []

    for user in users['USERS']:
        if user['access_key']:
            biliUser = BiliUser(user['access_key'], user['shared_uid'])
            initTasks.append(biliUser.init())
            startTasks.append(biliUser.start())
    await asyncio.gather(*initTasks)
    await asyncio.gather(*startTasks)


def run():
    loop = asyncio.new_event_loop()
    loop.run_until_complete(main())


if __name__ == '__main__':
    with open('users.yaml', 'r', encoding='utf-8') as f:
        users = yaml.load(f, Loader=yaml.FullLoader)
    cron = users.get('CRON', None)
    if cron:
        print('使用内置定时器,开启定时任务,等待时间到达后执行')
        schedulers = BlockingScheduler()
        schedulers.add_job(
            run,
            CronTrigger.from_crontab(cron),
        )
        schedulers.start()
    else:
        print('外部调用,开启任务')
        loop = asyncio.get_event_loop()
        loop.run_until_complete(main())
