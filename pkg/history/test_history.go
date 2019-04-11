
package main

import (
	"strings"
	"fmt"
	"context"
	"io/ioutil"
	H "github.com/AljabrIO/koalja-operator/pkg/history"
)

// ****************************************************************************
// SPLIT ! TOP
// 2. Koalja program starts BELOW ...
// ****************************************************************************

func main() {

	// 1. test cellibrium

	ctx := context.Background()
	ctx = H.LocationInfo(ctx, map[string]string{
		"Pod":     "A_pod_named_foo",
		"Process": "myApp_name2",  // insert instance data from env?
		"Version": "1.2.3",
	})

	MainLoop(ctx)

	// 2. test koalja, reads pipeline/container_description

}

//**************************************************************
// 1. Cellibrium application test - Non package code
//**************************************************************

func MainLoop(ctx context.Context){

	H.RefMarker(&ctx,"MainLoop start").
		PartOf(H.NR("main","function"))

	// Adaptive loop to update context by sensor/input activity
	// Context changes as a result of new environment detected

	// Start loop
	ctx = H.UpdateSensorContext(ctx)

        // ...other stuff happens
	mk := H.RefMarker(&ctx,"Beginning of test code").
		Role("Start process").
		Attributes(H.NR("cellibrium","go package"),H.N("example code"))
	// ...
	mk.Note(&ctx,"look up a name")

	// ...
	H.RefMarker(&ctx,"code signpost X"). // what you intended
	Intent("open file X").
		Used(H.NR("/etc/passed","file")).
		Used(H.NR("123.456.789.123","dns lookup")).
		FailedToUse(H.N("cc")).
		FailedIntent("xxx").
		PartOf(H.NR("main","function")).
		Contains(H.NR("Test1","test function"))

	// Pass ctx down for logging in lower levels
	Test1(ctx)

	// End loop
	H.RefMarker(&ctx,"The end!")
}

//**************************************************************

func Test1(ctx context.Context){

	m := H.RefMarker(&ctx,"TEST1---------").
		PartOf(H.N("Testing suite"))

	m.Intent("read whole file of data").
		Used(H.NR("file://URI","file"))

	_, err := ioutil.ReadFile("file://URI")

	if err != nil {
	        m.Note(&ctx,"file read failed").AddError(err)
	}
	
}

//**************************************************************

func I(level int) string {
	var indent string = strings.Repeat("  ",level)
	var s string
	s = fmt.Sprintf("%.3d:%s",level,indent)
	s = indent
	return s
}