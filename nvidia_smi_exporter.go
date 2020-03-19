package main

import (
    "bytes"
    "encoding/csv"
    "fmt"
    "net/http"
    "log"
    "os"
    "os/exec"
    "strings"
)


// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used

func metrics(response http.ResponseWriter, request *http.Request) {
    out, err := exec.Command(
        "nvidia-smi",
        "--query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used",
        "--format=csv,noheader,nounits").Output()
    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    csvReader := csv.NewReader(bytes.NewReader(out))
    csvReader.TrimLeadingSpace = true
    records, err := csvReader.ReadAll()

    if err != nil {
        fmt.Printf("%s\n", err)
        return
    }

    cmd := "nvidia-smi | grep 'python' | awk '{ print $3 }'"
    out1, err1 := exec.Command("bash", "-c", cmd).Output()

    if err1 != nil {
        fmt.Printf("%s\n", err)
        return
    }

    csvReader1 := csv.NewReader(bytes.NewReader(out1))
    csvReader1.TrimLeadingSpace = true
    records1, err1 := csvReader1.ReadAll()
    // print(records1)

    if err1 != nil {
        fmt.Printf("%s\n", err1)
        return
    }

    metricList := []string {
        "temperature.gpu", "utilization.gpu",
        "utilization.memory", "memory.total", "memory.free", "memory.used", "gpu.using.pid"}

    result := ""
    result_gpu_name := ""
    result_nth := ""
    for _, row := range records {
        result_gpu_name = row[0]
        result_nth = row[1]
        name := fmt.Sprintf("%s[%s]", result_gpu_name, result_nth)
        for idx, value := range row[2:] {
            result = fmt.Sprintf(
                "%s%s{gpu=\"%s\"} %s\n", result,
                metricList[idx], name, value)
        }
    }
    fl := 0
    for _, row := range records1 {
        fl = 1
        name := fmt.Sprintf("%s[%s]", result_gpu_name, result_nth)
        result = fmt.Sprintf(
            "%s%s{gpu=\"%s\"} %s\n", result,
            metricList[6], name, row[0])
    }
    if fl==0 {
        name := fmt.Sprintf("%s[%s]", result_gpu_name, result_nth)
        result = fmt.Sprintf(
            "%s%s{gpu=\"%s\"} %s\n", result,
            metricList[6], name, "0")
    }

    // print(result)

    fmt.Fprintf(response, strings.Replace(result, ".", "_", -1))
}

func main() {
    addr := ":9101"
    if len(os.Args) > 1 {
        addr = ":" + os.Args[1]
    }

    http.HandleFunc("/metrics/", metrics)
    err := http.ListenAndServe(addr, nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

