// generated by smg -output misc/fixture/e/model_search.go misc/fixture/e; DO NOT EDIT

package e

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/favclip/smg/smgutils"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
	"strconv"
	"time"
)

// InventorySearch best match Search API wrapper for Inventory.
type InventorySearch struct {
	src *Inventory

	ID                 string
	ProductName        string
	Description        string
	DescriptionUnigram string
	DescriptionBigram  string
	Stock              float64
	Price              string
	Barcode            string
	AdminNames         string
	Shops              string
	CreatedAt          time.Time
	UpdatedAt          time.Time
	UpdatedAtUnixTime  float64
}

// Load by search.LoadStruct.
func (s *InventorySearch) Load(fields []search.Field, metadata *search.DocumentMetadata) error {
	return search.LoadStruct(s, fields)
}

// Save with search.DocumentMetadata#Rank.
func (s *InventorySearch) Save() ([]search.Field, *search.DocumentMetadata, error) {
	fields, err := search.SaveStruct(s)
	if err != nil {
		return nil, nil, err
	}
	// https://cloud.google.com/appengine/docs/go/search/reference#DocumentMetadata
	// 0 origin value can not be correctly sorted.
	// Typically, the data is assumed to be 0's origin, 1 added every time.
	metadata := &search.DocumentMetadata{Rank: int(s.Stock) + 1}

	return fields, metadata, nil
}

// Searchfy converts *Inventory to *InventorySearch.
func (src *Inventory) Searchfy() (*InventorySearch, error) {
	if src == nil {
		return nil, nil
	}
	dest := &InventorySearch{}
	dest.src = src
	var err error
	var b []byte
	dest.ID = strconv.FormatInt(src.ID, 10)
	dest.ProductName = src.ProductName

	dest.Description = src.Description
	dest.DescriptionUnigram, err = smgutils.UnigramForSearch(src.Description)
	if err != nil {
		return nil, err
	}
	dest.DescriptionBigram, err = smgutils.BigramForSearch(src.Description)
	if err != nil {
		return nil, err
	}
	dest.Stock = float64(src.Stock)
	dest.Price = strconv.Itoa(src.Price)
	dest.Barcode = strconv.FormatInt(src.Barcode, 10)

	b, err = json.Marshal(src.AdminNames)
	if err != nil {
		return nil, err
	}
	if str := string(b); str != "" && str != "\"\"" {
		dest.AdminNames = str
	}

	b, err = json.Marshal(src.Shops)
	if err != nil {
		return nil, err
	}
	if str := string(b); str != "" && str != "\"\"" {
		dest.Shops = str
	}
	dest.CreatedAt = src.CreatedAt

	// Number Field is value between -2,147,483,647 and 2,147,483,647.
	// https://cloud.google.com/appengine/docs/go/search/#Go_Documents_and_fields
	unixtime := src.UpdatedAt.Unix()
	if unixtime < -2147483647 {
		unixtime = -2147483647
	} else if 2147483647 < unixtime {
		unixtime = 2147483647
	}
	dest.UpdatedAtUnixTime = float64(unixtime)
	dest.UpdatedAt = src.UpdatedAt
	return dest, nil
}

// NewInventorySearch create new *InventorySearchBuilder.
func NewInventorySearch() *InventorySearchBuilder {
	op := &smgutils.Op{}
	b := &InventorySearchBuilder{
		rootOp:    op,
		currentOp: op,
	}
	b.ProductName = &InventorySearchStringPropertyInfo{"ProductName", b}
	b.Description = &InventorySearchNgramStringPropertyInfo{InventorySearchStringPropertyInfo{"Description", b}}
	b.Stock = &InventorySearchNumberPropertyInfo{"Stock", b}
	b.Price = &InventorySearchStringPropertyInfo{"Price", b}
	b.Barcode = &InventorySearchStringPropertyInfo{"Barcode", b}
	b.AdminNames = &InventorySearchStringPropertyInfo{"AdminNames", b}
	b.Shops = &InventorySearchStringPropertyInfo{"Shops", b}
	b.CreatedAt = &InventorySearchTimePropertyInfo{"CreatedAt", b}
	b.UpdatedAt = &InventorySearchUnixTimePropertyInfo{"UpdatedAt", b}

	return b
}

