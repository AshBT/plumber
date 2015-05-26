from plugins.twitter import face_finder
from . utils import Runner

# from a subset of actual data
_test_nodes = [
    {"youtube": ["http://www.youtube.com/embed/oLuZywblDns?rel=0&wmode=opaque"]}

_expected_results = [
{node['youtube_user'] = "PHInestnupe",
        node['youtube_video_ids'] = ['ZWpE6fTEw24', 'qLQH1W-_sdc', 'Jp2WDiuJUVE', 'bLnXr-T-Jlg', 'EfdJLCV3b_c', 'Yp4w04gHdIw', 'Q9nMNJ2ehDM', 'i1nwMmEfFNQ', 'oLuZywblDns', '2etImOlRlB8', 'esOZWw0hkpg', 'gpsFheEkS94', 'cxiBSYNaY6Y']
        node['youtube_video_titles'] = ['Dynamic Duo ft. Mistah FAB & Hollywood Hov- "Pull ME Close"', 'Arthroscopic Labral Tear Debridement  and Decompression', 'NATOMAS SLIMM - NINER GOLD', "SHE GETT'IN IT", '"Love in Da Studio" Y.P. Ft. Stephanie Nicole SWAMP SoundZ Prod.', 'THIS FRIDAY FEB 15TH "KUPIDS KRUSH" RECAP PROMO', 'Sister Friends Casting Call Ad', 'Love in Da Studio "Y.P. Ft. Stephanie Nicole" by  SWAMP SOundZ', 'ADRIENNE AMOUR', 'Dynamic Duo- Pull Me Close (Grab My Pants Leg) behind the scene', 'CHOPPA ft. DANN-E and FOSTA CHILD- SNAKES', 'KUPIDS KRUSH VALENTINES BASH', 'NATOMAS SLIMM "ASS DROP" SNIPPET TEASER'],
        node['youtube_video_dates'] = ['2011-01-24T22:13:15.000Z', '2013-03-19T18:21:44.000Z', '2013-02-04T02:02:30.000Z', '2013-03-06T07:39:36.000Z', '2009-06-03T22:55:11.000Z', '2013-02-11T23:01:38.000Z', '2013-05-09T21:41:16.000Z', '2009-06-03T08:19:29.000Z', '2014-06-19T00:46:50.000Z', '2010-11-27T07:43:08.000Z', '2012-09-14T07:03:41.000Z', '2013-02-09T17:13:56.000Z', '2013-04-13T00:27:36.000Z'],
        node['youtube_video_urls'] = ['https://www.youtube.com/watch?v=ZWpE6fTEw24', 'https://www.youtube.com/watch?v=qLQH1W-_sdc', 'https://www.youtube.com/watch?v=Jp2WDiuJUVE', 'https://www.youtube.com/watch?v=bLnXr-T-Jlg', 'https://www.youtube.com/watch?v=EfdJLCV3b_c', 'https://www.youtube.com/watch?v=Yp4w04gHdIw', 'https://www.youtube.com/watch?v=Q9nMNJ2ehDM', 'https://www.youtube.com/watch?v=i1nwMmEfFNQ', 'https://www.youtube.com/watch?v=oLuZywblDns', 'https://www.youtube.com/watch?v=2etImOlRlB8', 'https://www.youtube.com/watch?v=esOZWw0hkpg', 'https://www.youtube.com/watch?v=gpsFheEkS94', 'https://www.youtube.com/watch?v=cxiBSYNaY6Y'],
        node['youtube_video_durations'] = ['PT4M45S', 'PT10M22S', 'PT3M15S', 'PT34S', 'PT3M46S', 'PT1M44S', 'PT1M41S', 'PT3M29S', 'PT1M58S', 'PT4M58S', 'PT4M57S', 'PT2M28S', 'PT36S'],
        node['youtube_video_view_counts'] = ['145', '1777', '172', '160', '357', '294', '170', '396', '10550', '636', '606', '226', '502']
        }]






        











class Test_face_finder(Runner):
    @classmethod
    def setup_class(cls):
        super(Test_Youtube, cls).setup_class(face_finder, _test_nodes)

    def test_run(self):
        for t in super(Test_Youtube, self).test_run(_expected_results):
            yield t



            