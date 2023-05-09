package types // nolint: dupl

import "github.com/FishtechCSOC/locomotive/v8/pkg/helpers"

type Tenant struct {
	ID        string `json:"id" mapstructure:"id"`
	Reference string `json:"reference" mapstructure:"reference"` // maps to 'Tenant Ref' in CMS
	Name      string `json:"name" mapstructure:"name"`
}

func (tenant *Tenant) MergeLeft(other Tenant) Tenant {
	return Tenant{
		Reference: helpers.DefaultString(other.Reference, tenant.Reference),
		ID:        helpers.DefaultString(other.ID, tenant.ID),
		Name:      helpers.DefaultString(other.Name, tenant.Name),
	}
}

func (tenant *Tenant) MergeRight(other Tenant) Tenant {
	return Tenant{
		Reference: helpers.DefaultString(tenant.Reference, other.Reference),
		ID:        helpers.DefaultString(tenant.ID, other.ID),
		Name:      helpers.DefaultString(tenant.Name, other.Name),
	}
}

func (tenant *Tenant) Empty() bool {
	return tenant.Reference == "" && tenant.ID == "" && tenant.Name == ""
}
