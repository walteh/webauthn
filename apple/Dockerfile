# syntax=docker/dockerfile:1

FROM public.ecr.aws/lambda/provided:al2 as builder

RUN yum install -y golang
RUN go env -w GOPROXY=direct
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /bin/handler main.go

FROM public.ecr.aws/lambda/provided:al2

COPY --from=builder /bin/handler ${LAMBDA_TASK_ROOT}

ENTRYPOINT ["/var/task/handler"]
