package web

import (
	"net/http"
	"notes-app/internal/core/domain"
	"notes-app/internal/core/ports/driver"
	"notes-app/internal/infrastructure/adapters/web/dto"

	"github.com/gin-gonic/gin"
)

type NoteHandler struct {
	noteService driver.NoteService
}

func NewNoteHandler(noteService driver.NoteService) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (h *NoteHandler) CreateNote(c *gin.Context) {
	var req dto.CreateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := &domain.Note{
		Title:   req.Title,
		Content: req.Content,
	}

	err := h.noteService.Create(note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note created successfully", "id": note.ID})
}

func (h *NoteHandler) GetNotes(c *gin.Context) {
	notes, err := h.noteService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func (h *NoteHandler) GetNote(c *gin.Context) {
	id := c.Param("id")
	note, err := h.noteService.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"note": note})
}

func (h *NoteHandler) UpdateNote(c *gin.Context) {
	var req dto.UpdateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id cannot be empty"})
		return
	}

	note := &domain.Note{
		ID:      id,
		Title:   req.Title,
		Content: req.Content,
	}

	err := h.noteService.Update(note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note updated successfully"})
}

func (h *NoteHandler) DeleteNote(c *gin.Context) {
	id := c.Param("id")
	err := h.noteService.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func (h *NoteHandler) DeleteAll(c *gin.Context) {
	err := h.noteService.DeleteAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notes deleted successfully"})
}
