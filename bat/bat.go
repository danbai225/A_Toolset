package bat

var (
	Ver         bool
	Form        bool
	Pretty      bool
	Download    bool
	InsecureSSL bool
	Auth        string
	Proxy       string
	PrintV      string
	PrintOption uint8
	Body        string
	Bench       bool
	BenchN      int
	BenchC           int
	Isjson           *bool
	Method           *string
	URL              *string
	Jsonmap          map[string]interface{}
	ContentJsonRegex = `application/(.*)json`
)

func init() {
	json:=true
	Isjson=&json
}
const (
	Version              = "0.1.0"
	PrintReqHeader uint8 = 1
	PrintReqBody
	PrintRespHeader
	PrintRespBody
)