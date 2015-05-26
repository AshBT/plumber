import py2neo
from . import enhancer

import time
import random
import logging
import collections

log = logging.getLogger("link.plugins.text_hash_linker")

#NUM_WORKERS = 8

def listify(x):
    return x if isinstance(x, list) else [x]

class NeighborLinker(enhancer.Enhancer):
    def __init__(self):
        super(NeighborLinker, self).__init__()

        # this cypher statement looks up all entities associated with
        # the matching ad and creates a relationship between those
        # entities and the origin
        self.__add_similar =  """
        Match (n1:Ad {id : {id}})--(e1:Entity)--(n2:Ad)--(e2:Entity)--(n3:Ad) WHERE NOT (n1)--(:Entity)--(n3)
        MERGE (n1)-[r1:NBHR]->(e2)
        ON CREATE SET r1.user='auto', r1.reason='nearest-neighbor',
        """
        

    def enhance(self, node):
        log.info("==> Handling record [%s]" % node['id'])

        t0 = time.time()

        tx = self.db.cypher.begin()

        parameters = {"id": node['id']
        }
        
        tx.append(self.__add_similar, parameters)
        # this batches up all cypher queries and executes them here
        tx.commit()

        t1 = time.time()
        log.info("<== Handled %s in %f seconds" % (str(node['id']), t1 - t0))
