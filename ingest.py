import os
import pprint
import json
import time
import warnings
from py2neo import Graph, Node, Relationship, watch, GraphError
import pymysql.cursors
import ftfy
import gevent
import gevent.queue

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

# add 100 at a time
batch_size = 100
workers = 8

# work queue for passing data around
work_queue = gevent.queue.Queue(maxsize=batch_size * workers)

def collect_ads(gid):
    nodes = []
    t0 = time.time()
    i = 0

    while True:
        ad = work_queue.get()
        i += 1
        if ad == "!!STOP!!":
            graph.create(*nodes)
            t1 = time.time()
            print "[Greenlet %02d] %09d: Elapsed" % (gid, i), (t1 - t0), "seconds"
            break
        else:
            nodes.append(ad)

        if len(nodes) == batch_size:
            graph.create(*nodes)
            t1 = time.time()
            nodes = []
            print "[Greenlet %02d] %09d: Elapsed" % (gid, i), (t1 - t0), "seconds"
            t0 = time.time()
            i = 0
        gevent.sleep(0)

threads = map(lambda x: gevent.spawn(collect_ads,x), range(workers))

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
        cursor=connection.cursor()
        sql = ["SELECT * FROM (SELECT * from ads ORDER BY id DESC LIMIT 1000) DUMMY",
               "SELECT * from ads LIMIT 1000",
               "SELECT * from ads where phone='4059285288'",
               "SELECT * from ads where phone='4047842015'",
               "SELECT * from ads where phone='8703950134'",
               "SELECT * from ads where phone='6984584328'",
               "SELECT * from ads where phone='5109788125'"]
        sql = " UNION ".join(sql)
        cursor.execute(sql)

        for i, result in enumerate(cursor):
            # fix some bad unicode encodings
            for k in result:
                if isinstance(result[k], unicode):
                    result[k] = ftfy.fix_encoding(result[k])
            # put result into work queue
            #print(result)
            work_queue.put(Node("Ad", "Datum", **result))
    finally:
        # send the stop signal to all workers
        for _ in range(workers):
            work_queue.put("!!STOP!!")
        gevent.joinall(threads)
        connection.close()
