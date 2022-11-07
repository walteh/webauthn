setup:
	mkdir -p ~/.aws-lambda-rie && \
	curl -Lo ~/.aws-lambda-rie/aws-lambda-rie \
		https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie-arm64 && \
	chmod +x ~/.aws-lambda-rie/aws-lambda-rie

run:
	docker compose build && docker compose up

invoke:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"entity":"Nuggs","path":"","host":"127.0.0.1:8181","headers":["X-ECHO-BODY:PAYLOAD","X-ECHO-CODE:200"],"scheme":"http"}'

invoke2:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"entity":"0x82403EF00000000000017750C1107DB042C0011"}'

