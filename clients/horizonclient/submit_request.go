package horizonclient

import (
	"net/url"

	"github.com/stellar/go/support/errors"
)

// BuildURL returns the url for submitting transactions to a running horizon instance
func (sr submitRequest) BuildURL() (string, error) {
	if sr.endpoint == "" {
		return "", errors.New("invalid request: too few parameters")
	}
	return sr.endpoint, nil
}

func (sr submitRequest) BuildBody() ([]byte, error) {
	if sr.transactionXdr == "" {
		return nil, errors.New("invalid request: submit request missing transaction xdr")
	}

	query := url.Values{}
	query.Set("tx", sr.transactionXdr)

	body := query.Encode()

	return []byte(body), nil
}
