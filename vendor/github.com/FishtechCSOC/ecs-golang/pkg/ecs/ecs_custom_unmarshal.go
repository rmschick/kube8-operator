// nolint: dupl, gocyclo, funlen, gocognit, cyclop, forcetypeassert
package ecs

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

const ElasticVersion = "1.7.0"

func nullcheck(obj map[string]interface{}, keyword string) bool {
	return obj[keyword] != nil
}

func safeCastInt(data interface{}) (int64, error) {
	if data == nil {
		return -1, errors.New("provided value in nil")
	}

	switch val := data.(type) {
	case int64, int32, int, uint64, uint32, uint:
		return val.(int64), nil
	case string:
		return strconv.ParseInt(val, 10, 64) //nolint: gomnd
	default:
		return -1, fmt.Errorf("failed to convert value to a int: %v", data)
	}
}

func safeCastFloat(data interface{}) (float64, error) {
	if data == nil {
		return -1, errors.New("provided value in nil")
	}

	switch val := data.(type) {
	case float64, float32, int64, int32, int, uint64, uint32, uint:
		return val.(float64), nil
	case string:
		return strconv.ParseFloat(val, 64) //nolint: gomnd
	default:
		return -1, fmt.Errorf("failed to convert value to a float: %v", data)
	}
}

func (b *Event) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "action") {
		b.Action = item["action"].(string)
	}

	if nullcheck(item, "category") {
		if reflect.TypeOf(item["category"]).Kind() == reflect.String {
			b.Category = append(b.Category, item["category"].(string))
		} else {
			for _, v := range item["category"].([]interface{}) {
				b.Category = append(b.Category, v.(string))
			}
		}
	}

	if nullcheck(item, "code") {
		b.Code = item["code"].(string)
	}

	if nullcheck(item, "created") {
		b.Created, _ = time.Parse(time.RFC3339, item["created"].(string))
	}

	if nullcheck(item, "dataset") {
		b.Dataset = item["dataset"].(string)
	}

	if nullcheck(item, "duration") {
		b.Duration, _ = time.ParseDuration(item["duration"].(string))
	}

	if nullcheck(item, "end") {
		b.End, _ = time.Parse(time.RFC3339, item["end"].(string))
	}

	if nullcheck(item, "hash") {
		b.Hash = item["hash"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "ingested") {
		b.Ingested, _ = time.Parse(time.RFC3339, item["ingested"].(string))
	}

	if nullcheck(item, "kind") {
		b.Kind = item["kind"].(string)
	}

	if nullcheck(item, "module") {
		b.Module = item["module"].(string)
	}

	if nullcheck(item, "original") {
		b.Original = item["original"].(string)
	}

	if nullcheck(item, "outcome") {
		b.Outcome = item["outcome"].(string)
	}

	if nullcheck(item, "provider") {
		b.Provider = item["provider"].(string)
	}

	if nullcheck(item, "reason") {
		b.Reason = item["reason"].(string)
	}

	if nullcheck(item, "reference") {
		b.Reference = item["reference"].(string)
	}

	if nullcheck(item, "risk_score") {
		val, err := safeCastFloat(item["risk_score"])
		if err != nil {
			return err
		}

		b.RiskScore = json.Number(fmt.Sprintf("%f", val))
	}

	if nullcheck(item, "risk_score_norm") {
		val, err := safeCastFloat(item["risk_score_norm"])
		if err != nil {
			return err
		}

		b.RiskScoreNorm = json.Number(fmt.Sprintf("%f", val))
	}

	if nullcheck(item, "sequence") {
		val, err := safeCastInt(item["sequence"])
		if err != nil {
			return err
		}

		b.Sequence = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "severity") {
		val, err := safeCastFloat(item["severity"])
		if err != nil {
			return err
		}

		b.Severity = json.Number(fmt.Sprintf("%f", val))
	}

	if nullcheck(item, "start") {
		b.Start, _ = time.Parse(time.RFC3339, item["start"].(string))
	}

	if nullcheck(item, "timezone") {
		b.Timezone = item["timezone"].(string)
	}

	if nullcheck(item, "type") {
		if reflect.TypeOf(item["type"]).Kind() == reflect.String {
			b.Type = append(b.Type, item["type"].(string))
		} else {
			for _, v := range item["type"].([]interface{}) {
				b.Type = append(b.Type, v.(string))
			}
		}
	}

	if nullcheck(item, "url") {
		b.URL = item["url"].(string)
	}

	return nil
}

func (b *ClientUser) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "domain") {
		b.Domain = item["domain"].(string)
	}

	if nullcheck(item, "email") {
		b.Email = item["email"].(string)
	}

	if nullcheck(item, "full_name") {
		b.FullName = item["full_name"].(string)
	}

	if nullcheck(item, "group") {
		jsonData, _ := json.Marshal(item["group"])

		err := json.Unmarshal(jsonData, &b.Group)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "hash") {
		b.Hash = item["hash"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "roles") {
		if reflect.TypeOf(item["roles"]).Kind() == reflect.String {
			b.Roles = append(b.Roles, item["roles"].(string))
		} else {
			for _, v := range item["roles"].([]interface{}) {
				b.Roles = append(b.Roles, v.(string))
			}
		}
	}

	return nil
}

