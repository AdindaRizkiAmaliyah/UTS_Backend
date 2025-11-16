package tests

import (
    "clean-archi/app/model"
    mockRepo "clean-archi/app/repository/MongoRepo/mock"
    "clean-archi/app/service"
    "context"
    "testing"
)

func TestAlumniService(t *testing.T) {
    ctx := context.Background()

    // Pakai mock repository
    repo := mockRepo.NewMockAlumniRepo()

    // Pakai service yang sudah menerima interface
    svc := service.NewAlumniService(repo)

    // ================== CREATE ==================
    alumni := &model.Alumni{
        Nama:  "Adinda",
        Email: "adinda@example.com",
        NIM:   "202112345",
    }

    created, err := svc.CreateAlumni(ctx, alumni)
    if err != nil {
        t.Fatalf("Create error: %v", err)
    }

    if created.Nama != "Adinda" {
        t.Fatalf("Expected name Adinda, got %s", created.Nama)
    }

    id := created.MongoID.Hex()

    // ================== GET ==================
    fetched, err := svc.GetAlumniByID(ctx, id)
    if err != nil {
        t.Fatalf("Get error: %v", err)
    }

    if fetched.Email != "adinda@example.com" {
        t.Fatalf("Expected email adinda@example.com, got %s", fetched.Email)
    }

    // ================== UPDATE ==================
    alumni.Nama = "Adinda Update"
    err = svc.UpdateAlumni(ctx, id, alumni)
    if err != nil {
        t.Fatalf("Update error: %v", err)
    }

    updated, _ := svc.GetAlumniByID(ctx, id)
    if updated.Nama != "Adinda Update" {
        t.Fatalf("Update failed, expected 'Adinda Update', got %s", updated.Nama)
    }

    // ================== DELETE ==================
    err = svc.DeleteAlumni(ctx, id)
    if err != nil {
        t.Fatalf("Delete error: %v", err)
    }

    // Setelah delete, Get harus error
    _, err = svc.GetAlumniByID(ctx, id)
    if err == nil {
        t.Fatalf("Expected error after delete, got nil")
    }
}
