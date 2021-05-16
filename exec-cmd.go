package main

// This example requires go-cmd v1.2 or newer

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/go-cmd/cmd"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"time"
)

var cmdlist []*cmd.Cmd

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func cmds(command string) {
	// Disable output buffering, enable streaming
	cmdOptions := cmd.Options{
		Buffered:  true,
		Streaming: true,
	}


	// Create Cmd with options
	var envCmd *cmd.Cmd
	sysType := runtime.GOOS

	if sysType == "windows" {
		envCmd = cmd.NewCmdOptions(cmdOptions, "cmd.exe", "/c", command)
	}else {
		envCmd = cmd.NewCmdOptions(cmdOptions, "/bin/bash", "-c", command)
	}

	cmdlist = append(cmdlist, envCmd)
	// Print STDOUT and STDERR lines streaming from Cmd
	doneChan := make(chan struct{})
	go func() {
		defer close(doneChan)
		// Done when both channels have been closed
		// https://dave.cheney.net/2013/04/30/curious-channels
		for envCmd.Stdout != nil || envCmd.Stderr != nil {
			lines := "\nOutput :\n"
			select {
			case line, open := <-envCmd.Stdout:
				if !open {
					envCmd.Stdout = nil
					continue
				}
				//fmt.Println(line)
				lines = lines + line + "\n"

			case line, open := <-envCmd.Stderr:
				if !open {
					envCmd.Stderr = nil
					continue
				}
				fmt.Println(os.Stderr, line)
			}
			gbk, _ := GbkToUtf8([]byte(lines))
			fmt.Println(string(gbk)+"\n")

		}
	}()


	// Run and wait for Cmd to return, discard Status
	//<-envCmd.Start()
	statusChan := envCmd.Start()
	_ = statusChan

	// Wait for goroutine to print everything
	<-doneChan
	DeleteSlice(envCmd)
}


func DeleteSlice(envCmd *cmd.Cmd){
	j := 0
	for _, val := range cmdlist {
		if val == envCmd {
			fmt.Println("Del: "+ strconv.Itoa(envCmd.Status().PID))
			break
		}
		j++
	}
	for i := j+1;i <len(cmdlist);i++{
		cmdlist[i-1] = cmdlist[i]
	}
	cmdlist = cmdlist[:len(cmdlist)-1]
}


func monitor(){
	var i *cmd.Cmd
	for {
		if len(cmdlist)>0 {
			for l := 0; l < len(cmdlist); l++ {
				i = cmdlist[l]
				str := "PID: " + strconv.Itoa(i.Status().PID)+"; CMD: " + i.Status().Cmd+"; Runtime: " + strconv.Itoa(int(i.Status().Runtime))
				fmt.Printf("%20s\n", str+"----------------------------------------------------")
			}
		}
		time.Sleep(1*time.Second)
	}

}

func getInput() string {
	//使用os.Stdin开启输入流
	//函数原型 func NewReader(rd io.Reader) *Reader
	//NewReader创建一个具有默认大小缓冲、从r读取的*Reader 结构见官方文档
	in := bufio.NewReader(os.Stdin)
	//in.ReadLine函数具有三个返回值 []byte bool error
	//分别为读取到的信息 是否数据太长导致缓冲区溢出 是否读取失败
	str, _, err := in.ReadLine()
	if err != nil {
		return err.Error()
	}
	return string(str)
}




func getInputByScanner() string {
	var str string
	//使用os.Stdin开启输入流
	in := bufio.NewScanner(os.Stdin)
	if in.Scan() {
		str = in.Text()
	} else {
		str = "Find input error"
	}
	return str
}


func main(){
	go monitor()
	for {
		go cmds(getInput())
	}
}
