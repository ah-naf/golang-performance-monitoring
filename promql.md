## 1. Process CPU

### 1.1 CPU time rate

```promql
rate(myapp_process_cpu_seconds_total[5m])
```

**What it does**
Computes how many CPU-seconds the process has consumed per second over the last five minutes. Internally, `rate()` takes the difference in the counter and divides by the elapsed time:

* **Math example**:

  * Counter at T–5 min = 1 000 s
  * Counter now = 1 060 s
  * Δcounter = 60 s, Δtime = 300 s → 60 / 300 = 0.2 CPU-s/s (i.e. 20% of one core)

### 1.2 Instant CPU usage percentage

```promql
myapp_process_cpu_usage_percent
```

**What it does**
A snapshot gauge of the process’s current CPU utilization as a percentage of one core. E.g. a value of 75 means the process is using 75% of a single CPU.

---

## 2. Process Memory

### 2.1 Resident memory (RSS)

```promql
myapp_process_resident_memory_bytes
```

**What it does**
Shows how many bytes of physical RAM the process currently occupies.

### 2.2 Virtual memory

```promql
myapp_process_virtual_memory_bytes
```

**What it does**
Total virtual address space allocated to the process (including RSS, swap, and memory-mapped files).

---

## 3. System & Go Runtime

### 3.1 Uptime

```promql
time() - myapp_process_start_time_seconds
```

**What it does**
Subtracts the Unix timestamp when the process started from the current Unix time, yielding total seconds running.

* **Math example**:

  * `myapp_process_start_time_seconds` = 1 650 000 000
  * `time()` = 1 650 000 600 → uptime = 600 s (10 minutes)

### 3.2 Open file descriptors

```promql
myapp_process_open_fds
```

**What it does**
Number of files, sockets, etc., the process has open right now.

### 3.3 Thread count

```promql
myapp_process_threads
```

**What it does**
Number of OS threads currently in use by the process.

### 3.4 95ᵗʰ-percentile GC pause

```promql
histogram_quantile(
  0.95,
  sum(rate(myapp_go_runtime_gc_pause_duration_seconds_bucket[5m]))
    by (le)
)
```

**What it does**

1. Converts GC pause-duration bucket counters into per-second rates over 5 min.
2. Aggregates all buckets into one histogram.
3. Computes the 95ᵗʰ-percentile pause time (in seconds) from that histogram.

---

## 4. HTTP Metrics

> All HTTP metrics include labels `method`, `endpoint`, and (where defined) `status_code`.

### 4.1 Request rate

```promql
rate(myapp_http_requests_total[1m])
```

**What it does**
Average number of HTTP requests the process receives per second over the last minute.

* **Math example**:

  * Counter at 12:00 = 10 000
  * Counter at 12:01 = 10 300 → (10 300 – 10 000)/60 = 5 rps

### 4.2 95ᵗʰ-percentile request latency

```promql
histogram_quantile(
  0.95,
  sum(rate(myapp_http_request_duration_seconds_bucket[5m]))
    by (le, endpoint)
)
```

**What it does**
Finds the 95ᵗʰ-percentile of HTTP request durations (seconds) per endpoint over the last five minutes.

### 4.3 95ᵗʰ-percentile request size

```promql
histogram_quantile(
  0.95,
  sum(rate(myapp_http_request_size_bytes_bucket[5m]))
    by (le)
)
```

**What it does**
Estimates the 95ᵗʰ-percentile of incoming request payload sizes (bytes) over five minutes.

### 4.4 95ᵗʰ-percentile response size

```promql
histogram_quantile(
  0.95,
  sum(rate(myapp_http_response_size_bytes_bucket[5m]))
    by (le)
)
```

**What it does**
Estimates the 95ᵗʰ-percentile of outgoing response payload sizes (bytes) over five minutes.

---

## Broad Filtering

You can narrow any of the above queries by adding a label selector in `{…}`. For example:

* **Exact match**:

  ```promql
  rate(myapp_http_requests_total{method="POST"}[1m])
  ```
* **Regex match**:

  ```promql
  rate(myapp_http_requests_total{endpoint=~"/api/.*"}[1m])
  ```

Use `{label="value"}` to include only series where that label equals `value`, or `{label=~"regex"}` to match multiple values via regular expressions.
