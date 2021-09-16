package database

import (
	"math"
	"reflect"

	"github.com/zhangsifeng92/geos/log"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type SkipSuffix int

const (
	_ SkipSuffix = iota
	SKIP_ONE
	SKIP_TWO
)

type LDataBase struct {
	db        *leveldb.DB
	stack     *deque
	path      string
	reversion int64
	nextId    map[string]int64
	logFlag   bool
	log       log.Logger
	batch     *leveldb.Batch
	count     int64
	isClosed  bool
}

/*
*	Create a database based on the provided path
*	successfully return the database handle
*	otherwise return error message
 */

func NewDataBase(path string, flag ...bool) (DataBase, error) {

	db, err := leveldb.OpenFile(path, &opt.Options{
		OpenFilesCacheCapacity: 16,
		BlockCacheCapacity:     16 / 2 * opt.MiB,
		WriteBuffer:            16 / 4 * opt.MiB, // Two of these are used internally
		Filter:                 filter.NewBloomFilter(10),
	})

	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		return nil, err
	}
	/*	read every type increment	*/
	nextId, err := readIncrementFromDb(db)
	if err != nil {
		log.Error("database init failed : %s", err.Error())
		panic("open database file failed : " + err.Error())
	}
	/* read reversion */
	reversion := readReversionFromDb(db)
	/* read stack */
	stack := readUndoStackFromDb(db)
	logFlag := false
	if len(flag) > 0 {
		logFlag = flag[0]
	}

	dbLog := log.New("db")
	if logFlag {
		h, _ := log.FileHandler(path+"/database_log.log", log.LogfmtFormat())
		dbLog.SetHandler(h)
	} else {
		dbLog.SetHandler(log.DiscardHandler())
	}
	return &LDataBase{db: db, stack: stack, path: path, nextId: nextId, logFlag: logFlag, log: dbLog, batch: new(leveldb.Batch), reversion: reversion}, nil
}
func readUndoStackFromDb(db *leveldb.DB) *deque {
	key := []byte(undoKey)
	val, err := db.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		//throw

	}
	dq := newDeque()
	if err == leveldb.ErrNotFound {
		return dq
	}
	values := [][]byte{}
	err = DecodeBytes(val, &values)
	if err != nil {
		// throw
	}
	cons := []undoContainer{}

	for index, _ := range values {
		con := undoContainer{}
		err = DecodeBytes(values[index], &con)
		if err != nil {
			//throw
		}
		cons = append(cons, con)
	}

	for index, _ := range cons {
		dq.Append(&cons[index])
	}

	return dq
}
func readReversionFromDb(db *leveldb.DB) int64 {
	key := []byte(dbReversion)
	val, err := db.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		//throw

	}
	if err == leveldb.ErrNotFound {
		return 0
	}
	var reversion int64
	err = DecodeBytes(val, &reversion)
	if err != nil {
		// throw
	}
	return reversion
}
func readIncrementFromDb(db *leveldb.DB) (map[string]int64, error) {
	nextId := make(map[string]int64)

	key := []byte(dbIncrement)
	val, err := db.Get(key, nil)
	if err != nil && err != leveldb.ErrNotFound {
		return nil, err
	}
	if err == nil && len(val) != 0 {
		err = DecodeBytes(val, &nextId)
		if err != nil {
			panic("database init failed : " + err.Error())
		}
	}
	err = db.Delete(key, nil)
	if err != nil {
		panic("database init failed : " + err.Error())
	}
	return nextId, nil
}

func (ldb *LDataBase) Close() {
	if ldb.isClosed {
		return
	}
	err := ldb.writeIncrementToDb()
	if err != nil {
		ldb.log.Error("database close failed : %s", err.Error())
	}
	err = ldb.db.Close()
	if err != nil {
		/* throw ? */
		ldb.log.Error("database close failed : %s", err.Error())
	} else {
		ldb.log.Info("----------------- database close -----------------")
	}
	ldb.isClosed = true
}

