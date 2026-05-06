package provider

import (
    "context"
    "encoding/json"
    "fmt"
    "reflect"

    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/types/basetypes"
    "github.com/hashicorp/terraform-plugin-go/tftypes"
)

// JSONSubsetType is a Plugin Framework string type for complex-JSON fields.
// Its semantic-equality returns true when the planned JSON is a structural
// subset of the actual JSON: same shape, same scalars at every leaf, but the
// actual side may have additional keys filled in by the server (e.g. defaults).
// This prevents "Provider produced inconsistent result after apply" without
// having to teach the provider about every server-side default.
type JSONSubsetType struct {
    basetypes.StringType
}

var _ basetypes.StringTypable = JSONSubsetType{}

func (t JSONSubsetType) Equal(o attr.Type) bool {
    other, ok := o.(JSONSubsetType)
    if !ok {
        return false
    }
    return t.StringType.Equal(other.StringType)
}

func (t JSONSubsetType) String() string {
    return "JSONSubsetType"
}

func (t JSONSubsetType) ValueType(_ context.Context) attr.Value {
    return JSONSubsetValue{}
}

func (t JSONSubsetType) ValueFromString(_ context.Context, in basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
    return JSONSubsetValue{StringValue: in}, nil
}

func (t JSONSubsetType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
    val, err := t.StringType.ValueFromTerraform(ctx, in)
    if err != nil {
        return nil, err
    }
    sv, ok := val.(basetypes.StringValue)
    if !ok {
        return nil, fmt.Errorf("unexpected base value type: %T", val)
    }
    return JSONSubsetValue{StringValue: sv}, nil
}

// JSONSubsetValue is the value half of JSONSubsetType. It embeds the standard
// StringValue, so all the usual accessors (ValueString, IsNull, IsUnknown)
// keep working without change at call sites.
type JSONSubsetValue struct {
    basetypes.StringValue
}

var _ basetypes.StringValuableWithSemanticEquals = JSONSubsetValue{}

func (v JSONSubsetValue) Type(_ context.Context) attr.Type {
    return JSONSubsetType{}
}

func (v JSONSubsetValue) Equal(o attr.Value) bool {
    other, ok := o.(JSONSubsetValue)
    if !ok {
        return false
    }
    return v.StringValue.Equal(other.StringValue)
}

// StringSemanticEquals reports v (the planned/prior value) and newV (the
// post-apply or refreshed value) as equal when the planned JSON is a subset
// of the actual JSON. Falls back to byte equality for null/unknown values
// or when either side is not valid JSON (e.g. plain strings, base64).
func (v JSONSubsetValue) StringSemanticEquals(_ context.Context, newV basetypes.StringValuable) (bool, diag.Diagnostics) {
    var diags diag.Diagnostics
    other, ok := newV.(JSONSubsetValue)
    if !ok {
        return false, diags
    }
    if v.IsNull() || v.IsUnknown() || other.IsNull() || other.IsUnknown() {
        return v.StringValue.Equal(other.StringValue), diags
    }
    var planned, actual interface{}
    if err := json.Unmarshal([]byte(v.ValueString()), &planned); err != nil {
        return v.StringValue.Equal(other.StringValue), diags
    }
    if err := json.Unmarshal([]byte(other.ValueString()), &actual); err != nil {
        return v.StringValue.Equal(other.StringValue), diags
    }
    return jsonIsSubset(planned, actual), diags
}

// NewJSONSubsetNull returns a typed null value.
func NewJSONSubsetNull() JSONSubsetValue {
    return JSONSubsetValue{StringValue: basetypes.NewStringNull()}
}

// NewJSONSubsetUnknown returns a typed unknown value.
func NewJSONSubsetUnknown() JSONSubsetValue {
    return JSONSubsetValue{StringValue: basetypes.NewStringUnknown()}
}

// NewJSONSubsetValue wraps a concrete string.
func NewJSONSubsetValue(s string) JSONSubsetValue {
    return JSONSubsetValue{StringValue: basetypes.NewStringValue(s)}
}

// jsonIsSubset returns true when every leaf in plan exists with the same
// value in actual. Maps: every key in plan must exist in actual (with a
// subset value); actual may have extra keys. Slices: same length, positional
// subset comparison. Scalars: reflect.DeepEqual.
func jsonIsSubset(plan, actual interface{}) bool {
    switch p := plan.(type) {
    case map[string]interface{}:
        a, ok := actual.(map[string]interface{})
        if !ok {
            return false
        }
        for key, pv := range p {
            av, exists := a[key]
            if !exists {
                return false
            }
            if !jsonIsSubset(pv, av) {
                return false
            }
        }
        return true
    case []interface{}:
        a, ok := actual.([]interface{})
        if !ok {
            return false
        }
        if len(p) != len(a) {
            return false
        }
        for i := range p {
            if !jsonIsSubset(p[i], a[i]) {
                return false
            }
        }
        return true
    default:
        return reflect.DeepEqual(plan, actual)
    }
}
