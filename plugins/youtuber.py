from apiclient.discovery import build
from apiclient.errors import HttpError
from oauth2client.tools import argparser
from bs4 import BeautifulSoup
from . import enhancer
import requests
import logging
import os
import sys
import traceback

log = logging.getLogger("link.plugins.youtube")

class Video:
	def __init__(self,vid_id,vid_title,vid_date,vid_url):
		self.vid_id = vid_id
		self.vid_title = vid_title
		self.vid_date = vid_date
		self.vid_url=vid_url
		
	def add_characteristics(self,youtube):
		try:
			vid_details = youtube.videos().list(
							id=self.vid_id,
							part="snippet, contentDetails, statistics").execute()
		except:
			exc_type, exc_value, exc_traceback = sys.exc_info()
			traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
		
		search_result =vid_details.get("items", [])[0]
		self.vid_duration = search_result["contentDetails"]["duration"]
		self.vid_view_count = search_result["statistics"]["viewCount"]
		
class Youtube(enhancer.Enhancer):

	def __init__(self, **kwargs):
		super(Youtube, self).__init__(**kwargs)
		YOUTUBE_DEVELOPER_KEY = os.environ['YOUTUBE_DEVELOPER_KEY']
		YOUTUBE_API_SERVICE_NAME = "youtube"
		YOUTUBE_API_VERSION = "v3"

		self.youtube = build(YOUTUBE_API_SERVICE_NAME, YOUTUBE_API_VERSION,
			developerKey=YOUTUBE_DEVELOPER_KEY)

	def get_username_from_video(self,identity):
		seed_html=requests.get(identity).content #identity represents embedded Youtube videos
		seedsoup=BeautifulSoup(seed_html)
		youtubelink=seedsoup.find_all("link")[0].get("href")
		vid_search = self.youtube.search().list(
			q=youtubelink,
			part="id,snippet",
			maxResults=1
			).execute()

		user = vid_search.get("items", [])[0]["snippet"]["channelTitle"]
		return user

	def append_attributes_to_node(self,node,user,video_list):
		node['youtube_user'] = user
		node['youtube_video_ids'] = [V.vid_id for V in video_list]
		node['youtube_video_titles'] = [V.vid_title for V in video_list]
		node['youtube_video_dates'] = [V.vid_date for V in video_list]
		node['youtube_video_urls'] = [V.vid_url for V in video_list]
		node['youtube_video_durations'] = [V.vid_duration for V in video_list]
		node['youtube_video_view_counts'] = [V.vid_view_count for V in video_list]


	def enhance(self, node):
		if 'youtube' in node:
			for identity in node['youtube']: 
				try:
					user = self.get_username_from_video(identity)
					
					user_search = self.youtube.search().list( #Get all the videos for this username
							q=user,
							part="id,snippet",
							maxResults=50
							).execute()

					# Get all the videos using a list comprehension
					video_list = [
					Video(search_result["id"]["videoId"],
								search_result["snippet"]["title"],
								search_result["snippet"]["publishedAt"],
								"https://www.youtube.com/watch?v="+search_result["id"]["videoId"])
								for search_result in user_search.get("items", []) 
								if search_result["id"]["kind"] == "youtube#video"]

					#Get additional details about the video
					for V in video_list:
						V.add_characteristics(self.youtube)

				except:
					exc_type, exc_value, exc_traceback = sys.exc_info()
					traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)

				try:
					self.append_attributes_to_node(node,user,video_list)

				except:
					exc_type, exc_value, exc_traceback = sys.exc_info()
					traceback.print_tb(exc_traceback, limit=1, file=sys.stdout)
					print sys.exc_info()

