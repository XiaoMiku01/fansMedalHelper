import json
import os
import sys
from loguru import logger
import warnings
import asyncio
import aiohttp
import itertools
from src import BiliUser

log_file = os.path.join(os.path.dirname(__file__), "log/fansMedalHelper_{time:YYYY-MM-DD}.log")
log_format="<green>{time:YYYY-MM-DD HH:mm:ss.SSS}</green> | <level>{level: <8}</level> | <cyan>{name}</cyan>:<cyan>{function}</cyan>:<cyan>{line}</cyan> - <level>{message}</level>"

logger.remove()
logger.add(
    sys.stdout,
    format=log_format,
    backtrace=True,
    diagnose=True,
    level="INFO"
)
log = logger.bind(user="B站粉丝牌助手")
__VERSION__ = "0.3.8"

warnings.filterwarnings(
    "ignore",
    message="The localize method is no longer necessary, as this time zone supports the fold attribute",
)
os.chdir(os.path.dirname(os.path.abspath(__file__)).split(__file__)[0])

try:
    if os.environ.get("USERS"):
        users = json.loads(os.environ.get("USERS"))
    else:
        import yaml

        with open("users.yaml", "r", encoding="utf-8") as f:
            users = yaml.load(f, Loader=yaml.FullLoader)
    if users.get("WRITE_LOG_FILE"):
        logger.add(
            log_file if users["WRITE_LOG_FILE"] == True else users["WRITE_LOG_FILE"],
            format=log_format,
            backtrace=True,
            diagnose=True,
            rotation="00:00",
            retention="30 days",
            level="DEBUG"
        )
    assert users["ASYNC"] in [0, 1], "ASYNC参数错误"
    assert users["LIKE_CD"] >= 0, "LIKE_CD参数错误"
    # assert users['SHARE_CD'] >= 0, "SHARE_CD参数错误"
    assert users["DANMAKU_CD"] >= 0, "DANMAKU_CD参数错误"
    assert users["DANMAKU_NUM"] >= 0, "DANMAKU_NUM参数错误"
    assert users["DANMAKU_CHECK_LIGHT"] in [0, 1], "DANMAKU_CHECK_LIGHT参数错误"
    assert users["DANMAKU_CHECK_LEVEL"] in [0, 1], "DANMAKU_CHECK_LEVEL参数错误"
    assert users["WATCHINGLIVE"] >= 0, "WATCHINGLIVE参数错误"
    assert users["WEARMEDAL"] in [0, 1], "WEARMEDAL参数错误"
    config = {
        "ASYNC": users["ASYNC"],
        "LIKE_CD": users["LIKE_CD"],
        # "SHARE_CD": users['SHARE_CD'],
        "DANMAKU_CD": users["DANMAKU_CD"],
        "DANMAKU_NUM": users["DANMAKU_NUM"],
        "DANMAKU_CHECK_LIGHT": users["DANMAKU_CHECK_LIGHT"],
        "DANMAKU_CHECK_LEVEL": users["DANMAKU_CHECK_LEVEL"],
        "WATCHINGLIVE": users["WATCHINGLIVE"],
        "WEARMEDAL": users["WEARMEDAL"],
        "SIGNINGROUP": users.get("SIGNINGROUP", 2),
        "PROXY": users.get("PROXY"),
        "STOPWATCHINGTIME": None,
    }
    stoptime = users.get("STOPWATCHINGTIME", None)
    if stoptime:
        import time
        now = int(time.time())
        if isinstance(stoptime, int):
            delay = now + int(stoptime)
        else:
            delay = int(time.mktime(time.strptime(f'{time.strftime("%Y-%m-%d", time.localtime(now))} {stoptime}', "%Y-%m-%d %H:%M:%S")))
            delay = delay if delay > now else delay + 86400
        config["STOPWATCHINGTIME"] = delay
        log.info(f"本轮任务将在 {time.strftime('%Y-%m-%d %H:%M:%S', time.localtime(config['STOPWATCHINGTIME']))} 结束")
