import shutil
import py2neo.server
import tempfile
import logging
import tarfile
import os
import subprocess
import requests
import warnings

__TEMP_DIR = tempfile.mkdtemp()
__NEO4J_TEST_DIRECTORY = os.getenv('NEO4J_TEST_DIRECTORY')
__NEO4J_SERVER = None

def setup_package():
    global __NEO4J_SERVER
    print "Copying server"
    shutil.copytree(__NEO4J_TEST_DIRECTORY, "%s/neo4j" % __TEMP_DIR)
    print "Running server"
    try:
        subprocess.check_output(["%s/neo4j/bin/neo4j" % __TEMP_DIR, "start"])
    except subprocess.CalledProcessError as e:
        print "Make sure an existing neo4j process is not running."
        teardown_package()
        raise e
    # change the password the first time
    payload = {"password": "password"}
    r = requests.post("http://localhost:7474/user/neo4j/password", auth=("neo4j", "neo4j"), data=payload)
    # ignore all warnings
    warnings.simplefilter("ignore")

def teardown_package():
    subprocess.check_output(["%s/neo4j/bin/neo4j" % __TEMP_DIR, "stop"])
    print "Removing temp directory"
    shutil.rmtree(__TEMP_DIR)
    warnings.resetwarnings()
