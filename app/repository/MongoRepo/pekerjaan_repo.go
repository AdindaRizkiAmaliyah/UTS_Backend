package MongoRepo

import (
	"clean-archi/app/model"
	"context"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PekerjaanMongoRepository struct {
	collection *mongo.Collection
}

func NewPekerjaanMongoRepository(db *mongo.Database, collectionName string) *PekerjaanMongoRepository {
	return &PekerjaanMongoRepository{
		collection: db.Collection(collectionName),
	}
}

// ========================= CRUD =========================

func (r *PekerjaanMongoRepository) GetAll() ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"is_deleted": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Pekerjaan
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func (r *PekerjaanMongoRepository) GetByID(id string) (*model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var pekerjaan model.Pekerjaan
	err = r.collection.FindOne(ctx, bson.M{"_id": objID, "is_deleted": false}).Decode(&pekerjaan)
	if err != nil {
		return nil, err
	}
	return &pekerjaan, nil
}

func (r *PekerjaanMongoRepository) Create(pekerjaan *model.Pekerjaan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pekerjaan.CreatedAt = time.Now()
	pekerjaan.UpdatedAt = time.Now()
	pekerjaan.IsDeleted = false

	if pekerjaan.MongoID.IsZero() {
		pekerjaan.MongoID = primitive.NewObjectID()
	}

	result, err := r.collection.InsertOne(ctx, pekerjaan)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		pekerjaan.MongoID = oid
	}

	return nil
}

func (r *PekerjaanMongoRepository) Update(id string, pekerjaan *model.Pekerjaan) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	pekerjaan.UpdatedAt = time.Now()
	update := bson.M{"$set": pekerjaan}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

func (r *PekerjaanMongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// ========================= SOFT DELETE =========================

func (r *PekerjaanMongoRepository) SoftDeleteByAlumni(jobID, alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}

	alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_by": alumniObjID,
			"updated_at": time.Now(),
		},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": jobObjID, "alumni_id": alumniObjID}, update)
	return err
}

func (r *PekerjaanMongoRepository) SoftDeleteAllByAdmin(alumniID string, adminID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return 0, err
	}

	adminObjID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return 0, err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"deleted_by": adminObjID,
			"updated_at": time.Now(),
		},
	}

	res, err := r.collection.UpdateMany(ctx, bson.M{"alumni_id": alumniObjID}, update)
	if err != nil {
		return 0, err
	}

	return res.ModifiedCount, nil
}

// ========================= TRASH =========================

func (r *PekerjaanMongoRepository) GetTrashedJobs() ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"is_deleted": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Pekerjaan
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *PekerjaanMongoRepository) GetTrashedJobsByAlumni(alumniID string) ([]model.Pekerjaan, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if !primitive.IsValidObjectID(alumniID) {
		return nil, fmt.Errorf("alumni_id '%s' bukan ObjectID valid", alumniID)
	}

	alumniObjID, _ := primitive.ObjectIDFromHex(alumniID)

	cursor, err := r.collection.Find(ctx, bson.M{
		"alumni_id":  alumniObjID,
		"is_deleted": true,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []model.Pekerjaan
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

// ========================= RESTORE =========================

func (r *PekerjaanMongoRepository) RestoreJob(jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set":   bson.M{"is_deleted": false, "updated_at": time.Now()},
		"$unset": bson.M{"deleted_by": ""},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": jobObjID}, update)
	return err
}

func (r *PekerjaanMongoRepository) RestoreJobByAlumni(jobID, alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}

	alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set":   bson.M{"is_deleted": false, "updated_at": time.Now()},
		"$unset": bson.M{"deleted_by": ""},
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": jobObjID, "alumni_id": alumniObjID}, update)
	return err
}

// ========================= HARD DELETE =========================

func (r *PekerjaanMongoRepository) HardDeleteJob(jobID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": jobObjID})
	return err
}

func (r *PekerjaanMongoRepository) HardDeleteJobByAlumni(jobID, alumniID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobObjID, err := primitive.ObjectIDFromHex(jobID)
	if err != nil {
		return err
	}

	alumniObjID, err := primitive.ObjectIDFromHex(alumniID)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": jobObjID, "alumni_id": alumniObjID})
	return err
}
