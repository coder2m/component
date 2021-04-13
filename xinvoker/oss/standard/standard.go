package standard

import (
	"encoding/xml"
	"io"
	"time"
)

type optionType string

const (
	optionParam optionType = "HTTPParameter" // URL parameter
	optionHTTP  optionType = "HTTPHeader"    // HTTP header
	optionArg   optionType = "FuncArgument"  // Function argument
)

// DeleteObjectsResult defines result of DeleteObjects request
type DeleteObjectsResult struct {
	Space          string
	Local          string
	DeletedObjects []string // Deleted object key list
}

// ObjectProperties defines Objecct properties
type ObjectProperties struct {
	Name         xml.Name  `xml:"Contents"`
	Key          string    `xml:"Key"`          // Object key
	Type         string    `xml:"Type"`         // Object type
	Size         int64     `xml:"Size"`         // Object size
	ETag         string    `xml:"ETag"`         // Object ETag
	Owner        Owner     `xml:"Owner"`        // Object owner information
	LastModified time.Time `xml:"LastModified"` // Object last modified time
	StorageClass string    `xml:"StorageClass"` // Object storage class (Standard, IA, Archive)
}

// Owner defines OssBucket/Object's owner
type Owner struct {
	Name        xml.Name `xml:"Owner"`
	ID          string   `xml:"ID"`          // Owner ID
	DisplayName string   `xml:"DisplayName"` // Owner's display name
}

// ListObjectsResult defines the result from ListObjects request
type ListObjectsResult struct {
	XMLName        xml.Name
	Prefix         string             // The object prefix
	Marker         string             // The marker filter.
	MaxKeys        int                // Max keys to return
	Delimiter      string             // The delimiter for grouping objects' name
	IsTruncated    bool               // Flag indicates if all results are returned (when it's false)
	NextMarker     string             // The start point of the next query
	Objects        []ObjectProperties // Object list
	CommonPrefixes []string           // You can think of commonprefixes as "folders" whose names end with the delimiter
}

type (
	optionValue struct {
		Value interface{}
		Type  optionType
	}

	Option func(map[string]optionValue) error
)

type Oss interface {
	PutObject(dstPath string, reader io.Reader, options ...Option) error
	PutObjectFromFile(dstPath, srcPath string, options ...Option) error
	GetObject(dstPath string, options ...Option) ([]byte, error)
	GetObjectToFile(dstPath, srcPath string, options ...Option) error
	DeleteObject(dstPath string) error
	DeleteObjects(dstPaths []string, options ...Option) (DeleteObjectsResult, error)
	IsObjectExist(dstPath string) (bool, error)
	ListObjects(options ...Option) (ListObjectsResult, error)
	SignURL(dstPath string, method string, expiredInSec int64, options ...Option) (string, error)
}
