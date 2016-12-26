code detector
=====
[![GoDoc](https://godoc.org/github.com/emirpasic/gods?status.svg)](https://godoc.org/github.com/chennqqi/codedetector/detector)

detect script language type by rules.

support user defined rules or interalrules.


## internal rule detectable languages
* jsp
* asp
* aspx
* Python
* php
## Install
```
go get github.com/chennqqi/codedetector
```

## Usage

	Usage of codedetector:
		-d string   
        	set working in dir mode(<-d/-f> is must   
		-f string  
        	set target file path(<-d/-f> is must   

## code

    import "fmt"
 	import "github.com/chennqqi/codedetector"

	function main() {
		det,_ := detector.LoadInteralRules()
		r,_ := det.Detect(targetTextCode)
        fmt.Println(r)
	}

## known issue
 
 if code type GBK/UTF-16 text may not work correctly.  

## License
Apache2