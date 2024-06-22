package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"src/petstore"
	"src/storageservice"
)

type PetHandler struct {
	Storage *storageservice.StorageService
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)
	w.Write([]byte(err.Error()))
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h PetHandler) FindPets(w http.ResponseWriter, r *http.Request, params petstore.FindPetsParams) {
	tags := *params.Tags
	fmt.Printf("%#v\n", tags)
	var limit int64
	if params.Limit != nil {
		limit = int64(*params.Limit)
	}

	pets, err := h.Storage.FindPets(r.Context(), limit, tags)
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
	fmt.Println("name", newPet.Name)
	for _, tag := range *newPet.Tags {
		fmt.Println("tag", *tag.Name)
	}
	fmt.Println("status", *newPet.Status)
	fmt.Println("photoUrls", *newPet.PhotoUrls)
	fmt.Println("category", *newPet.Category.Name)

	pet, err := h.Storage.CreatePet(r.Context(), &newPet)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusCreated, pet)
}

func (h PetHandler) DeletePet(w http.ResponseWriter, r *http.Request, id int64) {
	err := h.Storage.DeletePet(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h PetHandler) FindPetByID(w http.ResponseWriter, r *http.Request, id int64) {
	pet, err := h.Storage.FindPetByID(r.Context(), id)
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
