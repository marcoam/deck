package state

import (
	"testing"

	"github.com/hbagdi/go-kong/kong"
	"github.com/stretchr/testify/assert"
)

func TestTargetInsert(t *testing.T) {
	assert := assert.New(t)
	collection, err := NewTargetsCollection()
	assert.Nil(err)
	assert.NotNil(collection)
	var target Target
	target.Target.Target = kong.String("my-target")
	target.ID = kong.String("first")
	err = collection.Add(target)
	assert.NotNil(err)

	var target2 Target
	target2.Target.Target = kong.String("my-target")
	target2.ID = kong.String("first")
	target2.Upstream = &kong.Upstream{
		ID:   kong.String("upstream1-id"),
		Name: kong.String("upstream1-name"),
	}
	assert.NotNil(target2.Upstream)
	err = collection.Add(target2)
	assert.NotNil(target2.Upstream)
	assert.Nil(err)
}

func TestTargetGetUpdate(t *testing.T) {
	assert := assert.New(t)
	collection, err := NewTargetsCollection()
	assert.Nil(err)
	assert.NotNil(collection)
	var target Target
	target.Target.Target = kong.String("my-target")
	target.ID = kong.String("first")
	target.Upstream = &kong.Upstream{
		ID:   kong.String("upstream1-id"),
		Name: kong.String("upstream1-name"),
	}
	assert.NotNil(target.Upstream)
	err = collection.Add(target)
	assert.NotNil(target.Upstream)
	assert.Nil(err)

	re, err := collection.Get("first")
	assert.Nil(err)
	assert.NotNil(re)
	assert.Equal("my-target", *re.Target.Target)
	err = collection.Update(*re)
	assert.Nil(err)

	re, err = collection.Get("my-target")
	assert.Nil(err)
	assert.NotNil(re)
}

func TestTargetsInvalidType(t *testing.T) {
	assert := assert.New(t)

	collection, err := NewTargetsCollection()
	assert.Nil(err)
	assert.NotNil(collection)

	var upstream Upstream
	upstream.Name = kong.String("my-upstream")
	upstream.ID = kong.String("first")
	txn := collection.memdb.Txn(true)
	err = txn.Insert(targetTableName, &upstream)
	assert.NotNil(err)
	txn.Abort()

	type badTarget struct {
		kong.Target
		Meta
	}

	target := badTarget{
		Target: kong.Target{
			ID:     kong.String("id"),
			Target: kong.String("target"),
			Upstream: &kong.Upstream{
				ID:   kong.String("upstream-id"),
				Name: kong.String("upstream-name"),
			},
		},
	}

	txn = collection.memdb.Txn(true)
	err = txn.Insert(targetTableName, &target)
	assert.Nil(err)
	txn.Commit()

	assert.Panics(func() {
		collection.Get("target")
	})
}

func TestTargetDelete(t *testing.T) {
	assert := assert.New(t)
	collection, err := NewTargetsCollection()
	assert.Nil(err)
	assert.NotNil(collection)

	var target Target
	target.Target.Target = kong.String("my-target")
	target.ID = kong.String("first")
	target.Upstream = &kong.Upstream{
		ID:   kong.String("upstream1-id"),
		Name: kong.String("upstream1-name"),
	}
	err = collection.Add(target)
	assert.Nil(err)

	re, err := collection.Get("my-target")
	assert.Nil(err)
	assert.NotNil(re)

	err = collection.Delete(*re.ID)
	assert.Nil(err)

	err = collection.Delete(*re.ID)
	assert.NotNil(err)
}

func TestTargetGetAll(t *testing.T) {
	assert := assert.New(t)
	collection, err := NewTargetsCollection()
	assert.Nil(err)
	assert.NotNil(collection)

	var target Target
	target.Target.Target = kong.String("my-target1")
	target.ID = kong.String("first")
	target.Upstream = &kong.Upstream{
		ID:   kong.String("upstream1-id"),
		Name: kong.String("upstream1-name"),
	}
	err = collection.Add(target)
	assert.Nil(err)

	var target2 Target
	target2.Target.Target = kong.String("my-target2")
	target2.ID = kong.String("second")
	target2.Upstream = &kong.Upstream{
		ID:   kong.String("upstream1-id"),
		Name: kong.String("upstream1-name"),
	}
	err = collection.Add(target2)
	assert.Nil(err)

	targets, err := collection.GetAll()

	assert.Nil(err)
	assert.Equal(2, len(targets))
}

func TestTargetGetAllByUpstreamName(t *testing.T) {
	assert := assert.New(t)
	collection, err := NewTargetsCollection()
	assert.Nil(err)
	assert.NotNil(collection)

	targets := []*Target{
		{
			Target: kong.Target{
				ID:     kong.String("target1-id"),
				Target: kong.String("target1-name"),
				Upstream: &kong.Upstream{
					ID:   kong.String("upstream1-id"),
					Name: kong.String("upstream1-name"),
				},
			},
		},
		{
			Target: kong.Target{
				ID:     kong.String("target2-id"),
				Target: kong.String("target2-name"),
				Upstream: &kong.Upstream{
					ID:   kong.String("upstream1-id"),
					Name: kong.String("upstream1-name"),
				},
			},
		},
		{
			Target: kong.Target{
				ID:     kong.String("target3-id"),
				Target: kong.String("target3-name"),
				Upstream: &kong.Upstream{
					ID:   kong.String("upstream2-id"),
					Name: kong.String("upstream2-name"),
				},
			},
		},
		{
			Target: kong.Target{
				ID:     kong.String("target4-id"),
				Target: kong.String("target4-name"),
				Upstream: &kong.Upstream{
					ID:   kong.String("upstream2-id"),
					Name: kong.String("upstream2-name"),
				},
			},
		},
	}

	for _, target := range targets {
		err = collection.Add(*target)
		assert.Nil(err)
	}

	targets, err = collection.GetAllByUpstreamID("upstream1-id")
	assert.Nil(err)
	assert.Equal(2, len(targets))

	targets, err = collection.GetAllByUpstreamName("upstream2-name")
	assert.Nil(err)
	assert.Equal(2, len(targets))

	targets, err = collection.GetAllByUpstreamName("upstream1-id")
	assert.Nil(err)
	assert.Equal(0, len(targets))
}
