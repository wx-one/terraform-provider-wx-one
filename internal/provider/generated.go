// Code generated by github.com/Khan/genqlient, DO NOT EDIT.

package provider

import (
	"context"

	"github.com/Khan/genqlient/graphql"
)

// __createKeyInput is used internally by genqlient
type __createKeyInput struct {
	Name        string `json:"name"`
	PublicKey   string `json:"publicKey"`
	ProjectId   string `json:"projectId,omitempty"`
	ProjectWide bool   `json:"projectWide"`
}

// GetName returns __createKeyInput.Name, and is useful for accessing the field via an interface.
func (v *__createKeyInput) GetName() string { return v.Name }

// GetPublicKey returns __createKeyInput.PublicKey, and is useful for accessing the field via an interface.
func (v *__createKeyInput) GetPublicKey() string { return v.PublicKey }

// GetProjectId returns __createKeyInput.ProjectId, and is useful for accessing the field via an interface.
func (v *__createKeyInput) GetProjectId() string { return v.ProjectId }

// GetProjectWide returns __createKeyInput.ProjectWide, and is useful for accessing the field via an interface.
func (v *__createKeyInput) GetProjectWide() bool { return v.ProjectWide }

// __deleteKeyInput is used internally by genqlient
type __deleteKeyInput struct {
	Id        string `json:"id"`
	ProjectId string `json:"projectId,omitempty"`
}

// GetId returns __deleteKeyInput.Id, and is useful for accessing the field via an interface.
func (v *__deleteKeyInput) GetId() string { return v.Id }

// GetProjectId returns __deleteKeyInput.ProjectId, and is useful for accessing the field via an interface.
func (v *__deleteKeyInput) GetProjectId() string { return v.ProjectId }

// __getKeyInput is used internally by genqlient
type __getKeyInput struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	ProjectId   string `json:"projectId,omitempty"`
	ProjectWide *bool  `json:"projectWide,omitempty"`
}

// GetId returns __getKeyInput.Id, and is useful for accessing the field via an interface.
func (v *__getKeyInput) GetId() string { return v.Id }

// GetName returns __getKeyInput.Name, and is useful for accessing the field via an interface.
func (v *__getKeyInput) GetName() string { return v.Name }

// GetProjectId returns __getKeyInput.ProjectId, and is useful for accessing the field via an interface.
func (v *__getKeyInput) GetProjectId() string { return v.ProjectId }

// GetProjectWide returns __getKeyInput.ProjectWide, and is useful for accessing the field via an interface.
func (v *__getKeyInput) GetProjectWide() *bool { return v.ProjectWide }

// __updateKeyInput is used internally by genqlient
type __updateKeyInput struct {
	Id        string `json:"id"`
	ProjectId string `json:"projectId,omitempty"`
	Name      string `json:"name"`
}

// GetId returns __updateKeyInput.Id, and is useful for accessing the field via an interface.
func (v *__updateKeyInput) GetId() string { return v.Id }

// GetProjectId returns __updateKeyInput.ProjectId, and is useful for accessing the field via an interface.
func (v *__updateKeyInput) GetProjectId() string { return v.ProjectId }

// GetName returns __updateKeyInput.Name, and is useful for accessing the field via an interface.
func (v *__updateKeyInput) GetName() string { return v.Name }

// createKeyCreateKeyW1KeyResponse includes the requested fields of the GraphQL type W1KeyResponse.
// The GraphQL type's documentation follows.
//
// Key Response
type createKeyCreateKeyW1KeyResponse struct {
	// Return Code
	Code int `json:"code"`
	// Error Message
	Err string `json:"err"`
	// Success Message
	Msg createKeyCreateKeyW1KeyResponseMsgW1Key `json:"msg"`
}

// GetCode returns createKeyCreateKeyW1KeyResponse.Code, and is useful for accessing the field via an interface.
func (v *createKeyCreateKeyW1KeyResponse) GetCode() int { return v.Code }

