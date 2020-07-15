package bat

import (
	"log"
	"strings"
)

var methodList = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

func filter(args []string) []string {
	var i int
	if inSlice(strings.ToUpper(args[i]), methodList) {
		*Method = strings.ToUpper(args[i])
		i++
	} else if len(args) > 0 && *Method == "GET" {
		for _, v := range args[1:] {
			// defaults to either GET (with no request data) or POST (with request data).
			// Params
			strs := strings.Split(v, "=")
			if len(strs) == 2 {
				*Method = "POST"
				break
			}
			// files
			strs = strings.Split(v, "@")
			if len(strs) == 2 {
				*Method = "POST"
				break
			}
		}
	} else if *Method == "GET" && Body != "" {
		*Method = "POST"
	}
	if len(args) <= i {
		log.Fatal("Miss the URL")
	}
	*URL = args[i]
	i++
	return args[i:]
}
