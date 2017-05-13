# feserver
A server providing [F.E. examination](https://www.jitec.ipa.go.jp/1_11seido/fe.html) questions 
which have been appeared in the past.
The past questions are obtained from http://www.fe-siken.com.
Thanks to the great Web site!


## Installation

To install the command, run: 

```
go get -u https://github.com/mzki/feserver
```

To start server process, just type
```
feserver
```
then access `localhost:8080/r-question.json`,
you can get a json response which contains F.E. question randomly selected.   

## API

`github.com/mzki/feserver/src` provides API for getting the F.E. questions.

## License

The BSD 3-Clause license, the same as the [Go](https://golang.org/).
