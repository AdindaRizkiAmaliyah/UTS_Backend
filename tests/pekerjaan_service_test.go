package tests

import (
    "clean-archi/app/model"
    mockRepo "clean-archi/app/repository/MongoRepo/mock"
    "testing"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

func TestPekerjaanRepository(t *testing.T) {

    repo := mockRepo.NewMockPekerjaanRepo()

    // ===================== CREATE ======================

    alumniID := primitive.NewObjectID()
    pekerjaan := model.Pekerjaan{
        AlumniIDs:          []primitive.ObjectID{alumniID},
        NamaPerusahaan:     "Google",
        PosisiJabatan:      "Engineer",
        BidangIndustri:     "Teknologi",
        LokasiKerja:        "California",
        DeskripsiPekerjaan: "Software Engineer",
        TanggalMulaiKerja:  time.Now(),
        StatusPekerjaan:    "Aktif",
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    if err := repo.Create(&pekerjaan); err != nil {
        t.Fatalf("Create error: %v", err)
    }

    savedID := pekerjaan.MongoID.Hex()

    // ===================== GET BY ID ======================

    stored, err := repo.GetByID(savedID)
    if err != nil {
        t.Fatalf("GetByID error: %v", err)
    }
    if stored.NamaPerusahaan != "Google" {
        t.Fatalf("Expected Google, got %s", stored.NamaPerusahaan)
    }

    // ===================== UPDATE ======================

    pekerjaan.NamaPerusahaan = "Google Update"
    if err := repo.Update(savedID, &pekerjaan); err != nil {
        t.Fatalf("Update error: %v", err)
    }

    updated, _ := repo.GetByID(savedID)
    if updated.NamaPerusahaan != "Google Update" {
        t.Fatalf("Expected Google Update, got %s", updated.NamaPerusahaan)
    }

    // ===================== SOFT DELETE ======================

    if err := repo.SoftDeleteByAlumni(savedID, alumniID.Hex()); err != nil {
        t.Fatalf("SoftDelete error: %v", err)
    }

    sd := repo.Data[savedID]
    if !sd.IsDeleted {
        t.Fatalf("Soft delete failed")
    }

    // ===================== RESTORE ======================

    if err := repo.RestoreJob(savedID); err != nil {
        t.Fatalf("Restore error: %v", err)
    }

    restored := repo.Data[savedID]
    if restored.IsDeleted {
        t.Fatalf("Restore failed")
    }

    // ===================== HARD DELETE ======================

    if err := repo.HardDeleteJob(savedID); err != nil {
        t.Fatalf("HardDelete error: %v", err)
    }

    if _, ok := repo.Data[savedID]; ok {
        t.Fatalf("Hard delete did not remove data")
    }
}
