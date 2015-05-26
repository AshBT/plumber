from . import enhancer
import collections
import logging

log = logging.getLogger("link.plugins.imager")

class Imager(enhancer.SQLEnhancer):
    def __init__(self):
        super(Imager, self).__init__(db='memex_ht')

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
