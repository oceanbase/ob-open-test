#/bin/bash
go build -o obopentest cmd/obopentest.go
docker build -t obopentest:test .
kubectl run  -it -name obopentest --image obopentest:test -n obopentest