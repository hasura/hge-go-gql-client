package boolexpr

func And[T any](rhs []T) map[string]any {
	return arrRhsOps("_and", rhs)
}

func Or[T any](rhs []T) map[string]any {
	return arrRhsOps("_or", rhs)
}

func Not(rhs any) map[string]any {
	return anyRhsOp("_not", rhs)
}

func In[T any](rhs []T) map[string]any {
	return arrRhsOps("_in", rhs)
}

func NotIn[T any](rhs []T) map[string]any {
	return arrRhsOps("_nin", rhs)
}

func Eq(rhs any) map[string]any {
	return anyRhsOp("_eq", rhs)
}

func NotEq(rhs any) map[string]any {
	return anyRhsOp("_neq", rhs)
}

func Gt(rhs any) map[string]any {
	return anyRhsOp("_gt", rhs)
}

func GtEq(rhs any) map[string]any {
	return anyRhsOp("_gte", rhs)
}

func Lt(rhs any) map[string]any {
	return anyRhsOp("_lt", rhs)
}

func LtEq(rhs any) map[string]any {
	return anyRhsOp("_lte", rhs)
}

func IsNull(isNull bool) map[string]any {
	return anyRhsOp("_is_null", isNull)
}

func arrRhsOps[T any](op string, rhs []T) map[string]any {
	return map[string]any{
		op: rhs,
	}
}

func anyRhsOp(op string, rhs any) map[string]any {
	return map[string]any{
		op: rhs,
	}
}
