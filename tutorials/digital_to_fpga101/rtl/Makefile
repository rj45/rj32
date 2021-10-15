# Mostly copied from https://projectf.io/posts/building-ice40-fpga-toolchain/

FPGA_PKG = sg48
FPGA_TYPE = up5k
# PCF = icebreaker.pcf
PCF = icesugar.pcf

NEXTPNR_OPTS = --pcf-allow-unconstrained --opt-timing
IGN_WARNS = -Wno-WIDTH -Wno-DECLFILENAME -Wno-PINMISSING
# TRACE = --trace-fst

ADD_SRC = uart/uart.v uart/uart_rx.v uart/uart_tx.v uart/uart_baud_tick_gen.v bram.v
SIM_PLL = pll_sim.v

.PHONY: all clean lint sim upload

all: lint fpga101

fpga101: fpga101.rpt fpga101.bin

lint:
	verilator --lint-only -Wall ${IGN_WARNS} --top-module fpga101 fpga101.v $(ADD_SRC) ${SIM_PLL}

sim: obj_dir/Vfpga101
	(cd obj_dir && ./Vfpga101) || true

obj_dir/V%: %.v $(ADD_SRC) $(SIM_PLL) main.cpp telnet_uart.hpp
	verilator ${IGN_WARNS} ${TRACE} -I../ --top-module $(subst .v,,$<) -cc $< $(ADD_SRC) ${SIM_PLL} --exe main.cpp
	make -C ./obj_dir OPT=-O3 -f Vfpga101.mk

%.json: %.v $(ADD_SRC) $(FPGA_PLL)
	yosys -ql $(subst .json,,$@)-yosys.log -p 'synth_ice40 -top $(subst .json,,$@) -json $@' $< $(ADD_SRC) ${FPGA_PLL}

%.asc: %.json
	nextpnr-ice40 -ql $(subst .asc,,$@)-nextpnr.log ${NEXTPNR_OPTS} --${FPGA_TYPE} --package ${FPGA_PKG} --json $< --pcf ${PCF} --asc $@

%.rpt: %.asc
	icetime -d ${FPGA_TYPE} -P ${FPGA_PKG} -p ${PCF} -mtr $@ $<

%.bin: %.asc
	icepack $< $@

upload: fpga101.bin
	icesprog -w fpga101.bin

pll.v:
	icepll -q -m -f pll.v -i 12 -o 29.4912
	sed -i 's/SB_PLL40_CORE/SB_PLL40_PAD/g' pll.v
	sed -i 's/REFERENCECLK/PACKAGEPIN/g' pll.v

clean:
	rm -rf fpga101*.json fpga101*.asc fpga101*.rpt fpga101*.bin fpga101*yosys.log fpga101*nextpnr.log pll.v obj_dir
