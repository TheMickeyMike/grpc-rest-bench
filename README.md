# grpc-rest-bench
Benchmarking basic REST API based on HTTP and gRRPC.  
This project is part of my Master thesis.  


## Test 1 REST API

### HTTP/1.1 vs HTTTP/2
>Limitation of HTTP/1.1

![Time/Gorutines count](/docs/img/res-1.png)

>Number of retries during benchmark

![Number of retries/Gorutines count](/docs/img/res1-r.png)

### HTTTP/2 MAX
>HTTP/2 reaches its performance limit

![Tux, the Linux mascot](/docs/img/res1-m.png)


## Test 2 gRPC vs REST (HTTP/2)
>gRPC shows a 50-95% reduction in processing time.

![Time/Gorutines count](/docs/img/res-2.png)