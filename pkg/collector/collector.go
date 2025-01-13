package collector

const (
	namespace          = "process"
	cpuSubsystem       = "cpu"
	memorySubsystem    = "memory"
	networkSubsystem   = "network"
	ioSubsystem        = "io"
	ctxSwitchSubsystem = "ctx_switch"
	fdSubsystem        = "fds"
	threadSubsystem    = "threads"
)

var commonLabels = []string{"pid", "cmdline"}