// GetErr returns createKeyCreateKeyW1KeyResponse.Err, and is useful for accessing the field via an interface.
func (v *createKeyCreateKeyW1KeyResponse) GetErr() string { return v.Err }

// GetMsg returns createKeyCreateKeyW1KeyResponse.Msg, and is useful for accessing the field via an interface.
func (v *createKeyCreateKeyW1KeyResponse) GetMsg() createKeyCreateKeyW1KeyResponseMsgW1Key {
	return v.Msg
}

// createKeyCreateKeyW1KeyResponseMsgW1Key includes the requested fields of the GraphQL type W1Key.
// The GraphQL type's documentation follows.
//
// SSH Key
type createKeyCreateKeyW1KeyResponseMsgW1Key struct {
	// ID
	Id string `json:"id"`
}

// GetId returns createKeyCreateKeyW1KeyResponseMsgW1Key.Id, and is useful for accessing the field via an interface.
func (v *createKeyCreateKeyW1KeyResponseMsgW1Key) GetId() string { return v.Id }

// createKeyResponse is returned by createKey on success.
type createKeyResponse struct {
	// Create Key
	CreateKey createKeyCreateKeyW1KeyResponse `json:"createKey"`
}

// GetCreateKey returns createKeyResponse.CreateKey, and is useful for accessing the field via an interface.
func (v *createKeyResponse) GetCreateKey() createKeyCreateKeyW1KeyResponse { return v.CreateKey }

// deleteKeyDeleteKeyResponse includes the requested fields of the GraphQL type Response.
// The GraphQL type's documentation follows.
//
// Response
type deleteKeyDeleteKeyResponse struct {
	// Return Code
	Code int `json:"code"`
	// Error Message
	Err string `json:"err"`
	// Success Message
	Msg string `json:"msg"`
}

// GetCode returns deleteKeyDeleteKeyResponse.Code, and is useful for accessing the field via an interface.
func (v *deleteKeyDeleteKeyResponse) GetCode() int { return v.Code }

// GetErr returns deleteKeyDeleteKeyResponse.Err, and is useful for accessing the field via an interface.
func (v *deleteKeyDeleteKeyResponse) GetErr() string { return v.Err }

// GetMsg returns deleteKeyDeleteKeyResponse.Msg, and is useful for accessing the field via an interface.
func (v *deleteKeyDeleteKeyResponse) GetMsg() string { return v.Msg }

// deleteKeyResponse is returned by deleteKey on success.
type deleteKeyResponse struct {
	// Delete Key
	DeleteKey deleteKeyDeleteKeyResponse `json:"deleteKey"`
}

// GetDeleteKey returns deleteKeyResponse.DeleteKey, and is useful for accessing the field via an interface.
func (v *deleteKeyResponse) GetDeleteKey() deleteKeyDeleteKeyResponse { return v.DeleteKey }

// getDefaultProjectGetDefaultProjectProjectResponse includes the requested fields of the GraphQL type ProjectResponse.
// The GraphQL type's documentation follows.
//
// ProjectResponse
type getDefaultProjectGetDefaultProjectProjectResponse struct {
	// Return Code
	Code int `json:"code"`
	// Error Message
	Err string `json:"err"`
	// Success Message
	Msg getDefaultProjectGetDefaultProjectProjectResponseMsgProject `json:"msg"`
}

// GetCode returns getDefaultProjectGetDefaultProjectProjectResponse.Code, and is useful for accessing the field via an interface.
func (v *getDefaultProjectGetDefaultProjectProjectResponse) GetCode() int { return v.Code }

// GetErr returns getDefaultProjectGetDefaultProjectProjectResponse.Err, and is useful for accessing the field via an interface.
func (v *getDefaultProjectGetDefaultProjectProjectResponse) GetErr() string { return v.Err }