func (ldb *LDataBase) writeIncrementToDb() error {

	if len(ldb.nextId) == 0 {
		return nil
	}
	// write next id
	val, err := EncodeToBytes(ldb.nextId)
	if err != nil {
		ldb.log.Error("WriteIncrement rlp EncodeToBytes failed is : %s", err.Error())
		return err
	}
	err = ldb.db.Put([]byte(dbIncrement), val, nil)
	if err != nil {
		ldb.log.Error("WriteIncrement saveKey failed is : %s", err.Error())
		return err
	}
	// write reversion
	val, err = EncodeToBytes(ldb.reversion)
	if err != nil {
		return err
	}
	err = ldb.db.Put([]byte(dbReversion), val, nil)
	if err != nil {
		return err
	}
	// write undo stack
	ldb.writeUndoStack()
	return nil
}

func (ldb *LDataBase) undoStackToByte() []byte {
	u := ldb.getStack()
	uv := undoContainer{}
	uv = *u
	val, err := EncodeToBytes(uv)
	if err != nil {
		return nil
	}
	return val
}
func (ldb *LDataBase) writeUndoStack() {
	values := [][]byte{}
	for ldb.stack.Size() > 0 {
		val := ldb.undoStackToByte()
		values = append(values, val)
		ldb.stack.Pop()
	}
	if len(values) != 0 {
		val, err := EncodeToBytes(values)
		if err != nil {
			// throw
		}
		err = ldb.db.Put([]byte(undoKey), val, nil)
		if err != nil {
			//throw
		}
	}
}
func (ldb *LDataBase) Revision() int64 {
	ldb.log.Info("ldb reversion is : %d", ldb.reversion)
	return ldb.reversion
}

func (ldb *LDataBase) Undo() {

	stack := ldb.getStack()
	if stack == nil {
		return
	}

	ldb.nextId = stack.OldIds

	for _, undo := range stack.Undo {

		for _, value := range undo.OldValue { /* 	Undo edit */
			err := ldb.modifyKvToDb(value.NewKv, value.OldKv)
			if err != nil {
				//throw
			}
		}

		for _, value := range undo.NewValue { /*	Undo new */
			err := ldb.removeKvToDb(value.NewKv)
			if err != nil {
				//throw
			}
		}

		for _, value := range undo.RemoveValue { /*	Undo remove */
			err := ldb.insertKvToDb(value.NewKv)
			if err != nil {
				//throw
			}
		}
	}

	ldb.stack.Pop()
	ldb.reversion--
}

func (ldb *LDataBase) UndoAll() {
	// TODO wait Undo
	for ldb.enable() {
		ldb.Undo()
	}
}

func (ldb *LDataBase) Squash() {

	if ldb.stack.Size() == 1 {
		ldb.stack.Pop()
		ldb.reversion--
		return
	}
	stack := ldb.getStack()
	preStack := ldb.getSecond()

	for typeName, undo := range stack.Undo {

		if _, ok := preStack.Undo[typeName]; !ok {
			preStack.Undo[typeName] = newUndoState()
		}
		//TODO assert typeName
		undoStateSquash(undo, preStack.Undo[typeName])
	}

	ldb.stack.Pop()
	ldb.reversion--
}

func (ldb *LDataBase) StartSession() *Session {
	ldb.reversion++
	state := newUndoContainer(ldb.reversion, ldb.nextId)
	ldb.stack.Append(state)
	return &Session{db: ldb, apply: true, reversion: ldb.reversion}
}

func (ldb *LDataBase) Commit(reversion int64) {

	for {
		if ldb.stack.Size() == 0 {
			ldb.log.Info("database Commit stack is empty")
			return
		}

		stack := ldb.getFirstStack()
		if stack == nil {
			ldb.log.Info("database Commit stack is nil")
			break
		}
		ldb.log.Info("database Commit stack reversion is : %d", stack.Reversion)
		if stack.Reversion > reversion {
			break
		}

		ldb.stack.PopFront()
	}
}

