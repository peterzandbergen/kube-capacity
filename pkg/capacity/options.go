package capacity

// Options is a struct containing the command line options
// FetchAndPrint depends on
type Options struct {
	ShowContainers        bool
	ShowPods              bool
	ShowUtil              bool
	ShowPodCount          bool
	HideRequests          bool
	HideLimits            bool
	PodLabels             string
	NodeLabels            string
	NodeTaints            string
	ExcludeTainted        bool
	NamespaceLabels       string
	Namespaces            []string
	KubeContext           string
	KubeConfig            string
	InsecureSkipTLSVerify bool
	OutputFormat          string
	SortBy                string
	AvailableFormat       bool
	ImpersonateUser       string
	ImpersonateGroup      string
}
