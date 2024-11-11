"""
@Project   : onepush
@Author    : y1ndan
@Blog      : https://www.yindan.me
"""

# import logging
from loguru import logger

# import requests
from aiohttp import ClientSSLError, ClientSession, TCPConnector
from ssl import SSLCertVerificationError

# from requests.exceptions import SSLError

from .exceptions import NoSuchNotifierError
from .exceptions import OnePushException

# log = logging.getLogger('onepush')

log = None


class Provider(object):
    base_url = None
    site_url = None

    _params = None

    def __init__(self):
        self.method = 'post'
        self.datatype = 'data'
        self.url = None
        self.data = None
        self.proxy = None

    async def _prepare_url(self, **kwargs):
        ...

    async def _prepare_data(self, **kwargs):
        ...

    async def _send_message(self):
        if self.method.upper() == 'GET':
            response = self.request('get', self.url, params=self.data)
        elif self.method.upper() == 'POST':
            if self.datatype.lower() == 'json':
                response = await self.request('post', self.url, json=self.data)
            else:
                response = await self.request('post', self.url, data=self.data)
        else:
            raise OnePushException('Request method {} not supported.'.format(self.method))

        return response

    @property
    def params(self):
        return self._params

    @staticmethod
    def process_message(title, content):
        message = content
        if title and content:
            message = '{}\n\n{}'.format(title, content)
        if title and not content:
            message = title
        return message

    # @staticmethod
    async def request(self, method, url: str, **kwargs):
        if self.proxy:
            from aiohttp_socks import ProxyConnector
        # session = requests.Session()
        # session = (
        #     ClientSession() if not self.proxy else ClientSession(connector=TCPConnector(verify_ssl=False))
        # )
        # response = None
        try:
            sessions = []
            if self.proxy:
                connector = ProxyConnector.from_url(self.proxy)
                session = ClientSession(connector=connector, trust_env = True)
                sessions.append(session)
                response = await session.request(method, url, **kwargs)
            else:
                session = ClientSession(trust_env = True)
                sessions.append(session)
                response = await session.request(method, url, **kwargs)
            # log.debug('Response: {}'.format(response.text))
        except ClientSSLError as e:
            log.error(e)
            if self.proxy:
                connector = ProxyConnector.from_url(self.proxy, verify_ssl=False)
            else:
                connector = TCPConnector(verify_ssl=False)
            session = ClientSession(connector=connector, trust_env = True)
            sessions.append(session)
            response = await session.request(method, url.replace('https', 'http'), proxy=self.proxy, **kwargs)
            # log.debug('Response: {}'.format(response.text))
        except SSLCertVerificationError as e:
            log.error(e)
            if self.proxy:
                connector = ProxyConnector.from_url(self.proxy, verify_ssl=False)
            else:
                connector = TCPConnector(verify_ssl=False)
            session = ClientSession(connector=connector, trust_env = True)
            sessions.append(session)
            response = await session.request(method, url, proxy=self.proxy, **kwargs)
            # log.debug('Response: {}'.format(response.text))
        except Exception as e:
            log.error(e)
        finally:
            for session in sessions:
                await session.close()
            return response

    async def notify(self, **kwargs):
        self.proxy = kwargs.get('proxy')
        await self._prepare_url(**kwargs)
        await self._prepare_data(**kwargs)
        return await self._send_message()


from .providers import _all_providers  # noqa: E402


def all_providers():
    return list(_all_providers.keys())


def get_notifier(provider_name: str):
    if provider_name not in _all_providers:
        raise NoSuchNotifierError(provider_name)
    return _all_providers[provider_name]()


async def notify(provider_name: str, **kwargs):
    global log
    log = logger.bind(user=f"{provider_name} 推送")

    return await get_notifier(provider_name).notify(**kwargs)
