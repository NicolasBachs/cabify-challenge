package app

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_logrusLogger_Debug(t *testing.T) {
	type fields struct {
		env    string
		logger *logrus.Logger
	}
	type args struct {
		service string
		msg     string
		args    []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should print valid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%s'",
				args: []interface{}{
					"valid-string",
				},
			},
		},
		{
			name: "Should print invalid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%d'",
				args: []interface{}{
					"invalid-number",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &logrusLogger{
				env:    tt.fields.env,
				logger: tt.fields.logger,
			}
			log.Debug(tt.args.service, tt.args.msg, tt.args.args...)
		})
	}
}

func Test_logrusLogger_Info(t *testing.T) {
	type fields struct {
		env    string
		logger *logrus.Logger
	}
	type args struct {
		service string
		msg     string
		args    []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should print valid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%s'",
				args: []interface{}{
					"valid-string",
				},
			},
		},
		{
			name: "Should print invalid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%d'",
				args: []interface{}{
					"invalid-number",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &logrusLogger{
				env:    tt.fields.env,
				logger: tt.fields.logger,
			}
			log.Info(tt.args.service, tt.args.msg, tt.args.args...)
		})
	}
}

func Test_logrusLogger_Warn(t *testing.T) {
	type fields struct {
		env    string
		logger *logrus.Logger
	}
	type args struct {
		service string
		msg     string
		args    []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should print valid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%s'",
				args: []interface{}{
					"valid-string",
				},
			},
		},
		{
			name: "Should print invalid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%d'",
				args: []interface{}{
					"invalid-number",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &logrusLogger{
				env:    tt.fields.env,
				logger: tt.fields.logger,
			}
			log.Warn(tt.args.service, tt.args.msg, tt.args.args...)
		})
	}
}

func Test_logrusLogger_Error(t *testing.T) {
	type fields struct {
		env    string
		logger *logrus.Logger
	}
	type args struct {
		service string
		msg     string
		args    []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Should print valid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%s'",
				args: []interface{}{
					"valid-string",
				},
			},
		},
		{
			name: "Should print invalid msg format",
			fields: fields{
				env:    "production",
				logger: logrus.New(),
			},
			args: args{
				service: "some_service",
				msg:     "msg with this string: '%d'",
				args: []interface{}{
					"invalid-number",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := &logrusLogger{
				env:    tt.fields.env,
				logger: tt.fields.logger,
			}
			log.Error(tt.args.service, tt.args.msg, tt.args.args...)
		})
	}
}