func (ldb *LDataBase) SetRevision(reversion int64) {
	if ldb.stack.Size() != 0 {
		panic("cannot set reversion while there is an existing undo stack")
		// throw
	}
	if reversion > math.MaxInt64 {
		//throw
		panic("reversion to set is too high")
	}
	ldb.log.Info("database SetRevision  reversion is : %d", reversion)
	ldb.reversion = reversion
}

/*
*	Insert a piece of data into the database
*	returning null successfully
*	returns an error message
 */

func (ldb *LDataBase) Insert(in interface{}) error {
	err := ldb.insert(in)
	if err != nil {
		ldb.log.Error("error database insert failed : %s", err.Error())
		return err
	}
	return nil
}

func (ldb *LDataBase) insert(in interface{}, flag ...bool) error { /* struct cfg --> KV struct --> kv to db --> undo db */
	cfg, err := parseObjectToCfg(in) /* (struct cfg) parse object tag */
	if err != nil {
		ldb.log.Error("error database insert  parseObjectToCfg failed : %s", err.Error())
		return err
	}
	undoFlag := false
	if len(flag) > 0 {
		undoFlag = flag[0]
	}

	if !undoFlag {
		err = ldb.setIncrement(cfg) /* (kv.id) set increment id */
		if err != nil {
			ldb.log.Error("error database insert cfg setIncrement failed : %s,  cfg is : %v", err.Error(), cfg)
			return err
		}
	}

	dbKV := &dbKeyValue{}
	structKV(in, dbKV, cfg)

	err = ldb.insertKvToDb(dbKV) /* (kv to db) kv insert database (atomic) */
	if err != nil {
		ldb.log.Error("error database insert insertKvToDb failed : %s", err)
		return err
	}

	dbKV.id = cfg.id_
	m := new(modifyValue)
	m.NewKv = dbKV
	m.Id = dbKV.id
	m.OldKv = dbKV
	ldb.insertUndoState(cfg.Name, m, INSERT)
	return nil
}

func (ldb *LDataBase) insertKvToDb(dbKV *dbKeyValue) error {
	ldb.batch.Reset()
	defer ldb.batch.Reset()
	ldb.putBatch(dbKV)
	err := ldb.writeBatch()
	if err != nil {
		return err
	}

	return nil
}

func (ldb *LDataBase) setIncrement(cfg *structInfo) error {

	if _, ok := ldb.nextId[cfg.Name]; ok { // First insertion
		cfg.id_ = ldb.nextId[cfg.Name]
	}

	cfg.rId.Set(reflect.ValueOf(cfg.id_).Convert(cfg.rId.Type()))
	ldb.nextId[cfg.Name] = cfg.id_ + 1
	return nil
}

/*
*	Delete a piece of data from the database
*	returning null successfully
*	returns an error message
 */

func (ldb *LDataBase) Remove(in interface{}) error {
	err := ldb.remove(in)
	if err != nil {
		ldb.log.Error("failed : %s", err.Error())
		return err
	}
	return nil
}

func (ldb *LDataBase) remove(in interface{}) error {
	out := cloneInterface(in)
	err := ldb.find("id", in, out)
	if err != nil {
		return ErrNotFound
	}
	if !reflect.DeepEqual(in, out) {
		ldb.log.Error("remove DeepEqual error")
	}

	cfg, err := parseObjectToCfg(out)
	if err != nil {
		ldb.log.Error("failed : %s", err.Error())
		return err
	}

	dbKV := &dbKeyValue{}
	structKV(in, dbKV, cfg) /* (kv.index) all key and value*/

	dbKV.id = cfg.id_
	err = ldb.removeKvToDb(dbKV)
	if err != nil {
		ldb.log.Error("failed  : %s, dbKV is : %v", err, dbKV)
		return err
	}

	m := new(modifyValue)
	m.NewKv = dbKV
	m.Id = dbKV.id
	m.OldKv = dbKV
	ldb.insertUndoState(cfg.Name, m, REMOVE)
	return nil
}

func (ldb *LDataBase) removeKvToDb(dbKV *dbKeyValue) error {
	ldb.batch.Reset()
	defer ldb.batch.Reset()
	ldb.deleteBatch(dbKV)
	err := ldb.writeBatch()
	if err != nil {
		return err
	}
	return nil
}