func (b *DestinationUser) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "domain") {
		b.Domain = item["domain"].(string)
	}

	if nullcheck(item, "email") {
		b.Email = item["email"].(string)
	}

	if nullcheck(item, "full_name") {
		b.FullName = item["full_name"].(string)
	}

	if nullcheck(item, "group") {
		err := json.Unmarshal(data, &b.Group)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "hash") {
		b.Hash = item["hash"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "roles") {
		if reflect.TypeOf(item["roles"]).Kind() == reflect.String {
			b.Roles = append(b.Roles, item["roles"].(string))
		} else {
			for _, v := range item["roles"].([]interface{}) {
				b.Roles = append(b.Roles, v.(string))
			}
		}
	}

	return nil
}

func (b *DNS) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "answers") {
		jsonData, _ := json.Marshal(item["answers"])

		err := json.Unmarshal(jsonData, &b.Answers)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "header_flags") {
		if reflect.TypeOf(item["header_flags"]).Kind() == reflect.String {
			b.HeaderFlags = append(b.HeaderFlags, item["header_flags"].(string))
		} else {
			for _, v := range item["header_flags"].([]interface{}) {
				b.HeaderFlags = append(b.HeaderFlags, v.(string))
			}
		}
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "op_code") {
		b.OpCode = item["op_code"].(string)
	}

	if nullcheck(item, "question") {
		jsonData, _ := json.Marshal(item["question"])

		err := json.Unmarshal(jsonData, &b.Question)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "resolved_ip") {
		if reflect.TypeOf(item["resolved_ip"]).Kind() == reflect.String {
			b.ResolvedIP = append(b.ResolvedIP, item["resolved_ip"].(string))
		} else {
			for _, v := range item["resolved_ip"].([]interface{}) {
				b.ResolvedIP = append(b.ResolvedIP, v.(string))
			}
		}
	}

	if nullcheck(item, "response_code") {
		b.ResponseCode = item["response_code"].(string)
	}

	if nullcheck(item, "type") {
		b.Type = item["type"].(string)
	}

	return nil
}

func (b *File) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "accessed") {
		b.Accessed, _ = time.Parse(time.RFC3339, item["accessed"].(string))
	}

	if reflect.TypeOf(item["attributes"]).Kind() == reflect.String {
		b.Attributes = append(b.Attributes, item["attributes"].(string))
	}

	if nullcheck(item, "code_signature") {
		jsonData, _ := json.Marshal(item["code_signature"])

		err := json.Unmarshal(jsonData, &b.CodeSignature)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "created") {
		b.Created, _ = time.Parse(time.RFC3339, item["created"].(string))
	}

	if nullcheck(item, "ctime") {
		b.Ctime, _ = time.Parse(time.RFC3339, item["ctime"].(string))
	}

	if nullcheck(item, "device") {
		b.Device = item["device"].(string)
	}

	if nullcheck(item, "directory") {
		b.Directory = item["directory"].(string)
	}

	if nullcheck(item, "drive_letter") {
		b.DriveLetter = item["drive_letter"].(string)
	}

	if nullcheck(item, "extension") {
		b.Extension = item["extension"].(string)
	}

	if nullcheck(item, "group") {
		b.Group = item["group"].(string)
	}

	if nullcheck(item, "hash") {
		jsonData, _ := json.Marshal(item["hash"])

		err := json.Unmarshal(jsonData, &b.Hash)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "inode") {
		b.Inode = item["inode"].(string)
	}

	if nullcheck(item, "mime_type") {
		b.MIMEType = item["mime_type"].(string)
	}

	if nullcheck(item, "mode") {
		b.Mode = item["mode"].(string)
	}

	if nullcheck(item, "created") {
		b.Mtime, _ = time.Parse(time.RFC3339, item["created"].(string))
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "owner") {
		b.Owner = item["owner"].(string)
	}

	if nullcheck(item, "path") {
		b.Path = item["path"].(string)
	}

	if nullcheck(item, "pe") {
		jsonData, _ := json.Marshal(item["pe"])

		err := json.Unmarshal(jsonData, &b.PE)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "size") {
		val, err := safeCastInt(item["size"])
		if err != nil {
			return err
		}

		b.Size = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "target_path") {
		b.TargetPath = item["target_path"].(string)
	}

	if nullcheck(item, "type") {
		b.Type = item["type"].(string)
	}

	if nullcheck(item, "uid") {
		b.UID = item["uid"].(string)
	}

	if nullcheck(item, "x509") {
		jsonData, _ := json.Marshal(item["x509"])

		err := json.Unmarshal(jsonData, &b.X509)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *FileX509) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "alternative_names") {
		if reflect.TypeOf(item["alternative_names"]).Kind() == reflect.String {
			b.AlternativeNames = append(b.AlternativeNames, item["alternative_names"].(string))
		} else {
			for _, v := range item["alternative_names"].([]interface{}) {
				b.AlternativeNames = append(b.AlternativeNames, v.(string))
			}
		}
	}

	if nullcheck(item, "issuer") {
		jsonData, _ := json.Marshal(item["issuer"])

		err := json.Unmarshal(jsonData, &b.Issuer)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "not_after") {
		b.NotAfter, _ = time.Parse(time.RFC3339, item["not_after"].(string))
	}

	if nullcheck(item, "not_before") {
		b.NotBefore, _ = time.Parse(time.RFC3339, item["not_before"].(string))
	}

	if nullcheck(item, "public_key_algorithm") {
		b.PublicKeyAlgorithm = item["public_key_algorithm"].(string)
	}

	if nullcheck(item, "public_key_curve") {
		b.PublicKeyCurve = item["public_key_curve"].(string)
	}

	if nullcheck(item, "public_key_exponent") {
		val, err := safeCastInt(item["public_key_exponent"])
		if err != nil {
			return err
		}

		b.PublicKeyExponent = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "public_key_size") {
		val, err := safeCastInt(item["public_key_size"])
		if err != nil {
			return err
		}

		b.PublicKeySize = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "serial_number") {
		b.SerialNumber = item["serial_number"].(string)
	}

	if nullcheck(item, "signature_algorithm") {
		b.SignatureAlgorithm = item["signature_algorithm"].(string)
	}

	if nullcheck(item, "subject") {
		jsonData, _ := json.Marshal(item["subject"])

		err := json.Unmarshal(jsonData, &b.Subject)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "version_number") {
		b.VersionNumber = item["version_number"].(string)
	}

	return nil
}

