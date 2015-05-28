import cv2
import sys
import numpy as np
import urllib
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

log = logging.getLogger("link.plugins.face_finder")

class FaceFinder (enhancer.Enhancer):
	def __init__(self):
		super(FaceFinder, self).__init__()
		self.directory=os.getcwd()
		self.cascPath = os.getcwd() + "/plugins" + "/haarcascade_frontalface_default.xml"
		log.info("Haar cascade path: '{}'".format(self.cascPath))
		self.faceCascade = cv2.CascadeClassifier(self.cascPath)
		log.info("Cascade classifier loaded")

		S3_ACCESS_KEY=os.environ['S3_ACCESS_KEY']
		S3_SECRET_KEY=os.environ['S3_SECRET_KEY']
		self.conn = tinys3.Connection(S3_ACCESS_KEY,
			S3_SECRET_KEY,
			tls=True,
			endpoint='s3-us-west-1.amazonaws.com')

	def load_image(self,imagePath):
		resp = urllib.urlopen(imagePath)
		image = np.asarray(bytearray(resp.read()), dtype="uint8")
		image = cv2.imdecode(image, cv2.IMREAD_COLOR)
		return image

	def find_faces(self,image):
		gray = cv2.cvtColor(image, cv2.COLOR_BGR2GRAY)
		faces = self.faceCascade.detectMultiScale(
			    gray,
			    scaleFactor=1.1,
			    minNeighbors=5,
			    minSize=(30, 30),
			    flags = cv2.cv.CV_HAAR_SCALE_IMAGE
			)
		return faces

	def crop_face(self,image,faces,bucketname, image_location):
		for (x, y, w, h) in faces:
			crop_img = image[y:y+h, x:x+w]
			filepath=re.sub("httpss3amazonawscomroxyimages","",re.sub("[\W_]+","",str(image_location)))
			cv2.imwrite(self.directory + "/Plugins/Images/" + filepath + ".png",crop_img)
			f = open(self.directory + "/Plugins/Images/" + filepath + ".png",'rb')
			self.conn.upload(filepath + ".png",f,bucketname)
		return 'https://s3-us-west-1.amazonaws.com/' + bucketname + "/" + filepath + ".png"

	def enhance(self, node):
		if 'image_locations' in node:
			n_faces=0
			uploaded_urls=[]
			bucket_contents=self.conn.list("","memexadvertisements") #Get S3 bucket contents to avoid re-uploading images already in the database
			for files in bucket_contents:
				uploaded_urls.append(files["key"])

			node["face_image_url"] = []
			for identity in node['image_locations']:
				filepath=re.sub("httpss3amazonawscomroxyimages","",re.sub("[\W_]+","",str(identity))) + ".png"
				try:
					image=self.load_image(identity)
					faces=self.find_faces(image)
					if len(faces) > 0:
						n_faces = n_faces+1
						if filepath not in uploaded_urls:
							face_url=self.crop_face(image,faces,"memexadvertisements", identity)
							node["face_image_url"].append( face_url )
							print "Loaded New Face to S3"
						else:
							node["face_image_url"].append( 'https://s3-us-west-1.amazonaws.com/memexadvertisements/' + filepath )
				except Exception as e:
					exc_type, exc_value, exc_traceback = sys.exc_info()
					traceback.print_tb(exc_traceback, limit=10, file=sys.stdout)
					raise e

			node["n_faces"]=n_faces
