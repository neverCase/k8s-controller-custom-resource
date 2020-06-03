# k8s-controller-custom-resource

## features
- redis-operator: including a simple master-slave mode which was main based on the resources of **k8s.StatefulSet** and **k8s.Service**
- mysql-operator: the same with the redis-operator

## todo
1. core/v1/interfaces would add the storage plugins(e.g. nfs) later for dynamically creating pv and pvc
2. core/v1/interfaces need to add api/client(just like kubernetes/client-go) for creating instances instead of creating yaml files and executing `kubectl apply -f *.yaml`


## RedisOperator

### define resource
```sh
$ cat > redis-resource.yaml <<EOF
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: redisoperators.redisoperator.nevercase.io
spec:
  group: redisoperator.nevercase.io
  version: v1
  names:
    kind: RedisOperator
    plural: redisoperators
  scope: Namespaced
EOF
$ kubectl apply -f redis-resource.yaml
```

### define demo file
```sh
$ cat > example-redis.yaml <<EOF
apiVersion: redisoperator.nevercase.io/v1
kind: RedisOperator
metadata:
  name: example-redis
spec:
  masterSpec:
    spec:
      name: "redis-demo"
      replicas: 1
      image: domain/redis-slave:1.1
      imagePullSecrets:
        - name: private-secret
  slaveSpec:
    spec:
      name: "redis-demo"
      replicas: 4
      image: domain/redis-slave:1.1
      imagePullSecrets:
        - name: private-secret
EOF
$ kubectl apply -f example-redis.yaml
```

### run controller
```sh
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o redisoperatorcrd cmd/redisoperator/main.go
$ ./redisoperatorcrd -kubeconfig=$HOME/.kube/config -alsologtostderr=true
I0603 14:48:38.844075   20412 controller.go:72] Setting up event handlers
I0603 14:48:38.844243   20412 controller.go:195] Starting Foo controller
I0603 14:48:38.844249   20412 controller.go:198] Waiting for informer caches to sync
I0603 14:48:38.944352   20412 controller.go:209] Starting workers
I0603 14:48:38.944366   20412 controller.go:215] Started workers
...
I0603 14:48:47.721574   20412 event.go:255] Event(v1.ObjectReference{Kind:"RedisOperator", ... type: 'Normal' reason: 'Synced' Foo synced successfully
```

## mysql-operator

The usage was the same with the RedisOperator. 


## New custom-controller
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