/*
*	Modify an object
*	returning null successfully
*	returns an error message
 */

func (ldb *LDataBase) Modify(old interface{}, fn interface{}) error {
	err := ldb.modifyCallFn(old, fn)
	if err != nil {
		ldb.log.Error("%s", err.Error())
		return err
	}
	return nil
}

func (ldb *LDataBase) modifyCallFn(data interface{}, fn interface{}) error {

	dataRef := reflect.ValueOf(data)
	if dataRef.Kind() != reflect.Ptr {
		return ErrPtrNeeded
	}

	oldInter := cloneInterface(data)
	dataType := dataRef.Type()

	fnRef := reflect.ValueOf(fn)
	fnType := fnRef.Type()
	if fnType.NumIn() != 1 {
		ldb.log.Error("func Too many parameters : " + string(fnType.NumIn()))
		return errors.New("func Too many parameters ")
	}

	pType := fnType.In(0)

	if pType.Kind() != dataType.Kind() {
		ldb.log.Error("Parameter type does not match : " + pType.Kind().String() + " : " + dataType.Kind().String())
		return errors.New("Parameter type does not match ")
	}

	fnRef.Call([]reflect.Value{dataRef}) /*	call fn */
	// modify
	oldRef := reflect.ValueOf(oldInter)
	return ldb.modifyRefToKv(&oldRef, &dataRef)
}

func (ldb *LDataBase) modifyRefToKv(oldRef, newRef *reflect.Value) error {

	oldCfg, err := extractObjectTagInfo(oldRef)
	if err != nil {
		ldb.log.Error("extractObjectTagInfo oldRef failed : " + err.Error())
		return err
	}
	newCfg, err := extractObjectTagInfo(newRef)
	if err != nil {
		ldb.log.Error("extractObjectTagInfo newRef failed : " + err.Error())
		return err
	}
	if oldCfg.id_ != newCfg.id_ {
		ldb.log.Error("newCfg and oldCfg id failed,  newCfg id is :  %v,  oldCfg id is : %v", newCfg.id_, oldCfg.id_)
		return errors.New("newCfg and oldCfg id failed")
	}
	if oldCfg.Name != newCfg.Name {
		return errors.New("newCfg and oldCfg typeName failed")
	}
	newKV := &dbKeyValue{}
	oldKV := &dbKeyValue{}

	structKV(oldRef.Interface(), oldKV, oldCfg)
	structKV(newRef.Interface(), newKV, newCfg)

	err = ldb.modifyKvToDb(oldKV, newKV)
	if err != nil {
		return err
	}
	m := new(modifyValue)
	m.Id = oldCfg.id_
	//oldCfg.rId.Set(reflect.ValueOf(m.id).Convert(oldCfg.rId.Type()))
	m.NewKv = newKV
	m.OldKv = oldKV
	ldb.insertUndoState(oldCfg.Name, m, MODIFY)
	return nil
}

func (ldb *LDataBase) modifyKvToDb(oldKV, newKV *dbKeyValue) error {
	ldb.batch.Reset()
	defer ldb.batch.Reset()
	ldb.deleteBatch(oldKV)
	ldb.putBatch(newKV)
	err := ldb.writeBatch()
	if err != nil {
		ldb.log.Error("newKV is : %v", newKV)
		return err
	}

	return nil
}

/*
find object from database
@param tagName 		--> 	tag in field tags
@param in 			--> 	object
@param out 			-->		output(pointer)

@return
success 			-->		nil 	(out valid)
error 				-->		error 	(out invalid)

*/

func (ldb *LDataBase) Find(tagName string, in interface{}, out interface{}, skip ...SkipSuffix) error {
	return ldb.find(tagName, in, out, skip...)
}

