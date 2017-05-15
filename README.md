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

## Web API

feserver provides the following Web APIs:

* `[server-URL]:[server-Port]/r-question.json`

It returns json response which contains F.E. question randomly selected.

* `[server-URL]:[server-Port]/question.json?year=[year]&season=[haru|aki]&no=[no]`

It returns json response which contains F.E. question specified by the query parameters.

### Response 

The returned json-response has:

* `question`: Question text
* `selections`: Selections for Answer.
* `answer`: Answer Character, ア, イ, ウ, エ.
* `explanation`: Explanation for the Answer.
* `hasImage`: question, selections, or answer contain some images. These might not be represented by only text.
* `url`: Source URL in which the question is retrieved.
* `error`: Error message. Empty message indicates non-error.

## Configuration

feserver initially loads `config.toml` as configuration file.
See `config.toml` at top directory for more detail.

## Library

`github.com/mzki/feserver/src` provides the Go library for getting the F.E. questions.
you can write program using it.

## License

The BSD 3-Clause license, the same as the [Go](https://golang.org/).
