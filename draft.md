## memo

- [ ] _test全部読む
- [ ] Association(), Related()
- [ ] join table handler
- [ ] `INSERT INTO `post_tags` (`post_id`,`tag_id`) SELECT 1,1 FROM DUAL WHERE NOT EXISTS (SELECT * FROM `post_tags` WHERE `post_id` = 1 AND `tag_id` = 1)` このSQLしらべておく
- [ ] callbackのRegister
- [ ] search.Scope.Unscopedのやくわりについて
- [ ] callback一通りよんで、今のmodelで説明できるかかためる。そうするとmodelの紹介かける。




本記事では、goのORM library [gorm](https://github.com/jinzhu/gorm)のAPIについて解説していきます。
1 tableに対応したstructのCRUD処理はsimpleなのですが、relationやtagまわりは[document](https://gorm.io/)だけだと実際の挙動がわかりづらいことがあったので、sourceを読んで挙動を整理していきます。
はじめに、gormがSQLを発行する仕組みを見ていき、そのあとで各APIについてふれていきます。
なお、gormのversionは[`v1.9.11`](https://github.com/jinzhu/gorm/tree/v1.9.11)です。


## 準備
実際に動かしながら試せる環境を用意します。[source code](https://github.com/ymgyt/gorm-blog-post)

```console
git clone https://github.com/ymgyt/gorm-blog-post.git
cd gorm-blog-post
direnv allow
task

go run create.go
```

go(`1.13.1`)とdocker-composeを利用します。direnvはoptionalです。`.envrc`に定義してある環境変数が設定されていればよいです。
databaseにはmysql(`5.7`)を利用しています。

## 用語の整理

DB `gorm.DB`
database MySQLとかPostgres
model gormのAPIに渡す、table構造に対応したstruct



## 主な登場人物

gormの中心となるstructをみていきます。実際の利用方法はこのあとのAPIでふれていくので、ここでは役割の概要について説明します。

### `DB`

Memo: RowsAffectedはSelect文でもincrementされる。

### `Scope`

### `ModelStruct`

### `Field`

## Relation

gormではrelation(modelとmodelの関係)は以下の4つに分類されます。

* `belong_to`
* `has_one`
* `has_many`
* `many_to_many`

`belong_to`と`has_one`は1:1の関係を表し、`has_many`は1:N、`many_to_many`はN:Mを表します。

### `belong_to`と`has_one`の違い

どちらも、fieldの型がassociation先のmodelになる点では共通しています。
association元がassociation先のprimary keyをfieldに保持している場合、`belong_to`に分類されます。
association先がassociation元のprimary keyをfieldに保持している場合、`has_one`に分類されます。
// TODO: わかりづらいので、具体例だす。

`belong_to`と`has_one`の違いはCreate時の生成順に影響します。
`belong_to`の場合は、association先、association元の順に生成され、`has_one`の場合はassociation元、association先の順に生成されます。


## Tag

gormの挙動を制御するための sturct fieldのtagについてみていきます。
tagのkeyは`gorm`と`sql`です。
[https://github.com/jinzhu/gorm/blob/v1.9.11/model_struct.go#L640]

formatは`gorm:"key:value;key:value"`です。keyはcase insensitiveで、`gorm:"default"`と書いても、`gorm:"DEFAULT"`と書いても同じように扱われます。
[https://github.com/jinzhu/gorm/blob/v1.9.11/model_struct.go#L647]

valueをとらないkeyにvalueを書いても無視されます。
[https://github.com/jinzhu/gorm/blob/v1.9.11/model_struct.go#L649]



### `-`

gormはそのfieldを処理の対象にしなくなります。

[https://github.com/jinzhu/gorm/blob/81c17a7e2529c59efc4e74c5b32c1fb71fb12fa2/model_struct.go#L197-L198]

### `DEFAULT`

そのfieldにDatabase側で定義されたdefault値があることをgormに伝えます。
//TODO source
INSERT実行後に、SELECT文を発行してmodelにsetする処理が追加されます。
`gorm:"DEFAULT:JP"`のようにvalueに値を書いても、Database側のdefault値が利用されます。

### `PRIMARY_KEY`

modelのprimary keyとして扱われます。

### `EMBEDDED`

このタグを付与するとfiledの型として定義されているstructのfieldが展開されます。`EMBEDDED_PREFIX` keyでcolumnのprefixを指定できます。


### `COLUMN`

filedに対応するdatabaseのcolumn名を指定します。
[https://github.com/jinzhu/gorm/blob/81c17a7e2529c59efc4e74c5b32c1fb71fb12fa2/model_struct.go#L611]
指定しない場合はdefaultの変換ロジックが利用されます。

[https://github.com/jinzhu/gorm/blob/81c17a7e2529c59efc4e74c5b32c1fb71fb12fa2/naming.go#L71]


### `PRELOAD`

`DB.Set("gorm:auto_preload", true)`でauto preloadを有効にした際に、preloadさせたくない場合に利用。

[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_query_preload.go#L103-L110]

## Create

### callback

```go
// Define callbacks for creating
func init() {
	DefaultCallback.Create().Register("gorm:begin_transaction", beginTransactionCallback)
	DefaultCallback.Create().Register("gorm:before_create", beforeCreateCallback)
	DefaultCallback.Create().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	DefaultCallback.Create().Register("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	DefaultCallback.Create().Register("gorm:create", createCallback)
	DefaultCallback.Create().Register("gorm:force_reload_after_create", forceReloadAfterCreateCallback)
	DefaultCallback.Create().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	DefaultCallback.Create().Register("gorm:after_create", afterCreateCallback)
	DefaultCallback.Create().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}
```
create callbackの設定はこのように行われており、上から順番に実行されていきます。
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_create.go#L9-L18]


```go
func beginTransactionCallback(scope *Scope) {
	scope.Begin()
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_save.go#L8-L10]

transactionが張られていない場合はtransactionを張ります。

## Query





## Delete

### callback
```go
// Define callbacks for deleting
func init() {
	DefaultCallback.Delete().Register("gorm:begin_transaction", beginTransactionCallback)
	DefaultCallback.Delete().Register("gorm:before_delete", beforeDeleteCallback)
	DefaultCallback.Delete().Register("gorm:delete", deleteCallback)
	DefaultCallback.Delete().Register("gorm:after_delete", afterDeleteCallback)
	DefaultCallback.Delete().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_delete.go#L9-L15]

transactionをはってから

```go
// beforeDeleteCallback will invoke `BeforeDelete` method before deleting
func beforeDeleteCallback(scope *Scope) {
	if scope.DB().HasBlockGlobalUpdate() && !scope.hasConditions() {
		scope.Err(errors.New("missing WHERE clause while deleting"))
		return
	}
	if !scope.HasError() {
		scope.CallMethod("BeforeDelete")
	}
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_delete.go#L18-L27]

Global(条件を指定しない)Deleteを許可しないように設定しているかの確認をおこない、BeforeDelete Hookを呼び出します。


```go
// deleteCallback used to delete data from database or set deleted_at to current time (when using with soft delete)
func deleteCallback(scope *Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok {
			extraOption = fmt.Sprint(str)
		}

		deletedAtField, hasDeletedAtField := scope.FieldByName("DeletedAt")

		if !scope.Search.Unscoped && hasDeletedAtField {
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),
				scope.Quote(deletedAtField.DBName),
				scope.AddToVars(scope.db.nowFunc()),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else {
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_delete.go#L29-L56]

modelに`DeletedAt` fieldが定義されていれば、soft deleteを行います。定義されていなければScopeからWhere句を生成してDELETE文を発行します。
そのあとは、`AfterDelete` Hookを呼んで、CommitOrRollbackします。


## etc
どこかに書くものたち

### transaction

```go
type sqlDb interface {
	Begin() (*sql.Tx, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

// SQLDB return *sql.DB
func (scope *Scope) SQLDB() SQLCommon {
	return scope.db.db
}

```

```go
// Begin start a transaction
func (scope *Scope) Begin() *Scope {
	if db, ok := scope.SQLDB().(sqlDb); ok {
		if tx, err := db.Begin(); scope.Err(err) == nil {
			scope.db.db = interface{}(tx).(SQLCommon)
			scope.InstanceSet("gorm:started_transaction", true)
		}
	}
	return scope
}
```

DBが保持しているconnectionが`Begin()` methodをもっていれば、呼び出して以後のscopeで参照するようになっています。
DBやScopeで複数回`Begin()`を呼んでも、複数回transactionをはらないようになっています。

[https://github.com/jinzhu/gorm/blob/v1.9.11/scope.go#L403-L411]


### Hooks

modelに定義している、`BeforeSave()`や`AfterCreate()`をgormが呼んでくれる仕組みをみていきます。

```go

// callbackの中でhookを呼び出す
scope.CallMethod("BeforeSave")

// CallMethod call scope value's method, if it is a slice, will call its element's method one by one
func (scope *Scope) CallMethod(methodName string) {
	if scope.Value == nil {
		return
	}

	if indirectScopeValue := scope.IndirectValue(); indirectScopeValue.Kind() == reflect.Slice {
		for i := 0; i < indirectScopeValue.Len(); i++ {
			scope.callMethod(methodName, indirectScopeValue.Index(i))
		}
	} else {
		scope.callMethod(methodName, indirectScopeValue)
	}
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/scope.go#L242-L254]

`Scope.IndirectValue()`は`Scope`が保持しているvalueの`reflect.Value`で、pointerの場合は`Value.Elem()`をとったものです。
sliceの場合はelementそれぞれに対して呼び出しています。

```
func (scope *Scope) callMethod(methodName string, reflectValue reflect.Value) {
	// Only get address from non-pointer
	if reflectValue.CanAddr() && reflectValue.Kind() != reflect.Ptr {
		reflectValue = reflectValue.Addr()
	}

	if methodValue := reflectValue.MethodByName(methodName); methodValue.IsValid() {
		switch method := methodValue.Interface().(type) {
		case func():
			method()
		case func(*Scope):
			method(scope)
		case func(*DB):
			newDB := scope.NewDB()
			method(newDB)
			scope.Err(newDB.Error)
		case func() error:
			scope.Err(method())
		case func(*Scope) error:
			scope.Err(method(scope))
		case func(*DB) error:
			newDB := scope.NewDB()
			scope.Err(method(newDB))
			scope.Err(newDB.Error)
		default:
			scope.Err(fmt.Errorf("unsupported function %v", methodName))
		}
	}
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/scope.go#L432-L460]

reflectionを利用して、Goで動的にmethodを呼び出すcodeとして参考になりました。
この実装から、gormがsupportしているhookの関数のsignatureがわかります。

### GlobalUpdate/Delete
```go
// BlockGlobalUpdate if true, generates an error on update/delete without where clause.
// This is to prevent eventual error with empty objects updates/deletions
func (s *DB) BlockGlobalUpdate(enable bool) *DB {
	s.blockGlobalUpdate = enable
	return s
}
```
[https://github.com/jinzhu/gorm/blame/master/main.go#L180-L185]


条件を指定しないUpdate/Deleteを許可するかどうかを制御できます。gormの場合、意図しない全件削除がおきやすいのでtrueにしてもよいかもしれません。
DB単位の設定なので、必要な場合だけ許可するようにもできます。

### Scope setting

#### instance 単位
* `gorm:skip_query_callback` : ???
* `gorm:only_preload`: ???
* `gorm:order_by_primary_key` : ???
* `gorm:query_destination`: ???
* `gorm:started_transaction`: `Scope.db.db`(保持しているConnection)がtransactionのときtrue
* `gorm:skip_bindvar`: ???
* `gorm:skip_query_callback`: query callbackの処理をおこなわずにreturnする(どこで設定している???)
* `gorm:only_preload`:  query callbackの処理をおこなわずにreturnする(どこで設定???)
* `gorm:order_by_primary_key`:
* `gorm:query_destination`:  Find()の結果のbind先を変更。


#### DB単位

* `gorm:delete_option`: ???
* `gorm:query_option`:  発行するSelect文にappendされる文字列を指定できます。documentではFOR UPDATEが例としてあがっています。
* `gorm:auto_preload`:  なんらかのassociation(has_one,belong_to, has_many, many_to_many)をもつfieldをPreloadする。


