//------------------------------------------------------------------------------
// the code here is about building the config object
//------------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/xgfone/go-tools/file"
)

// getting the config from an existing JSON (.conf) file, or building it from the definition
func (commonDef *fileMap) getOrInitConfig(folderPath string, folderInfo os.FileInfo, configFileName string) (*j2tConfig, error) {

	// initialising the config object
	config := &j2tConfig{
		folderPath: folderPath,
		folderInfo: folderInfo,
	}

	// initialising the config, if not done yet (well, we know that its not done yet for now)
	commonDef.initConfigMap() // to change, with the lecture of the config file, if it exists

	// loading the config file, if present
	configFile := folderPath + "/" + configFileName
	if file.IsExist(configFile) {
		filebytes, err := ioutil.ReadFile(configFile)
		if err != nil {
			return nil, err
		}
		if errMarshall := json.Unmarshal(filebytes, config); errMarshall != nil {
			return nil, errMarshall
		}
	}

	return config, nil
}

// initialising the config file
func (commonDef *fileMap) initConfigMap() {

	for index, property := range commonDef.orderedProperties {

		// about to "compute" foregound and background colors
		var bg string

		// dealing with the top level, i.e. the first header line
		if commonDef.parent == nil {
			bg = getAdjustedColor(colors[index%len(colors)], 0, false)

		} else {

			// dealing with a sub-section level
			level := 10
			if index%2 == 1 {
				level = 30
			}
			bg = getAdjustedColor(commonDef.parent.chainedProperties[commonDef.name].conf.background, level*commonDef.getDepth(), false)
		}

		// new item
		newItem := &configItem{
			background: bg,
		}

		// linking the prop to its config
		commonDef.chainedProperties[property].conf = newItem

		// going deeper
		if subMap := commonDef.subMaps[property]; subMap != nil {
			subMap.initConfigMap()
		}
	}
}
