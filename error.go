package main

// Error is a constant-expression error
// (https://dave.cheney.net/2016/04/07/constant-errors)
type Error string

func (e Error) Error() string { return string(e) }
