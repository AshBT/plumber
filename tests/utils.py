from py2neo import Node
import nose.tools

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

class Runner(object):
    # subclasses have to set these fields
    ENHANCER = None
    OUTPUT = []
    INPUT = []
    IGNORE = [] # this is optional, setting a list of IGNORE
                # first checks that the fields exist before ignoring
                # their values. This is useful for instance, when
                # comparing "tweets" which may change over time.

    def check_ad_equals_expected(self, test_node, expected, ignore_fields=[]):
        try:
            enhancer = self.ENHANCER()
            enhancer.db.create(Node("Ad", **test_node))

            enhancer.run(skip_on_error=False)

            # use the ENHANCER's DB to check that results are as expected
            records = enhancer.db.cypher.execute("MATCH (ad) RETURN ad")

            assert len(records)==1, "We're expecting only one entry in the DB"

            ad = records[0].ad

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
                nose.tools.eq_(a, b)
        finally:
            enhancer.db.delete_all()

    def test_run(self):
        # Populates the DB one by one and runs the enhancer on it.
        # This structure lets us see captured log outputs.

        # check that the records contain the expected fields
        for test_node, exp in zip(self.INPUT, self.OUTPUT):
            yield self.check_ad_equals_expected, test_node, exp, self.IGNORE
