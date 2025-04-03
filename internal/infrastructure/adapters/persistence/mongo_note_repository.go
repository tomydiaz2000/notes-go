package persistence

import (
	"context"
	"errors"
	"notes-app/internal/core/domain"
	"notes-app/internal/core/ports/driven"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoNoteRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMongoNoteRepository(client *mongo.Client, dbName, collectionName string) driven.NoteRepository {
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	return &MongoNoteRepository{
		client:     client,
		database:   database,
		collection: collection,
	}
}

// NoteDocument es una representación de la entidad Note para MongoDB.
type NoteDocument struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Title        string             `bson:"title"`
	Content      string             `bson:"content"`
	CreateAt     time.Time          `bson:"created_at"`
	UpdateAt     time.Time          `bson:"updated_at"`
	ValidUntilAt time.Time          `bson:"valid_until_at"`
}

func toNoteDocument(noteDocument *domain.Note) *NoteDocument {
	objectID, _ := primitive.ObjectIDFromHex(noteDocument.ID)
	return &NoteDocument{
		ID:           objectID,
		Title:        noteDocument.Title,
		Content:      noteDocument.Content,
		CreateAt:     noteDocument.CreateAt,
		UpdateAt:     noteDocument.UpdateAt,
		ValidUntilAt: noteDocument.ValidUntilAt,
	}
}

func fromNoteDocument(note *NoteDocument) *domain.Note {
	return &domain.Note{
		ID:           note.ID.Hex(),
		Title:        note.Title,
		Content:      note.Content,
		ValidUntilAt: note.ValidUntilAt,
		CreateAt:     note.CreateAt,
		UpdateAt:     note.UpdateAt,
	}
}

func (r *MongoNoteRepository) Create(note *domain.Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	note.CreateAt = time.Now()
	doc := toNoteDocument(note)
	result, err := r.collection.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if objectID, ok := result.InsertedID.(primitive.ObjectID); ok {
		note.ID = objectID.Hex()
	}
	return nil
}

func (r *MongoNoteRepository) Update(note *domain.Note) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(note.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	var existingNote NoteDocument
	err = r.collection.FindOne(ctx, filter).Decode(&existingNote)
	if err == mongo.ErrNoDocuments {
		return errors.New("note not found")
	}
	if err != nil {
		return err
	}

	noteDocument := existingNote

	noteDocument.Title = note.Title
	noteDocument.Content = note.Content
	noteDocument.ValidUntilAt = note.ValidUntilAt
	noteDocument.UpdateAt = time.Now()

	_, err = r.collection.UpdateOne(ctx, filter, bson.M{"$set": noteDocument})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoNoteRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}

	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

// Delete all
func (r *MongoNoteRepository) DeleteAll() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.collection.DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoNoteRepository) FindById(id string) (*domain.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID} // Filter para buscar el documento por ID

	var result NoteDocument
	err = r.collection.FindOne(ctx, filter).Decode(&result)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	noteDocument := fromNoteDocument(&result)
	return noteDocument, nil
}

func (r *MongoNoteRepository) FindAll() ([]*domain.Note, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	notes := make([]*domain.Note, 0)
	for cursor.Next(ctx) {
		noteDocument := NoteDocument{}
		err := cursor.Decode(&noteDocument)
		if err != nil {
			return nil, err
		}
		notes = append(notes, fromNoteDocument(&noteDocument))
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}