// InventorySearchBuilder builds Search API query.
type InventorySearchBuilder struct {
	rootOp      *smgutils.Op
	currentOp   *smgutils.Op // for grouping
	opts        *search.SearchOptions
	query       string
	index       *search.Index
	ProductName *InventorySearchStringPropertyInfo
	Description *InventorySearchNgramStringPropertyInfo
	Stock       *InventorySearchNumberPropertyInfo
	Price       *InventorySearchStringPropertyInfo
	Barcode     *InventorySearchStringPropertyInfo
	AdminNames  *InventorySearchStringPropertyInfo
	Shops       *InventorySearchStringPropertyInfo
	CreatedAt   *InventorySearchTimePropertyInfo
	UpdatedAt   *InventorySearchUnixTimePropertyInfo
}

// And append new operant to query.
func (b *InventorySearchBuilder) And() *InventorySearchBuilder {
	b.currentOp.Children = append(b.currentOp.Children, &smgutils.Op{Type: smgutils.And})
	return b
}

// Or append new operant to query.
func (b *InventorySearchBuilder) Or() *InventorySearchBuilder {
	b.currentOp.Children = append(b.currentOp.Children, &smgutils.Op{Type: smgutils.Or})
	return b
}

// Group append new operant to query.
func (b *InventorySearchBuilder) Group(p func()) *InventorySearchBuilder {
	b.StartGroup()
	p()
	b.EndGroup()
	return b
}

// StartGroup append new operant to query.
func (b *InventorySearchBuilder) StartGroup() *InventorySearchBuilder {
	op := &smgutils.Op{Type: smgutils.Group, Parent: b.currentOp}
	b.currentOp.Children = append(b.currentOp.Children, op)
	b.currentOp = op
	return b
}

// EndGroup append new operant to query.
func (b *InventorySearchBuilder) EndGroup() *InventorySearchBuilder {
	b.currentOp = b.currentOp.Parent
	return b
}

// Put document to Index.
func (b *InventorySearchBuilder) Put(c context.Context, src *Inventory) (string, error) {
	doc, err := src.Searchfy()
	if err != nil {
		return "", err
	}
	return b.PutDocument(c, doc)
}

// PutDocument to Index.
func (b *InventorySearchBuilder) PutDocument(c context.Context, src *InventorySearch) (string, error) {
	index, err := search.Open("Inventory")
	if err != nil {
		return "", err
	}

	docID := ""
	if v, ok := interface{}(src).(smgutils.DocIDer); ok { // TODO can I shorten this cond expression?
		docID, err = v.DocID(c)
		if err != nil {
			return "", err
		}
		src.ID = docID
	}

	log.Debugf(c, "id: %#v, payload: %#v", docID, src)

	docID, err = index.Put(c, docID, src)
	if err != nil {
		return "", err
	}

	src.ID = docID

	return docID, nil
}

// Delete document from Index.
func (b *InventorySearchBuilder) Delete(c context.Context, src *Inventory) error {
	doc, err := src.Searchfy()
	if err != nil {
		return err
	}
	return b.DeleteDocument(c, doc)
}

// DeleteDocument from Index.
func (b *InventorySearchBuilder) DeleteDocument(c context.Context, src *InventorySearch) error {
	if v, ok := interface{}(src).(smgutils.DocIDer); ok { // TODO can I shorten this cond expression?
		docID, err := v.DocID(c)
		if err != nil {
			return err
		}
		return b.DeleteByDocID(c, docID)
	}

	return errors.New("src is not implemented DocIDer interface")
}

// DeleteByDocID from Index.
func (b *InventorySearchBuilder) DeleteByDocID(c context.Context, docID string) error {
	index, err := search.Open("Inventory")
	if err != nil {
		return err
	}

	return index.Delete(c, docID)
}

// Opts returns *InventorySearchOptions.
func (b *InventorySearchBuilder) Opts() *InventorySearchOptions {
	return &InventorySearchOptions{b: b}
}

// Search returns *InventorySearchIterator, It is result from Index.
func (b *InventorySearchBuilder) Search(c context.Context) (*InventorySearchIterator, error) {
	index, err := search.Open("Inventory")
	if err != nil {
		return nil, err
	}
	b.index = index

	buffer := &bytes.Buffer{}
	err = b.rootOp.Query(buffer)
	if err != nil {
		return nil, err
	}
	b.query = buffer.String()
	log.Debugf(c, "query: '%s', opts: %#v", b.query, b.opts)
	iter := b.index.Search(c, b.query, b.opts)

	return &InventorySearchIterator{b, iter}, nil
}

