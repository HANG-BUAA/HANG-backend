package custom_error

// OptimisticLockError 乐观锁冲突错误
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
