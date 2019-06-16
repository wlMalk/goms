package parser

import (
	"fmt"
	"strconv"
	strs "strings"

	"github.com/wlMalk/goms/parser/types"
)

func ParseVersion(ver string) (*types.Version, error) {
	v := &types.Version{}
	var err error
	v.Major, v.Minor, v.Patch, err = parseVersion(strs.ToLower(ver))
	if err != nil {
		return nil, err
	}
	return v, nil
}

func parseVersion(ver string) (int, int, int, error) {
	err := fmt.Errorf("cannot parse \"%s\" as a version", ver)
	if ver == "" {
		return 1, 0, 0, nil
	}
	matches := versionPattern.FindAllStringSubmatch(ver, -1)
	if len(matches) != 1 {
		return 0, 0, 0, err
	}
	major := 0
	minor := 0
	patch := 0
	var nerr error
	if matches[0][1] != "" {
		if major, nerr = strconv.Atoi(matches[0][1]); nerr != nil {
			return 0, 0, 0, err
		}
	}
	if matches[0][4] != "" {
		if minor, nerr = strconv.Atoi(matches[0][4]); nerr != nil {
			return 0, 0, 0, err
		}
	}
	if matches[0][7] != "" {
		if patch, nerr = strconv.Atoi(matches[0][7]); nerr != nil {
			return 0, 0, 0, err
		}
	}
	return major, minor, patch, nil
}
