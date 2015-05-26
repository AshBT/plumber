#!/usr/bin/env bash
set -e

export NEO4J_TEST_DIRECTORY="neo4j-test"
if [ -d $NEO4J_TEST_DIRECTORY ];
then
    echo "Using existing neo4j test directory."
else
    echo "Downloading neo4j."
    # in case someone ever uses windows, we use py2neo to print the
    # proper distribution name for our os
    ARCHIVE_NAME=`python -c "import py2neo.server; print py2neo.server.dist_archive_name('community', '2.2.1')"`
    DIST_NAME=`python -c "import py2neo.server; print py2neo.server.dist_name('community', '2.2.1')"`
    curl http://dist.neo4j.org/$ARCHIVE_NAME -o neo4j.tar.gz
    tar -xvzf neo4j.tar.gz
    mv $DIST_NAME $NEO4J_TEST_DIRECTORY
fi

# check if the user has the neo4j binary, if so, attempt to stop any
# existing servers
status=""
if hash neo4j 2>/dev/null
then
  echo "Stopping any existing neo4j servers"
  status=$(neo4j stop)
fi

echo "Installing virtualenv..."
virtualenv --system-site-packages test-env
source test-env/bin/activate
{
  pip install nose coverage && \
  pip install -r requirements.txt && \
  NEO_PASS=password python `which nosetests` --with-coverage --cover-package=plugins --logging-filter=link
}

echo "Deactivating test environment"
deactivate
if [[ "$status" =~ "Stopping" ]];
then
  echo "Restarting stopped neo4j instance"
  neo4j start
fi
