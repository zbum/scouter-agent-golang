package strace

import "net/http"

type ScouterRoundTripper struct {
	Proxied http.RoundTripper
}

func (s ScouterRoundTripper) RoundTrip(req *http.Request) (res *http.Response, err error) {
	ctx := req.Context()
	step := StartApiCall(ctx, req.Host, req.URL.String())
	res, err = s.Proxied.RoundTrip(req)
	EndApiCall(ctx, step, err)
	return
}
