# k8s-controller-custom-resource

#### redis-operator

debug:
```sh
./redisoperatorcrd -kubeconfig=$HOME/.kube/config -alsologtostderr=true
```


#### todo
1. core/v1/interfaces would add pv and pvc later

2. core/v1/interfaces need to add api for creating instances instead of using\
yaml files and executing `kubectl apply -f *.yaml`