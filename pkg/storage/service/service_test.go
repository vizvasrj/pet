package service

import (
	"context"
	"testing"

	"src/env"
	"src/pkg/storage/database"
	"src/proto_storage"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"go.uber.org/zap"
)

func TestStorageService(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Storage Service Suite")
}

var _ = ginkgo.Describe("StorageService", func() {
	var (
		service *StorageService
		logger  *zap.Logger
		testPet *proto_storage.NewPet
		ctx     context.Context
		// petID   int64 = 1
	)

	ginkgo.BeforeEach(func() {
		var err error
		e := env.GetEnvs()
		db := database.GetConnection(e)
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		logger, err = zap.NewDevelopment()
		gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

		service = &StorageService{
			Db:     db,
			Logger: logger,
		}

		testPet = &proto_storage.NewPet{
			Name: "Test Pet",
			Category: &proto_storage.Category{
				Name: "Test Category",
			},
			PhotoUrls: []string{"http://example.com/photo1.jpg", "http://example.com/photo2.jpg"},
			Tags: []*proto_storage.Tag{
				{Name: "Tag1"},
				{Name: "Tag2"},
			},
			Status: "available",
		}

		ctx = context.Background()
	})

	ginkgo.Describe("CreatePet", func() {
		ginkgo.It("should successfully create a pet with new category and tags", func() {
			// Call the service method
			pet, err := service.CreatePet(ctx, testPet)

			// Assertions
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(pet).ShouldNot(gomega.BeNil())
			gomega.Expect(pet.Name).Should(gomega.Equal(testPet.Name))
			gomega.Expect(pet.Category.Name).Should(gomega.Equal(testPet.Category.Name))
			gomega.Expect(pet.PhotoUrls).Should(gomega.Equal(testPet.PhotoUrls))
			gomega.Expect(pet.Status).Should(gomega.Equal(testPet.Status))
			gomega.Expect(pet.Tags).Should(gomega.HaveLen(2))
		})

	})

	ginkgo.Describe("FindPetById", func() {
		ginkgo.It("should successfully find a pet", func() {
			// Call the service method
			pet, err := service.FindPetById(ctx, &proto_storage.PetID{Id: 1})

			// Assertions
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(pet).ShouldNot(gomega.BeNil())
			gomega.Expect(pet.Name).Should(gomega.Equal("Test Pet"))
			gomega.Expect(pet.Category.Name).Should(gomega.Equal("Test Category"))
			gomega.Expect(pet.PhotoUrls).Should(gomega.HaveLen(2))
			gomega.Expect(pet.Tags).Should(gomega.HaveLen(2))
		})
	})

	ginkgo.Describe("FindPets", func() {
		ginkgo.It("should successfully find pets", func() {
			// Call the service method
			pets, err := service.FindPets(ctx, &proto_storage.FindPetsRequest{})

			// Assertions
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(pets).ShouldNot(gomega.BeNil())
			gomega.Expect(pets.Pets).ShouldNot(gomega.BeEmpty())
		})
	})

	ginkgo.Describe("DeletePet", func() {
		ginkgo.It("should successfully delete a pet", func() {
			// Call the service method
			_, err := service.DeletePet(ctx, &proto_storage.PetID{Id: 1})

			// Assertions
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
		})
	})

})
