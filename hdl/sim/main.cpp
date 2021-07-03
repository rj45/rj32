// (C)2021 rj45
// Mostly copied from:
// Project F: FPGA Graphics
// (C)2021 Will Green, open source software released under the MIT License
// Learn more at https://projectf.io

#include <stdio.h>
#include <SDL2/SDL.h>
#include <verilated.h>
#include "Vtop.h"

// screen dimensions
const int H_RES = 640;
const int V_RES = 399;

typedef struct Pixel
{            // for SDL texture
  uint8_t a; // transparency
  uint8_t b; // blue
  uint8_t g; // green
  uint8_t r; // red
} Pixel;

int main(int argc, char *argv[])
{
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
  Vtop *top = new Vtop;

  top->clk_cpu = 0;
  top->clk_vga = 0;
  top->eval();

  // top->rst = 1;
  // top->clk_pix = 0;
  // top->eval();
  // top->rst = 0;
  // top->eval();

  int vga_x = 0;
  int vga_y = 0;
  bool done = false;
  bool freerun = false;
  bool runstep = false;

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

    // vga clock double the cpu clock
    top->clk_vga = 0;
    top->eval();
    top->clk_vga = 1;
    top->eval();

    if (top->error) {
      printf("\nERROR!!!\n");
      done = true;
    }

    if (top->halt) {
      printf("\nSuccess!!\n");
      done = true;
    }

    // update pixel if not in blanking interval
    if (top->vga_de)
    {
      vga_x++;
      if (vga_x >= H_RES) {
        vga_y++;
        vga_x = 0;
        if (vga_y >= V_RES) {
          vga_y = 0;
        }
      }
      Pixel *p = &screenbuffer[vga_y * H_RES + vga_x];
      p->a = 0xFF; // transparency
      p->b = top->b << 4;
      p->g = top->g << 4;
      p->r = top->r << 4;
    }

    // update texture once per frame at start of blanking
    bool updated = false;
    if (vga_x >= (H_RES-1) && vga_y >= (V_RES-1))
    {
      if (!updated) {
        updated = true;
        if (freerun) {
          runstep = true;
        }
        SDL_UpdateTexture(sdl_texture, NULL, screenbuffer, H_RES * sizeof(Pixel));
        SDL_RenderClear(sdl_renderer);
        SDL_RenderCopy(sdl_renderer, sdl_texture, NULL, NULL);
        SDL_RenderPresent(sdl_renderer);
      }
    } else {
      updated = false;
    }
  }

  top->final(); // simulation done

  SDL_DestroyTexture(sdl_texture);
  SDL_DestroyRenderer(sdl_renderer);
  SDL_DestroyWindow(sdl_window);
  SDL_Quit();
  return 0;
}
