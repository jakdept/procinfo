package procinfo

type Process struct {
	Pid  uint32
	PPid uint32
}

func getProcessByPid(pid uint32) Process {

	return Process{}
}