func (b *FileX509Issuer) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "common_name") {
		if reflect.TypeOf(item["common_name"]).Kind() == reflect.String {
			b.CommonName = append(b.CommonName, item["common_name"].(string))
		} else {
			for _, v := range item["category"].([]interface{}) {
				b.CommonName = append(b.CommonName, v.(string))
			}
		}
	}

	if nullcheck(item, "country") {
		if reflect.TypeOf(item["country"]).Kind() == reflect.String {
			b.Country = append(b.Country, item["country"].(string))
		} else {
			for _, v := range item["country"].([]interface{}) {
				b.Country = append(b.Country, v.(string))
			}
		}
	}

	if nullcheck(item, "distinguished_name") {
		b.DistinguishedName = item["distinguished_name"].(string)
	}

	if nullcheck(item, "locality") {
		if reflect.TypeOf(item["locality"]).Kind() == reflect.String {
			b.Locality = append(b.Locality, item["locality"].(string))
		} else {
			for _, v := range item["locality"].([]interface{}) {
				b.Locality = append(b.Locality, v.(string))
			}
		}
	}

	if nullcheck(item, "organization") {
		if reflect.TypeOf(item["organization"]).Kind() == reflect.String {
			b.Organization = append(b.Organization, item["organization"].(string))
		} else {
			for _, v := range item["organization"].([]interface{}) {
				b.Organization = append(b.Organization, v.(string))
			}
		}
	}

	if nullcheck(item, "organizational_unit") {
		if reflect.TypeOf(item["organizational_unit"]).Kind() == reflect.String {
			b.OrganizationalUnit = append(b.OrganizationalUnit, item["organizational_unit"].(string))
		} else {
			for _, v := range item["organizational_unit"].([]interface{}) {
				b.OrganizationalUnit = append(b.OrganizationalUnit, v.(string))
			}
		}
	}

	if nullcheck(item, "state_or_province") {
		if reflect.TypeOf(item["state_or_province"]).Kind() == reflect.String {
			b.StateOrProvince = append(b.StateOrProvince, item["state_or_province"].(string))
		} else {
			for _, v := range item["state_or_province"].([]interface{}) {
				b.StateOrProvince = append(b.StateOrProvince, v.(string))
			}
		}
	}

	return nil
}

func (b *FileX509Subject) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "common_name") {
		if reflect.TypeOf(item["common_name"]).Kind() == reflect.String {
			b.CommonName = append(b.CommonName, item["common_name"].(string))
		} else {
			for _, v := range item["category"].([]interface{}) {
				b.CommonName = append(b.CommonName, v.(string))
			}
		}
	}

	if nullcheck(item, "country") {
		if reflect.TypeOf(item["country"]).Kind() == reflect.String {
			b.Country = append(b.Country, item["country"].(string))
		} else {
			for _, v := range item["country"].([]interface{}) {
				b.Country = append(b.Country, v.(string))
			}
		}
	}

	if nullcheck(item, "distinguished_name") {
		b.DistinguishedName = item["distinguished_name"].(string)
	}

	if nullcheck(item, "locality") {
		if reflect.TypeOf(item["locality"]).Kind() == reflect.String {
			b.Locality = append(b.Locality, item["locality"].(string))
		} else {
			for _, v := range item["locality"].([]interface{}) {
				b.Locality = append(b.Locality, v.(string))
			}
		}
	}

	if nullcheck(item, "organization") {
		if reflect.TypeOf(item["organization"]).Kind() == reflect.String {
			b.Organization = append(b.Organization, item["organization"].(string))
		} else {
			for _, v := range item["organization"].([]interface{}) {
				b.Organization = append(b.Organization, v.(string))
			}
		}
	}

	if nullcheck(item, "organizational_unit") {
		if reflect.TypeOf(item["organizational_unit"]).Kind() == reflect.String {
			b.OrganizationalUnit = append(b.OrganizationalUnit, item["organizational_unit"].(string))
		} else {
			for _, v := range item["organizational_unit"].([]interface{}) {
				b.OrganizationalUnit = append(b.OrganizationalUnit, v.(string))
			}
		}
	}

	if nullcheck(item, "state_or_province") {
		if reflect.TypeOf(item["state_or_province"]).Kind() == reflect.String {
			b.StateOrProvince = append(b.StateOrProvince, item["state_or_province"].(string))
		} else {
			for _, v := range item["state_or_province"].([]interface{}) {
				b.StateOrProvince = append(b.StateOrProvince, v.(string))
			}
		}
	}

	return nil
}

