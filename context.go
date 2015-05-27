package plumb

type Field struct {
  Name string
  Description string `yaml:",omitempty"`
  Type string
}

type PlumbContext struct {
  Language string
  Name string
  Inputs []Field `yaml:",flow"`
  Outputs []Field `yaml:",flow"`
  Env []string `yaml:",flow,omitempty"`
  Install []string `yaml:",flow,omitempty"`
}
