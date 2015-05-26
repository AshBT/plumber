from plugins.instagrammer import Instagram
from . utils import Runner

_test_nodes = [
    {"text": "Follow me on INSTAGRAM: http://instagram.com/VIPVeronica<br />"},
    {"text": "Follow us on INSTAGRAM;  nyc_gentlemensclub"},
    {"text": "Instagram name: BBWbiancalatinadreamgirl"},
    {"text": "<b>Follow us</b> on INSTAGRAM; <b> nyc_gentlemensclub</b><br>"},
    {"text": "INSTAGRAM:aleo_lover4202010"},
    {"text": "I'm 100% real and my Instagram is SUPERBRIANNAFREAKY"},
    {"text": "TWITTER@VIPMSE & INSTAGRAM@PORTIASBESTCHOICE"},
    {"text": '<a rel="nofollow" target="_blank" href="http://instagram.com/icandycane">Instagram.com/icandycane</a>'},
    {"text": '''CARTY 7083208795 OR PROOF VIDEO on Insta.gram as QUEENOSIRIS
<br><br>
Still think I'm too beautiful to be real?? <a target="_blank" href="http://www.instagram.com/QUEENOSIRIS"> PICS AND PROOF VIDEO</a>'''},
    {"text": "I don't actually have an instagram username."}
]

_expected_results = [
    {"text": "Follow me on INSTAGRAM: http://instagram.com/VIPVeronica<br />"},
    {"text": "Follow us on INSTAGRAM;  nyc_gentlemensclub"},
    {"text": "Instagram name: BBWbiancalatinadreamgirl"},
    {"text": "<b>Follow us</b> on INSTAGRAM; <b> nyc_gentlemensclub</b><br>"},
    {"text": "INSTAGRAM:aleo_lover4202010"},
    {"text": "I'm 100% real and my Instagram is SUPERBRIANNAFREAKY"},
    {"text": "TWITTER@VIPMSE & INSTAGRAM@PORTIASBESTCHOICE"},
    {"text": '<a rel="nofollow" target="_blank" href="http://instagram.com/icandycane">Instagram.com/icandycane</a>'},
    {"text": '''CARTY 7083208795 OR PROOF VIDEO on Insta.gram as QUEENOSIRIS
<br><br>
Still think I'm too beautiful to be real?? <a target="_blank" href="http://www.instagram.com/QUEENOSIRIS"> PICS AND PROOF VIDEO</a>'''},
    {"text": "I don't actually have an instagram username."}
]

class TestInstagrammer(Runner):
    @classmethod
    def setup_class(cls):
        super(TestInstagrammer, cls).setup_class(Instagram, _test_nodes)

    def test_run(self):
        for t in super(TestInstagrammer, self).test_run(_expected_results, test_debug=True):
            yield t
