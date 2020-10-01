package main

import (
	"flag"
	"github.com/dxyinme/LukaClient/UserOperator"
	"log"
	"net/http"
)

var (
	addr = flag.String("addr", ":11501", "the core serve port")
)

func main() {
	flag.Parse()
	var (
		err error
	)
	http.HandleFunc("/Connect", UserOperator.Connect)
	if err = http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}
