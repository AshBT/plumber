from plugins.face_finder import FaceFinder
from . utils import Runner

# This is a test on two separate nodes.
class TestFaceFinder(Runner):
    ENHANCER = FaceFinder
    INPUT = [
        {"image_locations": ['https://s3.amazonaws.com/roxyimages/53df35c33b244eba6a4a4e9c28e45fbe7d8afd31.jpg', 'https://s3.amazonaws.com/roxyimages/aa6bc06dd3667edd4d3b9a9818ceb51ca209448a.jpg']},
        {'image_locations': ['https://s3.amazonaws.com/roxyimages/4d15e7535a98a8da636e39e8eb0740aa0a7241ad.jpg']}
    ]

    OUTPUT = [
    {"face_image_url": [
        "https://s3-us-west-1.amazonaws.com/memexadvertisements/53df35c33b244eba6a4a4e9c28e45fbe7d8afd31.png"],
        "n_faces":1,
        "image_locations": ['https://s3.amazonaws.com/roxyimages/53df35c33b244eba6a4a4e9c28e45fbe7d8afd31.jpg', 'https://s3.amazonaws.com/roxyimages/aa6bc06dd3667edd4d3b9a9818ceb51ca209448a.jpg'],
    },
    {'image_locations': ['https://s3.amazonaws.com/roxyimages/4d15e7535a98a8da636e39e8eb0740aa0a7241ad.jpg'],
     "n_faces":0,
     "face_image_url": ""}
    ]
