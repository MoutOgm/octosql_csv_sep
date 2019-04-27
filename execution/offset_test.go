package execution

import (
	"github.com/cube2222/octosql"
	"testing"
)

func TestOffset_Get(t *testing.T) {
	const NO_ERROR = ""

	tests := []struct {
		name       string
		vars       octosql.Variables
		node       *Offset
		wantStream *InMemoryStream
		wantError  string
	}{
		{
			name:       "negative offset value",
			vars:       octosql.NoVariables(),
			node:       NewOffset(UtilNewDummyNode(nil), UtilNewDummyValue(-42)),
			wantStream: nil,
			wantError:  "negative offset value",
		},
		{
			name:       "offset value not int",
			vars:       octosql.NoVariables(),
			node:       NewOffset(UtilNewDummyNode(nil), UtilNewDummyValue(2.0)),
			wantStream: nil,
			wantError:  "offset value not int",
		},
		{
			name: "normal offset get",
			vars: octosql.NoVariables(),
			node: NewOffset(&DummyNode{
				[]*Record{
					UtilNewRecord(
						[]octosql.VariableName{
							"num",
						},
						[]interface{}{
							1e10,
						}),
					UtilNewRecord(
						[]octosql.VariableName{
							"num",
						},
						[]interface{}{
							3.21,
						}),
					UtilNewRecord(
						[]octosql.VariableName{
							"flag",
						},
						[]interface{}{
							false,
						}),
					UtilNewRecord(
						[]octosql.VariableName{
							"num",
						},
						[]interface{}{
							2.23e7,
						}),
				},
			}, &DummyValue{2}),
			wantStream: NewInMemoryStream([]*Record{
				UtilNewRecord(
					[]octosql.VariableName{
						"flag",
					},
					[]interface{}{
						false,
					}),
				UtilNewRecord(
					[]octosql.VariableName{
						"num",
					},
					[]interface{}{
						2.23e7,
					}),
			}),
			wantError: NO_ERROR,
		},
		{
			name: "offset bigger than number of rows",
			vars: octosql.NoVariables(),
			node: NewOffset(&DummyNode{
				[]*Record{
					UtilNewRecord(
						[]octosql.VariableName{
							"num",
						},
						[]interface{}{
							1,
						}),
					UtilNewRecord(
						[]octosql.VariableName{
							"num",
						},
						[]interface{}{
							2,
						}),
				},
			}, &DummyValue{4}),
			wantStream: NewInMemoryStream([]*Record{
			}),
			wantError: NO_ERROR,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs, err := tt.node.Get(tt.vars)

			if (err == nil) != (tt.wantError == NO_ERROR) {
				t.Errorf("exactly one of test.wantError, tt.node.Get() is not nil")
				return
			}

			if err != nil {
				if err.Error() != tt.wantError {
					t.Errorf("Unexpected error %v, wanted: %v", err.Error(), tt.wantError)
				}
				return
			}

			equal, err := AreStreamsEqual(rs, tt.wantStream)
			if !equal {
				t.Errorf("limitedStream doesn't work as expected")
			}
			if err != nil {
				t.Errorf("limitedStream comparison error: %v", err)
			}
		})
	}
}