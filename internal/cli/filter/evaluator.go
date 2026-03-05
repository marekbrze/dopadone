package filter

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Evaluator struct {
	expr Expr
}

func NewEvaluator(expr Expr) *Evaluator {
	return &Evaluator{expr: expr}
}

func (e *Evaluator) Evaluate(data map[string]interface{}) (bool, error) {
	if e.expr == nil {
		return true, nil
	}
	return e.evaluateExpr(e.expr, data)
}

func (e *Evaluator) evaluateExpr(expr Expr, data map[string]interface{}) (bool, error) {
	switch v := expr.(type) {
	case *LogicalExpr:
		return e.evaluateLogical(v, data)
	case *ComparisonExpr:
		return e.evaluateComparison(v, data)
	default:
		return false, fmt.Errorf("unknown expression type: %T", expr)
	}
}

func (e *Evaluator) evaluateLogical(expr *LogicalExpr, data map[string]interface{}) (bool, error) {
	left, err := e.evaluateExpr(expr.Left, data)
	if err != nil {
		return false, err
	}

	right, err := e.evaluateExpr(expr.Right, data)
	if err != nil {
		return false, err
	}

	switch expr.Operator {
	case TokenAND:
		return left && right, nil
	case TokenOR:
		return left || right, nil
	default:
		return false, fmt.Errorf("unknown logical operator: %v", expr.Operator)
	}
}

func (e *Evaluator) evaluateComparison(expr *ComparisonExpr, data map[string]interface{}) (bool, error) {
	fieldValue, exists := data[expr.Field]
	if !exists {
		return false, nil
	}

	return e.compareValues(fieldValue, expr.Operator, expr.Value)
}

func (e *Evaluator) compareValues(fieldValue interface{}, op TokenType, targetValue interface{}) (bool, error) {
	fieldStr := toString(fieldValue)
	targetStr := toString(targetValue)

	switch op {
	case TokenEQ:
		return strings.EqualFold(fieldStr, targetStr), nil
	case TokenNE:
		return !strings.EqualFold(fieldStr, targetStr), nil
	}

	fieldNum, fieldIsNum := toFloat64(fieldValue)
	targetNum, targetIsNum := toFloat64(targetValue)

	if fieldIsNum && targetIsNum {
		switch op {
		case TokenGT:
			return fieldNum > targetNum, nil
		case TokenGE:
			return fieldNum >= targetNum, nil
		case TokenLT:
			return fieldNum < targetNum, nil
		case TokenLE:
			return fieldNum <= targetNum, nil
		}
	}

	switch op {
	case TokenGT:
		return fieldStr > targetStr, nil
	case TokenGE:
		return fieldStr >= targetStr, nil
	case TokenLT:
		return fieldStr < targetStr, nil
	case TokenLE:
		return fieldStr <= targetStr, nil
	}

	return false, fmt.Errorf("unknown operator: %v", op)
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%v", val)
	case bool:
		return fmt.Sprintf("%t", val)
	case time.Time:
		return val.Format("2006-01-02")
	default:
		return fmt.Sprintf("%v", val)
	}
}

func toFloat64(v interface{}) (float64, bool) {
	if v == nil {
		return 0, false
	}

	switch val := v.(type) {
	case int:
		return float64(val), true
	case int64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Int || rv.Kind() == reflect.Int8 || rv.Kind() == reflect.Int16 || rv.Kind() == reflect.Int32 || rv.Kind() == reflect.Int64 {
			return float64(rv.Int()), true
		}
		if rv.Kind() == reflect.Uint || rv.Kind() == reflect.Uint8 || rv.Kind() == reflect.Uint16 || rv.Kind() == reflect.Uint32 || rv.Kind() == reflect.Uint64 {
			return float64(rv.Uint()), true
		}
		if rv.Kind() == reflect.Float32 || rv.Kind() == reflect.Float64 {
			return rv.Float(), true
		}
		return 0, false
	}
}

func EvaluateFilter(filter string, data []map[string]interface{}) ([]map[string]interface{}, error) {
	if filter == "" {
		return data, nil
	}

	expr, err := Parse(filter)
	if err != nil {
		return nil, fmt.Errorf("invalid filter syntax: %w", err)
	}

	eval := NewEvaluator(expr)
	var result []map[string]interface{}

	for _, item := range data {
		matches, err := eval.Evaluate(item)
		if err != nil {
			return nil, err
		}
		if matches {
			result = append(result, item)
		}
	}

	return result, nil
}
