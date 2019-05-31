package arbiter

// Arbiter is multi mutex manager for exclusive access.
type Arbiter interface {
	Lock(key string)
	Unlock(key string)
}
