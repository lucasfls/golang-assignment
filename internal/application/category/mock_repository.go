package category

import (
	"context"

	"github.com/mytheresa/go-hiring-challenge/internal/domain/category"
)

// MockRepository is a mock implementation of category.Repository for testing.
type MockRepository struct {
	findAllFunc func(ctx context.Context) ([]category.Category, error)
	createFunc  func(ctx context.Context, code, name string) (*category.Category, error)
}

// NewMockRepository creates a new MockRepository.
func NewMockRepository() *MockRepository {
	return &MockRepository{}
}

// WithFindAll sets the mock behavior for FindAll.
func (m *MockRepository) WithFindAll(fn func(ctx context.Context) ([]category.Category, error)) *MockRepository {
	m.findAllFunc = fn
	return m
}

// WithCreate sets the mock behavior for Create.
func (m *MockRepository) WithCreate(fn func(ctx context.Context, code, name string) (*category.Category, error)) *MockRepository {
	m.createFunc = fn
	return m
}

// FindAll implements category.Repository.
func (m *MockRepository) FindAll(ctx context.Context) ([]category.Category, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return []category.Category{}, nil
}

// Create implements category.Repository.
func (m *MockRepository) Create(ctx context.Context, code, name string) (*category.Category, error) {
	if m.createFunc != nil {
		return m.createFunc(ctx, code, name)
	}
	return &category.Category{ID: 1, Code: code, Name: name}, nil
}
