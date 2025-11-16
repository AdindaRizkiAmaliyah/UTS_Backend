package MongoRepo

import (
	"clean-archi/app/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(client *mongo.Client, dbName string) *FileRepository {
	return &FileRepository{
		collection: client.Database(dbName).Collection("files"),
	}
}

func (r *FileRepository) Create(file *model.File) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file.UploadedAt = time.Now()
	result, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		return err
	}
	file.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// Ambil semua file
func (r *FileRepository) GetAll() ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []model.File
	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}

// Ambil file berdasarkan ID
func (r *FileRepository) GetByID(id string) (*model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var file model.File
	err = r.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&file)
	if err == mongo.ErrNoDocuments {
		return nil, errors.New("file tidak ditemukan")
	}
	return &file, err
}

// Hapus file berdasarkan ID
func (r *FileRepository) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objID})
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return errors.New("file tidak ditemukan")
	}
	return nil
}

// Ambil file berdasarkan user ID
func (r *FileRepository) GetByUserID(userID string) ([]model.File, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []model.File
	if err := cursor.All(ctx, &files); err != nil {
		return nil, err
	}
	return files, nil
}