// GetMsg returns getDefaultProjectGetDefaultProjectProjectResponse.Msg, and is useful for accessing the field via an interface.
func (v *getDefaultProjectGetDefaultProjectProjectResponse) GetMsg() getDefaultProjectGetDefaultProjectProjectResponseMsgProject {
	return v.Msg
}

// getDefaultProjectGetDefaultProjectProjectResponseMsgProject includes the requested fields of the GraphQL type Project.
// The GraphQL type's documentation follows.
//
// Project
type getDefaultProjectGetDefaultProjectProjectResponseMsgProject struct {
	// ID
	Id string `json:"id"`
	// Name
	Name string `json:"name"`
}

// GetId returns getDefaultProjectGetDefaultProjectProjectResponseMsgProject.Id, and is useful for accessing the field via an interface.
func (v *getDefaultProjectGetDefaultProjectProjectResponseMsgProject) GetId() string { return v.Id }

// GetName returns getDefaultProjectGetDefaultProjectProjectResponseMsgProject.Name, and is useful for accessing the field via an interface.
func (v *getDefaultProjectGetDefaultProjectProjectResponseMsgProject) GetName() string { return v.Name }

// getDefaultProjectResponse is returned by getDefaultProject on success.
type getDefaultProjectResponse struct {
	// Get Default Project
	GetDefaultProject getDefaultProjectGetDefaultProjectProjectResponse `json:"getDefaultProject"`
}

// GetGetDefaultProject returns getDefaultProjectResponse.GetDefaultProject, and is useful for accessing the field via an interface.
func (v *getDefaultProjectResponse) GetGetDefaultProject() getDefaultProjectGetDefaultProjectProjectResponse {
	return v.GetDefaultProject
}

// getKeyGetKeyW1KeyResponse includes the requested fields of the GraphQL type W1KeyResponse.
// The GraphQL type's documentation follows.
//
// Key Response
type getKeyGetKeyW1KeyResponse struct {
	// Return Code
	Code int `json:"code"`
	// Error Message
	Err string `json:"err"`
	// Success Message
	Msg getKeyGetKeyW1KeyResponseMsgW1Key `json:"msg"`
}

// GetCode returns getKeyGetKeyW1KeyResponse.Code, and is useful for accessing the field via an interface.
func (v *getKeyGetKeyW1KeyResponse) GetCode() int { return v.Code }

// GetErr returns getKeyGetKeyW1KeyResponse.Err, and is useful for accessing the field via an interface.
func (v *getKeyGetKeyW1KeyResponse) GetErr() string { return v.Err }

// GetMsg returns getKeyGetKeyW1KeyResponse.Msg, and is useful for accessing the field via an interface.
func (v *getKeyGetKeyW1KeyResponse) GetMsg() getKeyGetKeyW1KeyResponseMsgW1Key { return v.Msg }

// getKeyGetKeyW1KeyResponseMsgW1Key includes the requested fields of the GraphQL type W1Key.
// The GraphQL type's documentation follows.
//
// SSH Key
type getKeyGetKeyW1KeyResponseMsgW1Key struct {
	// ID
	Id string `json:"id"`
	// Name
	Name string `json:"name"`
	// Public Key
	PublicKey string `json:"publicKey"`
	// Can this key be used for the complete project or only the user
	ProjectWide bool `json:"projectWide"`
}

// GetId returns getKeyGetKeyW1KeyResponseMsgW1Key.Id, and is useful for accessing the field via an interface.
func (v *getKeyGetKeyW1KeyResponseMsgW1Key) GetId() string { return v.Id }

// GetName returns getKeyGetKeyW1KeyResponseMsgW1Key.Name, and is useful for accessing the field via an interface.
func (v *getKeyGetKeyW1KeyResponseMsgW1Key) GetName() string { return v.Name }

// GetPublicKey returns getKeyGetKeyW1KeyResponseMsgW1Key.PublicKey, and is useful for accessing the field via an interface.
func (v *getKeyGetKeyW1KeyResponseMsgW1Key) GetPublicKey() string { return v.PublicKey }

