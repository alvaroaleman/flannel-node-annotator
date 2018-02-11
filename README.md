# flannel-node-annotator

A simple controller that adds a `flannel.alpha.coreos.com/public-ip-overwrite`
annotation containing a nodes `ExternalIP` if set.

May be helpfull if some of your nodes are behind a NAT.

## Quickstart

```bash
kubectl apply -f kubernetes/deployment.yaml
```
