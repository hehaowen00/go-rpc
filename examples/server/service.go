package main

import (
	"context"
	"errors"
	"strconv"
	"sync"

	"github.com/hehaowen00/go-rpc"
	"github.com/hehaowen00/go-rpc/examples/api"
)

type NotesService struct {
	notes   []api.Note
	counter int64
	mu      sync.Mutex
}

func (s *NotesService) Init(service *rpc.Service) {
	rpc.Register(service, "GetNotes", s.GetNotes)
	rpc.Register(service, "CreateNote", s.CreateNote)
	rpc.Register(service, "UpdateNote", s.UpdateNote)
	rpc.Register(service, "DeleteNote", s.DeleteNote)
	rpc.Register(service, "Error", s.Error)
}

func (s *NotesService) GetNotes(
	ctx context.Context,
	r *api.GetNotesReq,
) (*api.GetNotesRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	res := &api.GetNotesRes{
		Notes: s.notes,
	}

	return res, nil
}

func (s *NotesService) CreateNote(
	ctx context.Context,
	r *api.CreateNoteReq,
) (*api.CreateNoteRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	r.Note.ID = strconv.FormatInt(s.counter, 10)
	s.counter++
	s.notes = append(s.notes, r.Note)

	res := &api.CreateNoteRes{}

	return res, nil
}

func (s *NotesService) UpdateNote(
	ctx context.Context,
	r *api.UpdateNoteReq,
) (*api.UpdateNoteRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < len(s.notes); i++ {
		if s.notes[i].ID == r.Note.ID {
			s.notes[i] = r.Note
		}
	}

	return nil, nil
}

func (s *NotesService) DeleteNote(
	ctx context.Context,
	r *api.DeleteNoteReq,
) (*api.DeleteNoteRes, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := 0; i < len(s.notes); i++ {
		if s.notes[i].ID == r.Note.ID {
			s.notes = append(s.notes[:i], s.notes[i+1:]...)
			return &api.DeleteNoteRes{}, nil
		}
	}

	return nil, errors.New("note not found")
}

func (s *NotesService) Error(
	ctx context.Context,
	r *api.CreateNoteReq,
) (*api.CreateNoteReq, error) {
	return nil, errors.New("error processing request")
}