// InventorySearchOptions construct *search.SearchOptions.
type InventorySearchOptions struct {
	b *InventorySearchBuilder
}

// Limit setup opts.
func (b *InventorySearchOptions) Limit(value int) *InventorySearchOptions {
	if b.b.opts == nil {
		b.b.opts = &search.SearchOptions{}
	}
	b.b.opts.Limit = value
	return b
}

// IDsOnly setup opts.
func (b *InventorySearchOptions) IDsOnly() *InventorySearchOptions {
	if b.b.opts == nil {
		b.b.opts = &search.SearchOptions{}
	}
	b.b.opts.IDsOnly = true
	return b
}

// Cursor setup opts.
func (b *InventorySearchOptions) Cursor(cursor search.Cursor) *InventorySearchOptions {
	b.b.opts.Cursor = cursor
	return b
}

// Offset setup opts.
func (b *InventorySearchOptions) Offset(value int) *InventorySearchOptions {
	b.b.opts.Offset = value
	return b
}

// InventorySearchIterator can access to search result.
type InventorySearchIterator struct {
	b    *InventorySearchBuilder
	iter *search.Iterator
}

// Next returns next document from iter.
func (b *InventorySearchIterator) Next(c context.Context) (string, *InventorySearch, error) {
	var s *InventorySearch
	if b.b.opts == nil || b.b.opts.IDsOnly != true {
		s = &InventorySearch{}
	}

	docID, err := b.iter.Next(s)
	if err != nil {
		return "", nil, err
	}
	if s != nil {
		s.ID = docID
	}

	return docID, s, err
}

// Cursor returns cursor of search.
func (b *InventorySearchIterator) Cursor() search.Cursor {
	return b.iter.Cursor()
}

// InventorySearchStringPropertyInfo hold property info.
type InventorySearchStringPropertyInfo struct {
	Name string
	b    *InventorySearchBuilder
}

// Match add query operand.
func (p *InventorySearchStringPropertyInfo) Match(value string) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Match, Value: value})
	return p.b
}

// Asc add query operand.
func (p *InventorySearchStringPropertyInfo) Asc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name,
		Reverse: true,
	})

	return p.b
}

// Desc add query operand.
func (p *InventorySearchStringPropertyInfo) Desc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name,
		Reverse: false,
	})

	return p.b
}

// InventorySearchNgramStringPropertyInfo hold property info.
type InventorySearchNgramStringPropertyInfo struct {
	InventorySearchStringPropertyInfo
}

// NgramMatch add query operand.
func (p *InventorySearchNgramStringPropertyInfo) NgramMatch(value string) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.NgramMatch, Value: value})
	return p.b
}

// InventorySearchNumberPropertyInfo hold property info.
type InventorySearchNumberPropertyInfo struct {
	Name string
	b    *InventorySearchBuilder
}

// IntGreaterThanOrEqual add query operand.
func (p *InventorySearchNumberPropertyInfo) IntGreaterThanOrEqual(value int) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.GtEq, Value: value})
	return p.b
}

// IntGreaterThan add query operand.
func (p *InventorySearchNumberPropertyInfo) IntGreaterThan(value int) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Gt, Value: value})
	return p.b
}

// IntLessThanOrEqual add query operand.
func (p *InventorySearchNumberPropertyInfo) IntLessThanOrEqual(value int) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.LtEq, Value: value})
	return p.b
}

// IntLessThan add query operand.
func (p *InventorySearchNumberPropertyInfo) IntLessThan(value int) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Lt, Value: value})
	return p.b
}

// IntEqual add query operand.
func (p *InventorySearchNumberPropertyInfo) IntEqual(value int) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Eq, Value: value})
	return p.b
}

// Int64GreaterThanOrEqual add query operand.
func (p *InventorySearchNumberPropertyInfo) Int64GreaterThanOrEqual(value int64) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.GtEq, Value: value})
	return p.b
}

// Int64GreaterThan add query operand.
func (p *InventorySearchNumberPropertyInfo) Int64GreaterThan(value int64) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Gt, Value: value})
	return p.b
}

// Int64LessThanOrEqual add query operand.
func (p *InventorySearchNumberPropertyInfo) Int64LessThanOrEqual(value int64) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.LtEq, Value: value})
	return p.b
}

