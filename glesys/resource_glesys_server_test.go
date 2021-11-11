package glesys

import (
	"testing"

	"github.com/glesys/glesys-go/v3"
)

func Test_getTemplate(t *testing.T) {
	srv := &glesys.ServerDetails{}
	for _, tt := range []struct {
		name         string
		tfTemplate   string
		readTemplate string
		readTags     []string
		want         string
	}{
		{
			name:         "KVM_instance",
			tfTemplate:   "ubuntu-20-04",
			readTemplate: "Ubuntu 20.04 LTS (Focal Fossa)",
			readTags:     []string{"ubuntu", "ubuntu-lts", "ubuntu-20-04"},
			want:         "ubuntu-20-04",
		},
		{
			name:         "VMware_instance",
			tfTemplate:   "Ubuntu 20.04 LTS 64-bit",
			readTemplate: "Ubuntu 20.04 LTS 64-bit",
			readTags:     []string{},
			want:         "Ubuntu 20.04 LTS 64-bit",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			srv.Template = tt.readTemplate
			srv.InitialTemplate.Name = tt.readTemplate
			srv.InitialTemplate.CurrentTags = tt.readTags
			if got := getTemplate(tt.tfTemplate, srv); got != tt.want {
				t.Errorf("got: %v, want %v", got, tt.want)
			}
		})
	}
}
