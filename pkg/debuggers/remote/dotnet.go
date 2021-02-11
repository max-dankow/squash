package remote

import (
	"os"
	"os/exec"
	log "github.com/sirupsen/logrus"
)

const (
	SshPort  = 2222
)

type DotNetDbg struct {
}

type DotNetDbgLiveDebugSession struct {
	port    int
	process *os.Process
	cmd     *exec.Cmd
}

func (d *DotNetDbgLiveDebugSession) Detach() error {
	d.process.Kill()
	return nil
}

func (d *DotNetDbgLiveDebugSession) Port() int {
	return d.port
}

func (d *DotNetDbgLiveDebugSession) HostType() DebugHostType {
	return DebugHostTypeTarget
}

func (d *DotNetDbgLiveDebugSession) Cmd() *exec.Cmd {
	return d.cmd
}

func (d *DotNetDbg) Attach(pid int) (DebugServer, error) {
	
	cmd1 := exec.Command("service", "ssh", "restart")
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr

	cmd1.Start()

	log.Debug("StartDebugServer called")
	cmd := exec.Command("sleep", "infinity")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()

	if err != nil {
		log.WithField("err", err).Error("Failed to start stub")
		return nil, err
	}

	dnls := &DotNetDbgLiveDebugSession{
		port:    SshPort,
		process: cmd.Process,
		// store the cmd so we can Wait() for its compltion in the proxy
		cmd: cmd,
	}
	return dnls, nil
}
