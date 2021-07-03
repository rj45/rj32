IGNWARNINGS=-Wno-PINMISSING -Wno-DECLFILENAME -Wno-MODDUP -Wno-UNUSED -Wno-WIDTH

build:
	verilator $(IGNWARNINGS) -I../ --top-module top -cc top.v rj32.v vga_blinkenlights.v mem.v --exe main.cpp -LDFLAGS "`sdl2-config --libs`"
	cp ../*.hex ./obj_dir/
	make -C ./obj_dir OPT=-O3 -f Vtop.mk
	cd obj_dir && ./Vtop $(ARGS)

clean:
	rm -rf obj_dir
