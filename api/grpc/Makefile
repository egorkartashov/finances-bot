send-report:
	mkdir -p ../../internal/grpcapi/send_report
	protoc --go_out=../../internal/grpcapi/send_report \
			--go_opt=paths=source_relative \
			--go-grpc_out=../../internal/grpcapi/send_report \
			--go-grpc_opt=paths=source_relative \
			send_report.proto
