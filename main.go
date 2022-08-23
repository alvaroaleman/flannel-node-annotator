package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/alvaroaleman/flannel-node-annotator/controller"
	"github.com/kubermatic/machine-controller/pkg/signals"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig string
	var master string
	var addressType string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.StringVar(&addressType, "type", "ExternalIP", "which address (InternalIP or ExternalIP) should be used by flannel")
	flag.Parse()

	// creates the connection
	config, err := clientcmd.BuildConfigFromFlags(master, kubeconfig)
	if err != nil {
		glog.Fatal(err)
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatal(err)
	}

	stopChannel := signals.SetupSignalHandler()
	controller := controller.NewController(clientset, corev1.NodeAddressType(addressType), stopChannel)

	controller.Run(1, stopChannel)
}
