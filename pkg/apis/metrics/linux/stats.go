package linux

import (
	linux "github.com/c9s/goprocinfo/linux"
	"log"
)

/***
* Stats Regarding metrics we are interested in followin metrics
* - CPU usage (as a percentage)
* - Total memory
* - Available memory
* - Total disk
* - Available disk
 */
type Stats struct {
	MemStats  *linux.MemInfo
	DiskStats *linux.Disk
	CpuStats  *linux.CPUStat
	LoadStats *linux.LoadAvg
}

// MemTotalKb returns the total memory amount from the MemStats
func (stats *Stats) MemTotalKb() uint64 {
	return stats.MemStats.MemTotal
}

// MemAvailableKb returns the available memory from the MemStats
func (stats *Stats) MemAvailableKb() uint64 {
	return stats.MemStats.MemAvailable
}

// MemUsedKb returns the memory used from the MemStats
func (stats *Stats) MemUsedKb() uint64 {
	return stats.MemStats.MemTotal - stats.MemStats.MemAvailable
}

// MemUsedPercent returns the memory used from the MemStats in percent
func (stats *Stats) MemUsedPercent() uint64 {
	return stats.MemStats.MemAvailable / stats.MemStats.MemTotal
}

// DiskTotal returns the total Disk available from the DiskStats in percent
func (stats *Stats) DiskTotal() uint64 {
	return stats.DiskStats.All
}

// DiskUsed returns the Disk used from the DiskStats in percent
func (stats *Stats) DiskUsed() uint64 {
	return stats.DiskStats.Used
}

// DiskFree returns the Disk free from the DiskStats in percent
func (stats *Stats) DiskFree() uint64 {
	return stats.DiskStats.Free
}

// CpuUsage returns the cpu usage.
// https://stackoverflow.com/questions/5514119/accurately-calculating-cpu-utilization-in-linux-using-proc-stat
func (stats *Stats) CpuUsage() float64 {
	idle := stats.CpuStats.Idle + stats.CpuStats.IOWait
	nonIdle := stats.CpuStats.User + stats.CpuStats.IRQ + stats.CpuStats.SoftIRQ +
		stats.CpuStats.Nice + stats.CpuStats.System + stats.CpuStats.GuestNice
	total := nonIdle + idle

	if total == 0 {
		return 0.00
	}

	return (float64(total) - float64(idle)) / float64(total)
}

// GetStats returns the actual stats of the running task
func GetStats() *Stats {
	return &Stats{
		CpuStats:  GetCpuStats(),
		DiskStats: GetDiskStats(),
		MemStats:  GetMemInfo(),
		LoadStats: GetLoadAvg(),
	}
}

// GetMemInfo return information about linux memory
func GetMemInfo() *linux.MemInfo {
	memStats, err := linux.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Printf("Error reading from /proc/meminfo")
		return &linux.MemInfo{}
	}
	return memStats
}

// GetDiskStats return stats about the hardDisk
func GetDiskStats() *linux.Disk {
	disk, err := linux.ReadDisk("/")
	if err != nil {
		log.Printf("Error reading from disk")
		return &linux.Disk{}
	}
	return disk
}

// GetLoadAvg return infos about how last the task is.
func GetLoadAvg() *linux.LoadAvg {
	last, err := linux.ReadLoadAvg("/proc/loadavg")
	if err != nil {
		log.Printf("Error reading from /proc/loadavg")
		return &linux.LoadAvg{}
	}
	return last
}

// GetCpuStats return infos about the cpu(s) running the task.
func GetCpuStats() *linux.CPUStat {
	cpu, err := linux.ReadStat("/proc/stat")
	if err != nil {
		log.Printf("Error reading from /proc/stat")
		return &linux.CPUStat{}
	}
	return &cpu.CPUStatAll
}
