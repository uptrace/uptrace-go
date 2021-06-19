module github.com/uptrace/uptrace-go/example/go-pg

go 1.15

replace github.com/uptrace/uptrace-go => ../..

require (
	github.com/go-pg/pg/extra/pgotel/v10 v10.10.0
	github.com/go-pg/pg/v10 v10.10.0
	github.com/uptrace/uptrace-go v0.20.0
	go.opentelemetry.io/otel v1.0.0-RC1
	golang.org/x/crypto v0.0.0-20210616213533-5ff15b29337e // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
)
