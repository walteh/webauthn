package indexable

import (
	"reflect"
	"time"
)

type LastModifier interface {
	At() float64
	IsPermanent() bool
}

type LastModifiable interface {
	LastModifiedField() reflect.StructField
}

var _ LastModifier = (*SimpleLastModifier)(nil)

type SimpleLastModifier struct {
	Value     float64
	Permanent bool
}

func (me *SimpleLastModifier) IsPermanent() bool {
	return me.Permanent
}

func (me *SimpleLastModifier) At() float64 {
	return me.Value
}

func NewOnlyOnceModifier() *SimpleLastModifier {
	return &SimpleLastModifier{1, true}
}

func NewCustomLastModifier(val float64, permanent bool) *SimpleLastModifier {
	return &SimpleLastModifier{val, permanent}
}

type LastModiferPermanantOverride struct {
	parent    LastModifier
	permanent bool
}

func NewLastModifierPermanantOverride(lm LastModifier, permanent bool) LastModifier {
	return &LastModiferPermanantOverride{lm, permanent}
}

func (me *LastModiferPermanantOverride) At() float64 {
	return me.parent.At()
}

func (me *LastModiferPermanantOverride) IsPermanent() bool {
	return me.permanent
}

type LastModiferWithTTL interface {
	LastModifier
	Delete() bool
	TTL() time.Duration
}
