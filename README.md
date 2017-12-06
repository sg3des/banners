# Banners

Banner structure:

```go
type Banner struct {
	//URL to banner
	URL string

	//Count of prepaid shows
	Count int64

	//Available categories
	Categories []string

	//Lock flag 1 or 0
	Lock int32
}
```


## Storage mechanisms

### Lock storage

```go
type LockStorage struct {
	sync.RWMutex
	categories map[string][]*banners.Banner
	count      int64
}
```

On input requests lock `categories` map.
If display count of banner is zero, removed it from all categories.


### Tree storage

```go
type TreeStorage struct {
	root *branch
}

type branch struct {
	branches [256]*branch
	leafs    []leaf
}

type leaf *banners.Banner
```
This storage type does not required locking.
It is not binary tree!
Banners it is leafs.
Each branch may have 256 other branches by character.
Each branch contains all leafs from child branches.
Banners, with over display count, not removed. 

Appending banner by category, (where category is key), each byte in a key added new branch to tree. Looking for banner by categories, performing by same method.


### Slice storage

```go
type SliceStorage struct {
	list []*banners.Banner
}
```

It is very simple storage - all banners are contained in one slice, without categorizing by categories.
To find a random element from slice, generate random number less then length of slice and lookup at first from this number to end, and then from start to this number.


## Benchmarks

	BenchmarkLockStorage-8      	10000000	       153 ns/op
	BenchmarkLockStorageAny-8    	10000000	       194 ns/op

	BenchmarkTreeStorage-8      	 3000000	       386 ns/op
	BenchmarkTreeStorageAny-8   	10000000	       220 ns/op

	BenchmarkSliceStorage-8     	 5000000	       201 ns/op
	BenchmarkSliceStorageAny-8   	20000000	        97.5 ns/op

