package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/alvaroaleman/flannel-node-annotator/controller"
	"github.com/kubermatic/machine-controller/pkg/signals"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig string
	var master string

	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
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
	controller := controller.NewController(clientset, stopChannel)

	controller.Run(1, stopChannel)
}
