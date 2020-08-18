package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "net/http"
    "regexp"
    "time"
)

type Config struct {
    Host string `json:"http_host"`
    Keyword [] struct {
        Code string `json:"code"`
        Echo string `json:"echo"`
    } `json:"keyword"`
}

type Message struct {
    MsgId int `json:"message_id"`
    GroupId int `json:"group_id"`
    Msg string `json:"message"`
    MsgType string `json:"message_type"`
    PostType string `json:"post_type"`
    RawMsg string `json:"raw_message"`
    SelfId int `json:"self_id"`
    UserId int `json:"user_id"`
    SubType string `json:"sub_type"`
    Time int `json:"time"`
    Sender struct{
        Card string `json:"card"`
        NickName string `json:"nickname"`
        Role string `json:"role"`
        Title string `json:"title"`
        UserId int `json:"user_id"`
    } `json:"sender"`
}

type SendMsgGroup struct{
    GroupId int `json:"group_id"`
    Msg string `json:"message"`
}

func main() {
    JsonParse := NewJsonStruct()
    v := Config{}
    //下面使用的是相对路径，config.json文件和main.go文件处于同一目录下
    JsonParse.Load("./config.json", &v)
    fmt.Println(v.Keyword)

    r := gin.Default()
    r.POST("/bot", func(c *gin.Context) {
        data, _ := ioutil.ReadAll(c.Request.Body)
        m := Message{}
        JsonParse.decode(data,&m)
        for _, item := range v.Keyword {
            ok, err := regexp.MatchString(item.Code,m.Msg)
            if err != nil {  //解释失败，返回nil
                fmt.Println("regexp err")
                return
            }
            if ok {
                data := SendMsgGroup{
                    GroupId: m.GroupId,
                    Msg: item.Echo,
                }
                Post(v.Host+"/send_group_msg",data)
                return
            }
        }
    })
    r.Run(":8080")
}

type JsonStruct struct {
}

func NewJsonStruct() *JsonStruct {
    return &JsonStruct{}
}

func (jst *JsonStruct) Load(filename string, v interface{}) {
    //ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return
    }
    jst.decode(data,v)
}

func (jst *JsonStruct) decode(data []byte, v interface{})  {
    //读取的数据为json格式，需要进行解码
    err := json.Unmarshal(data, v)
    if err != nil {
        return
    }
}


// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
// content：     请求放回的内容
func Post(url string, data interface{}) bool {

    // 超时时间：5秒
    client := &http.Client{Timeout: 5 * time.Second}
    jsonStr, _ := json.Marshal(data)
    resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    return true
}
