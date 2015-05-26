import tweepy
from py2neo import Graph, Node, Relationship, watch, GraphError
from . import enhancer
import os
import logging
import time
import urlparse

log = logging.getLogger("link.plugins.twitter")

def is_rate_limited(error):
	return isinstance(error.message, list) and \
		len(error.message) == 1 and \
		isinstance(error.message, dict) and \
		error.message[0]['code'] == 88

def with_retries(f, *args, **kwargs):
	attempts = 0
	num_backoff = kwargs.pop('num_backoff', 10)
	result = None
	while result is None:
		try:
			result = f(*args, **kwargs)
		except tweepy.error.TweepError as e:
			if is_rate_limited(e) and attempts < num_backoff:
				attempts += 1
				log.debug("Retrying, attempt %d out of %d" % (attempts, num_backoff))
				time.sleep(60*attempts)
			else:
				log.error(e)
				# usually, the user is deleted or we don't have access
				break
	return result

class Twitter(enhancer.Enhancer):
	def __init__(self):
		super(Twitter, self).__init__()

		TWITTER_CONSUMER_KEY = os.environ['TWITTER_CONSUMER_KEY']
		TWITTER_CONSUMER_KEY_SECRET = os.environ['TWITTER_CONSUMER_SECRET']
		TWITTER_ACCESS_TOKEN = os.environ['TWITTER_ACCESS_TOKEN']
		TWITTER_ACCESS_TOKEN_SECRET = os.environ['TWITTER_ACCESS_TOKEN_SECRET']

		auth = tweepy.OAuthHandler(TWITTER_CONSUMER_KEY, TWITTER_CONSUMER_KEY_SECRET)
		auth.set_access_token(TWITTER_ACCESS_TOKEN, TWITTER_ACCESS_TOKEN_SECRET)
		self.api = tweepy.API(auth)


	def enhance(self, node, num_backoff=10):
		if 'twitter' in node:
			# make sure it's a list
			identities = node['twitter'] if isinstance(node['twitter'], list) else [node['twitter']]

			# node['twitter'] contains a *list* of possible twitter handles
			for identity in identities:
				# twitter entries look like "http://twitter.com/foobar".
				# so we split along "/" and take the last one
				identity = identity.split("/")[-1]

				# some "identities" then look like url queries, so we
				# parse them as 'urls' and only keep the path
				identity = urlparse.urlparse(identity).path

				# some twitter handles begin with '@', so we strip them
				# this might not be a good idea since the scraper may
				# have just picked up a retweet
				if "@" == identity[0]:
					identity = identity[1::]
				log.debug(identity)

				user = with_retries(self.api.get_user, identity, num_backoff=num_backoff)
				if user is None:
					# we abort the attempt to get twitter data
					continue

				log.debug(user)

				tweets = with_retries(self.api.user_timeline, identity, num_backoff=num_backoff)
				if tweets is None:
					continue

				friends = with_retries(lambda x: [f.screen_name for f in x.friends()], user, num_backoff=num_backoff)
				if friends is None:
					continue

				followers = with_retries(lambda x: [f.screen_name for f in x.followers()], user, num_backoff=num_backoff)
				if followers is None:
					continue

				node['twitter_friends'] = friends
				node['twitter_followers'] = followers
				node['twitter_background_pic'] = user.profile_background_image_url_https
				node['twitter_profile_pic'] = user.profile_image_url_https
				node['twitter_description'] = user.description
				node['twitter_name'] = user.name
				node['twitter_profile_url'] = user.entities['url']['urls'][0]['expanded_url']
				node['tweets'] = [f.text for f in tweets]
