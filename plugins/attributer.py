from . import enhancer
import os
import collections
import logging

log = logging.getLogger("link.plugins.attributer")

class Attributer(enhancer.SQLEnhancer):
    def __init__(self):
        super(Attributer, self).__init__(db='memex_ht')

    def enhance(self, node):
        # this makes a db query for each node
        cursor = self.connection.cursor()
        sql = "SELECT attribute, value from ads_attributes WHERE ads_id=%s"
        # use a set to keep unique values
        datum = collections.defaultdict(set)
        self.do_query(sql, node['id'], lambda key, value: datum[key].add(value))

        for key in datum:
            # append the node's value if the key already exists
            if key in node:
                datum[key].add(node[key])
            node[key] = datum[key]
