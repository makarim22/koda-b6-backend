package errors

import "fmt"

// Domain error types for consistent error handling across services and repositories

// NotFoundError represents a resource that doesn't exist
type NotFoundError struct {
	Resource string
	ID       interface{}
	Message  string
}

func (e *NotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("%s with ID %v not found", e.Resource, e.ID)
}

// ValidationError represents invalid input
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// ConflictError represents a resource that already exists or conflicts
type ConflictError struct {
	Resource string
	Message  string
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// InsufficientStockError represents when a product doesn't have enough stock
type InsufficientStockError struct {
	ProductID   int
	Available   int
	Requested   int
	ProductName string
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf("insufficient stock for product '%s': available %d, requested %d",
		e.ProductName, e.Available, e.Requested)
}

// InvalidStatusError represents an invalid status transition or value
type InvalidStatusError struct {
	CurrentStatus   string
	RequestedStatus string
	ResourceType    string
}

func (e *InvalidStatusError) Error() string {
	return fmt.Sprintf("invalid status transition for %s: cannot change from '%s' to '%s'",
		e.ResourceType, e.CurrentStatus, e.RequestedStatus)
}

// DatabaseError represents a database operation failure
type DatabaseError struct {
	Operation string
	Details   string
	Err       error
}

func (e *DatabaseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("database error during %s: %v", e.Operation, e.Err)
	}
	return fmt.Sprintf("database error during %s: %s", e.Operation, e.Details)
}

// Unwrap allows error wrapping/chaining
func (e *DatabaseError) Unwrap() error {
	return e.Err
}

// Constructor functions for easier error creation

// NewNotFoundError creates a new NotFoundError
func NewNotFoundError(resource string, id interface{}) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// NewNotFoundErrorWithMessage creates a NotFoundError with a custom message
func NewNotFoundErrorWithMessage(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}

// NewValidationError creates a new ValidationError
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// NewConflictError creates a new ConflictError
func NewConflictError(resource, message string) *ConflictError {
	return &ConflictError{
		Resource: resource,
		Message:  message,
	}
}

// NewInsufficientStockError creates a new InsufficientStockError
func NewInsufficientStockError(productID, available, requested int, productName string) *InsufficientStockError {
	return &InsufficientStockError{
		ProductID:   productID,
		Available:   available,
		Requested:   requested,
		ProductName: productName,
	}
}

// NewInvalidStatusError creates a new InvalidStatusError
func NewInvalidStatusError(current, requested, resourceType string) *InvalidStatusError {
	return &InvalidStatusError{
		CurrentStatus:   current,
		RequestedStatus: requested,
		ResourceType:    resourceType,
	}
}

// NewDatabaseError creates a new DatabaseError
func NewDatabaseError(operation string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}

// Type assertion helpers for checking error types

// IsNotFoundError checks if an error is a NotFoundError
func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}

// IsValidationError checks if an error is a ValidationError
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}

// IsConflictError checks if an error is a ConflictError
func IsConflictError(err error) bool {
	_, ok := err.(*ConflictError)
	return ok
}

// IsInsufficientStockError checks if an error is an InsufficientStockError
func IsInsufficientStockError(err error) bool {
	_, ok := err.(*InsufficientStockError)
	return ok
}

// IsInvalidStatusError checks if an error is an InvalidStatusError
func IsInvalidStatusError(err error) bool {
	_, ok := err.(*InvalidStatusError)
	return ok
}

// IsDatabaseError checks if an error is a DatabaseError
func IsDatabaseError(err error) bool {
	_, ok := err.(*DatabaseError)
	return ok
}

/// service level failure

// ServiceError represents a service-level operation failure
type ServiceError struct {
	Operation string
	Err       error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("service error during %s: %v", e.Operation, e.Err)
	}
	return fmt.Sprintf("service error during %s", e.Operation)
}

// Unwrap allows error wrapping/chaining
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError creates a new ServiceError
func NewServiceError(operation string, err error) *ServiceError {
	return &ServiceError{
		Operation: operation,
		Err:       err,
	}
}
