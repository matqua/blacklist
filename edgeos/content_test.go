package edgeos

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/britannic/testutils"
)

func TestGetContent(t *testing.T) {
	// var (
	// 	d     = new(HTTPserver)
	// 	dpage = "/domains.txt"
	// 	h     = new(HTTPserver)
	// 	hpage = "/hosts.txt"
	// )

	tests := []struct {
		leaf     string
		node     string
		page     string
		pageData string
		stype    string
		svr      *HTTPserver
		url      string
		want     string
	}{
		{
			leaf:     "malc0de",
			node:     domains,
			page:     "/domains.txt",
			pageData: HTTPDomainData,
			stype:    preConf,
			svr:      new(HTTPserver),
			want:     "adsrvr.org\nadtechus.net\nadvertising.com\ncentade.com\ndoubleclick.net\nfree-counter.co.uk\nintellitxt.com\nkiosked.com",
		},
		{
			leaf:     "adaway",
			node:     hosts,
			page:     "/hosts.txt",
			pageData: httpHostData,
			stype:    preConf,
			svr:      new(HTTPserver),
			want:     "beap.gemini.yahoo.com",
		},
	}

	c, err := ReadCfg(bytes.NewBufferString(Cfg))
	OK(t, err)
	NewParms(c).SetOpt(
		Dir("/tmp"),
		Ext("blacklist.conf"),
		Method("GET"),
		Nodes([]string{"domains", "hosts"}),
		STypes([]string{"pre-configured", "files", "urls"}),
	)

	for _, test := range tests {
		test.url = test.svr.NewHTTPServer().String() + test.page
		test.svr.Mux.HandleFunc(test.page,
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, test.pageData)
			},
		)
		c.bNodes[test.node].data[test.leaf].url = test.url
		for _, src := range *c.Get(test.node).Source(test.stype).GetContent() {
			got, err := ioutil.ReadAll(src.r)
			OK(t, err)
			Equals(t, test.want, string(got))
		}
	}
}

var (
	// Cfg contains a valid full EdgeOS blacklist configuration
	Cfg = `blacklist {
    disabled false
    dns-redirect-ip 0.0.0.0
    domains {
        dns-redirect-ip 0.0.0.0
        include adsrvr.org
        include adtechus.net
        include advertising.com
        include centade.com
        include doubleclick.net
        include free-counter.co.uk
        include intellitxt.com
        include kiosked.com
        source malc0de {
            description "List of zones serving malicious executables observed by malc0de.com/database/"
            prefix "zone "
            url http://malc0de.com/bl/ZONES
        }
    }
    exclude 122.2o7.net
    exclude 1e100.net
    exclude adobedtm.com
    exclude akamai.net
    exclude amazon.com
    exclude amazonaws.com
    exclude apple.com
    exclude ask.com
    exclude avast.com
    exclude bitdefender.com
    exclude cdn.visiblemeasures.com
    exclude cloudfront.net
    exclude coremetrics.com
    exclude edgesuite.net
    exclude freedns.afraid.org
    exclude github.com
    exclude githubusercontent.com
    exclude google.com
    exclude googleadservices.com
    exclude googleapis.com
    exclude googleusercontent.com
    exclude gstatic.com
    exclude gvt1.com
    exclude gvt1.net
    exclude hb.disney.go.com
    exclude hp.com
    exclude hulu.com
    exclude images-amazon.com
    exclude msdn.com
    exclude paypal.com
    exclude rackcdn.com
    exclude schema.org
    exclude skype.com
    exclude smacargo.com
    exclude sourceforge.net
    exclude ssl-on9.com
    exclude ssl-on9.net
    exclude static.chartbeat.com
    exclude storage.googleapis.com
    exclude windows.net
    exclude yimg.com
    exclude ytimg.com
    hosts {
        dns-redirect-ip 192.168.168.1
        include beap.gemini.yahoo.com
        source adaway {
            description "Blocking mobile ad providers and some analytics providers"
            prefix "127.0.0.1 "
            url http://adaway.org/hosts.txt
        }
    }
}`
)