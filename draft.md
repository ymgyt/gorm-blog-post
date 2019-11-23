## memo

- [ ] _test全部読む
- [ ] Association(), Related()
- [ ] join table handler
- [ ] `INSERT INTO` `post_tags` (`post_id`,`tag_id`) SELECT 1,1 FROM DUAL WHERE NOT EXISTS (SELECT * FROM `post_tags` WHERE `post_id` = 1 AND `tag_id` = 1)` このSQLしらべておく
- [ ] search.Scope.Unscopedのやくわりについて
- [ ] callback一通りよんで、今のmodelで説明できるかかためる。そうするとmodelの紹介かける。
- Association
- Delete Nested
- DB, Scope, GetModelStructの解説
- Select Saveの挙動


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

### Find

#### callback

```go
// Define callbacks for querying
func init() {
	DefaultCallback.Query().Register("gorm:query", queryCallback)
	DefaultCallback.Query().Register("gorm:preload", preloadCallback)
	DefaultCallback.Query().Register("gorm:after_query", afterQueryCallback)
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_query.go#L9-L14]

findのcallbackは3つで、Valueのcallback, associationを取得するpreload callback, hookを呼び出すcallbackで構成されています。

##### `queryCallback`

```go
// queryCallback used to query data from database
func queryCallback(scope *Scope) {
	if _, skip := scope.InstanceGet("gorm:skip_query_callback"); skip {
		return
	}

	//we are only preloading relations, dont touch base model
	if _, skip := scope.InstanceGet("gorm:only_preload"); skip {
		return
	}

	defer scope.trace(scope.db.nowFunc())

	var (
		isSlice, isPtr bool
		resultType     reflect.Type
		results        = scope.IndirectValue()
	)

	if orderBy, ok := scope.Get("gorm:order_by_primary_key"); ok {
		if primaryField := scope.PrimaryField(); primaryField != nil {
			scope.Search.Order(fmt.Sprintf("%v.%v %v", scope.QuotedTableName(), scope.Quote(primaryField.DBName), orderBy))
		}
	}

	if value, ok := scope.Get("gorm:query_destination"); ok {
		results = indirect(reflect.ValueOf(value))
	}

	if kind := results.Kind(); kind == reflect.Slice {
		isSlice = true
		resultType = results.Type().Elem()
		results.Set(reflect.MakeSlice(results.Type(), 0, 0))

		if resultType.Kind() == reflect.Ptr {
			isPtr = true
			resultType = resultType.Elem()
		}
	} else if kind != reflect.Struct {
		scope.Err(errors.New("unsupported destination, should be slice or struct"))
		return
	}

	scope.prepareQuerySQL()

	if !scope.HasError() {
		scope.db.RowsAffected = 0
		if str, ok := scope.Get("gorm:query_option"); ok {
			scope.SQL += addExtraSpaceIfExist(fmt.Sprint(str))
		}

		if rows, err := scope.SQLDB().Query(scope.SQL, scope.SQLVars...); scope.Err(err) == nil {
			defer rows.Close()

			columns, _ := rows.Columns()
			for rows.Next() {
				scope.db.RowsAffected++

				elem := results
				if isSlice {
					elem = reflect.New(resultType).Elem()
				}

				// Memo: resultsが[]stringだと scope.New(string).Fields()だけど問題ないのか
				// => scope.GetModelStructはvalueがstruct以外だと空を返すので、問題ない
				scope.scan(rows, columns, scope.New(elem.Addr().Interface()).Fields())

				if isSlice {
					if isPtr {
						results.Set(reflect.Append(results, elem.Addr()))
					} else {
						results.Set(reflect.Append(results, elem))
					}
				}
			}

			if err := rows.Err(); err != nil {
				scope.Err(err)
			} else if scope.db.RowsAffected == 0 && !isSlice {
				scope.Err(ErrRecordNotFound)
			}
		}
	}
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_query.go#L17-L98]

scopeに設定されている条件(Where,Join,Group...)からSQLを生成して、Databaseに発行します。Queryの場合でも、`RowsAffected`をincrementしています。sliceの場合見つからなくても`RecordNotFound`errorが返らない理由も最後のifからわかります。
`scope.scan()`の中では最終的には`database/sql.Rows.Scan()`を呼び出しています。
gormのreflectionを利用した、scanの実装も追いたいのですが、あまりにも長くなるので、この記事では諦めます。
(`**uint`のようなpointerのpointerをreflectionで扱う処理がでてくるのですが、追いきれませんでした。[https://github.com/jinzhu/gorm/blob/v1.9.11/scope.go#L500])


#### `preloadCallback`

```go
// preloadCallback used to preload associations
func preloadCallback(scope *Scope) {
	if _, skip := scope.InstanceGet("gorm:skip_query_callback"); skip {
		return
	}

	if ap, ok := scope.Get("gorm:auto_preload"); ok {
		// If gorm:auto_preload IS NOT a bool then auto preload.
		// Else if it IS a bool, use the value
		if apb, ok := ap.(bool); !ok {
			autoPreload(scope)
		} else if apb {
			autoPreload(scope)
		}
	}

	if scope.Search.preload == nil || scope.HasError() {
		return
	}

	var (
		preloadedMap = map[string]bool{}
		fields       = scope.Fields()
	)

	for _, preload := range scope.Search.preload {
		var (
			preloadFields = strings.Split(preload.schema, ".")
			currentScope  = scope
			currentFields = fields
		)

		for idx, preloadField := range preloadFields {
			var currentPreloadConditions []interface{}

			if currentScope == nil {
				continue
			}

			// if not preloaded
			if preloadKey := strings.Join(preloadFields[:idx+1], "."); !preloadedMap[preloadKey] {

				// assign search conditions to last preload
				if idx == len(preloadFields)-1 {
					currentPreloadConditions = preload.conditions
				}

				for _, field := range currentFields {
					if field.Name != preloadField || field.Relationship == nil {
						continue
					}


					switch field.Relationship.Kind {
					case "has_one":
						currentScope.handleHasOnePreload(field, currentPreloadConditions)
					case "has_many":
						currentScope.handleHasManyPreload(field, currentPreloadConditions)
					case "belongs_to":
						currentScope.handleBelongsToPreload(field, currentPreloadConditions)
					case "many_to_many":
						currentScope.handleManyToManyPreload(field, currentPreloadConditions)
					default:
						scope.Err(errors.New("unsupported relation"))
					}

					preloadedMap[preloadKey] = true
					break
				}

				if !preloadedMap[preloadKey] {
					scope.Err(fmt.Errorf("can't preload field %s for %s", preloadField, currentScope.GetModelStruct().ModelType))
					return
				}
			}

			// preload next level
			if idx < len(preloadFields)-1 {
				currentScope = currentScope.getColumnAsScope(preloadField)
				if currentScope != nil {
					currentFields = currentScope.Fields()
				}
			}
		}
	}
}

func autoPreload(scope *Scope) {
	for _, field := range scope.Fields() {
		if field.Relationship == nil {
			continue
		}

		if val, ok := field.TagSettingsGet("PRELOAD"); ok {
			if preload, err := strconv.ParseBool(val); err != nil {
				scope.Err(errors.New("invalid preload option"))
				return
			} else if !preload {
				continue
			}
		}

		scope.Search.Preload(field.Name)
	}
}

```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_query_preload.go#L12-L115]
association先のfiledを取得します。auto_preloadを指定しておくか、`db.Preload("Author.Reviews")`のようにassociation先のassociationをしておくと、Queryを発行してくれます。associationのfiledの`Relationship`がnilでないことが条件なので、うまくassociation先がloadされない場合は、`scope.GetModelStruct()`して、associationの分析が意図どおりにされているか確認すると原因がわかりそうです。

`preloadCallback`実行後は、`AfterFind` Hookを呼び出して、Findは完了です。


### `FirstOrInit`

```go
func (s *DB) FirstOrInit(out interface{}, where ...interface{}) *DB {
	c := s.clone()
	if result := c.First(out, where...); result.Error != nil {
		if !result.RecordNotFound() {
			return result
		}
		c.NewScope(out).inlineCondition(where...).initialize()
	} else {
		c.NewScope(out).updatedAttrsWithValues(c.search.assignAttrs)
	}
	return c
}

func (scope *Scope) initialize() *Scope {
	for _, clause := range scope.Search.whereConditions {
		scope.updatedAttrsWithValues(clause["query"])
	}
	scope.updatedAttrsWithValues(scope.Search.initAttrs)
	scope.updatedAttrsWithValues(scope.Search.assignAttrs)
	return scope
}

```

[https://github.com/jinzhu/gorm/blob/v1.9.11/main.go#L408-L419]

Where条件でSelectして、Attrs/Assignで指定したfiledをValueに代入します。
見つかっても、見つからなくてもINSERT処理は実行しません。
AttrsとAssignの共通点はどちらもSelect時の条件には反映されないこと、違いは見つかったか見つからなかったかに応じてValueに代入されるかが変わることです。
INSERT文を発行したくない場合には、`FirstOrCreate`ではなくこちらを利用するといいと思います。


｜Where句|Insert|Attrs|Assign|
| ----  | ---- | --- | ----- |
| Found | No   | No  |  Yes  |
| Not   | No   | Yes |  Yes  |


### `FirstOrCreate`

```go
func (s *DB) FirstOrCreate(out interface{}, where ...interface{}) *DB {
	c := s.clone()
	if result := s.First(out, where...); result.Error != nil {
		if !result.RecordNotFound() {
			return result
		}
		return c.NewScope(out).inlineCondition(where...).initialize().callCallbacks(c.parent.callbacks.creates).db
	} else if len(c.search.assignAttrs) > 0 {
		return c.NewScope(out).InstanceSet("gorm:update_interface", c.search.assignAttrs).callCallbacks(c.parent.callbacks.updates).db
	}
	return c
}

```
[https://github.com/jinzhu/gorm/blob/v1.9.11/main.go#L423-L434]

Where句の条件でみつからなかった場合INSERT文が発行されます。
INSERTする際は、Attrs/Assignで指定したfiledが反映されてCreateされます。
見つかった場合、Assignを指定しているとUPDATE文が発行されます。


｜Where句|Insert|Attrs|Assign|
| ----  | ---- | --- | ----- |
| Found | No   | No  |  Yes  |
| No    | Yes  | Yes |  Yes  |


### `SubQuery/QueryExpr`

SubQueryを作る際などScopeに設定した条件からSQLがほしいときに利用します。SubQueryとQueryExprの違いは結果を`()`でかこってくれるかどうかだけです。
メイン側の条件(Where句)が既に設定されている場合は、一度`DB.New()`を経由して条件をクリアしてから設定する必要があります。(したの場合は必要なし)


```go
db.Where("author_id = ?",
    db.Table("authors").Where(model.Author{Name: "ymgyt"}).Select("id").Limit(1).SubQuery()
).Set("gorm:auto_preload", true).First(&post2)

// => SELECT * FROM `posts`
//    WHERE (author_id = (
//      SELECT id FROM `authors`  WHERE (`authors`.`name` = 'ymgyt') LIMIT 1))
//    ORDER BY `posts`.`id` ASC LIMIT 1;
```

## Update

### callback

```go
// Define callbacks for updating
func init() {
	DefaultCallback.Update().Register("gorm:assign_updating_attributes", assignUpdatingAttributesCallback)
	DefaultCallback.Update().Register("gorm:begin_transaction", beginTransactionCallback)
	DefaultCallback.Update().Register("gorm:before_update", beforeUpdateCallback)
	DefaultCallback.Update().Register("gorm:save_before_associations", saveBeforeAssociationsCallback)
	DefaultCallback.Update().Register("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	DefaultCallback.Update().Register("gorm:update", updateCallback)
	DefaultCallback.Update().Register("gorm:save_after_associations", saveAfterAssociationsCallback)
	DefaultCallback.Update().Register("gorm:after_update", afterUpdateCallback)
	DefaultCallback.Update().Register("gorm:commit_or_rollback_transaction", commitOrRollbackTransactionCallback)
}
```
[https://github.com/jinzhu/gorm/blob/v1.9.11/callback_update.go#L11-L21]



## Delete

### callback
`eha``go
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

## Association Mode

- Association()呼ぶ前にModelに渡すstructにはprimary keyがあることが必須
- Associaation()に渡したcolumn名でFieldがひけることがひつよう
- column名のfieldにはRelationshipが必要。
- supportしているmethodは
  - Find
  - Append
  - Replace
  - Delete
  - Clear
  - Count




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
* `gorm:skip_bindvar`: SubQuery
* `gorm:skip_query_callback`: query callbackの処理をおこなわずにreturnする(どこで設定している???)

* `gorm:order_by_primary_key` : Firstのorder指定で利用
* `gorm:only_preload`:  query callbackの処理をおこなわずにreturnする(どこで設定???)
* `gorm:query_destination`:  Find()の結果のbind先を変更。

* `gorm:update_interface`: update時のfieldを指定。`DB.Update(), DB.UpdateColumns(), DB.FirstOrCreate()`時に利用
* `gorm:update_column`: ???
* `gorm:update_attrs`: ???
* `gorm:ignore_protected_attrs`: ??? Setするコードしか存在していない

* `gorm:started_transaction`: `Scope.db.db`(保持しているConnection)がtransactionのときtrue
* `skip_bindvar`: SQLのplaceholderを`?`固定にする。MySQLの場合はglobalでこのoptionをいれても問題なさそう。

#### DB単位

* `gorm:delete_option`: ???
* `gorm:query_option`:  発行するSelect文にappendされる文字列を指定できます。documentではFOR UPDATEが例としてあがっています。
* `gorm:auto_preload`:  なんらかのassociation(has_one,belong_to, has_many, many_to_many)をもつfieldをPreloadする。


