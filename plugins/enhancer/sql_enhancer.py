from . import enhancer
import pymysql
import os
import logging
import ftfy

log = logging.getLogger("link.sql_enhancer")

def fix_encoding(elem):
    if isinstance(elem, unicode):
        return ftfy.fix_encoding(elem)
    else:
        return elem

class SQLEnhancer(enhancer.Enhancer):
    def __init__(self, db, user=os.environ['SQL_USER'], password=os.environ['SQL_PASS'], host=os.environ['SQL_HOST'], **kwargs):
        super(SQLEnhancer, self).__init__(**kwargs)

        log.debug("Connecting to MySQL DB at %s" % host)
        self.connection = pymysql.connect(host=host,
                user=user,
                passwd=password,
                db=db,
                charset='utf8')

    def do_query(self, query, parameters=(), callback=lambda *x: log.debug(x)):
        #with self.connection.cursor() as cursor:
        cursor = self.connection.cursor()
        cursor.execute(query, parameters)
        for result in cursor:
            # fix unicode encoding errors
            result = map(fix_encoding, result)
            callback(*result)

    def cleanup(self):
        self.connection.close()
