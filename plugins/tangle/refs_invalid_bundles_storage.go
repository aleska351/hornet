package tangle

import (
	"time"

	"github.com/iotaledger/iota.go/trinary"

	"github.com/iotaledger/hive.go/objectstorage"

	"github.com/gohornet/hornet/packages/profile"
)

var (
	refsAnInvalidBundleStorage *objectstorage.ObjectStorage
)

type invalidBundleReference struct {
	objectstorage.StorableObjectFlags

	hashBytes []byte
}

// ObjectStorage interface

func (r *invalidBundleReference) Update(other objectstorage.StorableObject) {
	panic("invalidBundleReference should never be updated")
}

func (r *invalidBundleReference) GetStorageKey() []byte {
	return r.hashBytes
}

func (r *invalidBundleReference) MarshalBinary() (data []byte, err error) {
	return nil, nil
}

func (r *invalidBundleReference) UnmarshalBinary(data []byte) error {
	return nil
}

func invalidBundleFactory(key []byte) objectstorage.StorableObject {
	invalidBndl := &invalidBundleReference{
		hashBytes: make([]byte, len(key)),
	}
	copy(invalidBndl.hashBytes, key)
	return invalidBndl
}

func configureRefsAnInvalidBundleStorage() {
	opts := profile.GetProfile().Caches.RefsInvalidBundle

	refsAnInvalidBundleStorage = objectstorage.New(
		nil,
		nil,
		invalidBundleFactory,
		objectstorage.CacheTime(time.Duration(opts.CacheTimeMs)*time.Millisecond),
		objectstorage.PersistenceEnabled(false),
		objectstorage.LeakDetectionEnabled(opts.LeakDetectionOptions.Enabled,
			objectstorage.LeakDetectionOptions{
				MaxConsumersPerObject: opts.LeakDetectionOptions.MaxConsumersPerObject,
				MaxConsumerHoldTime:   time.Duration(opts.LeakDetectionOptions.MaxConsumerHoldTimeSec) * time.Second,
			}),
	)
}

func GetRefsAnInvalidBundleStorageSize() int {
	return refsAnInvalidBundleStorage.GetSize()
}

// +-0
func PutInvalidBundleReference(txHash trinary.Hash) {
	invalidBundleRef := invalidBundleFactory(trinary.MustTrytesToBytes(txHash)[:49])

	refsAnInvalidBundleStorage.ComputeIfAbsent(invalidBundleRef.GetStorageKey(), func(key []byte) objectstorage.StorableObject {
		invalidBundleRef.SetModified()
		return invalidBundleRef
	}).Release()
}

// +-0
func ContainsInvalidBundleReference(txHash trinary.Hash) bool {
	return refsAnInvalidBundleStorage.Contains(trinary.MustTrytesToBytes(txHash)[:49])
}