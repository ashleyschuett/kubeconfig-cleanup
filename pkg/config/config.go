package config

import (
	"fmt"
	"os"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kubeconfigutil "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"

	bootstrapapi "k8s.io/client-go/tools/bootstrap/token/api"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Manager struct {
	Original           *clientcmdapi.Config
	New                *clientcmdapi.Config
	prompter           *prompter
	path               string
	workqueue          chan ContextResult
	totalContexts      int
	contextedValidated int
}

func NewManager() *Manager {
	path := getKubeconfigPath()
	config, err := clientcmd.LoadFromFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ch := make(chan ContextResult)

	return &Manager{
		config,
		config.DeepCopy(),
		NewPrompter(),
		path,
		ch,
		len(config.Contexts),
		0,
	}
}

func getKubeconfigPath() string {
	defer fmt.Println()

	if kcfgp := os.Getenv("KUBECTL_PLUGINS_LOCAL_FLAG_KUBECONFIG"); kcfgp != "" {
		fmt.Printf("[kubeconfig] Using path '%s'\n", kcfgp)
		return kcfgp
	}

	fmt.Println("[kubeconfig] Using default path '$HOME/.kube/config'")
	home := os.Getenv("HOME")
	return fmt.Sprintf("%s/.kube/config", home)
}

func (m *Manager) Run() {

	for id, context := range m.Original.Contexts {
		go m.ValidateAndAddToWorkqueue(id, context)
	}

	m.runWorkqueue()
	m.RemoveUnusedUsers()
	m.Finish()
}

func (m *Manager) GetKubeconfigPath() string {
	return m.path
}

func (m *Manager) getContextsUser(context *clientcmdapi.Context) *clientcmdapi.AuthInfo {
	return m.Original.AuthInfos[context.AuthInfo]
}

func (m *Manager) getContextsCluster(context *clientcmdapi.Context) *clientcmdapi.Cluster {
	return m.Original.Clusters[context.Cluster]
}

func (m *Manager) removeCluster(cluster string) {
	delete(m.New.Clusters, cluster)
}

func (m *Manager) removeUser(user string) {
	delete(m.New.AuthInfos, user)
}

func (m *Manager) removeContext(context string) {
	delete(m.New.Contexts, context)

}

func (m *Manager) userIsInUse(user string) bool {
	count := 0

	for _, context := range m.New.Contexts {
		if context.AuthInfo != user {
			continue
		}

		count = count + 1
	}

	return count > 0
}

func (m *Manager) RemoveContext(id string, context *clientcmdapi.Context) {
	if m.prompter.RemoveContext(id) {
		m.removeCluster(context.Cluster)
		m.removeContext(id)
	}
}

func (m *Manager) RemoveUnusedUsers() {
	for user, _ := range m.New.AuthInfos {
		if !m.userIsInUse(user) && m.prompter.RemoveUser(user) {
			m.removeUser(user)
		}
	}
}

func (m *Manager) Finish() {
	config, _ := clientcmd.Write(*m.New)
	fmt.Println("----------- NEW KUBECONFIG --------------")
	fmt.Print(string(config))
	fmt.Println("-----------------------------------------")

	path := m.GetKubeconfigPath()
	if !m.prompter.WriteConfig() {
		if !m.prompter.WriteConfigToPath() {
			return
		}

		path = m.prompter.GetPath()
	}

	fmt.Println("Writing file to: ", path)
	err := clientcmd.WriteToFile(*m.New, path)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

func (m *Manager) Validate(context *clientcmdapi.Context) (bool, string) {
	// make request to server/healz
	cluster := m.getContextsCluster(context)
	configFromClusterInfo := kubeconfigutil.CreateBasic(
		cluster.Server,
		context.Cluster,
		context.AuthInfo,
		cluster.CertificateAuthorityData,
	)

	configFromClusterInfo.AuthInfos[context.AuthInfo] = m.getContextsUser(context)
	client, err := kubeconfigutil.ToClientSet(configFromClusterInfo)
	if err != nil {
		return false, fmt.Sprintf("error converting to clientset: %v", err)
	}
	_, err = client.CoreV1().ConfigMaps(metav1.NamespacePublic).Get(bootstrapapi.ConfigMapClusterInfo, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsForbidden(err) {
			// If the request is unauthorized, the cluster admin has not granted access to the cluster info configmap for unauthenticated users
			// In that case, trust the cluster admin and do not refresh the cluster-info credentials
			return true, fmt.Sprintf("[discovery] Could not access the %s ConfigMap for refreshing the cluster-info information, but the TLS cert is valid so proceeding...", bootstrapapi.ConfigMapClusterInfo)
		}

		return false, fmt.Sprintf("[discovery] Failed to validate the API Server's identity, will try again: [%v]", err)
	}

	return true, fmt.Sprintf("[discovery] Valid cluster associated with context")
}
