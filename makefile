## You can change these values
RELEASE_NAME_DB=db
RELEASE_NAME_MUTATOR=mutator
RELEASE_NAME_APP=app

.PHONY : all-darwin
all-darwin:
	make minikube update-darwin docker-integrity-sum update-patch docker-mutator

.PHONY : integrity-sum
docker-integrity-sum:
	make -C integrity-sum docker

.PHONY : minikube
minikube:
	minikube start

.PHONY : docker-mutator
docker-mutator:
	make -C integrity-mutator docker

.PHONY: update-patch
update-patch:
	cp patch-json-command.json integrity-mutator/

.PHONY: helm-all
helm-all:	helm-database helm-mutator helm-demo

.PHONY: helm-database
helm-database:
	helm dependency update helm-charts/database-to-integrity-sum
	helm install ${RELEASE_NAME_DB} helm-charts/database-to-integrity-sum

.PHONY: helm-mutator
helm-mutator:
	helm install ${RELEASE_NAME_MUTATOR} helm-charts/mutator

.PHONY: helm-demo
helm-demo:
	helm install ${RELEASE_NAME_APP} helm-charts/demo-apps-to-monitor

SECRET_DB="$$(grep 'secretName' helm-charts/database-to-integrity-sum/values.yaml | cut -d':' -f2 | tr -d '[:space:]')"
SECRET_HASHER="$$(grep 'secretNameDB' helm-charts/demo-apps-to-monitor/values.yaml | cut -d':' -f2 | tr -d '[:space:]')"
VALUE_RELEASE_NAME_APP="$$(grep 'releaseNameDB' helm-charts/demo-apps-to-monitor/values.yaml | cut -d':' -f2 | tr -d '[:space:]')"
SECRET_PATCH_NAME=${RELEASE_NAME_DB}-${SECRET_DB}
A="$$(grep 'secretRef' patch-json-command.json | cut -d'{' -f2 | tr -d '[:space:]')"
.PHONY: update-darwin
update-darwin:
#	echo ${SECRET_PATCH_NAME}
	echo $(A)
#	sed -i '' "s/${SECRET_HASHER}/${SECRET_DB}/" helm-charts/demo-apps-to-monitor/values.yaml >> helm-charts/demo-apps-to-monitor/values.yaml
#	sed	-i '' "s/${VALUE_RELEASE_NAME_APP}/${RELEASE_NAME_DB}/" helm-charts/demo-apps-to-monitor/values.yaml >> helm-charts/demo-apps-to-monitor/values.yaml
