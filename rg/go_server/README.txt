Test to run a simple Go server under Kubenetes using Minikube

1) Start Minikube

minikube start --driver=none

2) Build the Docker image

docker build -t go_server .

3) Generate yaml for creating a deployment

kubectl create deployment go-server-node --image=go_server --output=yaml > manifest.yaml

4) Edit the manifest.yaml file changing the imagePullPolicy from Always to Never and save

5) Delete the original deployment:

kubectl delete deployment go-server-node

6) Start the deployment again:

kubectl create -f manifest.yaml

7) Expose the Pod to the public internet:

kubectl expose deployment go-server-node --type=LoadBalancer --port=8080

8) Run the following command:

minikube service go-server-node

9) To delete all services, pods and deployments:

kubectl delete services,pods,deployments --all

10) To stop Minikube

minikube stop

11) To delete the existing Minikube VM:

minikube delete -p minikube