// Int64LessThan add query operand.
func (p *InventorySearchNumberPropertyInfo) Int64LessThan(value int64) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Lt, Value: value})
	return p.b
}

// Int64Equal add query operand.
func (p *InventorySearchNumberPropertyInfo) Int64Equal(value int64) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Eq, Value: value})
	return p.b
}

// Asc add query operand.
func (p *InventorySearchNumberPropertyInfo) Asc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name,
		Reverse: true,
	})

	return p.b
}

// Desc add query operand.
func (p *InventorySearchNumberPropertyInfo) Desc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name,
		Reverse: false,
	})

	return p.b
}

// InventorySearchBoolPropertyInfo hold property info.
type InventorySearchBoolPropertyInfo struct {
	Name string
	b    *InventorySearchBuilder
}

// Equal add query operand.
func (p *InventorySearchNumberPropertyInfo) Equal(value bool) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Eq, Value: value})
	return p.b
}

// InventorySearchTimePropertyInfo hold property info.
type InventorySearchTimePropertyInfo struct {
	Name string
	b    *InventorySearchBuilder
}

// query spec for time.Time.
// https://cloud.google.com/appengine/docs/go/search/query_strings#Go_Queries_on_date_fields
// It using date, not datetime.

// GreaterThanOrEqual add query operand.
func (p *InventorySearchTimePropertyInfo) GreaterThanOrEqual(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.GtEq, Value: value.UTC()})
	return p.b
}

// GreaterThan add query operand.
func (p *InventorySearchTimePropertyInfo) GreaterThan(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Gt, Value: value.UTC()})
	return p.b
}

// LessThanOrEqual add query operand.
func (p *InventorySearchTimePropertyInfo) LessThanOrEqual(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.LtEq, Value: value.UTC()})
	return p.b
}

// LessThan add query operand.
func (p *InventorySearchTimePropertyInfo) LessThan(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Lt, Value: value.UTC()})
	return p.b
}

// Equal add query operand.
func (p *InventorySearchTimePropertyInfo) Equal(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Eq, Value: value.UTC()})
	return p.b
}

// Asc add query operand.
func (p *InventorySearchTimePropertyInfo) Asc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name,
		Reverse: true,
	})

	return p.b
}

// Desc add query operand.
func (p *InventorySearchTimePropertyInfo) Desc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name,
		Reverse: false,
	})

	return p.b
}

// InventorySearchUnixTimePropertyInfo hold property info.
type InventorySearchUnixTimePropertyInfo struct {
	Name string
	b    *InventorySearchBuilder
}

// GreaterThanOrEqual add query operand.
func (p *InventorySearchUnixTimePropertyInfo) GreaterThanOrEqual(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name + "UnixTime", Type: smgutils.GtEq, Value: value.Unix()})
	return p.b
}

// GreaterThan add query operand.
func (p *InventorySearchUnixTimePropertyInfo) GreaterThan(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name + "UnixTime", Type: smgutils.Gt, Value: value.Unix()})
	return p.b
}

// LessThanOrEqual add query operand.
func (p *InventorySearchUnixTimePropertyInfo) LessThanOrEqual(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name + "UnixTime", Type: smgutils.LtEq, Value: value.Unix()})
	return p.b
}

// LessThan add query operand.
func (p *InventorySearchUnixTimePropertyInfo) LessThan(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name + "UnixTime", Type: smgutils.Lt, Value: value.Unix()})
	return p.b
}

// Equal add query operand.
func (p *InventorySearchUnixTimePropertyInfo) Equal(value time.Time) *InventorySearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name + "UnixTime", Type: smgutils.Eq, Value: value.Unix()})
	return p.b
}

// Asc add query operand.
func (p *InventorySearchUnixTimePropertyInfo) Asc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name + "UnixTime",
		Reverse: true,
	})

	return p.b
}

// Desc add query operand.
func (p *InventorySearchUnixTimePropertyInfo) Desc() *InventorySearchBuilder {
	if p.b.opts == nil {
		p.b.opts = &search.SearchOptions{}
	}
	if p.b.opts.Sort == nil {
		p.b.opts.Sort = &search.SortOptions{}
	}
	p.b.opts.Sort.Expressions = append(p.b.opts.Sort.Expressions, search.SortExpression{
		Expr:    p.Name + "UnixTime",
		Reverse: false,
	})

	return p.b
}
