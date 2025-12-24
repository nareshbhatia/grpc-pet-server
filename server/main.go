package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"connectrpc.com/connect"
	petService "github.com/nareshbhatia/grpc-pet-server/gen/go/pet/v1"
	petConnect "github.com/nareshbhatia/grpc-pet-server/gen/go/pet/v1/petv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const address = "localhost:8080"

func main() {
	mux := http.NewServeMux()
	path, handler := petConnect.NewPetServiceHandler(&petServiceHandler{})
	mux.Handle(path, handler)

	server := &http.Server{
		Addr:    address,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	// Start server in a goroutine
	go func() {
		fmt.Println("... Listening on", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	handleShutdown(server)
}

// handleShutdown handles Ctrl+C gracefully by canceling the context
func handleShutdown(server *http.Server) {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for interrupt signal
	<-sigChan
	fmt.Println("\nShutting down server...")

	// Gracefully shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server shutdown error: %v\n", err)
	} else {
		fmt.Println("Server stopped gracefully")
	}
}

// ------------------------------------------------------------
// petServiceHandler implements the PetService API.
// ------------------------------------------------------------
type petServiceHandler struct {
	petConnect.UnimplementedPetServiceHandler
}

// GetStatus returns the current status of the pet.
func (s *petServiceHandler) GetStatus(
	_ context.Context,
	req *connect.Request[petService.GetStatusRequest],
) (*connect.Response[petService.GetStatusResponse], error) {
	// Randomly select a status from the valid statuses
	randomStatus := validStatuses[rand.Intn(len(validStatuses))]

	return connect.NewResponse(&petService.GetStatusResponse{
		Status: randomStatus,
	}), nil
}

func (s *petServiceHandler) SubscribeHeartbeat(
	ctx context.Context,
	req *connect.Request[petService.SubscribeHeartbeatRequest],
	stream *connect.ServerStream[petService.SubscribeHeartbeatResponse],
) error {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := stream.Send(&petService.SubscribeHeartbeatResponse{
				TimestampMs: time.Now().UnixMilli(),
			}); err != nil {
				return err
			}
		}
	}
}

// valid statuses for the pet
var validStatuses = []petService.PetStatus{
	petService.PetStatus_PET_STATUS_SLEEPING,
	petService.PetStatus_PET_STATUS_EATING,
	petService.PetStatus_PET_STATUS_DRINKING,
	petService.PetStatus_PET_STATUS_CHEWING,
	petService.PetStatus_PET_STATUS_PLAYING,
	petService.PetStatus_PET_STATUS_RUNNING,
	petService.PetStatus_PET_STATUS_TRAINING,
	petService.PetStatus_PET_STATUS_BARKING,
	petService.PetStatus_PET_STATUS_CUDDLING,
	petService.PetStatus_PET_STATUS_LICKING,
}
