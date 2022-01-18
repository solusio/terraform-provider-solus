package provider

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/solusio/solus-go-sdk"
)

// validationIsDomainName checks that specified value is valid domain name.
// Example: example.com.
func validationIsDomainName(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if !isDomainName(v) {
		return nil, []error{errors.New("invalid domain name")}
	}
	return nil, nil
}

func isDomainName(v string) bool {
	// To simplify implementations, the total number of octets that represent a
	// domain name (i.e., the sum of all label octets and label lengths) is
	// limited to 255.
	// https://datatracker.ietf.org/doc/html/rfc1034#section-3.1
	return len(strings.ReplaceAll(v, ".", "")) <= 255 &&
		net.ParseIP(v) == nil &&
		rDomainName.MatchString(v)
}

//goland:noinspection RegExpRedundantEscape
var rDomainName = regexp.MustCompile(`^([a-zA-Z0-9_][a-zA-Z0-9_-]{0,62})(\.[a-zA-Z0-9_][a-zA-Z0-9_-]{0,62})*[\._]?$`)

func validationIsVirtualizationType(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}
	if !solus.IsValidVirtualizationType(v) {
		return nil, []error{fmt.Errorf("invalid virtualization type %q", v)}
	}
	return nil, nil
}