func (b *Host) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "architecture") {
		b.Architecture = item["architecture"].(string)
	}

	if nullcheck(item, "domain") {
		b.Domain = item["domain"].(string)
	}

	if nullcheck(item, "geo") {
		jsonData, _ := json.Marshal(item["geo"])

		err := json.Unmarshal(jsonData, &b.Geo)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "hostname") {
		b.Hostname = item["hostname"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "ip") {
		if reflect.TypeOf(item["ip"]).Kind() == reflect.String {
			b.IP = append(b.IP, item["ip"].(string))
		} else {
			for _, v := range item["ip"].([]interface{}) {
				b.IP = append(b.IP, v.(string))
			}
		}
	}

	if nullcheck(item, "mac") {
		if reflect.TypeOf(item["mac"]).Kind() == reflect.String {
			b.MAC = append(b.MAC, item["mac"].(string))
		} else {
			for _, v := range item["mac"].([]interface{}) {
				b.MAC = append(b.MAC, v.(string))
			}
		}
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "os") {
		jsonData, _ := json.Marshal(item["os"])

		err := json.Unmarshal(jsonData, &b.OS)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "type") {
		b.Type = item["type"].(string)
	}

	if nullcheck(item, "uptime") {
		val, err := safeCastInt(item["uptime"])
		if err != nil {
			return err
		}

		b.Uptime = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "user") {
		jsonData, _ := json.Marshal(item["user"])

		err := json.Unmarshal(jsonData, &b.User)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *HostUser) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "domain") {
		b.Domain = item["domain"].(string)
	}

	if nullcheck(item, "email") {
		b.Email = item["email"].(string)
	}

	if nullcheck(item, "full_name") {
		b.FullName = item["full_name"].(string)
	}

	if nullcheck(item, "group") {
		jsonData, _ := json.Marshal(item["group"])

		err := json.Unmarshal(jsonData, &b.Group)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "hash") {
		b.Hash = item["hash"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "roles") {
		if reflect.TypeOf(item["roles"]).Kind() == reflect.String {
			b.Roles = append(b.Roles, item["roles"].(string))
		} else {
			for _, v := range item["roles"].([]interface{}) {
				b.Roles = append(b.Roles, v.(string))
			}
		}
	}

	return nil
}

func (b *Observer) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "egress") {
		jsonData, _ := json.Marshal(item["egress"])

		err := json.Unmarshal(jsonData, &b.Egress)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "geo") {
		jsonData, _ := json.Marshal(item["geo"])

		err := json.Unmarshal(jsonData, &b.Geo)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "hostname") {
		b.Hostname = item["hostname"].(string)
	}

	if nullcheck(item, "ingress") {
		jsonData, _ := json.Marshal(item["ingress"])

		err := json.Unmarshal(jsonData, &b.Ingress)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "ip") {
		if reflect.TypeOf(item["ip"]).Kind() == reflect.String {
			b.IP = append(b.IP, item["ip"].(string))
		} else {
			for _, v := range item["ip"].([]interface{}) {
				b.IP = append(b.IP, v.(string))
			}
		}
	}

	if nullcheck(item, "mac") {
		if reflect.TypeOf(item["mac"]).Kind() == reflect.String {
			b.MAC = append(b.MAC, item["mac"].(string))
		} else {
			for _, v := range item["mac"].([]interface{}) {
				b.MAC = append(b.MAC, v.(string))
			}
		}
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "os") {
		jsonData, _ := json.Marshal(item["os"])

		err := json.Unmarshal(jsonData, &b.OS)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "product") {
		b.Product = item["product"].(string)
	}

	if nullcheck(item, "serial_number") {
		b.SerialNumber = item["serial_number"].(string)
	}

	if nullcheck(item, "type") {
		b.Type = item["type"].(string)
	}

	if nullcheck(item, "vendor") {
		b.Vendor = item["vendor"].(string)
	}

	if nullcheck(item, "version") {
		b.Version = item["version"].(string)
	}

	return nil
}

func (b *RegistryData) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "bytes") {
		b.Bytes = item["bytes"].(string)
	}

	if nullcheck(item, "strings") {
		if reflect.TypeOf(item["strings"]).Kind() == reflect.String {
			b.Strings = append(b.Strings, item["strings"].(string))
		} else {
			for _, v := range item["strings"].([]interface{}) {
				b.Strings = append(b.Strings, v.(string))
			}
		}
	}

	if nullcheck(item, "bytes") {
		b.Type = item["bytes"].(string)
	}

	return nil
}

func (b *Related) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "hash") {
		if reflect.TypeOf(item["hash"]).Kind() == reflect.String {
			b.Hash = append(b.Hash, item["hash"].(string))
		} else {
			for _, v := range item["hash"].([]interface{}) {
				b.Hash = append(b.Hash, v.(string))
			}
		}
	}

	if nullcheck(item, "hosts") {
		if reflect.TypeOf(item["hosts"]).Kind() == reflect.String {
			b.Hosts = append(b.Hosts, item["hosts"].(string))
		} else {
			for _, v := range item["hosts"].([]interface{}) {
				b.Hosts = append(b.Hosts, v.(string))
			}
		}
	}

	if nullcheck(item, "ip") {
		if reflect.TypeOf(item["ip"]).Kind() == reflect.String {
			b.IP = append(b.IP, item["ip"].(string))
		} else {
			for _, v := range item["ip"].([]interface{}) {
				b.IP = append(b.IP, v.(string))
			}
		}
	}

	if nullcheck(item, "user") {
		if reflect.TypeOf(item["user"]).Kind() == reflect.String {
			b.User = append(b.User, item["user"].(string))
		} else {
			for _, v := range item["user"].([]interface{}) {
				b.User = append(b.User, v.(string))
			}
		}
	}

	return nil
}

