let s:w2v_model = get(g:, 'w2v_model', 'data.model')

function! s:w2v_search(expr) abort
  let expr = a:expr == '' ? input('Expr: ') : a:expr
  if expr == ''
    return
  endif
  let cmd = printf('w2v-repl -f %s -q %s', shellescape(s:w2v_model), shellescape(expr))
  let result = system('w2v-repl -q ' . shellescape(expr))
  if v:shell_error
    redraw
    echohl Error | echom substitute(iconv(result, 'char', &encoding), '\n', '', 'g') | echohl None
    return
  endif
  let @/ = join(split(result, "\n"), '\|')
  call feedkeys('n', 'n')
endfunction

command! -nargs=* W2V call s:w2v_search(<q-args>)
