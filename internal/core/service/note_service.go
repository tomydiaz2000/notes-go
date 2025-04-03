package service

import (
	"errors"
	"notes-app/internal/core/domain"
	"notes-app/internal/core/ports/driven"
	"notes-app/internal/core/ports/driver"
	"time"
)

type NoteService struct {
	noteRepository driven.NoteRepository
}

func NewNoteService(noteRepository driven.NoteRepository) driver.NoteService {
	return &NoteService{
		noteRepository: noteRepository,
	}
}

func (s *NoteService) Create(note *domain.Note) error {
	if note == nil {
		return errors.New("note cannot be nil")
	}

	note.CreateAt = time.Now()
	err := s.noteRepository.Create(note)
	if err != nil {
		return err
	}
	return nil
}

func (s *NoteService) Update(note *domain.Note) error {
	if note == nil {
		return errors.New("note cannot be nil")
	}
	err := s.noteRepository.Update(note)
	if err != nil {
		return err
	}
	return nil
}

func (s *NoteService) Delete(id string) error {
	note, err := s.noteRepository.FindById(id)
	if err != nil {
		return err
	}
	err = s.noteRepository.Delete(note.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *NoteService) FindById(id string) (*domain.Note, error) {
	return s.noteRepository.FindById(id)
}

func (s *NoteService) FindAll() ([]*domain.Note, error) {
	return s.noteRepository.FindAll()
}

func (s *NoteService) DeleteAll() error {
	return s.noteRepository.DeleteAll()
}
