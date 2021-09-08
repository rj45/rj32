#!/bin/env ruby

# Digital emits verilog with all the built-in modules having the same
# names, which makes it difficult to have two different exports be
# included in the same run of yosys or verilator.
#
# This script just runs through and renames all non-top modules with
# a prefix to ensure they don't conflict with the names from another
# file.

verilog = File.read(ARGV[0])
top = ARGV[1]
prefix = ARGV[2]

verilog = verilog.split("\n")

translations = {}

verilog = verilog.map do |line|
  moddef = line.match(/^(\s*)module (\S+)(.*)$/)
  modref = line.match(/^(\s*)(\S+)( #\(| \S+ \()$/)
  if !moddef.nil?
    pre = moddef[1]
    name = moddef[2]
    post = moddef[3]

    if name == top
      line
    else
      newname = prefix + name

      translations[name] = newname

      "#{pre}module #{newname}#{post}"
    end
  elsif !modref.nil?
    pre = modref[1]
    name = modref[2]
    post = modref[3]

    if !translations[name].nil?
      "#{pre}#{translations[name]}#{post}"
    else
      line
    end
  else
    line
  end
end

print verilog.join("\n")