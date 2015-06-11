.PHONY: test clean
test:
	./scripts/test.sh

clean:
	rm -rf neo4j* test-env
