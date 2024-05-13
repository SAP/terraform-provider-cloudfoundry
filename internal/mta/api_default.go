package mta

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/antihax/optional"
)

const (
	POST          string = "POST"
	GET           string = "GET"
	MTA_NOT_FOUND string = "no mta found with given id"
)

type DefaultApiService service

type DefaultApiGetMtaOperationOpts struct {
	Embed optional.String
}

type DefaultApiStartMtaOperationOpts struct {
	Body optional.Interface
}

type DefaultApiUploadMtaFileOpts struct {
	Namespace optional.String
}

type DefaultApiGetMtaOperationsOpts struct {
	MtaId optional.String
	Last  optional.Int32
	State optional.Interface
}

type Request struct {
	path        string
	method      string
	postBody    any
	headers     map[string]string
	queryParams url.Values
	formParams  url.Values
	fileName    string
	fileBytes   []byte
}

func newRequestInfo() Request {
	return Request{
		headers:     make(map[string]string),
		queryParams: url.Values{},
		formParams:  url.Values{},
	}
}

/*
Executes a particular action over Multi-Target Application operation.
*/
func (a *DefaultApiService) ExecuteOperationAction(ctx context.Context, spaceGuid string, operationId string, actionId string) (string, *http.Response, error) {
	var (
		s       string
		request Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v1/spaces/" + spaceGuid + "/operations/" + operationId
	request.queryParams.Add("actionId", actionId)
	operationId, httpResponse, err := a.client.post(ctx, request, &s)
	return operationId, httpResponse, err
}

/*
Retrieves a csrf-token header.
*/
func (a *DefaultApiService) GetCsrfToken(ctx context.Context) (*http.Response, error) {
	var (
		s       string
		request Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v1/csrf-token"
	httpResponse, err := a.client.get(ctx, request, &s)
	return httpResponse, err
}

/*
Retrieves Multi-Target Application in a space.
*/
func (a *DefaultApiService) GetMta(ctx context.Context, spaceGuid string, mtaId string, namespace string) (Mta, *http.Response, error) {
	var (
		mtas    []Mta
		request Request = newRequestInfo()
	)
	//
	request.path = a.client.cfg.BasePath + "/api/v2/spaces/" + spaceGuid + "/mtas"
	request.queryParams.Add("name", mtaId)
	request.queryParams.Add("namespace", namespace)
	httpResponse, err := a.client.get(ctx, request, &mtas)
	if len(mtas) == 1 {
		return mtas[0], httpResponse, err
	}
	return Mta{}, httpResponse, fmt.Errorf(MTA_NOT_FOUND)
}

/*
Retrieves Multi-Target Application operation.
*/
func (a *DefaultApiService) GetMtaOperation(ctx context.Context, spaceGuid string, operationId string, localVarOptionals *DefaultApiGetMtaOperationOpts) (Operation, *http.Response, error) {
	var (
		operation Operation
		request   Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v1/spaces/" + spaceGuid + "/operations/" + operationId
	if localVarOptionals != nil && localVarOptionals.Embed.IsSet() {
		request.queryParams.Add("embed", parameterToString(localVarOptionals.Embed.Value(), ""))
	}
	httpResponse, err := a.client.get(ctx, request, &operation)
	return operation, httpResponse, err
}

/*
Retrieves Multi-Target Application operations.
*/
func (a *DefaultApiService) GetMtaOperations(ctx context.Context, spaceGuid string, localVarOptionals *DefaultApiGetMtaOperationsOpts) ([]Operation, *http.Response, error) {
	var (
		operations []Operation
		request    Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v1/spaces/" + spaceGuid + "/operations"
	if localVarOptionals != nil {
		if localVarOptionals.MtaId.IsSet() {
			request.queryParams.Add("mtaId", parameterToString(localVarOptionals.MtaId.Value(), ""))
		}
		if localVarOptionals.Last.IsSet() {
			request.queryParams.Add("last", parameterToString(localVarOptionals.Last.Value(), ""))
		}
		if localVarOptionals.State.IsSet() {
			request.queryParams.Add("state", parameterToString(localVarOptionals.State.Value(), "csv"))
		}
	}
	httpResponse, err := a.client.get(ctx, request, &operations)
	return operations, httpResponse, err
}

/*
Retrieves all Multi-Target Applications in a space.
*/
func (a *DefaultApiService) GetMtas(ctx context.Context, spaceGuid string, namespace *string, name string) ([]Mta, *http.Response, error) {
	var (
		mtas    []Mta
		request Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v2/spaces/" + spaceGuid + "/mtas"
	if namespace != nil {
		request.queryParams.Add("namespace", *namespace)
	}
	if name != "" {
		request.queryParams.Add("name", name)
	}
	httpResponse, err := a.client.get(ctx, request, &mtas)
	return mtas, httpResponse, err
}

/*
Starts execution of a Multi-Target Application operation.
*/
func (a *DefaultApiService) StartMtaOperation(ctx context.Context, spaceGuid string, localVarOptionals *DefaultApiStartMtaOperationOpts) (string, Operation, *http.Response, error) {
	var (
		operation Operation
		request   Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v1/spaces/" + spaceGuid + "/operations"
	if localVarOptionals != nil && localVarOptionals.Body.IsSet() {
		localVarOptionalBody, localVarOptionalBodyok := localVarOptionals.Body.Value().(Operation)
		if !localVarOptionalBodyok {
			return "", operation, nil, fmt.Errorf("body should be an Operation Request")
		}
		request.postBody = &localVarOptionalBody
	}
	operationId, httpResponse, err := a.client.post(ctx, request, &operation)
	return operationId, operation, httpResponse, err
}

/*
Uploads an Multi Target Application archive or an Extension Descriptor.
*/
func (a *DefaultApiService) UploadMtaFile(ctx context.Context, spaceGuid string, localVarOptionals *DefaultApiUploadMtaFileOpts, filePath string) (FileMetadata, *http.Response, error) {
	var (
		file    FileMetadata
		err     error
		request Request = newRequestInfo()
	)
	if filePath != "" {
		request.fileBytes, err = os.ReadFile(filePath)
		if err != nil {
			return file, nil, err
		}
		request.fileName = filepath.Base(filePath)
	} else {
		return file, nil, errors.New("filePath required for uploading")
	}
	request.path = a.client.cfg.BasePath + "/api/v1/spaces/" + spaceGuid + "/files"
	if localVarOptionals != nil && localVarOptionals.Namespace.IsSet() {
		request.queryParams.Add("namespace", parameterToString(localVarOptionals.Namespace.Value(), ""))
	}
	request.headers["Content-Type"] = "multipart/form-data"
	_, httpResponse, err := a.client.post(ctx, request, &file)
	return file, httpResponse, err
}

/*
Uploads an Multi Target Application archive from a remote URL.
*/
func (a *DefaultApiService) AsyncUploadFileFromURL(ctx context.Context, spaceGuid string, localVarOptionals *DefaultApiUploadMtaFileOpts, fileURL string) (string, *http.Response, error) {
	var (
		s       string
		request Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v1/spaces/" + spaceGuid + "/files/async"
	if localVarOptionals != nil && localVarOptionals.Namespace.IsSet() {
		request.queryParams.Add("namespace", parameterToString(localVarOptionals.Namespace.Value(), ""))
	}
	if fileURL != "" {
		request.postBody = &FileUrl{
			FileUrl: base64.URLEncoding.EncodeToString([]byte(fileURL)),
		}
	} else {
		return s, nil, errors.New("fileURL required for uploading")
	}
	operationId, httpResponse, err := a.client.post(ctx, request, &s)
	return operationId, httpResponse, err
}

/*
Gets the status of the upload job from the remote URL.
*/
func (a *DefaultApiService) GetAsyncUploadJob(ctx context.Context, spaceGuid string, jobId string, xInstance string, localVarOptionals *DefaultApiUploadMtaFileOpts) (UploadStatus, *http.Response, error) {
	var (
		uploadJob UploadStatus
		request   Request = newRequestInfo()
	)
	request.path = a.client.cfg.BasePath + "/api/v1/spaces/" + spaceGuid + "/files/jobs/" + jobId
	if localVarOptionals != nil && localVarOptionals.Namespace.IsSet() {
		request.queryParams.Add("namespace", parameterToString(localVarOptionals.Namespace.Value(), ""))
	}
	request.headers["x-cf-app-instance"] = xInstance
	httpResponse, err := a.client.get(ctx, request, &uploadJob)
	return uploadJob, httpResponse, err
}
