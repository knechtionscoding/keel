package http

import (
	"bytes"
	"net/http"

	"net/http/httptest"
	"testing"
)

var fakeRequest = `{
	"push_data": {
		"pushed_at": 1497467660,
		"images": [],
		"tag": "0.1.7",
		"pusher": "karolisr"
	},
	"callback_url": "https://registry.hub.docker.com/u/karolisr/keel/hook/22hagb51h1gfb4eefc5f1g4j3abi0beg4/",
	"repository": {
		"status": "Active",
		"description": "",
		"is_trusted": false,
		"full_description": "desc",
		"repo_url": "https://hub.docker.com/r/karolisr/keel",
		"owner": "karolisr",
		"is_official": false,
		"is_private": false,
		"name": "keel",
		"namespace": "karolisr",
		"star_count": 0,
		"comment_count": 0,
		"date_created": 1497032538,
		"dockerfile": "FROM golang:1.8.1-alpine\nCOPY . /go/src/github.com/datagravity-ai/keel\nWORKDIR /go/src/github.com/datagravity-ai/keel\nRUN apk add --no-cache git && go get\nRUN CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags -'w' -o keel .\n\nFROM alpine:latest\nRUN apk --no-cache add ca-certificates\nCOPY --from=0 /go/src/github.com/datagravity-ai/keel/keel /bin/keel\nENTRYPOINT [\"/bin/keel\"]\n\nEXPOSE 9300",
		"repo_name": "karolisr/keel"
	}
}`

func TestDockerhubWebhookHandler(t *testing.T) {

	fp := &fakeProvider{}
	srv, teardown := NewTestingServer(fp)
	defer teardown()

	req, err := http.NewRequest("POST", "/v1/webhooks/dockerhub", bytes.NewBuffer([]byte(fakeRequest)))
	if err != nil {
		t.Fatalf("failed to create req: %s", err)
	}

	//The response recorder used to record HTTP responses
	rec := httptest.NewRecorder()

	srv.router.ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Errorf("unexpected status code: %d", rec.Code)

		t.Log(rec.Body.String())
	}

	if len(fp.submitted) != 1 {
		t.Fatalf("unexpected number of events submitted: %d", len(fp.submitted))
	}

	if fp.submitted[0].Repository.Name != "karolisr/keel" {
		t.Errorf("expected karolisr/keel but got %s", fp.submitted[0].Repository.Name)
	}

	if fp.submitted[0].Repository.Tag != "0.1.7" {
		t.Errorf("expected 0.1.7 but got %s", fp.submitted[0].Repository.Tag)
	}
}