func (b *Rule) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "author") {
		if reflect.TypeOf(item["author"]).Kind() == reflect.String {
			b.Author = append(b.Author, item["author"].(string))
		} else {
			for _, v := range item["author"].([]interface{}) {
				b.Author = append(b.Author, v.(string))
			}
		}
	}

	if nullcheck(item, "category") {
		b.Category = item["category"].(string)
	}

	if nullcheck(item, "description") {
		b.Description = item["description"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "license") {
		b.License = item["license"].(string)
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "reference") {
		b.Reference = item["reference"].(string)
	}

	if nullcheck(item, "ruleset") {
		b.Ruleset = item["ruleset"].(string)
	}

	if nullcheck(item, "uuid") {
		b.UUID = item["uuid"].(string)
	}

	if nullcheck(item, "version") {
		b.Version = item["version"].(string)
	}

	return nil
}

func (b *ServerUser) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if reflect.TypeOf(item["roles"]).Kind() == reflect.String {
		b.Roles = append(b.Roles, item["roles"].(string))
	}

	return nil
}

func (b *SourceUser) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "domain") {
		b.Domain = item["domain"].(string)
	}

	if nullcheck(item, "email") {
		b.Email = item["email"].(string)
	}

	if nullcheck(item, "full_name") {
		b.FullName = item["full_name"].(string)
	}

	if nullcheck(item, "group") {
		jsonData, _ := json.Marshal(item["group"])

		err := json.Unmarshal(jsonData, &b.Group)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "hash") {
		b.Hash = item["hash"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "roles") {
		if reflect.TypeOf(item["roles"]).Kind() == reflect.String {
			b.Roles = append(b.Roles, item["roles"].(string))
		} else {
			for _, v := range item["roles"].([]interface{}) {
				b.Roles = append(b.Roles, v.(string))
			}
		}
	}

	return nil
}

func (b *ThreatTactic) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "id") {
		if reflect.TypeOf(item["id"]).Kind() == reflect.String {
			b.ID = append(b.ID, item["id"].(string))
		} else {
			for _, v := range item["id"].([]interface{}) {
				b.ID = append(b.ID, v.(string))
			}
		}
	}

	if nullcheck(item, "name") {
		if reflect.TypeOf(item["name"]).Kind() == reflect.String {
			b.Name = append(b.Name, item["name"].(string))
		} else {
			for _, v := range item["name"].([]interface{}) {
				b.Name = append(b.Name, v.(string))
			}
		}
	}

	if nullcheck(item, "reference") {
		if reflect.TypeOf(item["reference"]).Kind() == reflect.String {
			b.Reference = append(b.Reference, item["reference"].(string))
		} else {
			for _, v := range item["reference"].([]interface{}) {
				b.Reference = append(b.Reference, v.(string))
			}
		}
	}

	return nil
}

func (b *TLSClient) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "certificate") {
		b.Certificate = item["certificate"].(string)
	}

	if nullcheck(item, "certificate_chain") {
		if reflect.TypeOf(item["certificate_chain"]).Kind() == reflect.String {
			b.CertificateChain = append(b.CertificateChain, item["certificate_chain"].(string))
		} else {
			for _, v := range item["certificate_chain"].([]interface{}) {
				b.CertificateChain = append(b.CertificateChain, v.(string))
			}
		}
	}

	if nullcheck(item, "hash") {
		jsonData, _ := json.Marshal(item["hash"])

		err := json.Unmarshal(jsonData, &b.Hash)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "issuer") {
		b.Issuer = item["issuer"].(string)
	}

	if nullcheck(item, "ja3") {
		b.JA3 = item["ja3"].(string)
	}

	if nullcheck(item, "not_after") {
		b.NotAfter, _ = time.Parse(time.RFC3339, item["not_after"].(string))
	}

	if nullcheck(item, "not_before") {
		b.NotBefore, _ = time.Parse(time.RFC3339, item["not_before"].(string))
	}

	if nullcheck(item, "server_name") {
		b.ServerName = item["server_name"].(string)
	}

	if nullcheck(item, "subject") {
		b.Subject = item["subject"].(string)
	}

	if nullcheck(item, "supported_ciphers") {
		if reflect.TypeOf(item["supported_ciphers"]).Kind() == reflect.String {
			b.SupportedCiphers = append(b.SupportedCiphers, item["supported_ciphers"].(string))
		} else {
			for _, v := range item["supported_ciphers"].([]interface{}) {
				b.SupportedCiphers = append(b.SupportedCiphers, v.(string))
			}
		}
	}

	if nullcheck(item, "x509") {
		jsonData, _ := json.Marshal(item["x509"])

		err := json.Unmarshal(jsonData, &b.X509)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *TLSClientX509) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "alternative_names") {
		if reflect.TypeOf(item["alternative_names"]).Kind() == reflect.String {
			b.AlternativeNames = append(b.AlternativeNames, item["alternative_names"].(string))
		} else {
			for _, v := range item["alternative_names"].([]interface{}) {
				b.AlternativeNames = append(b.AlternativeNames, v.(string))
			}
		}
	}

	if nullcheck(item, "issuer") {
		jsonData, _ := json.Marshal(item["issuer"])

		err := json.Unmarshal(jsonData, &b.Issuer)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "not_after") {
		b.NotAfter, _ = time.Parse(time.RFC3339, item["not_after"].(string))
	}

	if nullcheck(item, "not_before") {
		b.NotBefore, _ = time.Parse(time.RFC3339, item["not_before"].(string))
	}

	if nullcheck(item, "public_key_algorithm") {
		b.PublicKeyAlgorithm = item["public_key_algorithm"].(string)
	}

	if nullcheck(item, "public_key_curve") {
		b.PublicKeyCurve = item["public_key_curve"].(string)
	}

	if nullcheck(item, "public_key_exponent") {
		val, err := safeCastInt(item["public_key_exponent"])
		if err != nil {
			return err
		}

		b.PublicKeyExponent = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "public_key_size") {
		val, err := safeCastInt(item["public_key_size"])
		if err != nil {
			return err
		}

		b.PublicKeySize = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "serial_number") {
		b.SerialNumber = item["serial_number"].(string)
	}

	if nullcheck(item, "signature_algorithm") {
		b.SignatureAlgorithm = item["signature_algorithm"].(string)
	}

	if nullcheck(item, "subject") {
		jsonData, _ := json.Marshal(item["subject"])

		err := json.Unmarshal(jsonData, &b.Subject)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "version_number") {
		b.VersionNumber = item["version_number"].(string)
	}

	return nil
}

