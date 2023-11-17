# codefixture
codefixture is a Go library designed to set up test fixtures through code.


## Motivation
Writing complex integration tests is always necessary but challenging.
This library assists in creating such tests programatically.
codefixture doesn't rely on external files such as JSON, YAML or SQL,
as it is more maintenance to define specific test conditions directly in the code.

## Use Case

codefixture is useful when you want to conduct tests in the following scenarios:

* Integration tests (Round-trip tests)
* Application layer tests (use case layer)
* Tests involving database access

## Installation

```bash
go install github.com/bluegreenhq/codefixture
```

## Usage

1. Import codefixture in your test file.
```go
import "github.com/bluegreenhq/codefixture"
```

2. Create a `FixtureBuilder`.
```go
builder := codefixture.NewFixtureBuilder()
```

3. Register writers with `FixtureBuilder`.

* Using Gorm
(Note: this is an example. codefixture is not dependent on any DB libraries.)

```go
conn := gorm.Open(...)
err := builder.RegisterWriter(&Person{}, func(m any) (any, error) {
    m, ok := m.(*Person)
    if !ok {
        return nil, errors.New("invalid type")
    }
    res := conn.Create(m)
    return m, res.Error
})
```

4. Add models and relations to `FixtureBuilder`.

```go
p, _ := builder.AddModel(&Person{Name: "John"})
g, _ := builder.AddModel(&Group{Name: "Family"})

builder.AddRelation(p, g, func(p, g any) {
    p.(*Person).GroupID = g.(*Group).ID
})
```

5. Build `Fixture` from `FixtureBuilder`.

```go
fixture, err := builder.Build()
```

6. Access your models from `Fixture`.

```go
fmt.Printf("Person name: %s\n", fixture.GetModel(p).(*Person).Name)
fmt.Printf("Group name: %s\n", fixture.GetModel(g).(*Group).Name)
```

## License

MIT
