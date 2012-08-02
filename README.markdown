dagon
=====

Very simple commandline tool to make LDAP queries against Columbia's
LDAP directory and output the result in JSON:

Just call it with a UNI as the first and only parameter

    $ dagon anp8
    {"sn":"Pearson","givenName":"Anders","telephoneNumber":"+1 212 854
    1813","telephonenumber":"+1 212 854
    1813","cuMiddlename":"N.","uid":"anp8,
    anders","firstname":"Anders","departmentNumber":"2209102","objectClass":"person,
    organizationalPerson, inetOrgPerson, cuPerson, cuRestricted,
    eduPerson","lastname":"Pearson","title":"Senior
    Programmer","mail":"anders@columbia.edu","campusphone":"MS
    4-1813","uni":"anp8","postalAddress":"505 Butler Library$Mail
    Code: 1130$United States","ou":"Ctr for New Media Teaching^Ctr for
    New Media Teaching^Ctr New Media Teach & Lrng","cn":"Anders
    N. Pearson","found":true}

(Yes, that's my real info. It's pretty much all publically available
info on CCNMTL's website anyway)

Note that one of the fields is a boolean called "found". That will be
set to false if the LDAP server didn't have a match for that
UNI. Otherwise, the fields are just what the LDAP server returns.

Motivation
----------

We don't really do much with LDAP at CCNMTL. Our main use-case is that
a user has logged into one of our apps via WIND so we have their UNI
and we'd like to know their full name and perhaps department so we can
set up an account for them. AFAIK, LDAP is the *only* way for us to
reliably get that info.

Previously, we've implemented this with a Python library. Python's
LDAP support is via a wrapper on some C libraries. Installing on
Ubuntu involves pulling in a whole mess of dependencies and installing
inside a virtualenv can be a pain and involve long compiles (just to
get their name!). Getting the python-ldap dependencies installed on
developer OS X machines seems to be even more complicated and painful.

So we set that up as a web-service that responds to a simple GET
request with JSON of the response. That's nice, since all the LDAP
dependencies only have to be dealt with on a single server, but
introduces yet another service that needs to be deployed and monitored
and can act as a single point of failure in our systems (LDAP web
gateway server is offline -> all our apps are broken). 

Dagon takes yet another approach. This time we have a small(ish)
statically compiled binary that can run pretty much anywhere and
handles the LDAP mess for us (outputting JSON that's identical to what
the web service generated) without any dependency
nightmares. Installing Dagon on a server involves copying the binary
to that server. Easy. 

It's even small enough that I plan to just bundle Linux and OS X
binaries of Dagon right in our Django apps so no one even has to think
about whether Dagon has been installed on a particular server. 