func (b *TLSClientX509Issuer) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "common_name") {
		if reflect.TypeOf(item["common_name"]).Kind() == reflect.String {
			b.CommonName = append(b.CommonName, item["common_name"].(string))
		} else {
			for _, v := range item["common_name"].([]interface{}) {
				b.CommonName = append(b.CommonName, v.(string))
			}
		}
	}

	if nullcheck(item, "country") {
		if reflect.TypeOf(item["country"]).Kind() == reflect.String {
			b.Country = append(b.Country, item["country"].(string))
		} else {
			for _, v := range item["country"].([]interface{}) {
				b.Country = append(b.Country, v.(string))
			}
		}
	}

	if nullcheck(item, "distinguished_name") {
		b.DistinguishedName = item["distinguished_name"].(string)
	}

	if nullcheck(item, "locality") {
		if reflect.TypeOf(item["locality"]).Kind() == reflect.String {
			b.Locality = append(b.Locality, item["locality"].(string))
		} else {
			for _, v := range item["locality"].([]interface{}) {
				b.Locality = append(b.Locality, v.(string))
			}
		}
	}

	if nullcheck(item, "organization") {
		if reflect.TypeOf(item["organization"]).Kind() == reflect.String {
			b.Organization = append(b.Organization, item["organization"].(string))
		} else {
			for _, v := range item["organization"].([]interface{}) {
				b.Organization = append(b.Organization, v.(string))
			}
		}
	}

	if nullcheck(item, "organizational_unit") {
		if reflect.TypeOf(item["organizational_unit"]).Kind() == reflect.String {
			b.OrganizationalUnit = append(b.OrganizationalUnit, item["organizational_unit"].(string))
		} else {
			for _, v := range item["organizational_unit"].([]interface{}) {
				b.OrganizationalUnit = append(b.OrganizationalUnit, v.(string))
			}
		}
	}

	if nullcheck(item, "state_or_province") {
		if reflect.TypeOf(item["state_or_province"]).Kind() == reflect.String {
			b.StateOrProvince = append(b.StateOrProvince, item["state_or_province"].(string))
		} else {
			for _, v := range item["state_or_province"].([]interface{}) {
				b.StateOrProvince = append(b.StateOrProvince, v.(string))
			}
		}
	}

	return nil
}

func (b *TLSClientX509Subject) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "common_name") {
		if reflect.TypeOf(item["common_name"]).Kind() == reflect.String {
			b.CommonName = append(b.CommonName, item["common_name"].(string))
		} else {
			for _, v := range item["common_name"].([]interface{}) {
				b.CommonName = append(b.CommonName, v.(string))
			}
		}
	}

	if nullcheck(item, "country") {
		if reflect.TypeOf(item["country"]).Kind() == reflect.String {
			b.Country = append(b.Country, item["country"].(string))
		} else {
			for _, v := range item["country"].([]interface{}) {
				b.Country = append(b.Country, v.(string))
			}
		}
	}

	if nullcheck(item, "distinguished_name") {
		b.DistinguishedName = item["distinguished_name"].(string)
	}

	if nullcheck(item, "locality") {
		if reflect.TypeOf(item["locality"]).Kind() == reflect.String {
			b.Locality = append(b.Locality, item["locality"].(string))
		} else {
			for _, v := range item["locality"].([]interface{}) {
				b.Locality = append(b.Locality, v.(string))
			}
		}
	}

	if nullcheck(item, "organization") {
		if reflect.TypeOf(item["organization"]).Kind() == reflect.String {
			b.Organization = append(b.Organization, item["organization"].(string))
		} else {
			for _, v := range item["organization"].([]interface{}) {
				b.Organization = append(b.Organization, v.(string))
			}
		}
	}

	if nullcheck(item, "organizational_unit") {
		if reflect.TypeOf(item["organizational_unit"]).Kind() == reflect.String {
			b.OrganizationalUnit = append(b.OrganizationalUnit, item["organizational_unit"].(string))
		} else {
			for _, v := range item["organizational_unit"].([]interface{}) {
				b.OrganizationalUnit = append(b.OrganizationalUnit, v.(string))
			}
		}
	}

	if nullcheck(item, "state_or_province") {
		if reflect.TypeOf(item["state_or_province"]).Kind() == reflect.String {
			b.StateOrProvince = append(b.StateOrProvince, item["state_or_province"].(string))
		} else {
			for _, v := range item["state_or_province"].([]interface{}) {
				b.StateOrProvince = append(b.StateOrProvince, v.(string))
			}
		}
	}

	return nil
}

