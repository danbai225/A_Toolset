package cmd

import (
	"bufio"
	"encoding/json"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/gfile"
	"github.com/spf13/cobra"
	"log"
	"os"
	"p00q.cn/A_Toolset/common"
	"p00q.cn/A_Toolset/utils"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	f2fPort      string
	f2fSend      bool   //是否为发送端 默认是的
	f2fIp        string //接收端IP
	f2fIsConn    int32  = 0
	f2fConn      *gtcp.Conn
	f2fPath      string
	f2fFiles     []File
	f2fFileIndex = 0
	f2fPathMap   = make(map[string]string)
	f2fFilesLog  bool
)

const (
	//目录结构
	Structure = "Structure"
	GetFile   = "GetFile"
	Over      = "Over"
)

type File struct {
	Path  string
	Data  []byte
	Size  int64
	Md5   string
	Begin int64
}

// tcpCmd represents the serve command
var f2fCmd = &cobra.Command{
	Use:   "f2f",
	Short: "f2f文件发送",
	Long:  `f2f将本机文件发送到接收端`, Example: "使用例子： A f2f D:\\test \n A f2f D:\\test -s=false ",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			path := args[0]
			var err error
			f2fPath = gfile.Abs(path)
			ips := utils.GetLANIps()
			print("本机局域网IPs: ")
			for _, ip := range ips {
				print(ip + " ")
			}
			if f2fSend {
				//发送端
				if f2fIp == "" {
					println("\n未指定IP,自动扫描局域网段,如果失败请-i加上连接IP")
					lanIps := utils.GetLANIps()
					for _, ip := range lanIps {
						atomic.AddInt32(&f2fIsConn, 1)
						go func(ip string) {
							split := strings.Split(ip, ".")
							for i := 1; i < 256; i++ {
								if atomic.LoadInt32(&f2fIsConn) < 0 {
									break
								}
								TempConn, err := gtcp.NewConn(split[0]+"."+split[1]+"."+split[2]+"."+strconv.Itoa(i)+":"+f2fPort, time.Millisecond*100)
								if TempConn != nil && err == nil {
									write, err := TempConn.Write([]byte("test\n"))
									if write > 0 && err == nil {
										f2fIp = split[0] + "." + split[1] + "." + split[2] + "." + strconv.Itoa(i)
										atomic.AddInt32(&f2fIsConn, -100*10000)
										f2fConn = TempConn
										break
									}
								}
							}
							atomic.AddInt32(&f2fIsConn, -1)
						}(ip)
					}
					for atomic.LoadInt32(&f2fIsConn) > 0 {
						time.Sleep(time.Millisecond * 100)
					}
				} else {
					f2fConn, err = gtcp.NewConn(f2fIp + ":" + f2fPort)
				}
				if f2fConn == nil || err != nil {
					println("\n未发现接收端,请开启接收端并指定IP和端口")
					return
				} else {
					println("\n连接成功!" + f2fIp + ":" + f2fPort)
					//开始发送
					send()
				}
			} else {
				//接收端
				gtcp.NewServer(":"+f2fPort, receive).Run()
			}
		} else {
			println(cmd.UsageString())
		}
	},
}

func init() {
	rootCmd.AddCommand(f2fCmd)
	f2fCmd.Flags().StringVarP(&f2fPort, "port", "p", "2252", "服务端口(双方一致或默认)")
	f2fCmd.Flags().StringVarP(&f2fIp, "ip", "i", "", "接收端IP")
	f2fCmd.Flags().BoolVarP(&f2fSend, "send", "s", true, "是否为发送端")
	f2fCmd.Flags().BoolVarP(&f2fFilesLog, "log", "l", false, "是否显示传输文件Log(文件数量多最好关闭)")
}

