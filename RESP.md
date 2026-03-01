# RESP (REdis Serialization Protocol)

1. `Simple Strings`: are used for short strings that have minimal overhead, typically used for short successful status replies.
The string content cannot contain newlines (`\r` or `\n`).
2. `Errors`: are identical to `Simple Strings` but act as an exception. Client implementations will parse these differently, usually throwing an error in the application code.
3. `Integers`: are simply numbers represented as strings.
4. `Bulk Strings`: are used to send binary safe data, since they can contain any byte (even newlines), they use a prefixed length so a parser knows the exact number of bytes to read before expecting the final CRLF.
    - `Null Bulk Strings`: If a key doesn't exist (like running `GET` on a missing key), RESP uses a special Bulk String length of -1 to represent a null value.
    - `Null Format`: `$-1\r\n`
5. `Arrays`: collection of other RESP types. Arrays can contain mixed types  (e.g. an integer and a bulk string in the same array). The most important usecase of Arrays is that all client commands are sent to the server as an Array of Bulk Strings.
    - `Null Arrays`: Similar to bulk strings, you can represent a Null Array.
    - `Null Format`: `*-1\r\n`


| Data Type | Prefix | Format | Example |
| :-- | :-- | :-- | :-- |
| Simple Strings | `+` | `+<string>\r\n` | A standard success response like `+OK\r\n` |
| Errors | `-` | `-<error message>\r\n` | `-ERR unknown command 'foobar'\r\n` |
| Integers | `:` | `:<number>\r\n` | `:1000\r\n` |
| Bulk Strings | `$` | `$<length>\r\n<actual string data>\r\n` | The string "hello" is encoded as `$5\r\nhello\r\n` |
| Arrays | `*` | `*<number of elements>\r\n<element-1><element-2>...` | `*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n` | 

