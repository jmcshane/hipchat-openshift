package ocexec

import (
	"bytes"
	"io"
	"os/exec"
	"time"

	"github.com/jmcshane/http-server/util"
)

//OcExecute executes and oc command and allows passing output to grep
func OcExecute(args []string) (bytes.Buffer, bytes.Buffer, error) {
	pipeIndex := util.SliceIndex(len(args), func(i int) bool { return args[i] == "|" })
	if pipeIndex == -1 {
		return standardCommand(args)
	} else {
		return pipedCommand(args, pipeIndex)
	}
}

func pipedCommand(args []string, pipeIndex int) (bytes.Buffer, bytes.Buffer, error) {
	var out, stderr2 bytes.Buffer
	c1 := exec.Command("oc", args[0:pipeIndex]...)
	c2 := exec.Command(args[pipeIndex+1], args[pipeIndex+2:]...)
	r, w := io.Pipe()
	c1.Stdout = w
	c2.Stdin = r
	c2.Stderr = &stderr2
	c2.Stdout = &out
	c1.Start()
	c2.Start()
	c1.Wait()
	w.Close()
	err := c2.Wait()
	if err != nil {
		return out, stderr2, err
	}
	return out, stderr2, nil
}

func standardCommand(args []string) (bytes.Buffer, bytes.Buffer, error) {
	cmd := exec.Command("oc", args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	go func() {
		time.Sleep(3000)
		cmd.Process.Kill()
	}()
	return out, stderr, err
}
