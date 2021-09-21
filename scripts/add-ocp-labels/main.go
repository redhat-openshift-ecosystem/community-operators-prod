package main

import (
	goerror "errors"
	"fmt"
	"github.com/blang/semver"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/yaml"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	apimanifests "github.com/operator-framework/api/pkg/manifests"
	apivalidation "github.com/operator-framework/api/pkg/validation"
	"github.com/operator-framework/api/pkg/validation/errors"
	log "github.com/sirupsen/logrus"
)

const ocpLabelAnnotation = "com.redhat.openshift.versions"

// OCP version where the apis v1beta1 is no longer supported
const ocpVerV1beta1Unsupported = "4.9"

type BundleAnnotations struct {
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

type UpdateGraph struct {
	UpdateGraph string `yaml:"annotations,omitempty"`
}

type BundleDir struct {
	Path  string
	Error string
}

type Output struct {
	Errors                       []BundleDir
	RequiresUpdate               []BundleDir
	RequiresUpdatePkgFormat      []BundleDir
	BundleUpdated                []BundleDir
	BundleUnload                 []BundleDir
	BundlesConfiguredAccordingly []BundleDir
	BundlesMigrated              []BundleDir
	UpdatedUsingServerMode       []BundleDir
	PkgManifestUsingServerMode   []BundleDir
	CurrentPath                  string
}

func main() {
	createOutput()
}

func createOutput() {
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	output := Output{CurrentPath: currentPath}
	output.updateAllBundles()
	output.checkIfWeCanLoadTheBundles()

	// Only to add info in the output report
	output.checkBundlesWhichRequiresUpdateAfterChangesMade()
	output.checkBundlesUpdateUsingSemver()
	output.checkPkgUsingSemver()

	output.sort()
	output.write()

}

func (o *Output) checkPkgUsingSemver() {
	// Check upgrade graph in 4.9+ of the bundles updated
	for _, bu := range o.RequiresUpdatePkgFormat {
		pkgDir := filepath.Clean(filepath.Join(bu.Path, ".."))

		// check if was not added already
		found := false
		for _, bh := range o.PkgManifestUsingServerMode {
			if bh.Path == pkgDir {
				found = true
				break
			}
		}
		if found {
			continue
		}

		err2 := filepath.Walk(filepath.Join(o.CurrentPath, pkgDir), func(pathPkg string, infoPkg os.FileInfo, err2 error) error {
			if err2 != nil {
				return err2
			}
			if infoPkg != nil && !infoPkg.IsDir() && strings.Contains(infoPkg.Name(), "ci.yaml") {
				filePath := filepath.Join(o.CurrentPath, pkgDir, infoPkg.Name())
				if _, err := os.Stat(filePath); err == nil && !os.IsNotExist(err) {
					infoPkg, err := ReadFile(filePath)
					if err != nil {
						log.Infof("unable to read infoPkg: %s", err)
						return nil
					}

					var upgradeGraph UpdateGraph
					if err := yaml.Unmarshal(infoPkg, &upgradeGraph); err == nil {
						if upgradeGraph.UpdateGraph == "semver-mode" {
							o.PkgManifestUsingServerMode = append(o.PkgManifestUsingServerMode, BundleDir{Path: pkgDir})
							return nil
						}
					}
				}
				return nil
			}
			return nil
		})
		if err2 != nil {
			log.Fatalf("checking server mode", err2)
		}

	}
}

func (o *Output) checkBundlesUpdateUsingSemver() {
	// Check upgrade graph in 4.9+ of the bundles updated
	for _, bu := range o.BundleUpdated {
		pkgDir := filepath.Clean(filepath.Join(bu.Path, ".."))

		// check if was not added already
		found := false
		for _, bh := range o.UpdatedUsingServerMode {
			if bh.Path == pkgDir {
				found = true
				break
			}
		}
		if found {
			continue
		}

		err2 := filepath.Walk(filepath.Join(o.CurrentPath, pkgDir), func(pathPkg string, infoPkg os.FileInfo, err2 error) error {
			if err2 != nil {
				return err2
			}
			if infoPkg != nil && !infoPkg.IsDir() && strings.Contains(infoPkg.Name(), "ci.yaml") {
				filePath := filepath.Join(o.CurrentPath, pkgDir, infoPkg.Name())
				if _, err := os.Stat(filePath); err == nil && !os.IsNotExist(err) {
					infoPkg, err := ReadFile(filePath)
					if err != nil {
						log.Infof("unable to read infoPkg: %s", err)
						return nil
					}

					var upgradeGraph UpdateGraph
					if err := yaml.Unmarshal(infoPkg, &upgradeGraph); err == nil {
						if upgradeGraph.UpdateGraph == "semver-mode" {
							o.UpdatedUsingServerMode = append(o.UpdatedUsingServerMode, BundleDir{Path: pkgDir})
							return nil
						}
					}
				}
				return nil
			}
			return nil
		})
		if err2 != nil {
			log.Fatalf("checking server mode", err2)
		}

	}
}

func (o *Output) checkBundlesWhichRequiresUpdateAfterChangesMade() {
	var stillRequiring []BundleDir
	for _, r := range o.RequiresUpdate {
		found := false
		for _, u := range o.BundleUpdated {
			if r.Path == u.Path {
				found = true
				break
			}
		}

		if !found {
			for _, u := range o.BundlesConfiguredAccordingly {
				if r.Path == u.Path {
					found = true
					break
				}
			}
		}

		if !found {
			for _, u := range o.RequiresUpdatePkgFormat {
				if r.Path == u.Path {
					found = true
					break
				}
			}
		}

		if !found {
			stillRequiring = append(stillRequiring, r)
		}
	}
	// Replace
	o.RequiresUpdate = stillRequiring
}

func (o *Output) updateAllBundles() {
	pathToWalk := filepath.Join(o.CurrentPath, "community-operators-prod/operators")

	err := filepath.Walk(pathToWalk, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		_, errVer := semver.ParseTolerant(info.Name())
		if pathToWalk != path && info != nil && info.IsDir() && info.Name() != "manifest" && info.Name() != "tests" && info.Name() != "scorecard" && errVer == nil {

			// Load the bundle
			bundle, err := apimanifests.GetBundleFromDir(path)
			if err != nil {
				o.Errors = append(o.Errors,
					BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, ""), Error: fmt.Sprintf("Unable to load the bundle: %s", err.Error())})
				return nil
			}

			results := runOperatorHubValidator(bundle)
			found := false
			for _, v := range results {
				for _, e := range v.Errors {
					if strings.Contains(e.Detail, "this bundle is using APIs which were deprecated and removed in v1.22.") {
						o.RequiresUpdate = append(o.RequiresUpdate,
							BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "")})
						found = true
						break
					}
				}
			}

			if !found {
				o.BundlesMigrated = append(o.BundlesMigrated,
					BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "")})
				return nil
			}

			// Check if has metadata/annotations.yaml
			annotationsPath := filepath.Join(path, "metadata", "annotations.yaml")
			if _, err := os.Stat(annotationsPath); err == nil && !os.IsNotExist(err) {
				annFile, err := ReadFile(annotationsPath)
				if err != nil {
					msg := fmt.Errorf("unable to read annotations.yaml: %s", err)
					o.Errors = append(o.Errors,
						BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, msg.Error())})
					return nil
				}

				var bundleAnnotations BundleAnnotations
				if err := yaml.Unmarshal(annFile, &bundleAnnotations); err != nil {
					msg := fmt.Errorf("unable to Unmarshal annotations.yaml from path: %s", err)
					o.Errors = append(o.Errors,
						BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, msg.Error())})
					return nil
				}

				if len(bundleAnnotations.Annotations) > 0 {
					ocpLabel := bundleAnnotations.Annotations[ocpLabelAnnotation]

					// If does not have the OCP label then add it
					if len(ocpLabel) == 0 {
						if err := InsertCode(annotationsPath, fragment); err != nil {
							o.Errors = append(o.Errors,
								BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, err.Error())})
							return nil
						}
						o.BundleUpdated = append(o.BundleUpdated,
							BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "")})
					} else {
						// Check if the range is < 4.9
						rangeContainsVersion(ocpLabel, ocpVerV1beta1Unsupported)
						isPartOfTarget, err := rangeContainsVersion(ocpLabel, ocpVerV1beta1Unsupported)
						if err != nil {
							o.Errors = append(o.Errors,
								BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, err.Error())})
							return nil
						}

						if !isPartOfTarget {
							o.BundlesConfiguredAccordingly = append(o.BundlesConfiguredAccordingly,
								BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "")})
						} else {
							value, err := GetOcpLabel(ocpLabel)
							if err != nil {
								o.Errors = append(o.Errors,
									BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, err.Error())})
								return nil
							}
							err = ReplaceInFile(annotationsPath, ocpLabel, value)
							if err == nil {
								o.BundleUpdated = append(o.BundleUpdated,
									BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "")})
							} else {
								o.Errors = append(o.Errors,
									BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, err.Error())})
							}
						}
						return nil
					}
				}
			} else if os.IsNotExist(err) {

				// Check if has the a ripsaw.package.yaml
				// Go to the previous dir
				previous := filepath.Clean(filepath.Join(path, ".."))
				isPkgFormat := false
				err2 := filepath.Walk(previous, func(pathPkg string, infoPkg os.FileInfo, err2 error) error {
					if err2 != nil {
						return err2
					}
					if infoPkg != nil && !infoPkg.IsDir() && strings.HasSuffix(infoPkg.Name(), "package.yaml") {
						isPkgFormat = true
						return nil
					}
					return nil
				})
				if err2 != nil {
					log.Fatal(err2)
				}

				if isPkgFormat {
					o.RequiresUpdatePkgFormat = append(o.RequiresUpdatePkgFormat,
						BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "")})
				} else {
					err = writeAnnotations(path)
					if err == nil {
						o.BundleUpdated = append(o.BundleUpdated,
							BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "")})
					} else {
						o.Errors = append(o.Errors,
							BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "to write the annotations: "+err.Error())})
					}
				}
			} else {
				o.Errors = append(o.Errors,
					BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, "unexpected: "+err.Error())})
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (o *Output) checkIfWeCanLoadTheBundles() {
	pathToWalk := filepath.Join(o.CurrentPath, "operators")

	err := filepath.Walk(pathToWalk, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		_, errVer := semver.ParseTolerant(info.Name())
		if pathToWalk != path && info != nil && info.IsDir() && info.Name() != "manifest" &&
			info.Name() != "metadata" && info.Name() != "tests" && info.Name() != "scorecard" && errVer == nil {
			// Load the bundle
			_, err := apimanifests.GetBundleFromDir(path)
			if err != nil {
				o.BundleUnload = append(o.BundleUnload,
					BundleDir{Path: strings.ReplaceAll(path, o.CurrentPath, ""), Error: fmt.Sprintf("Unable to load the bundle after the changes: %s", err.Error())})
				return nil
			}
			return nil
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (o *Output) sort() {
	sort.Slice(o.BundleUpdated[:], func(i, j int) bool {
		return o.BundleUpdated[i].Path < o.BundleUpdated[j].Path
	})

	sort.Slice(o.RequiresUpdate[:], func(i, j int) bool {
		return o.RequiresUpdate[i].Path < o.RequiresUpdate[j].Path
	})

	sort.Slice(o.BundlesMigrated[:], func(i, j int) bool {
		return o.BundlesMigrated[i].Path < o.BundlesMigrated[j].Path
	})

	sort.Slice(o.RequiresUpdatePkgFormat[:], func(i, j int) bool {
		return o.RequiresUpdatePkgFormat[i].Path < o.RequiresUpdatePkgFormat[j].Path
	})

	sort.Slice(o.BundleUnload[:], func(i, j int) bool {
		return o.BundleUnload[i].Path < o.BundleUnload[j].Path
	})

	sort.Slice(o.BundlesConfiguredAccordingly[:], func(i, j int) bool {
		return o.BundlesConfiguredAccordingly[i].Path < o.BundlesConfiguredAccordingly[j].Path
	})

	sort.Slice(o.PkgManifestUsingServerMode[:], func(i, j int) bool {
		return o.PkgManifestUsingServerMode[i].Path < o.PkgManifestUsingServerMode[j].Path
	})

	sort.Slice(o.UpdatedUsingServerMode[:], func(i, j int) bool {
		return o.UpdatedUsingServerMode[i].Path < o.UpdatedUsingServerMode[j].Path
	})

	sort.Slice(o.Errors[:], func(i, j int) bool {
		return o.Errors[i].Path < o.Errors[j].Path
	})
}

func (o *Output) write() {
	f, err := os.Create(filepath.Join(o.CurrentPath, "community-operators-prod/scripts/add-ocp-labels/output.yml"))
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	t := template.Must(template.ParseFiles(filepath.Join(o.CurrentPath, "community-operators-prod/scripts/add-ocp-labels/template.go.tmpl")))
	err = t.Execute(f, o)
	if err != nil {
		panic(err)
	}
}

func runOperatorHubValidator(bundle *apimanifests.Bundle) []errors.ManifestResult {

	validators := apivalidation.DefaultBundleValidators
	validators = validators.WithValidators(apivalidation.OperatorHubValidator)

	objs := bundle.ObjectsToValidate()
	m := make(map[string]string)
	m["k8s-version"] = "1.22"
	objs = append(objs, m)

	results := validators.Validate(objs...)
	nonEmptyResults := []errors.ManifestResult{}

	for _, result := range results {
		if result.HasError() || result.HasWarn() {
			nonEmptyResults = append(nonEmptyResults, result)
		}
	}

	return nonEmptyResults
}

// InsertCode searches target content in the file and insert `toInsert` after the target.
func InsertCode(filename, code string) error {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	out := string(contents) + code
	// nolint:gosec
	return ioutil.WriteFile(filename, []byte(out), 0644)
}

func GetOcpLabel(ocpLabel string) (string, error) {
	if strings.Contains(ocpLabel, "=v4.9") || strings.Contains(ocpLabel, "=4.9") {
		return "=v4.8" + stamp, nil
	}

	if strings.HasPrefix(ocpLabel, "v4.9") || strings.HasPrefix(ocpLabel, "4.9") {
		return "=v4.8" + stamp, nil
	}

	if strings.HasPrefix(ocpLabel, "v4.8") || strings.HasPrefix(ocpLabel, "4.8") {
		return "=v4.8" + stamp, nil
	}

	// That is no longer accepted. It is only to check the old register
	if strings.Contains(ocpLabel, ",") {
		versionsSet := strings.Split(ocpLabel, ",")
		return versionsSet[0] + "-v4.8" + stamp, nil
	}

	// if not has not the = then the value needs contains - value less < 4.9
	if strings.Contains(ocpLabel, "-") {
		versionsSet := strings.Split(ocpLabel, "-")
		return versionsSet[0] + "-v4.8" + stamp, nil
	}
	return ocpLabel + "-v4.8\"" + stamp, nil
}

func ReplaceInFile(path, old, new string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	if !strings.Contains(string(b), old) {
		return goerror.New("unable to find the content to be replaced")
	}
	s := strings.Replace(string(b), old, new, -1)
	err = ioutil.WriteFile(path, []byte(s), info.Mode())
	if err != nil {
		return err
	}
	return nil
}

func writeAnnotations(path string) error {

	if err := os.MkdirAll(filepath.Join(path, "metadata"), os.ModePerm); err != nil {
		return err
	}
	file := filepath.Join(path, "metadata", "annotations.yaml")
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(annotationsFragment)
	if err != nil {
		return err
	}

	return nil
}

// ReadFile will return the bites of file
func ReadFile(file string) ([]byte, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return []byte{}, err
	}
	defer jsonFile.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(jsonFile)
	if err != nil {
		return []byte{}, err
	}
	return byteValue, err
}

const stamp = `
  # Updating com.redhat.openshift.versions programmatically.
  # This annotation was update because was possible to identify that this distribution uses API(s),
  # which were removed in the k8s version 1.22 and OpenShift version OCP 4.9. Then, this annotation will prevent the bundle from being distributed on 4.9.
  # Please, ensure that your project is no longer using these API(s) and that you start to
  # distribute solutions which is compatible with OpenShift 4.9.
  # For further information, see: https://github.com/redhat-openshift-ecosystem/community-operators-prod/discussions/138`

const annotationsFragment = `annotations:
  # Adding annotations.yml file with the annotation com.redhat.openshift.versions programmatically
  # This annotation was added because was possible to identify that this distribution uses API(s),
  # which will be removed in the k8s version 1.22 and OpenShift version OCP 4.9. Then, it will prevent the bundle from being distributed on to 4.9.
  # Please, ensure that your project is no longer using these API(s) and that you start to
  # distribute solutions which is compatible with OpenShift 4.9.
  # For further information, see: https://github.com/redhat-openshift-ecosystem/community-operators-prod/discussions/138
  com.redhat.openshift.versions: "v4.6-v4.8"
`

const fragment = `
  # Adding com.redhat.openshift.versions programmatically
  # This annotation was added because was possible to identify that this distribution uses API(s),
  # which were removed in the k8s version 1.22 and OpenShift version OCP 4.9. Then, it will prevent the bundle from being distributed on 4.9.
  # Please, ensure that your project is no longer using these API(s) and that you start to
  # distribute solutions which is compatible with OpenShift 4.9.
  # For further information, see: https://github.com/redhat-openshift-ecosystem/community-operators-prod/discussions/138
  com.redhat.openshift.versions: "v4.6-v4.8"
`

// rangeContainsVersion expected the range and the targetVersion version and returns true
// when the targetVersion version contains in the range
func rangeContainsVersion(r string, v string) (bool, error) {
	if len(r) == 0 {
		return false, goerror.New("range is empty")
	}
	if len(v) == 0 {
		return false, goerror.New("version is empty")
	}

	v = strings.TrimPrefix(v, "v")
	compV, err := semver.Parse(v + ".0")
	if err != nil {
		return false, fmt.Errorf("invalid version %q: %v", v, err)
	}

	// special legacy cases
	if r == "v4.5,v4.6" || r == "v4.6,v4.5" {
		semverRange := semver.MustParseRange(">=4.5.0")
		return semverRange(compV), nil
	}

	var semverRange semver.Range
	rs := strings.SplitN(r, "-", 2)
	switch len(rs) {
	case 1:
		// Range specify exact version
		if strings.HasPrefix(r, "=") {
			trimmed := strings.TrimPrefix(r, "=v")
			semverRange, err = semver.ParseRange(fmt.Sprintf("%s.0", trimmed))
		} else {
			trimmed := strings.TrimPrefix(r, "v")
			// Range specifies minimum version
			semverRange, err = semver.ParseRange(fmt.Sprintf(">=%s.0", trimmed))
		}
		if err != nil {
			return false, fmt.Errorf("invalid range %q: %v", r, err)
		}
	case 2:
		min := rs[0]
		max := rs[1]
		if strings.HasPrefix(min, "=") || strings.HasPrefix(max, "=") {
			return false, fmt.Errorf("invalid range %q: cannot use equal prefix with range", r)
		}
		min = strings.TrimPrefix(min, "v")
		max = strings.TrimPrefix(max, "v")
		semverRangeStr := fmt.Sprintf(">=%s.0 <=%s.0", min, max)
		semverRange, err = semver.ParseRange(semverRangeStr)
		if err != nil {
			return false, fmt.Errorf("invalid range %q: %v", r, err)
		}
	default:
		return false, fmt.Errorf("invalid range %q", r)
	}
	return semverRange(compV), nil
}
