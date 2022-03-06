[root@host1 mode3]# yum install -y golang
[root@host1 mode3]# docker pull golang:1.17
1.17: Pulling from library/golang
0e29546d541c: Pull complete 
9b829c73b52b: Pull complete 
cb5b7ae36172: Pull complete 
6494e4811622: Pull complete 
6e1d20a8313e: Pull complete 
593823f101dc: Pull complete 
1b4aae56cdbe: Pull complete 
Digest: sha256:c72fa9afc50b3303e8044cf28fb358b48032a548e1825819420fd40155a131cb
Status: Downloaded newer image for golang:1.17
docker.io/library/golang:1.17
You have new mail in /var/spool/mail/root
[root@host1 mode3]# docker build . -t httpserver:0.01
Sending build context to Docker daemon  10.24kB
Step 1/13 : FROM golang:1.17 AS build
 ---> 276895edf967
Step 2/13 : WORKDIR /httpserver/
 ---> Running in 4550a62b4d54
Removing intermediate container 4550a62b4d54
 ---> 88ffe117ada4
Step 3/13 : COPY . .
 ---> b37359508bd1
Step 4/13 : ENV CGO_ENABLED=0
 ---> Running in 2a9d6bdfc097
Removing intermediate container 2a9d6bdfc097
 ---> 1906a936f402
Step 5/13 : ENV GO111MODULE=on
 ---> Running in 97f20044127e
Removing intermediate container 97f20044127e
 ---> bb2ae1c6d8fd
Step 6/13 : ENV GOPROXY=https://goproxy.cn.direct
 ---> Running in 535d841dcb00
Removing intermediate container 535d841dcb00
 ---> f12d215dfa6e
Step 7/13 : RUN GOOS=linux go build -installsuffix cgo -o httpserver main.go
 ---> Running in 5f6567df5d4f
Removing intermediate container 5f6567df5d4f
 ---> 9e55ffb322cf
Step 8/13 : FROM busybox
 ---> 829374d342ae
Step 9/13 : COPY --from=build /httpserver/httpserver /httpserver/httpserver
 ---> 85aeaf85233b
Step 10/13 : EXPOSE 8360
 ---> Running in bf0024c0a2fd
Removing intermediate container bf0024c0a2fd
 ---> f720ec7402f2
Step 11/13 : ENV ENV local
 ---> Running in 8382b0695d2d
Removing intermediate container 8382b0695d2d
 ---> 7bfab5796b85
Step 12/13 : WORKDIR /httpserver/
 ---> Running in 71969e10a334
Removing intermediate container 71969e10a334
 ---> 392347e569c9
Step 13/13 : ENTRYPOINT ["./httpserver"]
 ---> Running in a935a3fc428f
Removing intermediate container a935a3fc428f
 ---> eb13631ed45b
Successfully built eb13631ed45b
Successfully tagged httpserver:0.01
[root@host1 mode3]# docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
httpserver          0.01                eb13631ed45b        5 seconds ago       8.27MB
<none>              <none>              9e55ffb322cf        6 seconds ago       1GB
busybox             latest              829374d342ae        36 hours ago        1.24MB
golang              1.17                276895edf967        2 months ago        941MB
[root@host1 mode3]# docker tag httpserver:0.01 registry.cn-hangzhou.aliyuncs.com/cuiqiannan/httpserver:V0.0.1
[root@host1 mode3]# docker push registry.cn-hangzhou.aliyuncs.com/cuiqiannan/httpserver:V0.0.1
The push refers to repository [registry.cn-hangzhou.aliyuncs.com/cuiqiannan/httpserver]
c2ab92ae1a04: Pushed 
252fdf0c3b6a: Pushed 
V0.0.1: digest: sha256:9c85bcf87b0ce0f2dd64056786d44d5fb3855eef2be1d410ee512046c2faf0f6 size: 738
[root@host1 mode3]# docker run -d httpserver:0.01 
40442fb1c2b3717fed65abaf9da131fd234dbd2f4197b075156d81e0844cee9d
[root@host1 mode3]# docker ps
CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS               NAMES
40442fb1c2b3        httpserver:0.01     "./httpserver"      4 seconds ago       Up 3 seconds        8360/tcp            optimistic_chatelet
[root@host1 mode3]# PID=$(docker inspect --format "{{ .State.Pid }}" 40442fb1c2b3)
[root@host1 mode3]# echo $PID
6766
[root@host1 mode3]# nsenter -t $PID -n ip a
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
8: eth0@if9: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever