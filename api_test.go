package http_test

import (
	"fmt"
	api "github.com/phiskills/http-api.go"
	"io/ioutil"
	"net/http"
	"testing"
)

const httpPort = 8080

func TestApi(t *testing.T) {
	fmt.Println("# Initialize test API")
	a := api.New("test-api")
	a.UsePort(httpPort)
	if a.Status() != api.Unknown {
		t.Fatalf("New API failed: wrong status, %s should be %s", a.Status(), api.Unknown)
	}
	fmt.Println("# Start API")
	go a.Start()
	status := a.Change()
	if a.Status() != api.Serving || a.Status() != status {
		t.Fatalf("Start API failed: wrong status, a.status %s & a.change %s should be %s", a.Status(), status, api.Serving)
	}
	fmt.Println("# Validate health check services")
	validate := buildValidator(t, httpPort, "check")
	for _, name := range []string{"GET", "OPTIONS"} {
		validate(name, "{\"status\":\"SERVING\"}")
	}
	fmt.Println("# Register services")
	a.Register("/test", &api.Router{
		Head: func(context *api.Context) {
			context.Res.WriteHeader(http.StatusAccepted)
		},
		Get:     buildRoute("GET"),
		Post:    buildRoute("POST"),
		Put:     buildRoute("PUT"),
		Patch:   buildRoute("PATCH"),
		Delete:  buildRoute("DELETE"),
		Connect: buildRoute("CONNECT"),
		Options: buildRoute("OPTIONS"),
		Trace:   buildRoute("TRACE"),
	})
	fmt.Println("# Validate registered services")
	validate = buildValidator(t, httpPort, "test")
	for _, name := range []string{"GET", "POST", "PUT", "PATCH", "DELETE", "CONNECT", "OPTIONS", "TRACE"} {
		validate(name, fmt.Sprintf("%s Request", name))
	}
	fmt.Println("# Stop API")
	a.Stop()
	status = a.Change()
	if a.Status() != api.NotServing || a.Status() != status {
		t.Fatalf("Stop API failed: wrong status, a.status %s & a.change %s should be %s", a.Status(), status, api.NotServing)
	}
}

func buildRoute(name string) func(context *api.Context) {
	return func(context *api.Context) {
		message := fmt.Sprintf("%s Request", name)
		context.Res.Write([]byte(message))
	}
}

func buildValidator(t *testing.T, port int, target string) func(string, string) {
	client := &http.Client{}
	url := fmt.Sprintf("http://localhost:%d/%s", port, target)
	return func(name string, result string) {
		fmt.Printf("- Validate %s\t%s\n", url, name)
		req, err := http.NewRequest(name, url, nil)
		if err != nil {
			t.Errorf("%s failed: %v", name, err)
			return
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Errorf("%s failed: %v", name, err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s failed: wrong status %d, expected %d", name, resp.StatusCode, http.StatusOK)
			return
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("%s failed: could not read response: %v", name, err)
			return
		}
		if string(body) != result {
			t.Errorf("%s failed: wrong response %s, expected %s", name, body, result)
		}
	}
}
