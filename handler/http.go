package handler

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"src/petstore"
	"src/proto_storage"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type PetHandler struct {
	// Storage *storageservice.StorageService
	Client proto_storage.StorageServiceClient
}

// func NewPetHandler(storageURL string) (*PetHandler, error) {

// 	conn, err := grpc.NewClient(storageURL, grpc.WithTransportCredentials(insecure.NewCredentials()))

// 	if err != nil {
// 		return nil, err
// 	}
// 	return &PetHandler{
// 		Client: proto_storage.NewStorageServiceClient(conn),
// 	}, nil
// }

func NewPetHandler(storageURL string) (*PetHandler, error) {
	// Load the server's certificate
	serverCert, err := os.ReadFile("./cert/server.crt")
	if err != nil {
		return nil, err
	}

	// Create a certificate pool
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(serverCert) {
		return nil, fmt.Errorf("failed to append server certs")
	}

	// Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}
	creds := credentials.NewTLS(config)

	conn, err := grpc.NewClient(storageURL, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &PetHandler{
		Client: proto_storage.NewStorageServiceClient(conn),
	}, nil
}

func writeError(w http.ResponseWriter, code int32, err error) {
	if errStatus, ok := status.FromError(err); ok {
		if errStatus.Code() == codes.NotFound {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(petstore.Error{Code: http.StatusNotFound, Message: errStatus.Message()})
			return
		}

	}
	if code == 0 {
		code = 500
	}
	petErr := petstore.Error{
		Code:    code,
		Message: err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(code))
	err = json.NewEncoder(w).Encode(petErr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h PetHandler) FindPets(w http.ResponseWriter, r *http.Request, params petstore.FindPetsParams) {
	tags := []string{}
	if params.Tags != nil {
		tags = *params.Tags
	}
	var limit int64
	if params.Limit != nil {
		limit = int64(*params.Limit)
	}
	req := proto_storage.FindPetsRequest{
		Limit: int32(limit),
		Tags:  tags,
	}
	pets, err := h.Client.FindPets(r.Context(), &req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if pets == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, pets)
}

func (h PetHandler) AddPet(w http.ResponseWriter, r *http.Request) {
	var newPet petstore.NewPet
	err := json.NewDecoder(r.Body).Decode(&newPet)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	tags := []*proto_storage.Tag{}
	for _, tag := range newPet.Tags {
		tags = append(tags, &proto_storage.Tag{Id: tag.Id, Name: tag.Name})
	}

	pet, err := h.Client.CreatePet(r.Context(), &proto_storage.NewPet{
		Name:      newPet.Name,
		Category:  &proto_storage.Category{Name: newPet.Category.Name},
		Status:    newPet.Status,
		PhotoUrls: newPet.PhotoUrls,
		Tags:      tags,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, pet)
}

func (h PetHandler) DeletePet(w http.ResponseWriter, r *http.Request, id int64) {
	_, err := h.Client.DeletePet(r.Context(), &proto_storage.PetID{Id: id})
	if err != nil {
		writeError(w, 0, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PetHandler) FindPetByID(w http.ResponseWriter, r *http.Request, id int64) {
	pet, err := h.Client.FindPetById(r.Context(), &proto_storage.PetID{Id: id})
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	if pet == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, pet)
}
