package paxos

import (
	"runtime/debug"
	"strings"
)

// devVersion is reported when the SDK's own tagged module version can't be
// resolved from build info — e.g. running this module's own test suite, or
// a consumer building against it via a local `replace` directive.
const devVersion = "dev"

// sdkModulePath is this SDK's module path, used to find its own entry in
// runtime/debug.BuildInfo.
const sdkModulePath = "github.com/paxosglobal/paxos-sdk-go"

// Version returns the SDK version string sent in the Paxos-SDK-Version telemetry
// header, without the leading "v" of the module version. It is resolved from the
// consumer's build info (the tagged module version Go modules actually resolved),
// not hardcoded, so it can't drift from the version the caller depends on.
func Version() string { return version }

var version = resolveVersion()

func resolveVersion() string {
	info, ok := debug.ReadBuildInfo()
	return versionFromBuildInfo(info, ok)
}

// versionFromBuildInfo extracts the SDK's resolved module version from build
// info. Split out from resolveVersion so it can be tested against a
// fabricated *debug.BuildInfo, without depending on this repo actually being
// consumed as a tagged dependency.
func versionFromBuildInfo(info *debug.BuildInfo, ok bool) string {
	if !ok || info == nil {
		return devVersion
	}

	// Rare: the SDK itself is the main module being built (e.g. a one-off
	// binary built from inside this repo).
	if info.Main.Path == sdkModulePath && isResolvedVersion(info.Main.Version) {
		return strings.TrimPrefix(info.Main.Version, "v")
	}

	// Ordinary case: the SDK is a dependency of the consumer's main module.
	for _, dep := range info.Deps {
		if dep.Path != sdkModulePath {
			continue
		}
		if dep.Replace != nil {
			if isResolvedVersion(dep.Replace.Version) {
				return strings.TrimPrefix(dep.Replace.Version, "v")
			}
			return devVersion
		}
		if isResolvedVersion(dep.Version) {
			return strings.TrimPrefix(dep.Version, "v")
		}
		return devVersion
	}

	return devVersion
}

func isResolvedVersion(v string) bool {
	return v != "" && v != "(devel)"
}
