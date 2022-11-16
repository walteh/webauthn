setup:
	mkdir -p ~/.aws-lambda-rie && \
	curl -Lo ~/.aws-lambda-rie/aws-lambda-rie \
		https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/latest/download/aws-lambda-rie-arm64 && \
	chmod +x ~/.aws-lambda-rie/aws-lambda-rie

run:
	docker compose build && docker compose up

invoke:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" --json  "`jq -c . < $(FILE)`"

invoke2:
	curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -d '{"entity":"0x82403EF00000000000017750C1107DB042C0011"}'



invoke-register: # load json from fil
	 curl -XPOST "http://localhost:9000/2015-03-31/functions/function/invocations" -H 'Content-Type: application/json' --json '{"entity":"0x82403EF00000000000017750C1107DB042C0011"}' --cli-binary-format raw-in-base64-out

