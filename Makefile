IMAGE_EXPORT := $(shell make -s -C integrity-sum -f Makefile print-var-image-name)
ALG=$(shell make -s -C integrity-sum -f Makefile print-var-alg)
DIRS="app,bin"	# You can change these values

.PHONY : all
all: minikube start-minio-integrity-crd snapshots start-mutator helm-demo
	@echo "==> Successfully installed all systems"

.PHONY : start-minio-integrity-crd
start-minio-integrity-crd:
	make -C ./integrity-sum minio-install buildtools build docker crd-controller-build load-images crd-controller-deploy

.PHONY : snapshots
snapshots:
	make -C ./integrity-sum IMAGE_EXPORT=$(IMAGE_EXPORT) ALG=$(ALG) DIRS=$(DIRS) export-fs snapshot helm-snapshot

.PHONY : start-mutator
start-mutator:
	make -C ./integrity-mutator docker helm-mutator
	@echo "==> Successfully installed mutator"

.PHONY : minikube
minikube:
	minikube start

.PHONY: helm-demo
helm-demo:
	helm install ${RELEASE_NAME_APP} helm-charts/demo-apps-to-monitor
	@echo "==> Successfully installed demo-apps"
