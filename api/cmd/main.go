package main

import (
	"flag"
	"github.com/Shanghai-Lunara/pkg/zaplogger"
	"github.com/nevercase/k8s-controller-custom-resource/api/conf"
	"github.com/nevercase/k8s-controller-custom-resource/api/service"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
	"k8s.io/klog/v2"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()
	zaplogger.Sugar().Info("k8s-custom-api-server is starting")
	s := service.NewService(conf.Init())
	zaplogger.Sugar().Info("k8s-custom-api-server is running")
	<-stopCh
	zaplogger.Sugar().Info("k8s-custom-api-server trigger shutdown")
	s.Close()
	<-stopCh
	zaplogger.Sugar().Info("k8s-custom-api-server shutdown gracefully")
}