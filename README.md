gourl
=====

A simple clone of cURL client written in go. 

Install
=======

```
  go get github.com/andrix/gourl
```

Run
===

If $GOBIN it's on the $PATH, then 
```
 gourl [options] url
```

Changelog
=========

0.2
 - SSL connections implemented 
 - -k option implemented - non-secure connetions allowed
 - -o option implemented - write to file
 - -x option implemented - HTTP proxy support
 - -U option implemented - HTTP proxy basic auth
 - -u option implemented - Provide credentials for Server Authentication

0.1 - hello world!
 - lacks of all functionality
 - retrieve an url
 - -I option implemented
 - -A option implemented

