package options

import "strings"
import "fmt"
import "os"

/*

spec = `
haraway - a simple private cloud
Usage: haraway <flags> <command>
--
file=      --file=          Some file
--
prefix     HARAWAY_PREFIX   The installation prefix of haraway
--
exec       exec             Execute a command
shell      sh,shell         Run a shell
remote     remote           -
controller controller       -
*
`

*/

type Spec struct {
  usage     string

  allow_unknown_args bool

  options     map[string]string
  flags       map[string]bool
  required    map[string]bool
  environment map[string]string
  commands    map[string]string
}

type Options struct {
  options map[string]string
  Command string
  Args    []string
}

func New (desc string) (spec * Spec, err error) {
  spec = new(Spec)
  spec.options     = make(map[string]string, 0)
  spec.flags       = make(map[string]bool,   0)
  spec.required    = make(map[string]bool,   0)
  spec.commands    = make(map[string]string, 0)
  spec.environment = make(map[string]string, 0)
  spec.allow_unknown_args = false

  g_indent := -1
  indent   := -1
  section  := 0
  lines    := []string{}

  for _, line := range strings.Split(desc, "\n") {
    if g_indent == -1 {
      clean_line := strings.TrimLeft(line, " \t")
      if clean_line != "" {
        g_indent = len(line) - len(clean_line)
      }
    } else {
      line = line[g_indent:]
    }

    line := strings.TrimRight(line, " \t")

    if line == "" {
      if section != 1 && section != 2 && section != 3 {
        lines = append(lines, line)
      }
      continue
    }

    if section == 1 || section == 2 || section == 3 {
      if strings.HasPrefix(line, "#") {
        if indent == -1 {
          indent = len(line) - len(strings.TrimLeft(line[1:], " \t"))
        }

        if line == "#" {
          lines = append(lines, "")
        } else {
          line = line[indent:]
          lines = append(lines, line)
        }
        continue
      }
    }

    switch section {


    case 0: // usage
      if line == "--" {
        if len(lines) > 0 && lines[len(lines)-1] != "" {
          lines = append(lines, "")
        }
        section += 1;
        continue
      }

      lines = append(lines, line)


    case 1: // options
      if line == "--" {
        if len(lines) > 0 && lines[len(lines)-1] != "" {
          lines = append(lines, "")
        }
        section += 1;
        continue
      }


      parts := strings.SplitN(line, " ", 2)
      if len(parts) == 1 {
        err = fmt.Errorf("Invalid option spec: %s", line)
        return
      }
      if indent == -1 {
        indent = len(line) - len(strings.TrimLeft(parts[1], " \t"))
      }
      option := parts[0]
      line    = strings.Trim(parts[1], " \t")

      required := false
      flag     := true

      if strings.HasPrefix(option, "!") {
        option   = option[1:]
        required = true
      }

      if strings.HasSuffix(option, "=") {
        option = option[0:len(option)-1]
        flag   = false
      }

      spec.flags[option]    = flag
      spec.required[option] = required

      parts = strings.SplitN(line, " ", 2)
      if len(parts) == 1 { parts = append(parts, "-") }
      parts[1] = strings.Trim(parts[1], " \t")

      if parts[1] != "-" {
        lines  = append(lines, "  " + line)
      }

      parts = strings.Split(parts[0], ",")

      for _, part := range parts {
        part = strings.SplitN(part, "=", 2)[0]

        if strings.HasPrefix(part, "--") {
          spec.options[part] = option
          continue
        }

        if strings.HasPrefix(part, "-") {
          spec.options[part] = option
          continue
        }

        spec.environment[part] = option
      }


    case 2: // environment variables
      if line == "--" {
        if len(lines) > 0 && lines[len(lines)-1] != "" {
          lines = append(lines, "")
        }
        section += 1;
        continue
      }


      parts := strings.SplitN(line, " ", 2)
      if len(parts) == 1 {
        err = fmt.Errorf("Invalid env spec: %s", line)
        return
      }
      if indent == -1 {
        indent = len(line) - len(strings.TrimLeft(parts[1], " \t"))
      }
      env    := parts[0]
      line    = strings.Trim(parts[1], " \t")

      required := false
      flag     := true

      if strings.HasPrefix(env, "!") {
        env      = env[1:]
        required = true
      }

      if strings.HasSuffix(env, "=") {
        env  = env[0:len(env)-1]
        flag = false
      }

      spec.flags[env]    = flag
      spec.required[env] = required

      parts = strings.SplitN(line, " ", 2)
      if len(parts) == 1 { parts = append(parts, "-") }
      parts[1] = strings.Trim(parts[1], " \t")

      if parts[1] != "-" {
        lines  = append(lines, "  " + line)
      }

      parts = strings.Split(parts[0], ",")

      for _, part := range parts {
        part = strings.SplitN(part, "=", 2)[0]
        spec.environment[part] = env
      }


    case 3: // commands
      if line == "--" {
        if len(lines) > 0 && lines[len(lines)-1] != "" {
          lines = append(lines, "")
        }
        section += 1;
        continue
      }

      if line == "*" {
        spec.allow_unknown_args = true
        continue
      }

      parts := strings.SplitN(line, " ", 2)
      if len(parts) == 1 {
        err = fmt.Errorf("Invalid command spec: %s", line)
        return
      }
      if indent == -1 {
        indent = len(line) - len(strings.TrimLeft(parts[1], " \t"))
      }
      command := parts[0]
      line     = strings.Trim(parts[1], " \t")

      parts = strings.SplitN(line, " ", 2)
      if len(parts) == 1 { parts = append(parts, "-") }
      parts[1] = strings.Trim(parts[1], " \t")

      if parts[1] != "-" {
        lines = append(lines, "  " + line)
      }

      parts = strings.Split(parts[0], ",")
      for _, part := range parts {
        spec.commands[part] = command
      }


    case 4: // appendix
      if line == "--" {
        if len(lines) > 0 && lines[len(lines)-1] != "" {
          lines = append(lines, "")
        }
        section += 1;
        continue
      }

      lines = append(lines, line)


    }
  }

  spec.usage = strings.Join(lines, "\n") + "\n"
  spec.usage = strings.Trim(spec.usage, " \t\n")
  return
}

