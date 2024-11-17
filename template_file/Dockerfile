#FROM golang:1.22 AS go-builder
#
## 为我们的镜像设置必要的环境变量
#ENV GO111MODULE=on \
#GOPROXY=https://goproxy.cn,direct \
#CGO_ENABLED=0 \
#GOOS=linux \
#GOARCH=amd64
#
#WORKDIR /build
#
## pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY go.mod go.sum /build/
#RUN go mod download && go mod verify
#
#COPY . .
#RUN go build -o srbbs-app ./src/
#
#
#
#
#FROM centos:7
#ENV container docker
#RUN (cd /lib/systemd/system/sysinit.target.wants/; for i in *; do [ $i == \
#systemd-tmpfiles-setup.service ] || rm -f $i; done); \
#rm -f /lib/systemd/system/multi-user.target.wants/*;\
#rm -f /etc/systemd/system/*.wants/*;\
#rm -f /lib/systemd/system/local-fs.target.wants/*; \
#rm -f /lib/systemd/system/sockets.target.wants/*udev*; \
#rm -f /lib/systemd/system/sockets.target.wants/*initctl*; \
#rm -f /lib/systemd/system/basic.target.wants/*;\
#rm -f /lib/systemd/system/anaconda.target.wants/*;
#VOLUME [ "/sys/fs/cgroup" ]
##CMD ["/usr/sbin/init"]
#
#COPY --from=go-builder /build/ /build/
#
#EXPOSE 8081
#
#WORKDIR /build
#
#ENTRYPOINT ["./srbbs-app"]