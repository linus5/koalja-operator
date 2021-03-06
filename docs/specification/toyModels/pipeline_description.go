//
// Copyright © 2019 Aljabr, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// ********************************************************************************
//
// koalja pipeline language - POC parser, compiler
//
// ********************************************************************************

package main

import (
        "strings"
        "unicode"
        "io/ioutil"
        "bufio"
        "context"
        "fmt"
        "os"
	H "history"
)

type List []string
type BreadBoard map[string]List

type Container struct {
	alias string
	image_version string
	image string
	command string
	args []string
}

var CPIPELINE H.Concept

// ****************************************************************************

func main() {

	// 1. Initialize

	details := map[string]string{
		"Pod":        "Koalja_empty_pod",
		"Deployment": "Koalja pipeline description",
		"Version":    "0.1",
	}

	ctx := context.Background()
	ctx = H.SetLocationInfo(ctx,details)

	// 2. test koalja, reads pipeline/container_description

	ParsePipeline(ctx,details)

	// 3. run pipeline
}

//**************************************************************
//* Koalja PIPELINE SPEC package
//**************************************************************

func ParsePipeline(ctx context.Context, v map[string]string){

	m := H.SignPost(&ctx,"Parser").
		Intent("Parse a pipeline description file and generate YAML").
		PartOf(H.N("Koalja"))

	buffer, err := ioutil.ReadFile("./pipeline_description")

	if err != nil {
		m.FailedBecause("File read failed").AddError(err)
	}

	namespace, name, yaml := GetPipelineDefinition(ctx,buffer)

	// Print the k8s YAML

	fmt.Println(I(0),"apiVersion: koalja.aljabr.io/"+v["Version"])
	fmt.Println(I(0),"kind: Pipeline")
	fmt.Println(I(0),"metadata:")
	fmt.Println(I(1),"name: "+name)
	fmt.Println(I(1),"namespace: "+namespace)
	fmt.Println(I(0),"spec:")
	fmt.Println(I(1),"tasks:")

	for i := 0; i < len(yaml); i++ {
		fmt.Println(yaml[i])
	}
}

//**************************************************************

func DocumentPipelineByName(namespace,name string) {

	CPIPELINE = H.CreateConcept("Koalja pipeline "+name)
	H.ConceptLink(H.CreateConcept("Koalja pipeline"),H.GENERALIZES,CPIPELINE)
	cnamespace := H.CreateConcept("namespace "+name)
	H.ConceptLink(cnamespace,H.CONTAINS,CPIPELINE)
	kn := H.CreateConcept("kubernetes namespace")
	H.ConceptLink(kn,H.GENERALIZES,cnamespace)
}

//**************************************************************

func GetPipelineDefinition(ctx context.Context, input []byte) (string,string,[]string) {

	var yaml []string = make([]string, 0)
	var name, namespace string

	bb := make(BreadBoard)
	i := 0

	// Begin parsing the breadboard spec

NewTask:
	for i < len(input) {

		//var in,out,task string
		var parin int
		var operator []byte
		var paren[2] []byte
		var nsn []byte

		paren[0] = make([]byte, 0)
		paren[1] = make([]byte, 0)
		operator = make([]byte, 0)
		nsn = make([]byte, 0)

		parin = 0

		// j: get a single statement, k break line into components

		for j := i; j < len(input); j++ {

			switch input[j] {

			case byte('['): 

				for k := j+1; k < len(input); k++ {
					
					if name != "" {
						fmt.Println("Name of pipeline redefined")
					}

					if input[k] == byte(']') {
						j = k+1
						if strings.Contains(string(nsn),":") {
							a := strings.Split(string(nsn),":")
							namespace = a[0]
							name = a[1]
						} else {
							name =string(nsn)
							namespace = "default"
						}
						DocumentPipelineByName(namespace,name)
						break
					} else {
						nsn = append(nsn,input[k])
					}
				}

			case byte('('): 
				for k := j+1; k < len(input); k++ {

					if input[k] == byte(')') {
						k++
						j = k
						break
					}

					if input[k] == byte('\n') {
						continue
					}

					if parin < 2 {
						paren[parin] = append(paren[parin],input[k])
					}

					if k >= len(input){
						break NewTask
					}
				}
				parin++

			case byte('\n'): input[j] = byte(' ')

			case byte('#'): 
				for k := j+1; k < len(input); k++ {

					if input[k] == byte('\n') {
						j = k
						input[k] = byte(' ')
						break
					}

					if k >= len(input) {
						break NewTask
					}
				}
			}
			
			operator = append(operator,input[j])

			i = j+1

			if (parin == 2) {
				y := ReadOneTask(ctx,bb,paren[0],operator,paren[1])
				for i := 0; i < len(y); i++ {
					yaml = append(yaml,y[i])
				}
				break
			}
		}
	}

	// Suffix info

	yaml = append(yaml,I(2)+"links:")

	wires := AddWiring(ctx,bb)

	for i := 0; i < len(wires); i++ {
		yaml = append(yaml,wires[i])
	}
	
	// Suffix info

	yaml = append(yaml,I(2)+"types:")
	yaml = append(yaml,I(3)+"- name: singleTextFile")
	yaml = append(yaml,I(4)+  "protocol: File")
	yaml = append(yaml,I(4)+  "format: Text")

	return namespace, name, yaml
}

