from py2neo import Node

def check_equals(x, y):
    assert x == y

def check_exists(x, attribute):
    try:
        x[attribute]
    except KeyError:
        assert False, "missing key {}".format(attribute)

def check_not_exists(x, attribute):
    try:
        x[attribute]
    except KeyError:
        assert True, "key {} unexpectedly exists".format(attribute)

def check_ad_equals_expected(ad, expected, ignore_fields=[]):
    # check that all expected keys are present
    for key in expected:
        assert (ad.properties[key] is not None), ("{} is missing".format(key))

    for key in ad.properties:
        a = ad[key]
        check_exists(expected, key)
        b = expected[key]
        if key in ignore_fields:
            continue
        if isinstance(a, list):
            a.sort()
        if isinstance(b, list):
            b.sort()
        check_equals(a, b)

class Runner(object):
    @classmethod
    def setup_class(cls, EnhancerToTest, test_nodes):
        cls.ENHANCER = EnhancerToTest()
        # seed some data
        for node in test_nodes:
            cls.ENHANCER.db.create(Node("Ad", **node))

    @classmethod
    def teardown_class(cls):
        cls.ENHANCER.db.delete_all()

    def test_run(self, expected_results, ignore_fields=[]):
        """
        Setting a list of ignore_fields first checks that the fields
        exist before ignoring their values. This is useful for instance,
        when comparing "tweets" which may change over time.
        """
        self.ENHANCER.run(skip_on_error=False)

        # use the ENHANCER's DB to check that results are as expected
        records = self.ENHANCER.db.cypher.execute("MATCH (ad) RETURN ad")

        # first, check that we haven't *added* any records
        yield check_equals, len(records), len(expected_results)

        # second, check that the records contain the expected fields
        for result, exp in zip(records, expected_results):
            yield check_ad_equals_expected, result.ad, exp, ignore_fields
