import datetime
from typing import Tuple
from siftool.common import *
from siftool import cosmos, sifchain


def get_block_times(sifnoded: sifchain.Sifnoded, first_block: int, last_block: int) -> List[Tuple[int, datetime.datetime]]:
    result = [(block, cosmos.parse_iso_timestamp(sifnoded.query_block(block)["block"]["header"]["time"]))
        for block in range(first_block, last_block)]
    return result
