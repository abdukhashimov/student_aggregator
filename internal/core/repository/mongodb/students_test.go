package mongodb

import (
	"context"
	"github.com/abdukhashimov/student_aggregator/internal/core/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"reflect"
	"testing"
)

type key string

var k = key("withErr")

type collectionMock struct {
	data interface{}
}

func (m *collectionMock) InsertOne(ctx context.Context, document interface{},
	_ ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	if ctx.Value(k) == true {
		return &mongo.InsertOneResult{InsertedID: 0}, mongo.ErrNilDocument
	}
	m.data = document
	return &mongo.InsertOneResult{InsertedID: 1}, nil
}

func (m *collectionMock) FindOne(ctx context.Context, filter interface{},
	opts ...*options.FindOneOptions) *mongo.SingleResult {
	panic("implement me")
}

func (m *collectionMock) Find(ctx context.Context, filter interface{},
	opts ...*options.FindOptions) (cur *mongo.Cursor, err error) {
	panic("implement me")
}

func (m *collectionMock) UpdateOne(ctx context.Context, filter interface{}, update interface{},
	opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	panic("implement me")
}

func (m *collectionMock) DeleteOne(ctx context.Context, filter interface{},
	opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	panic("implement me")
}

func TestStudentsRepo_SaveRSS(t *testing.T) {
	type args struct {
		ctx     context.Context
		student domain.StudentRSS
	}
	type test struct {
		name    string
		args    args
		want    domain.StudentRecord
		wantErr bool
	}

	tests := []test{
		{
			"no error",
			args{
				context.WithValue(context.Background(), k, false),
				domain.StudentRSS{},
			},
			domain.StudentRecord{Source: domain.RSS},
			false,
		}, {
			"with error",
			args{
				context.WithValue(context.Background(), k, true),
				domain.StudentRSS{},
			},
			domain.StudentRecord{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := collectionMock{}
			repo := newRepo(&mock)

			_, err := repo.SaveRSS(tt.args.ctx, "", "", tt.args.student)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SaveRss() error expected but got %v", err)
				}
				return
			}

			if !reflect.DeepEqual(mock.data, tt.want) {
				t.Errorf("SaveRSS() got = %v, want %v", mock.data, tt.want)
			}
		})
	}
}

func TestStudentsRepo_SaveWAC(t *testing.T) {
	type args struct {
		ctx     context.Context
		student domain.StudentWAC
	}
	type test struct {
		name    string
		args    args
		want    domain.StudentRecord
		wantErr bool
	}
	tests := []test{
		{
			"no error",
			args{
				context.WithValue(context.Background(), k, false),
				domain.StudentWAC{},
			},
			domain.StudentRecord{Source: domain.WAC},
			false,
		}, {
			"with error",
			args{
				context.WithValue(context.Background(), k, true),
				domain.StudentWAC{},
			},
			domain.StudentRecord{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := collectionMock{}
			repo := newRepo(&mock)

			_, err := repo.SaveWAC(tt.args.ctx, "", "", tt.args.student)
			if tt.wantErr {
				if err == nil {
					t.Errorf("SaveWAC() error expected but got %v", err)
				}
				return
			}

			if !reflect.DeepEqual(mock.data, tt.want) {
				t.Errorf("SaveWAC() got = %v, want %v", mock.data, tt.want)
			}
		})
	}
}
