from . import enhancer
import collections
import logging
import os

log = logging.getLogger("link.plugins.imager")

class Imager(enhancer.SQLEnhancer):
    def __init__(self, **kwargs):
        db_name = os.environ.get("SQL_DB", "memex_ht")
        super(Imager, self).__init__(db=db_name, **kwargs)

    def enhance(self, node):
        # this makes a db query for each node
        sql = "SELECT id, url, location from images WHERE ads_id=%s"
        datum = collections.defaultdict(list)
        def row_handler(image_id, url, location):
            datum['image_urls'].append(url if url else '')
            datum['image_locations'].append(location if location else '')
            datum['image_ids'].append(image_id)
        self.do_query(sql, node['id'], callback=row_handler)

        for key in datum:
            node[key] = datum[key]
