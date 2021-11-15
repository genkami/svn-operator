module github.com/genkami/svn-operator

go 1.15

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.14.0
	go.uber.org/zap v1.18.1
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/yaml v1.2.0
)
