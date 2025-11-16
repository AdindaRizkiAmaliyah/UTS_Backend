package mock

import (
    "clean-archi/app/model"
    "clean-archi/app/repository"
    "errors"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

var _ repository.PekerjaanRepository = (*MockPekerjaanRepo)(nil)

type MockPekerjaanRepo struct {
    Data map[string]model.Pekerjaan
}

func NewMockPekerjaanRepo() *MockPekerjaanRepo {
    return &MockPekerjaanRepo{
        Data: make(map[string]model.Pekerjaan),
    }
}

// ======================== CRUD ========================

func (m *MockPekerjaanRepo) Create(p *model.Pekerjaan) error {
    p.MongoID = primitive.NewObjectID()
    m.Data[p.MongoID.Hex()] = *p
    return nil
}

func (m *MockPekerjaanRepo) GetAll() ([]model.Pekerjaan, error) {
    var out []model.Pekerjaan
    for _, v := range m.Data {
        if !v.IsDeleted {
            out = append(out, v)
        }
    }
    return out, nil
}

func (m *MockPekerjaanRepo) GetByID(id string) (*model.Pekerjaan, error) {
    if v, ok := m.Data[id]; ok && !v.IsDeleted {
        return &v, nil
    }
    return nil, errors.New("not found")
}

func (m *MockPekerjaanRepo) Update(id string, p *model.Pekerjaan) error {
    if _, ok := m.Data[id]; !ok {
        return errors.New("not found")
    }
    m.Data[id] = *p
    return nil
}

func (m *MockPekerjaanRepo) Delete(id string) error {
    delete(m.Data, id)
    return nil
}

// ======================== SOFT DELETE ========================

func (m *MockPekerjaanRepo) SoftDeleteByAlumni(jobID, alumniID string) error {
    p, ok := m.Data[jobID]
    if !ok {
        return errors.New("not found")
    }
    p.IsDeleted = true
    p.DeletedBy = &alumniID
    m.Data[jobID] = p
    return nil
}

func (m *MockPekerjaanRepo) SoftDeleteAllByAdmin(alumniID, adminID string) (int64, error) {
    var count int64 = 0
    for id, p := range m.Data {
        for _, a := range p.AlumniIDs {
            if a.Hex() == alumniID {
                p.IsDeleted = true
                p.DeletedBy = &adminID
                m.Data[id] = p
                count++
            }
        }
    }
    return count, nil
}

// ======================== TRASH ========================

func (m *MockPekerjaanRepo) GetTrashedJobs() ([]model.Pekerjaan, error) {
    var out []model.Pekerjaan
    for _, v := range m.Data {
        if v.IsDeleted {
            out = append(out, v)
        }
    }
    return out, nil
}

func (m *MockPekerjaanRepo) GetTrashedJobsByAlumni(alumniID string) ([]model.Pekerjaan, error) {
    var out []model.Pekerjaan
    for _, p := range m.Data {
        if p.IsDeleted {
            for _, a := range p.AlumniIDs {
                if a.Hex() == alumniID {
                    out = append(out, p)
                }
            }
        }
    }
    return out, nil
}

// ======================== RESTORE ========================

func (m *MockPekerjaanRepo) RestoreJob(jobID string) error {
    p, ok := m.Data[jobID]
    if !ok {
        return errors.New("not found")
    }
    p.IsDeleted = false
    p.DeletedBy = nil
    m.Data[jobID] = p
    return nil
}

func (m *MockPekerjaanRepo) RestoreJobByAlumni(jobID, alumniID string) error {
    p, ok := m.Data[jobID]
    if !ok {
        return errors.New("not found")
    }

    owned := false
    for _, id := range p.AlumniIDs {
        if id.Hex() == alumniID {
            owned = true
            break
        }
    }
    if !owned {
        return errors.New("tidak boleh restore pekerjaan milik orang lain")
    }

    p.IsDeleted = false
    p.DeletedBy = nil
    m.Data[jobID] = p
    return nil
}

// ======================== HARD DELETE ========================

func (m *MockPekerjaanRepo) HardDeleteJob(jobID string) error {
    delete(m.Data, jobID)
    return nil
}

func (m *MockPekerjaanRepo) HardDeleteJobByAlumni(jobID, alumniID string) error {
    p, ok := m.Data[jobID]
    if !ok {
        return errors.New("not found")
    }

    owner := false
    for _, id := range p.AlumniIDs {
        if id.Hex() == alumniID {
            owner = true
        }
    }
    if !owner {
        return errors.New("tidak memiliki pekerjaan ini")
    }

    delete(m.Data, jobID)
    return nil
}
