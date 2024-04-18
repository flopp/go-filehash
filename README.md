# go-filehash
Combute hash/checksum of files, insert file hash/checksum into file name while copying

```
import "filehash"

hash, err := filehash.Compute("myfile.txt")
// Returns "95c39c37ef89acb2", nil

newName, err := filehash.Copy("myfile.txt", "out/myfile-HASH.txt", "HASH")
// Copies "myfile.txt" to "out/myfile-95c39c37ef89acb2.txt"
// Returns "out/myfile-95c39c37ef89acb2.txt", nil
```
