package probe_test

import (
	"github.com/dangrier/alien/pkg/probe"
	"net/http"
	"testing"
	"time"
)

var testSets = []struct {
	result *probe.Result
	filter probe.ResultFilter
	expect bool
}{
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
			},
		},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
			},
		},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
				probe.FilterResponseCode(404),
			},
		},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAny{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
				probe.FilterResponseCode(404),
			},
		},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{
			Member: probe.FilterResponseCode(200),
		},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{
			Member: probe.FilterResponseCode(200),
		},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
			},
		}},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
			},
		}},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
				probe.FilterResponseCode(404),
			},
		}},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAny{
			Members: []probe.ResultFilter{
				probe.FilterResponseCode(200),
				probe.FilterResponseCode(404),
			},
		}},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupNot{
			Member: probe.FilterResponseCode(200),
		}},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupNot{
			Member: probe.FilterResponseCode(200),
		}},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
			},
		},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
			},
		},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
				probe.FilterResponseContains("tomato"),
			},
		},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupAny{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
				probe.FilterResponseContains("tomato"),
			},
		},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{
			Member: probe.FilterResponseContains("potato"),
		},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{
			Member: probe.FilterResponseContains("potato"),
		},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
			},
		}},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
			},
		}},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAll{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
				probe.FilterResponseContains("tomato"),
			},
		}},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupAny{
			Members: []probe.ResultFilter{
				probe.FilterResponseContains("potato"),
				probe.FilterResponseContains("tomato"),
			},
		}},
		expect: false,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      200,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupNot{
			Member: probe.FilterResponseContains("potato"),
		}},
		expect: true,
	},
	{
		result: &probe.Result{
			Timestamp: time.Now(),
			Probe:     nil,
			Code:      404,
			Body:      "This is a body containing potato",
			Headers:   http.Header{},
			Error:     nil,
		},
		filter: probe.FilterGroupNot{Member: probe.FilterGroupNot{
			Member: probe.FilterResponseContains("potato"),
		}},
		expect: true,
	},
}

func TestFilterSet(t *testing.T) {
	for _, ts := range testSets {
		check := ts.filter.Check(ts.result)
		if !(check == ts.expect) {
			t.Fatalf("Failed (want '%t' got '%t'): code %d: %+v", ts.expect, check, ts.result.Code, ts.filter)
		}
	}
}