func (spec * Spec) Parse (args []string, environ []string) (o * Options, err error) {
  opts := new(Options)
  opts.options = make(map[string]string, 0)
  opts.Args    = []string{}

  for _, env := range environ {
    parts := strings.SplitN(env, "=", 2)
    if option, present := spec.environment[parts[0]] ; present {
      opts.options[option] = parts[1]
    }
  }

  for i := 1; i < len(args); i ++ {
    arg := args[i]

    if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
      option := "-"
      value  := "true"

      parts := strings.SplitN(arg, "=", 2)

      if len(parts) == 2 {
        option = parts[0]
      } else {
        option = arg
      }

      if opt, present := spec.options[option] ; present {
        option = opt
      } else {
        err = fmt.Errorf("Invalid option: %s was not recognized", arg)
        return
      }

      if spec.flags[option] {
        if len(parts) == 2 {
          err = fmt.Errorf("Invalid option: %s was not recognized (doesn't take a value)", arg)
          return
        }
      } else {
        if len(parts) == 2 {
          value = parts[1]
        } else if len(args) > (i + 1) {
          value = args[i+1]
          i++
        } else {
          err = fmt.Errorf("Invalid option: %s was not recognized (requires a value)", arg)
          return
        }
      }

      opts.options[option] = value
      continue
    }

    if command, present := spec.commands[arg]; present {
      opts.Command = command
      opts.Args    = args[i:]
      opts.Args[0] = opts.Command
      break
    }

    if spec.allow_unknown_args {
      opts.Args = args[i:]
      break
    }

    err = fmt.Errorf("Invalid argument: %s was not recognized", arg)
    return
  }

  for option, required := range spec.required {
    if _, present := opts.options[option]; required && !present {
      err = fmt.Errorf("Missing option: %s", option)
      return
    }
  }

  for env, option := range spec.environment {
    if value, present := opts.options[option]; present {
      os.Setenv(env, value)
    }
  }

  o = opts
  return
}

func (spec * Spec) PrintUsage () {
  fmt.Fprintf(os.Stderr, "%s\n", spec.usage)
}

func (spec * Spec) PrintUsageAndExit () {
  spec.PrintUsage()
  os.Exit(62)
}

func (spec * Spec) PrintUsageWithError (err error) {
  fmt.Fprintf(os.Stderr, "error: %s\n%s\n", err, spec.usage)
  os.Exit(62)
}

func (opts * Options) Get (option string) (string) {
  if value, present := opts.options[option] ; present {
    return value
  }
  return ""
}

func (opts * Options) GetBool (option string) (bool) {
  return opts.Get(option) == "true"
}
