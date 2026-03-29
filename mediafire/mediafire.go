package mediafire

import (
	"crypto/tls"
	"io"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"golang.org/x/net/html"
)

func MediafireDownload(mediafireUrl string) (data io.ReadCloser, headers http.Header, name string, err error) {
	tr := &http3.Transport{
		TLSClientConfig: &tls.Config{},  // set a TLS client config, if desired
		QUICConfig:      &quic.Config{}, // QUIC connection options
	}
	defer tr.Close()
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", mediafireUrl, nil)
	if err != nil {
		return nil, nil, "", err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:133.0) Gecko/20100101 Firefox/133.0")
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, "", err
	}
	htmlBodyResp := res.Body
	defer htmlBodyResp.Close()
	doc, err := html.Parse(htmlBodyResp)
	if err != nil {
		return nil, nil, "", err
	}
	downloadUrl := findHref(doc, "Download file")

	nameNode := find(doc, "filename", "class")
	fileName := "Unknown name"
	if nameNode != nil && nameNode.FirstChild != nil {
		fileName = nameNode.FirstChild.Data
	} else {
		fileName, _ = url.QueryUnescape(filepath.Base(downloadUrl))
	}

	fileReq, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return nil, nil, fileName, err
	}

	fileReq.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:121.0) Gecko/20100101 Firefox/121.0")

	client = &http.Client{}
	fileRes, err := client.Do(fileReq)
	if err != nil {
		return nil, nil, fileName, err
	}

	return fileRes.Body, fileRes.Header, fileName, nil
}

func findHref(n *html.Node, val string) string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "aria-label" && attr.Val == val {
				for _, a := range n.Attr {
					if a.Key == "href" {
						return a.Val
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if href := findHref(c, val); href != "" {
			return href
		}
	}
	return ""
}

func find(n *html.Node, val string, key string) *html.Node {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == key && attr.Val == val {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if htmlK := find(c, val, key); htmlK != nil {
			return htmlK
		}
	}
	return nil
}