//**************************************************************

func LookupContainerDef(ctx context.Context, search string, in []byte, out []byte, ctask H.Concept) []string {

	m := H.SignPost(&ctx,"Container lookup").
		Intent("Import container description from database")

	file, err := os.Open("./container_description")
	
	if err != nil {
		m.Note("file read failed").AddError(err)
		fmt.Println("ERROR......")
	}
	
	defer file.Close()

	var c Container	
	var found bool = false

	scanner := bufio.NewScanner(file)
	yaml := make([]string,0)
	c.args = make([]string,0)

	for scanner.Scan() {
		var alias, key, value string
		line := string(scanner.Text())
		
		if strings.HasPrefix(line,"(") {

			fmt.Sscanf(line,"( %s", &alias)

			if alias == search || alias == "defaults" {

				if alias == search {
					found = true
					c.args = nil
				}

				for scanner.Scan() {
					line = string(scanner.Text())

					if strings.HasPrefix(line,"#") {
						continue
					}

					if strings.HasPrefix(line,")") {
						c.alias = alias
						goto Done
					}
					
					l := strings.Split(line,":")
					key = strings.TrimSpace(l[0])
					value = strings.TrimSpace(l[1])

					if key == "command" {
						c.command = value
					}

					if key == "image" {
						c.image = value
					}

					if key == "image_version" {
						c.image_version = value
					}

					if key == "arg" {
						c.args = append(c.args,value)
					}
					
				} 
			}
		}
	Done:
	}

	if err := scanner.Err(); err != nil {
		m.Note("reading line of file failed").AddError(err)
	}

	if !found {
		m.FailedBecause("No container definition found").AddError(err)
		return append(yaml,"<No Such Container>!! "+search)
	}

	// Get the wire endings to sub for IN and OUT, generate YAML

	var inarray, implicit_in = StripWires(in)
	var outarray, implicit_out = StripWires(out)	

	DocumentImplicit(ctask,implicit_in,implicit_out)

	var s string

	s = fmt.Sprintf("image: %s:%s", c.image,c.image_version)
	yaml = append(yaml,s)
	s = fmt.Sprintf("command:")
	yaml = append(yaml,s)

	// Expand IN and OUT wire ends, might be in command string or in args

	s = fmt.Sprintf("- %s", c.command)
	yaml = append(yaml,s)

	// Document the container and command
	nrimage := H.NR(c.image+":"+c.image_version,"container image")
	nrbinary:= H.NR(c.command,"container command for "+ctask.Name)

	H.ConceptLink(ctask,H.USES,nrbinary.Ccpt)
	H.ConceptLink(nrbinary.Ccpt,H.USES,nrimage.Ccpt)

	cpa := H.CreateConcept("process command argument")
	//

	for i := 0; i < len(c.args); i++ {

		s = c.args[i]

		cargstr := fmt.Sprintf("command arg[%d] %s",i,c.args[i])
		carg := H.CreateConcept(cargstr)

		H.ConceptLink(nrbinary.Ccpt,H.EXPRESSES,carg)
		H.ConceptLink(cpa,H.GENERALIZES,carg)

		// substitute IN/OUT if on separate line, then split the args one per line

		if s == "<IN>" {
			for j := 0; j < len(inarray); j++ {
				s = fmt.Sprintf("%s", "{{.inputs."+StripPolicy(inarray[j])+".path}}")
				c.args[i] = s
				s = fmt.Sprintf("- %s", s)
				yaml = append(yaml,s)

			}
		} else if strings.Contains(s,"<IN>") {
			var expand string
			for j := 0; j < len(inarray); j++ {
				expand = expand + "{{.inputs."+StripPolicy(inarray[j])+".path}} "
			}
			s = strings.Replace(s,"<IN>",expand,1)
			c.args[i] = s
			s = fmt.Sprintf("- %s", s)
			yaml = append(yaml,s)

		}
		
		if s == "<OUT>" {
			for j := 0; j < len(outarray); j++ {
				s = fmt.Sprintf("%s", "{{.outputs."+StripPolicy(outarray[j])+".path}}")
				c.args[i] = s
				s = fmt.Sprintf("- %s", s)
				yaml = append(yaml,s)
			}
		} else if strings.Contains(s,"<OUT>") {
			var expand string
			for j := 0; j < len(outarray); j++ {
				expand = expand + "{{.outputs."+StripPolicy(outarray[j])+".path}} "
				s = fmt.Sprintf("- %s", s)
				yaml = append(yaml,s)
			}
			s = strings.Replace(s,"<OUT>",expand,1)
			c.args[i] = s
		}
	}

	return yaml
}

