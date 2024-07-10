package main

import (
	"reflect"
	"testing"
)

func TestParseFilters(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected [][]Filter
	}{
		{
			name:  "Single filter",
			input: `[{"key":"age","field":"age","value":"30","relation":"="}]`,
			expected: [][]Filter{
				{
					{Key: "age", Field: "age", Value: "30", Relation: "="},
				},
			},
		},
		{
			name:  "Double filter",
			input: `[{"key":"age","field":"age","value":"30","relation":"="},{"key":"name","field":"name","value":"John","relation":"="}]`,
			expected: [][]Filter{
				{
					{Key: "age", Field: "age", Value: "30", Relation: "="},
					{Key: "name", Field: "name", Value: "John", Relation: "="},
				},
			},
		},
		{
			name:  "Multiple filters with OR",
			input: `[{"key":"age","field":"age","value":"30","relation":"="}, {"operator":"OR"}, {"key":"name","field":"name","value":"John","relation":"="}]`,
			expected: [][]Filter{
				{
					{Key: "age", Field: "age", Value: "30", Relation: "="},
				},
				{
					{Key: "name", Field: "name", Value: "John", Relation: "="},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseFilters(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestParseUserTags(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "Simple key-value pairs",
			input:    "{age: 30, name: John}",
			expected: map[string]string{"age": "30", "name": "John"},
		},
		{
			name:     "Empty input",
			input:    "{}",
			expected: map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseUserTags(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestParseUserLanguage(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected map[string]string
	}{
		{
			name:     "Simple language matching",
			input:    "{language: en}",
			expected: map[string]string{"language": "en"},
		},
		{
			name:     "Empty input",
			input:    "{}",
			expected: map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := parseUserTags(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestCheckConditionGroup(t *testing.T) {
	testCases := []struct {
		name         string
		group        []Filter
		userTags     map[string]string
		userLanguage string
		expected     bool
	}{
		{
			name: "Tag condition met",
			group: []Filter{
				{Key: "age", Field: "age", Value: "30", Relation: "="},
			},
			userTags: map[string]string{"age": "30"},
			expected: true,
		},
		{
			name: "Language condition met",
			group: []Filter{
				{Field: "language", Value: "en", Relation: "="},
			},
			userLanguage: "en",
			expected:     true,
		},
		{
			name: "Multiple conditions met",
			group: []Filter{
				{Key: "age", Field: "age", Value: "30", Relation: "="},
				{Key: "name", Field: "name", Value: "John", Relation: "="},
			},
			userTags: map[string]string{"age": "30", "name": "John"},
			expected: true,
		},
		{
			name: "One condition not met",
			group: []Filter{
				{Field: "language", Value: "en", Relation: "="},
				{Key: "name", Field: "name", Value: "John", Relation: "="},
			},
			userTags:     map[string]string{"age": "30", "name": "John"},
			userLanguage: "es",
			expected:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checkConditionGroup(tc.group, tc.userTags, tc.userLanguage)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestRelation(t *testing.T) {
	testCases := []struct {
		name         string
		group        []Filter
		userTags     map[string]string
		userLanguage string
		expected     bool
	}{
		{
			name: "Equals condition met",
			group: []Filter{
				{Key: "age", Field: "age", Value: "30", Relation: "="},
			},
			userTags: map[string]string{"age": "30"},
			expected: true,
		},
		{
			name: "Not Equals condition not met",
			group: []Filter{
				{Key: "age", Field: "age", Value: "30", Relation: "!="},
			},
			userTags: map[string]string{"age": "30", "name": "John"},
			expected: false,
		},
		{
			name: "Exists condition met",
			group: []Filter{
				{Key: "exists", Field: "tag", Relation: "exists"},
			},
			userTags: map[string]string{"age": "30", "name": "Jane", "exists": "1"},
			expected: true,
		},
		{
			name: "Not Exists condition met",
			group: []Filter{
				{Key: "not_exists", Field: "tag", Relation: "not_exists"},
			},
			userTags: map[string]string{"age": "30", "name": "Jane", "exists": "1"},
			expected: true,
		},
		{
			name: "Less than condition met",
			group: []Filter{
				{Key: "lt", Value: "9", Field: "tag", Relation: "<"},
			},
			userTags: map[string]string{"age": "30", "name": "Jane", "exists": "1", "lt": "8"},
			expected: true,
		},
		{
			name: "Greater than condition not met",
			group: []Filter{
				{Key: "lt", Value: "9", Field: "tag", Relation: ">"},
			},
			userTags: map[string]string{"age": "30", "name": "Jane", "exists": "1", "lt": "8"},
			expected: false,
		},
		{
			name: "Complex comparison",
			group: []Filter{
				{Value: "es", Field: "language", Relation: "!="},
				{Key: "gt", Value: "9", Field: "tag", Relation: ">"},
				{Key: "lt", Value: "9", Field: "tag", Relation: "<"},
				{Key: "exist", Field: "tag", Relation: "exists"},
				{Key: "nope", Field: "tag", Relation: "not_exists"},
				{Key: "atbat", Field: "tag", Value: "true", Relation: "="},
				{Key: "notatbat", Field: "tag", Value: "true", Relation: "!="},
			},
			userTags:     map[string]string{"gt": "30", "exist": "1", "lt": "8", "atbat": "true", "notatbat": "false"},
			userLanguage: "en",
			expected:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := checkConditionGroup(tc.group, tc.userTags, tc.userLanguage)
			if result != tc.expected {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
