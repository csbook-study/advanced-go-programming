// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: validator.proto

package _

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_Message_ImportantString = regexp.MustCompile(`^[a-z]{2,5}$`)

func (this *Message) Validate() error {
	if !_regex_Message_ImportantString.MatchString(this.ImportantString) {
		return github_com_mwitkow_go_proto_validators.FieldError("ImportantString", fmt.Errorf(`value '%v' must be a string conforming to regex "^[a-z]{2,5}$"`, this.ImportantString))
	}
	if !(this.Age > 0) {
		return github_com_mwitkow_go_proto_validators.FieldError("Age", fmt.Errorf(`value '%v' must be greater than '0'`, this.Age))
	}
	if !(this.Age < 100) {
		return github_com_mwitkow_go_proto_validators.FieldError("Age", fmt.Errorf(`value '%v' must be less than '100'`, this.Age))
	}
	return nil
}
