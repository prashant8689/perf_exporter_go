package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

// Util Funcs

func sliceInsert(arr []string, pos int, elem string) []string {
	out := make([]string, len(arr)+1)
	copy(out[:pos], arr[:pos])
	out[pos] = elem
	copy(out[pos+1:], arr[pos:])
	return out
}

func ioBySize() map[string]float64 {
	command, _ := exec.Command("./scripts/bitesize", "15").Output()
	commandOutput := string(command)
	commandOutputArr := strings.Split(commandOutput, "\n")
	commandOutputArr = commandOutputArr[:len(commandOutputArr)-1][3:]
	countMap := make(map[string]float64)
	for _, element := range commandOutputArr {
		lineOutputArr := strings.Fields(string(element))
		if lineOutputArr[0] == "->" {
			lineOutputArr = sliceInsert(lineOutputArr, 0, "0")
		}
		if lineOutputArr[2] == ":" {
			lineOutputArr = sliceInsert(lineOutputArr, 2, "")
		}
		floatKey, _ := strconv.ParseFloat(string(lineOutputArr[0]), 64)
		strKey := strconv.Itoa(int(floatKey))
		floatVal, _ := strconv.ParseFloat(string(lineOutputArr[4]), 64)
		countMap[strKey] = floatVal
	}
	return countMap
}

func recordMetrics() {
	go func() {
		for {
			ioBySizeCountMap := ioBySize()
			for key, val := range ioBySizeCountMap {
				metricMap["blockSizeGT"+key].Set(val)
			}
		}
	}()
}

var (
	blockSizeGT0 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "disk_io_by_block_size_gt0",
		Help: "Disk I/O by Block Size greater than 0 KB and less than 1 KB (Kilobytes)",
	})
	blockSizeGT1 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "disk_io_by_block_size_gt1",
		Help: "Disk I/O by Block Size greater than 1 KB and less than 8 KB (Kilobytes)",
	})
	blockSizeGT8 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "disk_io_by_block_size_gt8",
		Help: "Disk I/O by Block Size greater than 8 KB and less than 64 KB (Kilobytes)",
	})
	blockSizeGT64 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "disk_io_by_block_size_gt64",
		Help: "Disk I/O by Block Size greater than 64 KB and less than 128 KB (Kilobytes)",
	})
	blockSizeGT128 = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "disk_io_by_block_size_gt128",
		Help: "Disk I/O by Block Size greater than 128 KB (Kilobytes)",
	})
	metricMap = map[string]prometheus.Gauge{"blockSizeGT0": blockSizeGT0, "blockSizeGT1": blockSizeGT1, "blockSizeGT8": blockSizeGT8, "blockSizeGT64": blockSizeGT64, "blockSizeGT128": blockSizeGT128}
)

func main() {
	recordMetrics()
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9100", nil)
}
