helm dependency update
helm upgrade --install photobox --namespace $1 -f values.yaml . --atomic $2