func (ldb *LDataBase) find(tagName string, value interface{}, to interface{}, skip ...SkipSuffix) error {
	ldb.log.Info("tagName is: %v", tagName)
	fieldName := []byte(tagName)
	fields, err := getFieldInfo(tagName, value)
	if err != nil {
		ldb.log.Error("failed : %s", err.Error())
		return err
	}

	typeName := []byte(fields.typeName)

	suffix, err := fieldValueToByte(fields, skip...)
	if err != nil {
		ldb.log.Error("failed : %s", err.Error())
		return err
	}

	if len(suffix) == 0 {
		ldb.log.Error("failed : %s", err.Error())
		return ErrNotFound
	}

	key := splicingString(typeName, fieldName)

	key = append(key, suffix...)

	//ldb.log.Info("key is : %v", key)
	it := ldb.db.NewIterator(util.BytesPrefix(key), nil)

	//fmt.Println("ErrorNotFound key is : %v", key)
	//fmt.Println(it.Key(),"  ",it.Value())
	if !it.Next() {
		return ErrNotFound
	}
	return ldb.findFields(it.Value(), typeName, to)
}

func (ldb *LDataBase) getAllKv(key []byte) {

	it := ldb.db.NewIterator(nil, nil)
	for it.Next() {
		ldb.log.Info("key %v", it.Key(), "  value %v", it.Value())
	}
	it.Release()
	//os.Exit(65536)
}

func (ldb *LDataBase) findFields(key, typeName []byte, to interface{}) error {
	ldb.log.Info("key is : %v", key)

	val := splicingString(typeName, key)

	v, err := getDbKey(val, ldb.db)
	if err != nil {
		ldb.log.Error("failed : %s, val is %v", err.Error(), val)
		return err
	}
	err = DecodeBytes(v, to)
	if err != nil {
		ldb.log.Error("failed : %s, val is %v", err.Error(), val)
		return err
	}
	return nil
}

/*
@param tagName 		--> 	tag in field tags
@param in 			--> 	object

@return
success 			-->		iterator
error 				-->		error

*/
func (ldb *LDataBase) GetIndex(tagName string, in interface{}) (*MultiIndex, error) {
	return ldb.getIndex(tagName, in)
}

func (ldb *LDataBase) getIndex(tagName string, value interface{}) (*MultiIndex, error) {

	// fieldName == tagName --> Just different nextId
	fieldName := []byte(tagName)
	fields, err := getFieldInfo(tagName, value)
	if err != nil {
		ldb.log.Error("failed : %s, val is %v", err.Error(), fieldName)
		return nil, err
	}

	typeName := []byte(fields.typeName)
	begin := splicingString(typeName, fieldName)

	end := keyEnd(begin)

	ldb.log.Info("begin is : %v end is ", begin, end)
	it := newMultiIndex(typeName, fieldName, begin, end, ldb)
	return it, nil
}

func (ldb *LDataBase) GetMutableIndex(fieldName string, in interface{}) (*MultiIndex, error) {
	return ldb.GetIndex(fieldName, in)
}

func (ldb *LDataBase) LowerBound(begin, end, fieldName []byte, data interface{}, skip ...SkipSuffix) (*DbIterator, error) {
	key, typeName := ldb.dbPrefix(begin, fieldName, data, skip...)
	//return ldb.dbIterator(key,begin,end,typeName,false)
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Next() {
		return nil, ErrNotFound
	}

	if key != nil {
		if !it.Seek(key) {
			ldb.log.Error("key is : %v", key)
			return nil, ErrNotFound
		}
	}

	idx, err := newDbIterator(typeName, it, ldb.db)
	if err != nil {
		ldb.log.Error("failed is %s ,typeName is : %v", err.Error(), typeName)
		return nil, err
	}
	return idx, nil
}

func (ldb *LDataBase) UpperBound(begin, end, fieldName []byte, data interface{}, skip ...SkipSuffix) (*DbIterator, error) {

	key, typeName := ldb.dbPrefix(begin, fieldName, data, skip...)
	key = keyEnd(key)
	return ldb.dbIterator(key, begin, end, typeName, true)
}

