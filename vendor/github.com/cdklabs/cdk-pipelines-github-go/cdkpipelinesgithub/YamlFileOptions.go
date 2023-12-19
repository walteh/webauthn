package cdkpipelinesgithub


// Options for `YamlFile`.
// Experimental.
type YamlFileOptions struct {
	// The object that will be serialized.
	//
	// You can modify the object's contents
	// before synthesis.
	// Default: {} an empty object.
	//
	// Experimental.
	Obj interface{} `field:"optional" json:"obj" yaml:"obj"`
}

