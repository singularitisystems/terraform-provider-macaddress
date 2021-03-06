package main

import (
	"crypto/rand"
	"fmt"
	"strings"
        "errors"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAddress() *schema.Resource {
	return &schema.Resource{
		Create: resourceAddressCreate,
		Read:   readNil,
		Delete: schema.RemoveFromState,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"prefix": &schema.Schema{
				Type:     schema.TypeList,
				Elem: &schema.Schema{
					Type:     schema.TypeInt,
				},
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

const MAC_ADDRESS_LENGTH = 6

func resourceAddressCreate(d *schema.ResourceData, m interface{}) error {
	var groups []string
	buf := make([]byte, MAC_ADDRESS_LENGTH)

	_, err := rand.Read(buf)
	if err != nil {
		return errwrap.Wrapf("error generating random bytes: {{err}}", err)
	}

	prefix := d.Get("prefix").([]interface{})

	if len(prefix) > MAC_ADDRESS_LENGTH {
		return errors.New("error generating random mac address: prefix is too large")
	}

	for index, val := range prefix {
		if val.(int) > 255 {
			return errors.New("error generating random mac address: prefix segment must be in the range [0,256)")
		}
		buf[index] = byte(val.(int))
	}

	// Locally administered
	buf[0] |= 0x02

	// Unicast
	buf[0] &= 0xfe

	for _, i := range buf {
		groups = append(groups, fmt.Sprintf("%02x", i))
	}

	address := strings.Join(groups, ":")

	d.SetId(address)
	d.Set("address", address)

	return nil
}

func readNil(d *schema.ResourceData, meta interface{}) error {
	return nil
}
