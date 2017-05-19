# feserver

A server providing [F.E. examination](https://www.jitec.ipa.go.jp/1_11seido/fe.html) questions 
which have been appeared in the past.
The past questions are obtained from http://www.fe-siken.com.
Thanks to the great Web site!

feserver can also serve the questions from the other IPA examinations 
such as [A.P.](https://www.jitec.ipa.go.jp/1_11seido/ap.html) !

**IMPORTANT : Please use feserver personally. Copyrights for the past questions are reserved by 
http://www.fe-siken.com and [IPA](http://www.ipa.go.jp/index.html).**  

## Installation

To install the command, run: 

```
go get -u https://github.com/mzki/feserver
```

To start server process, just type
```
feserver
```

And then you will see the JSON response by 

```
curl http://localhost:8080/r-question.json
```

The JSON response contains F.E. question and its answer.
For more detail see JSON Response section.


## Web API

feserver provides the following Web APIs:

* `[server-address]/[sub-address]/r-question.json`

It returns json response which contains the question randomly selected.

* `[server-address]/[sub-address]/question.json?year=[year]&season=[haru|aki]&no=[no]`

It returns json response which contains the question specified by the query parameters.

## JSON Response 

The returned JSON response has:

* `question`: Question text
* `selections`: Selections for Answer.
* `answer`: Answer Character, ア, イ, ウ, エ.
* `explanation`: Explanation for the Answer.
* `hasImage`: question, selections, or answer contain some images. These might not be represented by only text.
* `url`: Source URL in which the question is retrieved.
* `error`: Error message. Empty message indicates non-error.

## Configuration

By default, feserver initially loads `config.toml` at the feserver's repository under `GOPATH`.
`config.toml` defines question source locations for serving the content.
See `config.toml` for more detail.

## Library

`src` directory provides the Go library for getting the F.E. questions or others.

Example:

```go
// make context with timeout for 3 second.
ctx := context.WithTimeout(context.Background(), 3*time.Second)

// get src.Response for the question randomly selected.
res, _ := src.GetRandom(ctx, src.MaxQueryRange)

```

You can create src.Getter with arbitrary source location.

```go
// construct Getter with A.P. examination source and 
// interval time for the request
g := src.NewGetter(src.AP, src.LeastIntervalTime)

// get A.P. question at H29, Spirng, No. 10.
res, _ := g.Get(context.Background(), Query{29, src.SeasonSpring, 10})

// get A.P. question randomly selected
res, _ = g.GetRandom(context.Background(), src.MaxQueryRange)
```


## License

The BSD 3-Clause license, the same as the [Go](https://golang.org/).