//**************************************************************

func StripPolicy(s string) string {

	array := strings.Split(s,"[")
	return array[0]
}

//**************************************************************

func DocumentImplicit(ctask H.Concept, in,out []string) {

	for i := range in {
		c := H.CreateConcept("implicit dataflow: "+in[i])
		H.ConceptLink(ctask,H.DEPENDS,c)
	}

	for o := range out {
		c := H.CreateConcept("implicit dataflow: "+out[o])
		H.ConceptLink(c,H.DEPENDS,ctask)
	}
}

//**************************************************************

func StripWires(array []byte) ([]string,[]string) {

	var a []string = strings.FieldsFunc(string(array),InOutSplit)
	var explicit []string
	var implicit []string

	var prev string = "missing identifier"

	for i := 0; i < len(a); i++ {

		if a[i] == "implicit" { // implicit DB/service query
			implicit = append(implicit,prev)
		} else {
			explicit = append(explicit,a[i])
		}
		prev = a[i]
	}

	return explicit,implicit
}

//**************************************************************

func ReadOneTask(ctx context.Context, bb BreadBoard, in []byte, op []byte,out []byte) []string {

	m := H.SignPost(&ctx,"ReadOneTask").
		Intent("Extract inputs, outputs, operator, and policy from a declaration")

	// Generate output assembly to yaml list
	var yaml []string = make([]string, 0)

	// Split: (in) op (out) into parts

	var opa []string = strings.FieldsFunc(string(op),InOutSplit)

	if len(opa) > 1 {
		m.FailedBecause("Too many operators in the pipeline junction")
		fmt.Println("Too many operators in junction: ",opa)
		os.Exit(1)
	}
		
	operator := string(opa[0])

	// get inputs ... if there's an "in" input, then add a filedrop/ingress
	// if the name starts with "in", create a filedrop for it

	drops := CheckFileDrops(ctx,bb,in)

	for i := 0; i < len(drops); i++ {
		yaml = append(yaml,drops[i])
	}

	// Generate the output
	// Start with the NAME of the transformer/task
	
	yaml = append(yaml,I(2)+"- name: "+string(operator))

	ctask := DocumentTask(string(operator))	

	// generate inputs

	var inputs []string = HandleIOPolicy(ctx,bb,operator,in,ctask)

	if (len(inputs)) > 0 {
		yaml = append(yaml,I(3)+"inputs:")
	}

	for i := 0; i < len(inputs); i++ {
		yaml = append(yaml,I(3)+inputs[i])
	}

	// generate outputs

	var outputs []string = HandleIOPolicy(ctx,bb,operator,out,ctask)

	if (len(outputs)) > 0 {
		yaml = append(yaml,I(3)+"outputs:")
	}

	for i := 0; i < len(outputs); i++ {
		yaml = append(yaml,I(3)+outputs[i])
	}

	// Append executor or filedrop ingress node
	
	yaml = append(yaml,I(3)+"executor:")

	containers := LookupContainerDef(ctx,operator,in,out,ctask)

	for i := 0; i < len(containers); i++ {
		yaml = append(yaml,I(4)+containers[i])
	}

	return yaml
}

