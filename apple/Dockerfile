# syntax=docker/dockerfile:1

FROM public.ecr.aws/lambda/provided:al2 as builder

RUN yum install -y golang
RUN go env -w GOPROXY=direct
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /bin/apple-apigw cmd/apple-apigw/main.go
RUN GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /bin/apple-appsync cmd/apple-appsync/main.go
RUN GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o /bin/challenge cmd/challenge/main.go

FROM public.ecr.aws/lambda/provided:al2 as apple-apigw

COPY --from=builder /bin/apple-apigw ${LAMBDA_TASK_ROOT}

ENTRYPOINT ["/var/task/handler"]

FROM public.ecr.aws/lambda/provided:al2 as apple-appsync

COPY --from=builder /bin/apple-appsync ${LAMBDA_TASK_ROOT}

ENTRYPOINT ["/var/task/handler"]

FROM public.ecr.aws/lambda/provided:al2 as challenge

COPY --from=builder /bin/challenge ${LAMBDA_TASK_ROOT}

ENTRYPOINT ["/var/task/handler"]