func (ldb *LDataBase) dbIterator(key, begin, end, typeName []byte, upper bool) (*DbIterator, error) {
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Next() {
		return nil, ErrNotFound
	}

	if key != nil {
		if !it.Seek(key) {
			ldb.log.Error("key is : %v", key)
			return nil, ErrNotFound
		}
	}

	idx, err := newDbIterator(typeName, it, ldb.db)
	if err != nil {
		ldb.log.Error("failed is %s ,typeName is : %v", err.Error(), typeName)
		return nil, err
	}
	return idx, nil
}

func (ldb *LDataBase) EndIterator(begin, end, typeName []byte) (*DbIterator, error) {

	ldb.log.Info("begin : %v, end : %v, typeName: %v", begin, end, typeName)

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Next() {
		// not found  --> iterator is nil  == end
		itr := &DbIterator{it: nil, db: ldb.db, first: false, typeName: typeName, currentStatus: itEND}
		return itr, nil
	}

	if it.Last() {
		it.Next()
		itr := &DbIterator{it: it, db: ldb.db, first: false, typeName: typeName, currentStatus: itEND}
		return itr, nil
	}

	return nil, ErrNotFound
}

func (ldb *LDataBase) dbPrefix(begin_, fieldName []byte, data interface{}, skip ...SkipSuffix) ([]byte, []byte) {
	begin := cloneByte(begin_)
	ldb.log.Info("begin : %v, end : %v, fieldName: %v", begin, fieldName)
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		ldb.log.Error("failed %s", err.Error())
		return nil, nil
	}

	prefix, err := fieldValueToByte(fields, skip...)
	if err != nil {
		ldb.log.Error("failed %s", err.Error())
		return nil, nil
	}

	if len(prefix) == 0 {
		ldb.log.Error("prefix is empty")
		return nil, nil
	}

	prefix = append(begin, prefix...)

	ldb.log.Info("prefix is : %v", prefix)
	return prefix, []byte(fields.typeName)
}

func (ldb *LDataBase) Empty(begin, end, fieldName []byte) bool {

	ldb.log.Info("begin : %v, end : %v, fieldName: %v ", begin, end, fieldName)
	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	defer it.Release()
	if it.Next() {
		return false
	}
	return true
}

func (ldb *LDataBase) IteratorTo(begin, end, fieldName []byte, data interface{}, skip ...SkipSuffix) (*DbIterator, error) {

	ldb.log.Info("begin : %v, end : %v, fieldName: %v , greater: %t", begin, end, fieldName)
	fields, err := getFieldInfo(string(fieldName), data)
	if err != nil {
		ldb.log.Error("failed %s", err.Error())
		return nil, err
	}
	prefix, err := fieldValueToByte(fields, skip...)
	if err != nil {
		ldb.log.Error("failed %s", err.Error())
		return nil, err
	}
	if len(prefix) == 0 {
		return nil, errors.New("Get Field Value Failed")
	}

	key := []byte{}
	key = append(begin, prefix...)

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Seek(key) {
		ldb.log.Error("seek failed key is %v", key)
		return nil, ErrNotFound
	}

	idk := splicingString([]byte(fields.typeName), it.Value())
	idv, err := getDbKey(idk, ldb.db)
	if err != nil {
		ldb.log.Error("failed %s", err.Error())
		return nil, err
	}

	itr := &DbIterator{it: it, db: ldb.db, first: false, value: idv, typeName: []byte(fields.typeName)}
	return itr, nil
}

func (ldb *LDataBase) BeginIterator(begin, end, typeName []byte) (*DbIterator, error) {

	ldb.log.Info("begin : %v, end : %v, typeName: %v ", begin, end, typeName)

	it := ldb.db.NewIterator(&util.Range{Start: begin, Limit: end}, nil)
	if !it.Next() {
		ldb.log.Error("Next Failed")
		return nil, ErrNotFound
	}

	k := splicingString([]byte(typeName), it.Value())
	val, err := getDbKey(k, ldb.db)
	if err != nil {
		ldb.log.Error("failed %s %v", err.Error(), k)
		return nil, ErrNotFound
	}

	itr := &DbIterator{it: it, db: ldb.db, first: false, typeName: typeName, value: val, currentStatus: itBEGIN}

	return itr, nil
}

