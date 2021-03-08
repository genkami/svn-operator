# svn-operator

![ci status](https://github.com/genkami/svn-operator/workflows/Test/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/genkami/svn-operator.svg)](https://pkg.go.dev/github.com/genkami/svn-operator)


The svn-operator is a simple, opinionated operator that has following features:

* Does not guarantee extremely high availability. (if an SVN server crashes, just wait for a second and retry.)
* Only basic authentication is allowed.
* Can manage users and repositories declaratively.
* Does not provide path-based authorization.

## Installation

```
$ kubectl apply -f https://github.com/genkami/svn-operator/releases/download/v0.2.1/svn-operator.yaml
```

## Examples

The following example creates following resources:

* An SVN server `svnserver-sample`.
* A single SVN repository named `svnrepository-sample`  whcih belongs to `svnserver-sample`.
* Two groups named `svngroup-sample-reader` and `svngroup-sample-writer`.  
  The former group has only read permission to the repository and the latter has full access to the repository.
* Two users `svnuser-sample-reader` who belongs to `svngroup-sample-reader`, and `svnuser-sample-writer` who belongs to `svngroup-sample-writer`.

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

Additionally, if you want to expose SVN server to the internet, you have to set up ingress (or LoadBalancer, etc.) like this:

``` yaml
apiVersion: v1
kind: Service
metadata:
  name: svnserver-sample-lb
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    targetPort: 80
    protocol: TCP
  selector:
    # These two labels are generated by svn-operator
    app: subversion
    svn.k8s.oyasumi.club/name: svnserver-sample
---
# WARNING: This configuration is INSECURE since svn-operator uses basic auth.
# You must use HTTPS in production environments.
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: svnserver-sample-lb
spec:
  backend:
    serviceName: svnserver-sample-lb
    servicePort: 80
```

Then you can checkout repositories served under `http://<INGRESS_IP_ADDRESS>/repos/`.

```
$ svn checkout http://<INGRESS_IP_ADDRESS>/repos/svnrepository-sample

Authentication realm: <http://<INGRESS_IP_ADDRESS>:80> SVN Server
Username: svnuser-sample-writer
Password for 'svnuser-sample-writer': ****

Checked out revision 0.
```

## Password Encryption
The `EncryptedPassword` field can be generated by using `htpasswd` command:

```
$ htpasswd -nB john | cut -d : -f 2-
New password: 
Re-type new password: 
$2y$05$sZw4te5XgfiRjNVNhLRVuO7cgiqbTAcdPRvzRog0r8Tj.lNAnpKyi
```

Or dedicated CLI tool that we provide:

```
$ go get github.com/genkami/svn-operator/cmd/svn-user-gen
$ svn-user-gen -user john
Password: 
Re-type Password: 
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNUser
metadata:
  name: john
spec:
  svnServer: TYPE_THE_SERVER_NAME_HERE
  encryptedPassword: $2a$10$teGKPe/vdxOvSRwpCN7iH.Neu.KH8sc.33ylcNSO3bDriKbua/48u
```


## License

Distributed under the Apache License Version 2.0. See LICENSE for more information.
