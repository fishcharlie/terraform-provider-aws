package create

import (
	"testing"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestStringHashcode(t *testing.T) {
	v := "hello, world"
	expected := StringHashcode(v)
	for i := 0; i < 100; i++ {
		actual := StringHashcode(v)
		if actual != expected {
			t.Fatalf("bad: %#v\n\t%#v", actual, expected)
		}
	}
}

func TestStringHashcode_positiveIndex(t *testing.T) {
	// "2338615298" hashes to uint32(2147483648) which is math.MinInt32
	ips := []string{"192.168.1.3", "192.168.1.5", "2338615298"}
	for _, ip := range ips {
		if index := StringHashcode(ip); index < 0 {
			t.Fatalf("Bad Index %#v for ip %s", index, ip)
		}
	}
}
