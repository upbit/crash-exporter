package websocket_test

import (
	"crash_exporter/models"
	"crash_exporter/websocket"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey" //nolint:revive // goconvey
)

func TestBaseCrash_MatchNormalLogTarget(t *testing.T) {
	testcases := []struct {
		name       string
		input      string
		expect     string
		expectType string
	}{
		{"普通日志，直接访问", "[TCP] 192.168.0.215:26858 --> 10.11.0.2:443 match GeoIP(CN) using DIRECT",
			"&{Src:192.168.0.215:26858 Dst:10.11.0.2:443 Match:GeoIP(CN) Type:DIRECT}", models.LogTypeNormal},
		{"普通日志，使用代理", "[TCP] 192.168.0.215:26859 --> encrypted-tbn0.gstatic.com:443 match DomainSuffix(gstatic.com) using UseProxy[Proxy xx]", //nolint:lll
			"&{Src:192.168.0.215:26859 Dst:encrypted-tbn0.gstatic.com:443 Match:DomainSuffix(gstatic.com) Type:UseProxy}", models.LogTypeNormal}, //nolint:lll
		{"异常日志，超时", "[TCP] dial DIRECT (match Match/) to extensions-auth.uc.r.appspot.com:443 error: dial tcp4 142.250.72.180:443: i/o timeout", //nolint:lll
			"&{Src: Dst:extensions-auth.uc.r.appspot.com:443 Match:i/o timeout Type:DIRECT}", models.LogTypeError},
		{"DNS日志", "[DNS] grafana.com --> 34.120.177.193", "&{Src: Dst:grafana.com Match:DNS Type:DIRECT}", models.LogTypeDNS}, //nolint:lll
	}

	Convey("MatchNormalLogTarget需要能正常提取日志", t, func() {
		crash, err := websocket.NewCrash("", "", nil, logrus.New())
		So(err, ShouldBeNil)
		So(crash, ShouldNotBeNil)

		for _, test := range testcases {
			actual, logType := crash.MatchLogTarget(test.input)
			So(fmt.Sprintf("%+v", actual), ShouldEqual, test.expect)
			So(logType, ShouldEqual, test.expectType)
		}
	})
}
