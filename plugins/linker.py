import py2neo
from . import enhancer

import time
import random
import logging
import collections
import sys
import traceback

log = logging.getLogger("link.plugins.linker")

#NUM_WORKERS = 8

def listify(x):
    return x if isinstance(x, list) else [x]

class Linker(enhancer.Enhancer):
    def __init__(self):
        super(Linker, self).__init__()

        # this cypher statement looks up all entities associated with
        # the matching ad and creates a relationship between those
        # entities and the origin
        
    def get_matching_node_ids(self, parameters):
        tx = self.db.cypher
        
        cypher_query = """
        MATCH (n:Ad {""" + str(parameters["field_name"]) + """: {field_value}})
        RETURN n.id
        """
        
        cypher_output = tx.execute(cypher_query,parameters)

        matching_node_ids = [x['n.id'] for x in cypher_output]

        if matching_node_ids[0] is not None:
            return matching_node_ids
        else:
            return

    def get_node_ids_parents(self, node_ids_list):
        tx = self.db.cypher
        
        parameters = {"node_ids_list" : node_ids_list}
        
        cypher_query = """
        MATCH (n1:Ad)
        WHERE has(n1.id) and (n1.id in {node_ids_list})
        OPTIONAL MATCH (n1)--(e1:Entity)
        RETURN e1.identifier
        """

        cypher_output = tx.execute(cypher_query,parameters)
        #print cypher_output
        parent_entities = [x['e1.identifier'] for x in cypher_output if x['e1.identifier'] is not None]
        parent_entities = list(set(parent_entities))

        return parent_entities

    def connect_node_to_matching_node_ids_parents(self, node, matching_node_parents):
        tx = self.db.cypher
        
        parameters = {"node_id": node['id'],
        "matching_node_parents" : matching_node_parents,
        "field_name" : "TXT"}
        
        #If matching_node_parents is none:
        #Create a new entity and join all ads to that

        #Else:
        #Join current ad to all given matching_node_parents
        if matching_node_parents is not None:
            cypher_query = """
            MATCH (n:Ad {id : {node_id}})
            MATCH (e:Entity)
            WHERE (e.identifier IN {matching_node_parents})
            MERGE (n)<-[r:TXT]-(e)
            ON CREATE SET r.user='auto', r.reason='TXT'
            ON MATCH SET r.user='auto', r.reason='TXT'
            """
            tx.execute(cypher_query,parameters)


    def enhance(self, node):
        if 'text_signature' in node:
            log.info("==> Handling record [%s]" % node['id'])

            t0 = time.time()

            #print self.get_node_parents(node)

            #If NONE of those nodes have a parent and this one does not, create one and add this to all of them
            #For every one of those nodes that has a parent, add this and all other match_ids to that one
            
            #Get all matching node_ids
            try:
                print str(node['id']) + '--------------------'
                parameters = {"field_name" : "text_signature",
                "field_value" : node['text_signature']
                }

                matching_node_ids = self.get_matching_node_ids(parameters)
                print matching_node_ids
                matching_node_parents = self.get_node_ids_parents(matching_node_ids)
                print matching_node_parents

                self.connect_node_to_matching_node_ids_parents(node,matching_node_parents)

            except Exception, err:
                print(traceback.format_exc())

        else:
            print 'no text hash'
        t1 = time.time()
        log.info("<== Handled %s in %f seconds" % (str(node['id']), t1 - t0))
