#!/bin/bash

mkdir -p /svn/repos

# TODO: remove this
cat /etc/svn-config/Repos | grep name: | cut -d ' ' -f 3 | while read repo; do
    echo "Creating $repo..."
    svnadmin create "/svn/repos/$repo"
done
