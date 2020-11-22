## Writing a Log Package

### Terms

The following terms are being used in the design and code:

| term | description |
| --- | --- |
| `Record` | The data stored in the log |
| `Store` | The file that stores the records. |
| `Index` | The file that stores the index entries. |
| `Segment` | The abstraction that ties a store and an index together. |
| `Log` | The abstraction that ties all the segments together |

