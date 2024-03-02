package main

import (
	"context"
	"log"

	"github.com/hehaowen00/go-rpc"
	"github.com/hehaowen00/go-rpc/examples/api"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := rpc.NewServiceClient(
		api.NotesService,
		"http://127.0.0.1:8080",
	)
	if err != nil {
		log.Println(err)
	}

	log.Println("call create note")
	{
		req := rpc.NewRequest[api.CreateNoteReq, api.CreateNoteRes](
			context.Background(),
			"CreateNote",
			&api.CreateNoteReq{
				Note: api.Note{
					Title: "title",
					Body:  "body",
				},
			})

		_, err := rpc.Call(client, req)
		if err != nil {
			log.Fatalln(err)
		}
	}

	log.Println("call get notes")
	{
		req := rpc.NewRequest[api.GetNotesReq, api.GetNotesRes](
			context.Background(),
			"GetNotes",
			&api.GetNotesReq{},
		)

		res, err := rpc.Call(client, req)
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("%+v\n", res)
	}

	log.Println("call error")
	{
		req := rpc.NewRequest[api.CreateNoteReq, api.CreateNoteRes](
			context.Background(),
			"Error",
			&api.CreateNoteReq{
				Note: api.Note{
					Title: "title",
					Body:  "body",
				},
			},
		)

		_, err := rpc.Call(client, req)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
