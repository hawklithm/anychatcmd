package jike

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/skratchdot/open-golang/open"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/48.0.2564.109 Safari/537.36"
)

var (
	SaveSubFolders = map[string]string{"webwxgeticon": "icons",
		"webwxgetheadimg": "headimgs",
		"webwxgetmsgimg":  "msgimgs",
		"webwxgetvideo":   "videos",
		"webwxgetvoice":   "voices",
		"_showQRCodeImg":  "qrcodes",
	}
	AppID             = "wx782c26e4c19acffb"
	Lang              = "zh_CN"
	LastCheckTs       = time.Now()
	SessionCreate     = "https://app.jike.ruguoapp.com/sessions.create"
	QrUrl1            = "jike://page.jk/web?displayHeader=false&displayFooter=false&url="
	QrUrl2            = "https://ruguoapp.com/account/scan?uuid="
	WaitForLoginUrl   = "https://app.jike.ruguoapp.com/sessions.wait_for_login?uuid="
	WaitForConfirmUrl = "https://app.jike.ruguoapp.com/sessions.wait_for_confirmation?uuid="
	TuringUrl         = "" //"http://www.tuling123.com/openapi/api"
	APIKEY            = "" //"391ad66ebad2477b908dce8e79f101e7"
	TUringUserId      = "" //"abc123"
)

type Jike struct {
	//User            User
	Root          string
	Debug         bool
	MessageNotify bool
	Uuid          string
	BaseUri       string
	RedirectedUri string
	Uin           string
	Sid           string
	Skey          string
	PassTicket    string
	DeviceId      string
	BaseRequest   map[string]string
	LowSyncKey    string
	SyncKeyStr    string
	SyncHost      string
	//SyncKey         SyncKey
	Users []string
	//InitContactList []User   //谈话的人
	//MemberList      []Member //
	//ContactList     []Member //好友
	GroupList []string //群
	//GroupMemberList []Member //群友
	//PublicUserList  []Member //公众号
	//SpecialUserList []Member //特殊账号

	AutoReplyMode bool //default false
	AutoOpen      bool
	Interactive   bool
	TotalMember   int
	TimeOut       int // 同步时间间隔   default:20
	MediaCount    int // -1
	SaveFolder    string
	QrImagePath   string
	Client        *http.Client
	//Request       *BaseRequest
	Log *log.Logger
	//MemberMap     map[string]Member
	ChatSet []string

	AutoReply    bool     //是否自动回复
	ReplyMsgs    []string // 回复的消息列表
	AutoReplySrc bool     //默认false,自动回复，列表。true调用AI机器人。
	lastCheckTs  time.Time
	SetCookie    []string
	AccessToken  string
}

func NewJike(logger *log.Logger) *Jike {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil
	}

	root, err := os.Getwd()
	transport := *(http.DefaultTransport.(*http.Transport))
	transport.ResponseHeaderTimeout = 1 * time.Minute
	transport.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return &Jike{
		Debug:         true,
		DeviceId:      "e123456789002237",
		AutoReplyMode: false,
		MessageNotify: true,
		Interactive:   false,
		AutoOpen:      false,
		MediaCount:    -1,
		Client: &http.Client{
			Transport: &transport,
			Jar:       jar,
			Timeout:   1 * time.Minute,
		},
		//Request:     new(BaseRequest),
		Root:        root,
		SaveFolder:  path.Join(root, "saved"),
		QrImagePath: filepath.Join(root, "qr.jpg"),
		Log:         logger,
		//MemberMap:   make(map[string]Member),
		SetCookie: []string{},
	}

}

func (w *Jike) onceTryLogin() (code string, err error) {
	code, tip := "", 1
	err = w.GetUUID()
	if err != nil {
		err = fmt.Errorf("get the uuid failed with error:%v", err)
		return "", err
	}
	err = w.GetQR()
	if err != nil {
		err = fmt.Errorf("创建二维码失败:%s", err.Error())
		return "", err
	}
	defer os.Remove(w.QrImagePath)
	w.Log.Println("扫描二维码登陆....")
	for code != "200" {
		w.RedirectedUri, code, tip, err = w.waitToLogin(w.Uuid, tip)
		if err != nil {
			err = fmt.Errorf("二维码登陆失败：%s", err.Error())
			return
		}
		if code == "TIMEOUT_EXPIRED" {
			break
		}
	}
	if code == "200" {
		code = ""
		for code != "200" {
			w.RedirectedUri, code, tip, err = w.waitToConfirm(w.Uuid, tip)
			if err != nil {
				err = fmt.Errorf("二维码登陆失败：%s", err.Error())
				return
			}
			if code == "TIMEOUT_EXPIRED" {
				break
			}
		}
	}
	return
}

func (w *Jike) WaitForLogin() (err error) {

	code := ""
	for code != "200" {
		if code, err = w.onceTryLogin(); err != nil {
			return
		}
	}
	return
}

func (w *Jike) waitToLogin(uuid string, tip int) (redirectUri, code string,
	rt int, err error) {
	var buf bytes.Buffer
	buf.WriteString(WaitForLoginUrl)
	buf.WriteString(w.Uuid)
	rt = tip
	response, err := w.Client.Get(buf.String())
	if err != nil {
		return
	}
	code = strconv.Itoa(response.StatusCode)
	switch code {
	case "404":
		reason := response.Header.Get("reason")
		if reason == "TIMEOUT_EXPIRED" {
			code = "TIMEOUT_EXPIRED"
			return
		}
	case "204":
		w.Log.Println("等待扫码")
	case "200":
		w.Log.Println("扫描成功，请在手机上点击确认登陆")
	default:
		err = errors.New("其它错误，请重启")

	}
	return
}

