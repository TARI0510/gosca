package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

type VulnDb struct {
	Module             string              `yaml:"module"`
	Package            string              `yaml:"package"`
	AdditionalPackages []VulnDb            `yaml:"additional_packages"`
	Versions           []map[string]string `yaml:"versions"`
	Description        string              `yaml:"description"`
	//Published          string              `yaml:"published"`
	Cves           []string `yaml:"cves"`
	Cvss3          []string `yaml:"cvss3"`
	Severities     []string `yaml:"severities"`
	Symbols        []string `yaml:"symbols"`
	DerivedSymbols []string `yaml:"derived_symbols"`
	Links          link     `yaml:"links"`
}

type link struct {
	Pr      string   `yaml:"pr"`
	Commit  string   `yaml:"commit"`
	Context []string `yaml:"context"`
}

type VulnDbMap map[string]VulnDb

func (v *VulnDb) String() string {
	return fmt.Sprintf("%s %s %s %s %s %s %s", v.Module, v.Package, v.AdditionalPackages, v.Versions, v.Description, v.Cves, v.Symbols)
}

func (vm VulnDbMap) ReadVulnDbYaml(vulnDBs map[string]string) VulDbIdxMap {
	var vulDbIdxMap VulDbIdxMap = make(VulDbIdxMap)

	for yamlPath, yamlContent := range vulnDBs {
		vuln := VulnDb{}
		err := yaml.Unmarshal([]byte(yamlContent), &vuln)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[skip] Error while parse "+yamlPath+": "+err.Error()+"\n")
			continue
		}
		id := strings.TrimRight(strings.TrimRight(filepath.Base(yamlPath), ".yaml"), ".yml")
		vm[id] = vuln

		if _, ok := vulDbIdxMap[vuln.Module]; ok {
			vulDbIdxMap[vuln.Module] = append(vulDbIdxMap[vuln.Module], VulnDbIdx{Id: id, Versions: vuln.Versions})
		} else {
			vulDbIdxMap[vuln.Module] = []VulnDbIdx{VulnDbIdx{Id: id, Versions: vuln.Versions}}
		}
	}

	return vulDbIdxMap
}
