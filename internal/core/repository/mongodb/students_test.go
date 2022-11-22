package mongodb

import (
	"context"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

type collectionMock struct {
	data interface{}
}

func (m *collectionMock) InsertOne(ctx context.Context, document interface{},
	_ ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if ctx.Value("withErr") == true {
		return &mongo.InsertOneResult{InsertedID: 0}, mongo.ErrNilDocument
	}
	m.data = document
	return &mongo.InsertOneResult{InsertedID: 1}, nil
}

type test struct {
	name    string
	args    args
	want    interface{}
	wantErr bool
}

type args struct {
	ctx     context.Context
	student domain.StudentRSS
}

func getTestCases(source string) []test {
	type key string
	k := key("withErr")
	return []test{
		{
			"no error",
			args{
				context.WithValue(context.Background(), k, false),
				domain.StudentRSS{},
			},
			domain.StudentRecord{Source: source},
			false,
		}, {
			"with error",
			args{
				context.WithValue(context.Background(), k, true),
				domain.StudentRSS{},
			},
			domain.StudentRecord{Source: source},
			true,
		},
	}
}

func TestStudentsRepo_SaveRSS(t *testing.T) {
	tests := getTestCases(domain.RSS)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := collectionMock{}
			repo := newRepo(&mock)

			_, err := repo.SaveRSS(tt.args.ctx, tt.args.student)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveRSS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(mock.data, tt.want) {
				t.Errorf("SaveRSS() got = %v, want %v", mock.data, tt.want)
			}
		})
	}
}

func TestStudentsRepo_SaveWAC(t *testing.T) {
	tests := getTestCases(domain.WAC)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := collectionMock{}
			repo := newRepo(&mock)

			_, err := repo.SaveRSS(tt.args.ctx, tt.args.student)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveWAC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.DeepEqual(mock.data, tt.want) {
				t.Errorf("SaveWAC() got = %v, want %v", mock.data, tt.want)
			}
		})
	}
}
