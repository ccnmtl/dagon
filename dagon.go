package main

import (
	"github.com/jmckaskill/goldap"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// this is the struct for pulling results out of LDAP
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
	Sn string `json:"sn"`
	GivenName string `json:"givenName"`
	TelephoneNumber string `json:"telephoneNumber"`
	Telephonenumber string `json:"telephonenumber"`
	CuMiddlename string `json:"cuMiddlename"`
	Uid string `json:"uid"`
	Firstname string `json:"firstname"`
	DepartmentNumber string `json:"departmentNumber"`
	ObjectClass string `json:"objectClass"`
	Lastname string `json:"lastname"`
	Title string `json:"title"`
	Mail string `json:"mail"`
	Campusphone string `json:"campusphone"`
	Uni string `json:"uni"`
	PostalAddress string `json:"postalAddress"`
	Ou string `json:"ou"`
	Cn string `json:"cn"`
	Found bool `json:"found"`
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

	b, err := json.Marshal(or)
	if err != nil {
    fmt.Println("error:", err)
	}
	os.Stdout.Write(b)
}
