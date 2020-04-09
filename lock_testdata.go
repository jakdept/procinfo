package procinfo

// +build linux darwin

var Testdata_Lock = []Lock{
	Lock{
		Priority:  41,
		ByteRange: false,
		Exclusive: true,
		PID:       5623,
		DevMajor:  0x00,
		DevMinor:  0x17,
		Inode:     810,
	},
	Lock{
		Priority:  42,
		ByteRange: true,
		Exclusive: false,
		PID:       5705,
		DevMajor:  0x00,
		DevMinor:  0x17,
		Inode:     1088,
	},
	Lock{
		Priority:  61,
		ByteRange: false,
		Exclusive: true,
		PID:       8819,
		DevMajor:  0xfd,
		DevMinor:  0x01,
		Inode:     23600664,
	},
}

var TestData_CheckInode_Processes = map[uint64][]Process{
	810: []Process{
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
	},
	1088: []Process{
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
	},
	23600664: []Process{
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
	},
}
