package repository

import (
	"clean-archi/app/model"
)

type PekerjaanRepository interface {
	// CRUD
	Create(p *model.Pekerjaan) error
	GetAll() ([]model.Pekerjaan, error)
	GetByID(id string) (*model.Pekerjaan, error)
	Update(id string, p *model.Pekerjaan) error
	Delete(id string) error

	// Soft Delete
	SoftDeleteByAlumni(jobID, alumniID string) error
	SoftDeleteAllByAdmin(alumniID, adminID string) (int64, error)

	// Trash
	GetTrashedJobs() ([]model.Pekerjaan, error)
	GetTrashedJobsByAlumni(alumniID string) ([]model.Pekerjaan, error)

	// Restore
	RestoreJob(jobID string) error
	RestoreJobByAlumni(jobID, alumniID string) error

	// Hard Delete
	HardDeleteJob(jobID string) error
	HardDeleteJobByAlumni(jobID, alumniID string) error
}
