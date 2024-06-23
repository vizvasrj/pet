package storageservice_test

import (
	"context"
	"src/env"
	"src/petstore"
	"src/storageservice"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCreatePet(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CreatePet Suite")
}

var _ = Describe("CreatePet", func() {
	var (
		s storageservice.StorageService
	)

	BeforeEach(func() {
		e := env.GetEnvs()
		db := storageservice.GetConnection(e)
		// defer db.Close()

		s = storageservice.StorageService{
			Db: db,
		}
	})

	AfterEach(func() {
		s.Db.Close()
	})

	It("creates a pet", func() {
		cId := int64(1)
		cName := "Cat"
		tagId := int64(0)
		tagName := "meao tag 4"
		tag := petstore.Tag{
			Id:   &tagId,
			Name: &tagName,
		}

		tag2Name := "meao tag 5"
		tag2 := petstore.Tag{
			Id:   &tagId,
			Name: &tag2Name,
		}

		tags := []petstore.Tag{}
		tags = append(tags, tag)
		tags = append(tags, tag2)

		_, err := s.CreatePet(context.Background(), &petstore.NewPet{
			Name: "cat meao 5",
			Category: &petstore.Category{
				Id:   &cId,
				Name: &cName,
			},
			PhotoUrls: &[]string{"mearourl_string"},
			Tags:      &tags,
		})

		Expect(err).NotTo(HaveOccurred())
	})
})