//发送目录结构
func sendStructure() {
	println("开始扫描文件")
	files := make([]File, 0)
	if utils.IsDir(f2fPath) {
		filesP, err := utils.GetAllFiles(f2fPath, "", true)
		if err == nil {
			for _, path := range filesP {
				files = append(files, File{Path: path})
			}
			for i, file := range files {
				files[i].Path = strings.ReplaceAll(file.Path, f2fPath, "")
				f2fPathMap[files[i].Path] = file.Path
				files[i].Size = utils.GetFileSize(file.Path)
				md5, _ := utils.CalcFileMD5(file.Path)
				files[i].Md5 = md5
			}

		} else {
			println("路径无效")
		}
	} else {
		md5, _ := utils.CalcFileMD5(f2fPath)
		file := File{Path: "", Size: utils.GetFileSize(f2fPath), Md5: md5}
		f2fPathMap[file.Path] = f2fPath
		files = append(files, file)
	}
	println("共需要发送", len(files), "个文件")
	common.SendStruct(f2fConn, Structure, "", files)
}
func send() {
	//发送目录结构
	sendStructure()
	defer f2fConn.Close()
	for {
		msg, err := bufio.NewReader(f2fConn).ReadString('\n')
		if err != nil {
			log.Printf("")
			break
		} else {
			msg = strings.Replace(msg, "\n", "", -1)
			b := []byte(msg)
			m := common.Msg{}
			err := json.Unmarshal(b, &m)
			if err == nil {
				switch m.Type {
				case GetFile:
					file := File{}
					json.Unmarshal([]byte(m.Data.(string)), &file)
					Read(&file)
					common.SendStruct(f2fConn, GetFile, "", file)
				case Over:
					println("发送完成")
					return
				}
			}
		}
	}
	println("发送完成")
}

//创建文件
func createFiles() {
	println("正在创建文件")
	for _, file := range f2fFiles {
		create, err := gfile.Create(f2fPath + file.Path)
		if err != nil {
			println("创建文件失败", f2fPath+file.Path)
		} else {
			create.Close()
		}
	}
	println("共需要接收", len(f2fFiles), "个文件")
}

func receive(conn *gtcp.Conn) {
	defer conn.Close()
	println("\n连接成功!", conn.RemoteAddr().String())
	for {
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("")
			break
		} else {
			msg = strings.Replace(msg, "\n", "", -1)
			b := []byte(msg)
			m := common.Msg{}
			err := json.Unmarshal(b, &m)
			if err == nil {
				switch m.Type {
				case Structure:
					var files []File
					json.Unmarshal([]byte(m.Data.(string)), &files)
					f2fFiles = files
					createFiles()
					if len(f2fFiles) > f2fFileIndex {
						common.SendStruct(conn, GetFile, "", f2fFiles[f2fFileIndex])
					}
				case GetFile:
					json.Unmarshal([]byte(m.Data.(string)), &f2fFiles[f2fFileIndex])
					Write(&f2fFiles[f2fFileIndex])
					if f2fFiles[f2fFileIndex].Begin == f2fFiles[f2fFileIndex].Size {
						f2fFileIndex++
						//传输完成
						if f2fFileIndex == len(f2fFiles) {
							//验证文件完整性
							Vmd5()
							if f2fFileIndex == len(f2fFiles) {
								common.SendStruct(conn, Over, "", nil)
								conn.Close()
								println("接收完成")
								os.Exit(1)
							} else {
								common.SendStruct(conn, GetFile, "", f2fFiles[f2fFileIndex])
							}
						}
						common.SendStruct(conn, GetFile, "", f2fFiles[f2fFileIndex])
					} else {
						common.SendStruct(conn, GetFile, "", f2fFiles[f2fFileIndex])
					}
				}
			}
		}
	}
}
func Read(file *File) {
	utils.CallClear()
	println("正在传输:" + file.Path)
	if file.Begin+1024*1024*8 < file.Size {
		file.Data = gfile.GetBytesByTwoOffsetsByPath(f2fPathMap[file.Path], file.Begin, file.Begin+1024*1024*8)
	} else {
		file.Data = gfile.GetBytesByTwoOffsetsByPath(f2fPathMap[file.Path], file.Begin, file.Size)
	}

}
func Vmd5() {
	f2fFilesErr := make([]File, 0)
	for _, file := range f2fFiles {
		md5, err := utils.CalcFileMD5(f2fPath + file.Path)
		if err == nil {
			if file.Md5 != md5 {
				file.Begin = 0
				f2fFilesErr = append(f2fFilesErr, file)
				gfile.Remove(f2fPath + file.Path)
			}
		}
	}
	if len(f2fFilesErr) > 0 {
		f2fFileIndex = 0
		f2fFiles = f2fFilesErr
	}
}
func Write(file *File) {
	utils.CallClear()
	println("正在传输:" + file.Path)
	f, err := os.OpenFile(f2fPath+file.Path, os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		println("打开文件失败", f2fPath+file.Path)
	} else {
		defer f.Close()
		l, _ := f.Write(file.Data)
		file.Begin += int64(l)
		file.Data = nil
	}
}
