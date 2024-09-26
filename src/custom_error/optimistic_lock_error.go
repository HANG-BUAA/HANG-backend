package custom_error

type OptimisticLockError struct {
}

func (e *OptimisticLockError) Error() string {
	return "optimistic lock constrained custom_error"
}

func NewOptimisticLockError() error {
	return &OptimisticLockError{}
}

func (e *OptimisticLockError) Is(target error) bool {
	_, ok := target.(*OptimisticLockError)
	return ok
}
