OWNER=graphql
IMAGE_NAME=auth-proxy
QNAME=$(OWNER)/$(IMAGE_NAME)

GIT_TAG=$(QNAME):$(TRAVIS_COMMIT)
BUILD_TAG=$(QNAME):$(TRAVIS_BUILD_NUMBER).$(TRAVIS_COMMIT)
TAG=$(QNAME):`echo $(TRAVIS_BRANCH) | sed 's/master/latest/;s/develop/unstable/'`

test:
	# DATABASE_URL=sqlite3://test.db $(IMAGE_NAME) server -p 8005
	DATABASE_URL="mysql://test:test@tcp(localhost:3306)/test?parseTime=true" godog