func (ldb *LDataBase) enable() bool { /* Whether the database enables the undo function*/
	return ldb.stack.Size() != 0
}

/*Batch operation*/

func (ldb *LDataBase) putBatch(dbKV *dbKeyValue) {
	for _, v := range dbKV.index {
		ldb.count++
		ldb.log.Debug("save key %v | %v", v.key, v.value)
		ldb.batch.Put(v.key, v.value)
	}
	ldb.count++
	ldb.log.Debug("save key %v | %v", dbKV.idk.key, dbKV.idk.value)
	ldb.batch.Put(dbKV.idk.key, dbKV.idk.value)
}

func (ldb *LDataBase) deleteBatch(dbKV *dbKeyValue) {
	for idx, _ := range dbKV.index {
		//ldb.log.Debug("delete key %v",dbKV.index[idx].key)
		ldb.batch.Delete(dbKV.index[idx].key)
	}
	//ldb.log.Debug("delete key %v",dbKV.idk.key)
	ldb.batch.Delete(dbKV.idk.key)
}

func (ldb *LDataBase) writeBatch() error {
	if ldb.batch.Len() > 0 {
		err := ldb.db.Write(ldb.batch, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

/* The following three functions are the undo functions provided by the database.*/

type undoOperatorType uint

const (
	INSERT = undoOperatorType(0)
	REMOVE = undoOperatorType(1)
	MODIFY = undoOperatorType(2)
)

func (ldb *LDataBase) insertUndoState(typeName string, value *modifyValue, oT undoOperatorType) {
	if !ldb.enable() {
		return
	}

	stack := ldb.getStack()
	if stack == nil {
		ldb.log.Warn("undo session empty")
		return
	}
	switch oT {
	case INSERT:
		stack.undoContainerInsert(typeName, value)
	case REMOVE:
		stack.undoContainerRemove(typeName, value)
	case MODIFY:
		stack.undoContainerModify(typeName, value)
	default:
		/* panic*/
	}
}

/*

The three functions here are the implementation of
the database in order to operate the double-ended queue
and the subsequent queues will provide corresponding operations

*/

func (ldb *LDataBase) getFirstStack() *undoContainer {
	stack := ldb.stack.First()
	switch typ := stack.(type) {
	case *undoContainer:
		return typ
	default:
		return nil
		//panic(TYPE_NOT_FOUND)
	}
	return nil
}

func (ldb *LDataBase) getSecond() *undoContainer {
	stack := ldb.stack.LastSecond()
	switch typ := stack.(type) {
	case *undoContainer:
		return typ
	default:
		return nil
		//panic(TYPE_NOT_FOUND)
	}
	return nil
}

func (ldb *LDataBase) getStack() *undoContainer {
	stack := ldb.stack.Last()
	if stack == nil {
		return nil
	}
	switch typ := stack.(type) {
	case *undoContainer:
		return typ
	default:
		return nil
	}
	return nil
}

func getDbKey(key []byte, db *leveldb.DB) ([]byte, error) {
	val, err := db.Get(key, nil)
	if err != nil {
		return nil, err
	}
	return val, err
}

func undoStateSquash(stack *undoState, preStack *undoState) {
	for key, value := range stack.OldValue {
		if _, ok := preStack.NewValue[key]; ok {
			continue
		}
		if _, ok := preStack.OldValue[key]; ok {
			continue
		}
		if _, ok := preStack.RemoveValue[key]; ok {
			// panic ?
		}
		preStack.OldValue[key] = value
	}

	for key, value := range stack.NewValue {
		preStack.NewValue[key] = value
	}

	for key, value := range stack.RemoveValue {

		if _, ok := preStack.NewValue[key]; ok {
			if _, ok := preStack.NewValue[key]; ok {
				delete(preStack.NewValue, key)
			}
			continue
		}
		if _, ok := preStack.OldValue[key]; ok {

			preStack.RemoveValue[key] = value
			if _, ok := preStack.OldValue[key]; ok {
				delete(preStack.OldValue, key)
			}
			continue
		}
		preStack.RemoveValue[key] = value
	}
}
