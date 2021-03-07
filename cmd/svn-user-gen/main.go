package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"
	"text/template"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	var err error
	var user string
	var cost int
	var svnServer string
	var svnGroups string
	flag.StringVar(&user, "user", "", "The name of the user")
	flag.StringVar(&svnServer, "svn-server", "TYPE_THE_SERVER_NAME_HERE", "The name of SVNServer resource")
	flag.StringVar(&svnGroups, "svn-groups", "", "Comma-separated list of SVNGroups that the user belongs to")
	flag.IntVar(&cost, "cost", bcrypt.DefaultCost, "The cost of bcrypt encryption")
	flag.Parse()

	if cost < bcrypt.MinCost || bcrypt.MaxCost < cost {
		fmt.Fprintf(os.Stderr, "cost must be between %d and %d\n", bcrypt.MinCost, bcrypt.MaxCost)
		os.Exit(1)
	}

	if user == "" {
		fmt.Fprint(os.Stderr, "Username: ")
		r := bufio.NewReader(os.Stdin)
		user, err = r.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading from stdin", err)
			os.Exit(1)
		}
		user = strings.TrimSuffix(user, "\n")
	}

	fmt.Fprint(os.Stderr, "Password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin", err)
		os.Exit(1)
	}

	fmt.Fprint(os.Stderr, "\nRe-type Password: ")
	password2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin", err)
		os.Exit(1)
	}

	if !bytes.Equal(password, password2) {
		fmt.Fprintln(os.Stderr, "password mismatch")
		os.Exit(1)
	}

	trimmedPassword := []byte(strings.TrimSuffix(string(password), "\n"))
	encryptedPassword, err := bcrypt.GenerateFromPassword(trimmedPassword, cost)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error encrypting password", err)
		os.Exit(1)
	}

	tmpl, err := template.New("svnuser.yaml").Parse(tmplSource)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse template", err)
		os.Exit(1)
	}
	groups := make([]string, 0)
	for _, s := range strings.Split(svnGroups, ",") {
		if len(s) > 0 {
			groups = append(groups, s)
		}
	}
	err = tmpl.Execute(os.Stdout, map[string]interface{}{
		"User":              user,
		"EncryptedPassword": string(encryptedPassword),
		"Server":            svnServer,
		"Groups":            groups,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to write SVNUser", err)
		os.Exit(1)
	}
}

const tmplSource = `
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNUser
metadata:
  name: {{ .User }}
spec:
  svnServer: {{ .Server }}
  encryptedPassword: {{ .EncryptedPassword }}
{{- if lt 0 (len .Groups) }}
  groups:
{{- range .Groups }}
  - name: {{ . }}
{{- end -}}
{{- end }}
`
