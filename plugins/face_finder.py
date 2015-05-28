import cv2
import sys
import numpy as np
import urllib
import urlparse
import re
import os
import csv
from . import enhancer
import requests
import logging
import os
import sys
import traceback
import tinys3
import tempfile

log = logging.getLogger("link.plugins.face_finder")

def get_filepath(image_url):
	path = urlparse.urlparse(image_url).path
	basename = os.path.basename(path)
	return os.path.splitext(basename)[0]

class FaceFinder (enhancer.Enhancer):
	def __init__(self):
		super(FaceFinder, self).__init__()
		self.image_directory = tempfile.mkdtemp()
		self.cascPath = os.getcwd() + "/plugins" + "/haarcascade_frontalface_default.xml"
		log.debug("Haar cascade path: '{}'".format(self.cascPath))
		self.faceCascade = cv2.CascadeClassifier(self.cascPath)
		log.debug("Cascade classifier loaded")

		S3_ACCESS_KEY=os.environ['S3_ACCESS_KEY']
		S3_SECRET_KEY=os.environ['S3_SECRET_KEY']
		self.conn = tinys3.Connection(S3_ACCESS_KEY,
			S3_SECRET_KEY,
			tls=True,
			endpoint='s3-us-west-1.amazonaws.com')

		uploaded_urls=[]
		log.debug("Caching bucket contents.")
		bucket_contents=self.conn.list("","memexadvertisements") #Get S3 bucket contents to avoid re-uploading images already in the database
		for files in bucket_contents:
			key = os.path.splitext(files["key"])[0]
			uploaded_urls.append(key)

		# convert the uploaded urls to a set for fast(er) inclusion checking
		self.uploaded_urls = set(uploaded_urls)
		log.debug("Bucket contents cached.")

	def load_image(self,imagePath):
		log.debug("Loading image")
		resp = urllib.urlopen(imagePath)
		log.debug(resp)
		image = np.asarray(bytearray(resp.read()), dtype="uint8")
		image = cv2.imdecode(image, cv2.IMREAD_COLOR)
		return image

	def find_faces(self,image):
		log.debug("Finding faces")
		gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
		faces = self.faceCascade.detectMultiScale(
			    gray,
			    scaleFactor=1.1,
			    minNeighbors=5,
			    minSize=(30, 30),
			    flags = cv2.cv.CV_HAAR_SCALE_IMAGE
			)
		log.debug(faces)
		return faces

	def crop_face(self,image,faces,bucketname, image_location):
		for (x, y, w, h) in faces:
			crop_img = image[y:y+h, x:x+w]
			filepath = get_filepath(image_location)
			cv2.imwrite(self.image_directory + "/" + filepath + ".png", crop_img)
			with open(self.image_directory + "/" + filepath + ".png",'rb') as f:
				self.conn.upload(filepath + ".png",f,bucketname)
		return 'https://s3-us-west-1.amazonaws.com/' + bucketname + "/" + filepath + ".png"

	def enhance(self, node):
		if 'image_locations' in node:
			n_faces=0

			node["face_image_url"] = []
			for identity in node['image_locations']:
				filepath = get_filepath(identity)
				try:
					image=self.load_image(identity)
					faces=self.find_faces(image)
					if len(faces) > 0:
						n_faces += 1
						if filepath not in self.uploaded_urls:
							face_url=self.crop_face(image,faces,"memexadvertisements", identity)
							node["face_image_url"].append( face_url )
							log.info("Loaded new face '{}' to S3".format(face_url))
						else:
							node["face_image_url"].append( 'https://s3-us-west-1.amazonaws.com/memexadvertisements/' + filepath + '.png')
				except Exception as e:
					exc_type, exc_value, exc_traceback = sys.exc_info()
					traceback.print_tb(exc_traceback, limit=10, file=sys.stdout)
					log.error(e)

			node["n_faces"]=n_faces
