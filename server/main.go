package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	petService "github.com/bufbuild/buf-examples/gen/pet/v1"
	petConnect "github.com/bufbuild/buf-examples/gen/pet/v1/petv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const address = "localhost:8080"

func main() {
	mux := http.NewServeMux()
	path, handler := petConnect.NewPetStoreServiceHandler(&petStoreServiceServer{})
	mux.Handle(path, handler)
	fmt.Println("... Listening on", address)
	http.ListenAndServe(
		address,
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}

// petStoreServiceServer implements the PetStoreService API.
type petStoreServiceServer struct {
	petConnect.UnimplementedPetStoreServiceHandler
}

// PutPet adds the pet associated with the given request into the PetStore.
func (s *petStoreServiceServer) PutPet(
	_ context.Context,
	req *petService.PutPetRequest,
) (*petService.PutPetResponse, error) {
	pet := &petService.Pet{Name: req.GetName(), PetType: req.PetType}
	log.Printf("PutPet received a %v named %s", pet.GetPetType(), pet.GetName())
	return &petService.PutPetResponse{Pet: pet}, nil
}
