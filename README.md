#Banners

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


##Storage mechanisms

###Lock storage

```go
type LockStorage struct {
	sync.RWMutex
	categories map[string][]*banners.Banner
	count      int64
}
```

On input requests lock `categories` map.
If display count of banner is zero, removed it from all categories.


###Chan storage

```go
type ChanStorage struct {
	categories map[string][]*banners.Banner
	count      int64

	chAppend    chan (*banners.Banner)
	chLookupIn  chan ([]string)
	chLookupOut chan (*banners.Banner)
}
```
It is storage look like a `LockStorage` but does not contains sync.Mutex(RWMutes) and does not required locks.
On create new instance of `ChanStorage` it start on background channels listener, input requests passed values to specify channels.


###Slice storage

```go
type SliceStorage struct {
	list []*banners.Banner
}
```

It is very simple storage - all banners are contained in one slice, without categorizing by categories.
To find a random element from slice, generate random number less then length of slice and lookup at first from this number to end, and then from start to this number.


##Benchmarks

	BenchmarkLockStorage-8     	10000000	       172 ns/op
	BenchmarkLockStorage2-8    	10000000	       188 ns/op

	BenchmarkChanStorage-8     	 1000000	      2748 ns/op
	BenchmarkChanStorage2-8    	  500000	      2717 ns/op

	BenchmarkSliceStorage-8    	10000000	       182 ns/op
	BenchmarkSliceStorage2-8   	20000000	        74.2 ns/op
