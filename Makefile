docker:
	docker build -f build/Dockerfile -t local/tarantool-operator:11.0.0 .

crds:
	operator-sdk generate crds
	cp -r deploy/* ./ci/helm-chart/templates/
