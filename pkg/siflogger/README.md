## Tendermint logger interface
This interface is used internally in tendermint as well as in all cosmos modules
```
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	With(keyvals ...interface{}) Logger
}
```

## To create
```
import github.com/Sifchain/sifnode/pkg/siflogger
logger := siflogger.New(siflogger.TDFmt) # for standard formatted messages
logger := siflogger.New(siflogger.JSON) # for JSON formatted messages
```

## Layers
Each layer consist of key/value which are used as postfix in message of every logging functions *Debug()*, *Info()*, *Error()*.
There is a chain of layers, so each new layer inherits key/value of the all previous ones.
There is a function *With()* which is used to create new layer. For example:
```
logger2 = logger.With("key1", "value1")
logger3 = logger2.With("key2", "value2")
logger4 = logger.With("key3", "value3")
```

So the postfixes which will be printed in the message of **logger2** are **"key1=value1**", for **logger3** are **"key1=value1 key2=value2"**, for **logger4** are **"key3=value3"**. 
Representation of these postfixes is dependent on the implementation (siflogger.TDFmt or siflogger.JSON above)

## Filters
Layers are used to filter by keys/values as well.

#### Levels which are used in filters:
```
const (
	Debug Level = iota
	Info
	Error
	None
)
```

#### Functions to create filters:
```
func (e *Logger) SetGlobalFilter(level Level)
```
Set global level filter for all layers which have not been set explicitly

```
func (e *Logger) SetFilterForLayer(level Level, keyvals ...interface{})
```
Set filter for specific layers explicitly using key/value.
Number parameters **must be** even (key + value) because of the nature of logging of tendermint internals.

#### Filter sequence:
* Each SetFilterForLayer() call is equivalent of AND operator
* Each pair of arguments in SetFilterForLayer() call is equivalent of OR operator
* Filter applies only on the last value of the same key in layers chain, for example:
```
logger2 = logger.With("key", "value1")
logger3 = logger2.With("key", "value2")
```
Filter of **logger3** will apply only on the last one **"key=value2**" while the full chain of key/value is stored and printed in logger as **"key=value1 key=value2"**.
*Key **"module"** is exceptional and it is printed as prefix always having only the last value*

#### Config:
There is original cosmos config for filter configuration in file:
```
$HOME/.sifnoded/config/config.toml
```
in variable **log_level**.
The format is the following:
```
<module name>:<info|debug|error|none>,...
```
or:
```
*:<info|debug|error|none>
```
for global filter. This configuration allows to filter all logs by modules and to be more precise by anything created as:
```
logger.With("module", <module name>)
```