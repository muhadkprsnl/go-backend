module local/qa-report

go 1.24.3

require (
	github.com/golang-jwt/jwt/v4 v4.5.2
	github.com/gorilla/mux v1.8.1
	go.mongodb.org/mongo-driver v1.17.4
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.26.0
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/text v0.17.0 // indirectn
)

replace (
	local/qa-report/internal/config => ./internal/config
	local/qa-report/internal/routes => ./internal/routes
	local/qa-report/pkg/middleware => ./pkg/middleware
)
