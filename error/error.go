// Copyright (c) 2018 Anton Semjonov
// Licensed under the MIT License

package error

// ConstError is a constant-expression error
// (https://dave.cheney.net/2016/04/07/constant-errors)
type ConstError string

// Error makes this const compatible with the Error interface
func (e ConstError) Error() string { return string(e) }
