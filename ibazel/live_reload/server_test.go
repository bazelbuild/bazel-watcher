package live_reload

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"

	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
	"github.com/golang/protobuf/proto"
)

func TestTaretDecider(t *testing.T) {
	l := New()
	defer l.Cleanup()
	rule := &blaze_query.Rule{
		Attribute: []*blaze_query.Attribute{
			&blaze_query.Attribute{
				Name: proto.String("name"),
			},
			&blaze_query.Attribute{
				Name:            proto.String("tags"),
				Type:            blaze_query.Attribute_STRING_LIST.Enum(),
				StringListValue: []string{"ibazel_live_reload"},
			},
		},
	}

	l.TargetDecider(rule)

	if l.lrserver == nil {
		t.Errorf("Should have started live reload server")
	}
}

func TestStartLiveReloadServer(t *testing.T) {
	l := New()
	defer l.Cleanup()
	l.startLiveReloadServer()

	livereloadUrl := os.Getenv("IBAZEL_LIVERELOAD_URL")
	validUrl := regexp.MustCompile("^http\\:\\/\\/localhost\\:[0-9]+\\/livereload\\.js\\?snipver\\=1$")
	if !validUrl.MatchString(livereloadUrl) {
		t.Errorf("Invalid livereload URL '%s'", livereloadUrl)
	}

	client := new(http.Client)
	resp, err := client.Get(livereloadUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	bodyString := string(body)
	validBodyStart := regexp.MustCompile("^\\(function e\\(t\\,n\\,r\\)")
	validBodyEnd := regexp.MustCompile("\\}\\,\\{\\}\\]\\}\\,\\{\\}\\,\\[8\\]\\)\\;$")
	if !validBodyStart.MatchString(bodyString) || !validBodyEnd.MatchString(bodyString) {
		t.Errorf("Invalid livereload.js")
	}
}
