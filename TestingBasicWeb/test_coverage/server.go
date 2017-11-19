package main

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

const (
	CorrectAccessToken = "d41d8cd98f00b204e9800998ecf8427e"
)

type Users struct {
	Users     []User `xml:"row"`
	sortField string
	orderBy   int
}

func (users Users) Len() int {
	return len(users.Users)
}

func (users Users) Less(i, j int) bool {
	valI := reflect.ValueOf(users.Users[i]).FieldByName(users.sortField)
	valJ := reflect.ValueOf(users.Users[j]).FieldByName(users.sortField)
	var res bool
	if valI.Kind() == reflect.String {
		if users.orderBy == OrderByAsc {
			res = (strings.Compare(valI.String(), valJ.String()) == -1)
		} else if users.orderBy == OrderByDesc {
			res = (strings.Compare(valI.String(), valJ.String()) == 1)
		} else {
			res = false
		}
	} else if valI.Kind() == reflect.Int {
		if users.orderBy == OrderByAsc {
			res = valI.Int() < valJ.Int()
		} else if users.orderBy == OrderByDesc {
			res = valI.Int() > valJ.Int()
		} else {
			res = false
		}
	}
	return res
}

func (users Users) Swap(i, j int) {
	reflect.Swapper(users.Users)(i, j)
}

var DatafilePath string = "dataset.xml"

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func SearchServer(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("AccessToken") != CorrectAccessToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	query := r.URL.Query().Get("query")
	order_field := r.URL.Query().Get("order_field")
	if len(order_field) == 0 {
		order_field = "Name"
	} else if order_field != "Id" && order_field != "Age" && order_field != "Name" {
		w.WriteHeader(http.StatusBadRequest)
		resp := SearchErrorResponse{
			Error: ErrorBadOrderField,
		}
		respData, err := json.Marshal(&resp)
		if err == nil {
			w.Write(respData)
		}
		return
	}
	order_by, err := strconv.Atoi(r.URL.Query().Get("order_by"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadFile(DatafilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	users := Users{}
	err = xml.Unmarshal(data, &users)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := Users{}
	if len(query) > 0 {
		for _, user := range users.Users {
			if user.Name == query || user.About == query {
				res.Users = append(res.Users, user)
			}
		}
	} else {
		res = users
	}

	res.sortField = order_field
	res.orderBy = order_by
	sort.Stable(res)

	if limit == 0 || offset < 0 || offset > len(res.Users)-1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	respData, err := json.Marshal(res.Users[offset:min(offset+limit, len(res.Users))])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respData)
}
