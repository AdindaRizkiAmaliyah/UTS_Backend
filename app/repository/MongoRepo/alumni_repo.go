package MongoRepo

import (
	"clean-archi/app/model"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Struct repository untuk MongoDB
type AlumniMongoRepository struct {
	collection *mongo.Collection
}

// Constructor
func NewAlumniMongoRepository(db *mongo.Database, collectionName string) *AlumniMongoRepository {
	return &AlumniMongoRepository{
		collection: db.Collection(collectionName),
	}
}

// Getter untuk akses koleksi (agar bisa dipakai AuthService)
func (r *AlumniMongoRepository) GetCollection() *mongo.Collection {
	return r.collection
}

// ========================= CRUD =========================

// Create - Tambah data alumni baru
func (r *AlumniMongoRepository) Create(ctx context.Context, alumni *model.Alumni) (*model.Alumni, error) {
	alumni.MongoID = primitive.NewObjectID()
	alumni.CreatedAt = time.Now()
	alumni.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, alumni)
	if err != nil {
		return nil, err
	}
	return alumni, nil
}

// GetAll - Ambil semua data alumni
func (r *AlumniMongoRepository) GetAll(ctx context.Context) ([]model.Alumni, error) {
	cur, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var alumniList []model.Alumni
	for cur.Next(ctx) {
		var a model.Alumni
		if err := cur.Decode(&a); err != nil {
			return nil, err
		}
		alumniList = append(alumniList, a)
	}
	return alumniList, nil
}

// GetByID - Ambil data alumni berdasarkan ID
func (r *AlumniMongoRepository) GetByID(ctx context.Context, id string) (*model.Alumni, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID: %v", err)
	}

	var alumni model.Alumni
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&alumni)
	if err != nil {
		return nil, err
	}
	return &alumni, nil
}

// Update - Update data alumni berdasarkan ID
func (r *AlumniMongoRepository) Update(ctx context.Context, id string, data *model.Alumni) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
	}

	data.UpdatedAt = time.Now()
	update := bson.M{"$set": data}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// Delete - Hapus data alumni berdasarkan ID
func (r *AlumniMongoRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID: %v", err)
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}
