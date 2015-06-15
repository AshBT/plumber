from py2neo import Graph, Node, Relationship, watch, GraphError
from . import enhancer
import os
import logging
import re
import requests
import sys
import traceback
import most_common_words_list

from nltk.util import ngrams
import hashlib
import csv

from htmlentitydefs import name2codepoint
# for some reason, python 2.5.2 doesn't have this one (apostrophe)
name2codepoint['#39'] = 39

log = logging.getLogger("link.plugins.twitter")

class TextHash(enhancer.Enhancer):
	def unescape(self,s):
	    "unescape HTML code refs; c.f. http://wiki.python.org/moin/EscapingHtml"
	    return re.sub('&(%s);' % '|'.join(name2codepoint),
	              lambda m: unichr(name2codepoint[m.group(1)]), s)
	    
	# Remove 100 most common words 
	def strip_common_words(self,string,common_words):
		return ' '.join(filter(lambda w: not w in common_words,string.split()))

	def clean_HTML(self,string):
		unescaped_string = self.unescape(string)
		unescaped_string = unescaped_string.encode("utf-8")
		clean_string = re.sub('(<.*?>|\n)',' ',unescaped_string)
		return clean_string
	
	def enhance(self, node):
		#most_common_words = self.get_most_common_words()
		most_common_words = most_common_words_list.most_common_words
		dirty_text = node['text']
		signature = ''
		#4, 6 gives good cross-phone number matching	
		for n in [4]:
			cleaned_text = self.strip_common_words(self.clean_HTML(dirty_text).lower(),most_common_words)
			if len(cleaned_text.split())>=4:
				n_grams = list(ngrams(cleaned_text.split(),n))
				minimum_hash = min([int(hashlib.sha256(' '.join(gram)).hexdigest(),16) for gram in n_grams])
				signature += str(minimum_hash)

		node['text_signature'] = signature
