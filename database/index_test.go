package database

import (
	"testing"

	"github.com/mamaar/features/champ/maps"
	"github.com/mamaar/features/feature"
)

func TestProp(t *testing.T) {
	type args struct {
		propName string
		feat     *feature.Feature
	}
	tests := []struct {
		name    string
		args    args
		want    Identifier
		wantErr bool
	}{
		{
			name: "prop does not exist on feature",
			args: args{
				propName: "non_existent_prop",
				feat:     feature.New(feature.Schema{}),
			},
			wantErr: true,
			want:    "",
		},
		{
			name: "prop exists on feature",
			args: args{
				propName: "id",
				feat:     feature.New(feature.Schema{}, feature.WithMap(maps.NewFromItems("id", "123"))),
			},
			wantErr: false,
			want:    "123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := Prop(tt.args.propName)(tt.args.feat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Prop() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if id != tt.want {
				t.Fatalf("Prop() = %v, want %v", id, tt.want)
			}
		})
	}
}

func TestComp(t *testing.T) {
	type args struct {
		items []IdentifierBuilder
		feat  *feature.Feature
	}
	tests := []struct {
		name    string
		args    args
		want    Identifier
		wantErr bool
	}{
		{
			name: "multiple literals",
			args: args{
				items: []IdentifierBuilder{
					Literal("order"),
					Literal("1"),
				},
				feat: feature.New(feature.Schema{}),
			},
			want:    "order:1",
			wantErr: false,
		},
		{
			name: "literal and prop",
			args: args{
				items: []IdentifierBuilder{
					Literal("order"),
					Prop("id"),
				},
				feat: feature.New(feature.Schema{}, feature.WithMap(maps.NewFromItems("id", "123"))),
			},
			want:    "order:123",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Comp(tt.args.items...)(tt.args.feat)
			if (err != nil) != tt.wantErr {
				t.Errorf("Literal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Fatalf("Literal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLiteral(t *testing.T) {
	type args struct {
		id Identifier
	}
	tests := []struct {
		name    string
		args    args
		want    Identifier
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
			got, err := Literal(tt.args.id)(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("Literal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Fatalf("Literal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateIndexedDatabasePayload(t *testing.T) {
	type args struct {
		feat     *feature.Feature
		indexMap map[string]IdentifierBuilder
	}
	tests := []struct {
		name string
		args args
		want *maps.Map
	}{
		{
			name: "empty index map",
			args: args{
				feat:     feature.New(feature.Schema{}),
				indexMap: map[string]IdentifierBuilder{},
			},
			want: maps.NewFromItems("payload", maps.New()),
		},
		{
			name: "single index map",
			args: args{
				feat: feature.New(feature.Schema{}, feature.WithMap(maps.NewFromItems("id", "123"))),
				indexMap: map[string]IdentifierBuilder{
					"pk": Prop("id"),
				},
			},
			want: maps.NewFromItems("pk", "123", "payload", maps.NewFromItems("id", "123")),
		},
		{
			name: "composite index map",
			args: args{
				feat: feature.New(feature.Schema{}, feature.WithMap(maps.NewFromItems("id", "123"))),
				indexMap: map[string]IdentifierBuilder{
					"pk": Comp(Literal("order"), Prop("id")),
					"sk": Literal("order"),
				},
			},
			want: maps.NewFromItems("pk", "order:123", "sk", "order", "payload", maps.NewFromItems("id", "123")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := CreateIndexedDatabasePayload(tt.args.feat, tt.args.indexMap)
			if err != nil {
				t.Fatalf("CreateIndexedDatabasePayload() error = %v", err)
			}
			diff, err := payload.Diff(tt.want)
			if err != nil {
				t.Fatalf("CreateIndexedDatabasePayload() error = %v", err)
			}
			if len(diff.Keys()) > 0 {
				t.Fatalf("CreateIndexedDatabasePayload() = %v, want %v", payload, tt.want)
			}
		})
	}
}
