# urDOS

## Version 2 design
 * Every widget can draw itself and its children
 * A widget's position and size is determined by its parent
 * A widget does not know its parent
 * Every widget has (hopefully immutable) reference to a theme which decides style, passed on by its parent by default
 * A widget may also have some styles overridden
 * A widget is told if it has focus via the constructor
 * Any changes to appearance or state must call `Rebuild()` on the widget changing state
 * All widgets derive from the Widget type
 * Focus: only a SetFocus method on all Widgets, parents set child's focuses
 * All Widgets have a HandleMouseEvent and HandleKeyEvent

### Widgets composed of embedded structs
```go
ui.Focusable -> Focused bool field
```

```go
// Throws the MainWidget into an event loop until that loop is exited
ui.RunApp(ui.App{
	Name: "My App",
	MainWidget: App,
})
```
