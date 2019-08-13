# Handlers

## List of currently implemented Handlers

- Go
- Brew

## Development

To write new handlers, one has to implement the `PackageHandler` interface
defined in `pkg/controller/controller.go`.
After that, the handler can be registered in `app/cmd/cmd.go`.

The interface should be well enough described to get going. Some informal
conventions are listed here.

### Logging

As the controller has not parsed the package and only delegates the work,
the handler should be logging it's actions accordingly. This means:

For every action, on every package, log accordingly.

#### Example

Installing package X and Y should output something like:

```
output.Info("📦 [HandlerName]\tInstalling Package %s", pkg)
```

On a successful operation, log accordingly, too:

```
output.Success("📦 [HandlerName]\tUpgraded Package %s", pkg)
```

Depending on how long the handler name is, add one or two tabs

### Errors

Do not log errors if the function returns the error. Instead,
use `errors.Wrapf` to add context to the error.

If there is the possibility of receiving a command for multiple packages,
use the [collection](../collection/) package. It allows to collect multiple
errors, merge errors together and more. This way, try to work through all
the given packages before returning an error.

See the `goget` directory for an example handler.
