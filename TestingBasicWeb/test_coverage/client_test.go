package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type TestCase struct {
	AccessToken    string
	Request        SearchRequest
	ExpectedResult *SearchResponse
	ExpectedError  error
	DataFilePath   string
}

func testOneCase(c *TestCase) error {
	ts := httptest.NewServer(http.HandlerFunc(SearchServer))
	client := SearchClient{
		AccessToken: c.AccessToken,
		URL:         ts.URL,
	}
	DatafilePath = c.DataFilePath
	res, err := client.FindUsers(c.Request)

	falsePositiveError := c.ExpectedError == nil && err != nil
	falseNegativeError := c.ExpectedError != nil && err == nil

	if falsePositiveError || falseNegativeError || (err != nil && c.ExpectedError.Error() != err.Error()) {
		reserr := fmt.Sprintf("Expected errors doesn't match\nError: %v, Expected: %v", err, c.ExpectedError)
		return errors.New(reserr)
	}

	if !reflect.DeepEqual(res, c.ExpectedResult) {
		reserr := fmt.Sprintf("Results doesn't match\nResult: %v, Expected: %v", res, c.ExpectedResult)
		return errors.New(reserr)
	}

	return nil
}

func TestBadToken(t *testing.T) {
	c := TestCase{
		AccessToken: "TotallyWrongToken",
		Request: SearchRequest{
			Limit:      1,
			Offset:     1,
			Query:      "",
			OrderField: "Name",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  fmt.Errorf("Bad AccessToken"),
		DataFilePath:   "test_data.xml",
	}

	err := testOneCase(&c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestIntSort(t *testing.T) {
	cases := []TestCase{
		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "",
				OrderField: "Age",
				OrderBy:    OrderByAsc,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     1,
						Name:   "Bob",
						Age:    21,
						About:  "Test data 3",
						Gender: "female",
					},
					User{
						Id:     2,
						Name:   "Zord",
						Age:    23,
						About:  "Test data 2",
						Gender: "male",
					},
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},

		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "",
				OrderField: "Age",
				OrderBy:    OrderByDesc,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
					User{
						Id:     2,
						Name:   "Zord",
						Age:    23,
						About:  "Test data 2",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "Bob",
						Age:    21,
						About:  "Test data 3",
						Gender: "female",
					},
				},
				NextPage: false,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},

		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "",
				OrderField: "Age",
				OrderBy:    OrderByAsIs,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "Bob",
						Age:    21,
						About:  "Test data 3",
						Gender: "female",
					},
					User{
						Id:     2,
						Name:   "Zord",
						Age:    23,
						About:  "Test data 2",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},
	}

	for _, c := range cases {
		err := testOneCase(&c)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestStringSort(t *testing.T) {
	cases := []TestCase{
		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "Bob",
						Age:    21,
						About:  "Test data 3",
						Gender: "female",
					},
					User{
						Id:     2,
						Name:   "Zord",
						Age:    23,
						About:  "Test data 2",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},

		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByDesc,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     2,
						Name:   "Zord",
						Age:    23,
						About:  "Test data 2",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "Bob",
						Age:    21,
						About:  "Test data 3",
						Gender: "female",
					},
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},

		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      3,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByAsIs,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "Bob",
						Age:    21,
						About:  "Test data 3",
						Gender: "female",
					},
					User{
						Id:     2,
						Name:   "Zord",
						Age:    23,
						About:  "Test data 2",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},
	}

	for _, c := range cases {
		err := testOneCase(&c)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestWrongField(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      3,
			Offset:     0,
			Query:      "",
			OrderField: "WrongField",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  fmt.Errorf("OrderField WrongField invalid"),
		DataFilePath:   "test_data.xml",
	}

	err := testOneCase(&c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestExactQuery(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Zord",
			OrderField: "",
			OrderBy:    1,
		},
		ExpectedResult: &SearchResponse{
			Users: []User{
				User{
					Id:     2,
					Name:   "Zord",
					Age:    23,
					About:  "Test data 2",
					Gender: "male",
				},
			},
			NextPage: false,
		},
		ExpectedError: nil,
		DataFilePath:  "test_data.xml",
	}

	err := testOneCase(&c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUnableToReadFile(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Zord",
			OrderField: "",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  fmt.Errorf("SearchServer fatal error"),
		DataFilePath:   "unexisting_data.xml",
	}

	err := testOneCase(&c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestUnableToParseFile(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Zord",
			OrderField: "",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  fmt.Errorf("SearchServer fatal error"),
		DataFilePath:   "wrong_data.xml",
	}

	err := testOneCase(&c)
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestWrongRequest(t *testing.T) {
	cases := []TestCase{
		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      -2,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByDesc,
			},
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("limit must be > 0"),
			DataFilePath:   "test_data.xml",
		},

		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      0,
				Offset:     -2,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("offset must be > 0"),
			DataFilePath:   "test_data.xml",
		},
	}

	for _, c := range cases {
		err := testOneCase(&c)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestLimits(t *testing.T) {
	cases := []TestCase{
		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      28,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
					User{
						Id:     1,
						Name:   "Bob",
						Age:    21,
						About:  "Test data 3",
						Gender: "female",
					},
					User{
						Id:     2,
						Name:   "Zord",
						Age:    23,
						About:  "Test data 2",
						Gender: "male",
					},
				},
				NextPage: false,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},

		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      1,
				Offset:     0,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			ExpectedResult: &SearchResponse{
				Users: []User{
					User{
						Id:     0,
						Name:   "Anna",
						Age:    24,
						About:  "Test data 1",
						Gender: "male",
					},
				},
				NextPage: true,
			},
			ExpectedError: nil,
			DataFilePath:  "test_data.xml",
		},

		TestCase{
			AccessToken: CorrectAccessToken,
			Request: SearchRequest{
				Limit:      0,
				Offset:     200,
				Query:      "",
				OrderField: "Name",
				OrderBy:    OrderByAsc,
			},
			ExpectedResult: nil,
			ExpectedError:  fmt.Errorf("cant unpack error json: unexpected end of JSON input"),
			DataFilePath:   "test_data.xml",
		},
	}

	for _, c := range cases {
		err := testOneCase(&c)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestNoConnectionToServer(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Zord",
			OrderField: "",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  nil,
		DataFilePath:   "wrong_data.xml",
	}

	client := SearchClient{
		AccessToken: c.AccessToken,
		URL:         "http://localhost:8000",
	}

	res, err := client.FindUsers(c.Request)
	if res != nil || err == nil {
		t.Errorf("Should have connection error")
	}
}

