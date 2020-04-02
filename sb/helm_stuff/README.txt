Test to run a go application under kubenetes using helm.

docker contains a go app and Dockerfile to build it into an image.  use doit.sh to build and tag the image.

The app is run as a server and one or more clients.  Start the server with
  go run chat.go --mode server

and the clients with
  go run chat.go --mode client

Then what is typed into one client appears in the others.

kubenetes contains the basic manifest files for running the app in minikube.
They are not needed for helm but were a first step.  They can be applied using
  kubectl apply -f <filename>

helm contains the helm charts to run the chat app.  To run, cd into helm then run
  helm install <release_name> ./chat-app

where release_name is a name of your choosing such as "wobbly-cat".

Multiple releases can be installed at the same time.  A release can be uninstalled with
  helm uninstall <release_name>

Once installed, use kubectl get services to find out the external TCP port
for the service.  You can then connect a client chat app to the server running
in kubenetes using
  go run chat.go --mode client --port <external TCP port>

N.B. much of this may need to be done as root