// GetProjectWide returns getKeyGetKeyW1KeyResponseMsgW1Key.ProjectWide, and is useful for accessing the field via an interface.
func (v *getKeyGetKeyW1KeyResponseMsgW1Key) GetProjectWide() bool { return v.ProjectWide }

// getKeyResponse is returned by getKey on success.
type getKeyResponse struct {
	// Get Key
	GetKey getKeyGetKeyW1KeyResponse `json:"getKey"`
}

// GetGetKey returns getKeyResponse.GetKey, and is useful for accessing the field via an interface.
func (v *getKeyResponse) GetGetKey() getKeyGetKeyW1KeyResponse { return v.GetKey }

// meMeUser includes the requested fields of the GraphQL type User.
// The GraphQL type's documentation follows.
//
// User
type meMeUser struct {
	// ID
	Id string `json:"id"`
	// Name
	Username string `json:"username"`
	// Role
	Role int `json:"role"`
}

// GetId returns meMeUser.Id, and is useful for accessing the field via an interface.
func (v *meMeUser) GetId() string { return v.Id }

// GetUsername returns meMeUser.Username, and is useful for accessing the field via an interface.
func (v *meMeUser) GetUsername() string { return v.Username }

// GetRole returns meMeUser.Role, and is useful for accessing the field via an interface.
func (v *meMeUser) GetRole() int { return v.Role }

// meResponse is returned by me on success.
type meResponse struct {
	// Get information about me
	Me meMeUser `json:"me"`
}

// GetMe returns meResponse.Me, and is useful for accessing the field via an interface.
func (v *meResponse) GetMe() meMeUser { return v.Me }

// updateKeyResponse is returned by updateKey on success.
type updateKeyResponse struct {
	// Update Key
	UpdateKey updateKeyUpdateKeyW1KeyResponse `json:"updateKey"`
}

// GetUpdateKey returns updateKeyResponse.UpdateKey, and is useful for accessing the field via an interface.
func (v *updateKeyResponse) GetUpdateKey() updateKeyUpdateKeyW1KeyResponse { return v.UpdateKey }

// updateKeyUpdateKeyW1KeyResponse includes the requested fields of the GraphQL type W1KeyResponse.
// The GraphQL type's documentation follows.
//
// Key Response
type updateKeyUpdateKeyW1KeyResponse struct {
	// Return Code
	Code int `json:"code"`
	// Error Message
	Err string `json:"err"`
	// Success Message
	Msg updateKeyUpdateKeyW1KeyResponseMsgW1Key `json:"msg"`
}

// GetCode returns updateKeyUpdateKeyW1KeyResponse.Code, and is useful for accessing the field via an interface.
func (v *updateKeyUpdateKeyW1KeyResponse) GetCode() int { return v.Code }

// GetErr returns updateKeyUpdateKeyW1KeyResponse.Err, and is useful for accessing the field via an interface.
func (v *updateKeyUpdateKeyW1KeyResponse) GetErr() string { return v.Err }

// GetMsg returns updateKeyUpdateKeyW1KeyResponse.Msg, and is useful for accessing the field via an interface.
func (v *updateKeyUpdateKeyW1KeyResponse) GetMsg() updateKeyUpdateKeyW1KeyResponseMsgW1Key {
	return v.Msg
}

// updateKeyUpdateKeyW1KeyResponseMsgW1Key includes the requested fields of the GraphQL type W1Key.
// The GraphQL type's documentation follows.
//
// SSH Key
type updateKeyUpdateKeyW1KeyResponseMsgW1Key struct {
	// ID
	Id string `json:"id"`
}

// GetId returns updateKeyUpdateKeyW1KeyResponseMsgW1Key.Id, and is useful for accessing the field via an interface.
func (v *updateKeyUpdateKeyW1KeyResponseMsgW1Key) GetId() string { return v.Id }

