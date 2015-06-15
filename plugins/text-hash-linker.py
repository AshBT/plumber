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

class TextHashLinker(enhancer.Enhancer):
    def __init__(self, **kwargs):
        super(TextHashLinker, self).__init__(**kwargs)

        # this cypher statement looks up all entities associated with
        # the matching ad and creates a relationship between those
        # entities and the origin
        self.__add_similar =  """
        MATCH (n1:Ad {id: {id}})<--(e1:Entity)
        MATCH (n2:Ad {text_signature: {text_signature}})<--(e2:Entity)  
        MERGE (e1)-[r1:BY_TXT]->(n2)
        ON CREATE SET r1.user='auto', r1.reason='contains-similar-text', r1.matches_to=n1.id
        ON MATCH SET r1.user='auto', r1.reason='contains-similar-text', r1.matches_to=n1.id
        MERGE (e2)-[r2:BY_TXT]->(n1)
        ON CREATE SET r2.user='auto', r2.reason='contains-similar-text', r2.matches_to=n2.id
        ON MATCH SET r1.user='auto', r1.reason='contains-similar-text', r1.matches_to=n1.id
        """
        

    def enhance(self, node):
        if 'text_signature' in node:
            log.info("==> Handling record [%s]" % node['id'])

            t0 = time.time()

            tx = self.db.cypher.begin()

            parameters = {"id": node['id'],
                "text_signature": node['text_signature'],
            }
            
            tx.append(self.__add_similar, parameters)
            # this batches up all cypher queries and executes them here
            tx.commit()

            t1 = time.time()
            log.info("<== Handled %s in %f seconds" % (str(node['id']), t1 - t0))
