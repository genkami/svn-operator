# svn-operator

The svn-operator is a simple, opinionated operator that has following features:

* Does not guarantee extremely high availability. (if an SVN server crashes, just wait for a second and retry.)
* Only basic authentication is allowed.
* Can manage users and repositories declaratively.
* Does not provide path-based authorization.

## Installation

```
$ kubectl apply -f https://github.com/genkami/svn-operator/releases/download/v0.0.1/svn-operator.crds.yaml
$ kubectl apply -f https://github.com/genkami/svn-operator/releases/download/v0.0.1/svn-operator.yaml
```

## Examples

The following example creates an SVN server `svnserver-sample` that has a single SVN repository named `svnrepository-sample`. It also creates two groups named `svngroup-sample-reader` and `svngroup-sample-writer`. The former group has only read permission to the repository and the latter has full access to the repository. Finally, it creates two users `svnuser-sample-reader` who belongs to `svngroup-sample-reader` and `svnuser-sample-writer` who belongs to `svngroup-sample-writer`.

``` yaml
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNServer
metadata:
  name: svnserver-sample
spec:
  volumeClaimTemplate:
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 512M
---
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNRepository
metadata:
  name: svnrepository-sample
spec:
  svnServer: svnserver-sample
---
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNGroup
metadata:
  name: svngroup-sample-reader
spec:
  svnServer: svnserver-sample
  permissions:
  - repository: svnrepository-sample
    permission: r
---
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNGroup
metadata:
  name: svngroup-sample-writer
spec:
  svnServer: svnserver-sample
  permissions:
  - repository: svnrepository-sample
    permission: rw
---
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNUser
metadata:
  name: svnuser-sample-reader
spec:
  svnServer: svnserver-sample
  groups:
    - name: svngroup-sample-reader
  # The password is 'foobar'
  encryptedPassword: $2y$05$lHorekjyyp9w2fXD/ppQLOJ2N1KmY.9yiJ0mZQlkIeUpUg8enPN4e
---
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNUser
metadata:
  name: svnuser-sample-writer
spec:
  svnServer: svnserver-sample
  groups:
    - name: svngroup-sample-writer
  # The password is 'quux'
  encryptedPassword: $2y$05$skzShfjCsTKCYcvr55ByIO5G7icGU8Lofs2CpmR5AoGho9OzBLb4O
```

## License

Distributed under the Apache License Version 2.0. See LICENSE for more information.
