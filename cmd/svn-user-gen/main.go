package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"syscall"

	flag "github.com/spf13/pflag"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var err error
	var user string
	flag.StringVarP(&user, "user", "u", "", "The name of the user")
	flag.Parse()

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
	password, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin", err)
		os.Exit(1)
	}

	fmt.Fprint(os.Stderr, "\nRe-type Password: ")
	password2, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin", err)
		os.Exit(1)
	}

	if !bytes.Equal(password, password2) {
		fmt.Fprintln(os.Stderr, "password mismatch")
		os.Exit(1)
	}

	trimmedPassword := []byte(strings.TrimSuffix(string(password), "\n"))
	encryptedPassword, err := bcrypt.GenerateFromPassword(trimmedPassword, 5) // TODO
	if err != nil {
		fmt.Println(os.Stderr, "error encrypting password", err)
		os.Exit(1)
	}

	fmt.Println(`
apiVersion: svn.k8s.oyasumi.club/v1alpha1
kind: SVNUser
metadata:
  name: ` + user + `
spec:
  encryptedPassword: ` + string(encryptedPassword) + `
`)
}
