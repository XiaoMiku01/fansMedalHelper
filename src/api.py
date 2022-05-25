
import asyncio
from hashlib import md5
import os
import random
import sys
import time
import json
import aiohttp
from typing import Union
from urllib.parse import urlencode


from aiohttp import ClientSession
sys.path.append(os.path.dirname(
    os.path.dirname(os.path.abspath(__file__))))


class Crypto:

    APPKEY = '1d8b6e7d45233436'
    APPSECRET = '560c52ccd288fed045859ed18bffd973'

    @staticmethod
    def md5(data: Union[str, bytes]) -> str:
        '''generates md5 hex dump of `str` or `bytes`'''
        if type(data) == str:
            return md5(data.encode()).hexdigest()
        return md5(data).hexdigest()

    @staticmethod
    def sign(data: Union[str, dict]) -> str:
        '''salted sign funtion for `dict`(converts to qs then parse) & `str`'''
        if isinstance(data, dict):
            _str = urlencode(data)
        elif type(data) != str:
            raise TypeError
        return Crypto.md5(_str + Crypto.APPSECRET)


class SingableDict(dict):
    @property
    def sorted(self):
        '''returns a alphabetically sorted version of `self`'''
        return dict(sorted(self.items()))

    @property
    def signed(self):
        '''returns our sorted self with calculated `sign` as a new key-value pair at the end'''
        _sorted = self.sorted
        return {**_sorted, 'sign': Crypto.sign(_sorted)}


class BiliApi:
    headers = {
        "User-Agent": "Mozilla/5.0 BiliDroid/6.73.1 (bbcallen@gmail.com) os/android model/Mi 10 Pro mobi_app/android build/6731100 channel/xiaomi innerVer/6731110 osVer/12 network/2",
    }
    from .user import BiliUser

    def __init__(self, u: BiliUser, s: ClientSession):
        self.u = u
        self.session = s

    def __check_response(self, resp: dict) -> dict:
        if resp['code'] != 0:
            raise Exception(resp['message'])
        return resp['data']

    async def getFansMedalandRoomID(self) -> dict:
        '''
        获取用户粉丝勋章和直播间ID
        '''
        url = "http://api.live.bilibili.com/xlive/app-ucenter/v1/fansMedal/panel"
        params = {
            "access_key": self.u.access_key,
            "actionKey": "appkey",
            "appkey": Crypto.APPKEY,
            "ts": int(time.time()),
            "page": 1,
            "page_size": 100,
        }
        first_flag = True
        while True:
            async with self.session.get(url, params=SingableDict(params).signed, headers=self.headers) as resp:
                data = self.__check_response(await resp.json())
                if first_flag and data['special_list']:
                    for item in data['special_list']:
                        yield item
                    first_flag = False
                for item in data['list']:
                    yield item
                if not data['list']:
                    break
                params['page'] += 1
    async def likeInteract(self, room_id: int):
        '''
        点赞 *3
        '''
        url = "http://api.live.bilibili.com/xlive/web-ucenter/v1/interact/likeInteract"
        data = {
            "access_key": self.u.access_key,
            "actionKey": "appkey",
            "appkey": Crypto.APPKEY,
            "ts": int(time.time()),
            "roomid": room_id,
        }
        for _ in range(3):
            try:
                async with self.session.post(url, data=SingableDict(data).signed, headers=self.headers.update({
                    "Content-Type": "application/x-www-form-urlencoded",
                })) as resp:
                    self.__check_response(await resp.json())
                    await asyncio.sleep(2)
            except aiohttp.ClientError:
                pass

    async def shareRoom(self, room_id: int):
        '''
        分享直播间
        '''
        url = "http://api.live.bilibili.com/xlive/app-room/v1/index/TrigerInteract"
        data = {
            "access_key": self.u.access_key,
            "actionKey": "appkey",
            "appkey": Crypto.APPKEY,
            "ts": int(time.time()),
            "interact_type": 3,
            "roomid": room_id,
        }
        for _ in range(5):
            try:
                async with self.session.post(url, data=SingableDict(data).signed, headers=self.headers.update({
                    "Content-Type": "application/x-www-form-urlencoded",
                })) as resp:
                    self.__check_response(await resp.json())
                await asyncio.sleep(3)
            except aiohttp.ClientError:
                pass
    async def sendDanmaku(self, room_id: int) -> str:
        '''
        发送弹幕
        '''
        url = "http://api.live.bilibili.com/xlive/app-room/v1/dM/sendmsg"
        danmakus = [
            "(⌒▽⌒).",
            "（￣▽￣）.",
            "(=・ω・=).",
            "(｀・ω・´).",
            "(〜￣△￣)〜.",
            "(･∀･).",
            "(°∀°)ﾉ.",
            "(￣3￣).",
            "╮(￣▽￣)╭.",
            "_(:3」∠)_.",
            "(^・ω・^ ).",
            "(●￣(ｴ)￣●).",
            "ε=ε=(ノ≧∇≦)ノ.",
            "⁄(⁄ ⁄•⁄ω⁄•⁄ ⁄)⁄.",
            "←◡←.",
        ]
        params = {
            "access_key": self.u.access_key,
            "actionKey": "appkey",
            "appkey": Crypto.APPKEY,
            "ts": int(time.time()),
        }
        data = {
            "cid": room_id,
            "msg": random.choice(danmakus),
            "rnd": int(time.time()),
            "color": "16777215",
            "fontsize": "25",
        }
        async with self.session.post(url, params=SingableDict(params).signed, data=data, headers=self.headers.update({
            "Content-Type": "application/x-www-form-urlencoded",
        })) as resp:
            return json.loads(self.__check_response(await resp.json())['mode_info']['extra'])['content']

    async def loginVerift(self):
        '''
        登录验证
        '''
        url = "http://app.bilibili.com/x/v2/account/mine"
        params = {
            "access_key": self.u.access_key,
            "actionKey": "appkey",
            "appkey": Crypto.APPKEY,
            "ts": int(time.time()),
        }
        async with self.session.get(url, params=SingableDict(params).signed, headers=self.headers) as resp:
            return self.__check_response(await resp.json())

    async def doSign(self):
        '''
        直播区签到
        '''
        url = "http://api.live.bilibili.com/rc/v1/Sign/doSign"
        params = {
            "access_key": self.u.access_key,
            "actionKey": "appkey",
            "appkey": Crypto.APPKEY,
            "ts": int(time.time()),
        }
        async with self.session.get(url, params=SingableDict(params).signed, headers=self.headers) as resp:
            return self.__check_response(await resp.json())

    async def getUserInfo(self):
        '''
        用户直播等级
        '''
        url = "http://api.live.bilibili.com/xlive/app-ucenter/v1/user/get_user_info"
        params = {
            "access_key": self.u.access_key,
            "actionKey": "appkey",
            "appkey": Crypto.APPKEY,
            "ts": int(time.time()),
        }
        async with self.session.get(url, params=SingableDict(params).signed, headers=self.headers) as resp:
            return self.__check_response(await resp.json())
        