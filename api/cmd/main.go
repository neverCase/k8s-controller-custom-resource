package main

import (
	"flag"

	"github.com/nevercase/k8s-controller-custom-resource/api/conf"
	"github.com/nevercase/k8s-controller-custom-resource/api/service"
	"github.com/nevercase/k8s-controller-custom-resource/pkg/signals"
	"k8s.io/klog"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()
	s := service.NewService(conf.Init())
	s.Listen()
	<-stopCh
}
