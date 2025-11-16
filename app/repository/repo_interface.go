package repository

import (
    "context"
    "clean-archi/app/model"
)

// gunakan string id di interface untuk konsistensi dengan code-mu
type AlumniRepository interface {
    Create(ctx context.Context, a *model.Alumni) (*model.Alumni, error)
    GetAll(ctx context.Context) ([]model.Alumni, error)
    GetByID(ctx context.Context, id string) (*model.Alumni, error)
    Update(ctx context.Context, id string, a *model.Alumni) error
    Delete(ctx context.Context, id string) error
}