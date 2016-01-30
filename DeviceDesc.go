package upnp

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

type DeviceDesc struct {
	upnp *Upnp
}

func (this *DeviceDesc) Send() (bool, error) {
	request, err := this.BuildRequest()
	if err != nil {
		return false, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return false, err
	}
	resultBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	if response.StatusCode == 200 {
		this.resolve(string(resultBody))
		return true, nil
	}
	return false, nil
}
func (this *DeviceDesc) BuildRequest() (*http.Request, error) {
	//请求头
	header := http.Header{}
	header.Set("Accept", "text/html, image/gif, image/jpeg, *; q=.2, */*; q=.2")
	header.Set("User-Agent", "preston")
	header.Set("Host", this.upnp.Gateway.Host)
	header.Set("Connection", "keep-alive")

	//请求
	request, err := http.NewRequest("GET", "http://" + this.upnp.Gateway.Host + this.upnp.Gateway.DeviceDescUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header = header
	// request := http.Request{Method: "GET", Proto: "HTTP/1.1",
	// 	Host: this.upnp.Gateway.Host, Url: this.upnp.Gateway.DeviceDescUrl, Header: header}
	return request, nil
}

func (this *DeviceDesc) resolve(resultStr string) {
	inputReader := strings.NewReader(resultStr)

	// 从文件读取，如可以如下：
	// content, err := ioutil.ReadFile("studygolang.xml")
	// decoder := xml.NewDecoder(bytes.NewBuffer(content))

	lastLabel := ""

	ISUpnpServer := false

	IScontrolURL := false
	var controlURL string //`controlURL`
	// var eventSubURL string //`eventSubURL`
	// var SCPDURL string     //`SCPDURL`

	decoder := xml.NewDecoder(inputReader)
	for t, err := decoder.Token(); err == nil && !IScontrolURL; t, err = decoder.Token() {
		switch token := t.(type) {
		// 处理元素开始（标签）
		case xml.StartElement:
			if ISUpnpServer {
				name := token.Name.Local
				lastLabel = name
			}

		// 处理元素结束（标签）
		case xml.EndElement:
		// log.Println("结束标记：", token.Name.Local)
		// 处理字符数据（这里就是元素的文本）
		case xml.CharData:
			//得到url后其他标记就不处理了
			content := string([]byte(token))

			//找到提供端口映射的服务
			if content == this.upnp.Gateway.ServiceType {
				ISUpnpServer = true
				continue
			}
			//urn:upnp-org:serviceId:WANIPConnection
			if ISUpnpServer {
				switch lastLabel {
				case "controlURL":

					controlURL = content
					IScontrolURL = true
				case "eventSubURL":
				// eventSubURL = content
				case "SCPDURL":
				// SCPDURL = content
				}
			}
		default:
		// ...
		}
	}
	this.upnp.CtrlUrl = controlURL
}
