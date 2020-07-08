package group

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

type fakeConfig struct {
	url      string
	admin    string
	password string
	timeout  int
}

var fc = fakeConfig{
	url:      "http://harbor.domain.com",
	admin:    "admin",
	password: "pwd",
	timeout:  10,
}

func TestNewHarbor(t *testing.T) {
	type args struct {
		url      string
		admin    string
		password string
	}
	tests := []struct {
		name string
		args args
		want Harbor
	}{
		{
			name: "case1",
			args: args{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
			},
			want: NewHarbor(fc.url, fc.admin, fc.password),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHarbor(tt.args.url, tt.args.admin, tt.args.password); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHarbor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_harbor_SystemInfo(t *testing.T) {
	type fields struct {
		url           string
		admin         string
		password      string
		timeout       int
		cookie        []*http.Cookie
		cookieTimeout time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:           tt.fields.url,
				admin:         tt.fields.admin,
				password:      tt.fields.password,
				timeout:       tt.fields.timeout,
				cookie:        tt.fields.cookie,
				cookieTimeout: tt.fields.cookieTimeout,
			}
			h.SystemInfo()
		})
	}
}

func Test_harbor_Login(t *testing.T) {
	type fields struct {
		url           string
		admin         string
		password      string
		timeout       int
		cookie        []*http.Cookie
		cookieTimeout time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:           tt.fields.url,
				admin:         tt.fields.admin,
				password:      tt.fields.password,
				timeout:       tt.fields.timeout,
				cookie:        tt.fields.cookie,
				cookieTimeout: tt.fields.cookieTimeout,
			}
			if err := h.Login(); (err != nil) != tt.wantErr {
				t.Errorf("harbor.Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_harbor_Projects(t *testing.T) {
	type fields struct {
		url           string
		admin         string
		password      string
		timeout       int
		cookie        []*http.Cookie
		cookieTimeout time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantRes []Project
		wantErr bool
	}{
		{
			name: "case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:           tt.fields.url,
				admin:         tt.fields.admin,
				password:      tt.fields.password,
				timeout:       tt.fields.timeout,
				cookie:        tt.fields.cookie,
				cookieTimeout: tt.fields.cookieTimeout,
			}
			_, err := h.Projects()
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Projects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.Projects() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}

func Test_harbor_Http(t *testing.T) {
	type fields struct {
		url           string
		admin         string
		password      string
		timeout       int
		cookie        []*http.Cookie
		cookieTimeout time.Time
	}
	type args struct {
		method string
		url    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes *http.Response
		wantErr bool
	}{
		{
			name: "http_case1",
			fields: fields{
				url:      fc.url,
				admin:    fc.admin,
				password: fc.password,
				timeout:  fc.timeout,
			},
			args: args{
				method: "GET",
				url:    fc.url,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:           tt.fields.url,
				admin:         tt.fields.admin,
				password:      tt.fields.password,
				timeout:       tt.fields.timeout,
				cookie:        tt.fields.cookie,
				cookieTimeout: tt.fields.cookieTimeout,
			}
			gotRes, err := h.Http(tt.args.method, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Http() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			_ = gotRes
			//if !reflect.DeepEqual(gotRes, tt.wantRes) {
			//	t.Errorf("harbor.Http() = %v, want %v", gotRes, tt.wantRes)
			//}
		})
	}
}

func Test_harbor_Repositories(t *testing.T) {
	type fields struct {
		url           string
		admin         string
		password      string
		timeout       int
		cookie        []*http.Cookie
		cookieTimeout time.Time
	}
	type args struct {
		projectId int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes []RepoRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &harbor{
				url:           tt.fields.url,
				admin:         tt.fields.admin,
				password:      tt.fields.password,
				timeout:       tt.fields.timeout,
				cookie:        tt.fields.cookie,
				cookieTimeout: tt.fields.cookieTimeout,
			}
			gotRes, err := h.Repositories(tt.args.projectId)
			if (err != nil) != tt.wantErr {
				t.Errorf("harbor.Repositories() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("harbor.Repositories() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
