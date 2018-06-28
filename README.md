# kubeconfig-cleanup

## About

### Description
A kubectl plugin to clean up your kubeconfig file

### Installation
Add `kubeconfig-cleanup` to your `kubectl` plugins directory. For more information about how plugins are loaded, please see the [official documentation](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).
```
git clone git@github.com:ashleyschuett/kubeconfig-cleanup.git ~/.kube/plugins/kubeconfig-cleanup
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
