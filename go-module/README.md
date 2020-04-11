Go Module使用步骤：

1.创建项目

2.Goland配置

File->Settings...->Go->GOPATH指定全局的一个GOPATH即可，作为代码仓库或代码缓存

File->Settings...->Go->Go Modules启用

3.项目目录执行

go mod init

4.需要什么包，执行对应的go get，go module会自动同步

5.go mod vendor从代码仓库提取代码到项目

