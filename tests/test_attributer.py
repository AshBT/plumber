from plugins.attributer import Attributer
from . utils import Runner

_test_nodes = [
    {"id": 1, "phone": "1234"}, # this ad has a phone attribute,
                                # we want to make sure the existing
                                # number is included

    {"id": "84", "phone": 4567},# this ad does not have a phone
                                # attribute. we want to make
                                # sure it stays undisturbed for
                                # now. maybe promote later?
    {"id": 5823}
]

_expected_results = [
    {"height":["165"],
     "id":1,
     "phone":["1234", "2055419574"]},
    {"availability":["Outcall"],
     "build":["Curvy"],
     "bust":["36"],
     "cup":["C"],
     "email":["mimi.adora@yahoo.com"],
     "ethnicity":["asian"],
     "eyes":["Brown"],
     "hair":["Blonde"],
     "height":["5'6''"],
     "id":"84",
     "username":["mimiadora"],
     "weight":["140"],
     "phone": 4567},
    {"id":5823,
     "phone":["9413027218"]}
]

class TestAttributer(Runner):
    @classmethod
    def setup_class(cls):
        super(TestAttributer, cls).setup_class(Attributer, _test_nodes)

    def test_run(self):
        for t in super(TestAttributer, self).test_run(_expected_results):
            yield t
