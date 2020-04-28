# k8s-controller-custom-resource

## features
- redis-operator: including simple master-slave mode
- mysql-operator

## core/v1
1. core/v1/interfaces would add pv and pvc later
2. core/v1/interfaces need to add api for creating instances instead of using\
yaml files and executing `kubectl apply -f *.yaml`

## redis-operator

**debug**
```sh
./redisoperatorcrd -kubeconfig=$HOME/.kube/config -alsologtostderr=true
```

## mysql-operator

**debug**
```sh
./mysqloperatorcrd -kubeconfig=$HOME/.kube/config -alsologtostderr=true
```

### todo