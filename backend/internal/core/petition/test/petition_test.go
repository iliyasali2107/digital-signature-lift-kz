package petition_test

import (
	"context"
	"os"
	"testing"

	"go.uber.org/zap"

	"mado/internal/core/petition"
)

// Mock repository for testing
type mockRepository struct{}

// GetPetitionPdfByID implements petition.Repository.
func (*mockRepository) GetPetitionPdfByID(ctx context.Context, pdfID *int) (*petition.PetitionData, error) {
	return &petition.PetitionData{PdfData: []byte("vfdvfd"), FileName: "file1.pdf"}, nil
}

// GetNextID implements petition.Repository.
func (*mockRepository) GetNextID(ctx context.Context) (*int, error) {
	id := 123
	return &id, nil
}

func (m *mockRepository) Save(ctx context.Context, dto *petition.PetitionData) (*petition.PetitionData, error) {
	return dto, nil
}

func TestGeneratePetitionPDF(t *testing.T) {

	// Create a mock repository for testing
	mockRepo := &mockRepository{}

	// Create a mock logger for testing
	mockLogger, _ := zap.NewDevelopment()

	// Create a new service using the mock repository and logger
	service := petition.NewService(mockRepo, mockLogger)

	// Create test data
	// id := 123
	testData := &petition.PetitionData{
		FileName: "output.pdf",
		// SheetNumber:       &id,
		// CreationDate:      "01 September 2023",
		Location:          "Apartment Building 123",
		ResponsiblePerson: "John Doe",
		Questions: []petition.Question{
			{Number: 1, Text: "Should we repaint the common areas?", Decision: "За"},
			{Number: 2, Text: "Should we install security cameras?", Decision: "Воздержусь"},
			{Number: 3, Text: "Should we increase the maintenance fee?", Decision: "Против"},
		},
		OwnerName:    "Alice Smith",
		OwnerAddress: "vfv4d6515",
	}

	// Call the function being tested
	_, err := service.GeneratePetitionPDF(testData)
	if err != nil {
		t.Errorf("Error generating PDF: %v", err)
	}

	// Check if the PDF file was created
	_, err = os.Stat(testData.FileName)
	if os.IsNotExist(err) {
		t.Error("PDF file was not created")
	}

	// Clean up: remove the generated PDF file
	// err = os.Remove(testData.FileName)
	// if err != nil {
	// 	t.Errorf("Error deleting PDF file: %v", err)
	// }

	// err = os.Remove("temp.html")
	// if err != nil {
	// 	t.Errorf("Error deleting html file: %v", err)
	// }
}

func NewService(mockRepo *mockRepository, mockLogger *zap.Logger) {
	panic("unimplemented")
}

/**
testing json:

{
    "file_name": "output.pdf",
    "location": "Apartment Building 123",
    "responsible_person": "John Doe",
    "questions": [
        {
            "number": 1,
            "text": "Should we repaint the common areas?",
            "description": "За"
        },
        {
            "number": 2,
            "text": "Should we install security cameras?",
            "description": "Воздержусь"
        },
        {
            "number": 3,
            "text": "Should we increase the maintenance fee?",
            "description": "Против"
        }
    ],
    "owner_name": "Alice Smith",
    "owner_address": "vfv4d6515"
}

*

*/
