## ESGBook

# Stack/Tools used:
1) Github Actions  - CI
2) Go App with ```/ping``` and ```/metrics``` routes
3) docker       - container runtime used, other alternaive is containerd
4) kubernetes   - orchestrator

## Testing app and metrics in local env:
- ```SERVICE__PORT=8080 TARGET=http://localhost:8080/ping  go run .```  #start and test app
- ```CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o pingpong```   #build executable
- ```docker build -t pingpong . -f Dockerfile.multistage```             #build image
- ```docker run -p 8080:8080 -p 9080:9080 --env-file .env  pingpong```  #run image
- Metrics:
- Added ```ping_requests_received_total``` and ```ping_failures_total``` metrics in metrics.go and modified startAppServer function in main.go accordingly
- Go to ```http://localhost:9080/-/metrics``` to view the metrics
- Simulating a failure request to observe ```ping_failures_total``` metrics using ```curl -X POST http://localhost:3030/ping```

## Deployment in local kubernetes
- Built the image and pushed to dockerhub ```https://hub.docker.com/repository/docker/kinnerakallam/esgbook/tags``` using Github actions. CI is ```https://github.com/kinnerakallam/esgbook-platform-engineer/actions```

- Deploying manifests using built image above. ```kubectl apply -f manifest.yaml```

### HTTPS Implementation
- Istio to be installed on the cluster
- TLS termination can be done at different levels like application, CDN, Loadbalancer, Ingress/Istio etc
- Application: 
    - Implemented it and can demo it during the interview. Not advised in general other than some exceptional cases. Why to load application with this when we have other secure options
- Ingress: 
    - Implemented TLS termination with nginx Ingress and can demo it during interview
- Istio: 
    - Require three major components for the implementation: Secret, Gateway and Virtual Service
    - Gateway and Virtual Service is in the manifest.
    - Secret
        - Created a self signed cert and key using ```openssl req -x509 -nodes -days 365 -newkey rsa:2048 
        -keyout server.key -out server.crt -subj \
        "/CN=test.localdev.me/O=test.localdev.me" ```
        - Ideally we store the cert and key in secrets manager but for this demo, I stored the data locally and created a secret in ```istio-system``` namespace 
        ```kubectl create secret tls httpbin-credential --cert=./server.crt --key= server.key -n istio-system```
#### Testing
- Portforwarding ingressgateway to localhost. Here forwarding both http and https ports
```sudo kubectl port-forward svc/istio-ingressgateway -n istio-system 8080:80 443:443``` 
- Hit ```https://test.localdev.me:443/ping``` in your browser or ```curl -v -k https://test.localdev.me:443/ping``` in your terminal - This is for pingpong-a 
- Hit ```https://test-b.localdev.me:443/ping``` in your browser or ```curl -v -k https://test-b.localdev.me:443/ping``` in your terminal - This is for pingpong-b
- You get a ```pong``` response for https request which proves that https is enabled

### mTLS
- Implemented it with Istio
- ```pingpong-a``` and ```pingpong-b``` namespaces have the istio-injection enabled on them
- PeerAuthentication with mTLS mode strict is created for both the namespaces
#### Testing
- For testing, created a namespace ```pingpong-c``` and a ```curlpod``` for testing. Istio-injection is not enabled on this namespace
    - When sent a request from pingpong-a to pingpong-b or vice-versa, the communication works in a secured way(we get a success response)
        - ```kubectl exec -it <pingpong-a pod> -n pingpong-a -- curl http://pingponger-b.pingpong-b.svc.cluster.local:8080/ping```
        - ```kubectl exec -it <pingpong-b pod> -n pingpong-b -- curl http://pingponger-a.pingpong-a.svc.cluster.local:8080/ping```
    - When sent a request from pingpong-c pod to other two namespaces pods, the communication fails which proves that mTLS is enabled
       - ```kubectl exec -it curlpod -n pingpong-c -- curl http://pingponger-b.pingpong-b.svc.cluster.local:8080/ping```
       - ```kubectl exec -it curlpod -n pingpong-c -- curl http://pingponger-a.pingpong-a.svc.cluster.local:8080/ping```

## Few features to make the cluster Production ready

### Security
- Pod Security Admission(Pod Security Policies newer version as its depricated)
    -  Implemented this in the current setup
    -  Created a namespace ```pingpong-d``` with a label ```pod-security.kubernetes.io/enforce: restricted``` which will not allow anything to run with privileged access
    -  Try creating ```privileged-pod``` with security context as privileged in ```pingpong-d``` and it fails to create.
- RBAC can be implemented
- Service Accounts
- Network Policies
    - Implemented this
    - Created a network policy which deny all the traffic(ingress and egress) on the namespace ```pingpong-e```
    - Run ```kubectl run -n pingpong-d egress-test --rm -it --image=busybox -- /bin/sh``` and after getting in curl any url and it fails

### Reliability
- Health Checks (Liveness Probe & Readiness Probe)
    - Implemented for deployment ```pingpong-a```
- Autoscaling (Pod disruption budgets, Horizontal Pod autoscaling)
    - Implemented for deployment ```pingpong-a```
- Rolling updates

### Observability
- Dashboards for the ease of monitoring

### Deployments
- Helm charts for kubernetes objects
- GitOps(ArgoCD)
- CI/CD pipelines

# Improvements:
- Code is not loading values from ```.env``` file. From my research, maybe koanf library doesn't reads environment variables from the system or from file, maybe we need to use different library ```godotenv```
- Two deployments for the same application is not required, we can have multiple replicas in a single deployment

