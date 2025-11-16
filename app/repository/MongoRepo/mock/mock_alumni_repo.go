package mock

import (
    "clean-archi/app/model"
    "context"
    "errors"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

type MockAlumniRepo struct {
    Data map[string]model.Alumni
}

func NewMockAlumniRepo() *MockAlumniRepo {
    return &MockAlumniRepo{
        Data: make(map[string]model.Alumni),
    }
}

func (m *MockAlumniRepo) Create(ctx context.Context, a *model.Alumni) (*model.Alumni, error) {
    a.MongoID = primitive.NewObjectID()
    m.Data[a.MongoID.Hex()] = *a
    return a, nil
}

func (m *MockAlumniRepo) GetAll(ctx context.Context) ([]model.Alumni, error) {
    var out []model.Alumni
    for _, v := range m.Data {
        out = append(out, v)
    }
    return out, nil
}

func (m *MockAlumniRepo) GetByID(ctx context.Context, id string) (*model.Alumni, error) {
    if v, ok := m.Data[id]; ok {
        return &v, nil
    }
    return nil, errors.New("alumni not found")
}

func (m *MockAlumniRepo) Update(ctx context.Context, id string, a *model.Alumni) error {
    if _, ok := m.Data[id]; !ok {
        return errors.New("not found")
    }
    m.Data[id] = *a
    return nil
}

func (m *MockAlumniRepo) Delete(ctx context.Context, id string) error {
    delete(m.Data, id)
    return nil
}
