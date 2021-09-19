#!/bin/env ruby

projpath = File.absolute_path(File.join(__dir__, ".."))

testsuite = File.read("#{projpath}/dig/testsuite.dig")
testsglob = "#{projpath}/programs/tests/*.asm"

Dir[testsglob].each do |filename|
  testname = File.basename(filename, ".asm")
  # puts testname

  testprog = "#{projpath}/scripts/testprog.rb"
  replacementtext = `ruby #{testprog} #{filename}`


  i = testsuite.index("<string>#{testname}</string>")
  if i.nil?
    warn "Could not find #{testname}!!!"
    print testsuite
    exit 1
  end

  ds = "<dataString>"

  first = testsuite.index(ds, i)+ds.length

  last = testsuite.index("</dataString>", i)

  texttoreplace = testsuite[first,last-first]

  encodedreplacement = replacementtext.
    encode(:xml => :text).
    gsub("\"", "&quot;").
    gsub("\'", "&apos;")

  testsuite = testsuite.sub(texttoreplace, encodedreplacement)

end

print testsuite



