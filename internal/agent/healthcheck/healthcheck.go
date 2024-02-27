package healthcheck

import "github.com/go-resty/resty/v2"

func MakeHealthcheck(addr string) bool {
	client := resty.New()
	_, err := client.R().
		Head(addr)
	return err == nil
}
