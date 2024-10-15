// Code generated by ObjectBox; DO NOT EDIT.
// Learn more about defining entities and generating this file - visit https://golang.objectbox.io/entity-annotations

package models

import (
	"errors"
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/objectbox/fbutils"
)

type encoderModel_EntityInfo struct {
	objectbox.Entity
	Uid uint64
}

var EncoderModelBinding = encoderModel_EntityInfo{
	Entity: objectbox.Entity{
		Id: 1,
	},
	Uid: 3730907597553742834,
}

// EncoderModel_ contains type-based Property helpers to facilitate some common operations such as Queries.
var EncoderModel_ = struct {
	Id          *objectbox.PropertyUint64
	EncoderType *objectbox.PropertyString
	SourceID    *objectbox.PropertyString
}{
	Id: &objectbox.PropertyUint64{
		BaseProperty: &objectbox.BaseProperty{
			Id:     1,
			Entity: &EncoderModelBinding.Entity,
		},
	},
	EncoderType: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     2,
			Entity: &EncoderModelBinding.Entity,
		},
	},
	SourceID: &objectbox.PropertyString{
		BaseProperty: &objectbox.BaseProperty{
			Id:     3,
			Entity: &EncoderModelBinding.Entity,
		},
	},
}

// GeneratorVersion is called by ObjectBox to verify the compatibility of the generator used to generate this code
func (encoderModel_EntityInfo) GeneratorVersion() int {
	return 6
}

// AddToModel is called by ObjectBox during model build
func (encoderModel_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("EncoderModel", 1, 3730907597553742834)
	model.Property("Id", 6, 1, 1439176170655404775)
	model.PropertyFlags(1)
	model.Property("EncoderType", 9, 2, 3793273354329931500)
	model.Property("SourceID", 9, 3, 4736056337548406699)
	model.EntityLastPropertyId(3, 4736056337548406699)
}

// GetId is called by ObjectBox during Put operations to check for existing ID on an object
func (encoderModel_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*EncoderModel).Id, nil
}

// SetId is called by ObjectBox during Put to update an ID on an object that has just been inserted
func (encoderModel_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*EncoderModel).Id = id
	return nil
}

// PutRelated is called by ObjectBox to put related entities before the object itself is flattened and put
func (encoderModel_EntityInfo) PutRelated(ob *objectbox.ObjectBox, object interface{}, id uint64) error {
	return nil
}

// Flatten is called by ObjectBox to transform an object to a FlatBuffer
func (encoderModel_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) error {
	obj := object.(*EncoderModel)
	var offsetEncoderType = fbutils.CreateStringOffset(fbb, obj.EncoderType)
	var offsetSourceID = fbutils.CreateStringOffset(fbb, obj.SourceID)

	// build the FlatBuffers object
	fbb.StartObject(3)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetUOffsetTSlot(fbb, 1, offsetEncoderType)
	fbutils.SetUOffsetTSlot(fbb, 2, offsetSourceID)
	return nil
}

// Load is called by ObjectBox to load an object from a FlatBuffer
func (encoderModel_EntityInfo) Load(ob *objectbox.ObjectBox, bytes []byte) (interface{}, error) {
	if len(bytes) == 0 { // sanity check, should "never" happen
		return nil, errors.New("can't deserialize an object of type 'EncoderModel' - no data received")
	}

	var table = &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	var propId = table.GetUint64Slot(4, 0)

	return &EncoderModel{
		Id:          propId,
		EncoderType: fbutils.GetStringSlot(table, 6),
		SourceID:    fbutils.GetStringSlot(table, 8),
	}, nil
}

// MakeSlice is called by ObjectBox to construct a new slice to hold the read objects
func (encoderModel_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*EncoderModel, 0, capacity)
}

// AppendToSlice is called by ObjectBox to fill the slice of the read objects
func (encoderModel_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	if object == nil {
		return append(slice.([]*EncoderModel), nil)
	}
	return append(slice.([]*EncoderModel), object.(*EncoderModel))
}

// Box provides CRUD access to EncoderModel objects
type EncoderModelBox struct {
	*objectbox.Box
}

// BoxForEncoderModel opens a box of EncoderModel objects
func BoxForEncoderModel(ob *objectbox.ObjectBox) *EncoderModelBox {
	return &EncoderModelBox{
		Box: ob.InternalBox(1),
	}
}

// Put synchronously inserts/updates a single object.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the EncoderModel.Id property on the passed object will be assigned the new ID as well.
func (box *EncoderModelBox) Put(object *EncoderModel) (uint64, error) {
	return box.Box.Put(object)
}

// Insert synchronously inserts a single object. As opposed to Put, Insert will fail if given an ID that already exists.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the EncoderModel.Id property on the passed object will be assigned the new ID as well.
func (box *EncoderModelBox) Insert(object *EncoderModel) (uint64, error) {
	return box.Box.Insert(object)
}

// Update synchronously updates a single object.
// As opposed to Put, Update will fail if an object with the same ID is not found in the database.
func (box *EncoderModelBox) Update(object *EncoderModel) error {
	return box.Box.Update(object)
}

// PutAsync asynchronously inserts/updates a single object.
// Deprecated: use box.Async().Put() instead
func (box *EncoderModelBox) PutAsync(object *EncoderModel) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutMany inserts multiple objects in single transaction.
// In case Ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the EncoderModel.Id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the EncoderModel.Id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *EncoderModelBox) PutMany(objects []*EncoderModel) ([]uint64, error) {
	return box.Box.PutMany(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *EncoderModelBox) Get(id uint64) (*EncoderModel, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*EncoderModel), nil
}

