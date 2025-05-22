## Как я понял задачу:
1. надо исследовать как ogen генерирует код
2. как он добавляет комментарии из конфика OpenApi (найти поле, к которому добавить описание, чтобы оно попало в комментарий к handler-методу)
3. написать yaml/json конфиг для OpenApi
4. проверить что методы содержат нужные комменты после генерации

## Как Ogen генерит код 
Примеры [отсюда](https://github.com/ogen-go/ogen/tree/main/examples) и [отсюда](https://ogen.dev/docs/intro/)  
Он парсит конфигурацию и создаюет:
  - отдельные структуры для запросов и ответов (то есть в хэндере мы будет тело POST-запроса/параметры GET-запроса парсит в стуктуру а не параметры по отдельности (парметры могут быть как `in: query` так и `in: path`)) и их парсеры/сераилайзеры
  - интерфейсы с handler-методами (сервисы из MVC)
  - маршрутизацию
  - для каждого endpoint'а из спецификации создается соответствующий метод в интерфейсе handler'а
  - маршализацию/анмаршализацию джейсонов кодом без рефлекшена
  

## Как комменты из конфига попадают непосредственно в код
В OAPI конфиге есть несколько мест, самое очевидно description в полях схем и в описании эндпоинтов, итого:
 - `description` для операции (в `paths.{path}.{method}.description`)
 - `summary` для операции (в `paths.{path}.{method}.summary`)
 - `description` для тегов (в `tags[].description`)
 - описания схем и запросов (для документации в swagger)

(по умолчанию используется `description` из операции для генерации коментов, если он не найден, то используется `summary` как менее подробный)

Кроме того `summary` проставляется в поле `OperationSummary` в [`oas_handlers_gen.go`](gen_api/oas_handlers_gen.go#L276) и [`oas_router_gen.go`](gen_api/oas_router_gen.go#306) 

Например `summary` `getUserBook`:
```bash
grep "Get book by it's id" gen_api/*
```
```
gen_api/oas_handlers_gen.go:                    OperationSummary: "Get book by it's id",
gen_api/oas_router_gen.go:                                                      r.summary = "Get book by it's id"
```

И `description` попадает в комментарии к коду:
```bash
grep "Returns a book by book's id" gen_api/*
```
```
gen_api/oas_client_gen.go:      // Returns a book by book's id.
gen_api/oas_client_gen.go:// Returns a book by book's id.
gen_api/oas_handlers_gen.go:// Returns a book by book's id.
gen_api/oas_server_gen.go:      // Returns a book by book's id.
gen_api/oas_unimplemented_gen.go:// Returns a book by book's id.
```

сами комментарии к сгенерированным методам можно посмореть в [oas_server_gen.go](gen_api/oas_server_gen.go#L10) 

например это выглядит так:
```yaml
paths:
  /users:
    get:
      tags: [users]
      description: returns the list of all registered users
      summary: get all users
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Users'
  /user:
    get:
      summary: get user by id
      parameters:
        - in: query
          name: id
          schema:
            type: integer
            minimum: 1
          description: id of user
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'

components:
  schemas:
    User:
      type: object
      description: user
      properties:
        id:
          type: integer
          description: id of the user
        ...
    Users:
      descriptioh: array of `User`
      type: array
      items:
        $ref: '#/components/schemas/User'
       
```
В приоритете для кодогенерации description комментарий (он подробнее и может содержать Markdown) 
и сгенероиванный код будет выглядит примерно так:

P.S. видимо название метода: путь к энпоинту + имя метода, например `/api/users/getall` -> `APIUsersGetAllGet` и `APIUsersGetAllGetParams/Res` если в поле эндпоинта `operationId` не указано обратное

Тогда снерегированный код будет примерно таким:
`ogen --package api --target ./gen_api ./spec.yaml`:
```go
type Handler interface {
	// GetUsers implements GET /users operation.
	//
	// returns the list of all registered users
	//
	// GET /users
	UsersGet(ctx context.Context) ([]User, error)

	// GetUser implements GET /user operation.
	//
	// returns the list of all registered users
	//
	// GET /user
	UserGet(ctx context.Context, params UserGetParams) (User, error)

    ...
}

// UserGetParams is parameters of GET /user operation.
type UserGetParams struct {
    // id of user
    Id OptInt
}

// UsersGet implements GET /users operation.
//
// returns the list of all registered users
//
// GET /users
func (UnimplementedHandler) UsersGet(ctx context.Context) (r []User, _ error) {
	return r, ht.ErrNotImplemented
}

// UserGet implements GET /user operation.
//
// get user by id
//
// GET /user
func (UnimplementedHandler) UserGet(ctx context.Context, paramd UserGetParams) (r User, _ error) {
	return r, ht.ErrNotImplemented
}

```
В routes будет:
```go
   r.Summary = "get all users"
```
И дальше мы может для каждого роута уже явно вызвать .Summary() и как-то использовать

### Итого:
   - `summary` для краткого однострочного описания, например в Swagger
   - `description` попадает в комментарий к сгенерированному коду (например [oas_server_gen.go](gen_api/oas_server_gen.go#L9)) и
     содержит уже более подробное описание поведение и ограничений 
   - `tags[].description` группирует эндпоинты по тегу в одну группу запросов по смыслу (в код не попадает, только для отображения документации API в том же Swagger)

## Напишем семпловый конфиг 
создамим суррогатный пример сервиса, например это будет сервис по чтению книг и будет для пользователя с неким ID хранить список 
(book_id, book_properties, page) книг, где page - номер страницы которую сейчас читает юзер
```go
type Book struct {
    ID        int       `json:"id"`
    Title     string    `json:"title"`
    Author    string    `json:"author"`
    Published time.Time `json:"published"`
    Page      int       `json:"page"`
}
```
Я написал [ямлик](api.yml) конфига (по примерам их документации) и [сервер](main.go) с [клиентом](client/main.go)
(для простоты нормальную структуру пакета соблюдать не буду)

комментарии можно увидеть в [handlers](gen_api/oas_handlers_gen.go#L33) и [server](gen_api/oas_server_gen.go#L10)

Я использовую `*Res`-структуры, можно включить convenient errors, но тогда управление ставновится менее очевидным:
Во первых мне придется делать обертку над `api.Error` чтобы она возвращада JSON строку на `Error()`, что не очень идиоматично,
Во вторых тогда если возвращать ошибку из метода то будет не ошибка 500, а моя кастомная ошибка, которая не устанавливает в заголовок нужный код в заголовок, что опять лишний код в middlewares:

вместо 
```go
func (s *serviceImpl) GetUserBook(ctx context.Context, params api.GetUserBookParams) (api.GetUserBookRes, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if books, ok := s.users[params.UserID]; !ok {
		return err(http.StatusNotFound, "user %d not found", params.UserID), nil
	} else if book, ok := books[params.BookID]; !ok {
		return err(http.StatusNotFound, "book %d not found for user %d", params.BookID, params.UserID), nil
	} else {
		return &book, nil
	}
}
```
было бы что то вроде
```go

type apiError struct{ e *api.Error } // иначе коллизия поля .Error и метода .Error()

func (e apiError) Error() string {
	return fmt.Sprintf(`{"status_code": %d, "message": %s`, e.e.StatusCode, e.e.Message)
}

func wrap(e *api.Error) error {
	return apiError{e}
}

func (s *serviceImpl) GetUserBook(ctx context.Context, params api.GetUserBookParams) (*api.Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if books, ok := s.users[params.UserID]; !ok {
		return nil, wrap(err(http.StatusNotFound, "user %d not found", params.UserID))
	} else if book, ok := books[params.BookID]; !ok {
		return nil, wrap(err(http.StatusNotFound, "book %d not found for user %d", params.BookID, params.UserID))
	} else {
		return &book, nil
	}
}
```

(Кроме GetUserBooks, там не может возникнуть ошибки в рамках сервиса)

## Как запустить
Собрать
```bash
make all
```

Сервер
```bash
./server
```

Клиент
```bash
./client/client i
```
