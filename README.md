# k8s-controller-custom-resource

## features
- redis-operator: including a simple master-slave mode which was main based on the resources of **k8s.StatefulSet** and **k8s.Service**
- mysql-operator: the same with the redis-operator

## core/v1
1. core/v1/interfaces would add the storage plugins(e.g. nfs) later for dynamically creating pv and pvc
2. core/v1/interfaces need to add api for creating instances instead of creating yaml files and executing `kubectl apply -f *.yaml`

## custom-controller
```go
opt := k8sCoreV1.NewOption(&mysqlOperatorV1.MysqlOperator{},
    controllerAgentName,
    operatorKindName,
    mysqlOperatorScheme.AddToScheme(scheme.Scheme),
    sampleclientset,
    fooInformer,
    fooInformer.Informer().HasSynced,
    fooInformer.Informer().AddEventHandler,
    CompareResourceVersion,
    Get,
    Sync)
opts := k8sCoreV1.NewOptions()
if err := opts.Add(opt); err != nil {
    klog.Fatal(err)
}
op := k8sCoreV1.NewKubernetesOperator(kubeclientset, stopCh, controllerAgentName, opts)
kc := k8sCoreV1.NewKubernetesController(op)
...
```


## redis-operator

**debug**
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o redisoperatorcrd cmd/redisoperator/main.go
./redisoperatorcrd -kubeconfig=$HOME/.kube/config -alsologtostderr=true
```

## mysql-operator

**debug**
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o mysqloperatorcrd cmd/mysqloperator/main.go
./mysqloperatorcrd -kubeconfig=$HOME/.kube/config -alsologtostderr=true
```

### todo