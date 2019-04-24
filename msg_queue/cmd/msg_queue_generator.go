package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

const msgQueueTemplate = `
type {{ .ClassName }} struct {
    // 消息队列
    msgQueue *msg_queue.MsgQueue
}

// Init{{ .ClassName }} {{ .ClassName }}构造函数
func Init{{ .ClassName }}() *{{ .ClassName }} {
    {{ .ObjectName }} := new({{ .ClassName }})

    // 初始化消息队列
	{{ .ObjectName }}.msgQueue = msg_queue.InitMsgQueue({{ .ObjectName }})

    return {{ .ObjectName }}
}

// Destroy{{ .ClassName }} {{ .ClassName }}析构函数
func Destroy{{ .ClassName }}({{ .ObjectName }} *{{ .ClassName }}) {
    // 销毁消息队列
	msg_queue.DestroyMsgQueue({{ .ObjectName }}.msgQueue)
	{{ .ObjectName }}.msgQueue = nil

    {{ .ObjectName }} = nil
}

// Start 开始运行
func ({{ .ObjectName }} *{{ .ClassName }}) Start() {
	{{ .ObjectName }}.msgQueue.Start()
}

// Stop 停止运行
func ({{ .ObjectName }} *{{ .ClassName }}) Stop() {
	{{ .ObjectName }}.msgQueue.Stop()
}

// OnStart 在消息循环启动前执行
func ({{ .ObjectName }} *{{ .ClassName }}) OnStart(q *msg_queue.MsgQueue) {
}

// OnStop 在消息循环停止前执行
func ({{ .ObjectName }} *{{ .ClassName }}) OnStop(q *msg_queue.MsgQueue) {
}

// OnMsgRecv 消息处理函数
func ({{ .ObjectName }} *{{ .ClassName }}) OnMsgRecv(q *msg_queue.MsgQueue, msg *msg_queue.Message) {
	switch msg.Msg {
	default:
		fmt.Println("{{ .ClassName }} no this msg")
	}
}

// OnDefaultRun 没有消息时循环执行
func ({{ .ObjectName }} *{{ .ClassName }}) OnDefaultRun(q *msg_queue.MsgQueue) {
}
`

var (
	className    string
	objectName   string
	filePathname string
)

func init() {
	flag.StringVar(&className, "c", "MyMsgQueue", "生成的类名")
	flag.StringVar(&objectName, "o", "q", "使用的方法对象名")
	flag.StringVar(&filePathname, "f", "my_msg_queue.go", "保存的文件名")
}

func main() {
	flag.Parse()

	param := struct {
		ClassName  string
		ObjectName string
	}{
		ClassName:  className,
		ObjectName: objectName,
	}

	t := template.Must(template.New("msg_queue").Parse(msgQueueTemplate))

	file, err := os.Create(filePathname)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	err = t.Execute(file, param)
	if err != nil {
		fmt.Println(err)
		return
	}
}
