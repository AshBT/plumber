import os
import pprint
import json
import time
import warnings
from py2neo import Graph, Node, Relationship, watch, GraphError
import pymysql.cursors
import ftfy

SQL_USER = os.environ['SQL_USER']
SQL_PASS = os.environ['SQL_PASS']
SQL_HOST = os.environ['SQL_HOST']
NEO_USER = os.environ['NEO_USER']
NEO_PASS = os.environ['NEO_PASS']
NEO_HOST = os.environ.get('NEO_HOST', 'localhost:7474')

# connect to local neo4j database
graph = Graph("http://{user}:{passwd}@{host}/db/data/".format(host=NEO_HOST, user=NEO_USER, passwd=NEO_PASS))
graph.delete_all()
try:
    graph.schema.create_uniqueness_constraint("Entity", "identifier")
except GraphError as e:
    print(e)
    pass

try:
    graph.schema.create_uniqueness_constraint("Ad", "id")
except GraphError as e:
    print(e)
    pass

with warnings.catch_warnings():
    warnings.simplefilter("once")
    connection = pymysql.connect(host=SQL_HOST,
            user=SQL_USER,
            passwd=SQL_PASS,
            db='memex_ht',
            cursorclass=pymysql.cursors.SSDictCursor,
            charset='utf8')

    count = 0
    watch("httpstream")
    try:
        #with connection.cursor() as cursor:
        cursor=connection.cursor()
        sql = "SELECT * from ads"
        cursor.execute(sql)
        for i, result in enumerate(cursor):
            # fix some bad unicode encodings
            for k in result:
                if isinstance(result[k], unicode):
                    result[k] = ftfy.fix_encoding(result[k])
            print(result)
            t0 = time.time()
            ad_node = Node("Ad", "Datum", **result)
            graph.create(ad_node)

            t1 = time.time()

            print "%09d: Elapsed" % i, (t1 - t0), "seconds"
    finally:
        connection.close()
