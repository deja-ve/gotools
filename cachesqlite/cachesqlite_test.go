package cachesqlite

import (
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

var (
	_defaultDBFile           = "./test_db/default_test.db"
	_defaultDBFileWithSpace  = "./test_db/default _ test.db"
	_defaultDBDirNoExists    = "./no_exists_dir/default _ test.db"
	_defaultTestKeyA         = "test_a"
	_defaultTestKeyB         = "test_b"
	_defaultTestKeyC         = "test_c"
	_defaultTestKeyAValue    = "test_value_a"
	_defaultTestKeyBValue    = "test_value_b"
	_defaultTestKeyCValue    = "test_value_c"
	_defaultTestKeyANoExists = "test_a_no_exists"
)

// ** 测试前请手动创建测试文件目录: ./test_db/ **
func TestNewCache(t *testing.T) {
	type args struct {
		storageFile string
	}
	tests := []struct {
		name string
		args args
		// want    *cacheSqlite
		wantErr bool
	}{
		{name: "ok, directory exists, no db file exists", args: args{storageFile: _defaultDBFile}, wantErr: false},
		{name: "ok, directory exists, db file exists", args: args{storageFile: _defaultDBFile}, wantErr: false},
		{name: "ok, directory exists, db file with spaces, no exists", args: args{storageFile: _defaultDBFileWithSpace}, wantErr: false},
		{name: "ok, directory exists, db file with spaces, exists ", args: args{storageFile: _defaultDBFileWithSpace}, wantErr: false},
		{name: "!ok, no directory exists", args: args{storageFile: _defaultDBDirNoExists}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCache(tt.args.storageFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// t.Errorf("NewCache() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func initCacheDB() *cacheSqlite {
	cacheS, err := NewCache(_defaultDBFile)
	if err != nil {
		panic(err)
	}
	if _, err := cacheS.Set(_defaultTestKeyA, _defaultTestKeyAValue); err != nil {
		panic(err)
	}
	if _, err := cacheS.Set(_defaultTestKeyB, _defaultTestKeyBValue); err != nil {
		panic(err)
	}
	if _, err := cacheS.Set(_defaultTestKeyC, _defaultTestKeyCValue); err != nil {
		panic(err)
	}
	return cacheS
}

func Test_cacheSqlite_Get(t *testing.T) {
	type args struct {
		key string
	}

	cacheS := initCacheDB()
	defer cacheS.Close()
	tests := []struct {
		name    string
		cacheDB *cacheSqlite
		args    args
		want    string
		wantErr bool
	}{
		{name: "ok, the key exist", cacheDB: cacheS, args: args{key: _defaultTestKeyA}, want: _defaultTestKeyAValue, wantErr: false},
		{name: "ok, the key exist", cacheDB: cacheS, args: args{key: _defaultTestKeyB}, want: _defaultTestKeyBValue, wantErr: false},
		{name: "ok, the key exist", cacheDB: cacheS, args: args{key: _defaultTestKeyC}, want: _defaultTestKeyCValue, wantErr: false},
		{name: "!ok, no the key exist", cacheDB: cacheS, args: args{key: _defaultTestKeyANoExists}, want: "", wantErr: true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cacheDB.Get(tt.args.key)
			fmt.Println(err)
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheSqlite.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cacheSqlite.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheSqlite_Del(t *testing.T) {
	type args struct {
		key string
	}
	cacheS := initCacheDB()
	defer cacheS.Close()

	tests := []struct {
		name    string
		c       *cacheSqlite
		args    args
		want    int64
		wantErr bool
	}{
		{name: "ok, the key test_a exists", c: cacheS, args: args{key: _defaultTestKeyA}, want: 1, wantErr: false},
		{name: "ok, the key test_b exists", c: cacheS, args: args{key: _defaultTestKeyB}, want: 1, wantErr: false},
		{name: "ok, the key test_c exists", c: cacheS, args: args{key: _defaultTestKeyC}, want: 1, wantErr: false},
		{name: "ok, no the key test_a_no_exists exists", c: cacheS, args: args{key: _defaultTestKeyANoExists}, want: 0, wantErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Del(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheSqlite.Del() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cacheSqlite.Del() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheSqlite_Set(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	cacheS, err := NewCache(_defaultDBFile)
	if err != nil {
		panic(err)
	}
	defer cacheS.Close()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name    string
		c       *cacheSqlite
		args    args
		want    int64
		wantErr bool
	}{
		{name: "ok, no the key test_a exists", c: cacheS, args: args{key: _defaultTestKeyA, value: _defaultTestKeyAValue}, want: 1, wantErr: false},
		{name: "ok, the key test_a exists", c: cacheS, args: args{key: _defaultTestKeyA, value: _defaultTestKeyAValue}, want: 1, wantErr: false},
		{name: "ok, the key test_a exists", c: cacheS, args: args{key: _defaultTestKeyA, value: "changed" + _defaultTestKeyAValue}, want: 1, wantErr: false},
		{name: "ok, no the key test_b exists", c: cacheS, args: args{key: _defaultTestKeyA, value: _defaultTestKeyAValue}, want: 1, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.c.Set(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("cacheSqlite.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("cacheSqlite.Set() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cacheSqlite_HasKey(t *testing.T) {
	type args struct {
		key string
	}
	cacheS := initCacheDB()
	defer cacheS.Close()
	tests := []struct {
		name string
		c    *cacheSqlite
		args args
		want bool
	}{
		{name: "ok, the key test_a exists", c: cacheS, args: args{key: _defaultTestKeyA}, want: true},
		{name: "ok, the key test_b exists", c: cacheS, args: args{key: _defaultTestKeyB}, want: true},
		{name: "ok, the key test_c exists", c: cacheS, args: args{key: _defaultTestKeyC}, want: true},
		{name: "!ok, no the key test_a_no_exists exists", c: cacheS, args: args{key: _defaultTestKeyANoExists},
			want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.HasKey(tt.args.key); got != tt.want {
				t.Errorf("cacheSqlite.HasKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
