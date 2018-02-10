package controller

import (
	"time"

	"github.com/golang/glog"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	addressLabelName = "flannel.alpha.coreos.com/public-ip-overwrite"
	addressType      = "ExternalIP"
)

type Controller struct {
	kubeClient  kubernetes.Interface
	workqueue   workqueue.RateLimitingInterface
	nodesLister listerscorev1.NodeLister
}

func NewController(clientset *kubernetes.Clientset, stopCh <-chan struct{}) *Controller {
	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(clientset, time.Second*30)
	nodeInformer := kubeInformerFactory.Core().V1().Nodes()
	controller := &Controller{
		kubeClient:  clientset,
		nodesLister: nodeInformer.Lister(),
		workqueue:   workqueue.NewNamedRateLimitingQueue(workqueue.NewItemFastSlowRateLimiter(2*time.Second, 10*time.Second, 5), "nodes"),
	}

	nodeInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				controller.workqueue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				controller.workqueue.Add(key)
			}
		},
	})
	go kubeInformerFactory.Start(stopCh)
	glog.Infof("Waiting for cache sync...")
	kubeInformerFactory.WaitForCacheSync(stopCh)
	glog.Infof("Cache sync done!")

	return controller

}

func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	key, quit := c.workqueue.Get()
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two pods with the same key are never processed in
	// parallel.
	defer c.workqueue.Done(key)

	// Invoke the method containing the business logic
	glog.Infof("Processing node '%s'", key)
	err := c.syncNode(key.(string))
	// Handle the error if something went wrong during the execution of the business logic
	c.handleErr(err, key)
	return true
}

func (c *Controller) syncNode(key string) error {
	listerNode, err := c.nodesLister.Get(key)
	if err != nil {
		glog.Infof("Error getting node '%s': '%v'", key, err)
		return nil
	}

	node := listerNode.DeepCopy()
	glog.Infof("Syncing node '%s'", node.Name)
	for _, address := range node.Status.Addresses {
		if address.Type == addressType {
			if err := c.ensureAddressLabel(node, address.Address); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) ensureAddressLabel(node *corev1.Node, address string) error {
	var updated bool
	var err error
	if value, exists := node.Labels[addressLabelName]; exists {
		updated = true
		node.Labels[addressLabelName] = address
	} else {
		if value != address {
			updated = true
			node.Labels[addressLabelName] = address
		}
	}
	if updated {
		glog.Infof("Updating label of node '%s'", node.Name)
		_, err = c.kubeClient.CoreV1().Nodes().Update(node)
	}
	return err
}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {
		// Forget about the #AddRateLimited history of the key on every successful synchronization.
		// This ensures that future processing of updates for this key is not delayed because of
		// an outdated error history.
		c.workqueue.Forget(key)
		return
	}

	// This controller retries 5 times if something goes wrong. After that, it stops trying.
	if c.workqueue.NumRequeues(key) < 5 {
		glog.Infof("Error syncing Node %v: %v", key, err)

		// Re-enqueue the key rate limited. Based on the rate limiter on the
		// queue and the re-enqueue history, the key will be processed later again.
		c.workqueue.AddRateLimited(key)
		return
	}

	c.workqueue.Forget(key)
	// Report to an external entity that, even after several retries, we could not successfully process this key
	runtime.HandleError(err)
	glog.Infof("Dropping Node %q out of the queue: %v", key, err)
}

func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	// Let the workers stop when we are done
	defer c.workqueue.ShutDown()
	glog.Info("Starting Node controller")

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	glog.Info("Stopping Node controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}