type ConfirmResponse struct {
	Confirmed    bool   `json:"confirmed"`
	Token        string `json:"token"`
	AccessToken  string `json:"x-jike-access-token"`
	RefreshToken string `json:"x-jike-refresh-token"`
}

type UuidResponse struct {
	Uuid string `json:"uuid"`
}

func (w *Jike) waitToConfirm(uuid string, tip int) (redirectUri, code string,
	rt int, err error) {
	var buf bytes.Buffer
	buf.WriteString(WaitForConfirmUrl)
	buf.WriteString(w.Uuid)
	rt = tip
	response, err := w.Client.Get(buf.String())
	if err != nil {
		return
	}
	code = strconv.Itoa(response.StatusCode)
	switch code {
	case "404":
		reason := response.Header.Get("reason")
		if reason == "TIMEOUT_EXPIRED" {
			code = "TIMEOUT_EXPIRED"
			return
		}
	case "200":
		defer response.Body.Close()
		respBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return "", code, rt, err
		}

		result := string(respBody)
		var data ConfirmResponse
		if err := json.Unmarshal([]byte(result), &data); err != nil {
			return "", code, rt, err
		}
		if data.Confirmed {
			w.Log.Println("登陆确认成功")
			w.AccessToken = data.AccessToken
		} else {
			w.Log.Println("登陆确认失败")
		}

	case "204":
		w.Log.Println("等待确认")
	default:
		err = errors.New("其它错误，请重启")

	}
	return
}

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

func (w *Jike) generateQR() (err error) {
	var buf = new(bytes.Buffer)
	buf.WriteString(QrUrl2)
	buf.WriteString(w.Uuid)

	url := encodeURIComponent(buf.String())

	buf = new(bytes.Buffer)
	buf.WriteString(QrUrl1)
	buf.WriteString(url)

	if qrCode, err := qr.Encode(buf.String(), qr.L, qr.Auto); err != nil {
		return err
	} else {
		// Scale the barcode to 200x200 pixels
		qrCode, err = barcode.Scale(qrCode, 200, 200)
		if err != nil {
			return err
		}

		// create the output file
		file, err := os.Create(w.QrImagePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// encode the barcode as png
		if err := png.Encode(file, qrCode); err != nil {
			return err
		}
		return nil
	}
}

func (w *Jike) GetQR() (err error) {
	if w.Uuid == "" {
		err = errors.New("no this uuid")
		return
	}

	w.generateQR()

	return open.Start(w.QrImagePath)
}

func (w *Jike) SetSynKey() {

}

func (w *Jike) AutoReplyMsg() string {
	if w.AutoReplySrc {
		return "" //not enabled
	} else {
		if len(w.ReplyMsgs) == 0 {
			return "未设置"
		}
		return w.ReplyMsgs[0]
	}

}

func (w *Jike) GetUUID() (err error) {
	params := url.Values{}
	datas := w.Get(SessionCreate, params)

	var r UuidResponse

	if err = json.Unmarshal([]byte(datas), &r); err != nil {
		return
	}

	fmt.Printf("%v\n", r)

	if r.Uuid != "" {
		w.Uuid = r.Uuid
		return
	} else {
		err = errors.New("get uuid failed")
		return
	}
}

func (w *Jike) Login() (err error) {
	//w.Log.Printf("the redirectedUri:%v", w.RedirectedUri)
	//
	//resp, err := w.Client.Get(w.RedirectedUri)
	//if err != nil {
	//	return
	//}
	//defer resp.Body.Close()
	//w.SetCookie = resp.Header["Set-Cookie"]
	//reader := resp.Body.(io.Reader)
	//if err = xml.NewDecoder(reader).Decode(w.Request); err != nil {
	//	return
	//}
	//
	//w.Request.DeviceID = w.DeviceId
	//
	//data, err := json.Marshal(Request{
	//	BaseRequest: w.Request,
	//})
	//if err != nil {
	//	return
	//}
	//
	//name := "webwxinit"
	//newResp := new(InitResp)
	//
	//index := strings.LastIndex(w.RedirectedUri, "/")
	//if index == -1 {
	//	index = len(w.RedirectedUri)
	//}
	//w.BaseUri = w.RedirectedUri[:index]
	//
	//apiUri := fmt.Sprintf("%s/%s?pass_ticket=%s&skey=%s&r=%d", w.BaseUri, name, w.Request.PassTicket, w.Request.Skey, int(time.Now().Unix()))
	//if err = w.Send(apiUri, bytes.NewReader(data), newResp); err != nil {
	//	return
	//}
	//w.Log.Printf("the webwxinit newResp:%#v", newResp)
	//for _, contact := range newResp.ContactList {
	//	w.InitContactList = append(w.InitContactList, contact)
	//}
	//
	//w.ChatSet = strings.Split(newResp.ChatSet, ",")
	//w.User = newResp.User
	//w.SyncKey = newResp.SyncKey
	//w.SyncKeyStr = ""
	//for i, item := range w.SyncKey.List {
	//
	//	if i == 0 {
	//		w.SyncKeyStr = strconv.Itoa(item.Key) + "_" + strconv.Itoa(item.Val)
	//		continue
	//	}
	//
	//	w.SyncKeyStr += "|" + strconv.Itoa(item.Key) + "_" + strconv.Itoa(item.Val)
	//
	//}
	//w.Log.Printf("the response:%+v\n", newResp)
	//w.Log.Printf("the sync key is %s\n", w.SyncKeyStr)
	//return
	return
}
