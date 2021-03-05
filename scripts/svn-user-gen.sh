#!/bin/bash

user="$1"
shift

if [ -z "$user" ]; then
    echo "usage: svn-user-gen.sh USERNAME [htpasswd OPTIONS...]" >&2
    exit 1
fi

if ! which htpasswd >/dev/null 2>&1; then
    echo "htpasswd not installed" >&2
    exit 1
fi

# NOTE: this script forces to use bcrypt in order to avoid using insecure passwords accidentally.
password=$(htpasswd -nB "$@" "$user" | cut -d : -f 2-)

cat <<EOF
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNUser
metadata:
  name: $user
spec:
  encryptedPassword: $password
EOF
