// Pager: Handles reading/writing fixed-size pages (4KB) to/from disk.
// This is the core abstraction for table row and index storage in letsgodb.
package storage

import (
	"fmt"
	"os"
)

const PageSize = 4096 // each page is 4kb

type Pager struct {
	file     *os.File          // file the stores data in Disk
	pages    map[uint32][]byte // in memory cache of page number
	pageSize int               // 4kb
	maxPage  int               //number of pages that exists in the disk
	pageNum  uint32            //to track the pages for allocation
}

/*
Note:

	Opens a file for read/write operations
	gets the file info like size and number of files in the disk
	returns the pager struct
	basically set up the Pager and
	it works like a constructor as it is called once
*/
func NewPager(filename string) *Pager {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(fmt.Sprintf("failed to open file: %v", err))
	}

	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	maxPage := int(fileSize / int64(PageSize))

	return &Pager{
		file:     file,
		pages:    make(map[uint32][]byte),
		pageSize: PageSize,
		maxPage:  maxPage,
	}
}

/*
ReadPage loads a page from disk into memory.
If the file is smaller than a full page, the remaining bytes are zeroed out.
This ensures every page is always PageSize bytes.
*/
func (p *Pager) ReadPage(pageNum int) ([]byte, error) {

	offset := int64(pageNum) * int64(p.pageSize)

	// Seek to the correct offset in the file
	_, err := p.file.Seek(offset, 0)
	if err != nil {
		return nil, err
	}

	// Allocate 4KB page buffer
	buf := make([]byte, p.pageSize)

	//the contents of the file is retrived and it fills the buf wiht the data
	//starting from the offset set by Seek
	n, err := p.file.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return nil, err
	}

	/*
		this block handles the case when the bytes read form the file doesn't
		fill the 4kb buffer. It fills the remaning space with 0 to avoid inconsistancy
		and partial page readings
	*/
	if n < p.pageSize {
		// Zero out the rest if the file is smaller than the page
		for i := n; i < p.pageSize; i++ {
			buf[i] = 0
		}
	}

	return buf, nil
}

/*
GetPage returns the page with the given page number.
If the page is not in the cache, it is loaded from disk and cached.
This is the main entry point for reading/modifying a page in memory.
*/
func (p *Pager) GetPage(pageNum uint32) []byte {
	//if the page is already in cache it returns immediately since no need to
	// extract from disk
	if page, ok := p.pages[pageNum]; ok {
		return page
	}

	//if it is still in disk it class the ReadPage function to do the thing
	// ReadPage does, to store data in bytes in page variable
	page, err := p.ReadPage(int(pageNum))
	if err != nil {
		panic(err)
	}
	//The page is stored in the pages map so future requests for the same
	//page can be served from memory (faster).
	p.pages[pageNum] = page
	return page
}

// WritePage writes the given data to the specified page number on disk.
// The data slice must be exactly PageSize bytes.
func (p *Pager) WritePage(pageNum int, data []byte) error {
	if len(data) > PageSize {
		return fmt.Errorf("data exceeds page size")
	}
	_, err := p.file.WriteAt(data, int64(pageNum)*int64(PageSize))
	return err
}

/*
FlushPage writes the given page data to disk at the specified page number.
This is a wrapper around WritePage for convenience.
*/
func (p *Pager) FlushPage(pageNum uint32, data []byte) error {
	return p.WritePage(int(pageNum), data)
}

/*
AllocatePage allocates a new blank page in memory, assigns it the next available page number,
and returns that page number. The new page is also cached in memory.
This should be called whenever a new row or B+Tree node needs to be stored.
*/
func (p *Pager) AllocatePage() uint32 {
	pageNum := uint32(p.maxPage)
	p.maxPage++ // Increment the page count for the next allocation
	p.pages[pageNum] = make([]byte, PageSize)
	return uint32(pageNum)
}

/*
File returns the internal file pointer used by the Pager for disk I/O.
This can be used for advanced operations or for closing the file.
*/
func (p *Pager) File() *os.File {
	return p.file
}
