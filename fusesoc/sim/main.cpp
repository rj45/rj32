// (C)2021 rj45
// Mostly copied from:
// Project F: FPGA Graphics
// (C)2021 Will Green, open source software released under the MIT License
// Learn more at https://projectf.io

#include <stdio.h>
#include <memory>
#include <SDL2/SDL.h>
#include <verilated.h>
#include "Vtop_verilator.h"

#if VM_TRACE
#include <verilated_fst_c.h>
#endif

// screen dimensions
const int H_RES = 640;
const int V_RES = 400;

typedef struct Pixel
{            // for SDL texture
  uint8_t a; // transparency
  uint8_t b; // blue
  uint8_t g; // green
  uint8_t r; // red
} Pixel;

int main(int argc, char *argv[])
{
  #if VM_TRACE
  const std::unique_ptr<VerilatedContext> contextp{new VerilatedContext};
  #endif

  Verilated::commandArgs(argc, argv);

  if (SDL_Init(SDL_INIT_VIDEO) < 0)
  {
    printf("SDL init failed.\n");
    return 1;
  }

  Pixel screenbuffer[H_RES * V_RES];

  SDL_Window *sdl_window = NULL;
  SDL_Renderer *sdl_renderer = NULL;
  SDL_Texture *sdl_texture = NULL;

  sdl_window = SDL_CreateWindow("rj32", SDL_WINDOWPOS_CENTERED,
                                SDL_WINDOWPOS_CENTERED, H_RES, V_RES, SDL_WINDOW_SHOWN);
  if (!sdl_window)
  {
    printf("Window creation failed: %s\n", SDL_GetError());
    return 1;
  }

  sdl_renderer = SDL_CreateRenderer(sdl_window, -1, SDL_RENDERER_ACCELERATED);
  if (!sdl_renderer)
  {
    printf("Renderer creation failed: %s\n", SDL_GetError());
    return 1;
  }

  sdl_texture = SDL_CreateTexture(sdl_renderer, SDL_PIXELFORMAT_RGBA8888,
                                  SDL_TEXTUREACCESS_TARGET, H_RES, V_RES);
  if (!sdl_texture)
  {
    printf("Texture creation failed: %s\n", SDL_GetError());
    return 1;
  }

  // initialize Verilog module
  Vtop_verilator *top = new Vtop_verilator;

#if VM_TRACE
  Verilated::traceEverOn(true);
  VerilatedFstC* tfp = new VerilatedFstC;
  top->trace(tfp, 4);
  tfp->open("run.fst");
#endif

  top->clk_cpu = 0;
  top->clk_vga = 0;
  top->eval();

#if VM_TRACE
  tfp->dump(contextp->time());
#endif

  // top->rst = 1;
  // top->clk_pix = 0;
  // top->eval();
  // top->rst = 0;
  // top->eval();

  int vga_x = 0;
  int vga_y = 0;
  int cycle = 0;
  bool done = false;
  bool freerun = false;
  bool runstep = false;
  bool vsynced = false;
  bool hsynced = false;
  bool updated = false;
  int blank = 0;
  int updatecount = 16;
  int pvs = top->vga_vs;
  int phs = top->vga_hs;

  if (argc > 1 && strcmp(argv[1], "run") == 0) {
    freerun = true;
  }

  while (!done)
  {
    top->step = 0;

    if (runstep) {
      top->step = 1;
      runstep = false;
    }

    // check for quit event
    SDL_Event e;
    if (SDL_PollEvent(&e))
    {
      switch(e.type) {
        case SDL_QUIT:
          done = true;
          break;
        case SDL_KEYDOWN:
          switch (e.key.keysym.scancode) {
            case SDL_SCANCODE_SPACE:
              top->step = 1;
              freerun = false;
              break;
            case SDL_SCANCODE_R:
              freerun = true;
              break;
          }
          break;
      }
    }

    // cpu on a half clock
    if (top->clk_cpu == 0)
    {
      top->clk_cpu = 1;
    }
    else
    {
      top->clk_cpu = 0;
    }

#if VM_TRACE
    contextp->timeInc(1);
#endif

    // vga clock double the cpu clock
    top->clk_vga = 0;
    top->eval();

#if VM_TRACE
    tfp->dump(contextp->time());
#endif

#if VM_TRACE
    contextp->timeInc(1);
#endif

    top->clk_vga = 1;
    top->eval();
    cycle++;

#if VM_TRACE
    tfp->dump(contextp->time());
#endif

    if (top->error) {
      printf("\nERROR!!!\n");
      done = true;
    }

    if (top->halt) {
      printf("\nSuccess!!\n");
      done = true;
    }

    if (!top->vga_de) {
      blank++;
      if(top->vga_vs != pvs) {
        pvs = top->vga_vs;
        vsynced = true;
        if (vga_y != 399) {
          printf("Error: expecting y to be %d, was %d\n", 399, vga_y);
        }
      } else if(top->vga_hs != phs) {
        phs = top->vga_hs;
        hsynced = true;

        if (vga_x != 639) {
          printf("Error: expecting x to be %d, was %d\n", 639, vga_x);
        }
      }
    }

    // update pixel if not in blanking interval
    if (top->vga_de)
    {
      if (vsynced) {
        vsynced = hsynced = false;
        if (blank != (800*(12+2+35))+160) {
          printf("Error: blank was %d but should have been %d!\n", blank, (800*(12+2+35))+160);
        }
        blank = 0;
        vga_y = vga_x = 0;
      } else if(hsynced) {
        hsynced = false;
        if (blank != 160) {
          printf("Error: blank was %d but should have been 160!\n", blank);
        }
        blank = 0;
        vga_x = 0;
        vga_y++;
      } else {
        vga_x++;
      }
      // if (vga_x >= H_RES) {
      //   vga_y++;
      //   vga_x = 0;
      //   if (vga_y >= V_RES) {
      //     vga_y = 0;
      //   }
      // }

      if (vga_x < H_RES && vga_y < V_RES) {
        Pixel *p = &screenbuffer[(vga_y * H_RES) + vga_x];
        p->a = 0xFF; // transparency
        p->b = top->b;
        p->g = top->g;
        p->r = top->r;
      } else {
        printf("out of bounds: %d %d\n", vga_x, vga_y);
      }
    }

    // update texture once in a while
    if (vga_x >= (H_RES-1))
    {
      if (!updated) {
        updated = true;
        if (freerun) {
          runstep = true;
        }
        updatecount--;
        if (updatecount <= 0) {
          updatecount = 16;
          SDL_UpdateTexture(sdl_texture, NULL, screenbuffer, H_RES * sizeof(Pixel));
          SDL_RenderClear(sdl_renderer);
          SDL_RenderCopy(sdl_renderer, sdl_texture, NULL, NULL);
          SDL_RenderPresent(sdl_renderer);
        }
      }
    } else {
      updated = false;
    }
  }

#if VM_TRACE
  tfp->close();
#endif

  top->final(); // simulation done

  SDL_DestroyTexture(sdl_texture);
  SDL_DestroyRenderer(sdl_renderer);
  SDL_DestroyWindow(sdl_window);
  SDL_Quit();
  return 0;
}
