DOCKERNAME=mrecco/mysqllog-collector
DOCKERTAG=v1.0.1
DOCKERFILE=./Dockerfile

build:
	@docker build -f $(DOCKERFILE) -t $(DOCKERNAME):$(DOCKERTAG) .

push:
	@docker push $(DOCKERNAME):$(DOCKERTAG)

rmi:
	@docker rmi $(DOCKERNAME):$(DOCKERTAG)

run: build
	# @docker run --rm -it --name mysqllog-collector-debug -v $(PWD):/opt/mysqllog:ro $(DOCKERNAME):$(DOCKERTAG) -general /opt/mysqllog/example.general.log -slowquery /opt/mysqllog/example.slowquery.log
	@docker run --rm -it --name mysqllog-collector-debug -v $(PWD):/opt/mysqllog:ro $(DOCKERNAME):$(DOCKERTAG) -slowquery /opt/mysqllog/example.slowquery.log