func (b *TLSServer) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "certificate") {
		b.Certificate = item["certificate"].(string)
	}

	if nullcheck(item, "certificate_chain") {
		if reflect.TypeOf(item["certificate_chain"]).Kind() == reflect.String {
			b.CertificateChain = append(b.CertificateChain, item["certificate_chain"].(string))
		} else {
			for _, v := range item["certificate_chain"].([]interface{}) {
				b.CertificateChain = append(b.CertificateChain, v.(string))
			}
		}
	}

	if nullcheck(item, "hash") {
		jsonData, _ := json.Marshal(item["hash"])

		err := json.Unmarshal(jsonData, &b.Hash)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "issuer") {
		b.Issuer = item["issuer"].(string)
	}

	if nullcheck(item, "ja3s") {
		b.JA3S = item["ja3s"].(string)
	}

	if nullcheck(item, "not_after") {
		b.NotAfter, _ = time.Parse(time.RFC3339, item["not_after"].(string))
	}

	if nullcheck(item, "not_before") {
		b.NotBefore, _ = time.Parse(time.RFC3339, item["not_before"].(string))
	}

	if nullcheck(item, "subject") {
		b.Subject = item["subject"].(string)
	}

	if nullcheck(item, "x509") {
		jsonData, _ := json.Marshal(item["x509"])

		err := json.Unmarshal(jsonData, &b.X509)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *TLSServerX509) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "alternative_names") {
		if reflect.TypeOf(item["alternative_names"]).Kind() == reflect.String {
			b.AlternativeNames = append(b.AlternativeNames, item["alternative_names"].(string))
		} else {
			for _, v := range item["alternative_names"].([]interface{}) {
				b.AlternativeNames = append(b.AlternativeNames, v.(string))
			}
		}
	}

	if nullcheck(item, "issuer") {
		jsonData, _ := json.Marshal(item["issuer"])

		err := json.Unmarshal(jsonData, &b.Issuer)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "not_after") {
		b.NotAfter, _ = time.Parse(time.RFC3339, item["not_after"].(string))
	}

	if nullcheck(item, "not_before") {
		b.NotBefore, _ = time.Parse(time.RFC3339, item["not_before"].(string))
	}

	if nullcheck(item, "public_key_algorithm") {
		b.PublicKeyAlgorithm = item["public_key_algorithm"].(string)
	}

	if nullcheck(item, "public_key_curve") {
		b.PublicKeyCurve = item["public_key_curve"].(string)
	}

	if nullcheck(item, "public_key_exponent") {
		val, err := safeCastInt(item["public_key_exponent"])
		if err != nil {
			return err
		}

		b.PublicKeyExponent = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "public_key_size") {
		val, err := safeCastInt(item["public_key_size"])
		if err != nil {
			return err
		}

		b.PublicKeySize = json.Number(fmt.Sprintf("%d", val))
	}

	if nullcheck(item, "serial_number") {
		b.SerialNumber = item["serial_number"].(string)
	}

	if nullcheck(item, "signature_algorithm") {
		b.SignatureAlgorithm = item["signature_algorithm"].(string)
	}

	if nullcheck(item, "subject") {
		jsonData, _ := json.Marshal(item["subject"])

		err := json.Unmarshal(jsonData, &b.Subject)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "version_number") {
		b.VersionNumber = item["version_number"].(string)
	}

	return nil
}

func (b *TLSServerX509Issuer) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "common_name") {
		if reflect.TypeOf(item["common_name"]).Kind() == reflect.String {
			b.CommonName = append(b.CommonName, item["common_name"].(string))
		} else {
			for _, v := range item["category"].([]interface{}) {
				b.CommonName = append(b.CommonName, v.(string))
			}
		}
	}

	if nullcheck(item, "country") {
		if reflect.TypeOf(item["country"]).Kind() == reflect.String {
			b.Country = append(b.Country, item["country"].(string))
		} else {
			for _, v := range item["country"].([]interface{}) {
				b.Country = append(b.Country, v.(string))
			}
		}
	}

	if nullcheck(item, "distinguished_name") {
		b.DistinguishedName = item["distinguished_name"].(string)
	}

	if nullcheck(item, "locality") {
		if reflect.TypeOf(item["locality"]).Kind() == reflect.String {
			b.Locality = append(b.Locality, item["locality"].(string))
		} else {
			for _, v := range item["locality"].([]interface{}) {
				b.Locality = append(b.Locality, v.(string))
			}
		}
	}

	if nullcheck(item, "organization") {
		if reflect.TypeOf(item["organization"]).Kind() == reflect.String {
			b.Organization = append(b.Organization, item["organization"].(string))
		} else {
			for _, v := range item["organization"].([]interface{}) {
				b.Organization = append(b.Organization, v.(string))
			}
		}
	}

	if nullcheck(item, "organizational_unit") {
		if reflect.TypeOf(item["organizational_unit"]).Kind() == reflect.String {
			b.OrganizationalUnit = append(b.OrganizationalUnit, item["organizational_unit"].(string))
		} else {
			for _, v := range item["organizational_unit"].([]interface{}) {
				b.OrganizationalUnit = append(b.OrganizationalUnit, v.(string))
			}
		}
	}

	if nullcheck(item, "state_or_province") {
		if reflect.TypeOf(item["state_or_province"]).Kind() == reflect.String {
			b.StateOrProvince = append(b.StateOrProvince, item["state_or_province"].(string))
		} else {
			for _, v := range item["state_or_province"].([]interface{}) {
				b.StateOrProvince = append(b.StateOrProvince, v.(string))
			}
		}
	}

	return nil
}

