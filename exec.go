package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

import "github.com/shirou/gopsutil/process"

type OutputStream struct {
	streamChan chan string
	bufSize    int
	buf        []byte
	lastChar   int
}

var (
	ErrLineBufferOverflow = errors.New("line buffer overflow")
)


// Write makes OutputStream implement the io.Writer interface.
func (rw *OutputStream) Write(p []byte) (n int, err error) {
	n = len(p) // end of buffer
	firstChar := 0

	for {
		newlineOffset := bytes.IndexByte(p[firstChar:], '\n')
		if newlineOffset < 0 {
			break // no newline in stream, next line incomplete
		}

		// End of line offset is start (nextLine) + newline offset. Like bufio.Scanner,
		// we allow \r\n but strip the \r too by decrementing the offset for that byte.
		lastChar := firstChar + newlineOffset // "line\n"
		if newlineOffset > 0 && p[newlineOffset-1] == '\r' {
			lastChar -= 1 // "line\r\n"
		}

		// Send the line, prepend line buffer if set
		var line string
		if rw.lastChar > 0 {
			line = string(rw.buf[0:rw.lastChar])
			rw.lastChar = 0 // reset buffer
		}
		line += string(p[firstChar:lastChar])
		rw.streamChan <- line // blocks if chan full

		// Next line offset is the first byte (+1) after the newline (i)
		firstChar += newlineOffset + 1
	}

	if firstChar < n {
		remain := len(p[firstChar:])
		bufFree := len(rw.buf[rw.lastChar:])
		if remain > bufFree {
			var line string
			if rw.lastChar > 0 {
				line = string(rw.buf[0:rw.lastChar])
			}
			line += string(p[firstChar:])
			err = ErrLineBufferOverflow
			n = firstChar
			return // implicit
		}
		copy(rw.buf[rw.lastChar:], p[firstChar:])
		rw.lastChar += remain
	}

	return // implicit
}

// NewOutputStream creates a new streaming output on the given channel.
func NewOutputStream(streamChan chan string) *OutputStream {
	out := &OutputStream{
		streamChan: streamChan,
		bufSize:    16384,
		buf:        make([]byte, 16384),
		lastChar:   0,
	}
	return out
}

func TestCheckStream(command string,timeout int,output bool) {
	stdoutChan := make(chan string, 100)
	if output {
		go func() {
			for line := range stdoutChan {
				gbk, _ := GbkToUtf8([]byte(line))
				fmt.Print("Output: ")
				fmt.Println(string(gbk))
			}
		}()
	}
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "cmd.exe", "/c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	stdout := NewOutputStream(stdoutChan)
	cmd.Stdout = stdout
	if err := cmd.Start(); err != nil {
		return
	}
	fmt.Println("PID: "+strconv.Itoa(cmd.Process.Pid))
	waitChan := make(chan struct{}, 1)
	defer close(waitChan)

	// 超时杀掉进程组 或正常退出
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Timeout Exit: "+strconv.Itoa(cmd.Process.Pid))
			p, _ := process.NewProcess(int32(cmd.Process.Pid)) // Specify process id of parent
			// handle error
			vp,_ := p.Children()
			var v *process.Process
			for _, v = range vp {
				_ = v.Kill() // Kill each child
				// handle error
			}
			p.Kill() // Kill the parent process
		case <-waitChan:
			fmt.Println("Success Exit: "+strconv.Itoa(cmd.Process.Pid))
		}
	}()

	if err := cmd.Wait(); err != nil {
		return
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


func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}


func main(){
	fmt.Print("set timeout: ")
	timeout,err := strconv.Atoi(getInput())
	if err != nil{
		timeout = 999999999
	}else if timeout == 0{
		timeout = 999999999
	}
	fmt.Println("Get Output: t/f?")
	fmt.Print("Output: ")
	output := true
	op := getInput()
	if op == "f"{
		output = false
	}else{
		output = true
	}
	for {
		cmd := getInput()
		go TestCheckStream(cmd,timeout,output)
	}
}