package feature

import (
	"testing"

	"github.com/mamaar/jsonchamp"
)

func TestPropKey(t *testing.T) {
	type args struct {
		propName string
		feat     *Feature
	}
	tests := []struct {
		name    string
		args    args
		want    Key
		wantErr bool
	}{
		{
			name: "prop does not exist on feature",
			args: args{
				propName: "non_existent_prop",
				feat:     New(Schema{}),
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "prop exists on feature",
			args: args{
				propName: "id",
				feat:     New(Schema{}, WithMap(jsonchamp.NewFromItems("id", "123"))),
			},
			wantErr: false,
			want:    "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := PropKey(tt.args.propName)(tt.args.feat)
			if (err != nil) != tt.wantErr {
				t.Errorf("PropKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id != tt.want {
				t.Fatalf("PropKey() = %v, want %v", id, tt.want)
			}
		})
	}
}

func TestCompositeKey(t *testing.T) {
	type args struct {
		items []KeyFunc
		feat  *Feature
	}
	tests := []struct {
		name    string
		args    args
		want    Key
		wantErr bool
	}{
		{
			name: "multiple literals",
			args: args{
				items: []KeyFunc{
					LiteralKey("order"),
					LiteralKey("1"),
				},
				feat: New(Schema{}),
			},
			want:    "order:1",
			wantErr: false,
		},
		{
			name: "literal and prop",
			args: args{
				items: []KeyFunc{
					LiteralKey("order"),
					PropKey("id"),
				},
				feat: New(Schema{}, WithMap(jsonchamp.NewFromItems("id", "123"))),
			},
			want:    "order:123",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CompositeKey(tt.args.items...)(tt.args.feat)
			if (err != nil) != tt.wantErr {
				t.Errorf("CompositeKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Fatalf("CompositeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLiteralKey(t *testing.T) {
	type args struct {
		id Key
	}
	tests := []struct {
		name    string
		args    args
		want    Key
		wantErr bool
	}{
		{
			name: "literal",
			args: args{
				id: "order:1",
			},
			want: "order:1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LiteralKey(tt.args.id)(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("LiteralKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Fatalf("LiteralKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateKeyedPayload(t *testing.T) {
	type args struct {
		feat     *Feature
		indexMap map[string]KeyFunc
	}
	tests := []struct {
		name string
		args args
		want *jsonchamp.Map
	}{
		{
			name: "empty index map",
			args: args{
				feat:     New(Schema{}),
				indexMap: map[string]KeyFunc{},
			},
			want: jsonchamp.NewFromItems("payload", jsonchamp.New()),
		},
		{
			name: "single index map",
			args: args{
				feat: New(Schema{}, WithMap(jsonchamp.NewFromItems("id", "123"))),
				indexMap: map[string]KeyFunc{
					"pk": PropKey("id"),
				},
			},
			want: jsonchamp.NewFromItems("pk", "123", "payload", jsonchamp.NewFromItems("id", "123")),
		},
		{
			name: "composite index map",
			args: args{
				feat: New(Schema{}, WithMap(jsonchamp.NewFromItems("id", "123"))),
				indexMap: map[string]KeyFunc{
					"pk": CompositeKey(LiteralKey("order"), PropKey("id")),
					"sk": LiteralKey("order"),
				},
			},
			want: jsonchamp.NewFromItems("pk", "order:123", "sk", "order", "payload", jsonchamp.NewFromItems("id", "123")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := CreateKeyedPayload(tt.args.feat, tt.args.indexMap)
			if err != nil {
				t.Fatalf("CreateKeyedPayload() error = %v", err)
			}
			diff := payload.Diff(tt.want)
			if len(diff.Keys()) > 0 {
				t.Fatalf("CreateKeyedPayload() = %v, want %v", payload, tt.want)
			}
		})
	}
}
