; Comentário
; x = a + b

.CODE
LDA VAR1   ; a
ADD VAR2   ; b
STA VAR3   ; x
HLT

.DATA
VAR1 DB 50
VAR2 DB 5
VAR3 DB 0