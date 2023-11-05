package helper_test

import (
	"fmt"
	"testing"

	"github.com/HoskeOwl/portscan/internal/helper"
	"github.com/stretchr/testify/assert"
)

func TestSinglePortRangeString(t *testing.T) {
	pr := helper.PortRange{Start: 22, End: 23}
	assert.Equal(t, "22", fmt.Sprintf("%v", pr), "Wrong single port to string")
}

func TestPortRangeString(t *testing.T) {
	pr := helper.PortRange{Start: 22, End: 24}
	assert.Equal(t, fmt.Sprintf("22%v23", helper.PortRangeSeparator), fmt.Sprintf("%v", pr), "Wrong ports to string")
}

func TestSinglePortRangePorts(t *testing.T) {
	pr := helper.PortRange{Start: 22, End: 23}
	exp := []int{22}
	assert.ElementsMatch(t, exp, pr.Ports(), "Wrong single port to array")
}

func TestPortRangePorts(t *testing.T) {
	pr := helper.PortRange{Start: 22, End: 25}
	exp := []int{22, 23, 24}
	assert.ElementsMatch(t, exp, pr.Ports(), "Wrong ports to array")
}

func TestParseSinglePort(t *testing.T) {
	s := "22"
	sExp := 22
	eExp := 23
	r, e := helper.ParseRange(s)
	assert.Nil(t, e)
	assert.Equal(t, sExp, r.Start, "Wrong start port in range")
	assert.Equal(t, eExp, r.End, "Wrong end port in range")
}

func TestParseSinglePortErr(t *testing.T) {
	s := "22w"
	_, e := helper.ParseRange(s)
	assert.NotNil(t, e)
}

func TestParsePortRangeWithOldSeparator(t *testing.T) {
	s := fmt.Sprintf("33%v35", helper.OldPortRangeSeparator)
	sExp := 33
	eExp := 36
	r, e := helper.ParseRange(s)
	assert.Nil(t, e)
	assert.Equal(t, sExp, r.Start, "Wrong start port in range")
	assert.Equal(t, eExp, r.End, "Wrong end port in range")
}

func TestParsePortRangeWithOldSeparatorErr(t *testing.T) {
	s := fmt.Sprintf("33%v35%v38", helper.OldPortRangeSeparator, helper.OldPortRangeSeparator)
	_, e := helper.ParseRange(s)
	assert.NotNil(t, e)
}

func TestParsePortRangeWithSeparator(t *testing.T) {
	s := fmt.Sprintf("38%v45", helper.OldPortRangeSeparator)
	sExp := 38
	eExp := 46
	r, e := helper.ParseRange(s)
	assert.Nil(t, e)
	assert.Equal(t, sExp, r.Start, "Wrong start port in range")
	assert.Equal(t, eExp, r.End, "Wrong end port in range")
}

func TestParsePortRangeWithSeparatorErr(t *testing.T) {
	s := fmt.Sprintf("33%v35%v38", helper.PortRangeSeparator, helper.PortRangeSeparator)
	_, e := helper.ParseRange(s)
	assert.NotNil(t, e)
}

func TestParsePortRangeWithSeparatorErr2(t *testing.T) {
	s := fmt.Sprintf("33%v35w", helper.PortRangeSeparator)
	_, e := helper.ParseRange(s)
	assert.NotNil(t, e)
}

func TestParseSinglePortRanges(t *testing.T) {
	s := "22"
	sExp := 22
	eExp := 23
	lExp := 1
	r, e := helper.ParsePortRanges(s)
	assert.Nil(t, e)
	assert.Len(t, r, lExp)
	assert.Equal(t, sExp, r[0].Start, "Wrong start port in range")
	assert.Equal(t, eExp, r[0].End, "Wrong end port in range")
}

func TestParseSinglesPortRanges(t *testing.T) {
	s := fmt.Sprintf("22%v23", helper.RangesSeparator)
	s1Exp := 22
	e1Exp := 23
	s2Exp := 23
	e2Exp := 24
	lExp := 2
	r, e := helper.ParsePortRanges(s)
	assert.Nil(t, e)
	assert.Len(t, r, lExp)
	assert.Equal(t, s1Exp, r[0].Start, "Wrong start port in range")
	assert.Equal(t, e1Exp, r[0].End, "Wrong end port in range")
	assert.Equal(t, s2Exp, r[1].Start, "Wrong start port in range")
	assert.Equal(t, e2Exp, r[1].End, "Wrong end port in range")
}

func TestParsePortRange(t *testing.T) {
	s := fmt.Sprintf("22%v25", helper.PortRangeSeparator)
	s1Exp := 22
	e1Exp := 26
	lExp := 1
	r, e := helper.ParsePortRanges(s)
	assert.Nil(t, e)
	assert.Len(t, r, lExp)
	assert.Equal(t, s1Exp, r[0].Start, "Wrong start port in range")
	assert.Equal(t, e1Exp, r[0].End, "Wrong end port in range")
}

func TestParsePortRanges(t *testing.T) {
	s := fmt.Sprintf("22%v26%v36%v40", helper.PortRangeSeparator, helper.RangesSeparator, helper.PortRangeSeparator)
	s1Exp := 22
	e1Exp := 27
	s2Exp := 36
	e2Exp := 41
	lExp := 2
	r, e := helper.ParsePortRanges(s)
	assert.Nil(t, e)
	assert.Len(t, r, lExp)
	assert.Equal(t, s1Exp, r[0].Start, "Wrong start port in range")
	assert.Equal(t, e1Exp, r[0].End, "Wrong end port in range")
	assert.Equal(t, s2Exp, r[1].Start, "Wrong start port in range")
	assert.Equal(t, e2Exp, r[1].End, "Wrong end port in range")
}

func TestParsePortRangesWithSingle(t *testing.T) {
	s := fmt.Sprintf("22%v26%v80%v36%v40", helper.PortRangeSeparator, helper.RangesSeparator, helper.RangesSeparator, helper.PortRangeSeparator)
	s1Exp := 22
	e1Exp := 27
	s2Exp := 80
	e2Exp := 81
	s3Exp := 36
	e3Exp := 41
	lExp := 3
	r, e := helper.ParsePortRanges(s)
	assert.Nil(t, e)
	assert.Len(t, r, lExp)
	assert.Equal(t, s1Exp, r[0].Start, "Wrong start port in range")
	assert.Equal(t, e1Exp, r[0].End, "Wrong end port in range")
	assert.Equal(t, s2Exp, r[1].Start, "Wrong start port in range")
	assert.Equal(t, e2Exp, r[1].End, "Wrong end port in range")
	assert.Equal(t, s3Exp, r[2].Start, "Wrong start port in range")
	assert.Equal(t, e3Exp, r[2].End, "Wrong end port in range")
}

func TestParseSinglePortRangesErr(t *testing.T) {
	s := "22w"
	r, e := helper.ParsePortRanges(s)
	assert.Nil(t, r)
	assert.NotNil(t, e)
}

func TestParsePortRangesErr(t *testing.T) {
	s := fmt.Sprintf("22%vq", helper.PortRangeSeparator)
	r, e := helper.ParsePortRanges(s)
	assert.Nil(t, r)
	assert.NotNil(t, e)
}
