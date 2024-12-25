rem Delete the ToDo-app
kubectl delete -f todo-app.yaml
rem Delete the CockroachDB Helm Chart
helm uninstall cockroachdb
rem Delete the kind cluster
kind delete cluster