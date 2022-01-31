package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/srimaln91/etcd-adminer/etcd"
	"github.com/srimaln91/etcd-adminer/http/response"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
)

func (jh *GenericHandler) GetUserList(rw http.ResponseWriter, r *http.Request) {

	user, pass, ok := r.BasicAuth()
	if !ok {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	endpointString := r.Header.Get("X-Endpoints")
	if endpointString == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	endpoints := strings.Split(endpointString, ",")
	if len(endpointString) < 1 {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	client, err := etcd.NewClient(endpoints, etcd.WithAuth(user, pass))
	if err != nil {
		if err == rpctypes.ErrAuthFailed {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	userList, err := client.UserList(r.Context())
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(userList.Users)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)

}

func (jh *GenericHandler) GetUserInfo(rw http.ResponseWriter, r *http.Request) {

	user, pass, ok := r.BasicAuth()
	if !ok {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	endpointString := r.Header.Get("X-Endpoints")
	if endpointString == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	endpoints := strings.Split(endpointString, ",")
	if len(endpointString) < 1 {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	vars := mux.Vars(r)

	if _, ok = vars["name"]; !ok {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	client, err := etcd.NewClient(endpoints, etcd.WithAuth(user, pass))
	if err != nil {
		if err == rpctypes.ErrAuthFailed {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	etcdUser, err := client.UserGet(r.Context(), vars["name"])
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseData := response.User{
		Name:  vars["name"],
		Roles: etcdUser.Roles,
	}
	response, err := json.Marshal(responseData)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(response)

}

func (jh *GenericHandler) AssignRole(rw http.ResponseWriter, r *http.Request) {

	user, pass, ok := r.BasicAuth()
	if !ok {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	endpointString := r.Header.Get("X-Endpoints")
	if endpointString == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	endpoints := strings.Split(endpointString, ",")
	if len(endpointString) < 1 {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	vars := mux.Vars(r)

	if _, ok = vars["name"]; !ok {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if _, ok = vars["role"]; !ok {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	client, err := etcd.NewClient(endpoints, etcd.WithAuth(user, pass))
	if err != nil {
		if err == rpctypes.ErrAuthFailed {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = client.UserGrantRole(r.Context(), vars["name"], vars["role"])
	if err != nil {
		if err == rpctypes.ErrPermissionDenied || err == rpctypes.ErrPermissionNotGranted {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}

func (jh *GenericHandler) UnassignRole(rw http.ResponseWriter, r *http.Request) {

	user, pass, ok := r.BasicAuth()
	if !ok {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	endpointString := r.Header.Get("X-Endpoints")
	if endpointString == "" {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	endpoints := strings.Split(endpointString, ",")
	if len(endpointString) < 1 {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	vars := mux.Vars(r)

	if _, ok = vars["name"]; !ok {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	if _, ok = vars["role"]; !ok {
		rw.WriteHeader(http.StatusNotAcceptable)
		return
	}

	client, err := etcd.NewClient(endpoints, etcd.WithAuth(user, pass))
	if err != nil {
		if err == rpctypes.ErrAuthFailed {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = client.UserRevokeRole(r.Context(), vars["name"], vars["role"])
	if err != nil {
		if err == rpctypes.ErrPermissionDenied || err == rpctypes.ErrPermissionNotGranted {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}