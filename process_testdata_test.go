package procinfo

// +build linux darwin

var Testdata_Process = []Process{
	{
		Pid:          5623,
		OriginalName: "slack",
		State:        'S',
		Cwd:          "/proc/lol/fd",
		PPid:         2770,
		UserTime:     1,
		KernelTime:   286,
		Nice:         20,
		Cmdline:      []string{"/usr/lib/slack/slack"},
		Env: []string{
			"slack",
			"",
			"",
			"\n",
		},
	},
	{
		Pid:          5705,
		OriginalName: "chrome",
		State:        'S',
		Cwd:          "/proc/5715/fdinfo",
		PPid:         5680,
		UserTime:     0,
		KernelTime:   2,
		Nice:         20,
		Cmdline: []string{
			"/opt/google/chrome/chrome --type=zygote --enable-crash-reporter=4d5673fa-4274-4096-a4f8-0c3f692622fa,",
			"",
		},
		Env: []string{
			"",
			"",
			"",
			"",
			"",
			"\n",
		},
	},
	{
		Pid:          8819,
		OriginalName: "code",
		State:        'S',
		Cwd:          "/proc/notthere",
		PPid:         1,
		UserTime:     1,
		KernelTime:   578,
		Nice:         20,
		Cmdline:      []string{"/usr/share/code/code .", ""},
		Env: []string{
			"",
			"",
			"",
			"",
			"",
			"\n",
		},
	},
}
