# kubeconfig-cleanup

## About

### Description
A kubectl plugin to clean up your kubeconfig file

### Installation
Add `kubeconfig-cleanup` to your `kubectl` plugins directory. For more information about how plugins are loaded, please see the [official documentation](https://kubernetes.io/docs/tasks/extend-kubectl/kubectl-plugins/).

See [Releases](https://github.com/ashleyschuett/kubeconfig-cleanup/releases) for the latest release matching your OS.

```bash
# Pick the right release, e.g. "v1.0.1 Linux x64"
$ export BINARY=kubeconfig-cleanup_1.0.1_Linux_x86_64.tar.gz
```

Then run `install.sh` by either cloning this repository or directly downloading the script.

**Note:** The script and plugin configuration is configured to use `~/.kube/plugins/kubeconfig-cleanup` as the plugin and installation folder.

```bash
$ git clone https://github.com/ashleyschuett/kubeconfig-cleanup
$ cd kubeconfig-cleanup

# make sure the script it's executable
$ chmod +x install.sh
$ ./install.sh
```

Alternatively with curl:

```bash
$ curl -LO https://raw.githubusercontent.com/ashleyschuett/kubeconfig-cleanup/master/install.sh

# make sure the script it's executable
$ chmod +x install.sh
$ ./install.sh
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
