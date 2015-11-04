// generated by smg -type Sample -output misc/fixture/c/model_search.go misc/fixture/c; DO NOT EDIT

package c

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/favclip/smg/smgutils"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/search"
)

// for Sample
type SampleSearch struct {
	src *Sample

	A string
	B string
}

func (src *Sample) Searchfy() (*SampleSearch, error) {
	if src == nil {
		return nil, nil
	}
	dest := &SampleSearch{}
	dest.src = src
	var err error
	var b []byte
	dest.A = src.A

	b, err = json.Marshal(src.B)
	if err != nil {
		return nil, err
	}
	if str := string(b); str != "" && str != "\"\"" {
		dest.B = str
	}
	return dest, nil
}

func NewSampleSearch() *SampleSearchBuilder {
	op := &smgutils.Op{}
	b := &SampleSearchBuilder{
		rootOp:    op,
		currentOp: op,
	}
	b.A = &SampleSearchStringPropertyInfo{"A", b}
	b.B = &SampleSearchStringPropertyInfo{"B", b}

	return b
}

type SampleSearchBuilder struct {
	rootOp    *smgutils.Op
	currentOp *smgutils.Op // for grouping
	opts      *search.SearchOptions
	query     string
	index     *search.Index
	A         *SampleSearchStringPropertyInfo
	B         *SampleSearchStringPropertyInfo
}

func (b *SampleSearchBuilder) And() *SampleSearchBuilder {
	b.currentOp.Children = append(b.currentOp.Children, &smgutils.Op{Type: smgutils.And})
	return b
}

func (b *SampleSearchBuilder) Or() *SampleSearchBuilder {
	b.currentOp.Children = append(b.currentOp.Children, &smgutils.Op{Type: smgutils.Or})
	return b
}

func (b *SampleSearchBuilder) Group(p func()) *SampleSearchBuilder {
	b.StartGroup()
	p()
	b.EndGroup()
	return b
}

func (b *SampleSearchBuilder) StartGroup() *SampleSearchBuilder {
	op := &smgutils.Op{Type: smgutils.Group, Parent: b.currentOp}
	b.currentOp.Children = append(b.currentOp.Children, op)
	b.currentOp = op
	return b
}

func (b *SampleSearchBuilder) EndGroup() *SampleSearchBuilder {
	b.currentOp = b.currentOp.Parent
	return b
}

func (b *SampleSearchBuilder) Put(c context.Context, src *Sample) (string, error) {
	doc, err := src.Searchfy()
	if err != nil {
		return "", err
	}
	return b.PutDocument(c, doc)
}

func (b *SampleSearchBuilder) PutDocument(c context.Context, src *SampleSearch) (string, error) {
	index, err := search.Open("Sample")
	if err != nil {
		return "", err
	}

	docID := ""
	if v, ok := interface{}(src).(smgutils.DocIDer); ok { // TODO can I shorten this cond expression?
		docID, err = v.DocID(c)
		if err != nil {
			return "", err
		}

	}

	log.Debugf(c, "id: %#v, payload: %#v", docID, src)

	docID, err = index.Put(c, docID, src)
	if err != nil {
		return "", err
	}

	return docID, nil
}

func (b *SampleSearchBuilder) Delete(c context.Context, src *Sample) error {
	doc, err := src.Searchfy()
	if err != nil {
		return err
	}
	return b.DeleteDocument(c, doc)
}

func (b *SampleSearchBuilder) DeleteDocument(c context.Context, src *SampleSearch) error {
	if v, ok := interface{}(src).(smgutils.DocIDer); ok { // TODO can I shorten this cond expression?
		docID, err := v.DocID(c)
		if err != nil {
			return err
		}
		return b.DeleteByDocID(c, docID)
	}

	return errors.New("src is not implemented DocIDer interface")
}

func (b *SampleSearchBuilder) DeleteByDocID(c context.Context, docID string) error {
	index, err := search.Open("Sample")
	if err != nil {
		return err
	}

	return index.Delete(c, docID)
}

func (b *SampleSearchBuilder) Opts() *SampleSearchOptions {
	return &SampleSearchOptions{b: b}
}

func (b *SampleSearchBuilder) Search(c context.Context) (*SampleSearchIterator, error) {
	index, err := search.Open("Sample")
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

	return &SampleSearchIterator{b, iter}, nil
}

type SampleSearchOptions struct {
	b *SampleSearchBuilder
}

func (b *SampleSearchOptions) Limit(value int) *SampleSearchOptions {
	if b.b.opts == nil {
		b.b.opts = &search.SearchOptions{}
	}
	b.b.opts.Limit = value
	return b
}

func (b *SampleSearchOptions) IDsOnly() *SampleSearchOptions {
	if b.b.opts == nil {
		b.b.opts = &search.SearchOptions{}
	}
	b.b.opts.IDsOnly = true
	return b
}

type SampleSearchIterator struct {
	b    *SampleSearchBuilder
	iter *search.Iterator
}

func (b *SampleSearchIterator) Next(c context.Context) (string, *SampleSearch, error) {
	var s *SampleSearch
	if b.b.opts == nil || b.b.opts.IDsOnly != true {
		s = &SampleSearch{}
	}

	docID, err := b.iter.Next(s)
	if err != nil {
		return "", nil, err
	}

	return docID, s, err
}

type SampleSearchStringPropertyInfo struct {
	Name string
	b    *SampleSearchBuilder
}

func (p *SampleSearchStringPropertyInfo) Match(value string) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Match, Value: value})
	return p.b
}

func (p *SampleSearchStringPropertyInfo) Asc() *SampleSearchBuilder {
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

func (p *SampleSearchStringPropertyInfo) Desc() *SampleSearchBuilder {
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

type SampleSearchNgramStringPropertyInfo struct {
	SampleSearchStringPropertyInfo
}

func (p *SampleSearchNgramStringPropertyInfo) NgramMatch(value string) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.NgramMatch, Value: value})
	return p.b
}

type SampleSearchNumberPropertyInfo struct {
	Name string
	b    *SampleSearchBuilder
}

func (p *SampleSearchNumberPropertyInfo) IntGreaterThanOrEqual(value int) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.GtEq, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) IntGreaterThan(value int) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Gt, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) IntLessThanOrEqual(value int) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.LtEq, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) IntLessThan(value int) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Lt, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) IntEqual(value int) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Eq, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) Int64GreaterThanOrEqual(value int64) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.GtEq, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) Int64GreaterThan(value int64) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Gt, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) Int64LessThanOrEqual(value int64) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.LtEq, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) Int64LessThan(value int64) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Lt, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) Int64Equal(value int64) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Eq, Value: value})
	return p.b
}

func (p *SampleSearchNumberPropertyInfo) Asc() *SampleSearchBuilder {
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

func (p *SampleSearchNumberPropertyInfo) Desc() *SampleSearchBuilder {
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

type SampleSearchBoolPropertyInfo struct {
	Name string
	b    *SampleSearchBuilder
}

func (p *SampleSearchNumberPropertyInfo) Equal(value bool) *SampleSearchBuilder {
	p.b.currentOp.Children = append(p.b.currentOp.Children, &smgutils.Op{FieldName: p.Name, Type: smgutils.Eq, Value: value})
	return p.b
}

type SampleSearchTimePropertyInfo struct {
	Name string
	b    *SampleSearchBuilder
}

func (p *SampleSearchTimePropertyInfo) Asc() *SampleSearchBuilder {
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

func (p *SampleSearchTimePropertyInfo) Desc() *SampleSearchBuilder {
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