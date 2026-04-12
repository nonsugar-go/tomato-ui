package paloalto

import (
	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func fillModel(app *model.App, palo *PaloAltoConfig) error {
	var err error

	if app.Tag, err = ToModelTags(palo.TagObject); err != nil {
		return err
	}

	if app.Addresses, err = ToModelAddresses(palo.Addresses); err != nil {
		return err
	}

	if app.AddressGroups, err = ToModelAddressGroups(palo.AddressGroups); err != nil {
		return err
	}

	// app.Services, err = ...
	// app.SecurityRules, err = ...
	// app.NATRules, err = ...

	return nil
}

// ToModelTags converts Palo Alto tag objects to model tags.
// NOTE: untested
func ToModelTags(tagObjects []ScopedTagObject) ([]model.Tag, error) {
	var result []model.Tag

	for _, st := range tagObjects {
		t := st.TagObject

		tag := model.Tag{
			Value:       t.Name,
			Color:       t.Color,
			Description: t.Comments,
		}

		result = append(result, tag)
	}

	return result, nil
}

func ToModelAddresses(scopedAddrs []ScopedAddress) ([]model.Address, error) {
	var result []model.Address

	for _, sa := range scopedAddrs {
		a := sa.Address

		addr := model.Address{
			Name:        a.Name,
			Description: a.Description,
		}

		switch {
		case a.IPNetmask != "":
			addr.Type = model.AddressTypeIPNetmask
			addr.Value = a.IPNetmask

		case a.FQDN != "":
			addr.Type = model.AddressTypeFQDN
			addr.Value = a.FQDN

		default:
			addr.Type = model.AddressTypeUnknown
			addr.Value = ""
		}

		result = append(result, addr)
	}

	return result, nil
}

func ToModelAddressGroups(scopedAddrGrps []ScopedAddressGroup) ([]model.AddressGroup, error) {
	var result []model.AddressGroup

	for _, sg := range scopedAddrGrps {
		g := sg.Group

		grp := model.AddressGroup{
			Name: g.Name,
			// Description: g.Description,
		}

		switch {
		case len(g.Static) != 0:
			grp.Members = g.Static
		case g.Dynamic != nil:
			grp.Members = []string{g.Dynamic.Filter}
		default:
			grp.Members = []string{}
		}

		result = append(result, grp)
	}

	return result, nil
}
