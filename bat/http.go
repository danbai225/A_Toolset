package bat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/astaxie/bat/httplib"
)

var DefaultSetting = httplib.BeegoHttpSettings{
	ShowDebug:        true,
	UserAgent:        "bat/" + Version,
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	DumpBody:         true,
}

func GetHTTP(method string, url string, args []string) (r *httplib.BeegoHttpRequest) {
	r = httplib.NewBeegoRequest(url, method)
	r.Setting(DefaultSetting)
	r.Header("Accept-Encoding", "gzip, deflate")
	if *Isjson {
		r.Header("Accept", "application/json")
	} else if Form || method == "GET" {
		r.Header("Accept", "*/*")
	} else {
		r.Header("Accept", "application/json")
	}
	for i := range args {
		// Headers
		strs := strings.Split(args[i], ":")
		if len(strs) >= 2 {
			if strs[0] == "Host" {
				r.SetHost(strings.Join(strs[1:], ":"))
			}
			r.Header(strs[0], strings.Join(strs[1:], ":"))
			continue
		}
		// files
		strs = strings.SplitN(args[i], "@", 2)
		if !*Isjson && len(strs) == 2 {
			if !Form {
				log.Fatal("file upload only support in forms style: -f=true")
			}
			r.PostFile(strs[0], strs[1])
			continue
		}
		// Json raws
		strs = strings.SplitN(args[i], ":=", 2)
		if len(strs) == 2 {
			if strings.HasPrefix(strs[1], "@") {
				f, err := os.Open(strings.TrimLeft(strs[1], "@"))
				if err != nil {
					log.Fatal("Read File", strings.TrimLeft(strs[1], "@"), err)
				}
				content, err := ioutil.ReadAll(f)
				if err != nil {
					log.Fatal("ReadAll from File", strings.TrimLeft(strs[1], "@"), err)
				}
				var j interface{}
				err = json.Unmarshal(content, &j)
				if err != nil {
					log.Fatal("Read from File", strings.TrimLeft(strs[1], "@"), "Unmarshal", err)
				}
				Jsonmap[strs[0]] = j
				continue
			}
			Jsonmap[strs[0]] = toRealType(strs[1])
			continue
		}
		// Params
		strs = strings.SplitN(args[i], "=", 2)
		if len(strs) == 2 {
			if strings.HasPrefix(strs[1], "@") {
				f, err := os.Open(strings.TrimLeft(strs[1], "@"))
				if err != nil {
					log.Fatal("Read File", strings.TrimLeft(strs[1], "@"), err)
				}
				content, err := ioutil.ReadAll(f)
				if err != nil {
					log.Fatal("ReadAll from File", strings.TrimLeft(strs[1], "@"), err)
				}
				strs[1] = string(content)
			}
			if Form || method == "GET" {
				r.Param(strs[0], strs[1])
			} else {
				Jsonmap[strs[0]] = strs[1]
			}
			continue
		}
	}
	if !Form && len(Jsonmap) > 0 {
		r.JsonBody(Jsonmap)
	}
	return
}

func FormatResponseBody(res *http.Response, httpreq *httplib.BeegoHttpRequest, pretty bool) string {
	body, err := httpreq.Bytes()
	if err != nil {
		log.Fatalln("can't get the url", err)
	}
	fmt.Println("")
	match, err := regexp.MatchString(ContentJsonRegex, res.Header.Get("Content-Type"))
	if err != nil {
		log.Fatalln("failed to compile regex", err)
	}
	if pretty && match {
		var output bytes.Buffer
		err := json.Indent(&output, body, "", "  ")
		if err != nil {
			log.Fatal("Response Json Indent: ", err)
		}

		return output.String()
	}

	return string(body)
}
