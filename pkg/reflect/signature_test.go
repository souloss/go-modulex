package reflect

import (
	"context"
	"reflect"
	"testing"
)

func TestSignatureSpec_Validate(t *testing.T) {
	tests := []struct {
		name    string
		spec    *SignatureSpec
		fn      interface{}
		wantErr bool
	}{
		{
			name:    "valid function with no type checkers",
			spec:    NewSignatureSpec(2, 1, nil, nil),
			fn:      func(a, b string) int { return 0 },
			wantErr: false,
		},
		{
			name: "valid function with type checkers",
			spec: NewSignatureSpec(2, 1,
				[]TypeChecker{IsString, IsInt},
				[]TypeChecker{IsError}),
			fn:      func(a string, b int) error { return nil },
			wantErr: false,
		},
		{
			name:    "invalid number of inputs",
			spec:    NewSignatureSpec(2, 1, nil, nil),
			fn:      func(a string) int { return 0 },
			wantErr: true,
		},
		{
			name:    "invalid number of outputs",
			spec:    NewSignatureSpec(1, 2, nil, nil),
			fn:      func(a string) int { return 0 },
			wantErr: true,
		},
		{
			name:    "nil function",
			spec:    NewSignatureSpec(1, 1, nil, nil),
			fn:      nil,
			wantErr: true,
		},
		{
			name:    "non-function input",
			spec:    NewSignatureSpec(1, 1, nil, nil),
			fn:      "not a function",
			wantErr: true,
		},
		{
			name: "invalid input type",
			spec: NewSignatureSpec(1, 1,
				[]TypeChecker{IsString},
				[]TypeChecker{IsError}),
			fn:      func(a int) error { return nil },
			wantErr: true,
		},
		{
			name: "invalid output type",
			spec: NewSignatureSpec(1, 1,
				[]TypeChecker{IsString},
				[]TypeChecker{IsError}),
			fn:      func(a string) int { return 0 },
			wantErr: true,
		},
		{
			name: "multiple valid type checkers",
			spec: NewSignatureSpec(1, 1,
				[]TypeChecker{IsString, IsInt},
				[]TypeChecker{IsError}),
			fn:      func(a int) error { return nil },
			wantErr: false,
		},
		{
			name:    "zero expected parameters",
			spec:    NewSignatureSpec(0, 0, nil, nil),
			fn:      func() {},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.spec.Validate(tt.fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignatureSpec.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommonTypeCheckers(t *testing.T) {
	tests := []struct {
		name    string
		checker TypeChecker
		type_   reflect.Type
		want    bool
	}{
		{
			name:    "IsString with string",
			checker: IsString,
			type_:   reflect.TypeOf(""),
			want:    true,
		},
		{
			name:    "IsString with int",
			checker: IsString,
			type_:   reflect.TypeOf(0),
			want:    false,
		},
		{
			name:    "IsInt with int",
			checker: IsInt,
			type_:   reflect.TypeOf(0),
			want:    true,
		},
		{
			name:    "IsInt with string",
			checker: IsInt,
			type_:   reflect.TypeOf(""),
			want:    false,
		},
		{
			name:    "IsError with error",
			checker: IsError,
			type_:   reflect.TypeOf((*error)(nil)).Elem(),
			want:    true,
		},
		{
			name:    "IsError with string",
			checker: IsError,
			type_:   reflect.TypeOf(""),
			want:    false,
		},
		{
			name:    "IsContext with context",
			checker: IsContext,
			type_:   reflect.TypeOf(context.Background()),
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.checker(tt.type_)
			if got != tt.want {
				t.Errorf("TypeChecker %v = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
