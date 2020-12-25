package remote

import (
	"fmt"
	"os"
	"os/exec"
	"time"
	log "github.com/sirupsen/logrus"
)

type DotNetDbg struct {
}

type DotNetDbgLiveDebugSession struct {
	client  *rpc1.RPCClient
	port    int
	process *os.Process
	cmd     *exec.Cmd
}

func (d *DotNetDbgLiveDebugSession) Detach() error {
	d.client.Detach(false)
	d.process.Kill()
	return nil
}

func (d *DotNetDbgLiveDebugSession) Port() int {
	return d.port
}

func (d *DotNetDbgLiveDebugSession) HostType() DebugHostType {
	return DebugHostTypeClient
}

func (d *DotNetDbgLiveDebugSession) Cmd() *exec.Cmd {
	return d.cmd
}

func (d *DotNetDbg) attachTo(pid int) (*DotNetDbgLiveDebugSession, error) {
	cmd, port, err := d.startDebugServer(pid)
	if err != nil {
		return nil, err
	}
	// use rpc1 client for vscode extension support
	client := rpc1.NewClient(fmt.Sprintf("localhost:%d", port))
	dls := &DotNetDbgLiveDebugSession{
		client:  client,
		port:    port,
		process: cmd.Process,
		// store the cmd so we can Wait() for its compltion in the proxy
		cmd: cmd,
	}
	return dls, nil
}

func (d *DotNetDbg) Attach(pid int) (DebugServer, error) {
	return d.attachTo(pid)
}

func (d *DotNetDbg) startDebugServer(pid int) (*exec.Cmd, int, error) {
	log.WithField("pid", pid).Debug("StartDebugServer called")
	// cmd := exec.Command("dlv", "attach", fmt.Sprintf("%d", pid), "--listen=127.0.0.1:0", "--accept-multiclient=true", "--api-version=2", "--headless", "--log")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.WithFields(log.Fields{"cmd": cmd, "args": cmd.Args}).Debug("dotnet fake command")

	err := cmd.Start()
	if err != nil {
		log.WithField("err", err).Error("Failed to start dotnet fake command")
		return nil, 0, err
	}

	// todo: remove
	log.Debug("starting headless dlv for user started, trying to get port")
	time.Sleep(2 * time.Second)
	port, err := GetPort(cmd.Process.Pid)
	if err != nil {
		log.WithField("err", err).Error("can't get headless dlv port")
		cmd.Process.Kill()
		cmd.Process.Release()
		return cmd, 0, err
	}

	return cmd, port, nil
}
