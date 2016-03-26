// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file defines the check for unused results of calls to certain
// pure functions.

package main

import (
	"flag"
	"go/ast"
	"go/token"
	"strings"

	"fmt"

	"golang.org/x/tools/go/types"
)

var unusedFuncsFlag = flag.String("unusedfuncs",
	"errors.New,fmt.Errorf,fmt.Sprintf,fmt.Sprint,sort.Reverse",
	"comma-separated list of functions whose results must be used")

var unusedStringMethodsFlag = flag.String("unusedstringmethods",
	"Error,String",
	"comma-separated list of names of methods of type func() string whose results must be used")

func init() {
	register("unusedresult",
		"check for unused result of calls to functions in -unusedfuncs list and methods in -unusedstringmethods list",
		checkUnusedResult,
		exprStmt)
}

// func() string
var sigNoArgsStringResult = types.NewSignature(nil, nil,
	types.NewTuple(types.NewVar(token.NoPos, nil, "", types.Typ[types.String])),
	false)

var unusedFuncs = make(map[string]bool)
var unusedStringMethods = make(map[string]bool)

func initUnusedFlags() {
	commaSplit := func(s string, m map[string]bool) {
		if s != "" {
			for _, name := range strings.Split(s, ",") {
				if len(name) == 0 {
					flag.Usage()
				}
				m[name] = true
			}
		}
	}
	fmt.Printf("unusedFuncsFlag: %s\n", *unusedFuncsFlag)
	fmt.Printf("unusedStringMethodsFlag: %s\n", *unusedStringMethodsFlag)
	commaSplit(*unusedFuncsFlag, unusedFuncs)
	commaSplit(*unusedStringMethodsFlag, unusedStringMethods)
}

func checkUnusedResult(f *File, n ast.Node) {
	fmt.Printf("\n[start]\n")
	call, ok := unparen(n.(*ast.ExprStmt).X).(*ast.CallExpr)
	if !ok {
		fmt.Printf("NO REPORT: not call\n")
		return // not a call statement
	}
	fun := unparen(call.Fun)

	fmt.Printf("[checking for conversion]\n")
	if f.pkg.types[fun].IsType() {
		fmt.Printf("NO REPORT: conversion\n")
		return // a conversion, not a call
	}

	fmt.Printf("[checking for (not method + unqualified)]\n")
	selector, ok := fun.(*ast.SelectorExpr)
	if !ok {
		fmt.Printf("NO REPORT: not method + unqualified\n")
		return // neither a method call nor a qualified ident
	}

	fmt.Printf("[checking pkg.selectors]\n")
	sel, ok := f.pkg.selectors[selector]
	if ok && sel.Kind() == types.MethodVal {
		fmt.Printf("[pkg.selectors okay and kind is methodval]\n")
		// method (e.g. foo.String())
		obj := sel.Obj().(*types.Func)
		sig := sel.Type().(*types.Signature)
		fmt.Printf("[checking type signature against sigNoArgsStringResult]\n")
		if types.Identical(sig, sigNoArgsStringResult) {
			fmt.Printf("[checking unusedStringMethods]\n")
			if unusedStringMethods[obj.Name()] {
				fmt.Printf("REPORT: method `%s` is in unusedStringMethods\n", obj.Name())
				f.Badf(call.Lparen, "result of (%s).%s call not used",
					sig.Recv().Type(), obj.Name())
			} else {
				fmt.Printf("NO REPORT BUT REPORTABLE: method `%s` is not in unusedStringMethods\n", obj.Name())
			}
		} else {
			fmt.Printf("NO REPORT: type signature is not sigNoArgsStringResult\n")
		}
	} else if !ok {
		fmt.Printf("[pkg.selectors not okay]\n")
		// package-qualified function (e.g. fmt.Errorf)
		fmt.Printf("[checking pkg.uses]\n")
		obj, _ := f.pkg.uses[selector.Sel]
		if obj, ok := obj.(*types.Func); ok {
			qname := obj.Pkg().Path() + "." + obj.Name()
			fmt.Printf("[checking for unusedFuncs]\n")
			if unusedFuncs[qname] {
				fmt.Printf("REPORT: qualified name `%s` is in unusedFuncs\n", qname)
				f.Badf(call.Lparen, "result of %v call not used", qname)
			} else {
				fmt.Printf("NO REPORT BUT REPORTABLE: qualified name `%s` is not in unusedFuncs\n", qname)
			}
		} else {
			fmt.Printf("[pkg.uses was not okay]\n")
		}
	} else {
		fmt.Printf("NO REPORT: pkg.selectors okay but kind is not MethodVal\n")
	}
}
