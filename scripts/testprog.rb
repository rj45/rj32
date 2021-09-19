#!/bin/env ruby

require "pathname"

# read the whole file specified in the first argument
prog = File.read(ARGV[0])

# split it into lines
prog = prog.split("\n")

# comment out each line with a `#`
prog = prog.map {|l| "# #{l}"}

# join the lines back together
prog = prog.join("\n")

# figure out a path that customasm will be happy with
absfilename = File.absolute_path(ARGV[0])
path = File.dirname(__dir__) # root of scripts folder
relfilename = Pathname.new(absfilename).
  relative_path_from(Pathname.new(path))

# run the program through customasm as well
hex = `cd #{path} && customasm -f logisim16 -p -q #{relfilename}`

# split into lines, take all but the first (header) line
_, *lines = hex.split("\n")

# split all the lines on spaces
machcode = lines.map {|x| x.split}

# flatten an array of arrays to a single large array
machcode = machcode.flatten

# add `0x` prefix to each hex number
machcode = machcode.map {|x| "0x#{x}"}

# print the template for digital test cases
puts <<~DONE
  clock run error halt

  #{prog}

  program(#{machcode.join(", ")})

  let i = 0;
  while(!(halt | error | (i >= 100)))
    let i = i + 1;
    0 1 0 0
    1 1 x x
  end while
  0 1 0 1
DONE
