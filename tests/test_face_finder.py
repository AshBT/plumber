from plugins.face_finder import FaceFinder
from . utils import Runner

# This is a test on two separate nodes.
_test_nodes = [
    {"image_locations": ['https://s3.amazonaws.com/roxyimages/795fd1d9bd22e30dc9d31c9379e859c19ef9fb27.jpg', 'https://s3.amazonaws.com/roxyimages/53df35c33b244eba6a4a4e9c28e45fbe7d8afd31.jpg', 'https://s3.amazonaws.com/roxyimages/aa6bc06dd3667edd4d3b9a9818ceb51ca209448a.jpg'],
'image_locations': ['https://s3.amazonaws.com/roxyimages/4d15e7535a98a8da636e39e8eb0740aa0a7241ad.jpg']}
]

_expected_results = [
{"face_image_url":"memexadvertisements.s3-website-us-west-1.amazonaws.com/memexadvertisements795fd1d9bd22e30dc9d31c9379e859c19ef9fb27jpg.png",
"face_image_url":"memexadvertisements.s3-website-us-west-1.amazonaws.com/memexadvertisements53df35c33b244eba6a4a4e9c28e45fbe7d8afd31jpg.png",
"n_faces":2},
{"n_faces":0}
]

class TestFaceFinder(Runner):
    @classmethod
    def setup_class(cls):
        super(TestFaceFinder, cls).setup_class(FaceFinder, _test_nodes)

    def test_run(self):
        for t in super(TestFaceFinder, self).test_run(_expected_results):
            yield t