// GetMany reads multiple objects at once.
// If any of the objects doesn't exist, its position in the return slice is nil
func (box *EncoderModelBox) GetMany(ids ...uint64) ([]*EncoderModel, error) {
	objects, err := box.Box.GetMany(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*EncoderModel), nil
}

// GetManyExisting reads multiple objects at once, skipping those that do not exist.
func (box *EncoderModelBox) GetManyExisting(ids ...uint64) ([]*EncoderModel, error) {
	objects, err := box.Box.GetManyExisting(ids...)
	if err != nil {
		return nil, err
	}
	return objects.([]*EncoderModel), nil
}

// GetAll reads all stored objects
func (box *EncoderModelBox) GetAll() ([]*EncoderModel, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*EncoderModel), nil
}

// Remove deletes a single object
func (box *EncoderModelBox) Remove(object *EncoderModel) error {
	return box.Box.Remove(object)
}

// RemoveMany deletes multiple objects at once.
// Returns the number of deleted object or error on failure.
// Note that this method will not fail if an object is not found (e.g. already removed).
// In case you need to strictly check whether all of the objects exist before removing them,
// you can execute multiple box.Contains() and box.Remove() inside a single write transaction.
func (box *EncoderModelBox) RemoveMany(objects ...*EncoderModel) (uint64, error) {
	var ids = make([]uint64, len(objects))
	for k, object := range objects {
		ids[k] = object.Id
	}
	return box.Box.RemoveIds(ids...)
}

// Creates a query with the given conditions. Use the fields of the EncoderModel_ struct to create conditions.
// Keep the *EncoderModelQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *EncoderModelBox) Query(conditions ...objectbox.Condition) *EncoderModelQuery {
	return &EncoderModelQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the EncoderModel_ struct to create conditions.
// Keep the *EncoderModelQuery if you intend to execute the query multiple times.
func (box *EncoderModelBox) QueryOrError(conditions ...objectbox.Condition) (*EncoderModelQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &EncoderModelQuery{query}, nil
	}
}

// Async provides access to the default Async Box for asynchronous operations. See EncoderModelAsyncBox for more information.
func (box *EncoderModelBox) Async() *EncoderModelAsyncBox {
	return &EncoderModelAsyncBox{AsyncBox: box.Box.Async()}
}

// EncoderModelAsyncBox provides asynchronous operations on EncoderModel objects.
//
// Asynchronous operations are executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "execute & forget:" you gain faster put/remove operations as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
// In situations with (extremely) high async load, an async method may be throttled (~1ms) or delayed up to 1 second.
// In the unlikely event that the object could still not be enqueued (full queue), an error will be returned.
//
// Note that async methods do not give you hard durability guarantees like the synchronous Box provides.
// There is a small time window in which the data may not have been committed durably yet.
type EncoderModelAsyncBox struct {
	*objectbox.AsyncBox
}

// AsyncBoxForEncoderModel creates a new async box with the given operation timeout in case an async queue is full.
// The returned struct must be freed explicitly using the Close() method.
// It's usually preferable to use EncoderModelBox::Async() which takes care of resource management and doesn't require closing.
func AsyncBoxForEncoderModel(ob *objectbox.ObjectBox, timeoutMs uint64) *EncoderModelAsyncBox {
	var async, err = objectbox.NewAsyncBox(ob, 1, timeoutMs)
	if err != nil {
		panic("Could not create async box for entity ID 1: %s" + err.Error())
	}
	return &EncoderModelAsyncBox{AsyncBox: async}
}

// Put inserts/updates a single object asynchronously.
// When inserting a new object, the Id property on the passed object will be assigned the new ID the entity would hold
// if the insert is ultimately successful. The newly assigned ID may not become valid if the insert fails.
func (asyncBox *EncoderModelAsyncBox) Put(object *EncoderModel) (uint64, error) {
	return asyncBox.AsyncBox.Put(object)
}

// Insert a single object asynchronously.
// The Id property on the passed object will be assigned the new ID the entity would hold if the insert is ultimately
// successful. The newly assigned ID may not become valid if the insert fails.
// Fails silently if an object with the same ID already exists (this error is not returned).
func (asyncBox *EncoderModelAsyncBox) Insert(object *EncoderModel) (id uint64, err error) {
	return asyncBox.AsyncBox.Insert(object)
}

// Update a single object asynchronously.
// The object must already exists or the update fails silently (without an error returned).
func (asyncBox *EncoderModelAsyncBox) Update(object *EncoderModel) error {
	return asyncBox.AsyncBox.Update(object)
}

// Remove deletes a single object asynchronously.
func (asyncBox *EncoderModelAsyncBox) Remove(object *EncoderModel) error {
	return asyncBox.AsyncBox.Remove(object)
}

// Query provides a way to search stored objects
//
// For example, you can find all EncoderModel which Id is either 42 or 47:
//
// box.Query(EncoderModel_.Id.In(42, 47)).Find()
type EncoderModelQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *EncoderModelQuery) Find() ([]*EncoderModel, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*EncoderModel), nil
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *EncoderModelQuery) Offset(offset uint64) *EncoderModelQuery {
	query.Query.Offset(offset)
	return query
}

// Limit sets the number of elements to process by the query
func (query *EncoderModelQuery) Limit(limit uint64) *EncoderModelQuery {
	query.Query.Limit(limit)
	return query
}
