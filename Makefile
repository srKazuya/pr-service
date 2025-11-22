.PHONY: gen-openapi

OPENAPI_FILE=internal/infrastructure/http/openapi/openapi.yml
OPENAPI_OUT=internal/infrastructure/http/openapi

gen-openapi:
	oapi-codegen -generate types -o $(OPENAPI_OUT)/types.gen.go $(OPENAPI_FILE) 
	oapi-codegen -generate chi-server -o $(OPENAPI_OUT)/server.gen.go $(OPENAPI_FILE)
