import py2neo
from . import enhancer

import os
import time
import random
import logging
import collections

log = logging.getLogger("link.plugins.image_linker")

#NUM_WORKERS = 8

def listify(x):
    return x if isinstance(x, list) else [x]

class ImageLinker(enhancer.SQLEnhancer):
    def __init__(self):
        db_name = os.environ.get("SQL_DB", "memex_ht")
        super(ImageLinker, self).__init__(db=db_name)
        # this cypher statement looks up all entities associated with
        # the current node and adds a relationship for each image in
        # the node to those entities
        #
        # with the merge keyword, any existing BY_IMG links are set to
        # a min_dist of 0. this usually happens when two Ads share a
        # phone number and similar images. the first Ad may share a
        # similar image to the second Ad, causing the edge to be created
        # with a minimum distance of, say, 0.1. then, the second Ad
        # links to the entity via the same image and sets the minimum
        # distance to 0.
        #
        # this query ensures that only one BY_IMG edge will exist
        # between an ad and its entity
        self.__add_self = """
        MATCH (e:Entity)-[:BY_PHONE]->(n:Ad {id:{id}})
        MERGE (e)-[r:BY_IMG]->(n)
        SET r.user='auto', r.reason='belongs-to', r.min_dist=0.0
        """

        # this cypher statement looks up all entities associated with
        # the matching ad and creates a relationship between those
        # entities and the origin
        #
        # if the new relationship matches any existing ones, then we
        # only keep the one with the smallest distance score.
        self.__add_similar = """
        MATCH (e:Entity)-[:BY_PHONE]->(:Ad {id:{id}}), (n:Ad)
        WHERE n.id = {origin}
        MERGE (e)-[r:BY_IMG]->(n)
        ON CREATE SET r.user='auto', r.reason='contains-similar', r.min_dist=toFloat({distance})
        ON MATCH SET r.user='auto', r.reason='contains-similar', r.min_dist=(
            CASE
                WHEN toFloat(r.min_dist) > toFloat({distance}) THEN toFloat({distance})
                ELSE toFloat(r.min_dist)
            END)
        """

    def find_images(self):
        cursor = self.connection.cursor()
        cursor.execute(find_images, ids)
        # list of similar images
        for matched_image, score in cursor:
            similar_images[matched_image].append(score)
        cursor.close()

    def enhance(self, node):
        if 'image_locations' in node and 'image_ids' in node:
            log.info("==> Handling record [%s]" % node['id'])

            #locations = listify(node['image_locations'])
            ids = listify(node['image_ids'])
            find_images = "SELECT match_id, score FROM images_sim WHERE images_id in (%s)"
            find_ads = "SELECT id, ads_id FROM images WHERE id in (%s)"

            t0 = time.time()
            if ids:
                similar_images = collections.defaultdict(list)
                self.do_query(find_images % ",".join(map(str,ids)),
                    callback = lambda match_id, score: similar_images[match_id].append(score))
                log.info("    Ad [%s] had %d images. Found %d similar ones." % (node['id'], len(ids), len(similar_images)))

                if similar_images:
                    # similar images now contains a dictionary with a list of
                    # all images similar to any image in the current ad
                    #
                    # we're going to find all the ads that contain these images
                    similar_ads = collections.defaultdict(list)

                    # handle a returned row
                    def handle_row(image_id, ad_id):
                        # if the ad_id is null or if it matches the current
                        # node id we skip it. we don't want false links
                        # between a node an entity based on similar images
                        # within the ad.
                        if ad_id is not None and ad_id != node['id']:
                            min_score = min(similar_images[image_id])
                            similar_ads[ad_id].append(min_score)
                    self.do_query(find_ads % ",".join(map(str,similar_images.keys())),
                        callback=handle_row)
                    log.info("    Ad [%s] had %d similar ads; linked by similar images." % (node['id'], len(similar_ads)))

                    tx = self.db.cypher.begin()

                    for ad_id in similar_ads:
                        parameters = {"id": ad_id,
                            "origin": node['id'],
                            "distance": min(similar_ads[ad_id])
                        }
                        tx.append(self.__add_similar, parameters)

                    # add yourself
                    parameters = {"id": node['id']}
                    tx.append(self.__add_self, parameters)

                    # this batches up all cypher queries and executes them here
                    tx.commit()

            t1 = time.time()
            log.info("<== Handled a record with %d images in %f seconds" % (len(ids), t1 - t0))
