OWNER=jakubknejzlik
IMAGE_NAME=godog-sql
QNAME=$(OWNER)/$(IMAGE_NAME)

GIT_TAG=$(QNAME):$(TRAVIS_COMMIT)
BUILD_TAG=$(QNAME):$(TRAVIS_BUILD_NUMBER).$(TRAVIS_COMMIT)
TAG=$(QNAME):`echo $(TRAVIS_BRANCH) | sed 's/master/latest/;s/develop/unstable/'`


lint:
	docker run -it --rm -v "$(PWD)/Dockerfile:/Dockerfile:ro" redcoolbeans/dockerlint

build:
	docker build -t $(GIT_TAG) .
	
tag:
	docker tag $(GIT_TAG) $(BUILD_TAG)
	docker tag $(GIT_TAG) $(TAG)
	
login:
	@docker login -u "$(DOCKER_USER)" -p "$(DOCKER_PASS)"
push: login
	# docker push $(GIT_TAG)
	# docker push $(BUILD_TAG)
	docker push $(TAG)

test:
	# DATABASE_URL=sqlite3://test.db $(IMAGE_NAME) server -p 8005
	DATABASE_URL="mysql://test:test@tcp(localhost:3306)/test?parseTime=true" godog
