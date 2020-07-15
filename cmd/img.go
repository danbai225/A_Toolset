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
	"fmt"
	"github.com/spf13/cobra"
	"github.com/golang/glog"
	"github.com/googege/collie/mem"
	"github.com/googege/gotools/id"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)
var (
	root    string
	outPath string
	width   int
	quality int
	recursion bool
	isFile bool
)
// imgCmd represents the serve command
var imgCmd = &cobra.Command{
	Use:   "img",
	Short: "压缩图片",
	Long: `该命令用于图片压缩支持格式:Png和Jpg`,Example: "使用例子： A img -r ./imgs -o ./newimgs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("开始压缩...🚀")
		DataProcessing(root, outPath, width, quality)
		fmt.Println("压缩结束 ☕️")
	},
}

func init() {
	rootCmd.AddCommand(imgCmd)
	imgCmd.Flags().StringVarP(&root,"root","r","./","需要压缩的图片目录或文件位置(默认执行文件目录)")
	imgCmd.Flags().StringVarP(&outPath,"out","o","./out","输出的目录！！！（默认执行文件/out）")
	imgCmd.Flags().IntVarP(&width,"width","w",0,"图片的宽（0为不压缩大小）")
	imgCmd.Flags().IntVarP(&quality,"quality","q",75,"图片压缩质量（20-100）")
	imgCmd.Flags().BoolVarP(&recursion,"recursion","R",false,"是否递归目录(默认false)")


}
// get file's path
func retrieveData(root string) (value chan string, err chan error) {

	err = make(chan error, 1)
	value = make(chan string)
	if !IsFile(root){
		last3 := root[len(root)-1:]
		if last3!=string(os.PathSeparator) {
			root+=string(os.PathSeparator)
		}
	}else {
		isFile=true
	}
	go func() {
		defer close(value)
		err <- filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// if the file is noe regular, it mean the file is done,you should return
			if !info.Mode().IsRegular() {
				return nil
			}
			//是否递归
			if recursion||isFile{
				value <- path
			}else if root==strings.ReplaceAll(path,info.Name(),"") {
				value <- path
			}
			return nil
		})
	}()
	return
}

// get file send to a chan.
func ReceiveData(file chan string, value chan io.Reader, wg *sync.WaitGroup) {
	for v := range file {
		dif, err := mem.MemDifference()
		if err != nil {
			fmt.Println(err)
		}
		if dif > 0.2 {
			time.Sleep(time.Second >> 1)
			fmt.Println("waiting for mem less.")
		}
		fi, err := os.Open(v)
		if err != nil {
			fmt.Println(err)
		} else {

			value <- fi}
	}
	wg.Done()
}
// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}
// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断文件夹是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

type imageFile struct {
	img image.Image
	path string
}
// resize and create a new photo with only id name.
func DataProcessing(root string, outputFile string, wid int, q int) {
	reader := make(chan io.Reader)
	b := make(chan imageFile)
	c := make(chan imageFile)
	value, err := retrieveData(root)

		exist,errs := PathExists(outputFile)
		if errs != nil {
			fmt.Printf("获取文件夹错误![%v]\n", errs)
			return
		}

		if exist {
			fmt.Printf("有文件夹![%v]\n", outputFile)
		} else {
			fmt.Printf("没有文件夹![%v]\n", outputFile)
			// 创建文件夹
			err := os.Mkdir(outputFile, os.ModePerm)
			if err != nil {
				fmt.Printf("创建文件夹失败![%v]\n", err)
			} else {
				fmt.Printf("创建文件夹成功!\n")
			}
	}
	wg := new(sync.WaitGroup)
	wg.Add(2)
	for i := 0; i < 2; i++ {
		mark(i, "获取路径")
		go ReceiveData(value, reader, wg)
	}
	go func() {
		wg.Wait()
		close(reader)
	}()
	//
	wg1 := new(sync.WaitGroup)
	wg1.Add(32)
	for i := 0; i < 32; i++ {
		go func(i int) {
			defer wg1.Done()
			mark(i, "文件判断")
			for r := range reader {
				v, ok := r.(*os.File)
				if !ok {
					glog.Errorln("不是图片")
				}
				_, name1 := filepath.Split(v.Name())
				name := findName(name1)
				if name == "" && name1 != ".DS_Store" {
					glog.Errorln("没有文件")
				}
				img, err := isJpg(name, r)
				if err != nil {
					glog.Errorln(err)
				} else {
					//println(v.Name())
					b <- imageFile{img:img,path: v.Name()}
				}
			}
		}(i)
	}
	go func() {
		wg1.Wait()
		close(b)
	}()
	//
	wg2 := new(sync.WaitGroup)
	wg2.Add(32)
	for i := 0; i < 32; i++ {
		go func(i int) {
			mark(i, "压缩")
			defer wg2.Done()
			for i := range b {
				c <- imageFile{img:resize.Resize(uint(wid), 0, i.img, resize.NearestNeighbor),path: i.path}
			}
		}(i)
	}
	go func() {
		wg2.Wait()
		close(c)
	}()
	//
	wg3 := new(sync.WaitGroup)
	wg3.Add(32)
	for i := 0; i < 32; i++ {
		go func(i int) {
			mark(i, "处理图片。。。")
			defer wg3.Done()
			for i := range c {
				file, err := os.Create(outputFile+ string(os.PathSeparator) + filepath.Base(i.path))
				if err != nil {
					fmt.Println(err)
				}
				if q < 20 {
					q = 20
				}
				if err := jpeg.Encode(file, i.img, &jpeg.Options{q}); err != nil {
					glog.Errorln("图片处理出错:", err)
				}
			}
		}(i)
	}
	//
	if er := <-err; er != nil {
		fmt.Println(er)
	}
	//
	wg3.Wait()
}

// workNode is the computer's name if you have so many computers.
func onlyID() string {
	snow, err := id.NewSnowFlake(1)
	if err != nil {
		fmt.Println(err)
	}
	glog.V(1).Info("use snowFlake")
	return strconv.FormatInt(snow.GetID(), 10)
}
func onlyID1() string {
	u, err := id.NewUUID(id.VERSION_1, nil)
	if err != nil {
		glog.Error(err)
	}
	return u.String()
}
func findName(name string) string {
	v := name[len(name)-4:]
	v1 := name[len(name)-3:]
	if v == "jpeg" {
		return v
	}
	if v1 == "jpg" || v1 == "png" || v1 == "gif" {
		return v1
	}
	return ""
}
func isJpg(name string, r io.Reader) (image.Image, error) {
	name = strings.ToLower(name)
	switch name {
	case "jpeg", "jpg":
		return jpeg.Decode(r)
	case "png":
		return png.Decode(r)
	case "gif":
		return gif.Decode(r)
	default:
		return nil, fmt.Errorf("只能压缩jpeg jpg png和gif")
	}
}

func mark(i int, name string) {
	if i == 0 {
		fmt.Printf("%s is runing...\n", name)
	}
}