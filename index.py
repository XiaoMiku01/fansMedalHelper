try:
    import aiohttp
except ImportError:
    import os
    import sys
    os.chdir(os.path.dirname(os.path.abspath(__file__)).split(__file__)[0])
    os.system('pip install -r requirements.txt -t .')

from main import run as handler