//**************************************************************

func DocumentTask(taskname string) H.Concept {

	// the official name/role of the container task
	nr := H.NR(taskname,"Koalja pipeline task")

	// the task is part of this pipeline
	H.ConceptLink(CPIPELINE,H.CONTAINS,nr.Ccpt)
	return nr.Ccpt
}

//**************************************************************

func CheckFileDrops(ctx context.Context, bb BreadBoard, in []byte) []string {

	var array []string = strings.FieldsFunc(string(in),InOutSplit)

	var yaml []string = make([]string, 0)

	for i := 0; i < len(array); i++ {

		s := strings.Split(array[i],"[")
		name := s[0]

		if strings.HasPrefix(name,"in") {
			yaml = append(yaml,I(2)+"- name: Drop"+name)
			yaml = append(yaml,I(3)+"outputs: ")
			yaml = append(yaml,I(3)+"- name: "+name)
			yaml = append(yaml,I(4)+"typeRef: singleTextFile")
			yaml = append(yaml,I(4)+"ready: Auto")
			yaml = append(yaml,I(3)+"type: FileDrop")

			if bb[name] == nil {
				bb[name] = make([]string,0)

				ctask := H.CreateConcept("Koalja pipeline task Drop"+name)
				H.ConceptLink(ctask,H.EXPRESSES,H.CreateConcept("Koalja pipeline ingress"))
				H.ConceptLink(CPIPELINE,H.CONTAINS,ctask)
			}

			bb[name] = append(bb["Drop"+name],"Drop"+name)
		}
	}

	return yaml
}

//**************************************************************

func HandleIOPolicy(ctx context.Context, bb BreadBoard, operator string, in []byte, ctask H.Concept) []string {

	var array []string = strings.FieldsFunc(string(in),InOutSplit)

	m := H.SignPost(&ctx,"HandleIOPolicy")

	var yaml []string = make([]string, 0)
	var window, min, max int
	var name,typename string

	// If there's a file input, add a filedrop stanza

	for i := 0; i < len(array); i++ {

		// Default type is file
		typename = "singleTextFile"
		name, min, max, window = GetInputPolicy(ctx,array[i])

		if max == 0 {
			m.Note("Bad separator in min-max range")
			fmt.Println("Bad separator in min-max range: ",array[i],"use - for range")
			os.Exit(1)
		}
		
		if max < min {
			m.Note("Policy error - max < min")
			fmt.Println("Bad separator - max < min: ",array[i])
			os.Exit(1)
		}
		
		if bb[name] == nil {
			bb[name] = make([]string,0)
		}
		
		bb[name] = append(bb[name],operator)
		
		yaml = append(yaml,"- name:"+" "+name)
		yaml = append(yaml,I(1)+"typeRef:"+" "+typename)

		if window > 0 {
			yaml = append(yaml,I(1)+fmt.Sprintf("slide: %d",window))

			cwin := H.CreateConcept(fmt.Sprintf("sliding window size %d",window))
			H.ConceptLink(ctask,H.EXPRESSES,cwin)
		}
		
		if min > 1 && max == min {
			yaml = append(yaml,I(1)+fmt.Sprintf("requiredSequenceLength: %d",min))
			c := H.CreateConcept(fmt.Sprintf("exact required batch size %d",min))
			H.ConceptLink(ctask,H.EXPRESSES,c)

		} else if min > 1 {
			yaml = append(yaml,I(1)+fmt.Sprintf("minSequenceLength: %d",min))
			c := H.CreateConcept(fmt.Sprintf("minimum required batch size %d",min))
			H.ConceptLink(ctask,H.EXPRESSES,c)
		}

		if max > 1 && max > min {
			yaml = append(yaml,I(1)+fmt.Sprintf("maxSequenceLength: %d",max))
			c := H.CreateConcept(fmt.Sprintf("maximum required batch size %d",max))
			H.ConceptLink(ctask,H.EXPRESSES,c)

		}
		
		if strings.HasPrefix(name,"in") {
			yaml = append(yaml,I(1)+"ready: Auto")
		} else {
			yaml = append(yaml,I(1)+"ready: Succeeded")
		}
	}

return yaml
}

