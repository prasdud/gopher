package internal

import "golang.org/x/sys/unix"

type Uptime struct {
	Hours   int64
	Minutes int64
	Seconds int64
}

func (u *Uptime) MetricName() string {
	return "uptime"
}

func GetUptime() (*Uptime, error) {
	var info unix.Sysinfo_t
	if err := unix.Sysinfo(&info); err != nil {
		return nil, err
	}
	seconds := info.Uptime
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return &Uptime{
		Hours:   hours,
		Minutes: minutes,
		Seconds: secs,
	}, nil
}
