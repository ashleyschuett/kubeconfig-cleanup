# kubeconfig-cleanup

## About

### Description
A kubectl plugin to clean up your kubeconfig file

### Installation
Add `kubeconfig-cleanup` to your `kubectl` plugins directory. For more information about how plugins are loaded, please see the [official documentation](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).
```
curl -LO https://github.com/ashleyschuett/kubeconfig-cleanup/releases/download/v1.0.1/cleanup && \
curl -LO https://raw.githubusercontent.com/ashleyschuett/kubeconfig-cleanup/v1.0.1/plugin.yaml && \
mkdir ~/.kube/plugins/kubeconfig-cleanup && mv cleanup plugin.yaml
```

### Usage
Parse through kubeconfig in default directory
```
kubectl plugin kubeconfig-cleanup
```

Parse through kubeconfig at custom location
```
kubectl plugin kubeconfig-cleanup --kubeconfig /custom/kubeconfig/location
```