func DummyHandlerWithError(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	resp := SearchErrorResponse{
		Error: "UnknownTestError",
	}
	respData, err := json.Marshal(&resp)
	if err == nil {
		w.Write(respData)
	}
	return
}

func DummyHandlerWithWrongData(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SomeByteMess"))
	return
}

func DummyHandlerTimeout(w http.ResponseWriter, r *http.Request) {
	time.Sleep(10 * time.Second)
	return
}

func TestUnknownBadRequestError(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Zord",
			OrderField: "",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  fmt.Errorf("unknown bad request error: UnknownTestError"),
		DataFilePath:   "wrong_data.xml",
	}

	ts := httptest.NewServer(http.HandlerFunc(DummyHandlerWithError))
	client := SearchClient{
		AccessToken: c.AccessToken,
		URL:         ts.URL,
	}

	_, err := client.FindUsers(c.Request)
	if err.Error() != c.ExpectedError.Error() {
		t.Errorf("Expected errors doesn't match\nError: %v, Expected: %v", err, c.ExpectedError)
	}
}

func TestWrongDataInResponse(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Zord",
			OrderField: "",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  fmt.Errorf("cant unpack result json: invalid character 'S' looking for beginning of value"),
		DataFilePath:   "wrong_data.xml",
	}

	ts := httptest.NewServer(http.HandlerFunc(DummyHandlerWithWrongData))
	client := SearchClient{
		AccessToken: c.AccessToken,
		URL:         ts.URL,
	}

	_, err := client.FindUsers(c.Request)
	if err.Error() != c.ExpectedError.Error() {
		t.Errorf("Expected errors doesn't match\nError: %v, Expected: %v", err, c.ExpectedError)
	}
}

func TestWTimeout(t *testing.T) {
	c := TestCase{
		AccessToken: CorrectAccessToken,
		Request: SearchRequest{
			Limit:      1,
			Offset:     0,
			Query:      "Zord",
			OrderField: "",
			OrderBy:    1,
		},
		ExpectedResult: nil,
		ExpectedError:  fmt.Errorf("timeout for limit=2&offset=0&order_by=1&order_field=&query=Zord"),
		DataFilePath:   "wrong_data.xml",
	}

	ts := httptest.NewServer(http.HandlerFunc(DummyHandlerTimeout))
	client := SearchClient{
		AccessToken: c.AccessToken,
		URL:         ts.URL,
	}

	_, err := client.FindUsers(c.Request)
	if err.Error() != c.ExpectedError.Error() {
		t.Errorf("Expected errors doesn't match\nError: %v, Expected: %v", err, c.ExpectedError)
	}
}
