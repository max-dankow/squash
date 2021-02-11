package local

import (
	"os"
	"os/exec"
)

type DotNetDbg struct {
}

func (d *DotNetDbg) GetRemoteConnectionCmd(plankName, plankNamespace, podName, podNamespace string, localPort, remotePort int) *exec.Cmd {
	// for DotNetDbg, we proxy through the debug container
	return GetPortForwardCmd(plankName, plankNamespace, localPort, remotePort)
}

func (d *DotNetDbg) GetEditorRemoteConnectionCmd(plankName, plankNamespace, podName, podNamespace string, remotePort int) string {
	// for DotNetDbg, we proxy through the debug container
	return getPortForwardWithRandomLocalCmd(plankName, plankNamespace, remotePort)
}

func (d *DotNetDbg) GetDebugCmd(localPort int) *exec.Cmd {
	// todo remove
	cmd := exec.Command("echo", "hello!")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd
}

func (d *DotNetDbg) ExpectRunningPlank() bool {
	return true
}

func (d *DotNetDbg) WindowsSupportWarning() string {
	return "Squash does not currently support the DotNetDbg interactive terminal on Windows. Please use the vscode extension or pass the --machine flag to squashctl."
}
