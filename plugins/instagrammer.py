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

log = logging.getLogger("link.plugins.twitter")

class Instagram(enhancer.Enhancer):
	def __init__(self):
		super(Instagram, self).__init__()
		self.INSTAGRAM_CLIENT_ID = os.environ['INSTAGRAM_CLIENT_ID']
		self.INSTAGRAM_CLIENT_SECRET = os.environ['INSTAGRAM_CLIENT_SECRET']
		self.INSTAGRAM_ACCESS_TOKEN = os.environ['INSTAGRAM_ACCESS_TOKEN']
		self.api = instagram.client.InstagramAPI(client_id=self.INSTAGRAM_CLIENT_ID,client_secret=self.INSTAGRAM_CLIENT_SECRET,access_token=self.INSTAGRAM_ACCESS_TOKEN)

	def get_instagram_username(self,node):
		if 'text' in node:
			text = node['text']
			expression = re.compile('instagram\s*:?\s*[^\s]*',re.IGNORECASE)
			try:
				raw = re.findall(expression,text)[0]
				instagram_user_name = None
				if '/' in raw:
					instagram_user_name = raw.split('/')[-1].strip()
				elif ':' in raw:
					instagram_user_name = raw.split(':')[-1].strip()
				return instagram_user_name
			except Exception as e:
				log.error(e)

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
		try:
			ig_username = self.get_instagram_username(node)
			
			if ig_username is not None:

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
						
					except Exception as e:
						exc_type, exc_value, exc_traceback = sys.exc_info()
						print exc_type
						print exc_value
						traceback.print_tb(exc_traceback, limit=10, file=sys.stdout)

				
		except:
			exc_type, exc_value, exc_traceback = sys.exc_info()
			traceback.print_tb(exc_traceback, limit=10, file=sys.stdout)
