# Authsv
**Version: 1.0**

Yet another role-based authentication service, written in Go.

Authsv communicates via JSON over HTTPS (Port 3890).

### Usage
```
authsv server|setup [options...]
```
Server options:
```
  -cf string
    	Specify certificate file
  -df string
    	Specify database file
  -kf string
    	Specify key file
```
Setup options:
```
  -f string
    	Specify output file
  -n string
    	Specify initial administrator (default "admin")
  -r string
    	Specify roles for administrator (default "auth.admin")
  -s string
    	Specify secret for administrator (default "password")
```
