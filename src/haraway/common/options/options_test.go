package options

import "strings"
import "testing"

func TestNew (t *testing.T) {
  _, err := New(`
    usage: haraway <flags>... <command> <args>...
    --
    root=     -r,--root=,HARAWAY_ROOT     Path to the haraway data root
    prefix=   -p,--prefix,HARAWAY_PREFIX  Path to the haraway install prefix.
    verbose   -v,--verbose                Show more info
    debug     -d,--debug,HARAWAY_DEBUG    Show debug info
    --
    exec      exec                        Execute a command within the haraway sanbox
    shell     sh,shell                    Open a shell within the haraway sanbox
  `)

  if err != nil {
    t.Error(err)
  }
}

func TestParse (t *testing.T) {
  spec, err := New(`
    usage: haraway <flags>... <command> <args>...
    --
    root=     -r,--root=,HARAWAY_ROOT     Path to the haraway data root
    prefix=   -p,--prefix,HARAWAY_PREFIX  Path to the haraway install prefix.
    verbose   -v,--verbose                Show more info
    debug     -d,--debug,HARAWAY_DEBUG    Show debug info
    --
    exec      c,exec                      Execute a command within the haraway sanbox
    shell     sh,shell                    Open a shell within the haraway sanbox
  `)
  if err != nil { t.Error(err) }

  opts, err := spec.Parse([]string{"haraway", "-p", "/usr/local", "-r=hello", "-v", "c", "ls"},[]string{})

  if err != nil { t.Fatal(err) }

  t.Errorf("SPEC: %+v", spec)
  t.Errorf("OPTS: %+v", opts)

  if opts.Get("root") != "hello" {
    t.Error("--root != hello")
  }

  if opts.Get("verbose") != "true" {
    t.Error("--verbose != true")
  }

  if strings.Join(opts.Args, " ") != "exec ls" {
    t.Errorf(".Args != [`exec`, `ls`] (was: %+v)", opts.Args)
  }

}
