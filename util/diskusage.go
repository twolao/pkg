package util

import "syscall"

/*
state := util.DiskUsage("/")
	fmt.Printf("All=%dM, Free=%dM, Available=%dM, Used=%dM, Usage=%d%%",
		state.All/diskstate.MB, state.Free/diskstate.MB, state.Available/diskstate.MB, state.Used/diskstate.MB, 100*state.Used/state.All)

*/

const (
	// B 1bytes
	B = 1
	// KB 1024bytes
	KB = 1024 * B
	// MB 1024 * 1024bytes
	MB = 1024 * KB
	// GB 1024 * 1024 * 1024bytes
	GB = 1024 * MB
)

// DiskStatus 磁盘使用情况
type DiskStatus struct {
	All       uint64 `json:"all"`
	Used      uint64 `json:"used"`
	Available uint64 `json:"available"`
	Free      uint64 `json:"free"`
}

// DiskUsage 获取磁盘使用情况
func DiskUsage(drive string) (disk DiskStatus) {
	sf := syscall.Statfs_t{}
	err := syscall.Statfs(drive, &sf)
	if err != nil {
		return
	}
	disk.All = sf.Blocks * uint64(sf.Bsize)
	disk.Free = sf.Bfree * uint64(sf.Bsize)
	disk.Available = sf.Bavail * uint64(sf.Bsize)
	disk.Used = disk.All - disk.Free
	return disk
}
