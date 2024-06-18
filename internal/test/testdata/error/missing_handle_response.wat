(module $missing_handle_response
  (memory (export "memory") 1 1 (; 1 page==64KB ;))

  (global $name i32 (i32.const 0))
  (data (i32.const 0) "Content-Type")
  (global $name_len i32 (i32.const 12))

  (global $value i32 (i32.const 16))
  (data (i32.const 16) "text/plain")
  (global $value_len i32 (i32.const 10))

  ;; handle_request returns non-zero to proceed to the next handler.
  (func $handle_request (export "handle_request") (result (; ctx_next ;) i64)
    ;; call the next handler
    (return (i64.const 1)))

)
