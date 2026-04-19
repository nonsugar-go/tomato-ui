package checkpoint

import (
	"testing"

	"github.com/nonsugar-go/tomato-ui/internal/utmconv/model"
)

func TestConvertService(t *testing.T) {
	tests := []struct {
		name    string
		service model.Service
		want    string
		wantErr bool
	}{
		{
			name: "TCP service",
			service: model.Service{
				Type:        model.ServiceTypeTCP,
				Name:        "HTTP-Service",
				Ports:       "80",
				Description: "HTTP service",
				Tags:        []string{"web", "tcp"},
			},
			want:    `add service-tcp name "HTTP-Service" port 80 comments "HTTP service" tags.1 "web" tags.2 "tcp"`,
			wantErr: false,
		},
		{
			name: "UDP service",
			service: model.Service{
				Type:        model.ServiceTypeUDP,
				Name:        "DNS-Service",
				Ports:       "53",
				Description: "DNS service",
				Tags:        []string{"dns", "udp"},
			},
			want:    `add service-udp name "DNS-Service" port 53 comments "DNS service" tags.1 "dns" tags.2 "udp"`,
			wantErr: false,
		},
		{
			name: "Unsupported service type",
			service: model.Service{
				Type:        model.ServiceTypeUnknown,
				Name:        "Unknown Service",
				Ports:       "1234",
				Description: "Unsupported service type",
				Tags:        []string{"unknown"},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertService(tt.service)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertServiceGroup(t *testing.T) {
	tests := []struct {
		name    string
		group   model.ServiceGroup
		want    string
		wantErr bool
	}{
		{
			name: "Service group with members",
			group: model.ServiceGroup{
				Name:        "Web-Services",
				Members:     []string{"HTTP-Service", "HTTPS-Service"},
				Description: "Group of web services",
				Tags:        []string{"web", "group"},
			},
			want:    `add service-group name "Web-Services" members.1 "HTTP-Service" members.2 "HTTPS-Service" comments "Group of web services" tags.1 "web" tags.2 "group"`,
			wantErr: false,
		},
		{
			name: "Service group with no members",
			group: model.ServiceGroup{
				Name:        "Empty-Group",
				Members:     []string{},
				Description: "This group has no members",
				Tags:        []string{"empty"},
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertServiceGroup(tt.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertServiceGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertServiceGroup() = %v, want %v", got, tt.want)
			}
		})
	}
}
