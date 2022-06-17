package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

// type PriorVersion struct {
// 	XMLName     xml.Name `xml:"priorVersions"`
// 	Version     string   `xml:"version"`
// 	Description string   `xml:",innerxml"`
// }

type PriorVersion struct {
	XMLName xml.Name `xml:"priorVersions"`

	Version []string `xml:"version"`
}

type Namespace struct {
	XMLName       xml.Name     `xml:"namespace"`
	Name          string       `xml:"name"`
	Version       string       `xml:"version"`
	PriorVersions PriorVersion `xml:"priorVersions"`
	Description   string       `xml:",innerxml"`
}
type Namespaces struct {
	XMLName     xml.Name  `xml:"namespaces"`
	Version     string    `xml:"version,attr"`
	Description string    `xml:",innerxml"`
	Namespace   Namespace `xml:"namespace"`
}

func main() {
	file, err := os.Open("./name.xml") // For read access.
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	v := Namespaces{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	// fmt.Println(v)

	fmt.Println(v.Namespace.Name)
	fmt.Println(v.Namespace.Version)
	fmt.Println((v.Namespace.PriorVersions))
	fmt.Println(len(v.Namespace.PriorVersions.Version))

}