except Exception as e:
    log.error(f"读取配置文件失败,请检查配置文件格式是否正确: {e}")
    exit(1)


@log.catch
async def main():
    messageList = []
    session = aiohttp.ClientSession(trust_env=True)
    try:
        log.warning("当前版本为: " + __VERSION__)
        resp = await (
            await session.get(
                "http://version.fansmedalhelper.1961584514352337.cn-hangzhou.fc.devsapp.net/"
            )
        ).json()
        if resp["version"] != __VERSION__:
            log.warning("新版本为: " + resp["version"] + ",请更新")
            log.warning("更新内容: " + resp["changelog"])
            messageList.append(f"当前版本: {__VERSION__} ,最新版本: {resp['version']}")
            messageList.append(f"更新内容: {resp['changelog']} ")
        if resp["notice"]:
            log.warning("公告: " + resp["notice"])
            messageList.append(f"公告: {resp['notice']}")
    except Exception as ex:
        messageList.append(f"检查版本失败，{ex}")
        log.warning(f"检查版本失败，{ex}")
    initTasks = []
    startTasks = []
    catchMsg = []
    for user in users["USERS"]:
        if user["access_key"]:
            biliUser = BiliUser(
                user["access_key"],
                user.get("white_uid", ""),
                user.get("banned_uid", ""),
                config,
            )
            initTasks.append(biliUser.init())
            startTasks.append(biliUser.start())
            catchMsg.append(biliUser.sendmsg())
    try:
        await asyncio.gather(*initTasks)
        await asyncio.gather(*startTasks)
    except Exception as e:
        log.exception(e)
        # messageList = messageList + list(itertools.chain.from_iterable(await asyncio.gather(*catchMsg)))
        messageList.append(f"任务执行失败: {e}")
    finally:
        messageList = messageList + list(
            itertools.chain.from_iterable(await asyncio.gather(*catchMsg))
        )
    [log.info(message) for message in messageList]
    if users.get("SENDKEY", ""):
        await push_message(session, users["SENDKEY"], "  \n".join(messageList))
    await session.close()
    if users.get("MOREPUSH", ""):
        from onepush import notify

        notifier = users["MOREPUSH"]["notifier"]
        params = users["MOREPUSH"]["params"]
        await notify(
            notifier,
            title=f"【B站粉丝牌助手推送】",
            content="  \n".join(messageList),
            **params,
            proxy=config.get("PROXY"),
        )
        log.info(f"{notifier} 已推送")


def run(*args, **kwargs):
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)
    loop.run_until_complete(main())
    log.info("任务结束，等待下一次执行。")


async def push_message(session, sendkey, message):
    url = f"https://sctapi.ftqq.com/{sendkey}.send"
    data = {"title": f"【B站粉丝牌助手推送】", "desp": message}
    await session.post(url, data=data)
    log.info("Server酱已推送")


if __name__ == "__main__":
    cron = users.get("CRON", None)

    if cron:
        from apscheduler.schedulers.blocking import BlockingScheduler
        from apscheduler.triggers.cron import CronTrigger

        log.info(f"使用内置定时器 {cron}，开启定时任务，等待时间到达后执行。")
        schedulers = BlockingScheduler()
        schedulers.add_job(run, CronTrigger.from_crontab(cron), misfire_grace_time=3600)
        schedulers.start()
    elif "--auto" in sys.argv:
        from apscheduler.schedulers.blocking import BlockingScheduler
        from apscheduler.triggers.interval import IntervalTrigger
        import datetime

        log.info("使用自动守护模式，每隔 24 小时运行一次。")
        scheduler = BlockingScheduler(timezone="Asia/Shanghai")
        scheduler.add_job(
            run,
            IntervalTrigger(hours=24),
            next_run_time=datetime.datetime.now(),
            misfire_grace_time=3600,
        )
        scheduler.start()
    else:
        log.info("未配置定时器，开启单次任务。")
        run()
        log.info("任务结束")