// The query or mutation executed by createKey.
const createKey_Operation = `
mutation createKey ($name: String!, $publicKey: String!, $projectId: UUID, $projectWide: Boolean) {
	createKey(name: $name, publicKey: $publicKey, projectId: $projectId, projectWide: $projectWide) {
		code
		err
		msg {
			id
		}
	}
}
`

func createKey(
	ctx_ context.Context,
	client_ graphql.Client,
	name string,
	publicKey string,
	projectId string,
	projectWide bool,
) (*createKeyResponse, error) {
	req_ := &graphql.Request{
		OpName: "createKey",
		Query:  createKey_Operation,
		Variables: &__createKeyInput{
			Name:        name,
			PublicKey:   publicKey,
			ProjectId:   projectId,
			ProjectWide: projectWide,
		},
	}
	var err_ error

	var data_ createKeyResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by deleteKey.
const deleteKey_Operation = `
mutation deleteKey ($id: UUID!, $projectId: UUID) {
	deleteKey(id: $id, projectId: $projectId) {
		code
		err
		msg
	}
}
`

func deleteKey(
	ctx_ context.Context,
	client_ graphql.Client,
	id string,
	projectId string,
) (*deleteKeyResponse, error) {
	req_ := &graphql.Request{
		OpName: "deleteKey",
		Query:  deleteKey_Operation,
		Variables: &__deleteKeyInput{
			Id:        id,
			ProjectId: projectId,
		},
	}
	var err_ error

	var data_ deleteKeyResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by getDefaultProject.
const getDefaultProject_Operation = `
query getDefaultProject {
	getDefaultProject {
		code
		err
		msg {
			id
			name
		}
	}
}
`

func getDefaultProject(
	ctx_ context.Context,
	client_ graphql.Client,
) (*getDefaultProjectResponse, error) {
	req_ := &graphql.Request{
		OpName: "getDefaultProject",
		Query:  getDefaultProject_Operation,
	}
	var err_ error

	var data_ getDefaultProjectResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by getKey.
const getKey_Operation = `
query getKey ($id: UUID, $name: String, $projectId: UUID, $projectWide: Boolean) {
	getKey(id: $id, name: $name, projectId: $projectId, projectWide: $projectWide) {
		code
		err
		msg {
			id
			name
			publicKey
			projectWide
		}
	}
}
`

func getKey(
	ctx_ context.Context,
	client_ graphql.Client,
	id string,
	name string,
	projectId string,
	projectWide *bool,
) (*getKeyResponse, error) {
	req_ := &graphql.Request{
		OpName: "getKey",
		Query:  getKey_Operation,
		Variables: &__getKeyInput{
			Id:          id,
			Name:        name,
			ProjectId:   projectId,
			ProjectWide: projectWide,
		},
	}
	var err_ error

	var data_ getKeyResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by me.
const me_Operation = `
query me {
	me {
		id
		username
		role
	}
}
`

func me(
	ctx_ context.Context,
	client_ graphql.Client,
) (*meResponse, error) {
	req_ := &graphql.Request{
		OpName: "me",
		Query:  me_Operation,
	}
	var err_ error

	var data_ meResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}

// The query or mutation executed by updateKey.
const updateKey_Operation = `
mutation updateKey ($id: UUID!, $projectId: UUID, $name: String!) {
	updateKey(id: $id, name: $name, projectId: $projectId) {
		code
		err
		msg {
			id
		}
	}
}
`

func updateKey(
	ctx_ context.Context,
	client_ graphql.Client,
	id string,
	projectId string,
	name string,
) (*updateKeyResponse, error) {
	req_ := &graphql.Request{
		OpName: "updateKey",
		Query:  updateKey_Operation,
		Variables: &__updateKeyInput{
			Id:        id,
			ProjectId: projectId,
			Name:      name,
		},
	}
	var err_ error

	var data_ updateKeyResponse
	resp_ := &graphql.Response{Data: &data_}

	err_ = client_.MakeRequest(
		ctx_,
		req_,
		resp_,
	)

	return &data_, err_
}
