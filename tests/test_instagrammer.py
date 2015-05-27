from plugins.instagrammer import Instagram
from . utils import Runner
import textwrap

class TestInstagrammer(Runner):
    ENHANCER = Instagram
    INPUT = [
        {"text": "Follow me on INSTAGRAM: http://instagram.com/VIPVeronica<br />"},
        {"text": "Follow us on INSTAGRAM;  nyc_gentlemensclub"},
        {"text": "Instagram name: BBWbiancalatinadreamgirl"},
        {"text": "<b>Follow us</b> on INSTAGRAM; <b> nyc_gentlemensclub</b><br>"},
        {"text": "INSTAGRAM:aleo_lover4202010"},
        {"text": "I'm 100% real and my Instagram is SUPERBRIANNAFREAKY"},
        {"text": "TWITTER@VIPMSE & INSTAGRAM@PORTIASBESTCHOICE"},
        {"text": '<a rel="nofollow" target="_blank" href="http://instagram.com/icandycane">Instagram.com/icandycane</a>'},
        {"text": textwrap.dedent('''CARTY 7083208795 OR PROOF VIDEO on Insta.gram as QUEENOSIRIS
            <br><br>
            Still think I'm too beautiful to be real?? <a target="_blank" href="http://www.instagram.com/QUEENOSIRIS"> PICS AND PROOF VIDEO</a>''')},
        {"text": "I don't actually have an instagram username."}
    ]
    OUTPUT = [
        {"text": "Follow me on INSTAGRAM: http://instagram.com/VIPVeronica<br />",
            "escaped_text": "Follow me on INSTAGRAM: http://instagram.com/VIPVeronica",
            'instagram': 'VIPVeronica',
            'instagram_profile_picture': 'https://igcdn-photos-a-a.akamaihd.net/hphotos-ak-xpf1/t51.2885-19/11084657_387545114783112_1103517643_a.jpg',
            "instagram_followers": "",
            "instagram_follows": "",
            "instagram_tags": "",
            "instagram_likers": "",
            "get_media_ids_and_posttimes": "",
            "get_commentors": ""
        },
        {"text": "Follow us on INSTAGRAM;  nyc_gentlemensclub",
            "escaped_text": "Follow us on INSTAGRAM;  nyc_gentlemensclub"},
        {"text": "Instagram name: BBWbiancalatinadreamgirl",
            "escaped_text": "Instagram name: BBWbiancalatinadreamgirl"},
        {"text": "<b>Follow us</b> on INSTAGRAM; <b> nyc_gentlemensclub</b><br>",
            "escaped_text": "Follow us on INSTAGRAM;  nyc_gentlemensclub"},
        {"text": "INSTAGRAM:aleo_lover4202010",
            "escaped_text": "INSTAGRAM:aleo_lover4202010",
            "instagram": "aleo_lover4202010"},
        {"text": "I'm 100% real and my Instagram is SUPERBRIANNAFREAKY",
            "escaped_text": "I'm 100% real and my Instagram is SUPERBRIANNAFREAKY"},
        {"text": "TWITTER@VIPMSE & INSTAGRAM@PORTIASBESTCHOICE",
            "escaped_text": 'TWITTER@VIPMSE & INSTAGRAM@PORTIASBESTCHOICE'},
        {"text": '<a rel="nofollow" target="_blank" href="http://instagram.com/icandycane">Instagram.com/icandycane</a>',
            "escaped_text": " http://instagram.com/icandycane Instagram.com/icandycane",
            "instagram": "icandycane",
            "instagram_error_message": "(400) APINotAllowedError-you cannot view this resource"},
        {"text": textwrap.dedent('''CARTY 7083208795 OR PROOF VIDEO on Insta.gram as QUEENOSIRIS
            <br><br>
            Still think I'm too beautiful to be real?? <a target="_blank" href="http://www.instagram.com/QUEENOSIRIS"> PICS AND PROOF VIDEO</a>'''),
            "escaped_text": "CARTY 7083208795 OR PROOF VIDEO on Insta.gram as QUEENOSIRIS\n            \n            Still think I'm too beautiful to be real??  http://www.instagram.com/QUEENOSIRIS  PICS AND PROOF VIDEO",
            "instagram": "QUEENOSIRIS",
            'instagram_profile_picture': 'https://instagramimages-a.akamaihd.net/profiles/profile_308829295_75sq_1387802497.jpg',
            "instagram_followers": "",
            "instagram_follows": "",
            "instagram_tags": "",
            "instagram_likers": "",
            "get_media_ids_and_posttimes": "",
            "get_commentors": ""},
        {"text": "I don't actually have an instagram username.",
         "escaped_text": "I don't actually have an instagram username."}
    ]
    IGNORE = set(["instagram_followers", "instagram_follows", "instagram_tags", "instagram_likers", "get_media_ids_and_posttimes", "get_commentors"])
