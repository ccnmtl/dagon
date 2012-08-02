package main

import (
	"github.com/jmckaskill/goldap"
	"fmt"
	"os"
)

func printer (r string) {
	fmt.Println("out!")
	fmt.Printf("%v\n", r)
}

type result struct {
	Sn string
	GivenName string
	TelephoneNumber string
}

func channel_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) {
	out := make(chan result)

	go func() {
		err := db.SearchTree(out, base_dn, ur)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	r := <- out
	fmt.Printf("%v\n", r)
}

func slice_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) {
	out := make([]result,0)

	err := db.SearchTree(&out, base_dn, ur)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v\n", out)
}

func func_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) {
	out := func (r result) {
		fmt.Printf("%v\n", r)
	}

	err := db.SearchTree(out, base_dn, ur)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func struct_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) {
	out := result{}

	err := db.SearchTree(&out, base_dn, ur)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("%v\n", out)
}


func main () {
	if len(os.Args) < 2 {
		fmt.Println("must specify a UNI")
		return
	}
	uni := os.Args[1]

	auth := make([]ldap.AuthMechanism,1)
	auth[0] = ldap.SimpleAuth{User:"",Pass:""}
	cfg := ldap.ClientConfig{Dial:nil,Auth:auth,TLS:nil}

	db := ldap.Open("ldap://ldap.columbia.edu", &cfg)
	base_dn := ldap.ObjectDN("o=Columbia University, c=us")
	ur := ldap.Equal{Attr:"uni",Value:[]byte(uni)}

	channel_query(db, base_dn, ur)
	slice_query(db, base_dn, ur)
	func_query(db, base_dn, ur)
	struct_query(db, base_dn, ur)
}
