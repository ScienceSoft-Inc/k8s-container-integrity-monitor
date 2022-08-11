## You can change these values
RELEASE_NAME_DB=db
RELEASE_NAME_MUTATOR=mutator
RELEASE_NAME_APP=app
TIMEOUT=30

.PHONY : all-darwin
all-darwin:
	make minikube update-darwin docker-integrity-sum update-patch docker-mutator helm-all

.PHONY : integrity-sum
docker-integrity-sum:
	make -C integrity-sum docker

.PHONY : minikube
minikube:
	minikube start

.PHONY : stop
stop:
	helm uninstall  ${RELEASE_NAME_APP}
	helm uninstall  ${RELEASE_NAME_MUTATOR}
	helm uninstall  ${RELEASE_NAME_DB}
	minikube stop

.PHONY : docker-mutator
docker-mutator: update-patch
	make -C integrity-mutator docker

.PHONY: update-patch
update-patch:
	cp patch-json-command.json integrity-mutator/

.PHONY: helm-all
helm-all:	helm-database helm-mutator timeout helm-demo

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

.PHONE: timeout
timeout:
	sleep ${TIMEOUT}

SECRET_DB="$$(grep 'secretName' helm-charts/database-to-integrity-sum/values.yaml | cut -d':' -f2 | tr -d '[:space:]')"
SECRET_HASHER="$$(grep 'secretNameDB' helm-charts/demo-apps-to-monitor/values.yaml | cut -d':' -f2 | tr -d '[:space:]')"
VALUE_RELEASE_NAME_APP="$$(grep 'releaseNameDB' helm-charts/demo-apps-to-monitor/values.yaml | cut -d':' -f2 | tr -d '[:space:]')"
PATCH_NAME="$$(grep -m2 'name' patch-json-command.json | tail -n1 | cut -d':' -f2 | tr -d '[:space:]' | sed 's/\"//g')"
SECRET_PATCH_NAME=${RELEASE_NAME_DB}-${SECRET_DB}
.PHONY: update-darwin
update-darwin:
	sed -i '' "s/${PATCH_NAME}/${SECRET_PATCH_NAME}/" patch-json-command.json >> patch-json-command.json
	sed -i '' "s/${SECRET_HASHER}/${SECRET_DB}/" helm-charts/demo-apps-to-monitor/values.yaml >> helm-charts/demo-apps-to-monitor/values.yaml
	sed	-i '' "s/${VALUE_RELEASE_NAME_APP}/${RELEASE_NAME_DB}/" helm-charts/demo-apps-to-monitor/values.yaml >> helm-charts/demo-apps-to-monitor/values.yaml
