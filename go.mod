module github.com/genkami/svn-operator

go 1.15

require (
	github.com/fsnotify/fsnotify v1.4.9
	github.com/go-logr/logr v0.4.0
	github.com/go-logr/zapr v0.4.0
	github.com/onsi/ginkgo v1.15.1
	github.com/onsi/gomega v1.11.0
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20210314154223-e6e6c4f2bb5b
	golang.org/x/term v0.0.0-20210220032956-6a3ed077a48d
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/yaml v1.2.0
)
