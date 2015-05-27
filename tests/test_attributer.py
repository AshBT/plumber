from plugins.attributer import Attributer
from . utils import Runner

class TestAttributer(Runner):
    ENHANCER = Attributer
    INPUT = [
        {"id": 1, "phone": "1234"}, # this ad has a phone attribute,
                                    # we want to make sure the existing
                                    # number is included

        {"id": "84", "phone": 4567},# this ad does not have a phone
                                    # attribute. we want to make
                                    # sure it stays undisturbed for
                                    # now. maybe promote later?
        {"id": 5823}
    ]
    OUTPUT = [
        {"height":["165"],
         "id":1,
         "phone":["1234", "5555559574"]},
        {"availability":["sometimes"],
         "build":["Success"],
         "bust":["boom"],
         "cup":["stanley"],
         "email":["johndoe@foobar.com"],
         "ethnicity":["ethnicity"],
         "eyes":["Brown"],
         "hair":["Bald"],
         "height":["5'6''"],
         "id":"84",
         "username":["johndoe"],
         "weight":["140"],
         "phone": 4567},
        {"id":5823,
         "phone":["5555557218"]}
    ]
