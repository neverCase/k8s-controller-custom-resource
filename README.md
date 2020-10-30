# k8s-controller-custom-resource

## Features
- k8s-api: supports watching and listing default resources (such as Service, Pod, Configmap) and another custom resources definition
- redis-operator: including a simple master-slave mode which was main based on the resources of **k8s.StatefulSet** and **k8s.Service**
- mysql-operator: the same with the redis-operator

## Operators to do
- core/v1/interfaces would add the storage plugins(e.g. nfs) later for dynamically creating pv and pvc

## Api to do
- ingress controller
 

## RedisOperator

### clone git, build docker image and compile controller
```sh
$ git clone https://github.com/neverCase/k8s-controller-custom-resource.git
$ cd k8s-controller-custom-resource 

# build image
$ make mysql
$ make redis

# compile controller
$ CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o multiplexcrd cmd/multiplex/main.go
```

### define the resource of the `RedisOperator`
```sh
$ cat > redis-resource.yaml <<EOF
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  name: redisoperators.nevercase.io
spec:
  group: nevercase.io
  version: v1
  names:
    kind: RedisOperator
    plural: redisoperators
  scope: Namespaced
EOF

$ kubectl apply -f redis-resource.yaml
customresourcedefinition.apiextensions.k8s.io/redisoperators.nevercase.io created
```

### define demo file
```sh
$ cat > example-redis.yaml <<EOF
apiVersion: nevercase.io/v1
kind: RedisOperator
metadata:
  name: example-redis
spec:
  masterSpec:
    spec:
      name: "redis-cn1"
      replicas: 1
      image: harbor.domain.com/helix-saga/redis-slave:1.1
      imagePullSecrets:
        - name: private-harbor
      volumePath: /mnt/nas1
      containerPorts:
        - containerPort: 6379
          protocol: TCP
      servicePorts:
        - port: 6379
          protocol: TCP
          targetPort: 6379
      resources:
        limits:
          memory: "1Gi"
          cpu: "100m"
        requests:
          memory: "0.5Gi"
          cpu: "100m"
  slaveSpec:
    spec:
      name: "redis-cn1"
      replicas: 4
      image: harbor.domain.com/helix-saga/redis-slave:1.1
      imagePullSecrets:
        - name: private-harbor
      volumePath: /mnt/nas1
      containerPorts:
        - containerPort: 6379
          protocol: TCP
      servicePorts:
        - port: 6379
          protocol: TCP
          targetPort: 6379
      resources:
        limits:
          memory: "1Gi"
          cpu: "100m"
        requests:
          memory: "0.5Gi"
          cpu: "100m"
EOF

$ kubectl apply -f example-redis.yaml
redisoperator.nevercase.io/example-redis created
```

### run the controller
```sh
$ ./multiplexcrd -kubeconfig=$HOME/.kube/config -alsologtostderr=true
I0603 14:48:38.844075   20412 controller.go:72] Setting up event handlers
I0603 14:48:38.844243   20412 controller.go:195] Starting Foo controller
I0603 14:48:38.844249   20412 controller.go:198] Waiting for informer caches to sync
I0603 14:48:38.944352   20412 controller.go:209] Starting workers
I0603 14:48:38.944366   20412 controller.go:215] Started workers
...
I0603 14:48:47.721574   20412 event.go:255] Event(v1.ObjectReference{Kind:"RedisOperator", ... type: 'Normal' reason: 'Synced' Foo synced successfully
```

### watch status
```sh
$ kubectl get statefulset
NAME                           READY   AGE
statefulset-redis-demo-master   1/1     46s
statefulset-redis-demo-slave    4/4     46s

$ kubectl get pod
NAME                                                    READY   STATUS      RESTARTS   AGE
statefulset-redis-demo-master-0                          1/1     Running     0          101s
statefulset-redis-demo-slave-0                           1/1     Running     0          101s
statefulset-redis-demo-slave-1                           1/1     Running     0          99s
statefulset-redis-demo-slave-2                           1/1     Running     0          98s
statefulset-redis-demo-slave-3                           1/1     Running     0          97s

$ kubectl get svc
NAME                       TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
service-redis-demo-master  ClusterIP   10.96.110.148   <none>        6379/TCP   4m38s
service-redis-demo-slave   ClusterIP   10.96.0.120     <none>        6379/TCP   4m38s
```

## MysqlOperator

The usage was the same with the RedisOperator. 


## New custom-controller
```go
opt := k8sCoreV1.NewOption(&mysqlOperatorV1.MysqlOperator{},
    controllerName,
    OperatorKindName,
    mysqlOperatorScheme.AddToScheme(scheme.Scheme),
    clientSet,
    fooInformer,
    fooInformer.Informer(),
    CompareResourceVersion,
    Get,
    Sync,
    SyncStatus)
opts := k8sCoreV1.NewOptions()
if err := opts.Add(opt); err != nil {
    klog.Fatal(err)
}
op := k8sCoreV1.NewKubernetesOperator(k8sClientSet, stopCh, controllerName, opts)
kc := k8sCoreV1.NewKubernetesController(op)
...
```