let SessionLoad = 1
let s:so_save = &g:so | let s:siso_save = &g:siso | setg so=0 siso=0 | setl so=-1 siso=-1
let v:this_session=expand("<sfile>:p")
silent only
silent tabonly
cd ~/playground/go/nicest-backend-framework/pokedex
if expand('%') == '' && !&modified && line('$') <= 1 && getline(1) == ''
  let s:wipebuf = bufnr('%')
endif
let s:shortmess_save = &shortmess
if &shortmess =~ 'A'
  set shortmess=aoOA
else
  set shortmess=aoO
endif
badd +763 ~/.config/nvim/init.lua
badd +45 pokemon.go
badd +97 pokemon_test.go
badd +49 /opt/homebrew/Cellar/go/1.24.5/libexec/src/encoding/json/stream.go
badd +1115 /opt/homebrew/Cellar/go/1.24.5/libexec/src/testing/testing.go
badd +4 data.json
badd +309 ~/go/pkg/mod/github.com/google/go-cmp@v0.7.0/cmp/compare.go
argglobal
%argdel
edit pokemon.go
let s:save_splitbelow = &splitbelow
let s:save_splitright = &splitright
set splitbelow splitright
wincmd _ | wincmd |
vsplit
1wincmd h
wincmd w
let &splitbelow = s:save_splitbelow
let &splitright = s:save_splitright
wincmd t
let s:save_winminheight = &winminheight
let s:save_winminwidth = &winminwidth
set winminheight=0
set winheight=1
set winminwidth=0
set winwidth=1
exe 'vert 1resize ' . ((&columns * 69 + 84) / 168)
exe 'vert 2resize ' . ((&columns * 98 + 84) / 168)
tcd ~/playground/go/nicest-backend-framework/pokedex
argglobal
balt ~/playground/go/nicest-backend-framework/pokedex/data.json
setlocal fdm=expr
setlocal fde=v:lua.vim.treesitter.foldexpr()
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=99
setlocal fml=1
setlocal fdn=4
setlocal fen
62
normal! zo
let s:l = 54 - ((25 * winheight(0) + 19) / 38)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 54
normal! 0
wincmd w
argglobal
if bufexists(fnamemodify("~/playground/go/nicest-backend-framework/pokedex/pokemon_test.go", ":p")) | buffer ~/playground/go/nicest-backend-framework/pokedex/pokemon_test.go | else | edit ~/playground/go/nicest-backend-framework/pokedex/pokemon_test.go | endif
if &buftype ==# 'terminal'
  silent file ~/playground/go/nicest-backend-framework/pokedex/pokemon_test.go
endif
balt ~/go/pkg/mod/github.com/google/go-cmp@v0.7.0/cmp/compare.go
setlocal fdm=expr
setlocal fde=v:lua.vim.treesitter.foldexpr()
setlocal fmr={{{,}}}
setlocal fdi=#
setlocal fdl=99
setlocal fml=1
setlocal fdn=4
setlocal fen
11
normal! zo
110
normal! zo
121
normal! zo
let s:l = 97 - ((12 * winheight(0) + 19) / 38)
if s:l < 1 | let s:l = 1 | endif
keepjumps exe s:l
normal! zt
keepjumps 97
normal! 0
wincmd w
2wincmd w
exe 'vert 1resize ' . ((&columns * 69 + 84) / 168)
exe 'vert 2resize ' . ((&columns * 98 + 84) / 168)
tabnext 1
if exists('s:wipebuf') && len(win_findbuf(s:wipebuf)) == 0 && getbufvar(s:wipebuf, '&buftype') isnot# 'terminal'
  silent exe 'bwipe ' . s:wipebuf
endif
unlet! s:wipebuf
set winheight=1 winwidth=20
let &shortmess = s:shortmess_save
let &winminheight = s:save_winminheight
let &winminwidth = s:save_winminwidth
let s:sx = expand("<sfile>:p:r")."x.vim"
if filereadable(s:sx)
  exe "source " . fnameescape(s:sx)
endif
let &g:so = s:so_save | let &g:siso = s:siso_save
set hlsearch
nohlsearch
doautoall SessionLoadPost
unlet SessionLoad
" vim: set ft=vim :
