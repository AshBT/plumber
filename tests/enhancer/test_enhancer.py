""" This tests the general functionality of our enhancer class. It is
    generic.
"""
from plugins.enhancer import Enhancer, SQLEnhancer
from py2neo import Node, Graph
import requests
from .. utils import check_equals, check_exists, check_not_exists

class MyEnhancer(Enhancer):
    def cleanup(self):
        self.cleanup_called = True

    def enhance(self, node):
        # add an attribute to the node
        node['foobar'] = "inserted"

class MySQLEnhancer(SQLEnhancer):
    def enhance(self, node):
        query = "SELECT title, url, website, phone FROM ads LIMIT 1"
        def callback(title, url, website, phone):
            node["title"] = title
            node["url"] = url
            node["website"] = website
            node["phone"] = phone
        self.do_query(query, callback=callback)

class TestEnhancer(object):
    @classmethod
    def setup_class(cls):
        cls.ENHANCER = MyEnhancer()
        cls.SQL_ENHANCER = MySQLEnhancer(db='memex_ht')
        # seed some data
        test_nodes = [
            {"phone": ["123", "456"],
             "image": ["abc"],
             "text": "hello"},
            {"phone": ["123", "872"],
             "text": "goodbye"}
        ]
        for node in test_nodes:
            cls.ENHANCER.db.create(Node("Ad", **node))

    @classmethod
    def teardown_class(cls):
        cls.ENHANCER.db.delete_all()

    def test_enhancer_run(self):
        self.ENHANCER.run(skip_on_error=False)

        # use the ENHANCER's DB to check that results are as expected
        records = self.ENHANCER.db.cypher.execute("MATCH (ad) RETURN ad")
        for result in records:
            yield check_equals, result.ad['foobar'], 'inserted'

        yield check_equals, self.ENHANCER.cleanup_called, True

    def test_sql_enhancer_run(self):
        self.SQL_ENHANCER.run(skip_on_error=False)

        # use the ENHANCER's DB to check that results are as expected
        records = self.ENHANCER.db.cypher.execute("MATCH (ad) WHERE has(ad.title) RETURN ad")
        result = records[0].ad
        # these attributes should exist
        yield check_exists, result, 'phone'
        yield check_exists, result, 'url'
        # these attributes should not exist (since they are null from the query)
        yield check_not_exists, result, 'title'
        yield check_not_exists, result, 'website'
        # check that the phone number is correctly populated
        yield check_equals, result['phone'], "2055419574"
