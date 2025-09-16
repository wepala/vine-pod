module github.com/wepala/vine-pod

go 1.23.0

toolchain go1.23.3

replace github.com/wepala/vine-os/core/pericarp => ../../core/pericarp

require (
	github.com/go-kratos/kratos/contrib/log/zap/v2 v2.0.0-20250912104010-25b6c0fb9f38
	github.com/go-kratos/kratos/v2 v2.9.1
	go.uber.org/zap v1.26.0
)

require (
	github.com/go-kratos/aegis v0.2.0 // indirect
	github.com/go-playground/form/v4 v4.2.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	go.uber.org/dig v1.17.0 // indirect
	go.uber.org/fx v1.20.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/grpc v1.61.2 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
