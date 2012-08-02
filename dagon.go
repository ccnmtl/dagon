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

// {
//  "telephoneNumber": "+1 212 854 1813",
//  "telephonenumber": "+1 212 854 1813",
//  "cuMiddlename": "N.",
//  "uid": "anp8, anders",
//  "firstname": "Anders",
//  "departmentNumber": "2209102",
//  "objectClass": "person, organizationalPerson, inetOrgPerson, cuPerson, cuRestricted, eduPerson",
//  "lastname": "Pearson",
//  "title": "Senior Programmer",
//  "mail": "anders@columbia.edu",
//  "campusphone": "MS 4-1813",
//  "sn": "Pearson",
//  "uni": "anp8",
//  "found": true,
//  "postalAddress": "505 Butler Library$Mail Code: 1130$United States",
//  "givenName": "Anders",
//   "ou": "Ctr for New Media Teaching^Ctr for New Media Teaching^Ctr New Media Teach & Lrng",
//   "cn": "Anders N. Pearson"
// }


type result struct {
	Sn string
	GivenName string
	TelephoneNumber string
	Telephonenumber string
	CuMiddlename string
	Uid string
	Firstname string
	DepartmentNumber string
	ObjectClass string
	LastName string
	Title string
	Mail string
	Campusphone string
	Uni string
	PostalAddress string
	Ou string
	Cn string
	Found bool
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

	fmt.Printf("%v\n",r)
//	fmt.Printf("%v\n",channel_query(db, base_dn, ur))
//	fmt.Printf("%v\n",slice_query(db, base_dn, ur))
//	fmt.Printf("%v\n",func_query(db, base_dn, ur)) // only works with patched goldap


    //             for k, v in values.items():
    //                 results_dict[k] = ", ".join(v)
    //                 if k == 'sn':
    //                     results_dict['lastname'] = v[0]
    //                 if k == 'givenname':
    //                     results_dict['firstname'] = v[0]
    //                 if k == 'givenName':
    //                     results_dict['firstname'] = v[0]
    //                 if k == 'telephoneNumber':
    //                     results_dict['telephonenumber'] = v[0]

    // if results_dict['lastname'] == "":
    //     results_dict['lastname'] = uni

}
