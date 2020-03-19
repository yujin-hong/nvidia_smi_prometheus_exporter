# nvidia_smi_exporter

nvidia-smi metrics exporter for Prometheus

## Build
```
> go build -v nvidia_smi_exporter
```

## Run
```
> nohup ./nvidia_smi_exporter &
```
Default port is 9101

## Check
```
> sudo netstat -tnlp | grep gpu_exporter
```

### localhost:9101/metrics
```
temperature_gpu{gpu="TITAN V[0]"} 40
utilization_gpu{gpu="TITAN V[0]"} 0
utilization_memory{gpu="TITAN V[0]"} 0
memory_total{gpu="TITAN V[0]"} 12036
memory_free{gpu="TITAN V[0]"} 12036
memory_used{gpu="TITAN V[0]"} 0
gpu_using_pid{gpu="TITAN V[0]"} 0
```

### Exact command

To get temperature_gpu, utilization_gpu, utilization_memory, memory_total, memory_free, memory_used
```
nvidia-smi --query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used --format=csv,noheader,nounits
```

To get gpu_using_pid
```
nvidia-smi | grep 'python' | awk '{ print $3 }'
```

### Prometheus example config

```
- job_name: "gpu_exporter"
  static_configs:
  - targets: ['localhost:9101']
```

