package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	var err error
	var user string
	var cost int
	flag.StringVar(&user, "user", "", "The name of the user")
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
	encryptedPassword, err := bcrypt.GenerateFromPassword(trimmedPassword, cost)
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
