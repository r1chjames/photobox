# Running the API locally
There are several ways to run the API locally;
- Kubernetes using YAML
- Kubernetes using Helm
- Docker using docker-compose

## Kubernetes using YAML
The API can be deployed locally onto a Kubernetes (or K3s, Minikube). It requires that a suitable load balancer is available within the cluster. Minikube and K3s provide this out of the box.

If you want to present a directory to Minikube:

```minikube mount minikube mount <local_directory>:<directory_inside_VM>```

Deploy the app from inside the `kubernetes` directory:

```kubectl apply -f . -n <namespace>```

The API should be avalable on the external IP address on the `service/photobox` resource. In Minikube run `minikube service photobox` to get the API URL.

## Kubernetes using Helm
A Helm chart exists that deploys the API onto a Kubernetes (or K3s, minikube). It requires that Helm is installed locally (and set up on the cluster if using anything pre-Helm v3) and that a suitable load balancer is available within the cluster. Minikube and K3s provide this out of the box.

If you want to present a directory to Minikube:

```minikube mount minikube mount <local_directory>:<directory_inside_VM>```

Deploy the app from inside the `chart` directory:

```chmod +x ./run.sh && ./run.sh <namespace>```

The API should be avalable on the external IP address on the `service/photobox` resource. In Minikube run `minikube service photobox` to get the API URL.

## Docker using docker-compose
Providing Docker is installed, execute the following command:

```docker-compose up```

The API should be avalable on `localhost:8080`.