package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Define config file location and expand ${HOME}
var configPath = "${HOME}/.config/xheads/config.yml"
func init() {
	configPath = os.ExpandEnv(configPath)
}

func isBooleanTrue(o OptionValue) bool {
	val, ok := o.(bool)
	return val && ok
}

func isNonEmpty(o OptionValue) bool {
	val, ok := o.(string)
	return ok && len(val) > 0
}

func optionString(o OptionValue) string {
	val, ok := o.(string)
	if !ok {
		panic(fmt.Errorf("Option value %v is not string", o))
	}
	return val
}

// OptionName is xrandr option key
type OptionName string

func (o *OptionName) String() string {
	if len(*o) == 1 {
		return fmt.Sprintf("-%s", o)
	}

	formattedOption := strings.Replace(string(*o), "_", "-", -1)
	return fmt.Sprintf("--%s", formattedOption)
}

// OptionValue is xrandr option value. Can be boolean or string
type OptionValue interface{}

// OutputName is name of monitor (output in terms of xrandr)
type OutputName string

func (o *OutputName) String() string {
	return fmt.Sprintf("--output %s", string(*o))
}

// OutputConfig is map of OptionName => OptionValue
type OutputConfig map[OptionName]OptionValue

func (o *OutputConfig) String() string {
	result := make([]string, 0)

	for optionName, optionValue := range *o {
		if isBooleanTrue(optionValue) {
			result = append(result, optionName.String())
		} else if isNonEmpty(optionValue) {
			str := optionString(optionValue)
			result = append(result, strings.Join([]string{
				optionName.String(), str}, " "))
		}
	}

	sort.Strings(result)
	return strings.Join(result, " ")
}

// LayoutName is name for desired monitors layout
type LayoutName string

// LayoutConfig is map of OutputName => OutputConfig
type LayoutConfig map[OutputName]OutputConfig

func (lc *LayoutConfig) String() string {
	result := make([]string, 0)
	fmt.Println("test")
	for outputName, outputConfig := range *lc {
		result = append(result,
			strings.Join([]string{outputName.String(), outputConfig.String()}, " "))
	}

	sort.Strings(result)
	return strings.Join(result, " ")
}

// Config is a collection of layouts
type Config struct {
	Layouts map[LayoutName]LayoutConfig
}

func configRead(path string) (cfg *Config, err error) {
	cfg = &Config{}
	buf, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	yaml.Unmarshal(buf, cfg)
	return cfg, err
}

func currentOutputs() []string {
	return scanConnectedOutputs()
}

func makeXrandrScanCommand() string {
	return fmt.Sprintf("xrandr | grep ' connected' | awk '{print $1}'")
}

func scanConnectedOutputs() (outputs []string) {
        cmd := makeXrandrScanCommand()
        out, err := exec.Command("sh", "-c", cmd).Output()
        if err != nil {
                fmt.Println("Cannot read connected displays:", err)
                return
        }
        lines := strings.Trim(string(out), "\n")
        outputs = strings.Split(lines, "\n")
        return
}

func main() {
	// Open and read Config file
	cfg, err := configRead(configPath)
	if err != nil {
		fmt.Printf("Can't read config %s: %v\n", configPath, err)
		os.Exit(1)
	}

	// Fetch currently connected monitors
	outputs := currentOutputs()

	// Print results
	fmt.Println("outputs: ", outputs)
	fmt.Println("configPath: ", configPath)
	fmt.Println("config: ", cfg)
}
