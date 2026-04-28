# serror

A structured error library for Go that integrates with `log/slog`. Errors carry a typed trait, a call-site trace, and arbitrary context — and log as structured groups with no extra glue.

## Install

```bash
go get github.com/0xnshd/serror
```

## Usage

Define your error kinds as package-level `ErrorTrait` values, then use `New` to wrap a root cause.

---

### ErrorTrait — typed error classification

```go
type ErrorTrait struct {
    Code  int
    Trait string
}
```

Define your application's error kinds once:

```go
var (
    ErrNotFound  = serror.ErrorTrait{Code: 404, Trait: "not_found"}
    ErrForbidden = serror.ErrorTrait{Code: 403, Trait: "forbidden"}
)
```

---

### New — create a structured error

Wraps a root cause with a trait, context, and an automatic call-site trace. Panics if `err` is `nil`.

```go
err := serror.New(errors.New("user does not exist"), ErrNotFound, map[string]any{
    "user_id": 42,
})
```

---

### Wrap — attach context after the fact

Merges additional key-value pairs into an existing `*ErrorRecord`. No-op on plain errors.

```go
serror.Wrap(map[string]any{"request_id": "abc-123"}, err)
```

---

### OfTrait — match by trait

Reports whether an error carries a specific trait. Matches on both `Code` and `Trait`.

```go
if serror.OfTrait(err, ErrNotFound) {
    // handle not-found
}
```

---

### E — log with slog

Returns an `slog.Attr` ready to pass to any `slog` call. For a plain `error`, falls back to `slog.Any("error", err)`. Returns an empty attr for `nil`.

```go
slog.Error("failed to fetch user", serror.E(err))
```

Produces:

```
level=ERROR msg="failed to fetch user" error.trait=not_found error.code=404 error.trace="main.getUser -> main.handleRequest" error.cause="user does not exist" user_id=42
```

---

## Full example

```go
var ErrNotFound = serror.ErrorTrait{Code: 404, Trait: "not_found"}

func getUser(id int) (*User, error) {
    u, err := db.Find(id)
    if err != nil {
        return nil, serror.New(err, ErrNotFound, map[string]any{"user_id": id})
    }
    return u, nil
}

func handleRequest(id int) {
    u, err := getUser(id)
    if err != nil {
        serror.Wrap(map[string]any{"request_id": requestID()}, err)

        if serror.OfTrait(err, ErrNotFound) {
            slog.Warn("user not found", serror.E(err))
            return
        }

        slog.Error("unexpected error", serror.E(err))
        return
    }

    _ = u
}
```
