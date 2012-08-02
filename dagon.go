package main

import (
	"github.com/jmckaskill/goldap"
	"fmt"
	"os"
	"strings"
)

type result struct {
	Sn string
	GivenName string
	TelephoneNumber string
	CuMiddlename string
	Uid []string
	Firstname string
	DepartmentNumber string
	ObjectClass []string
	Lastname string
	Title string
	Mail string
	Campusphone string
	Uni string
	PostalAddress string
	Ou string
	Cn string
}

func channel_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) result {
	out := make(chan result)

	go func() {
		err := db.SearchTree(out, base_dn, ur)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	return <- out
}

func slice_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) result {
	out := make([]result,0)

	err := db.SearchTree(&out, base_dn, ur)
	if err != nil {
		fmt.Println(err.Error())
	}
	return out[0]
}

func func_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) result {
	ch := make(chan result)
	out := func (r result) {
		ch <- r
	}

	go func() {
		err := db.SearchTree(out, base_dn, ur)
		if err != nil {
			fmt.Println(err.Error())
		}
	}()
	return <-ch	
}

func struct_query(db *ldap.DB, base_dn ldap.ObjectDN, ur ldap.Equal) result {
	out := result{}

	err := db.SearchTree(&out, base_dn, ur)
	if err != nil {
		fmt.Println(err.Error())
	}
	return out
}

// this struct is set up to match our old interface
type compatresult struct {
	Sn string
	GivenName string
	TelephoneNumber string
	Telephonenumber string
	CuMiddlename string
	Uid string
	Firstname string
	DepartmentNumber string
	ObjectClass string
	Lastname string
	Title string
	Mail string
	Campusphone string
	Uni string
	PostalAddress string
	Ou string
	Cn string
	Found bool
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

	r := struct_query(db, base_dn, ur)

	or := compatresult{}
	or.Sn = r.Sn
	or.GivenName = r.GivenName
	or.TelephoneNumber = r.TelephoneNumber
	or.CuMiddlename = r.CuMiddlename
	or.Uid = strings.Join(r.Uid, ", ")
	or.Firstname = r.Firstname
	or.DepartmentNumber = r.DepartmentNumber
	or.ObjectClass = strings.Join(r.ObjectClass, ", ")
	or.Lastname = r.Lastname
	or.Title = r.Title
	or.Mail = r.Mail
	or.Campusphone = r.Campusphone
	or.Uni = r.Uni
	or.PostalAddress = r.PostalAddress
	or.Ou = r.Ou
	or.Cn = r.Cn
	or.Found = true

	// duplicating the logic from previous python versions
	// this looks a bit silly, but ensures that we are generating
	// as close as possible to the previous output
	if r.Sn != "" {
		or.Lastname = r.Sn
	}
	if r.GivenName != "" {
		or.Firstname = r.GivenName
	}
	if r.TelephoneNumber != "" {
		or.Telephonenumber = r.TelephoneNumber
	}
	if or.Lastname == "" {
		// at least put a UNI in there
		or.Lastname = uni
	}

	fmt.Printf("%v\n",or)
}
