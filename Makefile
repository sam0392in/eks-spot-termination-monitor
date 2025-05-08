.SILENT: build
.PHONY:	release build install uninstall

APPNAME = eks-spot-termination-monitor
REPO = sam0392in/eks-spot-termination-monitor
NAMESPACE = kube-system
PLATFORM = amd64
TAG ?= latest

deploy: build push release

build:
	docker build --platform=linux/$(PLATFORM) -t $(REPO):$(TAG) .

push:
	docker push $(REPO):$(TAG)

release:
	helm upgrade -i $(APPNAME) -n $(NAMESPACE) ./k8s -f ./k8s/values.yaml \
	--set-string image.tag=$(TAG)

dryrun:
	helm upgrade -i $(APPNAME) -n $(NAMESPACE) ./k8s -f ./k8s/values.yaml \
	--set-string image.tag=$(TAG) \
	--dry-run

uninstall:
	helm uninstall $(APPNAME) -n $(NAMESPACE)