//**************************************************************

func InOutSplit (r rune) bool {
	switch r {
	case ',': return true
	default: return unicode.IsSpace(r)
	}
}

//**************************************************************

func GetInputPolicy(ctx context.Context,in string) (string,int,int,int) {

	m := H.SignPost(&ctx,"GetInputPolicy from file")

	var min,max int
	var ch rune
	var slide int = 0

	s := strings.Split(in,"[")

	name := s[0]

	if len(s) == 1 {

		min = 1
		max = 1

	} else {

		min = 1
		max = 0
		
		if strings.Contains(s[1],"-") && strings.Contains(s[1],"/") {
			m.Note("Can't combine min-max range with sliding window").
				Attributes(H.NR(in,"policy range"))
			os.Exit(1)
		}

		if strings.Contains(s[1],"-") {
			fmt.Sscanf(s[1],"%d%c%d]", &min,&ch,&max)
			
			// Syntax error - used a comma in the range specifier
			if ch == 0 {
				m.Note("Bad separator in min-max range").
					Attributes(H.NR(in,"policy range"))
				os.Exit(1)
			}

		} else if strings.Contains(s[1],"/") {

			fmt.Sscanf(s[1],"%d%c%d]", &min,&ch,&slide)
			max = min
		} else {
			fmt.Sscanf(s[1],"%d", &min)
			max = min
		}
	}

	return name, min, max, slide
}

//**************************************************************

func AddWiring(ctx context.Context, bb BreadBoard) []string {

	var yaml []string = make([]string, 0)

	m := H.SignPost(&ctx,"AddWiring").
		Intent("Trace wires to ensure they are all connected")

	for k,v := range bb {
		if strings.HasPrefix(k,"out") {
		} else {

			if len(v) < 2 {
				fmt.Println("Bad wire ("+k+") emanating from ",v)
				m.FailedBecause("Bad wire ("+k+") emanating from "+v[0])
				os.Exit(1)	
			}

			var name,source,dest string
			
			for j := 1; j < len(v); j++ {
				//fmt.Println("WIRES: ",k,v)
				name = v[0]+"2"+v[j]

				cfrom := H.CreateConcept("Koalja pipeline task "+v[0])
				cto := H.CreateConcept("Koalja pipeline task "+v[j])
				H.ConceptLink(cto,H.FOLLOWS,cfrom)

				source = fmt.Sprintf("%s/%s",v[0],k)
				dest = fmt.Sprintf("%s/%s",v[j],k)
				yaml = append(yaml,I(3)+"- name:"+" "+name)
				yaml = append(yaml,I(4)+"sourceRef:"+" "+source)
				yaml = append(yaml,I(4)+"destinationRef:"+" "+dest)
			}
		}
	}

	return yaml
}

//**************************************************************

func I(level int) string {
	var indent string = strings.Repeat("  ",level)
	var s string
	s = fmt.Sprintf("%.3d:%s",level,indent)
	s = indent
	return s
}