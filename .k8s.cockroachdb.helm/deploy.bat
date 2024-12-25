rem CockroachDB Helm Chart
helm repo add cockroachdb https://charts.cockroachdb.com/
helm repo update
helm install cockroachdb cockroachdb/cockroachdb -f values.yaml
rem ToDo-app
kind load docker-image todo-app:latest
kubectl apply -f todo-app.yaml
rem Port Forward to access the ToDo-app
kubectl port-forward service/todo-app 8000:8000