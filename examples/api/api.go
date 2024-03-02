package api

const NotesService = "NotesService"

type GetNotesReq struct {
}

type GetNotesRes struct {
	Notes []Note `json:"notes"`
}

type CreateNoteReq struct {
	Note Note `json:"note"`
}

type CreateNoteRes struct {
}

type UpdateNoteReq struct {
	Note Note `json:"note"`
}

type UpdateNoteRes struct {
}

type DeleteNoteReq struct {
	Note Note `json:"note"`
}

type DeleteNoteRes struct {
}

type Note struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
