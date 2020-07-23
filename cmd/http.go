/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"p00q.cn/A_Toolset/bat"
)

var GET = "GET"

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:     "http [flags] [METHOD] bat.URL [ITEM [ITEM]]",
	Short:   "HTTP交互命令行",
	Long:    `详情： https://github.com/astaxie/bat`,
	Example: "A http baidu.com",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			print(cmd.UsageString())
			os.Exit(1)
		}
		if len(args) == 1 {
			bat.URL = &args[0]
			bat.Method = &GET
		} else {
			bat.URL = &args[1]
			bat.Method = &args[0]
		}
		if bat.Ver {
			fmt.Println("Version:", bat.Version)
			os.Exit(2)
		}
		parsePrintOption(bat.PrintV)
		if bat.PrintOption&bat.PrintReqBody != bat.PrintReqBody {
			bat.DefaultSetting.DumpBody = false
		}
		var stdin []byte
		if runtime.GOOS != "windows" {
			fi, err := os.Stdin.Stat()
			if err != nil {
				panic(err)
			}
			if fi.Size() != 0 {
				stdin, err = ioutil.ReadAll(os.Stdin)
				if err != nil {
					log.Fatal("Read from Stdin", err)
				}
			}
		}

		if *bat.URL == "" {
			fmt.Println(cmd.UsageString())
		}
		if strings.HasPrefix(*bat.URL, ":") {
			urlb := []byte(*bat.URL)
			if *bat.URL == ":" {
				*bat.URL = "http://localhost/"
			} else if len(*bat.URL) > 1 && urlb[1] != '/' {
				*bat.URL = "http://localhost" + *bat.URL
			} else {
				*bat.URL = "http://localhost" + string(urlb[1:])
			}
		}
		if !strings.HasPrefix(*bat.URL, "http://") && !strings.HasPrefix(*bat.URL, "https://") {
			*bat.URL = "http://" + *bat.URL
		}
		u, err := url.Parse(*bat.URL)
		if err != nil {
			log.Fatal(err)
		}
		if bat.Auth != "" {
			userpass := strings.Split(bat.Auth, ":")
			if len(userpass) == 2 {
				u.User = url.UserPassword(userpass[0], userpass[1])
			} else {
				u.User = url.User(bat.Auth)
			}
		}
		*bat.URL = u.String()
		httpreq := bat.GetHTTP(*bat.Method, *bat.URL, args)
		if u.User != nil {
			password, _ := u.User.Password()
			httpreq.GetRequest().SetBasicAuth(u.User.Username(), password)
		}
		// Insecure SSL Support
		if bat.InsecureSSL {
			httpreq.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
		}
		// Proxy Support
		if bat.Proxy != "" {
			purl, err := url.Parse(bat.Proxy)
			if err != nil {
				log.Fatal("Proxy Url parse err", err)
			}
			httpreq.SetProxy(http.ProxyURL(purl))
		} else {
			eurl, err := http.ProxyFromEnvironment(httpreq.GetRequest())
			if err != nil {
				log.Fatal("Environment Proxy Url parse err", err)
			}
			httpreq.SetProxy(http.ProxyURL(eurl))
		}
		if bat.Body != "" {
			httpreq.Body(bat.Body)
		}
		if len(stdin) > 0 {
			var j interface{}
			d := json.NewDecoder(bytes.NewReader(stdin))
			d.UseNumber()
			err = d.Decode(&j)
			if err != nil {
				httpreq.Body(stdin)
			} else {
				httpreq.JsonBody(j)
			}
		}

		// AB bench
		if bat.Bench {
			httpreq.Debug(false)
			bat.RunBench(httpreq)
			return
		}
		res, err := httpreq.Response()
		if err != nil {
			log.Fatalln("can't get the url", err)
		}

		// download file
		if bat.Download {
			var fl string
			if disposition := res.Header.Get("Content-Disposition"); disposition != "" {
				fls := strings.Split(disposition, ";")
				for _, f := range fls {
					f = strings.TrimSpace(f)
					if strings.HasPrefix(f, "filename=") {
						// Remove 'filename='
						f = strings.TrimLeft(f, "filename=")

						// Remove quotes and spaces from either end
						f = strings.TrimLeft(f, "\"' ")
						fl = strings.TrimRight(f, "\"' ")
					}
				}
			}
			if fl == "" {
				_, fl = filepath.Split(u.Path)
			}
			if fl == "" {
				fl = u.Host
			}
			println(u.Path)
			println(fl)
			fd, err := os.OpenFile(fl, os.O_RDWR|os.O_CREATE, 0666)
			if err != nil {
				log.Fatal("can't create file", err)
			}
			if runtime.GOOS != "windows" {
				fmt.Println(bat.Color(res.Proto, bat.Magenta), bat.Color(res.Status, bat.Green))
				for k, v := range res.Header {
					fmt.Println(bat.Color(k, bat.Gray), ":", bat.Color(strings.Join(v, " "), bat.Cyan))
				}
			} else {
				fmt.Println(res.Proto, res.Status)
				for k, v := range res.Header {
					fmt.Println(k, ":", strings.Join(v, " "))
				}
			}
			fmt.Println("")
			contentLength := res.Header.Get("Content-Length")
			var total int64
			if contentLength != "" {
				total, _ = strconv.ParseInt(contentLength, 10, 64)
			}
			fmt.Printf("Downloading to \"%s\"\n", fl)
			pb := bat.NewProgressBar(total)
			pb.Start()
			multiWriter := io.MultiWriter(fd, pb)
			_, err = io.Copy(multiWriter, res.Body)
			if err != nil {
				log.Fatal("Can't Write the body into file", err)
			}
			pb.Finish()
			defer fd.Close()
			defer res.Body.Close()
			return
		}

		if runtime.GOOS != "windows" {
			fi, err := os.Stdout.Stat()
			if err != nil {
				panic(err)
			}
			if fi.Mode()&os.ModeDevice == os.ModeDevice {
				var dumpHeader, dumpBody []byte
				dump := httpreq.DumpRequest()
				dps := strings.Split(string(dump), "\n")
				for i, line := range dps {
					if len(strings.Trim(line, "\r\n ")) == 0 {
						dumpHeader = []byte(strings.Join(dps[:i], "\n"))
						dumpBody = []byte(strings.Join(dps[i:], "\n"))
						break
					}
				}
				if bat.PrintOption&bat.PrintReqHeader == bat.PrintReqHeader {
					fmt.Println(bat.ColorfulRequest(string(dumpHeader)))
					fmt.Println("")
				}
				if bat.PrintOption&bat.PrintReqBody == bat.PrintReqBody {
					if string(dumpBody) != "\r\n" {
						fmt.Println(string(dumpBody))
						fmt.Println("")
					}
				}
				if bat.PrintOption&bat.PrintRespHeader == bat.PrintRespHeader {
					fmt.Println(bat.Color(res.Proto, bat.Magenta), bat.Color(res.Status, bat.Green))
					for k, v := range res.Header {
						fmt.Printf("%s: %s\n", bat.Color(k, bat.Gray), bat.Color(strings.Join(v, " "), bat.Cyan))
					}
					fmt.Println("")
				}
				if bat.PrintOption&bat.PrintRespBody == bat.PrintRespBody {
					body := bat.FormatResponseBody(res, httpreq, bat.Pretty)
					fmt.Println(bat.ColorfulResponse(body, res.Header.Get("Content-Type")))
				}
			} else {
				body := bat.FormatResponseBody(res, httpreq, bat.Pretty)
				_, err = os.Stdout.WriteString(body)
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			var dumpHeader, dumpBody []byte
			dump := httpreq.DumpRequest()
			dps := strings.Split(string(dump), "\n")
			for i, line := range dps {
				if len(strings.Trim(line, "\r\n ")) == 0 {
					dumpHeader = []byte(strings.Join(dps[:i], "\n"))
					dumpBody = []byte(strings.Join(dps[i:], "\n"))
					break
				}
			}
			if bat.PrintOption&bat.PrintReqHeader == bat.PrintReqHeader {
				fmt.Println(string(dumpHeader))
				fmt.Println("")
			}
			if bat.PrintOption&bat.PrintReqBody == bat.PrintReqBody {
				fmt.Println(string(dumpBody))
				fmt.Println("")
			}
			if bat.PrintOption&bat.PrintRespHeader == bat.PrintRespHeader {
				fmt.Println(res.Proto, res.Status)
				for k, v := range res.Header {
					fmt.Println(k, ":", strings.Join(v, " "))
				}
				fmt.Println("")
			}
			if bat.PrintOption&bat.PrintRespBody == bat.PrintRespBody {
				body := bat.FormatResponseBody(res, httpreq, bat.Pretty)
				fmt.Println(body)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.Flags().BoolVarP(&bat.Ver, "version", "v", false, "输出bat（内置http库）版本")
	httpCmd.Flags().BoolVarP(&bat.Pretty, "pretty", "p", false, "打印漂亮的Json格式")
	httpCmd.Flags().StringVar(&bat.PrintV, "print", "A", "打印请求和响应")
	httpCmd.Flags().BoolVarP(&bat.Form, "form", "f", false, "以表单形式提交")
	httpCmd.Flags().BoolVarP(&bat.Download, "download", "d", false, "下载url内容作为文件")
	httpCmd.Flags().BoolVarP(&bat.InsecureSSL, "insecure", "i", false, "允许在没有证书的情况下连接到SSL站点")
	httpCmd.Flags().StringVarP(&bat.Auth, "bat.Auth", "a", "", "HTTP身份验证用户名:密码，用户[:PASS]")
	httpCmd.Flags().StringVar(&bat.Proxy, "bat.Proxy", "", "代理主机和端口，代理bat.URL")
	httpCmd.Flags().BoolVarP(&bat.Bench, "bench", "b", false, "向bat.URL发送bench请求")
	httpCmd.Flags().IntVar(&bat.BenchN, "b.N", 1000, "要发起的请求数")
	httpCmd.Flags().IntVar(&bat.BenchC, "b.C", 100, "要并发运行的请求数。")
	httpCmd.Flags().StringVar(&bat.Body, "body", "", "原始数据作为主体发送")
	httpCmd.Flags().BoolVarP(bat.Isjson, "json", "j", true, "以JSON对象的形式发送数据")
}
func parsePrintOption(s string) {
	if strings.ContainsRune(s, 'A') {
		bat.PrintOption = bat.PrintReqHeader | bat.PrintReqBody | bat.PrintRespHeader | bat.PrintRespBody
		return
	}

	if strings.ContainsRune(s, 'H') {
		bat.PrintOption |= bat.PrintReqHeader
	}
	if strings.ContainsRune(s, 'B') {
		bat.PrintOption |= bat.PrintReqBody
	}
	if strings.ContainsRune(s, 'h') {
		bat.PrintOption |= bat.PrintRespHeader
	}
	if strings.ContainsRune(s, 'b') {
		bat.PrintOption |= bat.PrintRespBody
	}
	return
}
