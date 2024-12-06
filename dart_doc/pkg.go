package dart_doc

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/radical-ui/flywheel/helpers"
	"github.com/charmbracelet/huh/spinner"
)

var nodeKindMethod = 10
var nodeKindConstructor = 2
var nodeKindClass = 3
var nodeKindField = 16

type enclosedBy struct {
	Name string `json:"name"`
	Kind int    `json:"kind"`
	Href string `json:"href"`
}

type docNode struct {
	Name          string     `json:"name"`
	QualifiedName string     `json:"qualifiedName"`
	Href          string     `json:"href"`
	Kind          int        `json:"kind"`
	Desc          string     `json:"desc"`
	EnclosedBy    enclosedBy `json:"enclosedBy"`
}

type DartDoc struct {
	apiDir string
	nodes  []docNode
}

func NewDartDoc(dartPath string) (*DartDoc, error) {
	objectsPath := path.Join(dartPath, "objects")

	libLastModifiedAt, err := helpers.GetLatestModifiedTime(path.Join(objectsPath, "lib"))
	if err != nil {
		return nil, err
	}

	docLastModifiedAt, _ := helpers.GetLatestModifiedTime(path.Join(objectsPath, "doc"))

	if libLastModifiedAt.After(docLastModifiedAt) {
		if err := generateDartDoc(objectsPath); err != nil {
			return nil, err
		}
	}

	return loadDartDoc(objectsPath)
}

func generateDartDoc(libPath string) error {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.Command("dart", "doc")
	cmd.Dir = libPath
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := errors.New("failed to run spinner")

	spinner.New().
		Title("Parsing dart objects ... this can take a while").
		Action(func() {
			err = cmd.Run()
		}).
		Run()

	if err != nil {
		stderrText := strings.TrimSpace(stderr.String())
		if len(stderrText) > 0 {
			return errors.Join(err, errors.New(stderrText))
		}

		return err
	}

	return nil
}

func loadDartDoc(libPath string) (*DartDoc, error) {
	var nodes []docNode

	apiDir := path.Join(libPath, "doc", "api")
	bytes, err := os.ReadFile(path.Join(apiDir, "index.json"))
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &nodes); err != nil {
		return nil, err
	}

	return &DartDoc{apiDir, nodes}, nil
}
