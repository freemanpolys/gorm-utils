package gorm

import (
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestFromRequest_INOperator(t *testing.T) {
	u := &url.URL{
		RawQuery: "page=2&limit=5&sort=code desc&filter=code::in::D42||L12~~price::gt::150",
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

	if filters[1].Field != "price" || filters[1].Operator != "gt" || filters[1].Value != "150" {
		t.Errorf("expected price filter with operator 'gt' and value '150', got %+v", filters[1])
	}
}

func TestFromRequest_OtherOperators(t *testing.T) {
	u := &url.URL{
		RawQuery: "page=1&limit=10&sort=qty asc&filter=name::like::foo~~qty::gte::10~~price::between::100||200",
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

	if len(filters) != 3 {
		t.Fatalf("expected 3 filters, got %d", len(filters))
	}

	if filters[0].Field != "name" || filters[0].Operator != "like" || filters[0].Value != "foo" {
		t.Errorf("expected like filter for name, got %+v", filters[0])
	}

	if filters[1].Field != "qty" || filters[1].Operator != "gte" || filters[1].Value != "10" {
		t.Errorf("expected gte filter for qty, got %+v", filters[1])
	}

	if filters[2].Field != "price" || filters[2].Operator != "between" || !reflect.DeepEqual(filters[2].Value, []string{"100", "200"}) {
		t.Errorf("expected between filter for price, got %+v", filters[2])
	}
}

func TestFromRequest_ComplexValues(t *testing.T) {
	u := &url.URL{
		RawQuery: "filter=product_name::eq::Product, Awesome: New|Version~~description::contains::This item has no issues.",
	}
	r := &http.Request{URL: u}
	_, filters := FromRequest(r)

	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}

	if filters[0].Field != "product_name" || filters[0].Operator != "eq" || filters[0].Value != "Product, Awesome: New|Version" {
		t.Errorf("expected eq filter for product_name, got %+v", filters[0])
	}

	if filters[1].Field != "description" || filters[1].Operator != "contains" || filters[1].Value != "This item has no issues." {
		t.Errorf("expected contains filter for description, got %+v", filters[1])
	}
}

func TestFromRequest_BetweenOperator(t *testing.T) {
	u := &url.URL{
		RawQuery: "filter=price::between::100||200~~id::between::1||10",
	}
	r := &http.Request{URL: u}
	_, filters := FromRequest(r)

	if len(filters) != 2 {
		t.Fatalf("expected 2 filters, got %d", len(filters))
	}

	if filters[0].Field != "price" || filters[0].Operator != "between" || !reflect.DeepEqual(filters[0].Value, []string{"100", "200"}) {
		t.Errorf("expected between filter for price, got %+v", filters[0])
	}

	if filters[1].Field != "id" || filters[1].Operator != "between" || !reflect.DeepEqual(filters[1].Value, []string{"1", "10"}) {
		t.Errorf("expected between filter for id, got %+v", filters[1])
	}
}

func TestFromRequest_MixedOperators(t *testing.T) {
	u := &url.URL{
		RawQuery: "filter=code::in::A||B||C~~price::between::10||20~~name::like::foo",
	}
	r := &http.Request{URL: u}
	_, filters := FromRequest(r)

	if len(filters) != 3 {
		t.Fatalf("expected 3 filters, got %d", len(filters))
	}

	if filters[0].Field != "code" || filters[0].Operator != "in" || !reflect.DeepEqual(filters[0].Value, []string{"A", "B", "C"}) {
		t.Errorf("expected in filter for code, got %+v", filters[0])
	}

	if filters[1].Field != "price" || filters[1].Operator != "between" || !reflect.DeepEqual(filters[1].Value, []string{"10", "20"}) {
		t.Errorf("expected between filter for price, got %+v", filters[1])
	}

	if filters[2].Field != "name" || filters[2].Operator != "like" || filters[2].Value != "foo" {
		t.Errorf("expected like filter for name, got %+v", filters[2])
	}
}
