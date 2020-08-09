helm dependency update

kubectl create namespace $1
helm upgrade --install photobox --namespace $1 -f values.yaml . --atomic $2
