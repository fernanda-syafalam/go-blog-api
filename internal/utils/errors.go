package utils

import (
	"errors"
	"fmt"
)

type ErrNotFound string

func (e ErrNotFound) Error() string {
	if e == "" {
		return "Sumber daya tidak ditemukan"
	}
	return fmt.Sprintf("%s tidak ditemukan", string(e))
}

func IsErrNotFound(err error) bool {
	var e ErrNotFound
	return errors.As(err, &e)
}

type ErrForbidden string

func (e ErrForbidden) Error() string {
	if e == "" {
		return "Anda tidak memiliki izin untuk melakukan tindakan ini"
	}
	return string(e)
}

func IsErrForbidden(err error) bool {
	var e ErrForbidden
	return errors.As(err, &e)
}

type ErrValidation string

func (e ErrValidation) Error() string {
	if e == "" {
		return "Data tidak valid"
	}
	return string(e)
}

func IsErrValidation(err error) bool {
	var e ErrValidation
	return errors.As(err, &e)
}

type ErrInternal string

func (e ErrInternal) Error() string {
	if e == "" {
		return "Terjadi kesalahan internal"
	}
	return string(e)
}

func IsErrInternal(err error) bool {
	var e ErrInternal
	return errors.As(err, &e)
}

