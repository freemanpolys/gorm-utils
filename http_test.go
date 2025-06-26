package gorm

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestFromRequest_INOperator(t *testing.T) {
	u := &url.URL{
		RawQuery: "page=2&limit=5&sort=code desc&filter[code]=in,D42,L12&filter[price]=>,150",
	}
	r := &http.Request{URL: u}
	pagination, filters := FromRequest(r)

	// Check pagination
	if pagination.Page != 2 {
		t.Errorf("expected page 2, got %d", pagination.Page)
	}
	if pagination.Limit != 5 {
		t.Errorf("expected limit 5, got %d", pagination.Limit)
	}
	if pagination.Sort != "code desc" {
		t.Errorf("expected sort 'code desc', got '%s'", pagination.Sort)
	}

	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}

	if filters[0].Field != "code" || filters[0].Operator != "in" || !reflect.DeepEqual(filters[0].Value, []string{"D42", "L12"}) {
		t.Errorf("expected IN filter for code, got %+v", filters[0])
	}

	if filters[1].Field != "price" || filters[1].Operator != ">" || filters[1].Value != "150" {
		t.Errorf("expected price filter with operator '>' and value '150', got %+v", filters[1])
	}
}

func TestFromRequest_OtherOperators(t *testing.T) {
	u := &url.URL{
		RawQuery: "page=1&limit=10&sort=qty asc&filter[name]=like,foo&filter[qty]=>=,10",
	}
	r := &http.Request{URL: u}
	pagination, filters := FromRequest(r)

	// Check pagination
	if pagination.Page != 1 {
		t.Errorf("expected page 1, got %d", pagination.Page)
	}
	if pagination.Limit != 10 {
		t.Errorf("expected limit 10, got %d", pagination.Limit)
	}
	if pagination.Sort != "qty asc" {
		t.Errorf("expected sort 'qty asc', got '%s'", pagination.Sort)
	}

	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}

	if filters[0].Field != "name" || filters[0].Operator != "like" || filters[0].Value != "foo" {
		t.Errorf("expected like filter for name, got %+v", filters[0])
	}

	if filters[1].Field != "qty" || filters[1].Operator != ">=" || filters[1].Value != "10" {
		t.Errorf("expected >= filter for qty, got %+v", filters[1])
	}
}
