from plugins.youtuber import Youtube
from . utils import Runner

# from a subset of actual data




class TestYoutube(Runner):
    ENHANCER = Youtube
    INPUT = [
        {"youtube": ["http://www.youtube.com/embed/oLuZywblDns?rel=0&wmode=opaque"]}
    ]
    OUTPUT = [
        {"youtube": ["http://www.youtube.com/embed/oLuZywblDns?rel=0&wmode=opaque"],
        'youtube_user': "PHInestnupe",
         'youtube_video_ids': "",
         'youtube_video_titles': "",
         'youtube_video_dates': "",
         'youtube_video_urls': "",
         'youtube_video_durations': "",
         'youtube_video_view_counts': ""
    }]
    IGNORE = set([
        "youtube_video_view_counts",
        "youtube_video_durations",
        "youtube_video_ids",
        "youtube_video_titles",
        "youtube_video_dates",
        "youtube_video_urls"
    ])
