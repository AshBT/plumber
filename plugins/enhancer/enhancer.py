import os
import logging
import traceback
from abc import ABCMeta, abstractmethod
from py2neo import Graph, Node, Relationship, watch, GraphError

from py2neo.packages.httpstream import http
http.socket_timeout = 9999

log = logging.getLogger("link.enhancer")

class Enhancer(object):
	__metaclass__ = ABCMeta

	def __init__(self, graph = None):
		NEO_USER = os.environ['NEO_USER']
		NEO_PASS = os.environ['NEO_PASS']
		NEO_HOST = os.environ.get('NEO_HOST', 'localhost:7474')

		self.db = Graph("http://{user}:{passwd}@{host}/db/data/".format(host=NEO_HOST, user=NEO_USER,passwd=NEO_PASS)) if graph is None else graph

	@abstractmethod
	def enhance(self, node):
		""" The enhance method takes a node as an argument and does
			anything with the node. The graph can be accessed with
			self.db. Furthermore, any changes to the node are
			pushed to the db after this function is called, so you do
			not need to invoke node.push() in your code.
		"""
		pass

	def cleanup(self):
		pass

	def run(self, start=0, stop=None, batch_size=100, skip_on_error=True):
		""" This grabs 100 nodes at a time and process them in memory
			before moving on to the next 100, until there are no more
			nodes left.
		"""
		try:
			counter = 0
			nodes = []
			stop_now = False
			while True:
				result = self.db.cypher.execute('MATCH (ad:Ad) RETURN ad SKIP %d LIMIT %d' % (counter + start, batch_size))
				for record in result:
					if counter % 10 == 0:
						log.info("Handled %d records." % counter)
					counter += 1

					node = record.ad

					log.debug("==> %s" % str(node))
					try:
						self.enhance(node)	# this might modify node
						# remove any empty lists--replace with empty string?
						# see py2neo issue: https://github.com/nigelsmall/py2neo/issues/395
						for p in node.properties:
							elem = node.properties[p]
							if isinstance(elem, list) and not elem:
								node.properties[p] = ""
						nodes.append(node)
						log.debug("<== %s" % str(node))
					except Exception as e:
						log.error("[Exception] Skipping record; caught an exception during handling: '%s'" % e)
						#log.error(traceback.format_exc())
						if skip_on_error:
							continue
						else:
							raise e
					finally:
						if stop is not None:
							if counter >= (stop - start):
								stop_now = True
								break
				try:					
					self.db.push(*nodes)
				except Exception as e:
					log.error("[Exception] '%s', pushing '%d' nodes serially" % (e, len(nodes)))
					# try to update serially
					for n in nodes:
						try:
							n.push()
						except Exception as e:
							log.error("[Exception] Skipping node: '%s'" % str(n))
							continue
				nodes = []

				if stop_now or len(result) < batch_size:
					# if the number of returned results is less than the
					# batch size, then we now we've exhausted all nodes
					log.info("Finished!")
					break

			log.info("[COMPLETED] Handled %d records." % counter)
		finally:
			self.cleanup()
