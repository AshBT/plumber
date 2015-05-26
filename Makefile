.PHONY: test clean
test:
	sh scripts/test.sh

clean:
	rm -rf neo4j* test-env
