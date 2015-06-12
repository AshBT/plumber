import instagram
from py2neo import Graph, Node, Relationship, watch, GraphError
from . import enhancer
import os
import logging
import re
import requests
import sys
import traceback
from compiler.ast import flatten
from HTMLParser import HTMLParser

log = logging.getLogger("link.plugins.instagram")

# strip all HTML tags from the text (except <a href="..."></a>)
# modified from http://stackoverflow.com/questions/753052/strip-html-from-strings-in-python
class HTMLStripper(HTMLParser):
	def __init__(self):
		self.reset()
		self.fed = []

	def handle_starttag(self, tag, attrs):
		attrs = dict(attrs)
		if tag == 'a' and 'href' in attrs:
			self.fed.append(" " + attrs['href'] + " ")

	def handle_data(self, d):
		self.fed.append(d)

	def get_data(self):
		return ''.join(self.fed)

class Instagram(enhancer.Enhancer):
	def __init__(self):
		super(Instagram, self).__init__()
		self.INSTAGRAM_CLIENT_ID = os.environ['INSTAGRAM_CLIENT_ID']
		self.INSTAGRAM_CLIENT_SECRET = os.environ['INSTAGRAM_CLIENT_SECRET']
		self.INSTAGRAM_ACCESS_TOKEN = os.environ['INSTAGRAM_ACCESS_TOKEN']
		self.api = instagram.client.InstagramAPI(client_id=self.INSTAGRAM_CLIENT_ID,client_secret=self.INSTAGRAM_CLIENT_SECRET,access_token=self.INSTAGRAM_ACCESS_TOKEN)

	def get_instagram_username(self,node):
		if 'text' in node:

			s = HTMLStripper()
			text = node['text']
			s.feed(text)
			text = s.get_data()
			node['escaped_text'] = text

			expression = re.compile('instagram\s*:?\s*[^\s]*',re.IGNORECASE)
			instagram_user_name = None
			try:
				regexp_match = re.findall(expression,text)
				if len(regexp_match) > 0:
					raw = regexp_match[0]
					if '/' in raw:
						instagram_user_name = raw.split('/')[-1].strip()
					elif ':' in raw:
						instagram_user_name = raw.split(':')[-1].strip()
			finally:
				return instagram_user_name
		return None

	def get_instagram_id(self,instagram_user_name):
		#It seems difficult to find the id from a username using instagramAPI - figure out if this is possible later
		instagram_request = requests.get('https://api.instagram.com/v1/users/search?access_token=' + str(self.INSTAGRAM_ACCESS_TOKEN)+'&q='+str(instagram_user_name))
		try:
			id_num = instagram_request.json()['data'][0]['id']
		except:
			return 0
		return id_num

	def get_profile_picture(self,instagram_id):
		profile_picture_request = requests.get('https://api.instagram.com/v1/users/'+str(instagram_id)+'?access_token=' + self.INSTAGRAM_ACCESS_TOKEN)
		try:
			profile_picture = profile_picture_request.json()['data']['profile_picture']
		except:
			return 0
		return profile_picture

	def get_recent_media(self,instagram_id):
		recent_media_request = requests.get('https://api.instagram.com/v1/users/'+str(instagram_id)+'/media/recent?access_token=' + self.INSTAGRAM_ACCESS_TOKEN)
		try:
			recent_media = recent_media_request.json()
		except:
			return 0
		return recent_media

	def get_likers(self,recent_media):
		return [[x['username'],x['id'],x['full_name'],x['profile_picture']] for datum in recent_media['data'] for x in datum['likes']['data']]

	def get_all_instagram_tags(self,recent_media):
		return list(set(flatten([x['tags'] for x in recent_media['data']])))

	def get_media_ids_and_posttimes(self,recent_media):
		return [[x['id'],x['created_time']] for x in recent_media['data']]

	def get_commentors(self,recent_media):
		commentors = []
		for media in recent_media['data']:
			try:
				media_request = requests.get('https://api.instagram.com/v1/media/'+str(media['id'])+'/comments?access_token=' + self.INSTAGRAM_ACCESS_TOKEN)
				media_json = media_request.json()
				commentors += [
				[datum['from']['username'],
				datum['from']['id'],
				datum['from']['full_name'],
				datum['from']['profile_picture'],
				datum['text']]
				for datum in media_json['data']]
			except:
				exc_type, exc_value, exc_traceback = sys.exc_info()
				print exc_type
				print exc_value
				traceback.print_tb(exc_traceback, limit=10, file=sys.stdout)
		return commentors

	def enhance(self, node):
		#print('hi')
		ig_username = self.get_instagram_username(node)

		# the first checks if the username is not None, the second
		# checks if the username is "blank" (empty string)
		if ig_username is not None and ig_username:
			node['instagram'] = ig_username
			ig_id = int(self.get_instagram_id(ig_username))

			if ig_id != 0:
				try:
					#ig_id = 372991597
					recent_media = self.get_recent_media(ig_id)
					node['instagram_followers'] = ';'.join([str(f)for f in self.api.user_followed_by(str(ig_id))[0]])
					node['instagram_follows'] = ';'.join([str(f)for f in self.api.user_follows(str(ig_id))[0]])
					node['instagram_tags'] = self.get_all_instagram_tags(recent_media)
					node['instagram_profile_picture'] = self.get_profile_picture(ig_id)

					node['instagram_likers'] = ','.join(flatten(self.get_likers(recent_media)))
					node['get_media_ids_and_posttimes'] = ','.join(flatten(self.get_media_ids_and_posttimes(recent_media)))
					node['get_commentors'] = ','.join(flatten(self.get_commentors(recent_media)))
				except instagram.InstagramAPIError as e:
					# make a note that the user exists but we aren't allowed to access info
					if e.error_type == "APINotAllowedError":
						node['instagram_error_message'] = str(e)
					else:
						raise e
				except Exception as e:
					exc_type, exc_value, exc_traceback = sys.exc_info()
					print exc_type
					print exc_value
					traceback.print_tb(exc_traceback, limit=10, file=sys.stdout)
					raise e
