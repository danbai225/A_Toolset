package itself

import (
	"bytes"
	"encoding/json"
	"github.com/gogf/gf/os/gfile"
	"io/ioutil"
	"os"
	"os/exec"
	"p00q.cn/A_Toolset/utils"
	"path/filepath"
)

var (
	appPath           = ""
	appFileSize int64 = 0
	DataLength  int64 = 0
	mapData     map[string]string
)

func Get(key string) string {
	str := mapData[key]
	if !utils.IsNil(str) {
		return str
	}
	return ""
}
func Put(key string, val string) {
	mapData[key] = val
	marshalMap()
}
func getDataLength() int64 {
	path := gfile.GetBytesByTwoOffsetsByPath(execPath(), execFileSize()-8, execFileSize())
	DataLength = utils.BytesToInt64(path)
	return DataLength
}
func Init() {
	loadMapData()
}
func marshalMap() {
	if utils.IsNil(mapData) {
		loadMapData()
	}
	data, err := json.Marshal(mapData)
	utils.Check(err)
	//新的数据长度
	oldLength := DataLength
	DataLength = int64(len(data) + 8)
	data = bytesCombine(data, utils.Int64ToBytes(DataLength))
	//读取当前文件全部数据
	readFile, err := ioutil.ReadFile(execPath())
	utils.Check(err)
	//追加新的数据
	readFile = bytesCombine(readFile[:appFileSize-oldLength], data)
	//替换文件
	_ = gfile.Rename(execPath(), execPath()+"-old")
	_ = gfile.PutBytesAppend(execPath()+"-new", readFile)
	_ = gfile.Rename(execPath()+"-new", execPath())
}

//BytesCombine 多个[]byte数组合并成一个[]byte
func bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

//加载map数据
func loadMapData() {
	if gfile.IsFile(execPath() + "-old") {
		_ = gfile.Remove(execPath() + "-old")
	}
	if getDataLength() > 8 {
		mapDataBytes := gfile.GetBytesByTwoOffsetsByPath(execPath(), execFileSize()-getDataLength(), execFileSize()-8)
		err := json.Unmarshal(mapDataBytes, &mapData)
		utils.Check(err)
	} else {
		mapData = make(map[string]string)
	}
}

/**
程序路径
*/
func execPath() string {
	if appPath == "" {
		file, err := exec.LookPath(os.Args[0])
		utils.Check(err)
		appPath, _ = filepath.Abs(file)
	}
	return appPath
}

func execFileSize() int64 {
	if appFileSize == 0 {
		fileInfo, err := os.Stat(execPath())
		utils.Check(err)
		appFileSize = fileInfo.Size()
	}
	return appFileSize
}
func Remove(key string) {
	delete(mapData, key)
}
