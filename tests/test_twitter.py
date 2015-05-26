from plugins.twitter import Twitter
from . utils import Runner

# from a subset of actual data
_test_nodes = [
    {"twitter": "https://twitter.com/FantasysEscorts"},
    {"twitter": "https://twitter.com/@FantasysEscorts"},
    {"twitter": "https://twitter.com/escortdenver"},
    {"twitter": "https://twitter.com/https://twitter.com/KristenBaaby"},
    {"twitter": "https://twitter.com/&#8450;all Gi&#8477;&#8466;s in De&#8466;hi | 99IIII2O5I | Dhaula Kuan &#8495;sc"},
    {"twitter": "https://twitter.com/www.twitter.com/theirishmartini"},
    {"twitter": "https://twitter.com/https://twitter.com/#!/nycasianescor"},
    {"twitter": "https://twitter.com/Call Girls In Delhi 9811539366 Escort Service In D"},
    {"twitter": "https://twitter.com/https://mobile.twitter.com/monicasweets85?original_referer=http%3A%2F%2Fstatic.parastorage.com%2Fser"},
    {"twitter": """/phone/561-502-8304" title="Escort Phone 561-502-8304" class="fontyellow">561-502-8304</a></h2>

<div class="margintop5"></div>
For the mature gentlemen who loves paradise away for an hour . For more information call Paris @ 5615028304 in"""}
]

_expected_results = [
    {"twitter": "https://twitter.com/FantasysEscorts", "tweets": [], "twitter_followers": [], "twitter_friends": [],
     "twitter_name":"Fantasys Escorts",
     "twitter_profile_pic":"https://pbs.twimg.com/profile_images/3207171569/2476364f95748cdcecb4666bfaa706a5_normal.jpeg",
     "twitter_profile_url":"http://FantasysEscortService.com",
     "twitter_description":"",
     "twitter_background_pic":"https://pbs.twimg.com/profile_background_images/187044774/A-FantasyEscorts_.jpg"},

    {"twitter": "https://twitter.com/@FantasysEscorts", "tweets": [], "twitter_followers": [], "twitter_friends": [],
     "twitter_name":"Fantasys Escorts",
     "twitter_profile_pic":"https://pbs.twimg.com/profile_images/3207171569/2476364f95748cdcecb4666bfaa706a5_normal.jpeg",
     "twitter_profile_url":"http://FantasysEscortService.com",
     "twitter_description":"",
     "twitter_background_pic":"https://pbs.twimg.com/profile_background_images/187044774/A-FantasyEscorts_.jpg"},

    {"twitter": "https://twitter.com/escortdenver", "tweets": [], "twitter_followers": [], "twitter_friends": [],
     "twitter_name":"Denver Escorts",
     "twitter_background_pic": "https://pbs.twimg.com/profile_background_images/378800000061126191/8fe30a2f2b0b73e57ceaecc8c86c400b.jpeg",
     "twitter_description": "",
     'twitter_profile_url': 'http://EscortsinDenver.com',
     'twitter_profile_pic': 'https://pbs.twimg.com/profile_images/3594414986/196b1c7cfa898628c5476ea56bddc52b_normal.jpeg'},

    {"twitter": "https://twitter.com/https://twitter.com/KristenBaaby"},    # we aren't authorized to see this
    {"twitter": "https://twitter.com/&#8450;all Gi&#8477;&#8466;s in De&#8466;hi | 99IIII2O5I | Dhaula Kuan &#8495;sc"}, # this should totally fail

    {"twitter": "https://twitter.com/www.twitter.com/theirishmartini", "tweets": [], "twitter_followers": [], "twitter_friends": [],
     "twitter_name":"Irish Martini",
     'twitter_background_pic': 'https://pbs.twimg.com/profile_background_images/438000397/the_best.jpg',
     "twitter_description": "",
     "twitter_profile_url":"http://www.facebook.com/IrishMartini",
     'twitter_profile_pic': 'https://pbs.twimg.com/profile_images/1815703410/bust_reasonably_small_normal.jpg'},

    {"twitter": "https://twitter.com/https://twitter.com/#!/nycasianescor", "tweets": [], "twitter_followers": [], "twitter_friends": [],
     'twitter_background_pic': 'https://abs.twimg.com/images/themes/theme1/bg.png',
     'twitter_description': "",
     'twitter_profile_pic': 'https://pbs.twimg.com/profile_images/1763442408/2_normal.jpg',
     'twitter_name': 'Kanya Song',
     'twitter_profile_url': 'http://nycasianescorts.com/'},

    {"twitter": "https://twitter.com/Call Girls In Delhi 9811539366 Escort Service In D"},

    {"twitter": "https://twitter.com/https://mobile.twitter.com/monicasweets85?original_referer=http%3A%2F%2Fstatic.parastorage.com%2Fser",
     "tweets": [], "twitter_followers": [], "twitter_friends": [],
     'twitter_background_pic': 'https://pbs.twimg.com/profile_background_images/378800000081264488/37cb1e873bccb0320b45e910c4e446ae.jpeg',
     'twitter_description': "",
     'twitter_profile_pic': 'https://pbs.twimg.com/profile_images/378800000500705053/b14fd2e6ce8427be1f408afae7d4f306_normal.jpeg',
     'twitter_name': 'Monica Sweets',
     'twitter_profile_url': 'http://sweet85monica.wix.com/vipms'},

    {"twitter": """/phone/561-502-8304" title="Escort Phone 561-502-8304" class="fontyellow">561-502-8304</a></h2>

<div class="margintop5"></div>
For the mature gentlemen who loves paradise away for an hour . For more information call Paris @ 5615028304 in"""}
]

class TestTwitter(Runner):
    @classmethod
    def setup_class(cls):
        super(TestTwitter, cls).setup_class(Twitter, _test_nodes)

    def test_run(self):
        for t in super(TestTwitter, self).test_run(_expected_results, ignore_fields=set(["tweets", "twitter_followers", "twitter_friends", "twitter_description"])):
            yield t