func (b *TLSServerX509Subject) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "common_name") {
		if reflect.TypeOf(item["common_name"]).Kind() == reflect.String {
			b.CommonName = append(b.CommonName, item["common_name"].(string))
		} else {
			for _, v := range item["category"].([]interface{}) {
				b.CommonName = append(b.CommonName, v.(string))
			}
		}
	}

	if nullcheck(item, "country") {
		if reflect.TypeOf(item["country"]).Kind() == reflect.String {
			b.Country = append(b.Country, item["country"].(string))
		} else {
			for _, v := range item["country"].([]interface{}) {
				b.Country = append(b.Country, v.(string))
			}
		}
	}

	if nullcheck(item, "distinguished_name") {
		b.DistinguishedName = item["distinguished_name"].(string)
	}

	if nullcheck(item, "locality") {
		if reflect.TypeOf(item["locality"]).Kind() == reflect.String {
			b.Locality = append(b.Locality, item["locality"].(string))
		} else {
			for _, v := range item["locality"].([]interface{}) {
				b.Locality = append(b.Locality, v.(string))
			}
		}
	}

	if nullcheck(item, "organization") {
		if reflect.TypeOf(item["organization"]).Kind() == reflect.String {
			b.Organization = append(b.Organization, item["organization"].(string))
		} else {
			for _, v := range item["organization"].([]interface{}) {
				b.Organization = append(b.Organization, v.(string))
			}
		}
	}

	if nullcheck(item, "organizational_unit") {
		if reflect.TypeOf(item["organizational_unit"]).Kind() == reflect.String {
			b.OrganizationalUnit = append(b.OrganizationalUnit, item["organizational_unit"].(string))
		} else {
			for _, v := range item["organizational_unit"].([]interface{}) {
				b.OrganizationalUnit = append(b.OrganizationalUnit, v.(string))
			}
		}
	}

	if nullcheck(item, "state_or_province") {
		if reflect.TypeOf(item["state_or_province"]).Kind() == reflect.String {
			b.StateOrProvince = append(b.StateOrProvince, item["state_or_province"].(string))
		} else {
			for _, v := range item["state_or_province"].([]interface{}) {
				b.StateOrProvince = append(b.StateOrProvince, v.(string))
			}
		}
	}

	return nil
}

func (b *User) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "domain") {
		b.Domain = item["domain"].(string)
	}

	if nullcheck(item, "email") {
		b.Email = item["email"].(string)
	}

	if nullcheck(item, "full_name") {
		b.FullName = item["full_name"].(string)
	}

	if nullcheck(item, "group") {
		jsonData, _ := json.Marshal(item["group"])

		err := json.Unmarshal(jsonData, &b.Group)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "hash") {
		b.Hash = item["hash"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "name") {
		b.Name = item["name"].(string)
	}

	if nullcheck(item, "roles") {
		if reflect.TypeOf(item["roles"]).Kind() == reflect.String {
			b.Roles = append(b.Roles, item["roles"].(string))
		} else {
			for _, v := range item["roles"].([]interface{}) {
				b.Roles = append(b.Roles, v.(string))
			}
		}
	}

	return nil
}

func (b *Vulnerability) UnmarshalJSON(data []byte) error {
	var item map[string]interface{}

	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	if nullcheck(item, "category") {
		if reflect.TypeOf(item["category"]).Kind() == reflect.String {
			b.Category = append(b.Category, item["category"].(string))
		} else {
			for _, v := range item["category"].([]interface{}) {
				b.Category = append(b.Category, v.(string))
			}
		}
	}

	if nullcheck(item, "classification") {
		b.Classification = item["classification"].(string)
	}

	if nullcheck(item, "description") {
		b.Description = item["description"].(string)
	}

	if nullcheck(item, "enumeration") {
		b.Enumeration = item["enumeration"].(string)
	}

	if nullcheck(item, "id") {
		b.ID = item["id"].(string)
	}

	if nullcheck(item, "reference") {
		b.Reference = item["reference"].(string)
	}

	if nullcheck(item, "report_id") {
		b.ReportID = item["report_id"].(string)
	}

	if nullcheck(item, "scanner") {
		jsonData, _ := json.Marshal(item["scanner"])

		err := json.Unmarshal(jsonData, &b.Scanner)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "score") {
		jsonData, _ := json.Marshal(item["score"])

		err := json.Unmarshal(jsonData, &b.Score)
		if err != nil {
			return err
		}
	}

	if nullcheck(item, "severity") {
		b.Severity = item["severity"].(string)
	}

	return nil
}
