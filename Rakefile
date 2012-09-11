OS      = %w( darwin linux )
ARCH    = %w( 386 amd64 )
VERSION = File.read("VERSION").strip

LIBS    = []
TOOLS   = []
DEPS    = []

ENV['GOPATH'] = Dir.pwd
ENV['CGO_ENABLED'] = '0'

specs = {}
package_data = `go list -f '{{.ImportPath}}:{{.Name}}{{range .Imports}}:{{.}}{{end}}' ./...`
package_data.strip.split("\n").each do |line|
  path, name, *deps = *line.split(":")
  specs[path] = deps

  if name == "main"
    TOOLS << path
  else
    LIBS  << path
  end
end

specs.each do |path, deps|
  specs[path] = LIBS & deps
end

file "src/heroku-keepalive/version.go" => "VERSION" do
  File.open("src/heroku-keepalive/version.go", "w+", 0644) do |f|
    f.puts "package main"
    f.puts "const Version = `#{VERSION}`"
  end
end

OS.each do |os|
  ARCH.each do |arch|
    prefix = "dist/heroku-keepalive-#{VERSION}-#{os}-#{arch}"
    bindir = "#{prefix}/bin"
    directory bindir

    pkg_deps = []

    LIBS.each do |lib|
      deps  = specs[lib].map {|d| "pkg/#{os}_#{arch}/#{d}.a" }
      deps += FileList["src/#{lib}/*.go"]
      file "pkg/#{os}_#{arch}/#{lib}.a" => deps do
        ENV['GOARCH'] = arch
        ENV['GOOS']   = os
        sh "go install #{lib}"
      end
    end

    TOOLS.each do |tool|
      deps  = specs[tool].map {|d| "pkg/#{os}_#{arch}/#{d}.a" }
      deps += FileList["src/#{tool}/*.go"]
      deps += [bindir]
      bin = tool.gsub("/", "-").gsub(/-cli$/, '')

      if tool == "heroku-keepalive"
        deps << "src/heroku-keepalive/version.go"
      end

      file "#{bindir}/#{bin}" => deps do
        ENV['GOARCH'] = arch
        ENV['GOOS']   = os
        sh "go build -o #{bindir}/#{bin} src/#{tool}/*.go"
      end
      pkg_deps << "#{bindir}/#{bin}"
    end

    file "#{prefix}/README.md" => ["README.md", prefix] do
      cp "README.md", "#{prefix}/README.md"
    end
    pkg_deps << "#{prefix}/README.md"

    file "dist/heroku-keepalive-#{VERSION}-#{os}-#{arch}.tar.gz" => pkg_deps do
      Dir.chdir "dist" do
        sh "tar -czf heroku-keepalive-#{VERSION}-#{os}-#{arch}.tar.gz heroku-keepalive-#{VERSION}-#{os}-#{arch}"
      end
    end
    DEPS << "dist/heroku-keepalive-#{VERSION}-#{os}-#{arch}.tar.gz"

  end
end


desc "Remove all build targets"
task :clean do
  rm_rf 'pkg'
  rm_rf 'dist'
end


desc "Build all targets"
task :build => DEPS

desc "Build a new release"
task :dist => [:clean, :build]

task :default => :build
