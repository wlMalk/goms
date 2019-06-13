package generate_service

import (
	"os"
	strs "strings"

	"github.com/wlMalk/goms/generator/files"
	"github.com/wlMalk/goms/generator/strings"
	"github.com/wlMalk/goms/parser/types"
)

func GenerateDockerFile(base string, service *types.Service) *files.TextFile {
	serviceNameKebab := strings.ToLower(strings.ToKebabCase(service.Name))
	file := files.NewTextFile(base, "", "Dockerfile", "", false, false)
	path := "/go" + strs.TrimPrefix(file.Base(), os.Getenv("GOPATH"))
	file.Pf("FROM golang:1.12 as builder")
	file.Pf("ADD . %s", path)
	file.Pf("WORKDIR %s/cmd/%s", path, serviceNameKebab)
	file.Pf("RUN go get -d -v ./...")
	file.Pf("RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags \"-static\"' -o %s .", serviceNameKebab)
	file.Pf("FROM scratch")
	file.Pf("COPY --from=builder %s/cmd/%s/%s /%s/", path, serviceNameKebab, serviceNameKebab, serviceNameKebab)
	file.Pf("WORKDIR /%s", serviceNameKebab)
	file.Pf("EXPOSE 8080")
	file.Pf("ENTRYPOINT [\"./%s\"]", serviceNameKebab)
	return file
}
